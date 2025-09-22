package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/harshitrajsinha/medi-go/config"
	driver "github.com/harshitrajsinha/medi-go/internal/db"
	middleware "github.com/harshitrajsinha/medi-go/internal/middleware"
	apiRoutesV1 "github.com/harshitrajsinha/medi-go/internal/routes/api/v1"
	"github.com/harshitrajsinha/medi-go/internal/store"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
)

var db *driver.DBClient
var rdb *redis.Client

const dbDriver string = "postgres"

type RedisConfig struct {
	Host string
	Port string
	Pass string
}

func init() {

	var err error
	var dbConnStr string
	var redisConfig RedisConfig

	_ = godotenv.Load() // load env vars from .env file to program's environment

	switch os.Getenv("ENVIRONMENT") {
	case "development":
		{
			// local database
			localDBConfig, err := config.DBConfig()
			if err != nil {
				log.Fatal(err)
			}
			dbConnStr = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=30", localDBConfig.User, localDBConfig.Pass, localDBConfig.Host, localDBConfig.Port, localDBConfig.Name)

			// local redis (docker)
			redisConfig.Host = "redis"
			redisConfig.Pass = ""
			redisConfig.Port = "6379"
		}
	case "production":
		{
			// cloud database
			neonDbConfig, err := config.NeonDBConfig()
			if err != nil {
				log.Fatal(err)
			}
			dbConnStr = neonDbConfig.NeonConnStr

			// cloud redis
			redisConfiguartion, err := config.RedisConfig()
			if err != nil {
				log.Println(err)
			}

			redisConfig.Host = redisConfiguartion.Host
			redisConfig.Pass = redisConfiguartion.Pass
			redisConfig.Port = redisConfiguartion.Port
		}
	default:
		{
			log.Fatal("error recognizing ENVIRONMENT: Cannot load ENV variable")
		}
	}

	// Get database client
	db, err = driver.InitDB(dbDriver, dbConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// Load data into database
	err = db.LoadDataToDatabase("schema.sql")
	if err != nil {
		log.Fatalf("Failed to SQL file %v", err)
	} else {
		log.Println("SQL file executed successfully!")
	}

	// setup redis connection
	rdb, err = driver.InitRedis(redisConfig.Host, redisConfig.Pass, redisConfig.Port)
	if err != nil {
		log.Printf("Failed to connect to redis: %v", err)
	}

}

func main() {

	var allowedOrigin string
	defer db.Close()

	// Setup mux server for routing
	router := mux.NewRouter()

	// Dependency Injection for modularity
	patientStore := store.NewStore(db.DB, rdb)
	apiRoutes := apiRoutesV1.NewAPIRoutes(patientStore)

	// endpoint to check server health
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		log.Println("health check API triggered")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Server is functioning"})

	}).Methods(http.MethodGet)

	// endpoint to check database stats
	router.HandleFunc("/db/stats", func(w http.ResponseWriter, r *http.Request) {

		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

		if token != os.Getenv("MONITOR_KEY") {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			log.Println("DB stats API triggered but unauthorized")
			return
		}

		stats := db.Stats()

		log.Println("DB stats API triggered")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(stats)
	}).Methods(http.MethodGet)

	if os.Getenv("ENVIRONMENT") != "development" {
		router.Use(middleware.OriginValidator)
	}

	// Public route for patient details
	router.HandleFunc("/api/v1/patients/{token_id}", apiRoutes.GetPatientByTokenID).Methods(http.MethodGet)

	// Login Route
	router.HandleFunc("/api/v1/login", apiRoutes.LoginHandler).Methods(http.MethodPost)

	// Protected Routes
	protectedRouter := router.PathPrefix("/api/v1").Subrouter() // creating subrouter for path "/" that will require authentication
	protectedRouter.Use(middleware.AuthMiddleware)
	protectedRouter.HandleFunc("/patients", apiRoutes.GetAllPatients).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/patients", apiRoutes.CreatePatient).Methods(http.MethodPost)
	protectedRouter.HandleFunc("/patients/{token_id}", apiRoutes.UpdatePatient).Methods(http.MethodPut)
	protectedRouter.HandleFunc("/patients/{token_id}", apiRoutes.UpdatePatientPartial).Methods(http.MethodPatch)
	protectedRouter.HandleFunc("/patients/{token_id}", apiRoutes.DeletePatient).Methods(http.MethodDelete)
	protectedRouter.HandleFunc("/doctors/{doctor_id}", apiRoutes.GetAllPatientsByDocID).Methods(http.MethodGet)

	// Enable CORS
	allowedOriginWebsite, err := config.AllowedOrigin()
	if err != nil {
		fmt.Println(err)
		allowedOrigin = ""
	}
	allowedOrigin = allowedOriginWebsite.AllowedOrigin

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{allowedOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	done := make(chan struct{})

	// Delegate server startup to listen for interrupt signal and server shutdown
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)

			// sending an interrupt to trigger gracefull shutdown
			p, _ := os.FindProcess(os.Getpid())
			_ = p.Signal(os.Interrupt)
		}
		close(done) // signal server has exited
	}()

	shutdownCtx, shutdownCancel := gracefulShutdown()
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	<-done // Wait until ListenAndServe() exits
	log.Println("Server exited")

}

func gracefulShutdown() (context.Context, context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	return shutdownCtx, shutdownCancel
}

// function to generate hashed passwords and verify it
func _() {

	password := []byte("mountain.lucy@medigo")

	// Hashing the password
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	fmt.Println("Hashed:", string(hashedPassword))

	// Verifying a password
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte("mountain.lucy@medigo"))
	if err != nil {
		fmt.Println("Invalid password")
	} else {
		fmt.Println("Password is correct")
	}
}

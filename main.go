package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/gorilla/mux"
	"github.com/harshitrajsinha/medi-go/driver"
	"github.com/harshitrajsinha/medi-go/middleware"
	loginRoutes "github.com/harshitrajsinha/medi-go/routes"
	routesV1 "github.com/harshitrajsinha/medi-go/routes/api/v1"
	"github.com/harshitrajsinha/medi-go/store"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
)

//go:embed store/schema.sql
var schemaFS embed.FS
var db *sql.DB
var rdb *redis.Client

type dbConfig struct {
	User string `envconfig:"DB_USER"` // looks for 'DB_USER' in environment, tag key should match .env key
	Host string `envconfig:"DB_HOST"`
	Port string `envconfig:"DB_PORT"`
	Pass string `envconfig:"DB_PASS"`
	Name string `envconfig:"DB_NAME"`
}

// Function to load data to database via schema file
func loadDataToDatabase(dbClient *sql.DB) error {

	// Read file content
	sqlFile, err := schemaFS.ReadFile("store/schema.sql")
	fmt.Println("...loading schema file")
	if err != nil {
		return err
	}

	// Execute file content (queries)
	_, err = dbClient.Exec(string(sqlFile))
	if err != nil {
		return err
	}
	return nil
}

func init() {

	var cfg dbConfig
	var err error

	_ = godotenv.Load()               // load env from .env file to program's environment
	err = envconfig.Process("", &cfg) // load env from program's environment to declared struct
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	dbConnStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=30", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name)
	dbDriver := "postgres"

	// Get db client
	db, err = driver.InitDB(dbDriver, dbConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// Load data into database
	err = loadDataToDatabase(db)
	if err != nil {
		panic(err)
	} else {
		log.Println("SQL file executed successfully!")
	}

	// setup redis connection
	rdb, err = driver.InitRedis()
	if err != nil {
		log.Printf("Failed to connect to redis: %v", err)
	}

}

func main() {

	defer db.Close()

	// panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	// Setup mux server for routing
	router := mux.NewRouter()

	// Dependency Injection for modularity
	patientStore := store.NewStore(db, rdb)
	loginRoutes := loginRoutes.NewLoginRoutes(patientStore)
	patientRoutes := routesV1.NewPatientRoutes(patientStore)

	// endpoint to check server health
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Server is functioning"})
	}).Methods("GET")

	// Public routes for patient details
	router.HandleFunc("/api/v1/patient/{token_id}", patientRoutes.GetPatientByTokenID).Methods(http.MethodGet)

	// Protected Routes for patients
	router.HandleFunc("/api/v1/login", loginRoutes.LoginHandler).Methods(http.MethodPost)
	protectedRouter := router.PathPrefix("/").Subrouter() // creating subrouter for path "/" that will require authentication
	protectedRouter.Use(middleware.AuthMiddleware)

	protectedRouter.HandleFunc("/api/v1/patients", patientRoutes.GetAllPatients).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/api/v1/patientbydoc/{doctor_id}", patientRoutes.GetAllPatientsByDoc).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/api/v1/patient", patientRoutes.CreatePatient).Methods(http.MethodPost)
	protectedRouter.HandleFunc("/api/v1/patient/{token_id}", patientRoutes.UpdatePatient).Methods(http.MethodPut)
	protectedRouter.HandleFunc("/api/v1/patient/{token_id}", patientRoutes.UpdatePatientPartial).Methods(http.MethodPatch)
	protectedRouter.HandleFunc("/api/v1/patient/{token_id}", patientRoutes.DeletePatient).Methods(http.MethodDelete)

	// Enable CORS
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://medigo-frontend.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Println("Server listening on PORT ", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))

}

func _() {

	// this function is used to generate hashed passwords and verify it
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

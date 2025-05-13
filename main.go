package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"

	"github.com/harshitrajsinha/medi-go/driver"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/crypto/bcrypt"
)

//go:embed store/schema.sql
var schemaFS embed.FS
var db *sql.DB

type dbConfig struct {
	User string `envconfig:"DB_USER"`
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

	// load environment variables
	_ = godotenv.Load()
	err = envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	fmt.Println(cfg)
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=30", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name)
	dbDriver := "postgres"

	// Get db client
	db, err = driver.InitDB(dbDriver, connStr)
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

}

func main() {

	defer db.Close()
}

func hash() {
	password := []byte("priya@medigo")

	// Hashing the password
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	fmt.Println("Hashed:", string(hashedPassword))

	// Verifying a password
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte("priya@medigo"))
	if err != nil {
		fmt.Println("Invalid password")
	} else {
		fmt.Println("Password is correct")
	}
}

package driver

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/harshitrajsinha/medi-go/internal/store"
	_ "github.com/lib/pq"
)

type DBClient struct {
	*sql.DB
}

func InitDB(dbDriver string, connString string) (*DBClient, error) {
	var db *sql.DB
	var err error

	fmt.Println("Waiting for db startup ...")

	// Open connection pool
	db, err = sql.Open(dbDriver, connString)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	// connection timeout at application level
	// (If db.PingContext(ctx) does not complete within 30 seconds, it will be canceled automatically.)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// check if db connection is successful
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		log.Printf("Error connecting to the database: %v", err)
		return nil, err
	}

	fmt.Println("Successfully connected to database")

	return &DBClient{DB: db}, nil

}

// Function to load data to database via schema file
func (rec *DBClient) LoadDataToDatabase(filename string) error {

	// Read file content
	sqlFile, err := store.SchemaFS.ReadFile(filename)
	fmt.Println("...loading schema file")
	if err != nil {
		return err
	}

	// Execute file content (queries)
	_, err = rec.Exec(string(sqlFile))
	if err != nil {
		return err
	}
	return nil
}

package driver

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func InitDB(dbDriver string, connString string) (*sql.DB, error) {

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

	return db, nil

}

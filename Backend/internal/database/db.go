// database connection code

package database

import (
	"database/sql"
	"fmt"
	"time"
)

// postgres ql database

type DB struct {
	*sql.DB
}

// function for connection with db

func New(connectionString string) (*DB, error) {

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// setting connections-pool for db

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// verification of the created connection
	// Ping() to verify the valid connection

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connectiong to database: %w", err)
	}

	fmt.Println("Database connection established")

	return &DB{db}, nil

}

// functions to close the database and checking is it working correctly using health

func (db *DB) Close() error {
	return db.DB.Close()
}

// checking connection strength

func (db *DB) Health() error {
	return db.Ping()
}

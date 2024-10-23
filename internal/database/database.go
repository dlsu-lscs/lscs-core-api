package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func Connect() *sql.DB {
	// MySQL connection string format: username:password@tcp(host:port)/dbname
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)

	dbconn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to the database (DSN parse error or initialization error): %v", err)
	}
	dbconn.SetConnMaxLifetime(0)
	dbconn.SetMaxIdleConns(50)
	dbconn.SetMaxOpenConns(50)

	// test connection
	if err := dbconn.Ping(); err != nil {
		log.Fatalf("Unable to connect to database (Cannot ping...): %v", err)
	}

	return dbconn
}

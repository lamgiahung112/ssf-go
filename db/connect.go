package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDatabase() error {
	var err error
	// Capture connection properties.
	cfg := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DB_ADDR"),
		Logger:               mysql.NewConfig().Logger,
		DBName:               os.Getenv("DB_NAME"),
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	// Get a database handle.
	DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return fmt.Errorf("error opening database connection: %v", err)
	}

	// Test the connection to the database.
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("error pinging database: %v", err)
	}

	log.Println("Connected to MySQL database!")
	return nil
}

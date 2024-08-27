package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var Database *sql.DB

func ConnectDB() {
	var err error

	Database, err = sql.Open("mysql", "root:nandha@tcp(localhost:3306)/time_table")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = Database.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
	fmt.Println("DB connected")
}

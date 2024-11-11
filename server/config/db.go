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

// package config

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"os"

// 	_ "github.com/go-sql-driver/mysql"
// )

// var Database *sql.DB

// func ConnectDB() {
// 	var err error

// 	dbHost := os.Getenv("DB_HOST")
// 	dbPort := os.Getenv("DB_PORT")
// 	dbUser := os.Getenv("DB_USER")
// 	dbPassword := os.Getenv("DB_PASSWORD")
// 	dbName := os.Getenv("DB_NAME")

// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

// 	Database, err = sql.Open("mysql", dsn)
// 	if err != nil {
// 		log.Fatalf("Error opening database: %v", err)
// 	}

// 	err = Database.Ping()
// 	if err != nil {
// 		log.Fatalf("Error pinging database: %v", err)
// 	}
// 	fmt.Println("DB connected")
// }

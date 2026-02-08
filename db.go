package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error

	connStr := "postgres://postgres:123@localhost:5432/configapi?sslmode=disable"

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Database not reachable: ", err)
	}

	log.Println("Connected to PostgreSQL !")
}

package main

import (
	"database/sql"
	"log"
)

// LogError : create a log of the error passed
func LogError(err error) {
	if err != nil {
		log.Println(err)
	}
}

// FatalError : kill process with a log of the error passed
func FatalError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// CreateTableIfDoesntExist : create a new table at the beginning if it doesn't already exist
func CreateTableIfDoesntExist(db *sql.DB) {
	query, err := db.Prepare("CREATE TABLE IF NOT EXISTS todo (id INTEGER PRIMARY KEY, text TEXT, done INTEGER, channel TEXT)")
	LogError(err)

	query.Exec()
}

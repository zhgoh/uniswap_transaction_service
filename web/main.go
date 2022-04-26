package main

import (
	"log"
	"os"
	"strconv"
)

// cryptoTransaction is a data type that is used to store the current transactions info after processing.
type cryptoTransaction struct {
	Hash string  `json:"hash"`
	Fee  float64 `json:"fee"`
}

var client *DBClient

// TODO: Store Transactions in DB, and remove global variables
var db []cryptoTransaction
var latestHash string

// main is just the entry point of our backend service, we will run a goroutine that is polling
// live transactions in the background.
func main() {
	createTableStmt := `CREATE TABLE Transactions(
	"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	"hash" TEXT NOT NULL UNIQUE,
	"fee" REAL NOT NULL
    );` // SQL Statement for Create Table

	var err error
	client, err = makeDBClient("transactions.db", createTableStmt)
	if err != nil {
		log.Panic(err.Error())
	}

	defer client.db.Close()

	freq := 60
	if len(os.Args) > 1 {
		val, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Printf("Usage: %s [frequency]", os.Args[0])
		}
		freq = val
	}
	q := make(chan bool)
	go pollTransactions(q, freq)

	db = []cryptoTransaction{}

	// Serve
	serve("5050")
}

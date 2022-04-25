package main

import (
	"log"
)

type cryptoTransaction struct {
	Hash string  `json:"hash"`
	Fee  float64 `json:"fee"`
}

// TODO: Store Transactions in DB
var db []cryptoTransaction
var latestHash string

func main() {
	log.Println("Transactions service.")

	// TODO: Have a goroutine running in the background to poll live data
	// db, err := sql.Open("sqlite", "./test.sqlite")
	// if err != nil {
	// 	log.Print("Error: opening SQLite db.")

	// }
	// defer db.Close()

	q := make(chan bool)
	go PollTransactions(q)

	db = []cryptoTransaction{}

	// Serve
	Serve("5050")
}

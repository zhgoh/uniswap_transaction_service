package main

import (
	"log"
)

// cryptoTransaction is a data type that is used to store the current transactions info after processing.
type cryptoTransaction struct {
	Hash string  `json:"hash"`
	Fee  float64 `json:"fee"`
}

// TODO: Store Transactions in DB, and remove global variables
var db []cryptoTransaction
var latestHash string

// main is just the entry point of our backend service, we will run a goroutine that is polling
// live transactions in the background.
func main() {
	log.Println("Transactions service.")

	q := make(chan bool)
	go pollTransactions(q)

	db = []cryptoTransaction{}

	// Serve
	serve("5050")
}

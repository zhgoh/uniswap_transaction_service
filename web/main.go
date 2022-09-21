package main

import (
	"log"
	"math/big"
	"os"
	"strconv"
)

// cryptoTransaction is a data type that is used to store the current transactions info after processing.
type cryptoTransaction struct {
	Hash string     `json:"hash"`
	USDC *big.Float `json:"usdc"`
	ETH  *big.Float `json:"eth"`
	Fee  float64    `json:"fee"`
}

var db *dbClient
var latestHash string

// main is just the entry point of our backend service, we will run a goroutine that is polling
// live transactions in the background.
func main() {
	var err error
	db, err = makeDBClient("transactions.db")
	if err != nil {
		log.Panic(err.Error())
	}

	defer db.db.Close()

	// Check if etherscan api key is set
	if os.Getenv("etherscan_api") == "" {
		log.Print("Error: Please specify etherscan_api env variable before running.")
		log.Print("Closing service...")
		return
	}

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

	// Serve
	serve("5050")
}

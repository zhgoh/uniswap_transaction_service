package main

import (
	"log"
	"time"
)

func pollTransactions(quit chan bool, freq int) {
	log.Printf("Polling live transactions every %ds", freq)

	etherClient, err := makeEtherscan()
	if err != nil {
		log.Print("Error: did not create etherscan client properly.")
		log.Print("Shutting down live transactions fetching.")
		return
	}

	binanceClient := makeBinanceClient()

	for {
		select {
		case <-quit:
			log.Print("Polling stopped.")
			return
		default:
			log.Print("Checking for transactions.")

			// Will fetch latest price from order books and used it to store in latest transactions
			prices, err := binanceClient.getOrderBook("ETHUSDT", 1)
			if err != nil {
				log.Print("Error: getting prices, will try again later in 5s ", err.Error())
				time.Sleep(5 * time.Second)
				continue
			}

			etherTransactions, err := etherClient.fetchTransactions(0, 0)
			if err != nil {
				// Log error and try again later
				log.Print("Error: Failed to fetch etherscan transaction, trying again in 5s ", err.Error())
				time.Sleep(5 * time.Second)
				continue
			}

			if err := db.addLiveTransactions(etherTransactions, prices); err != nil {
				log.Print("Error: getting transactions, will try again later in 5s ", err.Error())
				time.Sleep(5 * time.Second)
				continue
			}

			// Try fetching again
			time.Sleep(time.Duration(freq) * time.Second)
		}
	}
}

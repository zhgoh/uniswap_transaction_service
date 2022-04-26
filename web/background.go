package main

import (
	"fmt"
	"log"
	"strconv"
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
				log.Print("Error: getting prices, will try again later")
				log.Print(err)
				continue
			}

			etherTransactions, err := etherClient.fetchTransactions(0, 0)
			if err != nil {
				// Log error and try again later
				log.Print("Error: Failed to fetch etherscan transaction")
				log.Print(err)
				continue
			}

			if err := addLiveTransactions(etherTransactions, prices); err != nil {
				log.Print("Error: getting transactions, will try again later")
				log.Print(err)
				continue
			}

			// Try fetching again
			time.Sleep(time.Duration(freq) * time.Second)
		}
	}
}

func addLiveTransactions(etherTransactions []etherscanTransaction, prices float64) error {
	if len(etherTransactions) == 0 {
		return fmt.Errorf("no transactions provided")
	}

	for _, v := range etherTransactions {
		if len(v.Hash) == 0 {
			return fmt.Errorf("hash is empty.")
		}

		if v.Hash == latestHash {
			break
		}

		err := addSingleTransaction(v, prices)
		if err != nil {
			return err
		}

	}
	latestHash = etherTransactions[0].Hash
	return nil
}

func addSingleTransaction(transaction etherscanTransaction, prices float64) error {
	for _, v := range db {
		if v.Hash == transaction.Hash {
			// Ignore duplicates
			return nil
		}
	}

	// Compute prices
	gasPrice, err := strconv.Atoi(transaction.GasPrice)
	if err != nil {
		log.Print("Error: failed to convert gas price to integer.")
		return err
	}

	gasUsed, err := strconv.Atoi(transaction.GasUsed)
	if err != nil {
		log.Print("Error: failed to convert gas used to integer.")
		return err
	}

	// Fees in eth
	// Note: no idea if division or multiplying would be faster here, probably same
	// fees := float64(gasPrice*gasUsed) / 1000000000000000000
	fees := float64(gasPrice*gasUsed) * 0.000000000000000001
	fees *= prices

	// Convert to price in USDT
	db = append(db, cryptoTransaction{transaction.Hash, fees})

	timeStamp, err := strconv.Atoi(transaction.TimeStamp)
	if err != nil {
		log.Print("Error: failed to convert timeStamp.")
		return err
	}

	// TODO: Add to DB
	log.Printf("Hash: %s, Time: %d, Fees: $%.2f", transaction.Hash, timeStamp, fees)
	return nil
}

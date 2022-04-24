package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

func PollTransactions(quit chan bool) {
	log.Print("Polling live transactions.")

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

			etherTransactions, err := etherClient.fetchTransactions()
			if err != nil {
				// Log error and try again later
				log.Print("Error: Failed to fetch etherscan transaction")
				log.Print(err)
				continue
			}

			if err := addTransactions(etherTransactions, prices); err != nil {
				log.Print("Error: getting transactions, will try again later")
				log.Print(err)
				continue
			}

			// Try fetching again
			time.Sleep(60 * time.Second)
		}
	}
}

//func getDailyPrice(client *BinanceClient, time int64) (map[int]float64, error) {
//	// Get daily prices data from 0 to current time
//	klineResp, err := client.getKlines("ETHUSDT", 1, Days, 0, time, 0)
//	if err != nil {
//		log.Print("Error: Failed to get kline results")
//	}
//
//	// Collate the price from kline api
//	prices := make(map[int]float64)
//	for _, v := range klineResp {
//		close, err := strconv.ParseFloat(v.Close, 64)
//		if err != nil {
//			log.Print("Error: failed to convert closing price")
//			return nil, err
//		}
//		prices[v.CloseTime] = close
//	}
//	return prices, nil
//}

func addTransactions(etherTransactions []etherscanTransaction, prices float64) error {
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

		// Compute prices
		gasPrice, err := strconv.Atoi(v.GasPrice)
		if err != nil {
			log.Print("Error: failed to convert gas price to integer.")
			return err
		}

		gasUsed, err := strconv.Atoi(v.GasUsed)
		if err != nil {
			log.Print("Error: failed to convert gas used to integer.")
			return err
		}

		timeStamp, err := strconv.Atoi(v.TimeStamp)
		if err != nil {
			log.Print("Error: failed to convert timeStamp.")
			return err
		}

		// Fees in eth
		// Note: no idea if division or multiplying would be faster here, probably same
		// fees := float64(gasPrice*gasUsed) / 1000000000000000000
		fees := float64(gasPrice*gasUsed) * 0.000000000000000001

		// Convert to price in USDT
		fees *= prices

		// TODO: Add to DB
		log.Printf("Hash: %s, Time: %d, Fees: $%.2f", v.Hash, timeStamp, fees)
		transactions = append(transactions, cryptoTransaction{v.Hash, fees})
	}
	latestHash = etherTransactions[0].Hash
	return nil
}

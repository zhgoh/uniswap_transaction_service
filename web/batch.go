package main

import (
	"log"
	"math"
	"strconv"
	"time"
)

type klinePrice struct {
	timeStamp time.Time
	price     float64
}

func batch(startTime, endTime time.Time) error {
	log.Print("Starting batch job.")
	log.Printf("Start: %d, End: %d", startTime.Unix(), endTime.Unix())

	// Split the code into internal scope for easier to manage variables
	var transactions []etherscanTransaction
	{
		// Pull etherscan transactions by date (2 steps)
		// Get block num by start and end timestamp
		etherScanClient, err := makeEtherscan()
		if err != nil {
			log.Print("Error: batch job failed to create etherscan client")
			return err
		}

		startBlock, err := etherScanClient.getBlockNumberByTimestamp(before, startTime)
		if err != nil {
			log.Print("Error: batch job failed to get start block number by timestamp")
			return err
		}

		endBlock, err := etherScanClient.getBlockNumberByTimestamp(before, endTime)
		if err != nil {
			log.Print("Error: batch job failed to get end block number by timestamp")
			return err
		}

		// Get transactions limited by start and end block number
		transactions, err = etherScanClient.fetchTransactions(startBlock, endBlock)
		if err != nil {
			log.Print("Error: batch job failed to get transactions.")
			return err
		}

		for _, v := range transactions {
			log.Print(v.TimeStamp)
		}
	}

	// Binance kline api pull
	{
		binanceClient := makeBinanceClient()
		// Limitations of kline api is there is max of 1000 entries, however the time given might
		// be longer than that, hence I propose we pull all the kline data first, and then go through the
		// transactions data

		// Keep pulling kline data (up to 1000 results), till batch job is done
		klineData := []klinePrice{}
		for startTime.Before(endTime) {
			klineResp, err := binanceClient.getKlines("ETHUSDT", 1, Minutes, startTime, endTime, 1000)

			if err != nil {
				log.Print("Error: when pulling kline data from binance")
				continue
			}

			if len(klineResp) == 0 {
				log.Print("Finished getting all kline data.")
				// Break when no kline response, probably finished
				break
			}

			for _, v := range klineResp {
				closePrice, err := strconv.ParseFloat(v.Close, 64)
				if err != nil {
					log.Print("Error: while converting kline closing price")
					continue
				}
				klineData = append(klineData, klinePrice{price: closePrice, timeStamp: time.Unix(v.CloseTime, 0)})
			}

			log.Printf("Start: %d, End: %d", startTime.Unix(), endTime.Unix())
			startTime = time.UnixMilli(klineResp[len(klineResp)-1].CloseTime)
		}

		// Once we get all the kline data, iterate through all possible transactions
		for _, v := range transactions {
			timeStamp, err := strconv.ParseInt(v.TimeStamp, 10, 64)
			if err != nil {
				log.Print("Error: converting timestamp in batch")
				return err
			}

			log.Printf("Timestamp: %v, kline: %v", timeStamp, klineData[0].timeStamp)
			price := binarySearchKline(klineData, time.Unix(timeStamp, 0))

			err = addSingleTransaction(v, price)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

const float64EqualityThreshold = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

// binarySearchKline is used for finding sorted kline data as an optimized way to search
func binarySearchKline(arr []klinePrice, timeStamp time.Time) float64 {
	if len(arr) == 0 {
		return 0
	}

	l, r := 0, len(arr)-1
	price := arr[0].price
	for l <= r {
		mid := (l + r) / 2
		if timeStamp.Before(arr[mid].timeStamp) {
			price = arr[mid].price
			r = mid - 1
		} else {
			l = mid + 1
		}
	}

	return price
}

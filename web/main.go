package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Transaction struct {
	Hash string  `json:"hash"`
	Fee  float64 `json:"fee"`
}

type TransactionRequest struct {
	Hash string `json:"hash"`
}

type TransactionResponse struct {
	ErrorCode    int         `json:"errorcode"`
	Transactions Transaction `json:"transactions"`
	Message      string      `json:"message"`
}

type BatchRequest struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type BatchResponse struct {
	ErrorCode int    `json:"errorcode"`
	Message   string `json:"message"`
}

// TODO: Store Transactions in DB
var Transactions []Transaction
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
	go pollTransactions(q)

	Transactions = []Transaction{}

	// Serve
	serve("5050")
}

func serve(port string) {
	// TODO: Provide REST api for the following,
	http.HandleFunc("/", homePage)

	// TODO: API: batch job
	http.HandleFunc("/batch", batch)

	// TODO: API: Get transaction fee given transaction hash
	http.HandleFunc("/transactions", transaction)
	http.HandleFunc("/transactions/all", allTransaction)

	// TODO: Bonus API: Decode actual Uniswap swap price

	// Listen to port
	log.Printf("Running server at port at %s.\n", port)
	port = fmt.Sprintf(":%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Endpoints\n")
	fmt.Fprintf(w, "GET  /transactions\n")
	fmt.Fprintf(w, "PUT  /batch\n")
}

func batch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(404)
		log.Print("Only PUT method supported on batch.")
		return
	}

	batchResp := BatchResponse{
		ErrorCode: 0,
		Message:   "Successfully process batch request.",
	}
	log.Print("Batch job request")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print("Error: Cannot read batch request.")
		batchResp.ErrorCode = 1
		batchResp.Message = "Error: cannot read batch request."
		return
	}

	var batchReq BatchRequest
	err = json.Unmarshal(reqBody, &batchReq)
	if err != nil {
		log.Print("Error: Cannot unmarshal batch request.")
		batchResp.ErrorCode = 1
		batchResp.Message = "Error: cannot unmarshal batch request."
		return
	}

	log.Printf("Processing transactions between %s and %s", batchReq.Start, batchReq.End)
	json.NewEncoder(w).Encode(batchResp)
}

func transaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(404)
		log.Print("Only GET method supported on transactions.")
		return
	}

	log.Print("Getting transactions.")
	transactionResp := TransactionResponse{
		Message:   "No transactions hash found",
		ErrorCode: 1,
	}

	hashes, ok := r.URL.Query()["hash"]
	if !ok || len(hashes[0]) < 1 {
		log.Print("Error: url param hash is missing.")
		json.NewEncoder(w).Encode(transactionResp)
		return
	}

	// Get the transaction
	for _, transaction := range Transactions {
		if transaction.Hash == hashes[0] {
			transactionResp.Message = "Found transactions."
			transactionResp.ErrorCode = 0
			transactionResp.Transactions = transaction
			break
		}
	}
	json.NewEncoder(w).Encode(transactionResp)
}

func allTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(404)
		log.Print("Only GET method supported on transactions.")
		return
	}

	log.Print("Getting all transactions.")
	json.NewEncoder(w).Encode(Transactions)
}

func pollTransactions(quit chan bool) {
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

func addTransactions(etherTransactions []EtherscanTransaction, prices float64) error {
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
		Transactions = append(Transactions, Transaction{v.Hash, fees})
	}
	latestHash = etherTransactions[0].Hash
	return nil
}

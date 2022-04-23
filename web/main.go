package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
	// "database/sql"
	//_ "modernc.org/sqlite"
)

type Transaction struct {
	Hash string `json:"hash"`
	Fee  string `json:"fee"`
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

	// q := make(chan bool)
	// go pollTransactions(q)

	Transactions = []Transaction{
		{Hash: "0x1", Fee: "10"},
	}

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

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print("Error: Cannot read transactions request.")
		return
	}

	var transactionReq TransactionRequest
	err = json.Unmarshal(reqBody, &transactionReq)
	if err != nil {
		log.Print("Error: Cannot unmarshal transactions.")
		return
	}

	transactionResp := TransactionResponse{
		Message:   "No transactions hash found",
		ErrorCode: 1,
	}

	// Get the transaction
	for _, transaction := range Transactions {
		if transaction.Hash == transactionReq.Hash {
			transactionResp.Message = "Found transactions."
			transactionResp.ErrorCode = 0
			transactionResp.Transactions = transaction
			break
		}
	}
	json.NewEncoder(w).Encode(transactionResp)
}

func pollTransactions(quit chan bool) {
	log.Print("Polling live transactions.")
	etherClient, err := makeEtherscan()
	if err != nil {
		log.Print("Error: did not create etherscan client properly.")
		return
	}

	//for {
	select {
	case <-quit:
		log.Print("Polling stopped.")
		return
	default:
		log.Print("Checking for transactions.")
		transactions, err := etherClient.fetchTransactions()
		if err != nil {
			// Log error and try again later
			log.Print("Error: Failed to fetch transaction")
			log.Print(err)
		}
		// Add transactions to db
		addTransactions(transactions)
		time.Sleep(60 * time.Second)
	}
	//}
}

func addTransactions(transactions []EtherscanTransaction) {
	for _, v := range transactions {
		if v.Hash == latestHash {
			latestHash = transactions[0].Hash
			break
		}

		gasPrice, err := strconv.Atoi(v.GasPrice)
		if err != nil {
			log.Print("Error: failed to convert gas price to integer.")
			continue
		}

		gasUsed, err := strconv.Atoi(v.GasUsed)
		if err != nil {
			log.Print("Error: failed to convert gas used to integer.")
			continue
		}

		price := float64(gasPrice*gasUsed) / 100000000000000000
		log.Printf("Hash: %s, Price: %f\n", v.Hash, price)

		//Transactions = append(Transactions, Transaction{Hash: v.Hash, Fee: v.})
	}

	if len(latestHash) == 0 {
		latestHash = transactions[0].Hash
	}

	log.Printf("Latest hash: %s\n", latestHash)
}

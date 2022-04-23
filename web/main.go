package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Transaction struct {
	Hash string
	Fee  string
}

type TransactionRequest struct {
	Hash string
}

type TransactionResponse struct {
	ErrorCode    int
	Transactions Transaction
	Message      string
}

// TODO: Store Transactions in DB
var Transactions []Transaction

func main() {
	log.Println("Transactions service")
	// TODO: Have a goroutine running in the background to poll live data

	Transactions = []Transaction{
		{Hash: "0x1", Fee: "10"},
		{Hash: "0x2", Fee: "10"},
		{Hash: "0x3", Fee: "10"},
	}

	// Serve
	serve("5050")
}

func serve(port string) {
	// TODO: Provide REST api for the following,

	//myRouter := mux.NewRouter().StrictSlash(true)
	//myRouter.HandleFunc("/", homePage)
	http.HandleFunc("/", homePage)

	// TODO: API: batch job
	// myRouter.HandleFunc("/batch", batch).Methods("POST")
	http.HandleFunc("/batch", batch)

	// TODO: API: Get transaction fee given transaction hash
	// myRouter.HandleFunc("/transactions", transaction).Methods("POST")
	http.HandleFunc("/transactions", transaction)

	// TODO: Bonus API: Decode actual Uniswap swap price

	// Listen to port
	log.Printf("Running server at port at %s\n", port)
	port = fmt.Sprintf(":%s", port)
	//log.Fatal(http.ListenAndServe(port, myRouter))
	log.Fatal(http.ListenAndServe(port, nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Endpoints\n")
	fmt.Fprintf(w, "POST /batch\n")
	fmt.Fprintf(w, "POST  /transactions\n")
}

func batch(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Batch\n")
}

func transaction(w http.ResponseWriter, r *http.Request) {
	log.Print("Getting transactions...")

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
			transactionResp.Message = "Found transactions"
			transactionResp.ErrorCode = 0
			transactionResp.Transactions = transaction
			break
		}
	}
	json.NewEncoder(w).Encode(transactionResp)
}

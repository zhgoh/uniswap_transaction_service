package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type transactionResponse struct {
	ErrorCode    int               `json:"errorcode"`
	Transactions cryptoTransaction `json:"transactions"`
	Message      string            `json:"message"`
}

type batchRequest struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type batchResponse struct {
	ErrorCode int    `json:"errorcode"`
	Message   string `json:"message"`
}

func Serve(port string) {
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

	batchResp := batchResponse{
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

	var batchReq batchRequest
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
	transactionResp := transactionResponse{
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
	for _, transaction := range transactions {
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
	json.NewEncoder(w).Encode(transactions)
}

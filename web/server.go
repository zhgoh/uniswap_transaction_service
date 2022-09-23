package main

import (
	"encoding/json"
	"fmt"
	"io"
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

func serve(port string) {
	// TODO: Provide REST api for the following,
	http.HandleFunc("/", homePageHandler)

	// TODO: API: batch job
	http.HandleFunc("/batch", batchHandler)

	// TODO: API: Get transaction fee and (swap details) given transaction hash
	http.HandleFunc("/transactions", transactionHandler)
	http.HandleFunc("/transactions/all", allTransactionHandler)

	// Listen to port
	log.Printf("Running server at port at %s.\n", port)
	port = fmt.Sprintf(":%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Endpoints\n")
	fmt.Fprintf(w, "GET  /transactions?id=\n")
	fmt.Fprintf(w, "GET  /transactions/all\n")
	fmt.Fprintf(w, "PUT  /batch\n")
}

func batchHandler(w http.ResponseWriter, r *http.Request) {
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

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Print("Error: Cannot read batch request.")
		batchResp.ErrorCode = 1
		batchResp.Message = "Error: cannot read batch request."
		return
	}

	var batchReq batchRequest
	err = json.Unmarshal(reqBody, &batchReq)
	if err != nil {
		log.Print("Error: cannot unmarshal batch request.")
		batchResp.ErrorCode = 1
		batchResp.Message = "Error: cannot unmarshal batch request."
		return
	}

	log.Printf("Processing transactions between %s and %s", batchReq.Start, batchReq.End)
	err = batch(batchReq.Start, batchReq.End)
	if err != nil {
		log.Print(err)
		batchResp.ErrorCode = 1
		batchResp.Message = err.Error()
	}

	json.NewEncoder(w).Encode(batchResp)
}

func transactionHandler(w http.ResponseWriter, r *http.Request) {
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

	hashes, ok := r.URL.Query()["id"]
	if !ok || len(hashes[0]) < 1 {
		log.Print("Error: url param id is missing.")
		json.NewEncoder(w).Encode(transactionResp)
		return
	}

	// Get the transaction
	res, err := db.getTransaction(hashes[0])
	if err != nil {
		log.Print("Error getting transaction: ", err.Error())
		return
	}

	if res != nil {
		transactionResp.Message = "Found transactions."
		transactionResp.ErrorCode = 0
		transactionResp.Transactions = *res
	}
	json.NewEncoder(w).Encode(transactionResp)
}

func allTransactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(404)
		log.Print("Only GET method supported on transactions.")
		return
	}

	log.Print("Getting all transactions.")
	res, err := db.getAllTransactions()
	if err != nil {
		log.Print("Error: getting all transactions from DB")
		json.NewEncoder(w).Encode(`{"ErrorCode": 1, "Message": "Error getting all transactions from DB."}`)
		return
	}
	json.NewEncoder(w).Encode(res)
}

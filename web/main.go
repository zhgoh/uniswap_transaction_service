package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Print("Transactions service")
	// TODO: Have a goroutine running in the background to poll live data

	// Serve
	serve("5050")
}

func serve(port string) {
	// TODO: Provide REST api for the following,
	// Handle route using handler func
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Endpoints \n")
		fmt.Fprintf(w, "POST /batch")
		fmt.Fprintf(w, "GET  /transactions/<hash>")
	})

	// TODO: API: batch job
	http.HandleFunc("/batch", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Batch")
	})

	// TODO: API: Get transaction fee given transaction hash
	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Transactions")
	})

	// TODO: Bonus API: Decode actual Uniswap swap price

	// Listen to port
	fmt.Printf("Running server at port at %s", port)
	port = fmt.Sprintf(":%s", port)
	http.ListenAndServe(port, nil)
}

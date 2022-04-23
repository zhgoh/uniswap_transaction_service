package main

import (
	"fmt"
	"net/http"
)

func main() {
	// TODO: Have a goroutine running in the background to poll live data
	// TODO: Provide REST api for the following,
	// TODO: API: batch job
	// TODO: API: Get transaction fee given transaction hash
	// TODO: Bonus API: Decode actual Uniswap swap price
	// Handle route using handler func
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to new server!")
	})

	// Listen to port
	const port = ":5050"
	fmt.Printf("Running server at port at %s", port)
	http.ListenAndServe(port, nil)
}

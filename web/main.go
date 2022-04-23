package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Handle route using handler func
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to new server!")
	})

	// Listen to port
	const port = ":5050"
	fmt.Printf("Running server at port at %s", port)
	http.ListenAndServe(port, nil)
}

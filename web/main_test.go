package main

import (
	"testing"
)

func Test_AddTransaction(t *testing.T) {

	client, err := makeEtherscan()
	if err != nil {
		t.Fatal(err)
	}

	transaction, err := client.fetchTransactions()
	if err != nil {
		t.Error(err)
	}
	addTransactions(transaction)

	if len(transaction) != 100 {
		t.Fail()
	}
}

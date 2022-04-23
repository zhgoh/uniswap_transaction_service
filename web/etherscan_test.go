package main

import (
	"testing"
)

func Test_FetchTransaction(t *testing.T) {

	client, err := makeEtherscan()
	if err != nil {
		t.Fatal(err)
	}

	transaction, err := client.fetchTransactions()
	if err != nil {
		t.Error(err)
	}

	if len(transaction) != 100 {
		t.Fail()
	}
}

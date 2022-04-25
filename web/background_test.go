package main

import "testing"

func Test_addTransactions(t *testing.T) {
	db = []cryptoTransaction{}

	{
		transactions := []etherscanTransaction{{
			BlockNumber: "0x1",
			GasPrice:    "48335977034",
			GasUsed:     "250621",
			TimeStamp:   "1650727781",
			Hash:        "0x6",
		}, {
			BlockNumber: "0x2",
			GasPrice:    "48335977034",
			GasUsed:     "250621",
			TimeStamp:   "1650727730",
			Hash:        "0x4",
		}}

		prices := 100.1
		err := addLiveTransactions(transactions, prices)
		if err != nil {
			t.Fatal(err)
		}

		got := len(db)
		want := 2

		if got != want {
			t.Errorf("Got %d, Want: %d", got, want)
		}
	}
}

func Test_addSingleTransaction(t *testing.T) {
	db = []cryptoTransaction{}

	{
		transaction := etherscanTransaction{
			BlockNumber: "0x1",
			GasPrice:    "44901991519",
			GasUsed:     "159030",
			TimeStamp:   "1650727793",
		}
		prices := 1000.0
		err := addSingleTransaction(transaction, prices)
		if err != nil {
			t.Fatal(err)
		}

		got := len(db)
		want := 1

		if got != want {
			t.Errorf("Got %d, Want: %d", got, want)
		}
	}
}

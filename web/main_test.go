package main

import (
	"testing"
)

func Test_AddTransaction(t *testing.T) {
	etherScanTransactions := []EtherscanTransaction{
		{
			Hash:      "0xf5bc869730283da55772add53c542ad1cb9d9f8452d20c62fb4141224812cabc",
			GasPrice:  "44401991519",
			GasUsed:   "149542",
			TimeStamp: "1650727793",
		},
	}
	var prices float64 = 2948.71

	if err := addTransactions(etherScanTransactions, prices); err != nil {
		t.Fatal("Error adding transactions")
	}
}

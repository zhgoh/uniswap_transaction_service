package main

import (
	"testing"
)

func Test_AddTransaction(t *testing.T) {
	etherScanTransactions := []etherscanTransaction{
		{
			Hash:      "0x48888e465a61d4f9908dab1d18d9ab67d8227d72a44f58ecb00750b719df9b9c",
			GasPrice:  "56741962048",
			GasUsed:   "1000694",
			TimeStamp: "1650727773",
		},
		{
			Hash:      "0xcbdf28fe5ddf07938f137aba50b85ba146d107707db0356a4b582395909f3f1f",
			GasPrice:  "58158832546",
			GasUsed:   "189032",
			TimeStamp: "1650727726",
		},
		{
			Hash:      "0x69aaa97b540fe8aeef5e35fdfc1d74dfc4f6e13b449d58772c301bdced1e1133",
			GasPrice:  "58158832546",
			GasUsed:   "138248",
			TimeStamp: "1650727726",
		},
	}
	var prices float64 = 2948.71

	if err := addTransactions(etherScanTransactions, prices); err != nil {
		t.Fatal("Error adding transactions")
	}

	etherScanTransactions = []etherscanTransaction{
		{
			Hash:      "0xf5bc869730283da55772add53c542ad1cb9d9f8452d20c62fb4141224812cabc",
			GasPrice:  "44901991519",
			GasUsed:   "159030",
			TimeStamp: "1650727793",
		},
		{
			Hash:      "0x90d3d525aa2ec5b5f0a644640002e7d40e8521b218b20856dc47f466536eddc6",
			GasPrice:  "48335977034",
			GasUsed:   "250621",
			TimeStamp: "1650727781",
		},
		{
			Hash:      "0x48888e465a61d4f9908dab1d18d9ab67d8227d72a44f58ecb00750b719df9b9c",
			GasPrice:  "56741962048",
			GasUsed:   "1000694",
			TimeStamp: "1650727773",
		},
	}

	if err := addTransactions(etherScanTransactions, prices); err != nil {
		t.Fatal("Error adding transactions.")
	}

	if len(db) != 5 {
		t.Fatal("Error when adding transactions. Some entries might be duplicated.")
	}
}

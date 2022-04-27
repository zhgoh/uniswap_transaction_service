package main

import "testing"

func Test_addTransactions(t *testing.T) {
	// Test adding 2 transactions
	//
	var err error
	db, err = makeDBClient("Test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.db.Close()

	{
		err = db.clearTable()
		if err != nil {
			t.Fatal(err)
		}

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
		err := db.addLiveTransactions(transactions, prices)
		if err != nil {
			t.Fatal(err)
		}

		got, err := db.getAllTransactions()
		if err != nil {
			t.Fatal(err)
		}
		want := 2

		if len(got) != want {
			t.Errorf("Got %d, Want: %d", len(got), want)
		}
	}

	// Test adding 3 transactions, with some overlap to see latest hash work or not
	{
		err = db.clearTable()
		if err != nil {
			t.Fatal(err)
		}
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

		if err := db.addLiveTransactions(etherScanTransactions, prices); err != nil {
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

		if err := db.addLiveTransactions(etherScanTransactions, prices); err != nil {
			t.Fatal("Error adding transactions.")
		}

		got, err := db.getAllTransactions()
		if err != nil {
			t.Fatal(err)
		}
		want := 5

		if len(got) != want {
			t.Errorf("Got %d, Want: %d", len(got), want)
		}
	}
}

func Test_addSingleTransaction(t *testing.T) {
	var err error
	db, err = makeDBClient("Test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.db.Close()

	{
		err = db.clearTable()
		if err != nil {
			t.Fatal(err)
		}
		transaction := etherscanTransaction{
			BlockNumber: "0x1",
			GasPrice:    "44901991519",
			GasUsed:     "159030",
			TimeStamp:   "1650727793",
			Hash:        "0x4123",
		}
		prices := 1000.0
		err := db.addSingleTransaction(transaction, prices)
		if err != nil {
			t.Fatal(err)
		}

		got, err := db.getAllTransactions()
		if err != nil {
			t.Fatal(err)
		}
		want := 1

		if len(got) != want {
			t.Errorf("Got %d, Want: %d", len(got), want)
		}
	}

	{
		// Should not be adding duplicate entry
		transaction := etherscanTransaction{
			BlockNumber: "0x1",
			GasPrice:    "44901991519",
			GasUsed:     "159030",
			TimeStamp:   "1650727793",
			Hash:        "0x4123",
		}
		prices := 1000.0
		err := db.addSingleTransaction(transaction, prices)
		if err != nil {
			t.Fatal(err)
		}

		got, err := db.getAllTransactions()
		if err != nil {
			t.Fatal(err)
		}
		want := 1

		if len(got) != want {
			t.Errorf("Got %d, Want: %d", len(got), want)
		}
	}
}

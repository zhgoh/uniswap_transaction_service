package main

import (
	"math/big"
	"os"
	"testing"
)

func Test_initDB(t *testing.T) {
	fileName := "Test.db"
	var err error
	db, err = makeDBClient(fileName)
	if err != nil {
		t.Fatal(err)
	}
	defer db.db.Close()
	// Remove before we do any testing
	if _, err := os.Stat(fileName); err == nil {
		t.Log("Removing files")
		os.Remove(fileName)
	}

}

func Test_addTransactionToDB(t *testing.T) {
	var err error
	db, err = makeDBClient("Test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.db.Close()
	err = db.clearTable()
	if err != nil {
		t.Fatal(err)
	}

	{
		want := 0.4
		err = db.addTransaction(cryptoTransaction{"0x1", big.NewFloat(0), big.NewFloat(0), want})
		if err != nil {
			t.Fatal(err.Error())
		}

		got, err := db.getTransaction("0x1")
		if err != nil || got == nil {
			t.Fatal(err.Error())
		}

		if got.Fee != want {
			t.Fatalf("Want: %f, Got: %f", got.Fee, want)
		}
	}

	{
		err := db.addTransaction(cryptoTransaction{"0x2", big.NewFloat(0), big.NewFloat(0), 0.5})
		if err != nil {
			t.Fatal(err.Error())
		}

		got, err := db.getAllTransactions()
		if err != nil || len(got) == 0 {
			t.Fatal(err.Error())
		}

		want := 2
		if len(got) != want {
			t.Fatalf("Want: %d, Got: %d", len(got), want)
		}
	}

	{
		err = db.clearTable()
		if err != nil {
			t.Fatal(err)
		}

		want := 0
		got, err := db.getAllTransactions()
		if err != nil {
			t.Fatal(err.Error())
		}

		if len(got) != want {
			t.Fatalf("Want: %d, Got: %d", len(got), want)
		}
	}
}

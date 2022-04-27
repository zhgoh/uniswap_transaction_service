package main

import (
	"os"
	"testing"
)

func Test_initDB(t *testing.T) {
	const fileName = "test.db"
	{
		// Remove before we do any testing
		if _, err := os.Stat(fileName); err == nil {
			t.Log("Removing files")
			os.Remove(fileName)
		}

		dbClient, err := makeDBClient(fileName)
		if err != nil {
			t.Fatal(err)
		}
		defer dbClient.db.Close()

		want := 0.4
		err = dbClient.addTransaction(cryptoTransaction{"0x1", want})
		if err != nil {
			t.Fatal(err.Error())
		}

		got, err := dbClient.getTransaction("0x1")
		if err != nil || got == nil {
			t.Fatal(err.Error())
		}

		if got.Fee != want {
			t.Fatalf("Want: %f, Got: %f", got.Fee, want)
		}
	}
}

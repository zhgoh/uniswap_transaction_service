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

		createTableStmt := `CREATE TABLE Transactions(
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"hash" TEXT NOT NULL UNIQUE,
		"fee" REAL NOT NULL
		);` // SQL Statement for Create Table

		dbClient, err := makeDBClient(fileName, createTableStmt)
		if err != nil {
			t.Fatal(err)
		}
		defer dbClient.db.Close()

		want := 0.4
		err = dbClient.addTransactions(cryptoTransaction{"0x1", want})
		if err != nil {
			t.Fatal(err.Error())
		}

		got, err := dbClient.getTransactions("0x1")
		if err != nil {
			t.Fatal(err.Error())
		}

		if got != want {
			t.Fatalf("Want: %f, Got: %f", got, want)

		}

	}
}

package main

import (
	"testing"
	"time"
)

func Test_FetchTransaction(t *testing.T) {
	client, err := makeEtherscan()
	if err != nil {
		t.Fatal(err)
	}

	{
		transaction, err := client.fetchTransactions(0, 0)
		if err != nil {
			t.Error(err)
		}
		want := 100
		if len(transaction) != want {
			t.Errorf("Want: %d, got %d", want, len(transaction))
		}

	}

	{
		transaction, err := client.fetchTransactions(14653342, 14653342)
		if err != nil {
			t.Error(err)
		}

		want := 4
		if len(transaction) != want {
			t.Errorf("Want: %d, got %d", want, len(transaction))
		}
	}
}

func Test_FetchBlock(t *testing.T) {
	client, err := makeEtherscan()
	if err != nil {
		t.Fatal(err)
	}

	timeStamp, err := time.Parse(time.RFC3339, "2022-01-02T06:28:43.000Z")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(timeStamp.Unix())

	got, err := client.getBlockNumberByTimestamp(before, timeStamp)
	if err != nil {
		t.Error(err)
	}

	want := 13924429
	if got == 0 || got != want {
		t.Errorf("Want: %d, got %d", want, got)
	}
}

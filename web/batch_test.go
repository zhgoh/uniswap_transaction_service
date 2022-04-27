package main

import (
	"testing"
	"time"
)

func Test_binarySearch(t *testing.T) {
	arr := []klinePrice{
		{price: 0.9, timeStamp: time.Unix(1645779980, 0)},
		{price: 0.8, timeStamp: time.Unix(1645779983, 0)},
		{price: 0.7, timeStamp: time.Unix(1645779986, 0)},
		{price: 0.6, timeStamp: time.Unix(1645779989, 0)},
		{price: 0.5, timeStamp: time.Unix(1645779991, 0)},
		{price: 0.4, timeStamp: time.Unix(1645779993, 0)},
		{price: 0.3, timeStamp: time.Unix(1645779996, 0)},
	}

	{
		got := binarySearchKline(arr, time.Unix(1645779982, 0))
		want := 0.8
		if got != want {
			t.Fatalf("Want %f, Got %f", want, got)
		}
	}

	{
		got := binarySearchKline(arr, time.Unix(1645779995, 0))
		want := 0.3
		if got != want {
			t.Fatalf("Want %f, Got %f", want, got)
		}
	}

	{
		got := binarySearchKline(arr, time.Unix(1645779979, 0))
		want := 0.9
		if got != want {
			t.Fatalf("Want %f, Got %f", want, got)
		}
	}

	{
		got := binarySearchKline(arr, time.Unix(1645779988, 0))
		want := 0.6
		if got != want {
			t.Fatalf("Want %f, Got %f", want, got)
		}
	}
}

func Test_batch(t *testing.T) {
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

	start, err := time.Parse(time.RFC3339, "2022-01-04T00:50:10.770Z")
	if err != nil {
		t.Fatal(err)
	}

	end, err := time.Parse(time.RFC3339, "2022-01-04T00:55:10.770Z")
	if err != nil {
		t.Fatal(err)
	}

	err = batch(start, end)
	if err != nil {
		t.Fatal(err)
	}
}

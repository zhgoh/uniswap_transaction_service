package main

import (
	"testing"
	"time"
)

func Test_getKlines_default(t *testing.T) {
	// Test binance getkline function to see if it can pull data using the api
	client := makeBinanceClient()
	got, err := client.getKlines("ETHUSDT", 1, Days, time.Time{}, time.Time{}, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) == 0 || got[0].OpenTime == 0 {
		t.Fatal("Error: No results from kline api.")
	}
}

func Test_getKlines_start(t *testing.T) {
	// Test binance getkline function with start time
	client := makeBinanceClient()
	got, err := client.getKlines("ETHUSDT", 1, Days, time.UnixMilli(1502928000000), time.Time{}, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) == 0 || got[0].OpenTime == 0 {
		t.Fatal("Error: No results from kline api.")
	}

	{
		var want int64 = 1502928000000
		if got[0].OpenTime != want {
			t.Errorf("got %q, wanted %q", got, want)
		}
	}

	{
		want := "301.13000000"
		if got[0].Open != want {
			t.Errorf("got %q, wanted %q", got, want)
		}
	}
}

func Test_getKlines_end(t *testing.T) {
	// Test binance getkline function with end time
	client := makeBinanceClient()
	got, err := client.getKlines("ETHUSDT", 1, Days, time.Time{}, time.UnixMilli(1502928000000), 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) == 0 || got[0].OpenTime == 0 {
		t.Fatal("Error: No results from kline api.")
	}

	{
		var want int64 = 1502928000000
		if got[len(got)-1].OpenTime != want {
			t.Errorf("got %q, wanted %q", got, want)
		}
	}

	{
		want := "301.13000000"
		if got[len(got)-1].Open != want {
			t.Errorf("got %q, wanted %q", got, want)
		}
	}
}

func Test_getOrderBooks(t *testing.T) {
	// Test binance getkline function with end time
	client := makeBinanceClient()
	got, err := client.getOrderBook("ETHUSDT", 1)
	if err != nil {
		t.Fatal(err)
	}

	if almostEqual(0.0, got) {
		t.Fatal("Error: No results from orderbook api.")
	}
}

func Test_getSymbolPrice(t *testing.T) {
	// Test binance get symbol price function
	client := makeBinanceClient()
	got, err := client.getSymbolPrice("ETHUSDT")
	if err != nil {
		t.Fatal(err)
	}

	if almostEqual(0.0, got) {
		t.Fatal("Error: No results from price api.")
	}
}

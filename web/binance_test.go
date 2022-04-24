package main

import (
	"math"
	"testing"
)

func Test_getKlines_default(t *testing.T) {
	// Test binance getkline function to see if it can pull data using the api
	client := makeBinanceClient()
	got, err := client.getKlines("ETHUSDT", 1, Days, 0, 0, 0)
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
	got, err := client.getKlines("ETHUSDT", 1, Days, 1502928000000, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) == 0 || got[0].OpenTime == 0 {
		t.Fatal("Error: No results from kline api.")
	}

	{
		want := 1502928000000
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
	got, err := client.getKlines("ETHUSDT", 1, Days, 0, 1502928000000, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) == 0 || got[0].OpenTime == 0 {
		t.Fatal("Error: No results from kline api.")
	}

	{
		want := 1502928000000
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

const float64EqualityThreshold = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
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
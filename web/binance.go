package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type binanceClient struct{}

func makeBinanceClient() *binanceClient {
	return &binanceClient{}
}

type chartInterval int64

const (
	Minutes chartInterval = iota
	Hours
	Days
	Weeks
	Months
)

func (chart chartInterval) String() string {
	switch chart {
	case Minutes:
		return "m"
	case Hours:
		return "h"
	case Days:
		return "d"
	case Weeks:
		return "w"
	case Months:
		return "M"
	}

	log.Print("Error: Unknown chart interval.")
	return "unknown"
}

// Representation of the kline response
type klineResponse struct {
	OpenTime              int64
	Open                  string
	High                  string
	Low                   string
	Close                 string
	Volume                string
	CloseTime             int64
	QuoteAssetVol         string
	NumTrades             int
	TakerBuyBaseAssetVol  string
	TakerBuyQuoteAssetVol string
	Ignore                string
}

// Custom unmarshal json to support unmarshalling arrays
func (k *klineResponse) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{
		&k.OpenTime,
		&k.Open,
		&k.High,
		&k.Low,
		&k.Close,
		&k.Volume,
		&k.CloseTime,
		&k.QuoteAssetVol,
		&k.NumTrades,
		&k.TakerBuyBaseAssetVol,
		&k.TakerBuyQuoteAssetVol,
		&k.Ignore,
	}

	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("wrong number of fields in KlineResponse: %d != %d", g, e)
	}
	return nil
}

func (client *binanceClient) getKlines(symbol string, freq int, interval chartInterval, startTime, endTime time.Time, limit int64) ([]klineResponse, error) {
	queries := make(map[string]int64)
	if !startTime.IsZero() {
		queries["startTime"] = startTime.UnixMilli()
	}

	if !endTime.IsZero() {
		queries["endTime"] = endTime.UnixMilli()
	}

	if limit > 0 {
		queries["limit"] = limit
	}

	// Get the kline api and unmarshal using custom func to our klineResponse struct
	api := fmt.Sprintf("https://api.binance.com/api/v3/klines?symbol=%s&interval=%d%s", symbol, freq, interval)
	for k, v := range queries {
		api = fmt.Sprintf("%s&%s=%d", api, k, v)
	}

	resp, err := http.Get(api)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var klineResp []klineResponse
	if err := json.Unmarshal(body, &klineResp); err != nil {
		return nil, err
	}
	return klineResp, nil
}

// Representation of the order book response
type orderBookResponse struct {
	LastUpdateId int
	Bids         [][]string
	Asks         [][]string
}

func (client *binanceClient) getOrderBook(symbol string, limit int) (float64, error) {
	if limit < 1 {
		log.Print("Error: invalid limit for order books, defaulting to 1")
		limit = 1
	}

	// Get the kline api and unmarshal using custom func to our klineResponse struct
	api := fmt.Sprintf("https://api.binance.com/api/v3/depth?symbol=%s&limit=%d", symbol, limit)
	resp, err := http.Get(api)
	if err != nil {
		return 0.0, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0.0, err
	}

	var orderBookResp orderBookResponse
	if err := json.Unmarshal(body, &orderBookResp); err != nil {
		return 0.0, err
	}

	// TODO: Find a way to compute better pricing
	price, err := strconv.ParseFloat(orderBookResp.Bids[0][0], 64)
	if err != nil {
		return 0.0, err
	}

	return price, nil
}

func (client *binanceClient) getSymbolPrice(symbol string) (float64, error) {
	// Get the kline api and unmarshal using custom func to our klineResponse struct
	api := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", symbol)
	resp, err := http.Get(api)
	if err != nil {
		return 0.0, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0.0, err
	}

	var priceResp map[string]string
	if err := json.Unmarshal(body, &priceResp); err != nil {
		return 0.0, err
	}

	price, err := strconv.ParseFloat(priceResp["price"], 64)
	if err != nil {
		return 0.0, err
	}
	return price, nil
}

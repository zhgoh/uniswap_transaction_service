package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type ChartInterval int64

const (
	Minutes ChartInterval = iota
	Hours
	Days
	Weeks
	Months
)

func (chart ChartInterval) String() string {
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

type binanceClient struct{}

func makeBinanceClient() *binanceClient {
	return &binanceClient{}
}

// Representation of the kline response
type KlineResponse struct {
	OpenTime              int
	Open                  string
	High                  string
	Low                   string
	Close                 string
	Volume                string
	CloseTime             int
	QuoteAssetVol         string
	NumTrades             int
	TakerBuyBaseAssetVol  string
	TakerBuyQuoteAssetVol string
	Ignore                string
}

// Custom unmarshal json to support unmarshalling arrays
func (k *KlineResponse) UnmarshalJSON(buf []byte) error {
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

func (client *binanceClient) getKlines(symbol string, freq int, interval ChartInterval, startTime, endTime, limit int64) ([]KlineResponse, error) {
	queries := make(map[string]int64)
	if startTime > 0 {
		queries["startTime"] = startTime
	}
	if endTime > 0 {
		queries["endTime"] = endTime
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var klineResp []KlineResponse
	if err := json.Unmarshal(body, &klineResp); err != nil {
		return nil, err
	}
	return klineResp, nil
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type etherscanClient struct {
	apiKey string
}

func makeEtherscan() (*etherscanClient, error) {
	apiKey := os.Getenv("etherscan_api")
	if len(apiKey) == 0 {
		log.Print("Error: API key etherscan_api not found in env variables.")
		return nil, fmt.Errorf("error: API key etherscan_api not found in env variables")
	}
	return &etherscanClient{apiKey: apiKey}, nil
}

type etherscanTransactionResponse struct {
	Status  string
	Message string
	Result  []etherscanTransaction
}

type etherscanTransaction struct {
	BlockNumber       string
	TimeStamp         string
	Hash              string
	Nonce             string
	BlockHash         string
	From              string
	ContractAddress   string
	To                string
	Value             string
	TokenName         string
	TokenSymbol       string
	TokenDecimal      string
	TransactionIndex  string
	Gas               string
	GasPrice          string
	GasUsed           string
	CumulativeGasUsed string
	Input             string
	Confirmations     string
}

func (client *etherscanClient) fetchTransactions(startBlock, endBlock int) ([]etherscanTransaction, error) {
	queries := map[string]string{
		"module":          "account",
		"action":          "tokentx",
		"contractaddress": "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		"address":         "0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640",
		"page":            "1",
		"offset":          "100",
		"sort":            "desc",
	}

	if startBlock > 0 {
		queries["startblock"] = fmt.Sprint(startBlock)
	}

	if endBlock > 0 {
		queries["endblock"] = fmt.Sprint(endBlock)
	}

	api := fmt.Sprintf("https://api.etherscan.io/api?apikey=%s", client.apiKey)
	for k, v := range queries {
		api = fmt.Sprintf("%s&%s=%s", api, k, v)
	}

	resp, err := http.Get(api)
	if err != nil {
		log.Print("Error: Unable to fetch transactions.")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print("Error: could not read body of fetched transactions.")
		return nil, err
	}

	var transactionResp etherscanTransactionResponse
	err = json.Unmarshal(body, &transactionResp)
	if err != nil {
		log.Print("Error: could not unmarshal etherscan results.")
		return nil, err
	}

	if transactionResp.Message != "OK" {
		log.Print("Error: problem with getting transactions.")
		log.Printf("Error: %s.\n", transactionResp.Message)
		return nil, err
	}

	//log.Print("Successfully get transactions.")
	//log.Printf("Got %d entries.", len(transactionResp.Result))
	return transactionResp.Result, nil
}

type closestBlock int64

const (
	before closestBlock = iota
	after
)

func (closest closestBlock) String() string {
	switch closest {
	case before:
		return "before"
	case after:
		return "after"
	}

	log.Print("Error: unknown closest block.")
	return "unknown"
}

type etherscanBlockResponse struct {
	Status  string
	Message string
	Result  string
}

// getBlockNumberByTimestamp will call the etherscan's api to fetch closest blocks
func (client *etherscanClient) getBlockNumberByTimestamp(closest closestBlock, timestamp time.Time) (int, error) {
	queries := map[string]string{
		"module":    "block",
		"action":    "getblocknobytime",
		"timestamp": fmt.Sprint(timestamp.Unix()),
		"closest":   fmt.Sprint(closest),
	}
	api := fmt.Sprintf("https://api.etherscan.io/api?apikey=%s", client.apiKey)
	for k, v := range queries {
		api = fmt.Sprintf("%s&%s=%s", api, k, v)
	}

	resp, err := http.Get(api)
	if err != nil {
		log.Print("Error: Unable to fetch block info from etherscan.")
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print("Error: could not read body of fetched block response.")
		return 0, err
	}

	var blockResp etherscanBlockResponse
	err = json.Unmarshal(body, &blockResp)
	if err != nil {
		log.Print("Error: could not unmarshal block response results.")
		return 0, err
	}

	if blockResp.Message != "OK" {
		log.Print("Error: problem with getting block information.")
		log.Printf("Error: %s\n", blockResp.Message)
		return 0, err
	}

	// log.Print("Successfully get block information.")
	// log.Printf("Got %d entries.", len(blockResp.Result))

	res, err := strconv.Atoi(blockResp.Result)
	if err != nil {
		log.Print("Error: problem with string conversion.")
		return 0, err
	}
	return res, nil
}

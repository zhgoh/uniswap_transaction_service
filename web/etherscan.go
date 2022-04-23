package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type EtherscanAPI struct {
	Status  string
	Message string
	Result  []EtherscanTransaction
}

type EtherscanTransaction struct {
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

type EtherscanClient struct {
	apiKey string
}

//func (client *EtherscanClient) fetchDatedTransactions() ([]EtherscanTransaction, error) {
//}

func (client *EtherscanClient) fetchTransactions() ([]EtherscanTransaction, error) {
	queries := map[string]string{
		"action":          "tokentx",
		"module":          "account",
		"contractaddress": "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		"address":         "0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640",
		"page":            "1",
		"offset":          "100",
		"sort":            "desc",
		//"startblock":      "0",
		//"endblock":        "0",
	}
	url := fmt.Sprintf("https://api.etherscan.io/api?apikey=%s", client.apiKey)
	for k, v := range queries {
		url = fmt.Sprintf("%s&%s=%s", url, k, v)
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Print("Error: Unable to fetch transactions.")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("Error: could not read body of fetched transactions.")
		return nil, err
	}

	var etherScan EtherscanAPI
	err = json.Unmarshal(body, &etherScan)
	if err != nil {
		log.Print("Error: could not unmarshal etherscan results.")
		return nil, err
	}

	if etherScan.Message != "OK" {
		log.Print("Error: problem with getting transactions.")
		log.Printf("Error: %s\n", etherScan.Message)
		return nil, err
	}

	log.Print("Successfully get transactions")
	log.Printf("Got %d entries", len(etherScan.Result))
	return etherScan.Result, nil
}

func makeEtherscan() (*EtherscanClient, error) {
	apiKey := os.Getenv("etherscan_api")
	if len(apiKey) == 0 {
		log.Print("Error: API key etherscan_api not found in env variables.")
		return nil, fmt.Errorf("Error: API key etherscan_api not found in env variables.")
	}
	return &EtherscanClient{apiKey: apiKey}, nil
}

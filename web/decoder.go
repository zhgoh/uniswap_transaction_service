package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"example.com/backend/v2/uniswap"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type swapAmounts struct {
	usdc *big.Float
	eth  *big.Float
}

func decodeTransaction(transaction_id string) ([]swapAmounts, error) {

	infura_node := os.Getenv("INFURA_NODE")
	if infura_node == "" {
		return nil, fmt.Errorf("error getting INFURA NODE env var")
	}

	client, err := ethclient.Dial(infura_node)
	if err != nil {
		log.Print("Failed to connect to infura.")
		return nil, err
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(uniswap.UniswapMetaData.ABI)))
	if err != nil {
		log.Print("Unable to create ABI ", err.Error())
		return nil, err
	}

	ctx := context.Background()

	txHash := common.HexToHash(transaction_id)
	receipts, err := client.TransactionReceipt(ctx, txHash)
	if err != nil {
		log.Print("Error: failed to process transaction receipt.")
		return nil, err
	}

	contractAddress := common.HexToAddress("0x88e6A0c2dDD26FEEb64F039a2c41296FcB3f5640")
	query := ethereum.FilterQuery{
		FromBlock: receipts.BlockNumber,
		ToBlock:   receipts.BlockNumber,
		Addresses: []common.Address{
			contractAddress,
		},
	}

	logs, err := client.FilterLogs(ctx, query)
	if err != nil {
		log.Print("Unable to filter Logs ", err.Error())
		return nil, err
	}

	var swapEvent uniswap.UniswapSwap

	var results []swapAmounts
	for _, vLog := range logs {
		err := contractAbi.UnpackIntoInterface(&swapEvent, "Swap", vLog.Data)
		if err != nil {
			continue
		}
		amount0 := big.NewFloat(0).SetInt(swapEvent.Amount0)
		amount1 := big.NewFloat(0).SetInt(swapEvent.Amount1)
		results = append(results, swapAmounts{amount0, amount1})
	}
	return results, nil
}

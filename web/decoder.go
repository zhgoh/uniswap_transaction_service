package main

import (
	"context"
	"fmt"
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
		return nil, err
	}

	txHash := common.HexToHash(transaction_id)
	receipts, err := client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

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
		return nil, err
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(uniswap.UniswapABI)))
	if err != nil {
		return nil, err
	}

	var swapEvent uniswap.UniswapSwap

	var results []swapAmounts
	for _, vLog := range logs {
		err := contractAbi.UnpackIntoInterface(&swapEvent, "Swap", vLog.Data)
		if err != nil {
			return nil, err
		}
		amount0 := big.NewFloat(0).SetInt(swapEvent.Amount0)
		amount1 := big.NewFloat(0).SetInt(swapEvent.Amount1)
		results = append(results, swapAmounts{amount0, amount1})

	}
	return results, nil
}

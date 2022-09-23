package main

import (
	"context"
	"log"
	"math/big"
	"os"
	"strings"
	"testing"

	"example.com/backend/v2/uniswap"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Test_decodetransaction(t *testing.T) {
	infura_node := os.Getenv("INFURA_NODE")
	if infura_node == "" {
		log.Fatal("Cannot find INFURA NODE")
	}

	client, err := ethclient.Dial(infura_node)
	if err != nil {
		log.Fatal(err)
	}

	// get block
	txHash := common.HexToHash("0x87339c2c394f132069213466e671aba46e3f9d22287ea8d16cf80d1884855321")

	receipts, err := client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		log.Fatal("Cannot get transaction receipt")
	}
	log.Print(receipts.BlockNumber)

	contractAddress := common.HexToAddress("0x88e6A0c2dDD26FEEb64F039a2c41296FcB3f5640")
	query := ethereum.FilterQuery{
		FromBlock: receipts.BlockNumber,
		ToBlock:   receipts.BlockNumber,
		Addresses: []common.Address{
			contractAddress,
		},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(uniswap.UniswapABI)))
	if err != nil {
		log.Fatal(err)
	}

	var swapEvent struct {
		Sender       common.Address
		Recipient    common.Address
		Amount0      *big.Int
		Amount1      *big.Int
		SqrtPriceX96 *big.Int
		Liquidity    *big.Int
		Tick         *big.Int
	}

	for _, vLog := range logs {
		err := contractAbi.UnpackIntoInterface(&swapEvent, "Swap", vLog.Data)
		if err != nil {
			// log.Fatal(err)
			log.Print("Skipping")
			continue
		}

		log.Print(swapEvent)

		amount0 := big.NewFloat(0).SetInt(swapEvent.Amount0)
		amount1 := big.NewFloat(0).SetInt(swapEvent.Amount1)

		log.Print(amount0, amount1)
	}
}

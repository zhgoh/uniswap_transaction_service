package main

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"

	"database/sql"

	_ "github.com/glebarez/go-sqlite"
)

type dbClient struct {
	db *sql.DB
}

func makeDBClient(fileName string) (*dbClient, error) {
	db, err := sql.Open("sqlite", fileName)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	client := &dbClient{db}

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		// Create new file
		file, err := os.Create(fileName)
		if err != nil {
			log.Fatal(err.Error())
			return nil, err
		}
		file.Close()

		log.Print("Creating DB")
		createTableStmt := `CREATE TABLE Transactions(
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"transaction_id" TEXT NOT NULL,
		"usdc" REAL NOT NULL,
		"eth" REAL NOT NULL,
		"fee" REAL NOT NULL
		);` // SQL Statement for Create Table

		if _, err := client.execStatement(createTableStmt); err != nil {
			log.Fatal(err.Error())
			return nil, err
		}
	}
	return client, nil
}

func (client *dbClient) execStatement(stmt string) (sql.Result, error) {
	statement, err := client.db.Prepare(stmt)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	res, err := statement.Exec() // Execute SQL Statements
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	return res, nil
}

func (client *dbClient) addTransaction(transaction cryptoTransaction) error {
	statement := fmt.Sprintf(
		"INSERT INTO Transactions (transaction_id, usdc, eth, fee) VALUES (\"%s\", \"%f\", \"%f\", \"%f\");",
		transaction.Hash,
		transaction.USDC,
		transaction.ETH,
		transaction.Fee)

	_, err := client.execStatement(statement)
	return err
}

func (client *dbClient) clearTable() error {
	statement := "Delete FROM Transactions"
	_, err := client.execStatement(statement)
	return err
}

func (client *dbClient) getTransaction(hash string) (*cryptoTransaction, error) {
	query := fmt.Sprintf("SELECT transaction_id, usdc, eth, fee from Transactions where transaction_id=\"%s\";", hash)
	row, err := client.db.Query(query)
	if err != nil {
		log.Print("Error: querying DB.")
		return nil, err
	}

	defer row.Close()

	var fee float64
	var usdc float64
	var eth float64

	var res *cryptoTransaction
	for row.Next() {
		if err = row.Scan(&hash, &usdc, &eth, &fee); err != nil {
			log.Print("Error: getting row info ", err.Error())
			continue
		}
		res = &cryptoTransaction{hash, big.NewFloat(usdc), big.NewFloat(eth), fee}
	}
	return res, nil
}

func (client *dbClient) getAllTransactions() ([]cryptoTransaction, error) {
	res := []cryptoTransaction{}
	query := "SELECT * FROM Transactions;"
	row, err := client.db.Query(query)
	if err != nil {
		return res, err
	}

	defer row.Close()

	var fee float64
	var hash string
	var usdc *big.Float
	var eth *big.Float

	for row.Next() {
		if err = row.Scan(&hash, &usdc, &eth, &fee); err != nil {
			log.Print("Error: getting row info")
			continue
		}
		res = append(res, cryptoTransaction{hash, usdc, eth, fee})
	}
	return res, nil
}

func (client *dbClient) addLiveTransactions(etherTransactions []etherscanTransaction, prices float64) error {
	if len(etherTransactions) == 0 {
		return fmt.Errorf("no transactions provided")
	}

	for _, v := range etherTransactions {
		if len(v.Hash) == 0 {
			return fmt.Errorf("hash is empty")
		}

		if v.Hash == latestHash {
			break
		}

		err := client.addSingleTransaction(v, prices)
		if err != nil {
			log.Print("Failed to add single transaction ", err.Error())
			return err
		}

	}
	latestHash = etherTransactions[0].Hash
	return nil
}

func (client *dbClient) addSingleTransaction(transaction etherscanTransaction, prices float64) error {
	res, err := db.getTransaction(transaction.Hash)
	if err != nil {
		log.Print("Failed to get transaction ", err.Error())
		return err
	}

	if res != nil {
		return nil
	}

	// Compute prices
	gasPrice, err := strconv.Atoi(transaction.GasPrice)
	if err != nil {
		log.Print("Error: failed to convert gas price to integer.")
		return err
	}

	gasUsed, err := strconv.Atoi(transaction.GasUsed)
	if err != nil {
		log.Print("Error: failed to convert gas used to integer.")
		return err
	}

	// Fees in eth
	// Note: no idea if division or multiplying would be faster here, probably same
	// fees := float64(gasPrice*gasUsed) / 1000000000000000000
	fees := float64(gasPrice*gasUsed) * 0.000000000000000001
	fees *= prices

	// Convert to price in USDT

	swapAmounts, err := decodeTransaction(transaction.Hash)
	if err != nil {
		log.Print("Error: failed to decode transaction ", err.Error())
		return err
	}

	for _, value := range swapAmounts {
		db.addTransaction(cryptoTransaction{transaction.Hash, value.usdc, value.eth, fees})
	}

	timeStamp, err := strconv.Atoi(transaction.TimeStamp)
	if err != nil {
		log.Print("Error: failed to convert timeStamp.")
		return err
	}

	log.Printf("Added Hash: %s, Time: %d, Fees: $%.2f", transaction.Hash, timeStamp, fees)
	return nil
}

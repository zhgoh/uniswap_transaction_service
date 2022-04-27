package main

import (
	"fmt"
	"log"
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

		// log.Print("Creating DB")
		createTableStmt := `CREATE TABLE Transactions(
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"hash" TEXT NOT NULL UNIQUE,
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
		"INSERT INTO Transactions (hash, fee) VALUES (\"%s\", %f);",
		transaction.Hash,
		transaction.Fee)

	_, err := client.execStatement(statement)
	return err
}

func (client *dbClient) clearTable() error {
	statement := fmt.Sprint("Delete FROM Transactions")
	_, err := client.execStatement(statement)
	return err
}

func (client *dbClient) getTransaction(hash string) (*cryptoTransaction, error) {
	query := fmt.Sprintf("SELECT fee from Transactions where hash=\"%s\";", hash)
	row, err := client.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	var fee float64
	var res *cryptoTransaction
	for row.Next() {
		if err = row.Scan(&fee); err != nil {
			log.Print("Error: getting row info")
			continue
		}
		res = &cryptoTransaction{hash, fee}
	}
	return res, nil
}

func (client *dbClient) getAllTransactions() ([]cryptoTransaction, error) {
	res := []cryptoTransaction{}
	query := fmt.Sprintf("SELECT hash, fee FROM Transactions;")
	row, err := client.db.Query(query)
	if err != nil {
		return res, err
	}

	defer row.Close()

	var fee float64
	var hash string
	for row.Next() {
		if err = row.Scan(&hash, &fee); err != nil {
			log.Print("Error: getting row info")
			continue
		}
		res = append(res, cryptoTransaction{hash, fee})
	}
	return res, nil
}

func (client *dbClient) addLiveTransactions(etherTransactions []etherscanTransaction, prices float64) error {
	if len(etherTransactions) == 0 {
		return fmt.Errorf("no transactions provided")
	}

	for _, v := range etherTransactions {
		if len(v.Hash) == 0 {
			return fmt.Errorf("hash is empty.")
		}

		if v.Hash == latestHash {
			break
		}

		err := client.addSingleTransaction(v, prices)
		if err != nil {
			return err
		}

	}
	latestHash = etherTransactions[0].Hash
	return nil
}

func (client *dbClient) addSingleTransaction(transaction etherscanTransaction, prices float64) error {
	res, err := db.getTransaction(transaction.Hash)
	if err != nil {
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
	db.addTransaction(cryptoTransaction{transaction.Hash, fees})

	timeStamp, err := strconv.Atoi(transaction.TimeStamp)
	if err != nil {
		log.Print("Error: failed to convert timeStamp.")
		return err
	}

	// TODO: Add to DB
	log.Printf("Hash: %s, Time: %d, Fees: $%.2f", transaction.Hash, timeStamp, fees)
	return nil
}

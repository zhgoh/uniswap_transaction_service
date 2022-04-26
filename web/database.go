package main

import (
	"fmt"
	"log"
	"os"

	"database/sql"

	_ "github.com/glebarez/go-sqlite"
)

type DBClient struct {
	db *sql.DB
}

func makeDBClient(fileName string, creationStatment string) (*DBClient, error) {
	db, err := sql.Open("sqlite", fileName)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	client := &DBClient{db}

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		// Create new file
		file, err := os.Create(fileName)
		if err != nil {
			log.Fatal(err.Error())
			return nil, err
		}
		file.Close()

		log.Print("Creating DB")
		if _, err := client.execStatement(creationStatment); err != nil {
			log.Fatal(err.Error())
			return nil, err
		}
	}
	return client, nil
}

func (client *DBClient) execStatement(stmt string) (sql.Result, error) {
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

func (client *DBClient) addTransactions(transaction cryptoTransaction) error {
	statement := fmt.Sprintf(
		"INSERT INTO Transactions (hash, fee) VALUES (\"%s\", %f);",
		transaction.Hash,
		transaction.Fee)

	_, err := client.execStatement(statement)
	return err
}

func (client *DBClient) getTransactions(hash string) (float64, error) {
	query := fmt.Sprintf("SELECT fee from Transactions where hash=\"%s\"", hash)
	row, err := client.db.Query(query)
	if err != nil {
		return 0, err
	}

	defer row.Close()

	var fee float64
	for row.Next() {
		row.Scan(&fee)
	}
	return fee, nil
}

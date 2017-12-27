package main

import (
	"regexp"
	"strconv"
	_ "github.com/lib/pq"
	"fmt"
	"github.com/ybbus/jsonrpc"
	"time"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Validates request from the TCP client
func validateData(from, to, amount string) (result []string, status bool) {
	var hexRegex = regexp.MustCompile("0[xX][0-9a-fA-F]+")
	status = false

	if hexRegex.Find([]byte(from)) == nil {
		result = append(result, "Sender has a wrong wallet")
	}

	if hexRegex.Find([]byte(to)) == nil {
		result = append(result, "Recipient has a wrong wallet")
	}

	if amountConverted, err := strconv.ParseFloat(amount, 64); err == nil {
		if amountConverted <= 0 {
			result = append(result, "Amount is zero or lower than")
		}
	} else {
		result = append(result, "Wrong amount")
	}

	if len(result) == 0 {
		status = true
	}

	return
}

// Check transaction delivery in goroutine
func CatchDeliveryTime(transactionHash string) {
	for {
		rpcClient := jsonrpc.NewRPCClient(fmt.Sprintf("http://%s:%s", RPCHOST, RPCPORT))
		response, _ := rpcClient.Call("eth_getTransactionByHash", transactionHash)
		if response.Result.(map[string]interface{})["blockNumber"] != nil {
			stmt, err := db.Prepare("UPDATE transactions SET date=now(), completed = TRUE where transaction_hash=$1")
			checkErr(err)
			res, err := stmt.Exec(transactionHash)
			checkErr(err)
			_, err = res.RowsAffected()
			checkErr(err)

			if to := response.Result.(map[string]interface{})["to"]; to != nil {
				MakeBalanceRecord(transactionHash, to.(string))
			}

			return
		}

		// Just delay, goroutine
		time.Sleep(time.Millisecond * 200)
	}
}

func Store(from, to, transaction string, amount int) (message string, isStored bool) {
	var lastInsertId int

	err := db.QueryRow("INSERT INTO transactions(sender,recipient,amount,transaction_hash) VALUES($1,$2,$3,$4) returning id;", from, to, amount, transaction).Scan(&lastInsertId)
	checkErr(err)

	if lastInsertId > 0 {
		fmt.Println("Entry created")
		return "Success", true
	}

	return "Couldn't create record", false
}

func GetLast(c chan bool) (lastTransactions []LastTransactions) {
	rpcClient := jsonrpc.NewRPCClient(fmt.Sprintf("http://%s:%s", RPCHOST, RPCPORT))
	responseHeight, _ := rpcClient.Call("eth_blockNumber")

	bHeight, err := strconv.ParseUint(responseHeight.Result.(string)[2:], 16, 32)
	checkErr(err)
	rows, err := db.Query("SELECT date, recipient, amount, show, transaction_hash FROM transactions WHERE show = FALSE and date NOTNULL")
	checkErr(err)

	for rows.Next() {
		var t LastTransactions
		err := rows.Scan(&t.Date, &t.Recipient, &t.Amount, &t.Show, &t.TransactionHash)
		transactionBlock := GetTransactionBlockByHash(t.TransactionHash)
		numOfConfirmations := bHeight - transactionBlock + 1
		checkErr(err)
		t.Confirmations = numOfConfirmations

		if t.Show == false || numOfConfirmations < 3 {
			lastTransactions = append(lastTransactions, t)
			if !t.Show && numOfConfirmations >= 3 {
				MarkAsShown(t.TransactionHash)
			}
		} else {
			continue
		}
	}

	c <- true

	return
}

func MarkAsShown(transactionHash string) {
	stmt, err := db.Prepare("UPDATE transactions SET show = TRUE where transaction_hash=$1")
	checkErr(err)
	_, err = stmt.Exec(transactionHash)
	checkErr(err)
}

func GetTransactionBlockByHash(transactionHash string) (uint64) {
	rpcClient := jsonrpc.NewRPCClient(fmt.Sprintf("http://%s:%s", RPCHOST, RPCPORT))
	response, _ := rpcClient.Call("eth_getTransactionByHash", transactionHash)
	if blockNumber := response.Result.(map[string]interface{})["blockNumber"]; blockNumber != nil {
		result, err := strconv.ParseUint(blockNumber.(string)[2:], 16, 32)
		checkErr(err)
		return result
	}

	return 0
}

func CheckTransactions() {
	rows, err := db.Query("SELECT transaction_hash FROM transactions WHERE completed = FALSE")
	checkErr(err)
	for rows.Next() {
		var transactionHash string
		err = rows.Scan(&transactionHash)
		checkErr(err)

		CatchDeliveryTime(transactionHash)
	}

	return
}

func MakeBalanceRecord(transactionHash string, recipient string) {
	rpcClient := jsonrpc.NewRPCClient(fmt.Sprintf("http://%s:%s", RPCHOST, RPCPORT))
	response, _ := rpcClient.Call("eth_getBalance", recipient, "latest")
	balance := response.Result.(string)

	stmt, err := db.Prepare("INSERT INTO balance(wallet, balance) VALUES($1, $2)")
	checkErr(err)
	_, err = stmt.Exec(transactionHash, string(balance))
	checkErr(err)
}

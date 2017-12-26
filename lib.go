package main

import (
	"regexp"
	"strconv"
	_ "github.com/lib/pq"
	"fmt"
	"time"
	"github.com/ybbus/jsonrpc"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

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
			return
		}
		time.Sleep(time.Millisecond * 50)
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
	rows, err := db.Query("SELECT date, recipient, amount, shown, shown_count, transaction_hash FROM transactions WHERE date NOTNULL")
	checkErr(err)

	for rows.Next() {
		var t LastTransactions
		err := rows.Scan(&t.Date, &t.Recipient, &t.Amount, &t.Shown, &t.ShownCount, &t.TransactionHash)
		checkErr(err)
		transactionBlock := GetTransactionBlockByHash(t.TransactionHash)

		numOfConfirmations := bHeight - transactionBlock + 1

		t.ShownCount.Int64 += 1
		count := t.ShownCount.Int64
		IncrementShownCount(t.TransactionHash, count)

		if t.Shown == false || numOfConfirmations < 3 {
			lastTransactions = append(lastTransactions, t)
			if !t.ShownCount.Valid {
				MarkAsShown(t.TransactionHash)
			}
		} else {
			continue
		}
	}

	c <- true

	return
}

func IncrementShownCount(transactionHash string, count int64) {
	stmt, err := db.Prepare("UPDATE transactions SET shown_count = $1 where transaction_hash=$2")
	checkErr(err)
	_, err = stmt.Exec(count, transactionHash)
	checkErr(err)
}

func MarkAsShown(transactionHash string) {
	stmt, err := db.Prepare("UPDATE transactions SET shown = TRUE where transaction_hash=$1")
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

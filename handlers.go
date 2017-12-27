package main

import (
	"net"
	"strings"
	"strconv"
	"fmt"
	"encoding/json"
)

func SendHandler(action string, data []string, conn net.Conn) {
	defer conn.Close()
	validatedResult, status := validateData(data[1], data[2], data[3])
	if !status {
		conn.Write([]byte("Please, check your data: " + strings.Join(validatedResult, "; ") ))
		conn.Close()
		return
	}

	password := data[len(data)-1]

	hexAmount, err := strconv.Atoi(data[3])
	checkErr(err)

	result, status := SendEth(data[1], data[2], fmt.Sprintf("0x%X", hexAmount), password)

	if status {
		res := []byte(fmt.Sprintf("Success: %s", action))
		conn.Write(res)
		message, isStored := Store(data[1], data[2], result, hexAmount)
		fmt.Println(message)

		go CatchDeliveryTime(result)

		if !isStored {
			conn.Write(res)
		}
	} else {
		conn.Write([]byte("Transaction error: " + result))
	}
}

func GetLastHandler(err error, conn net.Conn) {
	var lastTransactions []LastTransactions
	c := make(chan bool)

	go func() {
		lastTransactions = GetLast(c)
	}()
	<-c

	checkErr(err)
	if len(lastTransactions) > 0 {
		marshaledTransactions, err := json.Marshal(lastTransactions)
		checkErr(err)
		conn.Write(marshaledTransactions)
	} else {
		conn.Write([]byte("I have no recent transactions"))
	}
	conn.Close()
}

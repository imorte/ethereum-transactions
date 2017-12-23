package main

import (
	"fmt"
	"encoding/json"
)

func SendEth() {

}

func GetBalance() {
	marshaledResult, err := json.Marshal(Balance{
		Jsonrpc: "2.0",
		Method: "eth_getBalance",
		Params: BalanceParams{
			Data: "0x532bce52569bd8181fc1cadeeb18c2ae4e58cf0c",
			Quantity: "latest",
		},
		Id: 1,
	})
	checkErr(err)

	fmt.Println(marshaledResult)

	fmt.Println("balance is ")
}
package main

import (
	"github.com/ybbus/jsonrpc"
	"fmt"
)

func SendEth(from string, to string, amount string, password string) (string, bool) {
	rpcClient := jsonrpc.NewRPCClient(fmt.Sprintf("http://%s:%s", RPCHOST, RPCPORT))

	response, _ := rpcClient.Call("personal_sendTransaction", Transaction{from, to, amount}, password)

	transaction := response.Result.(string)

	if response.Error != nil {
		return fmt.Sprintf("An error occurred: %s", response.Error), false
	} else {
		return transaction, true
	}
}

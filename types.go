package main

type Balance struct {
	Jsonrpc string
	Method string
	Params BalanceParams
	Id int
}

type BalanceParams struct {
	Data string
	Quantity string
}
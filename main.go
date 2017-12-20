package main

import (
	_ "net/rpc/jsonrpc"
	_ "github.com/buger/jsonparser"
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
)

const (
	DBUSER     = "postgres"
	DBPASSWORD = "1111"
	DBNAME     = "etherium"
)

func main() {
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DBUSER, DBPASSWORD, DBNAME)

	db, err := sql.Open("postgres", dbInfo)
	checkErr(err)
	defer db.Close()

	fmt.Println("Established connection with db...")


}


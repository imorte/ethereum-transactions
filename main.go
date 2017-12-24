package main

import (
	_ "net/rpc/jsonrpc"
	_ "github.com/buger/jsonparser"
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"runtime"
)

var db *sql.DB

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var err error

	dbConnect := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DBUSER, DBPASSWORD, DBNAME)

	db, err = sql.Open("postgres", dbConnect)
	checkErr(err)

	defer db.Close()

	fmt.Println("Established connection with db...")

	ListenTcp()
}

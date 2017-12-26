package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
)

var db *sql.DB

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

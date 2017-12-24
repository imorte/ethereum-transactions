package main

import (
	"regexp"
	"strconv"
	_ "github.com/lib/pq"
	"fmt"
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

func Store(from, to, transaction string, amount int) (message string, isStored bool) {
	var lastInsertId int

	err := db.QueryRow("INSERT INTO transactions(sender,recipient,amount,transaction_label) VALUES($1,$2,$3,$4) returning id;", from, to, amount, transaction).Scan(&lastInsertId)
	checkErr(err)

	if lastInsertId > 0 {
		fmt.Println("Entry created")
		return "Success", true
	}

	return "Couldn't create record", false
}

func GetLast() {

}
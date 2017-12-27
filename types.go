package main

import (
	"database/sql"
	"encoding/json"
)

type Transaction struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}

type LastTransactions struct {
	Date            JsonNullString `json:"date"`
	Recipient       string         `json:"recipient"`
	Amount          int            `json:"amount"`
	Show            bool           `json:"-"`
	TransactionHash string         `json:"-"`
	Confirmations   uint64         `json:"confirmations"`
}

type JsonNullString struct {
	sql.NullString
}

func (v JsonNullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	} else {
		return json.Marshal(nil)
	}
}

func (v *JsonNullString) UnmarshalJSON(data []byte) error {
	var x *string
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.String = *x
	} else {
		v.Valid = false
	}

	return nil
}

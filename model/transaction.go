package model

import "time"

type Transaction struct {
	Id              string    `json:"id"`
	TransactionType int       `json:"transaction_type"`
	TransactionBy   string    `json:"deposited_by"`
	Status          bool      `json:"status"`
	TransactionAt   time.Time `json:"deposited_at"`
	Amount          int       `json:"amount"`
	ReferenceId     string    `json:"reference_id"`
}

package models

import "time"

type Transaction struct {
	Date          time.Time `json:"date"`
	Type          string    `json:"type"`
	Description   string    `json:"description"`
	Merchant      string    `json:"merchant"`
	NFC           bool      `json:"nfc"`
	Amount        float64   `json:"amount"`
	Balance       float64   `json:"balance"`
	DailySequence int       `json:"daily_sequence"`
	RawText       string    `json:"raw_text"`
}

type AccountStatement struct {
	AccountHolder string        `json:"account_holder"`
	AccountNumber string        `json:"account_number"`
	IBAN          string        `json:"iban"`
	StartDate     time.Time     `json:"start_date"`
	EndDate       time.Time     `json:"end_date"`
	Transactions  []Transaction `json:"transactions"`
}

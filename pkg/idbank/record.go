package idbank

import (
	"cloud.google.com/go/civil"
)

type Record struct {
	Date   civil.Date `json:"date"`
	Amount float64    `json:"amount"`
	Unit   string     `json:"unit"`
	Debit  bool       `json:"debit"`
	Notes  string     `json:"notes"`
}

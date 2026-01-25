package cashRegister

import "time"

type CashRegisterContract struct {
	Id           int       `json:"id"`
	Description  string    `json:"description"`
	Customer     string    `json:"customer"`
	Specie       string    `json:"specie"`
	InputValue   float64   `json:"input_value"`
	OutputValue  float64   `json:"output_value"`
	TotalBalance float64   `json:"total_balance"`
	DateCreated  time.Time `json:"date_created"`
}

package product

import "time"

type ProductContract struct {
	Id             int       `json:"id"`
	Name           string    `json:"name"`
	Un             string    `json:"un"`
	Qtde           int       `json:"qtde"`
	Returned       int       `json:"returned"`
	Saled          int       `json:"saled"`
	DateOfPurchase time.Time `json:"date_of_purchase"`
}

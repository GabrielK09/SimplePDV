package shopping

import "time"

type ShoppingContract struct {
	Id             int       `json:"id"`
	Load           int16     `json:"load"`
	Operation      string    `json:"operation"`
	DateOfPurchase time.Time `json:"date_of_purchase"`
}

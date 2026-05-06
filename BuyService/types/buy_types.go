package types

import "time"

type BuyReport struct {
	Id       string    `json:"id"`
	CartId   string    `json:"cartId"`
	SaleDate time.Time `json:"saleDate"`
}

type BuyReportDto struct {
	Id       string    `json:"id"`
	Cart     Cart      `json:"cart"`
	SaleDate time.Time `json:"saleDate"`
}

type Order struct {
	Id               string `json:"id"`
	CartId           string `json:"cartId"`
	OrderedProductId string `json:"orderedProductId"`
	Quantity         int    `json:"quantity"`
}

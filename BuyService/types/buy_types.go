package types

import "time"

type BuyReport struct {
	Id       string    `json:"id"`
	CartId   string    `json:"cartId"`
	SaleDate time.Time `json:"saleDate"`
}

type BuyReportDto struct {
	Id       string     `json:"id"`
	Cart     CartEntity `json:"cart"`
	SaleDate time.Time  `json:"saleDate"`
}

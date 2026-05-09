package types

type Place struct {
	Id          string `json:"id"`
	Address     string `json:"address"`
	WorkingTime string `json:"workingTime"`
}
type Cart struct {
	Id          string  `json:"id"`
	UserId      string  `json:"userId"`
	PlaceId     *string `json:"placeId"`
	AmountToPay int     `json:"amountToPay"`
	IsConfirmed bool    `json:"isConfirmed"`
	IsBought    bool    `json:"isBought"`
}
type Order struct {
	Id               string `json:"id"`
	CartId           string `json:"cartId"`
	OrderedProductId string `json:"orderedProductId"`
	Quantity         int    `json:"quantity"`
}

type OrderDto struct {
	CartId           string `json:"cartId"`
	OrderedProductId string `json:"orderedProductId"`
	Quantity         int    `json:"quantity"`
}

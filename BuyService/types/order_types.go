package types

type CartDto struct {
	Id          string `json:"id"`
	UserId      string `json:"userId"`
	PlaceId     string `json:"placeId"`
	AmountToPay int    `json:"amountToPay"`
	IsConfirmed bool   `json:"isConfirmed"`
	IsBought    bool   `json:"isBought"`
}

type CartEntity struct {
	Id          string        `json:"id"`
	User        UserDto       `json:"user"`
	Place       PlaceDto      `json:"place"`
	Orders      []OrderEntity `json:"orders"`
	AmountToPay int           `json:"amountToPay"`
	IsConfirmed bool          `json:"isConfirmed"`
	IsBought    bool          `json:"isBought"`
}

type PlaceDto struct {
	Id          string `json:"id"`
	Address     string `json:"address"`
	WorkingTime string `json:"workingTime"`
}

type OrderDto struct {
	Id               string `json:"id"`
	CartId           string `json:"cartId"`
	OrderedProductId string `json:"orderedProductId"`
	Quantity         int    `json:"quantity"`
}

type OrderEntity struct {
	Id             string     `json:"id"`
	OrderedProduct ProductDto `json:"product"`
	Quantity       int        `json:"quantity"`
}

type ProductDto struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Cost int    `json:"cost"`
}

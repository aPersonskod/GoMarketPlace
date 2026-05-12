package types

type ErrorResponse struct {
	Error string `json:"error"`
}

type CartDto struct {
	Id          string  `json:"id"`
	UserId      string  `json:"userId"`
	PlaceId     *string `json:"placeId"`
	AmountToPay int     `json:"amountToPay"`
	IsConfirmed bool    `json:"isConfirmed"`
	IsBought    bool    `json:"isBought"`
}

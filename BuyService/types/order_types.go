package types

type Cart struct {
	Id          string `json:"id"`
	UserId      string `json:"userId"`
	PlaceId     string `json:"placeId"`
	AmountToPay int    `json:"amountToPay"`
	IsConfirmed string `json:"isConfirmed"`
	IsBought    string `json:"isBought"`
}

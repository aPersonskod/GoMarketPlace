package types

type UserDto struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Wallet int    `json:"wallet"`
	Role   string `json:"role"`
}

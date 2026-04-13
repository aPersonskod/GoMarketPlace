package types

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Wallet   int    `json:"wallet"`
	Role     string `json:"role"`
}

type UpdateUserDto struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UserDto struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Wallet int    `json:"wallet"`
	Role   string `json:"role"`
}

type UserDtoAdapter interface {
	GetUserDto() UserDto
}

func (user *User) GetUserDto() UserDto {
	return UserDto{
		Id:     user.Id,
		Name:   user.Name,
		Email:  user.Email,
		Wallet: user.Wallet,
		Role:   user.Role,
	}
}

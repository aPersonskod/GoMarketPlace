package services

import (
	"encoding/json"
	"fmt"
	"order_service/configs"
	"order_service/types"
)

type IUserService interface {
	GetUser() (*types.UserDto, error)
}

type UserService struct {
	AuthHeader string
}

func (s UserService) GetUser() (*types.UserDto, error) {
	url := fmt.Sprintf("%s/api/user-service/", configs.Env.UserServiceAddressDev)
	resp, err := ServiceHelper{}.RunRequest("GET", url, &s.AuthHeader, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	userDto := types.UserDto{}
	err = json.NewDecoder(resp.Body).Decode(&userDto)
	if err != nil {
		return nil, err
	}
	return &userDto, nil
}

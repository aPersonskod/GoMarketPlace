package services

import (
	"buy_service/configs"
	"buy_service/types"
	"encoding/json"
	"fmt"
)

type IUserService interface {
	GetUser() (*types.UserDto, error)
	SpendMoney(id string, money int) (*types.UserDto, error)
}

type UserService struct {
	AuthHeader string
}

func (s UserService) SpendMoney(id string, money int) (*types.UserDto, error) {
	url := fmt.Sprintf("%s/api/user-service/spend-money?money=%d", configs.Env.UserServiceAddressDev, money)
	resp, err := ServiceHelper{}.RunRequest("POST", url, &s.AuthHeader, nil)
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

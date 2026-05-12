package services

import (
	"encoding/json"
	"fmt"
	"orchestrator_service/configs"
	"orchestrator_service/types"
)

type IUserService interface {
	SpendMoney(money int) (*types.UserDto, error)          // action
	WalletReplenishment(money int) (*types.UserDto, error) // compensation
}

type UserService struct {
	AuthHeader string
}

func (s UserService) SpendMoney(money int) (*types.UserDto, error) {
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

func (s UserService) WalletReplenishment(money int) (*types.UserDto, error) {
	url := fmt.Sprintf("%s/api/user-service/wallet-replenishment?money=%d", configs.Env.UserServiceAddressDev, money)
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

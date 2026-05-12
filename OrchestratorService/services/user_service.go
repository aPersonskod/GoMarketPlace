package services

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/user-service/spend-money?money=%d", configs.Env.UserServiceAddressDev, money)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", s.AuthHeader)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error sending request: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		userDto := types.UserDto{}
		err = json.NewDecoder(resp.Body).Decode(&userDto)
		if err != nil {
			return nil, err
		}
		return &userDto, nil
	}
	errResp := types.ErrorResponse{}
	json.NewDecoder(resp.Body).Decode(&errResp)
	return nil, fmt.Errorf("Server returned error: %s (Status: %d)", errResp.Error, resp.StatusCode)
}

func (s UserService) WalletReplenishment(money int) (*types.UserDto, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/user-service/wallet-replenishment?money=%d", configs.Env.UserServiceAddressDev, money)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", s.AuthHeader)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error sending request: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		userDto := types.UserDto{}
		err = json.NewDecoder(resp.Body).Decode(&userDto)
		if err != nil {
			return nil, err
		}
		return &userDto, nil
	}
	errResp := types.ErrorResponse{}
	json.NewDecoder(resp.Body).Decode(&errResp)
	return nil, fmt.Errorf("Server returned error: %s (Status: %d)", errResp.Error, resp.StatusCode)
}

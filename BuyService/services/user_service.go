package services

import (
	"buy_service/configs"
	"buy_service/types"
	"encoding/json"
	"fmt"
	"net/http"
)

type IUserService interface {
	SpendMoney(id string, money int) (*types.UserDto, error)
}

type UserService struct {
	AuthHeader string
}

func (s UserService) SpendMoney(id string, money int) (*types.UserDto, error) {
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
	return nil, fmt.Errorf("Error, status code: %d", resp.StatusCode)
}

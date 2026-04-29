package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"order_service/types"

	"github.com/gin-gonic/gin"
)

type IUserService interface {
	GetUser() (*types.UserDto, error)
}

type UserService struct {
	Context *gin.Context
}

func (s UserService) GetUser() (*types.UserDto, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8080/api/user-service/", nil) // TODO add to env
	if err != nil {
		return nil, err
	}
	authHeader := s.Context.GetHeader("Authorization")
	if authHeader == "" {
		return nil, errors.New("Auth header required")
	}
	req.Header.Set("Authorization", authHeader)
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
	return nil, errors.New("Bad request")

}

package services

import (
	"encoding/json"
	"fmt"
	"orchestrator_service/configs"
	"orchestrator_service/types"
)

type IBuyService interface {
	BuyCart(cart types.CartDto) error // action, do not need compensation, because last
}

type BuyService struct {
	AuthHeader string
}

func (s BuyService) BuyCart(cart types.CartDto) error {
	jsonCart, err := json.Marshal(cart)
	if err != nil {
		return fmt.Errorf("Error to convert cart to json !!!")
	}
	url := fmt.Sprintf("%s/api/buy-service/buy-cart", configs.Env.BuyServiceAddressDev)
	resp, err := ServiceHelper{}.RunRequest("POST", url, &s.AuthHeader, jsonCart)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

package services

import (
	"encoding/json"
	"fmt"
	"orchestrator_service/configs"
	"orchestrator_service/types"
)

type IOrderService interface {
	ConfirmCart(placeId string) (*types.CartDto, error) // action
	UnconfirmCart() (*types.CartDto, error)             // compensation

	MarkCartAsBought(cartId string) error    // action
	MarkCartAsNotBought(cartId string) error // compensation
}

type OrderService struct {
	AuthHeader string
}

func (s OrderService) ConfirmCart(placeId string) (*types.CartDto, error) {
	url := fmt.Sprintf("%s/api/order-service/confirm-cart?placeId=%s", configs.Env.OrderServiceAddressDev, placeId)
	resp, err := ServiceHelper{}.RunRequest("POST", url, &s.AuthHeader, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	cartDto := types.CartDto{}
	err = json.NewDecoder(resp.Body).Decode(&cartDto)
	if err != nil {
		return nil, err
	}
	return &cartDto, nil
}

func (s OrderService) UnconfirmCart() (*types.CartDto, error) {
	url := fmt.Sprintf("%s/api/order-service/unconfirm-cart", configs.Env.OrderServiceAddressDev)
	resp, err := ServiceHelper{}.RunRequest("POST", url, &s.AuthHeader, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	cartDto := types.CartDto{}
	err = json.NewDecoder(resp.Body).Decode(&cartDto)
	if err != nil {
		return nil, err
	}
	return &cartDto, nil
}

func (s OrderService) MarkCartAsBought(cartid string) error {
	url := fmt.Sprintf("%s/api/order-service/mark-cart-as-bought?cartId=%s", configs.Env.OrderServiceAddressDev, cartid)
	resp, err := ServiceHelper{}.RunRequest("POST", url, &s.AuthHeader, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (s OrderService) MarkCartAsNotBought(cartid string) error {
	url := fmt.Sprintf("%s/api/order-service/mark-cart-as-not-bought?cartId=%s", configs.Env.OrderServiceAddressDev, cartid)
	resp, err := ServiceHelper{}.RunRequest("POST", url, &s.AuthHeader, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

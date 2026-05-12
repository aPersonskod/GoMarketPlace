package services

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/order-service/confirm-cart?placeId=%s", configs.Env.OrderServiceAddressDev, placeId)
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
		cartDto := types.CartDto{}
		err = json.NewDecoder(resp.Body).Decode(&cartDto)
		if err != nil {
			return nil, err
		}
		return &cartDto, nil
	}
	errResp := types.ErrorResponse{}
	json.NewDecoder(resp.Body).Decode(&errResp)
	return nil, fmt.Errorf("Server returned error: %s (Status: %d)", errResp.Error, resp.StatusCode)
}

func (s OrderService) UnconfirmCart() (*types.CartDto, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/order-service/unconfirm-cart", configs.Env.OrderServiceAddressDev)
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
		cartDto := types.CartDto{}
		err = json.NewDecoder(resp.Body).Decode(&cartDto)
		if err != nil {
			return nil, err
		}
		return &cartDto, nil
	}
	errResp := types.ErrorResponse{}
	json.NewDecoder(resp.Body).Decode(&errResp)
	return nil, fmt.Errorf("Server returned error: %s (Status: %d)", errResp.Error, resp.StatusCode)
}

func (s OrderService) MarkCartAsBought(cartid string) error {
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/order-service/mark-cart-as-bought?cartId=%s", configs.Env.OrderServiceAddressDev, cartid)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", s.AuthHeader)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending request: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}
	errResp := types.ErrorResponse{}
	json.NewDecoder(resp.Body).Decode(&errResp)
	return fmt.Errorf("Server returned error: %s (Status: %d)", errResp.Error, resp.StatusCode)
}

func (s OrderService) MarkCartAsNotBought(cartid string) error {
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/order-service/mark-cart-as-not-bought?cartId=%s", configs.Env.OrderServiceAddressDev, cartid)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", s.AuthHeader)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending request: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}
	errResp := types.ErrorResponse{}
	json.NewDecoder(resp.Body).Decode(&errResp)
	return fmt.Errorf("Server returned error: %s (Status: %d)", errResp.Error, resp.StatusCode)
}

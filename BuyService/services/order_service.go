package services

import (
	"buy_service/configs"
	"buy_service/types"
	"encoding/json"
	"fmt"
	"net/http"
)

type IOrderService interface {
	GetCart(userId string) (*types.Cart, error)
	GetBoughtCarts(userId string) ([]types.Cart, error)
	GetOrders(cartId string) ([]types.Order, error)
	MarkCartAsBought(cartid string) error
}

type OrderService struct {
	AuthHeader string
}

func (s OrderService) GetCart(userId string) (*types.Cart, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/order-service/get-cart", configs.Env.OrderServiceAddressDev)
	req, err := http.NewRequest("GET", url, nil)
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
		cart := types.Cart{}
		err = json.NewDecoder(resp.Body).Decode(&cart)
		if err != nil {
			return nil, err
		}
		return &cart, nil
	}
	return nil, fmt.Errorf("Error, status code: %d", resp.StatusCode)
}
func (s OrderService) GetBoughtCarts(userId string) ([]types.Cart, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/order-service/get-bought-carts", configs.Env.OrderServiceAddressDev)
	req, err := http.NewRequest("GET", url, nil)
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
		carts := []types.Cart{}
		err = json.NewDecoder(resp.Body).Decode(&carts)
		if err != nil {
			return nil, err
		}
		return carts, nil
	}
	return nil, fmt.Errorf("Error, status code: %d", resp.StatusCode)
}
func (s OrderService) GetOrders(cartId string) ([]types.Order, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/order-service/get-cart-orders/%s", configs.Env.OrderServiceAddressDev, cartId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error sending request: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		orders := []types.Order{}
		err = json.NewDecoder(resp.Body).Decode(&orders)
		if err != nil {
			return nil, err
		}
		return orders, nil
	}
	return nil, fmt.Errorf("Error, status code: %d", resp.StatusCode)
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
	return fmt.Errorf("Error, status code: %d", resp.StatusCode)
}

package services

import (
	"buy_service/configs"
	"buy_service/types"
	"encoding/json"
	"fmt"
)

type IOrderService interface {
	GetCart(userId string) (*types.CartDto, error)
	GetBoughtCarts(userId string) ([]types.CartDto, error)
	GetOrders(cartId string) ([]types.OrderDto, error)
	MarkCartAsBought(cartid string) error
	GetPlace(placeId string) (*types.PlaceDto, error)

	GetProduct(productId string) (*types.ProductDto, error)
}

type OrderService struct {
	AuthHeader string
}

func (s OrderService) GetCart(userId string) (*types.CartDto, error) {
	url := fmt.Sprintf("%s/api/order-service/get-cart", configs.Env.OrderServiceAddressDev)
	resp, err := ServiceHelper{}.RunRequest("GET", url, &s.AuthHeader, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	cart := types.CartDto{}
	err = json.NewDecoder(resp.Body).Decode(&cart)
	if err != nil {
		return nil, err
	}
	return &cart, nil
}
func (s OrderService) GetBoughtCarts(userId string) ([]types.CartDto, error) {
	url := fmt.Sprintf("%s/api/order-service/get-bought-carts", configs.Env.OrderServiceAddressDev)
	resp, err := ServiceHelper{}.RunRequest("GET", url, &s.AuthHeader, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	carts := []types.CartDto{}
	err = json.NewDecoder(resp.Body).Decode(&carts)
	if err != nil {
		return nil, err
	}
	return carts, nil
}
func (s OrderService) GetOrders(cartId string) ([]types.OrderDto, error) {
	url := fmt.Sprintf("%s/api/order-service/get-cart-orders/%s", configs.Env.OrderServiceAddressDev, cartId)
	resp, err := ServiceHelper{}.RunRequest("GET", url, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	orders := []types.OrderDto{}
	err = json.NewDecoder(resp.Body).Decode(&orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
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
func (s OrderService) GetPlace(placeId string) (*types.PlaceDto, error) {
	url := fmt.Sprintf("%s/api/order-service/get-place/%s", configs.Env.OrderServiceAddressDev, placeId)
	resp, err := ServiceHelper{}.RunRequest("GET", url, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	place := types.PlaceDto{}
	err = json.NewDecoder(resp.Body).Decode(&place)
	if err != nil {
		return nil, err
	}
	return &place, nil
}
func (s OrderService) GetProduct(productId string) (*types.ProductDto, error) {
	url := fmt.Sprintf("%s/api/product-service/%s", configs.Env.ProductServiceAddressDev, productId)
	resp, err := ServiceHelper{}.RunRequest("GET", url, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	product := types.ProductDto{}
	err = json.NewDecoder(resp.Body).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

package services

import (
	"encoding/json"
	"fmt"
	"order_service/configs"
	"order_service/types"
)

type IProductService interface {
	GetProduct(productId string) (*types.Product, error)
}

type ProductService struct{}

func (s ProductService) GetProduct(productId string) (*types.Product, error) {
	url := fmt.Sprintf("%s/api/product-service/%s", configs.Env.ProductServiceAddressDev, productId)
	resp, err := ServiceHelper{}.RunRequest("GET", url, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	product := types.Product{}
	err = json.NewDecoder(resp.Body).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

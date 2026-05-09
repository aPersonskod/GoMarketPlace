package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"order_service/types"
)

type IProductService interface {
	GetProduct(productId string) (*types.Product, error)
}

type ProductService struct{}

func (s ProductService) GetProduct(productId string) (*types.Product, error) {
	client := &http.Client{}
	url := fmt.Sprintf("http://localhost:8081/api/product-service/%s", productId) // TODO add to env
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
		p := types.Product{}
		err = json.NewDecoder(resp.Body).Decode(&p)
		if err != nil {
			return nil, err
		}
		return &p, nil
	}
	return nil, errors.New("Bad request")
}

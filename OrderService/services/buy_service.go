package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"order_service/types"
)

type IBuyService interface {
	BuyCart(cart types.Cart) error
}

type BuyService struct {
	AuthHeader string
}

func (s BuyService) BuyCart(cart types.Cart) error {
	jsonCart, err := json.Marshal(cart)
	if err != nil {
		return fmt.Errorf("Error to convert cart to json !!!")
	}
	body := bytes.NewBuffer(jsonCart)
	req, err := http.NewRequest("POST", "http://localhost:8083/api/buy-service/buy-cart", body) // TODO add to env
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", s.AuthHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending request: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned %d: %s", resp.StatusCode, string(body))
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return fmt.Errorf("Error code: %d", resp.StatusCode)
}

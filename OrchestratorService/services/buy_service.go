package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	body := bytes.NewBuffer(jsonCart)
	url := fmt.Sprintf("%s/api/buy-service/buy-cart", configs.Env.BuyServiceAddressDev)
	req, err := http.NewRequest("POST", url, body)
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
	errResp := types.ErrorResponse{}
	json.NewDecoder(resp.Body).Decode(&errResp)
	return fmt.Errorf("Server returned error: %s (Status: %d)", errResp.Error, resp.StatusCode)
}

package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type ServiceHelper struct{}

func (s ServiceHelper) RunRequest(method, url string, authHeader *string, jsonMarshal []byte) (*http.Response, error) {
	var body io.Reader
	if jsonMarshal != nil {
		body = bytes.NewBuffer(jsonMarshal)
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if authHeader != nil {
		req.Header.Set("Authorization", *authHeader)
		req.Header.Set("Content-Type", "application/json")
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error sending request: %s", err.Error())
	}
	//defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return resp, nil
	}
	errResp := ErrorResponse{}
	json.NewDecoder(resp.Body).Decode(&errResp)
	return nil, fmt.Errorf("Server returned error: %s (Status: %d)", errResp.Error, resp.StatusCode)
}

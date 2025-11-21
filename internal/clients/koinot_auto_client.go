package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"loan_service/configs"
	"loan_service/internal/dto"
	"net/http"
	"time"
)

type KoinotAutoClient struct {
	httpClient *http.Client
	baseURL    string
	token      string
}

func NewKoinotAutoClient(cfg configs.HTTPClientConfig) (*KoinotAutoClient, error) {
	timeout, err := time.ParseDuration(cfg.Timeout)
	if err != nil {
		return nil, fmt.Errorf("Invalid timeout format koinot auto client: %w", err)
	}

	return &KoinotAutoClient{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		baseURL: cfg.BaseURL,
		token:   cfg.Token,
	}, nil
}

func (c *KoinotAutoClient) ListVehicles(ctx context.Context) ([]dto.Vehicle, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/vehicles", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("koinot auto returned status %d", resp.StatusCode)
	}

	var vehicles []dto.Vehicle
	if err := json.NewDecoder(resp.Body).Decode(&vehicles); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return vehicles, nil
}

func (c *KoinotAutoClient) SendLoanApplication(ctx context.Context, loanApp *dto.LoanApplication) error {

	jsonData, err := json.Marshal(loanApp)
	if err != nil {
		return fmt.Errorf("could not marshall request body: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/loan-application", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request:: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("koinot auto returned status %d", resp.StatusCode)
	}

	return nil
}

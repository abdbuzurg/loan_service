package clients

import (
	"fmt"
	"loan_service/configs"
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

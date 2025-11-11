package clients

import (
	"fmt"
	"loan_service/configs"
	"net/http"
	"time"
)

type AsrLeasingClient struct {
	httpClient *http.Client
	baseURL    string
	token      string
}

func NewAsrLeasingClient(cfg configs.HTTPClientConfig) (*AsrLeasingClient, error) {
	timeout, err := time.ParseDuration(cfg.Timeout)
	if err != nil {
		return nil, fmt.Errorf("Invalid timeout format for asr leasing client: %w", err)
	}

	return &AsrLeasingClient{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		baseURL: cfg.BaseURL,
		token:   cfg.Token,
	}, nil
}

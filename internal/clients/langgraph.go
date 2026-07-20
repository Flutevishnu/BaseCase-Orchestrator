package clients

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

type LangGraphClient struct {
	baseURL       string
	gatewayURL    string
	webhookSecret string
	httpClient    *http.Client
}

func NewLangGraphClient() *LangGraphClient {
	url := os.Getenv("LANGGRAPH_URL")
	if url == "" {
		url = "http://localhost:8123"
	}
	return &LangGraphClient{
		baseURL:       url,
		gatewayURL:    os.Getenv("GATEWAY_URL"),
		webhookSecret: os.Getenv("LANGGRAPH_WEBHOOK_SECRET"),
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// CreateThread calls LangGraph to initialize an AI thread
func (c *LangGraphClient) CreateThread(ctx context.Context, userId string) (map[string]interface{}, error) {

	// Mock response for scaffolding
	return map[string]interface{}{
		"thread_id": "th_mock_123",
	}, nil
}

func (c *LangGraphClient) StreamEvaluationRun(ctx context.Context, userId, submissionId string) (map[string]interface{}, error) {

	webhookURL := fmt.Sprintf("%s/api/webhooks/evaluations?secret=%s", c.gatewayURL, c.webhookSecret)

	fmt.Printf("[Mock] LangGraphClient: Configured webhook payload: %s\n", webhookURL)

	// Mock response for scaffolding
	return map[string]interface{}{
		"run_id": "run_mock_123",
	}, nil
}

package clients

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

type ExpressClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewExpressClient() *ExpressClient {
	url := os.Getenv("EXPRESS_URL")
	if url == "" {
		url = "http://localhost:3000"
	}
	return &ExpressClient{
		baseURL: url,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *ExpressClient) CreateSession(ctx context.Context, userId string) (map[string]interface{}, error) {

	// Mock response for scaffolding
	return map[string]interface{}{
		"id":    "sess_mock_123",
		"isNew": true,
	}, nil
}

func (c *ExpressClient) MarkSessionFailed(ctx context.Context, sessionId string) error {
	fmt.Printf("[Mock] ExpressClient: Marking session %s as FAILED\n", sessionId)
	return nil
}

func (c *ExpressClient) CreateSubmission(ctx context.Context, userId string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"id": "sub_mock_123",
	}, nil
}

func (c *ExpressClient) MarkSubmissionFailed(ctx context.Context, submissionId string) error {

	fmt.Printf("[Mock] ExpressClient: Marking submission %s as FAILED\n", submissionId)
	return nil
}

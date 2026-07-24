package clients

import (
	"bytes"
	"context"
	"encoding/json"
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
		url = "http://localhost:3001"
	}
	return &ExpressClient{
		baseURL: url,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *ExpressClient) CreateSession(ctx context.Context, userId string, problemId string, restart bool) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/problems/%s/sessions", c.baseURL, problemId)
	body, _ := json.Marshal(map[string]bool{"restart": restart})
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", userId)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("express returned status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *ExpressClient) MarkSessionFailed(ctx context.Context, sessionId string) error {
	url := fmt.Sprintf("%s/api/v1/sessions/%s/status", c.baseURL, sessionId)
	body, _ := json.Marshal(map[string]string{"status": "FAILED"})
	
	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("express returned status: %d", resp.StatusCode)
	}
	return nil
}

func (c *ExpressClient) CreateSubmission(ctx context.Context, userId string, sessionId string, code string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/sessions/%s/submissions", c.baseURL, sessionId)
	body, _ := json.Marshal(map[string]string{"code": code})
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", userId)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("express returned status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *ExpressClient) MarkSubmissionFailed(ctx context.Context, submissionId string) error {
	url := fmt.Sprintf("%s/api/v1/submissions/%s/status", c.baseURL, submissionId)
	body, _ := json.Marshal(map[string]string{"status": "failed"})
	
	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("express returned status: %d", resp.StatusCode)
	}
	return nil
}

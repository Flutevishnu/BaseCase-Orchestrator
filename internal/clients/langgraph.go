package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

type LangGraphClient struct {
	baseURL       string
	gatewayURL    string
	webhookSecret string
	httpClient    *http.Client
	// streamClient is for runs/stream. http.Client.Timeout also covers reading the
	// body, so the 15 s httpClient would cut every SSE stream mid-run. Streams are
	// bounded by the request context instead (cancelled when the caller disconnects).
	streamClient *http.Client
}

func NewLangGraphClient() *LangGraphClient {
	url := os.Getenv("LANGGRAPH_URL")
	if url == "" {
		url = "http://localhost:2024"
	}
	return &LangGraphClient{
		baseURL:       url,
		gatewayURL:    os.Getenv("GATEWAY_URL"),
		webhookSecret: os.Getenv("LANGGRAPH_WEBHOOK_SECRET"),
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		streamClient: &http.Client{
			Transport: &http.Transport{
				Proxy:                 http.ProxyFromEnvironment,
				DialContext:           (&net.Dialer{Timeout: 10 * time.Second}).DialContext,
				ResponseHeaderTimeout: 30 * time.Second,
			},
		},
	}
}

// CreateThread calls LangGraph to initialize an AI thread
func (c *LangGraphClient) CreateThread(ctx context.Context, userId string, threadId string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/threads", c.baseURL)
	
	payload := map[string]interface{}{
		"thread_id": threadId,
		"if_exists": "do_nothing",
		"metadata": map[string]string{
			"userId":   userId,
			"graph_id": os.Getenv("LANGGRAPH_ASSISTANT_ID"),
		},
	}
	body, _ := json.Marshal(payload)
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("langgraph returned status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *LangGraphClient) StreamEvaluationRun(ctx context.Context, threadId string, submissionId string, code string) (*http.Response, error) {
	webhookURL := fmt.Sprintf("%s/api/webhooks/evaluations?secret=%s", c.gatewayURL, c.webhookSecret)

	url := fmt.Sprintf("%s/threads/%s/runs/stream", c.baseURL, threadId)
	
	payload := map[string]interface{}{
		"webhook": webhookURL,
		"input": map[string]interface{}{
			"submission_id": submissionId,
			"code_submission": code,
		},
		"stream_mode":  []string{"messages", "values"},
		"assistant_id": os.Getenv("LANGGRAPH_ASSISTANT_ID"),
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.streamClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		errBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("langgraph returned status: %d, body: %s", resp.StatusCode, string(errBody))
	}

	return resp, nil
}

func (c *LangGraphClient) UpdateState(ctx context.Context, threadId string, values map[string]interface{}) error {
	url := fmt.Sprintf("%s/threads/%s/state", c.baseURL, threadId)
	
	payload := map[string]interface{}{
		"values": values,
	}
	body, _ := json.Marshal(payload)
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
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
		errBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("langgraph returned status: %d, body: %s", resp.StatusCode, string(errBody))
	}

	return nil
}

func (c *LangGraphClient) GetState(ctx context.Context, threadId string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/threads/%s/state", c.baseURL, threadId)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("langgraph returned status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *LangGraphClient) StreamChat(ctx context.Context, threadId string, input map[string]interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s/threads/%s/runs/stream", c.baseURL, threadId)
	
	payload := map[string]interface{}{
		"input": input,
		"stream_mode": []string{"messages", "values"},
		"assistant_id": os.Getenv("LANGGRAPH_ASSISTANT_ID"),
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.streamClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		resp.Body.Close()
		return nil, fmt.Errorf("langgraph returned status: %d", resp.StatusCode)
	}

	return resp, nil
}

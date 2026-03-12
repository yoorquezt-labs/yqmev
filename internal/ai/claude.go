package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const claudeAPI = "https://api.anthropic.com/v1/messages"

// Claude implements Provider using the Anthropic Messages API with SSE streaming.
type Claude struct {
	apiKey string
	model  string
	client *http.Client
}

// NewClaude creates a Claude provider. Default model: claude-sonnet-4-20250514.
func NewClaude(apiKey, model string) *Claude {
	if model == "" {
		model = "claude-sonnet-4-20250514"
	}
	return &Claude{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{},
	}
}

func (c *Claude) Name() string { return "claude" }

func (c *Claude) Analyze(ctx context.Context, req AnalyzeRequest) (string, error) {
	messages := buildMessages(req)

	body := map[string]any{
		"model":      c.model,
		"max_tokens": 2048,
		"system":     mevSystemPrompt(req.MEVContext),
		"messages":   messages,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", claudeAPI, bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("api call: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("api error %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}

	if len(result.Content) == 0 {
		return "", fmt.Errorf("empty response")
	}

	return result.Content[0].Text, nil
}

func (c *Claude) Stream(ctx context.Context, req AnalyzeRequest, out chan<- string) error {
	defer close(out)

	messages := buildMessages(req)

	body := map[string]any{
		"model":      c.model,
		"max_tokens": 2048,
		"system":     mevSystemPrompt(req.MEVContext),
		"messages":   messages,
		"stream":     true,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", claudeAPI, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("api call: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("api error %d: %s", resp.StatusCode, string(respBody))
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "event:") {
			continue
		}
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := line[6:]
		if data == "[DONE]" {
			return nil
		}

		var event struct {
			Type  string `json:"type"`
			Delta struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"delta"`
		}
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			continue
		}

		if event.Type == "content_block_delta" && event.Delta.Text != "" {
			select {
			case out <- event.Delta.Text:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return scanner.Err()
}

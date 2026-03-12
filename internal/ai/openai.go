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

// OpenAI implements Provider using the OpenAI Chat Completions API.
// Also works with Ollama by setting baseURL to http://localhost:11434/v1.
type OpenAI struct {
	apiKey  string
	model   string
	baseURL string
	client  *http.Client
}

// NewOpenAI creates an OpenAI provider. Default model: gpt-4o.
func NewOpenAI(apiKey, model, baseURL string) *OpenAI {
	if model == "" {
		model = "gpt-4o"
	}
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	return &OpenAI{
		apiKey:  apiKey,
		model:   model,
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (o *OpenAI) Name() string { return "openai" }

func (o *OpenAI) Analyze(ctx context.Context, req AnalyzeRequest) (string, error) {
	var messages []map[string]string
	messages = append(messages, map[string]string{
		"role":    "system",
		"content": mevSystemPrompt(req.MEVContext),
	})
	for _, m := range req.History {
		messages = append(messages, map[string]string{
			"role":    m.Role,
			"content": m.Content,
		})
	}
	messages = append(messages, map[string]string{
		"role":    "user",
		"content": req.Question,
	})

	body := map[string]any{
		"model":      o.model,
		"max_tokens": 2048,
		"messages":   messages,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshal: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", o.baseURL+"/chat/completions", bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if o.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+o.apiKey)
	}

	resp, err := o.client.Do(httpReq)
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
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("parse: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("empty response")
	}

	return result.Choices[0].Message.Content, nil
}

func (o *OpenAI) Stream(ctx context.Context, req AnalyzeRequest, out chan<- string) error {
	defer close(out)

	var messages []map[string]string
	messages = append(messages, map[string]string{
		"role":    "system",
		"content": mevSystemPrompt(req.MEVContext),
	})
	for _, m := range req.History {
		messages = append(messages, map[string]string{
			"role":    m.Role,
			"content": m.Content,
		})
	}
	messages = append(messages, map[string]string{
		"role":    "user",
		"content": req.Question,
	})

	body := map[string]any{
		"model":      o.model,
		"max_tokens": 2048,
		"messages":   messages,
		"stream":     true,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", o.baseURL+"/chat/completions", bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if o.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+o.apiKey)
	}

	resp, err := o.client.Do(httpReq)
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
		if line == "" || !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := line[6:]
		if data == "[DONE]" {
			return nil
		}

		var event struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			continue
		}

		if len(event.Choices) > 0 && event.Choices[0].Delta.Content != "" {
			select {
			case out <- event.Choices[0].Delta.Content:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return scanner.Err()
}

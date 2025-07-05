package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// OpenAIAgent implements the Agent interface using OpenAI's API.
type OpenAIAgent struct {
	APIKey  string
	Model   string
	BaseURL string
	Client  interface {
		Do(req *http.Request) (*http.Response, error)
	}
}

// NewOpenAIAgent creates a new OpenAIAgent, reading config from environment variables.
func NewOpenAIAgent() *OpenAIAgent {
	key := os.Getenv("OPENAI_API_KEY")
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-3.5-turbo"
	}
	base := os.Getenv("OPENAI_API_BASE")
	if base == "" {
		base = "https://api.openai.com/v1"
	}
	return &OpenAIAgent{
		APIKey:  key,
		Model:   model,
		BaseURL: base,
		Client:  &http.Client{},
	}
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIRequest struct {
	Model     string          `json:"model"`
	Messages  []openAIMessage `json:"messages"`
	MaxTokens int             `json:"max_tokens,omitempty"`
}

type openAIResponse struct {
	Choices []struct {
		Message openAIMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Respond sends the prompt to OpenAI and returns the reply.
func (a *OpenAIAgent) Respond(prompt string) (string, error) {
	debug := os.Getenv("BINKS_DEBUG_AI") == "1"
	if debug {
		fmt.Fprintf(os.Stderr, "[OpenAIAgent] Received prompt: %q\n", prompt)
	}
	if a.APIKey == "" {
		return "", errors.New("AI is not configured. Set OPENAI_API_KEY environment variable")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	url := a.BaseURL + "/chat/completions"
	payload := openAIRequest{
		Model:    a.Model,
		Messages: []openAIMessage{{Role: "user", Content: prompt}},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	if debug {
		fmt.Fprintf(os.Stderr, "[OpenAIAgent] Sending request to %s: %s\n", url, string(body))
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+a.APIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := a.Client.Do(req)
	if err != nil {
		return "", errors.New("AI error: " + err.Error())
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if debug {
		fmt.Fprintf(os.Stderr, "[OpenAIAgent] Raw response: %s\n", string(respBody))
	}
	var aiResp openAIResponse
	if err := json.Unmarshal(respBody, &aiResp); err != nil {
		return "", errors.New("AI error: failed to parse response")
	}
	if aiResp.Error != nil {
		return "", errors.New("OpenAI API error: " + aiResp.Error.Message)
	}
	if len(aiResp.Choices) == 0 {
		return "", errors.New("AI error: no response from model")
	}
	content := aiResp.Choices[0].Message.Content
	return strings.TrimRight(content, "\n\r "), nil
}

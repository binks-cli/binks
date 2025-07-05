package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"
)

type fakeHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (f *fakeHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return f.DoFunc(req)
}

func TestOpenAIAgent_Respond_Success(t *testing.T) {
	agent := NewOpenAIAgent()
	agent.APIKey = "test-key"
	agent.Model = "gpt-3.5-turbo"
	agent.BaseURL = "https://api.openai.com/v1"
	agent.Client = &fakeHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			resp := openAIResponse{
				Choices: []struct{ Message openAIMessage "json:\"message\"" }{{Message: openAIMessage{Role: "assistant", Content: "Hello!"}}},
			}
			b, _ := json.Marshal(resp)
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(b)),
			}, nil
		},
	}
	resp, err := agent.Respond("Hi")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp != "Hello!" {
		t.Errorf("expected 'Hello!', got '%s'", resp)
	}
}

func TestOpenAIAgent_Respond_NoAPIKey(t *testing.T) {
	agent := NewOpenAIAgent()
	agent.APIKey = ""
	_, err := agent.Respond("Hi")
	if err == nil || err.Error() != "AI is not configured. Set OPENAI_API_KEY environment variable" {
		t.Errorf("expected missing key error, got %v", err)
	}
}

func TestOpenAIAgent_Respond_APIError(t *testing.T) {
	agent := NewOpenAIAgent()
	agent.APIKey = "test-key"
	agent.Client = &fakeHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			resp := openAIResponse{
				Error: &struct{ Message string "json:\"message\"" }{Message: "unauthorized"},
			}
			b, _ := json.Marshal(resp)
			return &http.Response{
				StatusCode: 401,
				Body:       io.NopCloser(bytes.NewReader(b)),
			}, nil
		},
	}
	_, err := agent.Respond("Hi")
	if err == nil || err.Error() != "OpenAI API error: unauthorized" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestOpenAIAgent_Respond_Timeout(t *testing.T) {
	agent := NewOpenAIAgent()
	agent.APIKey = "test-key"
	agent.Client = &fakeHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			time.Sleep(20 * time.Millisecond)
			return nil, context.DeadlineExceeded
		},
	}
	_, err := agent.Respond("Hi")
	if err == nil || err.Error() == "" {
		t.Errorf("expected timeout error, got %v", err)
	}
}

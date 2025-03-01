package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/liushuangls/go-anthropic/v2"
	"github.com/sashabaranov/go-openai"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
)

type OpenAIParams struct {
	System string
	Model string
}

type OpenAIChat struct {
    client *openai.Client
    params OpenAIParams
}

var foo = 1337
var temperature float32 = 0.55
const shapeMessage = `Response Format MUST BE line separated Row,Col tuples
Example output with 2 positions specified
20,3
12,4

This specified position of row 20 col 3 and second, separated by new line, row 12 and col 4
`

func (o *OpenAIChat) chat(chat string, ctx context.Context) (string, error) {
	slog.Error("OpenAIChat chat", "model", o.params.Model, "chat", chat)
    resp, err := o.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
        Model: o.params.Model,
        Seed: &foo,
        Messages: []openai.ChatCompletionMessage{
            {
                Role: openai.ChatMessageRoleSystem,
                Content: o.params.System,
            },
            {
                Role: openai.ChatMessageRoleSystem,
                Content: shapeMessage,
            },
            {
                Role: openai.ChatMessageRoleUser,
                Content: chat,
            },
        },
    })
    if err != nil {
        return "", err
    }

    return resp.Choices[0].Message.Content, nil
}

func NewOpenAIChat(secret string, params OpenAIParams) *OpenAIChat {
    client := openai.NewClient(secret)
    return &OpenAIChat{
        params: params,
        client: client,
    }
}

type StatefulOpenAIChat struct {
    ai *OpenAIChat
    ctx context.Context
}

func (s *StatefulOpenAIChat) Name() string {
    return "openai"
}

func NewStatefulOpenAIChat(params OpenAIParams, ctx context.Context) *StatefulOpenAIChat {
    return &StatefulOpenAIChat{
        ai: NewOpenAIChat(os.Getenv("OPENAI_API_KEY"), params),
        ctx: ctx,
    }
}

func (s *StatefulOpenAIChat) prompt(p string, ctx context.Context) (string, error) {
    str, err := s.ai.chat(p, ctx)
    if err != nil {
        return "", err
    }
    return str, err
}

func (s *StatefulOpenAIChat) ReadWithTimeout(p string, t time.Duration) (string, error) {
    ctx, cancel := context.WithCancel(s.ctx)
    go func() {
        <-time.NewTimer(t).C
        cancel()
    }()

    return s.ai.chat(p, ctx)
}

type ClaudeSonnetParams struct {
	System string
	Model string
}

type ClaudeSonnet struct {
    client *anthropic.Client
    params ClaudeSonnetParams
    ctx context.Context
}

func NewClaudeSonnet(params ClaudeSonnetParams, ctx context.Context) *ClaudeSonnet {
    client := anthropic.NewClient(os.Getenv("ANTHROPIC_API_KEY"))
	params.System = params.System + "\n" + shapeMessage
    return &ClaudeSonnet{
        ctx: ctx,
        client: client,
        params: params,
    }
}

type AnthropicContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type AnthropicResponse struct {
	Content []AnthropicContent `json:"content"`
}

func getResults(thinkingData AnthropicResponse) (string, error) {
	for _, content := range thinkingData.Content {
		if content.Type == "text" {
			return content.Text, nil
		}
	}
	return "", fmt.Errorf("no text content found")
}

func (c *ClaudeSonnet) makeRequest(msg string) (string, error) {
	payload, err := json.Marshal(map[string]interface{}{
		"model": c.params.Model,
		"max_tokens": 8192,
		"system": c.params.System,
		"messages": []map[string]any{
			{"role": "user", "content": msg},
		},
	})
	assert.NoError(err, "unable to marshal payload")

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("unable to make request: %s", err)
	}
	req.Header.Set("x-api-key", os.Getenv("ANTHROPIC_API_KEY"))
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("unable to make client request: %w", err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non 200 status code: %d -- %s", resp.StatusCode, data)
	}

	defer resp.Body.Close()

	var result AnthropicResponse
	err = json.Unmarshal(data, &result)
	if err != nil {
		return "", fmt.Errorf("unable to unmarshal response data(%w): %d %s", err, resp.StatusCode, data)
	}

	return getResults(result)
}

func (o *ClaudeSonnet) chat(chat string, ctx context.Context) (string, error) {
	slog.Error("ClaudeSonnet chat", "chat", chat)
	resp, err := o.makeRequest(chat)

	return resp, err
}

func (s *ClaudeSonnet) Name() string {
    return "anthropic"
}

func (s *ClaudeSonnet) ReadWithTimeout(p string, t time.Duration) (string, error) {
    ctx, cancel := context.WithCancel(s.ctx)
    go func() {
        <-time.NewTimer(t).C
        cancel()
    }()

    return s.chat(p, ctx)
}


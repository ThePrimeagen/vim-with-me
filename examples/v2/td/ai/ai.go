package ai

import (
	"context"
	"os"
	"time"

	"github.com/liushuangls/go-anthropic/v2"
	"github.com/sashabaranov/go-openai"
)

type OpenAIChat struct {
    client *openai.Client
    system string
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
    resp, err := o.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
        Model: openai.GPT4oMini20240718,
        Temperature: temperature,
        Seed: &foo,
        Messages: []openai.ChatCompletionMessage{
            {
                Role: openai.ChatMessageRoleSystem,
                Content: o.system,
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

func NewOpenAIChat(secret string, system string) *OpenAIChat {
    client := openai.NewClient(secret)
    return &OpenAIChat{
        system: system,
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

func NewStatefulOpenAIChat(system string, ctx context.Context) *StatefulOpenAIChat {
    return &StatefulOpenAIChat{
        ai: NewOpenAIChat(os.Getenv("OPENAI_API_KEY"), system),
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

type ClaudeSonnet struct {
    client *anthropic.Client
    system string
    ctx context.Context
}

func NewClaudeSonnet(system string, ctx context.Context) *ClaudeSonnet {
    client := anthropic.NewClient(os.Getenv("ANTHROPIC_API_KEY"))
    return &ClaudeSonnet{
        ctx: ctx,
        client: client,
        system: system + "\n" + shapeMessage,
    }
}

func (o *ClaudeSonnet) chat(chat string, ctx context.Context) (string, error) {

    resp, err := o.client.CreateMessages(ctx, anthropic.MessagesRequest{
        Model: anthropic.ModelClaude3Sonnet20240229,
        Temperature: &temperature,
        MaxTokens: 1000,
        System: o.system,
        Messages: []anthropic.Message{
            anthropic.Message{
                Role: anthropic.RoleUser,
                Content: []anthropic.MessageContent{
                    anthropic.NewTextMessageContent(chat),
                },
            },
        },
    })

    if err != nil {
        return "", err
    }

    return resp.Content[0].GetText(), nil
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


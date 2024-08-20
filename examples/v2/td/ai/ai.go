package ai

import (
	"context"
	"time"

	"github.com/sashabaranov/go-openai"
)

type OpenAIChat struct {
    client *openai.Client
    system string
}

var foo = 1337

func (o *OpenAIChat) chat(chat string, ctx context.Context) (string, error) {
    resp, err := o.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
        Model: openai.GPT4oMini20240718,
        Temperature: 0.55,
        Seed: &foo,
        Messages: []openai.ChatCompletionMessage{
            {
                Role: openai.ChatMessageRoleSystem,
                Content: o.system,
            },
            {
                Role: openai.ChatMessageRoleSystem,
                Content: `Response Format MUST BE line separated Row,Col tuples
Example output with 2 positions specified
20,3
12,4

This specified position of row 20 col 3 and second, separated by new line, row 12 and col 4
`,
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

func NewStatefulOpenAIChat(secret string, system string, ctx context.Context) *StatefulOpenAIChat {
    return &StatefulOpenAIChat{
        ai: NewOpenAIChat(secret, system),
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



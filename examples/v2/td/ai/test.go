package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/liushuangls/go-anthropic/v2"
)

func main() {
    godotenv.Load()

    client := anthropic.NewClient(os.Getenv("ANTHROPIC_API_KEY"))
	resp, err := client.CreateMessages(context.Background(), anthropic.MessagesRequest{
		Model: anthropic.ModelClaude3Sonnet20240229,
		Messages: []anthropic.Message{
			anthropic.NewUserTextMessage("What is your name?"),
		},
		MaxTokens: 1000,
	})

	if err != nil {
		var e *anthropic.APIError
		if errors.As(err, &e) {
			fmt.Printf("Messages error, type: %s, message: %s", e.Type, e.Message)
		} else {
			fmt.Printf("Messages error: %v\n", err)
        }
		return
	}

	fmt.Println(resp.Content[0].GetText())
}


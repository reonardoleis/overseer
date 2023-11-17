package ai

import (
	"context"
	"errors"
	"log"

	openai "github.com/sashabaranov/go-openai"
)

var (
	cli *openai.Client
)

func Init(apiKey string) {
	cli = openai.NewClient(apiKey)
}

func Generate(prompt string) (string, error) {
	resp, err := cli.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     openai.GPT3Dot5Turbo,
			MaxTokens: 200,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		log.Println("chatgpt: error generating text:", err)
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no choices returned")
	}

	return resp.Choices[0].Message.Content, nil
}

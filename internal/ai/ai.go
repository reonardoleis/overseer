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

type MessageContext struct {
	Role    string
	Content string
}

func Init(apiKey string) {
	cli = openai.NewClient(apiKey)
}

func Generate(prompt string, messageContext []MessageContext) (string, error) {
	messages := make([]openai.ChatCompletionMessage, len(messageContext)+1)
	for i, message := range messageContext {
		messages[i] = openai.ChatCompletionMessage{
			Role:    message.Role,
			Content: message.Content,
		}
	}

	messages[len(messageContext)] = openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	}

	resp, err := cli.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     openai.GPT3Dot5Turbo,
			MaxTokens: 500,
			Messages:  messages,
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

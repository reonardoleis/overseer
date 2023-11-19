package ai

import (
	"context"
	"errors"
	"io"
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

func Generate(prompt string, messageContext []MessageContext, maxTokens ...int) (string, error) {
	_maxTokens := 500
	if len(maxTokens) > 0 {
		_maxTokens = maxTokens[0]
	}

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
			MaxTokens: _maxTokens,
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

func TTS(text string) (io.ReadCloser, error) {
	resp, err := cli.CreateSpeech(context.TODO(), openai.CreateSpeechRequest{
		Model: openai.TTSModel1,
		Voice: openai.VoiceAlloy,
		Input: text,
	})

	if err != nil {
		log.Println("chatgpt: error generating text:", err)
		return nil, err
	}

	return resp, nil
}

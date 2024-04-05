package ai

import (
	"bytes"
	"context"
	"encoding/base64"
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
			Model:       openai.GPT3Dot5Turbo,
			MaxTokens:   _maxTokens,
			Messages:    messages,
			Temperature: 0.85,
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

func Transcribe(audio io.Reader) (string, error) {
	resp, err := cli.CreateTranscription(context.TODO(), openai.AudioRequest{
		Reader: audio,
	})

	if err != nil {
		log.Println("chatgpt: error generating text:", err)
		return "", err
	}

	return resp.Text, nil
}

func Image(prompt string) (io.ReadCloser, error) {
	resp, err := cli.CreateImage(context.TODO(), openai.ImageRequest{
		Prompt:         prompt,
		Model:          openai.CreateImageModelDallE3,
		N:              1,
		Quality:        openai.CreateImageQualityStandard,
		Size:           openai.CreateImageSize1024x1024,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
		Style:          openai.CreateImageStyleNatural,
		User:           "overseer",
	})

	if err != nil {
		log.Println("chatgpt: error generating text:", err)
		return nil, err
	}

	imageBase64 := resp.Data[0].B64JSON

	imageBytes, err := base64.StdEncoding.DecodeString(imageBase64)
	if err != nil {
		log.Println("chatgpt: error decoding base64:", err)
		return nil, err
	}

	return io.NopCloser(bytes.NewReader(imageBytes)), nil
}

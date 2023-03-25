package openai

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
	"todo-reminder/repository/bsoncodec"
	"todo-reminder/util"
)

var (
	client *openai.Client
)

func init() {
	if !viper.GetBool("chatgpt.enabled") {
		return
	}
	config := openai.DefaultConfig(os.Getenv("OPENAI_SK"))
	proxyUrl, err := url.Parse(viper.GetString("chatgpt.proxyUrl"))
	if err != nil {
		panic(err)
	}
	config.HTTPClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
		Timeout: time.Minute,
	}
	client = openai.NewClientWithConfig(config)
}

func ChatCompletion(ctx context.Context, input string) (string, error) {
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: input,
			},
		},
	})
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func GenImage(ctx context.Context, input string) (string, string, error) {
	resp, err := client.CreateImage(ctx, openai.ImageRequest{
		Prompt:         input,
		N:              1,
		Size:           openai.CreateImageSize512x512,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
	})
	if err != nil {
		return "", "", err
	}
	fileName := fmt.Sprintf("%s.jpg", bsoncodec.NewObjectId().Hex())
	file, err := os.Create(fileName)
	if err != nil {
		return "", "", err
	}
	defer file.Close()
	bytes, err := base64.StdEncoding.DecodeString(resp.Data[0].B64JSON)
	if err != nil {
		return "", "", err
	}
	_, err = file.Write(bytes)
	if err != nil {
		return "", "", err
	}
	abs, err := filepath.Abs(file.Name())
	return abs, fileName, err
}

func GenImageVariation(ctx context.Context, imagePath string) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	resp, err := client.CreateVariImage(ctx, openai.ImageVariRequest{
		N:     1,
		Size:  openai.CreateImageSize512x512,
		Image: file,
	})
	if err != nil {
		return "", err
	}
	// TODO: use base64 instead
	bytes, err := util.DownloadImage(ctx, resp.Data[0].URL, viper.GetString("chatgpt.proxyUrl"))
	if err != nil {
		return "", err
	}
	file, err = os.Create(fmt.Sprintf("%s.jpg", bsoncodec.NewObjectId().Hex()))
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = file.Write(bytes)
	if err != nil {
		return "", err
	}
	return filepath.Abs(file.Name())
}

func GetResponseWithContext(ctx context.Context, input string, inputs, receivedMessages []string) (string, error) {
	messages := make([]openai.ChatCompletionMessage, 0, len(receivedMessages)+len(inputs)+1)
	for i, _ := range inputs {
		if i >= len(receivedMessages) {
			break
		}
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: inputs[i],
		})
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: receivedMessages[i],
		})
	}
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: input,
	})
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: messages,
	})
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

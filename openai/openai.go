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
	"todo-reminder/log"
	"todo-reminder/repository/bsoncodec"
	"todo-reminder/util"
)

var (
	client *openAIClient
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
	client = &openAIClient{
		client: openai.NewClientWithConfig(config),
	}
}

func GetOpenAIClient() OpenAI {
	if viper.GetBool("chatgpt.enabled") {
		return client
	}
	return openAIEmpty{}
}

type OpenAI interface {
	ChatCompletion(ctx context.Context, input string) (string, error)
	GenImage(ctx context.Context, input string) (string, string, error)
	GenImageVariation(ctx context.Context, imagePath string) (string, error)
	GetResponseWithContext(ctx context.Context, input string, inputs, receivedMessages []string) (string, error)
}

type openAIEmpty struct {
}

func (o openAIEmpty) ChatCompletion(ctx context.Context, input string) (string, error) {
	log.Warn("Calling ChatCompletion", map[string]interface{}{
		"input": input,
	})
	return "", nil
}

func (o openAIEmpty) GenImage(ctx context.Context, input string) (string, string, error) {
	log.Warn("Calling GenImage", map[string]interface{}{
		"input": input,
	})
	return "", "", nil
}

func (o openAIEmpty) GenImageVariation(ctx context.Context, imagePath string) (string, error) {
	log.Warn("Calling GenImageVariation", map[string]interface{}{
		"imagePath": imagePath,
	})
	return "", nil
}

func (o openAIEmpty) GetResponseWithContext(ctx context.Context, input string, inputs, receivedMessages []string) (string, error) {
	log.Warn("Calling GetResponseWithContext", map[string]interface{}{
		"input":            input,
		"inputs":           inputs,
		"receivedMessages": receivedMessages,
	})
	return "", nil
}

type openAIClient struct {
	client *openai.Client
}

func (c *openAIClient) ChatCompletion(ctx context.Context, input string) (string, error) {
	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
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

func (c *openAIClient) GenImage(ctx context.Context, input string) (string, string, error) {
	resp, err := c.client.CreateImage(ctx, openai.ImageRequest{
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

func (c *openAIClient) GenImageVariation(ctx context.Context, imagePath string) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	resp, err := c.client.CreateVariImage(ctx, openai.ImageVariRequest{
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

func (c *openAIClient) GetResponseWithContext(ctx context.Context, input string, inputs, receivedMessages []string) (string, error) {
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
	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: messages,
	})
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

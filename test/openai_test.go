package test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"todo-reminder/openai"
)

func TestGenImage(t *testing.T) {
	_, _, err := openai.GetOpenAIClient().GenImage(context.Background(), "cat")
	assert.NoError(t, err)
}

func TestChatCompletion(t *testing.T) {
	resp, err := openai.GetOpenAIClient().ChatCompletion(context.Background(), "what's csrf attack")
	assert.NoError(t, err)
	log.Println("TestChatCompletion", map[string]interface{}{
		"response": resp,
	})
}

func TestConversation(t *testing.T) {
	ctx := context.Background()
	resp1, err := openai.GetOpenAIClient().GetResponseWithContext(ctx, "list 20 random words", nil, nil)
	assert.NoError(t, err)
	log.Println(resp1)
	resp2, err := openai.GetOpenAIClient().GetResponseWithContext(ctx, "use those words to write a letter", []string{"list 20 random words"}, []string{resp1})
	assert.NoError(t, err)
	log.Println(resp2)
}

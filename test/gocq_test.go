package test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	_ "todo-reminder/conf"
	"todo-reminder/gocq"
)

func TestWSListFriends(t *testing.T) {
	friends, err := gocq.GoCqWebsocket.ListFriends(context.Background())
	assert.NoError(t, err)
	for _, friend := range friends {
		log.Println(friend)
	}
}

func TestWSAtInGroup(t *testing.T) {
	err := gocq.GoCqWebsocket.SendAtInGroup(context.Background(), "484122864", "1658272229", "测试")
	assert.NoError(t, err)
}

func TestWSSendImageInGroup(t *testing.T) {
	err := gocq.GoCqWebsocket.SendGroupImageMessage(context.Background(), "484122864", "test.png", "/home/user/Pictures/test.png")
	assert.NoError(t, err)
}

func TestWSSendPrivateMessage(t *testing.T) {
	err := gocq.GoCqWebsocket.SendPrivateStringMessage(context.Background(), "123", "1658272229")
	assert.NoError(t, err)
}

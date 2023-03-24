package test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	_ "todo-reminder/conf"
	"todo-reminder/gocq"
	"todo-reminder/util"
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

func TestCQCheck(t *testing.T) {
	isCQCode := util.IsCQCode(`[CQ:image,file=e71a2c702348c22f595312aa8532bff0.image,subType=1,url=https://gchat.qpic.cn/gchatpic_new/1658272229/3974122864-2732474974-E71A2C702348C22F595312AA8532BFF0/0?term=2&amp;is_origin=0]`)
	assert.True(t, isCQCode)
	isCQCode = util.IsCQCode(`”@”`)
	assert.False(t, isCQCode)
}

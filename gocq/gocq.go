package gocq

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"todo-reminder/log"
	"todo-reminder/util"
)

type goCq interface {
	ListFriends(ctx context.Context) ([]string, error)
	SendPrivateStringMessage(ctx context.Context, message, userId string) error
	SendGroupImageMessage(ctx context.Context, groupId string, fileName, filePath string) error
	SendAtInGroup(ctx context.Context, groupId, userId string, message string) error
}

type gocqEmpty struct {
}

func (g gocqEmpty) ListFriends(ctx context.Context) ([]string, error) {
	log.Warn("Calling ListFriends", nil)
	return nil, nil
}

func (g gocqEmpty) SendPrivateStringMessage(ctx context.Context, message, userId string) error {
	log.Warn("Calling SendPrivateStringMessage", map[string]interface{}{
		"message": message,
		"userId":  userId,
	})
	return nil
}

func (g gocqEmpty) SendGroupImageMessage(ctx context.Context, groupId string, fileName, filePath string) error {
	log.Warn("Calling SendGroupImageMessage", map[string]interface{}{
		"groupId":  groupId,
		"fileName": fileName,
		"filePath": filePath,
	})
	return nil
}

func (g gocqEmpty) SendAtInGroup(ctx context.Context, groupId, userId string, message string) error {
	log.Warn("Calling SendGroupImageMessage", map[string]interface{}{
		"groupId": groupId,
		"userId":  userId,
		"message": message,
	})
	return nil
}

type goCqHttp struct {
}

const (
	GET_FRIEND_LIST_ENDPOINT      = "get_friend_list"
	SEND_PRIVATE_MESSAGE_ENDPOINT = "send_private_msg"
	SEND_GROUP_MESSAGE_ENDPOINT   = "send_group_msg"
	GET_LOGIN_INFO                = "get_login_info"
)

type BaseResponse[T any] struct {
	Code   int    `json:"retcode"`
	Status string `json:"status"`
	Data   T      `json:"data"`
}

type FriendItem struct {
	NickName string `json:"nickName"`
	Remark   string `json:"remark"`
	UserId   int64  `json:"user_id"`
}

type SendMessageResponse struct {
	MessageId int64 `json:"message_id"`
}

func (g goCqHttp) genUrl(endpoint string) string {
	return fmt.Sprintf("%s/%s", viper.GetString("goCq.uri"), endpoint)
}

func (g goCqHttp) ListFriends(ctx context.Context) ([]string, error) {
	client := util.GetRestClient[BaseResponse[[]FriendItem]]()
	resp, err := client.Get(ctx, g.genUrl(GET_FRIEND_LIST_ENDPOINT), nil, nil)
	if err != nil {
		return nil, err
	}
	userIds := make([]string, 0, len(resp.Data))
	for _, item := range resp.Data {
		userIds = append(userIds, cast.ToString(item.UserId))
	}
	return userIds, nil
}

func (g goCqHttp) SendPrivateStringMessage(ctx context.Context, message, userId string) error {
	client := util.GetRestClient[SendMessageResponse]()
	_, err := client.PostJSON(ctx, g.genUrl(SEND_PRIVATE_MESSAGE_ENDPOINT), nil, map[string]interface{}{
		"user_id":     cast.ToInt64(userId),
		"message":     message,
		"auto_escape": true,
	})
	return err
}

func (g goCqHttp) SendGroupImageMessage(ctx context.Context, groupId string, fileName, filePath string) error {
	client := util.GetRestClient[SendMessageResponse]()
	_, err := client.PostJSON(ctx, g.genUrl(SEND_GROUP_MESSAGE_ENDPOINT), nil, map[string]interface{}{
		"group_id": cast.ToInt64(groupId),
		"message": map[string]interface{}{
			"type": "image",
			"data": map[string]interface{}{
				"file": fileName,
				"url":  util.GenFileURI(filePath),
			},
		},
	})
	return err
}

func (g goCqHttp) SendAtInGroup(ctx context.Context, groupId, userId string, message string) error {
	client := util.GetRestClient[SendMessageResponse]()
	_, err := client.PostJSON(ctx, g.genUrl(SEND_GROUP_MESSAGE_ENDPOINT), nil, map[string]interface{}{
		"group_id": cast.ToInt64(groupId),
		"message":  fmt.Sprintf("[CQ:at,qq=%s] %s", userId, message),
	})
	return err
}

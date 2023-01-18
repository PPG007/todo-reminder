package gocq

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"todo-reminder/util"
)

var (
	GoCq goCq = &goCqHttp{}
)

type goCq interface {
	ListFriends(ctx context.Context) ([]string, error)
}

type goCqHttp struct {
}

const (
	GET_FRIEND_LIST_ENDPOINT = "/get_friend_list"
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

func (g goCqHttp) genUrl(endpoint string) string {
	return fmt.Sprintf("%s%s", viper.GetString("gocq.uri"), endpoint)
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

package gocq

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"todo-reminder/repository/bsoncodec"
	"todo-reminder/util"
)

var (
	conn *websocket.Conn

	GoCqWebsocket goCq
)

const (
	POST_TYPE_MESSAGE    = "message"
	POST_TYPE_REQUEST    = "request"
	POST_TYPE_NOTICE     = "notice"
	POST_TYPE_META_EVENT = "meta_event"

	NOTICE_TYPE_FRIEND_ADD = "friend_add"
)

func init() {
	if viper.GetString("goCq.type") != "ws" {
		return
	}
	url := viper.GetString("goCq.websocketUri")
	action, _, err := websocket.DefaultDialer.DialContext(context.Background(), fmt.Sprintf("%s/api", url), nil)
	if err != nil {
		panic(err)
	}
	event, _, err := websocket.DefaultDialer.DialContext(context.Background(), fmt.Sprintf("%s/event", url), nil)
	if err != nil {
		panic(err)
	}
	g := &goCqWebsocket{
		friendList: make(chan []FriendItem),
		action:     action,
		event:      event,
	}
	go g.listenEventResponse(context.Background())
	go g.listenActionResponse(context.Background())
	GoCqWebsocket = g
}

type WebsocketRequest struct {
	Action string                 `json:"action"`
	Echo   string                 `json:"echo"`
	Params map[string]interface{} `json:"params"`
}

type WebsocketActionResponse struct {
	Status  string                 `json:"status"`
	RetCode int64                  `json:"retCode"`
	Message string                 `json:"msg"`
	Wording string                 `json:"wording"`
	Data    map[string]interface{} `json:"data"`
	Echo    string                 `json:"echo"`
}

type EventBody struct {
	UnixTime    int64  `json:"time"`
	SelfId      int64  `json:"selfId"`
	PostType    string `json:"post_type"`
	NoticeType  string `json:"notice_type"`
	UserId      int64  `json:"user_id"`
	Echo        string `json:"echo"`
	MessageType string `json:"message_type,omitempty"`
	SubType     string `json:"sub_type,omitempty"`
	MessageId   int32  `json:"message_id,omitempty"`
}

type goCqWebsocket struct {
	action     *websocket.Conn
	event      *websocket.Conn
	friendList chan []FriendItem
}

func (g goCqWebsocket) ListFriends(ctx context.Context) ([]string, error) {
	req := WebsocketRequest{
		Action: GET_FRIEND_LIST_ENDPOINT,
		Echo:   GET_FRIEND_LIST_ENDPOINT,
	}
	err := g.action.WriteJSON(req)
	if err != nil {
		return nil, err
	}
	list := <-g.friendList
	var result []string
	for _, item := range list {
		result = append(result, cast.ToString(item.UserId))
	}
	return result, nil
}

func (g goCqWebsocket) SendPrivateStringMessage(ctx context.Context, message, userId string) error {
	req := WebsocketRequest{
		Action: SEND_PRIVATE_MESSAGE_ENDPOINT,
		Echo:   SEND_PRIVATE_MESSAGE_ENDPOINT + bsoncodec.NewObjectId().Hex(),
		Params: map[string]interface{}{
			"user_id":     cast.ToInt64(userId),
			"message":     message,
			"auto_escape": true,
		},
	}
	return g.action.WriteJSON(req)
}

func (g goCqWebsocket) SendGroupImageMessage(ctx context.Context, groupId string, fileName, filePath string) error {
	req := WebsocketRequest{
		Action: SEND_GROUP_MESSAGE_ENDPOINT,
		Echo:   SEND_GROUP_MESSAGE_ENDPOINT + bsoncodec.NewObjectId().Hex(),
		Params: map[string]interface{}{
			"group_id": cast.ToInt64(groupId),
			"message": map[string]interface{}{
				"type": "image",
				"data": map[string]interface{}{
					"file": fileName,
					"url":  util.GenFileURI(filePath),
				},
			},
			"auto_escape": true,
		},
	}
	return g.action.WriteJSON(req)
}

func (g goCqWebsocket) SendAtInGroup(ctx context.Context, groupId, userId string, message string) error {
	req := WebsocketRequest{
		Action: SEND_GROUP_MESSAGE_ENDPOINT,
		Echo:   SEND_GROUP_MESSAGE_ENDPOINT + bsoncodec.NewObjectId().Hex(),
		Params: map[string]interface{}{
			"group_id": cast.ToInt64(groupId),
			"message":  fmt.Sprintf("[CQ:at,qq=%s] %s", userId, message),
		},
	}
	return g.action.WriteJSON(req)
}

func (g goCqWebsocket) listenActionResponse(ctx context.Context) {
	for {
		resp := &WebsocketActionResponse{}
		err := g.action.ReadJSON(resp)
		if err != nil {
			continue
		}
		switch resp.Echo {
		case GET_FRIEND_LIST_ENDPOINT:
			var list []FriendItem
			err := util.CopyByJson(resp.Data, &list)
			if err != nil {
				continue
			}
			g.friendList <- list
		}
	}
}

func (g goCqWebsocket) listenEventResponse(ctx context.Context) {

}

type WebsocketEventResponse struct {
}

func StartWebsocketClient(ctx context.Context) {
	connection, _, err := websocket.DefaultDialer.DialContext(ctx, viper.GetString("goCq.websocketUri"), nil)
	if err != nil {
		panic(err)
	}
	conn = connection
}

func ShutdownWSConnection() {
	conn.Close()
}

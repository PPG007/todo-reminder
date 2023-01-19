package gocq

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

var (
	conn *websocket.Conn
)

const (
	POST_TYPE_MESSAGE    = "message"
	POST_TYPE_REQUEST    = "request"
	POST_TYPE_NOTICE     = "notice"
	POST_TYPE_META_EVENT = "meta_event"

	NOTICE_TYPE_FRIEND_ADD = "friend_add"
)

type goCqMessage interface {
	string
}

type SendPrivateMessageRequest[T goCqMessage] struct {
	UserId     int64 `json:"user_id"`
	AutoEscape bool  `json:"auto_escape"`
	Message    T     `json:"message"`
}

type WebsocketRequest struct {
	Action string      `json:"action"`
	Echo   string      `json:"echo"`
	Params interface{} `json:"params"`
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

type EventBody struct {
	UnixTime   int64  `json:"time"`
	SelfId     int64  `json:"selfId"`
	PostType   string `json:"post_type"`
	NoticeType string `json:"notice_type"`
	UserId     int64  `json:"user_id"`
	Echo       string `json:"echo"`
}

func (e *EventBody) Handle(ctx context.Context) {
	switch e.PostType {
	case POST_TYPE_NOTICE:
		e.HandleNoticeEvent(ctx)
	}
}

func (e *EventBody) HandleNoticeEvent(ctx context.Context) {
	switch e.NoticeType {
	case NOTICE_TYPE_FRIEND_ADD:
		//user := model.User{
		//	UserId: cast.ToString(e.UserId),
		//}
		//user.UpsertWithoutPassword(ctx)
	}
}

func Listen(ctx context.Context) {
	for {
		eventBody := EventBody{}
		err := conn.ReadJSON(&eventBody)
		if err != nil {
			continue
		}
		eventBody.Handle(ctx)
	}

}

func SendPrivateMessage[T goCqMessage](ctx context.Context, req SendPrivateMessageRequest[T]) error {
	wsreq := WebsocketRequest{
		Action: SEND_PRIVATE_MESSAGE_ENDPOINT,
		Params: req,
	}
	return conn.WriteJSON(wsreq)
}

package gocq

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"strings"
	"time"
	"todo-reminder/log"
	"todo-reminder/repository/bsoncodec"
	"todo-reminder/util"
)

var (
	GoCqWebsocket goCq
)

const (
	POST_TYPE_MESSAGE      = "message"
	POST_TYPE_MESSAGE_SENT = "message_sent"
	POST_TYPE_REQUEST      = "request"
	POST_TYPE_NOTICE       = "notice"
	POST_TYPE_META_EVENT   = "meta_event"

	META_EVENT_TYPE_HEART_BEAT = "heartbeat"
	META_EVENT_TYPE_LIFE_CYCLE = "lifecycle"

	NOTICE_TYPE_FRIEND_ADD = "friend_add"

	MESSAGE_TYPE_PRIVATE = "private"
	MESSAGE_TYPE_GROUP   = "group"

	MESSAGE_SUB_TYPE_FRIEND = "friend"
	MESSAGE_SUB_TYPE_NORMAL = "normal"
)

type goCqWebsocket struct {
	action        *websocket.Conn
	event         *websocket.Conn
	friendList    chan []FriendItem
	heartBeat     chan HeartBeatStatus
	lastAlertTime time.Time
}

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
		heartBeat:  make(chan HeartBeatStatus),
		action:     action,
		event:      event,
	}
	go g.listenEventResponse(context.Background())
	go g.listenActionResponse(context.Background())
	go g.HeartBeat(context.Background())
	GoCqWebsocket = g
}

type WebsocketRequest struct {
	Action string                 `json:"action"`
	Echo   string                 `json:"echo"`
	Params map[string]interface{} `json:"params"`
}

type WebsocketActionResponse struct {
	Status  string      `json:"status"`
	RetCode int64       `json:"retCode"`
	Message string      `json:"msg"`
	Wording string      `json:"wording"`
	Data    interface{} `json:"data"`
	Echo    string      `json:"echo"`
}

type EventBody struct {
	UnixTime        int64           `json:"time"`
	SelfId          int64           `json:"selfId"`
	PostType        string          `json:"post_type"`
	MetaEventType   string          `json:"meta_event_type,omitempty"`
	HeartBeatStatus HeartBeatStatus `json:"status,omitempty"`
	NoticeType      string          `json:"notice_type,omitempty"`
	UserId          int64           `json:"user_id,omitempty"`
	MessageType     string          `json:"message_type,omitempty"`
	SubType         string          `json:"sub_type,omitempty"`
	MessageId       int64           `json:"message_id"`
	RawMessage      string          `json:"raw_message,omitempty"`
	Sender          Sender          `json:"sender,omitempty"`
	GroupId         int64           `json:"group_id,omitempty"`
}

type HeartBeatStatus struct {
	AppInitialized bool `json:"app_initialized"`
	AppEnabled     bool `json:"app_enabled"`
	AppGood        bool `json:"app_good"`
	Online         bool `json:"online"`
}

type Sender struct {
	UserId   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Sex      string `json:"sex"`
	Age      int64  `json:"age"`
	GroupId  int64  `json:"group_id"`
	Card     string `json:"card"`
	Area     string `json:"area"`
	Level    string `json:"level"`
	Role     string `json:"role"`
	Title    string `json:"title"`
}

func (s *HeartBeatStatus) Check() bool {
	return s.AppInitialized && s.AppEnabled && s.AppGood && s.Online
}

func (g *goCqWebsocket) ListFriends(ctx context.Context) ([]string, error) {
	req := WebsocketRequest{
		Action: GET_FRIEND_LIST_ENDPOINT,
		Echo:   GET_FRIEND_LIST_ENDPOINT,
	}
	err := g.action.WriteJSON(req)
	if err != nil {
		return nil, err
	}
	timer := time.NewTimer(time.Second * 3)
	for {
		select {
		case list := <-g.friendList:
			var result []string
			for _, item := range list {
				result = append(result, cast.ToString(item.UserId))
			}
			return result, nil
		case <-timer.C:
			return nil, errors.New("context deadline exceed")
		}
	}
}

func (g *goCqWebsocket) HeartBeat(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case status := <-g.heartBeat:
			if status.Check() {
				ticker.Reset(time.Second * 10)
			} else {
				log.Warn("Check heart beat failed", map[string]interface{}{
					"status": status,
				})
			}
		case <-ticker.C:
			log.Warn("HeartBeat deadline exceeded", map[string]interface{}{})
			if time.Now().Sub(g.lastAlertTime) >= time.Hour {
				util.SendEmail(ctx, viper.GetString("alert.receiver"), "Alert", "check gocq heart beat failed")
				g.lastAlertTime = time.Now()
			}
		}
	}
}

func (g *goCqWebsocket) SendPrivateStringMessage(ctx context.Context, message, userId string) error {
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

func (g *goCqWebsocket) SendGroupImageMessage(ctx context.Context, groupId string, fileName, filePath string) error {
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

func (g *goCqWebsocket) SendAtInGroup(ctx context.Context, groupId, userId string, message string) error {
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

func (g *goCqWebsocket) listenActionResponse(ctx context.Context) {
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

func (g *goCqWebsocket) listenEventResponse(ctx context.Context) {
	for {
		event := &EventBody{}
		err := g.event.ReadJSON(event)
		if err != nil {
			log.Warn("Failed to read event", map[string]interface{}{
				"error": err.Error(),
			})
			continue
		}
		util.Submit(func() {
			err := event.handleEvent(ctx, g)
			if err != nil {
				log.Warn("Failed to handle event", map[string]interface{}{
					"error": err.Error(),
				})
			}
		})
	}
}

func (e *EventBody) handleEvent(ctx context.Context, ws *goCqWebsocket) error {
	switch e.PostType {
	case POST_TYPE_MESSAGE:
		return e.handleMessageEvent(ctx, ws)
	case POST_TYPE_MESSAGE_SENT:
		return e.handleMessageSentEvent(ctx, ws)
	case POST_TYPE_REQUEST:
		return e.handleRequestEvent(ctx, ws)
	case POST_TYPE_NOTICE:
		return e.handleNoticeEvent(ctx, ws)
	case POST_TYPE_META_EVENT:
		return e.handleMetaInfoEvent(ctx, ws)
	default:
		return errors.New("unsupported event post type")
	}
}

func (e *EventBody) handleMessageEvent(ctx context.Context, ws *goCqWebsocket) error {
	switch e.MessageType {
	case MESSAGE_TYPE_PRIVATE:

	case MESSAGE_TYPE_GROUP:
		if !util.IsCQCode(e.RawMessage) {
			return nil
		}
		params, prefix, suffix := util.GetCQCodeParams(e.RawMessage)
		if params["type"] != "at" {
			return nil
		}
		content := prefix
		if content == "" {
			content = suffix
		}
		if content == "" {
			return nil
		}
		if strings.Contains(content, "图片") {
			absPath, name, rpcErr := util.GetImageFile(ctx, content)
			if rpcErr != nil {
				return ws.SendAtInGroup(ctx, cast.ToString(e.GroupId), cast.ToString(e.UserId), rpcErr.Error())
			}
			return ws.SendGroupImageMessage(ctx, cast.ToString(e.GroupId), name, absPath)
		} else {
			message, err := util.GetTextResponse(ctx, content)
			if err != nil {
				message = err.Error()
			}
			return ws.SendAtInGroup(ctx, cast.ToString(e.GroupId), cast.ToString(e.UserId), message)
		}
	default:
		return errors.New("unsupported message type")
	}
	return nil
}

func (e *EventBody) handleMessageSentEvent(ctx context.Context, ws *goCqWebsocket) error {
	return nil
}

func (e *EventBody) handleRequestEvent(ctx context.Context, ws *goCqWebsocket) error {
	return nil
}

func (e *EventBody) handleNoticeEvent(ctx context.Context, ws *goCqWebsocket) error {
	return nil
}

func (e *EventBody) handleMetaInfoEvent(ctx context.Context, ws *goCqWebsocket) error {
	switch e.MetaEventType {
	case META_EVENT_TYPE_HEART_BEAT:
		ws.heartBeat <- e.HeartBeatStatus
	case META_EVENT_TYPE_LIFE_CYCLE:
	}
	return nil
}

type WebsocketEventResponse struct {
}

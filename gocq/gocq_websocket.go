package gocq

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"strings"
	"sync"
	"time"
	"todo-reminder/log"
	"todo-reminder/model"
	"todo-reminder/openai"
	"todo-reminder/repository/bsoncodec"
	"todo-reminder/util"
)

var (
	goCqWs *goCqWebsocket
)

func GetGocqInstance() goCq {
	t := viper.GetString("goCq.type")
	switch t {
	case "ws":
		return goCqWs
	case "http":
		return &goCqHttp{}
	default:
		return gocqEmpty{}
	}
}

const (
	POST_TYPE_MESSAGE      = "message"
	POST_TYPE_MESSAGE_SENT = "message_sent"
	POST_TYPE_REQUEST      = "request"
	POST_TYPE_NOTICE       = "notice"
	POST_TYPE_META_EVENT   = "meta_event"

	META_EVENT_TYPE_HEART_BEAT = "heartbeat"
	META_EVENT_TYPE_LIFE_CYCLE = "lifecycle"

	MESSAGE_TYPE_PRIVATE = "private"
	MESSAGE_TYPE_GROUP   = "group"

	MESSAGE_SUB_TYPE_FRIEND = "friend"
	MESSAGE_SUB_TYPE_GROUP  = "group"
	MESSAGE_SUB_TYPE_NORMAL = "normal"

	REQUEST_TYPE_FRIEND = "friend"
	REQUEST_TYPE_GROUP  = "group"

	NOTICE_TYPE_FRIEND_ADDED = "friend_add"
)

type goCqWebsocket struct {
	action              *websocket.Conn
	event               *websocket.Conn
	friendList          chan []FriendItem
	messages            sync.Map
	self                LoginInfo
	heartBeat           chan HeartBeatStatus
	lastAlertTime       time.Time
	conversations       map[int64][]Conversation
	lastReceivedTimeMap map[int64]time.Time
	lock                *sync.Mutex
	admin               model.User
}

type Conversation struct {
	input    string
	response string
}

func init() {
	if viper.GetString("goCq.type") != "ws" {
		return
	}
	g := &goCqWebsocket{
		friendList:          make(chan []FriendItem),
		heartBeat:           make(chan HeartBeatStatus),
		conversations:       make(map[int64][]Conversation),
		lastReceivedTimeMap: make(map[int64]time.Time),
		lock:                &sync.Mutex{},
		messages:            sync.Map{},
	}
	admin, err := model.CUser.GetAdmin(context.Background())
	if err != nil {
		panic(err)
	}
	g.admin = admin
	err = g.dial(context.Background())
	if err != nil {
		panic(err)
	}
	go g.listenEventResponse(context.Background())
	go g.listenActionResponse(context.Background())
	go g.HeartBeat(context.Background())
	go g.InitSelfInfo(context.Background())
	goCqWs = g
}

func (g *goCqWebsocket) dial(ctx context.Context) error {
	url := viper.GetString("goCq.websocketUri")
	action, _, err := websocket.DefaultDialer.DialContext(context.Background(), fmt.Sprintf("%s/api", url), nil)
	if err != nil {
		return err
	}
	event, _, err := websocket.DefaultDialer.DialContext(context.Background(), fmt.Sprintf("%s/event", url), nil)
	if err != nil {
		return err
	}
	g.action = action
	g.event = event
	return nil
}

func (g *goCqWebsocket) close() error {
	err := g.action.Close()
	if err != nil {
		return err
	}
	return g.event.Close()
}

func (g *goCqWebsocket) retry() {
	g.close()
	ctx := context.Background()
	for {
		if err := g.dial(ctx); err != nil {
			log.Warn("Failed to connect gocq websocket, retrying...", map[string]interface{}{
				"error": err.Error(),
			})
			time.Sleep(time.Second * 5)
			continue
		}
		break
	}
	log.Warn("Gocq websocket connected", nil)
	go g.listenActionResponse(ctx)
	go g.listenEventResponse(ctx)
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
	RequestType     string          `json:"request_type,omitempty"`
	Flag            string          `json:"flag,omitempty"`
	Comment         string          `json:"comment,omitempty"`
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

type MessageDetail struct {
	IsGroup     bool   `json:"group"`
	GroupId     int64  `json:"group_id"`
	MessageId   int64  `json:"message_id"`
	RealId      int64  `json:"real_id"`
	MessageType string `json:"message_type"`
	Sender      Sender `json:"sender"`
	Time        int64  `json:"time"`
	Message     string `json:"message"`
	RawMessage  string `json:"raw_message"`
}

type LoginInfo struct {
	UserId int64 `json:"user_id"`
}

func (s *HeartBeatStatus) Check() bool {
	return s.AppInitialized && s.AppEnabled && s.AppGood && s.Online
}

func (g *goCqWebsocket) deleteFriend(ctx context.Context, userId string) error {
	req := WebsocketRequest{
		Action: DELETE_FRIEND_ENDPOINT,
		Params: map[string]interface{}{
			"user_id": cast.ToInt64(userId),
		},
	}
	return g.action.WriteJSON(req)
}

func (g *goCqWebsocket) handleFriendRequest(ctx context.Context, userId string, approve bool) error {
	waitingUser, err := model.CWaitingUser.GetWaitingOne(ctx, userId)
	if err != nil {
		return err
	}
	req := WebsocketRequest{
		Action: HANDLE_FRIEND_ENDPOINT,
		Params: map[string]interface{}{
			"flag":    waitingUser.Flag,
			"approve": approve,
		},
	}
	return g.action.WriteJSON(req)
}

func (g *goCqWebsocket) getMessageDetail(ctx context.Context, messageId int64) (MessageDetail, error) {
	req := WebsocketRequest{
		Action: GET_MESSAGE_DETAIL_ENDPOINT,
		Echo:   GET_MESSAGE_DETAIL_ENDPOINT,
		Params: map[string]interface{}{
			"message_id": messageId,
		},
	}
	err := g.action.WriteJSON(req)
	if err != nil {
		return MessageDetail{}, err
	}
	start := time.Now()
	for {
		value, ok := g.messages.LoadAndDelete(messageId)
		if !ok {
			if time.Now().Sub(start) > time.Second*3 {
				return MessageDetail{}, errors.New("deadline exceed")
			}
			continue
		}
		return value.(MessageDetail), nil
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

func (g *goCqWebsocket) ListFriends(ctx context.Context) ([]string, error) {
	resp, err := g.ListFriendsWithFullInfo(ctx)
	if err != nil {
		return nil, err
	}
	userIds := make([]string, 0, len(resp))
	for _, item := range resp {
		userIds = append(userIds, cast.ToString(item.UserId))
	}
	return userIds, nil
}

func (g *goCqWebsocket) ListFriendsWithFullInfo(ctx context.Context) ([]FriendItem, error) {
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
			return list, nil
		case <-timer.C:
			return nil, errors.New("context deadline exceed")
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

func (g *goCqWebsocket) sendPrivateStringMessageWithGroup(ctx context.Context, message, userId, groupId string) error {
	req := WebsocketRequest{
		Action: SEND_PRIVATE_MESSAGE_ENDPOINT,
		Echo:   SEND_PRIVATE_MESSAGE_ENDPOINT + bsoncodec.NewObjectId().Hex(),
		Params: map[string]interface{}{
			"user_id":     cast.ToInt64(userId),
			"group_id":    cast.ToInt64(groupId),
			"message":     message,
			"auto_escape": true,
		},
	}
	return g.action.WriteJSON(req)
}

func (g *goCqWebsocket) SendPrivateImageMessage(ctx context.Context, userId string, fileName, fileUrl string) error {
	req := WebsocketRequest{
		Action: SEND_PRIVATE_MESSAGE_ENDPOINT,
		Params: map[string]interface{}{
			"user_id": cast.ToInt64(userId),
			"message": map[string]interface{}{
				"type": "image",
				"data": map[string]interface{}{
					"file": fileName,
					"url":  fileUrl,
				},
			},
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

func (g *goCqWebsocket) InitSelfInfo(ctx context.Context) error {
	req := WebsocketRequest{
		Action: GET_LOGIN_INFO,
		Echo:   GET_LOGIN_INFO,
	}
	return g.action.WriteJSON(req)
}

func (g *goCqWebsocket) listenActionResponse(ctx context.Context) {
	for {
		resp := &WebsocketActionResponse{}
		err := g.action.ReadJSON(resp)
		if err != nil {
			if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
				return
			}
			log.Warn("Failed to read action response", map[string]interface{}{
				"error": err.Error(),
			})
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
		case GET_LOGIN_INFO:
			loginInfo := LoginInfo{}
			err := util.CopyByJson(resp.Data, &loginInfo)
			if err != nil {
				continue
			}
			g.self = loginInfo
		case GET_MESSAGE_DETAIL_ENDPOINT:
			detail := MessageDetail{}
			err := util.CopyByJson(resp.Data, &detail)
			if err != nil {
				continue
			}
			g.messages.Store(detail.MessageId, detail)
		}
	}
}

func (g *goCqWebsocket) listenEventResponse(ctx context.Context) {
	for {
		event := &EventBody{}
		err := g.event.ReadJSON(event)
		if err != nil {
			if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
				go g.retry()
				log.Warn("Gocq websocket connection closed", nil)
				return
			}
			log.Warn("Failed to read event response", map[string]interface{}{
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
		return e.handlePrivateMessage(ctx, ws)
	case MESSAGE_TYPE_GROUP:
		return e.handleGroupMessage(ctx, ws)
	default:
		return errors.New("unsupported message type")
	}
}

func (e *EventBody) handlePrivateMessage(ctx context.Context, ws *goCqWebsocket) error {
	if e.SubType == MESSAGE_SUB_TYPE_GROUP {
		return ws.sendPrivateStringMessageWithGroup(ctx, "please add friend first", cast.ToString(e.Sender.UserId), cast.ToString(e.Sender.GroupId))
	}
	if e.SubType != MESSAGE_SUB_TYPE_FRIEND {
		return nil
	}
	if util.IsCQCode(e.RawMessage) {
		return ws.SendPrivateStringMessage(ctx, "only accept plain text private message", cast.ToString(e.Sender.UserId))
	}
	if IsCommand(e.RawMessage) {
		return e.handleCommand(ctx, ws)
	}
	user, err := model.CUser.GetByUserId(ctx, cast.ToString(e.Sender.UserId))
	if err != nil {
		return ws.SendPrivateStringMessage(ctx, err.Error(), cast.ToString(e.Sender.UserId))
	}
	// check if use is admin or has openAI approved
	if !user.OpenAIApproved && ws.admin.UserId != user.UserId {
		return ws.SendPrivateStringMessage(ctx, fmt.Sprintf("permission denied, please contact the admin %s(%s)", ws.admin.UserId, ws.admin.Nickname), cast.ToString(e.Sender.UserId))
	}
	messageDetail, err := ws.getMessageDetail(ctx, e.MessageId)
	if err != nil {
		return ws.SendPrivateStringMessage(ctx, err.Error(), cast.ToString(e.Sender.UserId))
	}
	e.replyMessage(ctx, "processing...", messageDetail, ws)
	resp, err := openai.GetOpenAIClient().GetStreamResponse(ctx, e.RawMessage)
	if err != nil {
		resp = err.Error()
	}
	return e.replyMessage(ctx, resp, messageDetail, ws)
}

func (e *EventBody) handleCommand(ctx context.Context, ws *goCqWebsocket) error {
	messageDetail, err := ws.getMessageDetail(ctx, e.MessageId)
	if err != nil {
		return err
	}
	messageDetail.RawMessage = e.RawMessage
	if cast.ToString(e.Sender.UserId) != ws.admin.UserId {
		return e.replyMessage(ctx, "only admin can use command", messageDetail, ws)
	}
	if util.IsCQCode(messageDetail.RawMessage) {
		return e.replyMessage(ctx, "command must use in plain text", messageDetail, ws)
	}
	cmd, err := ParseCommand(messageDetail.RawMessage)
	if err != nil {
		return e.replyMessage(ctx, err.Error(), messageDetail, ws)
	}
	reply, err := cmd.Run(ctx)
	if err != nil {
		reply = err.Error()
	}
	return e.replyMessage(ctx, reply, messageDetail, ws)
}

func (e *EventBody) sendPrivateResponse(ctx context.Context, ws *goCqWebsocket, message string) error {
	return ws.SendPrivateStringMessage(ctx, message, cast.ToString(e.Sender.UserId))
}

func (e *EventBody) handleGroupMessage(ctx context.Context, ws *goCqWebsocket) error {
	if !util.IsCQCode(e.RawMessage) {
		return nil
	}
	params, prefix, suffix := util.GetCQCodeParams(e.RawMessage)
	if params["type"] != "at" {
		return nil
	}
	if cast.ToInt64(params["qq"]) != ws.self.UserId {
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
		absPath, fileName, err := openai.GetOpenAIClient().GenImage(ctx, strings.Join(strings.Split(content, "图片"), ""))
		if err != nil {
			return ws.SendAtInGroup(ctx, cast.ToString(e.GroupId), cast.ToString(e.UserId), err.Error())
		}
		return ws.SendGroupImageMessage(ctx, cast.ToString(e.GroupId), fileName, absPath)
	} else {
		ws.lock.Lock()
		lastReceivedTime := ws.lastReceivedTimeMap[e.UserId]
		if time.Now().Sub(lastReceivedTime) > time.Hour {
			delete(ws.conversations, e.UserId)
		}
		ws.lastReceivedTimeMap[e.UserId] = time.Now()
		ws.lock.Unlock()
		conversations := ws.conversations[e.UserId]
		inputs, receivedMessages := transConversationsToSlice(conversations)
		message, err := openai.GetOpenAIClient().GetResponseWithContext(ctx, content, inputs, receivedMessages)
		if err != nil {
			message = err.Error()
			ws.lock.Lock()
			delete(ws.conversations, e.UserId)
			ws.lock.Unlock()
		} else {
			conversations = append(conversations, Conversation{
				input:    content,
				response: message,
			})
			ws.lock.Lock()
			ws.conversations[e.UserId] = conversations
			ws.lock.Unlock()
		}
		return ws.SendAtInGroup(ctx, cast.ToString(e.GroupId), cast.ToString(e.UserId), message)
	}
}

func (e *EventBody) handleMessageSentEvent(ctx context.Context, ws *goCqWebsocket) error {
	return nil
}

func (e *EventBody) handleRequestEvent(ctx context.Context, ws *goCqWebsocket) error {
	switch e.RequestType {
	case REQUEST_TYPE_FRIEND:
		user, err := model.CWaitingUser.Upsert(ctx, cast.ToString(e.UserId), e.Flag)
		if err != nil {
			return err
		}
		if user.Notified {
			return nil
		}
		user.MarkAsNotified(ctx)
		return ws.SendPrivateStringMessage(ctx, fmt.Sprintf("You have a friend request from %d with comment: %s", e.UserId, e.Comment), ws.admin.UserId)
	case REQUEST_TYPE_GROUP:
		return nil
	}
	return nil
}

func (e *EventBody) handleNoticeEvent(ctx context.Context, ws *goCqWebsocket) error {
	switch e.NoticeType {
	case NOTICE_TYPE_FRIEND_ADDED:
		model.CUser.HandleFriendAdded(ctx, cast.ToString(e.UserId))
		model.CWaitingUser.HandleFriendAdded(ctx, cast.ToString(e.UserId))
		return nil
	}
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

func (e *EventBody) replyMessage(ctx context.Context, replyText string, message MessageDetail, ws *goCqWebsocket) error {
	action := SEND_PRIVATE_MESSAGE_ENDPOINT
	if message.IsGroup {
		action = SEND_GROUP_MESSAGE_ENDPOINT
	}
	req := WebsocketRequest{
		Action: action,
		Params: map[string]interface{}{
			"user_id": message.Sender.UserId,
			"message": getReplayCQCode(message.MessageId, message.RealId, replyText),
		},
	}
	return ws.action.WriteJSON(req)
}

func getReplayCQCode(messageId, seq int64, replay string) string {
	return fmt.Sprintf("[CQ:reply,id=%d,seq=%d] %s", messageId, seq, replay)
}

func isReplay(message string) bool {
	params, _, _ := util.GetCQCodeParams(message)
	return params["type"] == "reply"
}

func transConversationsToSlice(conversations []Conversation) ([]string, []string) {
	s1, s2 := make([]string, 0, len(conversations)), make([]string, 0, len(conversations))
	for _, conversation := range conversations {
		s1 = append(s1, conversation.input)
		s2 = append(s2, conversation.response)
	}
	return s1, s2
}

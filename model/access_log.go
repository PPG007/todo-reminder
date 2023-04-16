package model

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"io"
	"net"
	"strings"
	"time"
	"todo-reminder/constant"
	"todo-reminder/repository"
	"todo-reminder/repository/bsoncodec"
	"todo-reminder/util"
	"unicode/utf8"
)

const (
	C_ACCESS_LOG = "accessLog"
)

var (
	CAccessLog = &AccessLog{}
)

type AccessLog struct {
	Id            bsoncodec.ObjectId `bson:"_id"`
	Body          string             `bson:"body"`
	BodySize      int64              `bson:"bodySize"`
	Method        string             `bson:"method"`
	URL           string             `bson:"URL"`
	RemoteAddress string             `bson:"remoteAddress,omitempty"`
	RemotePort    string             `bson:"remotePort,omitempty"`
	Referer       string             `bson:"referer,omitempty"`
	StatusCode    int                `bson:"statusCode"`
	RequestId     string             `bson:"requestId"`
	UserId        string             `bson:"userId,omitempty"`
	UserAgent     string             `bson:"userAgent"`
	StartTime     time.Time          `bson:"startTime"`
	EndTime       time.Time          `bson:"endTime"`
	ResponseBody  string             `bson:"responseBody,omitempty"`
}

func (*AccessLog) Init(ctx *gin.Context) AccessLog {
	buf := bytes.Buffer{}
	size, _ := buf.ReadFrom(ctx.Request.Body)
	ctx.Request.Body = io.NopCloser(&buf)
	host, port, _ := net.SplitHostPort(strings.TrimSpace(ctx.Request.RemoteAddr))
	return AccessLog{
		Id: bsoncodec.NewObjectId(),
		Body: func() string {
			if utf8.Valid(buf.Bytes()) {
				return buf.String()
			}
			return ""
		}(),
		BodySize:      size,
		Method:        ctx.Request.Method,
		URL:           ctx.Request.URL.RequestURI(),
		RemoteAddress: host,
		RemotePort:    port,
		Referer:       ctx.Request.Referer(),
		RequestId:     util.ExtractRequestId(ctx),
		UserId:        util.ExtractUserId(ctx),
		UserAgent:     ctx.Request.UserAgent(),
		StartTime:     time.Now(),
	}
}

func (log *AccessLog) Record(ctx *gin.Context) {
	log.EndTime = time.Now()
	log.StatusCode = ctx.Writer.Status()
	value, exists := ctx.Get(constant.GIN_KEY_RESPONSE_BODY)
	if exists {
		log.ResponseBody = cast.ToString(value)
	}
	log.Create(ctx)
}

func (log *AccessLog) Create(ctx context.Context) error {
	return repository.Mongo.Insert(ctx, C_ACCESS_LOG, log)
}

package util

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
	"todo-reminder/constant"
)

var (
	weekdayMap = map[time.Weekday]int{
		time.Monday:    1,
		time.Tuesday:   2,
		time.Wednesday: 3,
		time.Thursday:  4,
		time.Friday:    5,
		time.Saturday:  6,
		time.Sunday:    7,
	}

	randomStringPool = []rune{
		'1', '2', '3', '4', '5', '6', '7', '8', '9', '0',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'G', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'g', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	}
)

func StrInArray(str string, arr *[]string) bool {
	if arr == nil {
		return false
	}
	for _, s := range *arr {
		if str == s {
			return true
		}
	}
	return false
}

func GetStartTimeOfYear(argTime time.Time) time.Time {
	y, _, _ := argTime.Date()
	return time.Date(y, 1, 1, 0, 0, 0, 0, time.Local)
}

func GetEndTimeOfYear(argTime time.Time) time.Time {
	y, _, _ := argTime.Date()
	return time.Date(y+1, 1, 1, 23, 59, 59, 999999999, time.Local).AddDate(0, 0, -1)
}

func GetStartTimeOfMonth(argTime time.Time) time.Time {
	y, m, _ := argTime.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, time.Local)
}

func GetEndTimeOfMonth(argTime time.Time) time.Time {
	y, m, _ := argTime.Date()
	return time.Date(y, m+1, 0, 23, 59, 59, 999999999, time.Local)
}

func PtrValue[T any](value T) *T {
	v := new(T)
	v = &value
	return v
}

func GetNextWeekday(t time.Time, weekday int) time.Time {
	for weekdayMap[t.Weekday()] != weekday {
		t = t.AddDate(0, 0, 1)
	}
	return t
}

func GenRandomString(length int) string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	result := ""
	for i := 0; i < length; i++ {
		result += string(randomStringPool[r.Intn(len(randomStringPool))])
	}
	return result
}

func GenFileURI(path string) string {
	return fmt.Sprintf("file:///%s", path)
}

func CopyByJson(src, dst interface{}) error {
	bytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dst)
}

func IsCQCode(message string) bool {
	compile := regexp.MustCompile(`\[CQ:(.*?)]`)
	return len(compile.FindStringSubmatch(message)) >= 2
}

func GetCQCodeParams(rawMessage string) (map[string]string, string, string) {
	strs := regexp.MustCompile(`\[CQ:(.*?)]`).FindStringSubmatch(rawMessage)
	if len(strs) < 2 {
		return nil, "", ""
	}
	paramPairs := strings.Split(strs[1], ",")
	params := map[string]string{
		"type": paramPairs[0],
	}
	for i, pair := range paramPairs {
		if i == 0 {
			continue
		}
		kv := strings.Split(pair, "=")
		params[kv[0]] = strings.Join(kv[1:], "=")
	}
	index := strings.Index(rawMessage, strs[1])
	prefix := rawMessage[0 : index-4]
	suffix := rawMessage[index+len(strs[1])+1:]
	return params, strings.TrimSpace(prefix), strings.TrimSpace(suffix)
}

func DownloadImage(ctx context.Context, url, proxy string) ([]byte, error) {
	req := gorequest.New()
	resp, bytes, errs := req.Proxy(proxy).Get(url).EndBytes()
	if len(errs) > 0 {
		return nil, errs[0]
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unknown error")
	}
	return bytes, nil
}

func GetAllCQParams(rawMessage string) ([]map[string]string, string) {
	var (
		result    []map[string]string
		plainText string
	)
	for IsCQCode(rawMessage) {
		params, prefix, suffix := GetCQCodeParams(rawMessage)
		rawMessage = suffix
		result = append(result, params)
		plainText = suffix
		if plainText == "" {
			plainText = prefix
		}
	}
	return result, plainText
}

func ExtractRequestId(ctx context.Context) string {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		return ginCtx.GetHeader(constant.HEADER_REQUEST_ID)
	}
	return ""
}

func ExtractUserId(ctx context.Context) string {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		return ginCtx.GetString(constant.GIN_KEY_USER_ID)
	}
	return ""
}

package util

import (
	"context"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"net/url"
)

func GetRestClient[T any]() httpClient[T] {
	return goRequester[T]{}
}

type httpClient[T any] interface {
	Get(ctx context.Context, url string, headers map[string]string, params *url.Values) (T, error)
	PostJSON(ctx context.Context, url string, headers map[string]string, params map[string]interface{}) (T, error)
}

type goRequester[T any] struct {
}

func (g goRequester[T]) Get(ctx context.Context, url string, headers map[string]string, params *url.Values) (T, error) {
	req := gorequest.New()
	for k, v := range headers {
		req = req.Set(k, v)
	}
	if params != nil {
		url = fmt.Sprintf("%s?%s", url, params.Encode())
	}
	result := new(T)
	_, _, errs := req.Get(url).EndStruct(result)
	if len(errs) > 0 {
		return *result, errs[0]
	}
	return *result, nil
}

func (g goRequester[T]) PostJSON(ctx context.Context, url string, headers map[string]string, params map[string]interface{}) (T, error) {
	req := gorequest.New()
	for k, v := range headers {
		req = req.Set(k, v)
	}
	result := new(T)
	jsonBody := MarshalToJson(params)
	_, _, errs := req.Post(url).Send(jsonBody).EndStruct(result)
	if len(errs) > 0 {
		return *result, errs[0]
	}
	return *result, nil
}

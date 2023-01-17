package cron

import (
	"context"
	"github.com/spf13/cast"
	"net/url"
	"time"
	"todo-reminder/model"
	"todo-reminder/util"
)

const (
	TimorApi = "https://timor.tech/api/holiday/year/%d/"

	ZhiHuApi = "https://api.apihubs.cn/holiday/get"
)

// TimorHolidayResponse https://timor.tech/api/holiday
type TimorHolidayResponse struct {
	Code    int                         `json:"code"`
	Holiday map[string]TimorHolidayInfo `json:"holiday"`
}

type TimorHolidayInfo struct {
	IsHoliday bool   `json:"holiday"`
	Name      string `json:"name"`
	Date      string `json:"date"`
}

// ZhiHuHolidayResponse https://zhuanlan.zhihu.com/p/343863993
type ZhiHuHolidayResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"msg"`
	Data    ZhiHuHolidayData `json:"data"`
}

type ZhiHuHolidayData struct {
	Page  int                `json:"page"`
	Size  int                `json:"size"`
	Total int                `json:"total"`
	List  []ZhiHuHolidayInfo `json:"list"`
}

type ZhiHuHolidayInfo struct {
	Date    int `json:"date"`
	WorkDay int `json:"workDay"`
}

func RefreshHoliday() {
	ctx := context.Background()
	year := time.Now().Year()
	if time.Now().Month() == time.December {
		year++
	}
	req := &url.Values{}
	req.Set("year", cast.ToString(year))
	req.Set("size", "400")
	resp, err := util.GetRestClient[ZhiHuHolidayResponse]().Get(ctx, ZhiHuApi, nil, req)
	if err != nil || resp.Code != 0 {
		return
	}
	for _, info := range resp.Data.List {
		holiday := &model.ChinaHoliday{
			IsWorkingDay: info.WorkDay == 1,
			Date:         cast.ToString(info.Date),
		}
		holiday.Create(ctx)
	}
}

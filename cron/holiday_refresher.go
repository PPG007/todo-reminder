package cron

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
	"todo-reminder/model"
	"todo-reminder/util"
)

const (
	TimorApi = "https://timor.tech/api/holiday/year/%d/"
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

func refreshHoliday() {
	ctx := context.Background()
	year := time.Now().Year()
	uri := fmt.Sprintf(TimorApi, year)
	rawResp, err := http.Get(uri)
	if err != nil {
		return
	}
	defer rawResp.Body.Close()
	var bytes []byte
	for {
		var temp []byte
		n, err := rawResp.Body.Read(temp)
		bytes = append(bytes, temp[:n]...)
		if err == io.EOF {
			break
		}
	}
	resp := util.MustUnmarshalFromJson[TimorHolidayResponse](string(bytes))
	if resp.Code != 0 {
		return
	}
	for _, holidayInfo := range resp.Holiday {
		holiday := model.ChinaHoliday{
			Date:         holidayInfo.Date,
			IsWorkingDay: !holidayInfo.IsHoliday,
		}
		holiday.Create(ctx)
	}
}

package test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	_ "todo-reminder/conf"
	"todo-reminder/model"
)

func TestGenTodoRecordWithoutRemind(t *testing.T) {
	ctx := context.Background()
	td := model.Todo{
		NeedRemind: false,
		Content:    "test",
		UserId:     "test_user_id",
	}
	err := td.Create(ctx)
	assert.NoError(t, err)
}

func TestGenTodoRecordWithOneRemind(t *testing.T) {
	ctx := context.Background()
	td := model.Todo{
		NeedRemind: true,
		Content:    "test",
		UserId:     "test_user_id",
		RemindSetting: model.RemindSetting{
			RemindAt:     time.Now().Add(time.Hour),
			IsRepeatable: false,
		},
	}
	err := td.Create(ctx)
	assert.NoError(t, err)
}

func TestGenTodoRecordWithDailyRepeat(t *testing.T) {
	ctx := context.Background()
	td := model.Todo{
		NeedRemind: true,
		Content:    "test",
		UserId:     "test_user_id",
		RemindSetting: model.RemindSetting{
			RemindAt:     time.Now().Add(time.Hour),
			IsRepeatable: true,
			RepeatSetting: model.RepeatSetting{
				Type:       model.REPEAT_TYPE_DAILY,
				DateOffset: 1,
			},
		},
	}
	err := td.Create(ctx)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
}

func TestGenTodoRecordWithDailyRepeatEveryTwoDays(t *testing.T) {
	ctx := context.Background()
	td := model.Todo{
		NeedRemind: true,
		Content:    "test",
		UserId:     "test_user_id",
		RemindSetting: model.RemindSetting{
			RemindAt:     time.Now().Add(time.Hour),
			IsRepeatable: true,
			RepeatSetting: model.RepeatSetting{
				Type:       model.REPEAT_TYPE_DAILY,
				DateOffset: 2,
			},
		},
	}
	err := td.Create(ctx)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
}

func TestGenTodoRecordWithDailyRepeatTimePassed(t *testing.T) {
	ctx := context.Background()
	td := model.Todo{
		NeedRemind: true,
		Content:    "test",
		UserId:     "test_user_id",
		RemindSetting: model.RemindSetting{
			RemindAt:     time.Now().Add(-time.Hour),
			IsRepeatable: true,
			RepeatSetting: model.RepeatSetting{
				Type:       model.REPEAT_TYPE_DAILY,
				DateOffset: 1,
			},
		},
	}
	err := td.Create(ctx)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
}

func TestGenTodoRecordWithDailyRepeatTimePassedEveryFiveDays(t *testing.T) {
	ctx := context.Background()
	td := model.Todo{
		NeedRemind: true,
		Content:    "test",
		UserId:     "test_user_id",
		RemindSetting: model.RemindSetting{
			RemindAt:     time.Now().Add(-time.Hour),
			IsRepeatable: true,
			RepeatSetting: model.RepeatSetting{
				Type:       model.REPEAT_TYPE_DAILY,
				DateOffset: 5,
			},
		},
	}
	err := td.Create(ctx)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
}

func TestGenTodoRecordWithWeeklyRepeat(t *testing.T) {
	ctx := context.Background()
	td := model.Todo{
		NeedRemind: true,
		Content:    "test",
		UserId:     "test_user_id",
		RemindSetting: model.RemindSetting{
			RemindAt:     time.Now().Add(time.Hour),
			IsRepeatable: true,
			RepeatSetting: model.RepeatSetting{
				Type:       model.REPEAT_TYPE_WEEKLY,
				DateOffset: 1,
			},
		},
	}
	err := td.Create(ctx)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
}

func TestGenTodoRecordWithWeeklyRepeatEveryTwoWeeks(t *testing.T) {
	ctx := context.Background()
	td := model.Todo{
		NeedRemind: true,
		Content:    "test",
		UserId:     "test_user_id",
		RemindSetting: model.RemindSetting{
			RemindAt:     time.Now().Add(time.Hour),
			IsRepeatable: true,
			RepeatSetting: model.RepeatSetting{
				Type:       model.REPEAT_TYPE_WEEKLY,
				DateOffset: 2,
			},
		},
	}
	err := td.Create(ctx)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
}

func TestGenTodoRecordWithWeeklyRepeatTimePassed(t *testing.T) {
	ctx := context.Background()
	td := model.Todo{
		NeedRemind: true,
		Content:    "test",
		UserId:     "test_user_id",
		RemindSetting: model.RemindSetting{
			RemindAt:     time.Now().Add(-time.Hour),
			IsRepeatable: true,
			RepeatSetting: model.RepeatSetting{
				Type:       model.REPEAT_TYPE_WEEKLY,
				DateOffset: 1,
			},
		},
	}
	err := td.Create(ctx)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
}

func TestGenTodoRecordWithWeeklyRepeatTimePassedEveryTwoWeeks(t *testing.T) {
	ctx := context.Background()
	td := model.Todo{
		NeedRemind: true,
		Content:    "test",
		UserId:     "test_user_id",
		RemindSetting: model.RemindSetting{
			RemindAt:     time.Now().Add(-time.Hour),
			IsRepeatable: true,
			RepeatSetting: model.RepeatSetting{
				Type:       model.REPEAT_TYPE_WEEKLY,
				DateOffset: 2,
			},
		},
	}
	err := td.Create(ctx)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
}

func TestGenTodoRecordWithMonthlyRepeat(t *testing.T) {
	ctx := context.Background()
	td := model.Todo{
		NeedRemind: true,
		Content:    "test",
		UserId:     "test_user_id",
		RemindSetting: model.RemindSetting{
			RemindAt:     time.Now().Add(time.Hour),
			IsRepeatable: true,
			RepeatSetting: model.RepeatSetting{
				Type:       model.REPEAT_TYPE_MONTHLY,
				DateOffset: 1,
			},
		},
	}
	err := td.Create(ctx)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
}

func TestGenTodoRecordWithMonthlyRepeatEveryThreeMonths(t *testing.T) {
	ctx := context.Background()
	td := model.Todo{
		NeedRemind: true,
		Content:    "test",
		UserId:     "test_user_id",
		RemindSetting: model.RemindSetting{
			RemindAt:     time.Now().Add(time.Hour),
			IsRepeatable: true,
			RepeatSetting: model.RepeatSetting{
				Type:       model.REPEAT_TYPE_MONTHLY,
				DateOffset: 3,
			},
		},
	}
	err := td.Create(ctx)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
}

func TestGenTodoRecordWithMonthlyRepeatTimePassed(t *testing.T) {
	ctx := context.Background()
	td := model.Todo{
		NeedRemind: true,
		Content:    "test",
		UserId:     "test_user_id",
		RemindSetting: model.RemindSetting{
			RemindAt:     time.Now().Add(-time.Hour),
			IsRepeatable: true,
			RepeatSetting: model.RepeatSetting{
				Type:       model.REPEAT_TYPE_MONTHLY,
				DateOffset: 1,
			},
		},
	}
	err := td.Create(ctx)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
}

func TestGenTodoRecordWithMonthlyRepeatTimePassedEveryFourMonths(t *testing.T) {
	ctx := context.Background()
	td := model.Todo{
		NeedRemind: true,
		Content:    "test",
		UserId:     "test_user_id",
		RemindSetting: model.RemindSetting{
			RemindAt:     time.Now().Add(-time.Hour),
			IsRepeatable: true,
			RepeatSetting: model.RepeatSetting{
				Type:       model.REPEAT_TYPE_MONTHLY,
				DateOffset: 3,
			},
		},
	}
	err := td.Create(ctx)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
}

func TestGenTodoRecordWithYearlyRepeat(t *testing.T) {
	ctx := context.Background()
	td := model.Todo{
		NeedRemind: true,
		Content:    "test",
		UserId:     "test_user_id",
		RemindSetting: model.RemindSetting{
			RemindAt:     time.Now().Add(time.Hour),
			IsRepeatable: true,
			RepeatSetting: model.RepeatSetting{
				Type:       model.REPEAT_TYPE_YEARLY,
				DateOffset: 1,
			},
		},
	}
	err := td.Create(ctx)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
}

func TestGenTodoRecordWithYearlyRepeatEveryTenYears(t *testing.T) {
	ctx := context.Background()
	td := model.Todo{
		NeedRemind: true,
		Content:    "test",
		UserId:     "test_user_id",
		RemindSetting: model.RemindSetting{
			RemindAt:     time.Now().Add(time.Hour),
			IsRepeatable: true,
			RepeatSetting: model.RepeatSetting{
				Type:       model.REPEAT_TYPE_YEARLY,
				DateOffset: 10,
			},
		},
	}
	err := td.Create(ctx)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
}

func TestGenTodoRecordWithYearlyRepeatTimePassed(t *testing.T) {
	ctx := context.Background()
	td := model.Todo{
		NeedRemind: true,
		Content:    "test",
		UserId:     "test_user_id",
		RemindSetting: model.RemindSetting{
			RemindAt:     time.Now().Add(-time.Hour),
			IsRepeatable: true,
			RepeatSetting: model.RepeatSetting{
				Type:       model.REPEAT_TYPE_YEARLY,
				DateOffset: 1,
			},
		},
	}
	err := td.Create(ctx)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
}

func TestGenTodoRecordWithYearlyRepeatTimePassedEveryTenYears(t *testing.T) {
	ctx := context.Background()
	td := model.Todo{
		NeedRemind: true,
		Content:    "test",
		UserId:     "test_user_id",
		RemindSetting: model.RemindSetting{
			RemindAt:     time.Now().Add(-time.Hour),
			IsRepeatable: true,
			RepeatSetting: model.RepeatSetting{
				Type:       model.REPEAT_TYPE_YEARLY,
				DateOffset: 10,
			},
		},
	}
	err := td.Create(ctx)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
}

func TestGenTodoRecordWithWorkingDayRepeat(t *testing.T) {
	ctx := context.Background()
	td := model.Todo{
		NeedRemind: true,
		Content:    "test",
		UserId:     "test_user_id",
		RemindSetting: model.RemindSetting{
			RemindAt:     time.Now().Add(time.Hour),
			IsRepeatable: true,
			RepeatSetting: model.RepeatSetting{
				Type: model.REPEAT_TYPE_WORKING_DAY,
			},
		},
	}
	err := td.Create(ctx)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
}

func TestGenTodoRecordWithWorkingDayRepeatTimePassed(t *testing.T) {
	ctx := context.Background()
	td := model.Todo{
		NeedRemind: true,
		Content:    "test",
		UserId:     "test_user_id",
		RemindSetting: model.RemindSetting{
			RemindAt:     time.Now().Add(-time.Hour),
			IsRepeatable: true,
			RepeatSetting: model.RepeatSetting{
				Type: model.REPEAT_TYPE_WORKING_DAY,
			},
		},
	}
	err := td.Create(ctx)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
}

func TestGenTodoRecordWithHolidayRepeat(t *testing.T) {
	ctx := context.Background()
	td := model.Todo{
		NeedRemind: true,
		Content:    "test",
		UserId:     "test_user_id",
		RemindSetting: model.RemindSetting{
			RemindAt:     time.Now().Add(time.Hour),
			IsRepeatable: true,
			RepeatSetting: model.RepeatSetting{
				Type: model.REPEAT_TYPE_HOLIDAY,
			},
		},
	}
	err := td.Create(ctx)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
}

func TestGenTodoRecordWithHolidayRepeatTimePassed(t *testing.T) {
	ctx := context.Background()
	td := model.Todo{
		NeedRemind: true,
		Content:    "test",
		UserId:     "test_user_id",
		RemindSetting: model.RemindSetting{
			RemindAt:     time.Now().Add(-time.Hour),
			IsRepeatable: true,
			RepeatSetting: model.RepeatSetting{
				Type: model.REPEAT_TYPE_HOLIDAY,
			},
		},
	}
	err := td.Create(ctx)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
	err = td.GenNextRecord(ctx, td.Id, false)
	assert.NoError(t, err)
}

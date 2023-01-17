# MongoDB schema

## user

```js
{
    _id: ObjectId,
    userId: String, // QQ 号
    password: String, // hashed password
    createdAt: DateTime,
    updatedAt: DateTime,
    isDeleted: Boolean,
}
```

## todo

```js
{
    _id: ObjectId,
    isDeleted: Boolean,
    CreatedAt: DateTime,
    UpdatedAt: DateTime,
    needRemind: Boolean,
    content: String,
    userId: String,
    remindSetting: {
        remindAt: DateTime, // 提醒时间
        isRepeatable: Boolean, // 是否重复
        repeatSetting: {
            type: String, // 重复类型，daily（每天）、weekly（每周）、monthly（每月）、yearly（每年）、workingDay（工作日）、holiday（节假日）
            month: Long, // 几月
            weekday: Long, // 周几
            day: Long, // 几号
            dateOffset: Long, // 每多少天、周、月、年
        }
    }
}
```

## todoRecord

```js
{
    _id: ObjectId,
    isDeleted: Boolean,
    createdAt: DateTime,
    updatedAt: DateTime,
    remindAt: dateTime,
    hasBeenDone: Boolean,
    content: String,
    todoId: ObjectId,
    doneAt: DateTime,
    needRemind: Boolean,
}
```

## chinaHoliday

```js
{
    _id: ObjectId,
    date: String, // YYYYMMDD 格式的日期字符串
    isWorkingDay: Boolean, // 是否是工作日
}
```

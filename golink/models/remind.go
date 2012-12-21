package models

import (
    "errors"
    "fmt"
    "github.com/QLeelulu/goku"
)

type RemindType int

const (
    REMIND_COMMENT RemindType = 1 // 新增评论提醒
    REMIND_FANS    RemindType = 2 // 新增粉丝提醒
)

var remindTypeKey map[RemindType]string = map[RemindType]string{
    REMIND_COMMENT: "c",
    REMIND_FANS:    "f",
}

// 评论、粉丝等提醒信息
type RemindInfo struct {
    Comments int
    Fans     int
}

// 增加用户的提醒数
func Remind_Inc(userId int64, t RemindType) error {
    field, ok := remindTypeKey[t]
    if !ok {
        return errors.New("错误提醒类型")
    }
    redisClient := GetRedis()
    defer redisClient.Quit()

    key := fmt.Sprintf("rd:%d", userId)
    _, err := redisClient.Hincrby(key, field, 1)
    if err != nil {
        goku.Logger().Errorln(err.Error())
    }
    return err
}

func Remind_ForUser(userId int64) (r RemindInfo, err error) {
    redisClient := GetRedis()
    defer redisClient.Quit()

    key := fmt.Sprintf("rd:%d", userId)
    res, err_ := redisClient.Hmget(key,
        remindTypeKey[REMIND_COMMENT],
        remindTypeKey[REMIND_FANS],
    )
    if err_ != nil {
        if err_.Error() != "Nonexisting key" {
            err = err_
            goku.Logger().Errorln(err.Error())
        }
        return
    }
    r2 := res.IntArray()
    r.Comments = int(r2[0])
    r.Fans = int(r2[1])
    return
}

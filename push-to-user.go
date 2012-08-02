package main

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink"
    "github.com/QLeelulu/ohlala/golink/models"
    "strconv"
    "strings"
    "time"
)

func main() {
    redisClient := models.GetRedis()
    defer redisClient.Quit()
    // 从推送队列取出处理
    for {
        // 格式: pushtype,userid(topicid),linkid,timestamp
        item, err := redisClient.Rpop(golink.KEY_LIST_PUSH_TO_USER)
        if err != nil {
            if strings.Index(err.Error(), "Nonexisting key") < 0 {
                goku.Logger().Errorln(err.Error())
            }
            time.Sleep(30 * time.Second)
        } else if item != nil {
            items := strings.Split(item.String(), ",")
            if len(items) > 1 {
                kid, err := strconv.ParseInt(items[1], 10, 64)
                if err != nil {
                    goku.Logger().Errorln(err.Error())
                    continue
                }
                linkId, err := strconv.ParseInt(items[2], 10, 64)
                if err != nil {
                    goku.Logger().Errorln(err.Error())
                    continue
                }
                if items[0] == "1" {
                    // 推给粉丝
                    err = models.LinkForUser_ToUserFollowers(kid, linkId)
                } else if items[0] == "2" {
                    // 推给话题关注者
                    err = models.LinkForUser_ToTopicidFollowers(kid, linkId)
                } else {
                    goku.Logger().Errorln("值错误:", item.String())
                }
                if err != nil {
                    goku.Logger().Errorln(err.Error())
                }
            }
        } else {
            time.Sleep(30 * time.Second)
        }
    }
}

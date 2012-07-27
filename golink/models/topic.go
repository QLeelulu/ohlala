package models

import (
    "github.com/QLeelulu/goku"
    "strings"
)

type Topic struct {
    Id        int64
    Name      string
    NameLower string // topic 名称小写，唯一索引
    Desc      string // 话题的描述
    Pic       string // 话题的图片
    Clicks    int64  // 话题点击次数
    Followers int    // 话题的关注者数量
    Links     int    // 添加到该话题的链接数量
}

type TopicToLink struct {
    TopicId int64
    LinkId  int64
}

// 保持topic到数据库，同时建立topic与link的关系表
// 全部成功则返回true
func Topic_SaveTopics(topics string, linkId int64) bool {
    if topics == "" {
        return true
    }
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    success := true
    topicList := strings.Split(topics, ",")
    for _, topic := range topicList {
        topicLower := strings.ToLower(topic)
        t := new(Topic)
        err := db.GetStruct(t, "`name_lower`=?", topic)
        if err != nil {
            goku.Logger().Logln(topic)
            goku.Logger().Errorln(err.Error())
            success = false
            continue
        }
        if t.Id < 1 {
            t.Name = topic
            t.NameLower = topicLower
            _, err = db.InsertStruct(t)
            if err != nil {
                goku.Logger().Errorln(err.Error())
                success = false
                continue
            }
        }
        if t.Id > 0 && linkId > 0 {
            _, err = db.Insert("topic_link", map[string]interface{}{"topic_id": t.Id, "link_id": linkId})
            if err != nil {
                goku.Logger().Errorln(err.Error())
                success = false
            }
        }
    }
    return success
}

func Topic_GetByName(name string) (*Topic, error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    t := new(Topic)
    err := db.GetStruct(t, "`name`=?", strings.ToLower(name))
    if err != nil || t.Id == 0 {
        if err != nil {
            goku.Logger().Errorln(err.Error())
        }
        t = nil
    }
    return t, err
}

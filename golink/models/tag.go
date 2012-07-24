package models

import (
    "github.com/QLeelulu/goku"
    "strings"
)

type Tag struct {
    Id         int64
    Name       string
    ClickCount int64
}

type TagToLink struct {
    TagId  int64
    LinkId int64
}

// 保持tag到数据库，同时建立tag与link的关系表
// 全部成功则返回true
func Tag_SaveTags(tags string, linkId int64) bool {
    if tags == "" {
        return true
    }
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    success := true
    tagList := strings.Split(tags, ",")
    for _, tag := range tagList {
        tag = strings.ToLower(tag)
        t := new(Tag)
        err := db.GetStruct(t, "`name`=?", tag)
        if err != nil {
            goku.Logger().Logln(tag)
            goku.Logger().Errorln(err.Error())
            success = false
            continue
        }
        if t.Id < 1 {
            t.Name = tag
            _, err = db.InsertStruct(t)
            if err != nil {
                goku.Logger().Errorln(err.Error())
                success = false
                continue
            }
        }
        if t.Id > 0 && linkId > 0 {
            _, err = db.Insert("tag_to_link", map[string]interface{}{"tag_id": t.Id, "link_id": linkId})
            if err != nil {
                goku.Logger().Errorln(err.Error())
                success = false
            }
        }
    }
    return success
}

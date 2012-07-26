package models

import (
    "github.com/QLeelulu/goku"
    "time"
)

/**
 * 链接推送给用户
 */

func LinkForUser_Add(userId, linkId int64) error {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    m := map[string]interface{}{
        "user_id":     userId,
        "link_id":     linkId,
        "create_time": time.Now(),
    }

    _, err := db.Insert("user_link", m)
    return err
}

// 将linkid推送给userid的所有粉丝
func LinkForUser_ToUserFollowers(userId, linkId int64) error {
    db := GetDB()
    defer db.Close()

    qi := goku.SqlQueryInfo{}
    qi.Fields = "`user_id`"
    qi.Where = "`follow_id`=?"
    qi.Params = []interface{}{userId}
    rows, err := db.Select("user_follow", qi)
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return err
    }

    var uid int64
    for rows.Next() {
        err = rows.Scan(&uid)
        if err == nil && uid > 0 {
            LinkForUser_Add(uid, linkId)
        }
    }
    return nil
}

// 将linkid推送给tag的所有关注者
func LinkForUser_ToTagFollowers(tag string, linkId int64) error {
    db := GetDB()
    defer db.Close()

    t, err := Tag_GetByName(tag)
    if err != nil {
        return err
    }
    if t == nil {
        return nil
    }

    qi := goku.SqlQueryInfo{}
    qi.Fields = "`user_id`"
    qi.Where = "`tag_id`=?"
    qi.Params = []interface{}{t.Id}
    rows, err := db.Select("user_to_tag", qi)
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return err
    }

    var uid int64
    for rows.Next() {
        err = rows.Scan(&uid)
        if err == nil && uid > 0 {
            LinkForUser_Add(uid, linkId)
        }
    }
    return nil
}

// 用户(userId)关注用户(followId)时，
// 将用户(followId)的链接推送给用户(userId)
func LinkForUser_FollowUser(userId, followId int64) {
    if userId < 1 {
        return
    }
    db := GetDB()
    defer db.Close()
    limit := 800 // 只导入最新的N条
    _, err := db.Exec("insert ingore into link_for_user(user_id,link_id,create_time) (select ?,id, NOW() from link where `user_id`=? order by `id` desc limit ?)", userId, followId, limit)
    if err != nil {
        goku.Logger().Errorln(err.Error())
    }
}

// 用户(userId)关注Tag时，
// 将Tag的链接推送给用户(userId)
func LinkForUser_FollowTag(userId, tagId int64) {
    if userId < 1 {
        return
    }
    db := GetDB()
    defer db.Close()
    limit := 800 // 只导入最新的N条
    _, err := db.Exec("insert ingore into link_for_user(user_id,link_id,create_time) (select ?,link_id, NOW() from tag_link where `tag_id`=? order by link_id desc limit ?)", userId, tagId, limit)
    if err != nil {
        goku.Logger().Errorln(err.Error())
    }
}

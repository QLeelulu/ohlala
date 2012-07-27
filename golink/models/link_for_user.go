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

// 将linkid推送给topic的所有关注者
func LinkForUser_ToTopicFollowers(topic string, linkId int64) error {
    db := GetDB()
    defer db.Close()

    t, err := Topic_GetByName(topic)
    if err != nil {
        return err
    }
    if t == nil {
        return nil
    }

    qi := goku.SqlQueryInfo{}
    qi.Fields = "`user_id`"
    qi.Where = "`topic_id`=?"
    qi.Params = []interface{}{t.Id}
    rows, err := db.Select("user_to_topic", qi)
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
    _, err := db.Exec("insert ignore into link_for_user(user_id,link_id,create_time) (select ?,id, NOW() from link where `user_id`=? order by `id` desc limit ?)", userId, followId, limit)
    if err != nil {
        goku.Logger().Errorln(err.Error())
    }
}

// 用户(userId)关注Topic时，
// 将Topic的链接推送给用户(userId)
func LinkForUser_FollowTopic(userId, topicId int64) {
    if userId < 1 {
        return
    }
    db := GetDB()
    defer db.Close()
    limit := 800 // 只导入最新的N条
    _, err := db.Exec("insert ignore into link_for_user(user_id,link_id,create_time) (select ?,T.link_id, NOW() from topic_link as T where T.`topic_id`=? order by T.link_id desc limit ?)", userId, topicId, limit)
    if err != nil {
        goku.Logger().Errorln(err.Error())
    }
}

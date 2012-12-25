package models

import (
    "fmt"
    "github.com/QLeelulu/goku"
    "strings"
    "time"
)

/**
 * 处理用户关注、链接推送
 */

const (
    LinkForUser_ByUser  = 1 // 由于关注用户而引发的推送
    LinkForUser_ByTopic = 2 // 由于关注话题而引发的推送
)

// 按用户id分表
func LinkForUser_TableName(userId int64) string {
    return fmt.Sprintf("link_for_user_%v", userId%24)
}

// 推送链接给用户
// @t: 推送类型， 1:关注的用户, 2:关注的话题
func LinkForUser_Add(userId, linkId int64, t int) error {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    return linkForUser_AddWithDb(db, userId, linkId, t)
}

// 减少DB操作
// @t: 推送类型， 1:关注的用户, 2:关注的话题
func linkForUser_AddWithDb(db *goku.MysqlDB, userId, linkId int64, t int) error {
    m := map[string]interface{}{
        "user_id":     userId,
        "link_id":     linkId,
        "create_time": time.Now(),
    }
    if t == 1 {
        m["user_count"] = 1
    } else {
        m["topic_count"] = 1
    }

    _, err := db.Insert(LinkForUser_TableName(userId), m)
    if err != nil {
        if strings.Index(err.Error(), "Duplicate entry") > -1 {
            m := map[string]interface{}{}
            if t == 1 {
                m["user_count"] = 1
            } else {
                m["topic_count"] = 1
            }
            _, err = db.Update(LinkForUser_TableName(userId), m, "user_id=? and link_id=?", userId, linkId)
            if err != nil {
                goku.Logger().Errorln(err.Error())
            }
        } else {
            goku.Logger().Errorln(err.Error())
        }
    }
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
            linkForUser_AddWithDb(db, uid, linkId, LinkForUser_ByUser)
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
    if t == nil || t.Id < 1 {
        return nil
    }

    return LinkForUser_ToTopicidFollowers(t.Id, linkId)
}

func LinkForUser_ToTopicidFollowers(topicId int64, linkId int64) error {
    db := GetDB()
    defer db.Close()

    qi := goku.SqlQueryInfo{}
    qi.Fields = "`user_id`"
    qi.Where = "`topic_id`=?"
    qi.Params = []interface{}{topicId}
    rows, err := db.Select("topic_follow", qi)
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return err
    }

    var uid int64
    for rows.Next() {
        err = rows.Scan(&uid)
        if err == nil && uid > 0 {
            linkForUser_AddWithDb(db, uid, linkId, LinkForUser_ByTopic)
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
    // 先更新
    _, err := db.Exec("update "+LinkForUser_TableName(userId)+
        " set user_count=user_count+1 where link_id in "+
        "    (select l.id from (select `id` from `link` where `user_id`=? limit 10000) as l)",
        followId)
    if err != nil {
        goku.Logger().Errorln(err.Error())
    }
    // 插入
    _, err = db.Exec("insert ignore into "+LinkForUser_TableName(userId)+" (user_id,link_id,user_count,create_time) (select ?,id, 1, NOW() from link where `user_id`=? order by `id` desc limit ?)", userId, followId, limit)
    if err != nil {
        goku.Logger().Errorln(err.Error())
    }
}

// 用户(userId)取消关注用户(followId)时，
// 将用户(followId)的链接从用户(userId)的推送列表中移除
func LinkForUser_UnFollowUser(userId, followId int64) {
    if userId < 1 {
        return
    }
    db := GetDB()
    defer db.Close()

    // 先更新计数
    _, err := db.Exec("update "+LinkForUser_TableName(userId)+
        " set user_count=user_count-1 where link_id in "+
        "    (select l.id from (select `id` from `link` where `user_id`=? limit 10000) as l)",
        followId)
    if err != nil {
        goku.Logger().Errorln(err.Error())
    }
    // 删除
    _, err = db.Exec("delete from "+LinkForUser_TableName(userId)+
        " where user_id=? and user_count=0 and topic_count=0", userId)
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
    // 先更新
    _, err := db.Exec("update "+LinkForUser_TableName(userId)+
        " set topic_count=topic_count+1 where link_id in "+
        "    ( select tl.link_id from (select link_id from `topic_link` where `topic_id`=? limit 10000 ) as tl )",
        topicId)
    if err != nil {
        goku.Logger().Errorln(err.Error())
    }
    // 插入
    _, err = db.Exec("insert ignore into "+LinkForUser_TableName(userId)+" (user_id,link_id,topic_count,create_time) (select ?, T.link_id, 1, NOW() from topic_link as T where T.`topic_id`=? order by T.link_id desc limit ?)", userId, topicId, limit)
    if err != nil {
        goku.Logger().Errorln(err.Error())
    }
}

// 用户(userId)取消关注Topic时，
// 将Topic的链接从用户(userId)的推送列表里面移除
func LinkForUser_UnFollowTopic(userId, topicId int64) {
    if userId < 1 {
        return
    }
    db := GetDB()
    defer db.Close()

    // 先更新
    _, err := db.Exec("update "+LinkForUser_TableName(userId)+
        " set topic_count=topic_count-1 where link_id in "+
        "    ( select tl.link_id from (select link_id from `topic_link` where `topic_id`=? limit 10000 ) as tl )",
        topicId)
    if err != nil {
        goku.Logger().Errorln(err.Error())
    }
    // 删除
    _, err = db.Exec("delete from "+LinkForUser_TableName(userId)+
        " where user_id=? and user_count=0 and topic_count=0", userId)
    if err != nil {
        goku.Logger().Errorln(err.Error())
    }
}

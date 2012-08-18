package models

import (
    "database/sql"
    "errors"
    "fmt"
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink"
    "strings"
    "time"
)

type Topic struct {
    Id            int64
    Name          string
    NameLower     string // topic 名称小写，唯一索引
    Description   string // 话题的描述
    Pic           string // 话题的图片
    ClickCount    int64  // 话题点击次数
    FollowerCount int64  // 话题的关注者数量
    LinkCount     int64  // 添加到该话题的链接数量
}

func (t Topic) PicPath() string {
    if t.Pic == "" {
        return "/assets/img/avatar/topic/topic_default.png"
    }
    return "/assets/img/avatar/topic/" + t.Pic
}

type TopicToLink struct {
    TopicId int64
    LinkId  int64
}

type VTopic struct {
    *Topic
    IsFollowed bool // 是否已关注
}

// 转换为用于view的用户类型
func Topic_ToVTopic(t *Topic, ctx *goku.HttpContext) *VTopic {
    if t == nil {
        return nil
    }
    vt := &VTopic{Topic: t}
    var userId int64
    if user, ok := ctx.Data["user"].(*User); ok && user != nil {
        userId = user.Id
    }
    if userId > 0 {
        vt.IsFollowed = Topic_CheckFollow(userId, vt.Id)
    }

    return vt
}

// 检查用户是否已经关注话题，
// @isFollowed: 是否已经关注话题
func Topic_CheckFollow(userId, topicId int64) (isFollowed bool) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    rows, err := db.Query("select * from `topic_follow` where `user_id`=? and `topic_id`=? limit 1",
        userId, topicId)
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return
    }
    defer rows.Close()
    if rows.Next() {
        isFollowed = true
    }

    return
}

// 保持topic到数据库，同时建立topic与link的关系表
// 如果topic已经存在，则直接建立与link的关联
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
            } else {
                // 成功，更新话题的链接数量统计
                Topic_IncCount(db, t.Id, "link_count", 1)

                redisClient := GetRedis()
                defer redisClient.Quit()
                // 加入推送队列
                // 格式: pushtype,topicid,linkid,timestamp
                qv := fmt.Sprintf("%v,%v,%v,%v", LinkForUser_ByTopic, t.Id, linkId, time.Now().Unix())
                _, err = redisClient.Lpush(golink.KEY_LIST_PUSH_TO_USER, qv)
                if err != nil {
                    goku.Logger().Errorln(err.Error())
                }
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

// 用户userId 关注 话题topicId
func Topic_Follow(userId, topicId int64) (bool, error) {
    if userId < 1 || topicId < 1 {
        return false, errors.New("参数错误")
    }
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    vals := map[string]interface{}{
        "user_id":     userId,
        "topic_id":    topicId,
        "create_time": time.Now(),
    }
    r, err := db.Insert("topic_follow", vals)
    if err != nil {
        if strings.Index(err.Error(), "Duplicate entry") > -1 {
            return false, errors.New("已经关注该话题")
        } else {
            goku.Logger().Errorln(err.Error())
            return false, err
        }
    }

    var afrow int64
    afrow, err = r.RowsAffected()
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return false, err
    }

    if afrow > 0 {
        // 关注话题成功，将话题的链接推送给用户
        LinkForUser_FollowTopic(userId, topicId)
        // 更新用户关注话题的数量
        User_IncCount(db, userId, "ftopic_count", 1)
        // 更新话题的关注用户数
        Topic_IncCount(db, topicId, "follower_count", 1)
        return true, nil
    }
    return false, nil
}

// 获取关注topicId的用户列表
func Topic_GetFollowers(topicId int64, page, pagesize int) ([]User, error) {
    if page < 1 {
        page = 1
    }
    page = page - 1
    if pagesize == 0 {
        pagesize = 20
    }
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    qi := goku.SqlQueryInfo{}
    qi.Fields = "u.id, u.name, u.email, u.user_pic"
    qi.Join = " tf INNER JOIN `user` u ON tf.user_id=u.id"
    qi.Where = "tf.topic_id=?"
    qi.Params = []interface{}{topicId}
    qi.Limit = pagesize
    qi.Offset = pagesize * page
    qi.Order = "u.id desc"

    rows, err := db.Select("topic_follow", qi)

    if err != nil {
        goku.Logger().Errorln(err.Error())
        return nil, err
    }
    defer rows.Close()

    users := make([]User, 0)
    for rows.Next() {
        user := User{}
        err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.UserPic)
        if err != nil {
            goku.Logger().Errorln(err.Error())
            return nil, err
        }
        users = append(users, user)
    }
    return users, nil
}

// 加（减）话题信息里面的统计数据
// @field: 要修改的字段
// @inc: 要增加或减少的值
func Topic_IncCount(db *goku.MysqlDB, topicId int64, field string, inc int) (sql.Result, error) {
    // m := map[string]interface{}{field: fmt.Sprintf("%v+%v", field, inc)}
    // r, err := db.Update("user", m, "id=?", userid)
    r, err := db.Exec(fmt.Sprintf("UPDATE `topic` SET %s=%s+? WHERE id=?;", field, field), inc, topicId)
    if err != nil {
        goku.Logger().Errorln(err.Error())
    }
    return r, err
}

func Topic_UpdatePic(id int64, pic string) (sql.Result, error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    m := map[string]interface{}{"pic": pic}
    r, err := db.Update("topic", m, "id=?", id)
    if err != nil {
        goku.Logger().Errorln(err.Error())
    }
    return r, err
}

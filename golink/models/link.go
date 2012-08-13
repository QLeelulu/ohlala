package models

import (
    "fmt"
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/goku/form"
    "github.com/QLeelulu/ohlala/golink"
    "github.com/QLeelulu/ohlala/golink/utils"
    "strings"
    "time"
)

type Link struct {
    Id           int64
    UserId       int64
    Title        string
    Context      string // 如为链接，则为url地址
    ContextType  int    // 1: url
    Topics       string
    VoteUp       int64
    VoteDown     int64
    RedditScore  float64
    ViewCount    int
    CommentCount int
    CreateTime   time.Time

    user *User `db:"exclude"`
}

func (l *Link) User() *User {
    if l.user == nil {
        l.user = User_GetById(l.UserId)
    }
    return l.user
}

func (l *Link) TopicList() []string {
    if l.Topics == "" {
        return nil
    }
    return strings.Split(l.Topics, ",")
}

func (l *Link) SinceTime() string {
    return utils.SmcTimeSince(l.CreateTime)
}

// 保存link到数据库，如果成功，则返回link的id
func Link_SaveMap(m map[string]interface{}) int64 {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()
    m["create_time"] = time.Now()
    //新增link默认投票1次,显示的时候默认减一
    m["vote_up"] = 1
    m["reddit_score"] = utils.RedditSortAlgorithm(m["create_time"].(time.Time), int64(1), int64(0))

    r, err := db.Insert("link", m)
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return 0
    }
    var id int64
    id, err = r.LastInsertId()
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return 0
    }

    if id > 0 {
        uid := m["user_id"].(int64)
        // 直接推送给自己，自己必须看到
        LinkForUser_Add(uid, id, LinkForUser_ByUser)

        // 存入`tui_link_for_handle` 链接处理队列表
        db.Query("INSERT ignore INTO tui_link_for_handle(link_id,create_time,user_id,insert_time,data_type) VALUES (?, ?, ?, NOW(), ?)",
            id, m["create_time"].(time.Time), uid, 1)

        redisClient := GetRedis()
        defer redisClient.Quit()
        // 加入推送队列
        // 格式: pushtype,userid,linkid,timestamp
        qv := fmt.Sprintf("%v,%v,%v,%v", LinkForUser_ByUser, uid, id, time.Now().Unix())
        _, err = redisClient.Lpush(golink.KEY_LIST_PUSH_TO_USER, qv)
        if err != nil {
            goku.Logger().Errorln(err.Error())
            return 0
        }

    }

    return id
}

// 如果保持失败，则返回错误信息
func Link_SaveForm(f *form.Form, userId int64) (bool, []string) {
    errorMsgs := make([]string, 0)
    if f.Valid() {
        m := f.CleanValues()
        m["topics"] = buildTopics(m["topics"].(string))
        m["user_id"] = userId

        id := Link_SaveMap(m)
        if id > 0 {
            Topic_SaveTopics(m["topics"].(string), id)
        } else {
            errorMsgs = append(errorMsgs, golink.ERROR_DATABASE)
        }
    } else {
        errs := f.Errors()
        for _, v := range errs {
            errorMsgs = append(errorMsgs, v[0]+": "+v[1])
        }
    }
    if len(errorMsgs) < 1 {
        return true, nil
    }
    return false, errorMsgs
}

// topic可以用英文逗号或者空格分隔
// 过滤重复topic，最终返回的topic列表只用英文逗号分隔
func buildTopics(topics string) string {
    if topics == "" {
        return ""
    }
    m := make(map[string]string)
    t := strings.Split(topics, ",")
    for _, topic := range t {
        topic = strings.TrimSpace(topic)
        if topic != "" {
            t2 := strings.Split(topic, " ")
            for _, topic2 := range t2 {
                topic2 = strings.TrimSpace(topic2)
                if topic2 != "" {
                    m[strings.ToLower(topic2)] = topic2
                }
            }
        }
    }
    r := ""
    for _, v := range m {
        if r != "" {
            r += ","
        }
        r += v
    }
    return r
}

func Link_GetById(id int64) (*Link, error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    l := new(Link)
    err := db.GetStruct(l, "id=?", id)
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return nil, err
    }
    if l.Id > 0 {
        return l, nil
    }
    return nil, nil
}

// @page: 从1开始
func Link_GetByPage(page, pagesize int) []Link {
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
    qi.Limit = pagesize
    qi.Offset = page * pagesize
    qi.Order = "id desc"
    var links []Link
    err := db.GetStructs(&links, qi)
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return nil
    }
    return links
}

// 获取由用户发布的link
// @page: 从1开始
func Link_ByUser(userId int64, page, pagesize int) []Link {
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
    qi.Limit = pagesize
    qi.Offset = page * pagesize
    qi.Where = "`user_id`=?"
    qi.Params = []interface{}{userId}
    var links []Link
    err := db.GetStructs(&links, qi)
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return nil
    }
    return links
}

// 获取属于某话题的link
// @page: 从1开始
func Link_ForTopic(topicId int64, page, pagesize int) ([]Link, error) {
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
    qi.Fields = "l.id, l.title, l.context, l.topics"
    qi.Join = " tl INNER JOIN `link` l ON tl.link_id=l.id"
    qi.Where = "tl.topic_id=?"
    qi.Params = []interface{}{topicId}
    qi.Limit = pagesize
    qi.Offset = pagesize * page
    qi.Order = "l.id desc"

    rows, err := db.Select("topic_link", qi)

    if err != nil {
        goku.Logger().Errorln(err.Error())
        return nil, err
    }
    defer rows.Close()

    links := make([]Link, 0)
    for rows.Next() {
        link := Link{}
        err = rows.Scan(&link.Id, &link.Title, &link.Context, &link.Topics)
        if err != nil {
            goku.Logger().Errorln(err.Error())
            return nil, err
        }
        links = append(links, link)
    }
    return links, nil
}

// 获取属于某用户的link
// @page: 从1开始
// @orderType: 排序类型, hot:热门, hotc:热议, time:最新, vote:投票得分
func Link_ForUser(userId int64, orderType string, page, pagesize int) ([]Link, error) {
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
    qi.Fields = "l.id, l.user_id, l.title, l.context, l.topics, l.vote_up, l.vote_down, l.view_count, l.comment_count, l.create_time"
    qi.Join = " ul INNER JOIN `link` l ON ul.link_id=l.id"
    qi.Where = "ul.user_id=?"
    qi.Params = []interface{}{userId}
    qi.Limit = pagesize
    qi.Offset = pagesize * page
    switch orderType {
    case "time":
        qi.Order = "l.id desc"
    case "hotc":
        qi.Order = "ABS(L.vote_up-L.vote_down) asc,L.vote_up+L.vote_down desc, id desc"
    case "vote":
        qi.Order = "l.vote_up desc, id desc"
    default:
        qi.Order = "l.reddit_score, id desc"
    }

    rows, err := db.Select(LinkForUser_TableName(userId), qi)

    if err != nil {
        goku.Logger().Errorln(err.Error())
        return nil, err
    }
    defer rows.Close()

    links := make([]Link, 0)
    for rows.Next() {
        link := Link{}
        err = rows.Scan(&link.Id, &link.UserId, &link.Title, &link.Context, &link.Topics,
            &link.VoteUp, &link.VoteDown, &link.ViewCount, &link.CommentCount, &link.CreateTime)
        if err != nil {
            goku.Logger().Errorln(err.Error())
            return nil, err
        }
        links = append(links, link)
    }
    return links, nil
}

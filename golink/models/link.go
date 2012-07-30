package models

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/goku/form"
    "github.com/QLeelulu/ohlala/golink"
    "strings"
    "time"
)

type Link struct {
    Id          int64
    UserId      int64
    Title       string
    Context     string // 如为链接，则为url地址
    ContextType int    // 1: url
    Topics      string
    VoteUp      int64
    VoteDown    int64
    RedditScore float64
    CreateTime  time.Time

    user *User `db:"exclude"`
}

func (l *Link) User() *User {
    if l.user == nil {
        l.user = User_GetById(l.UserId)
    }
    return l.user
}

// 保存link到数据库，如果成功，则返回link的id
func Link_SaveMap(m map[string]interface{}) int64 {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()
    m["create_time"] = time.Now()
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
    var links []Link
    err := db.GetStructs(&links, qi)
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return nil
    }
    return links
}

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

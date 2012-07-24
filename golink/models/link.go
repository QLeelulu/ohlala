package models

import (
    "github.com/QLeelulu/goku"
    "time"
)

type Link struct {
    Id          int64
    UserId      int64
    Title       string
    Context     string // 如为链接，则为url地址
    ContextType int    // 1: url
    Tags        string
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

// @page: 从1开始
func Link_GetByPage(page, pagesize int) []*Link {
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
    links, err := db.GetStructs(Link{}, qi)
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return nil
    }
    r := make([]*Link, len(links))
    for i, l := range links {
        r[i] = l.(*Link)
    }
    return r
}

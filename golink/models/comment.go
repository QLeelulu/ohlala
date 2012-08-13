package models

import (
    "errors"
    "fmt"
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/goku/form"
    "github.com/QLeelulu/ohlala/golink"
    "github.com/QLeelulu/ohlala/golink/utils"
    "time"
)

const Table_Comment = "comment"

type Comment struct {
    Id            int64
    LinkId        int64
    UserId        int64
    Status        int // 评论状态：1代表正常、2代表删除
    Content       string
    ParentId      int64
    TopParentId   int64
    ParentPath    string
    ChildrenCount int
    VoteUp        int
    VoteDown      int
    RedditScore   float64
    CreateTime    time.Time

    user *User `db:"exclude"`
}

func (c *Comment) User() *User {
    if c.user == nil {
        c.user = User_GetById(c.UserId)
    }
    return c.user
}

func (c *Comment) SinceTime() string {
    return utils.SmcTimeSince(c.CreateTime)
}

// 保存评论到数据库，如果成功，则返回comment的id
func Comment_SaveMap(m map[string]interface{}) (int64, error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    // TODO: 链接评论的链接存不存在？

    // 检查父评论是否存在
    var pComment *Comment
    var err error
    if id, ok := m["parent_id"].(int64); ok && id > 0 {
        pComment, err = Comment_GetById(id)
        if err != nil {
            goku.Logger().Errorln(err.Error())
            return int64(0), err
        }
        // 指定了父评论的id但是数据库中没有
        if pComment == nil {
            return int64(0), errors.New("指定的父评论不存在")
        }
    }

    // 路径相关
    if pComment == nil {
        m["parent_id"] = 0
        m["top_parent_id"] = 0
        m["parent_path"] = "/"
    } else {
        m["parent_id"] = pComment.Id
        if pComment.TopParentId == 0 {
            m["top_parent_id"] = pComment.Id
        } else {
            m["top_parent_id"] = pComment.TopParentId
        }
        m["parent_path"] = fmt.Sprintf("%s%s/", pComment.ParentPath, pComment.Id)
    }

    m["status"] = 1
    m["create_time"] = time.Now()
    //新增comment默认投票1次,显示的时候默认减一
    m["vote_up"] = 1
    m["reddit_score"] = utils.RedditSortAlgorithm(m["create_time"].(time.Time), int64(1), int64(0))

    r, err := db.Insert(Table_Comment, m)
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return 0, err
    }
    var id int64
    id, err = r.LastInsertId()
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return 0, err
    }

    if id > 0 {
        if pComment != nil {
            IncCountById(db, Table_Comment, pComment.Id, "children_count", 1)
        }
    }

    return id, nil
}

// 如果保存失败，则返回错误信息
func Comment_SaveForm(f *form.Form, userId int64) (bool, []string) {
    errorMsgs := make([]string, 0)
    if f.Valid() {
        m := f.CleanValues()
        m["user_id"] = userId

        id, err := Comment_SaveMap(m)
        if err != nil || id < 1 {
            errorMsgs = append(errorMsgs, golink.ERROR_DATABASE)
        }
    } else {
        errs := f.Errors()
        for _, v := range errs {
            errorMsgs = append(errorMsgs, v[1])
        }
    }
    if len(errorMsgs) < 1 {
        return true, nil
    }
    return false, errorMsgs
}

func Comment_GetById(id int64) (*Comment, error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    c := new(Comment)
    err := db.GetStruct(c, "id=?", id)
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return nil, err
    }
    if c.Id > 0 {
        return c, nil
    }
    return nil, nil
}

// @page: 从1开始
func Comment_GetByPage(page, pagesize int) []Comment {
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
    var comments []Comment
    err := db.GetStructs(&comments, qi)
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return nil
    }
    return comments
}

// 获取由用户发布的评论
// @page: 从1开始
func Comment_ByUser(userId int64, page, pagesize int) []Comment {
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
    var comments []Comment
    err := db.GetStructs(&comments, qi)
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return nil
    }
    return comments
}

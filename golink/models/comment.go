package models

import (
    // "bytes"
    "errors"
    "fmt"
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/goku/form"
    "github.com/QLeelulu/ohlala/golink"
    "github.com/QLeelulu/ohlala/golink/utils"
    // "html/template"
    // "strings"
    "time"
)

const (
    COMMENT_STATUS_NORMAL = 1
    COMMENT_STATUS_DEL    = 2
)

var commentStatus map[int]string = map[int]string{
    COMMENT_STATUS_NORMAL: "正常",
    COMMENT_STATUS_DEL:    "删除",
}

type Comment struct {
    Id            int64
    LinkId        int64
    UserId        int64
    Status        int // 评论状态：1代表正常、2代表删除
    Content       string
    ParentId      int64
    Deep          int
    TopParentId   int64
    ParentPath    string
    ChildrenCount int
    VoteUp        int
    VoteDown      int
    RedditScore   float64
    CreateTime    time.Time

    user *User `db:"exclude"`
    link *Link `db:"exclude"`
}

// 评论的用户信息
func (c Comment) User() *User {
    if c.user == nil {
        c.user = User_GetById(c.UserId)
    }
    return c.user
}

// 评论的链接信息
func (c Comment) Link() *Link {
    if c.link == nil {
        c.link, _ = Link_GetById(c.LinkId)
    }
    return c.link
}

// 投票得分
func (c Comment) VoteScore() int {
    return c.VoteUp - c.VoteDown
}

func (c Comment) SinceTime() string {
    return utils.SmcTimeSince(c.CreateTime)
}

// 评论状态
func (c Comment) StatusName() string {
    name, ok := commentStatus[c.Status]
    if !ok {
        return "未知状态"
    }
    return name
}

type CommentList struct {
    Comment *Comment
    Childs  []*CommentList
}

// 保存评论到数据库，如果成功，则返回comment的id
func Comment_SaveMap(m map[string]interface{}) (int64, error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    var err error
    // 评论的链接存不存在？
    linkId := m["link_id"].(int64)
    link, err := Link_GetById(linkId)
    if err != nil {
        return 0, err
    } else if link.Id < 1 {
        return int64(0), errors.New("评论的链接不存在")
    }

    // 检查父评论是否存在
    var pComment *Comment
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
    } else if !ok {
        m["parent_id"] = int64(0)
    }

    // 路径相关
    if pComment == nil {
        m["parent_id"] = int64(0)
        m["top_parent_id"] = int64(0)
        m["parent_path"] = "/"
        m["deep"] = 0
    } else {
        m["parent_id"] = pComment.Id
        if pComment.TopParentId == 0 {
            m["top_parent_id"] = pComment.Id
        } else {
            m["top_parent_id"] = pComment.TopParentId
        }
        m["parent_path"] = fmt.Sprintf("%v%v/", pComment.ParentPath, pComment.Id)
        m["deep"] = pComment.Deep + 1
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
        // 更新Link的计数器
        IncCountById(db, Table_Link, linkId, "comment_count", 1)
        if pComment != nil {
            IncCountById(db, Table_Comment, pComment.Id, "children_count", 1)
        } else {
            IncCountById(db, Table_Link, linkId, "comment_root_count", 1)
        }

        // 通知评论用户
        userId := m["user_id"].(int64)
        if userId != link.UserId {
            comment := Comment{}
            comment.Id = id
            comment.UserId = userId
            comment.LinkId = linkId
            comment.ParentId = m["parent_id"].(int64)
            comment.CreateTime = m["create_time"].(time.Time)
            CommentForUser_Add(link.UserId, comment)

            Remind_Inc(link.UserId, REMIND_COMMENT)
        }
    }

    return id, nil
}

// 如果保存失败，则返回错误信息.
// @return success, commentId, errors
func Comment_SaveForm(f *form.Form, userId int64) (bool, int64, []string) {
    errorMsgs := make([]string, 0)
    var commentId int64
    if f.Valid() {
        m := f.CleanValues()
        m["user_id"] = userId

        id, err := Comment_SaveMap(m)
        if err != nil || id < 1 {
            errorMsgs = append(errorMsgs, golink.ERROR_DATABASE)
        } else {
            commentId = id
        }
    } else {
        errs := f.Errors()
        for _, v := range errs {
            errorMsgs = append(errorMsgs, v[1])
        }
    }
    if len(errorMsgs) < 1 {
        return true, commentId, nil
    }
    return false, commentId, errorMsgs
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
// @return: comments, total-count, err
func Comment_GetByPage(page, pagesize int, order string) ([]Comment, int64, error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    page, pagesize = utils.PageCheck(page, pagesize)

    qi := goku.SqlQueryInfo{}
    qi.Limit = pagesize
    qi.Offset = page * pagesize
    if order == "" {
        qi.Order = "id desc"
    } else {
        qi.Order = order
    }
    var comments []Comment
    err := db.GetStructs(&comments, qi)
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return nil, 0, err
    }

    total, err := db.Count("comment", "")
    if err != nil {
        goku.Logger().Errorln(err.Error())
    }

    return comments, total, nil
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
    qi.Order = "id desc"
    var comments []Comment
    err := db.GetStructs(&comments, qi)
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return nil
    }
    return comments
}

// 获取link的评论
func Comment_ForLink(linkId int64) []Comment {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    qi := goku.SqlQueryInfo{}
    qi.Where = "`link_id`=?"
    qi.Params = []interface{}{linkId}
    qi.Order = "id asc"
    var comments []Comment
    err := db.GetStructs(&comments, qi)
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return nil
    }
    return comments
}

// 获取排好序的link的评论
// @sort: hot, vote
func Comment_SortForLink(linkId int64, sort string) []*CommentList {
    comments := Comment_ForLink(linkId)
    if comments == nil {
        return nil
    }
    var cl []*CommentList
    switch sort {
    case golink.ORDER_TYPE_HOT: //热门
        cl = comment_SortByHot(comments)
    case golink.ORDER_TYPE_VOTE: //得分
    default:
        cl = comment_SortByHot(comments)
    }
    return cl
}

// TODO: 内存优化？
// TODO: 指针指得我头晕 =。=
// @comments: 按id升序排序的评论列表
func comment_SortByHot(comments []Comment) []*CommentList {
    if comments == nil {
        return nil
    }
    index := map[int64]*CommentList{}
    cl := make([]*CommentList, 0, 1)
    var pcl *[]*CommentList

    for j, _ := range comments {
        // c不能写在 for 里面，否则取地址的时候都是取到同一个地址
        c := comments[j]
        // 是否是回复评论
        if c.ParentId < 1 {
            pcl = &cl
        } else {
            // 查找父节点
            tempCl := index[c.ParentId]
            if tempCl.Childs == nil {
                tempCl.Childs = make([]*CommentList, 0, 1)
            }
            pcl = &tempCl.Childs
        }

        ncl := &CommentList{
            Comment: &c,
        }
        index[c.Id] = ncl
        if len(*pcl) > 0 {
            for i, _cl := range *pcl {
                if c.RedditScore > _cl.Comment.RedditScore {
                    if i == 0 {
                        *pcl = append([]*CommentList{ncl}, *pcl...)
                    } else {
                        *pcl = append((*pcl)[:i], append([]*CommentList{ncl}, (*pcl)[i:]...)...)
                    }
                    goto FEND
                }
            }
        }
        *pcl = append(*pcl, ncl)
    FEND:
    }

    return cl
}

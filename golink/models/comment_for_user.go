package models

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink/utils"
)

var table_CommentForUser string = "comment_for_user"

// 添加一条推送评论到被评论的用户,
func CommentForUser_Add(userId int64, comment Comment) error {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    m := map[string]interface{}{}

    m["user_id"] = userId
    m["comment_id"] = comment.Id
    m["link_id"] = comment.LinkId
    m["pcomment_id"] = comment.ParentId
    m["create_time"] = comment.CreateTime

    _, err := db.Insert(table_CommentForUser, m)
    if err != nil {
        goku.Logger().Errorln(err.Error())
    }
    return err
}

// 获取收到的评论列表
// @page: 从1开始
// @return: comments, total-count, err
func CommentForUser_GetByPage(userId int64, page, pagesize int, order string) ([]Comment, int64, error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    page, pagesize = utils.PageCheck(page, pagesize)

    qi := goku.SqlQueryInfo{}
    qi.Limit = pagesize
    qi.Offset = page * pagesize
    qi.Where = "cfu.user_id=?"
    qi.Join = " cfu INNER JOIN `comment` c ON cfu.comment_id=c.id"
    qi.Fields = `c.id, c.user_id, c.link_id, c.parent_path, c.children_count, c.top_parent_id,
                c.parent_id, c.deep, c.status, c.content, c.create_time, c.vote_up, c.vote_down, c.reddit_score`

    if order == "" {
        qi.Order = "create_time desc"
    } else {
        qi.Order = order
    }

    qi.Params = []interface{}{userId}
    rows, err := db.Select("comment_for_user", qi)

    if err != nil {
        goku.Logger().Errorln(err.Error())
        return nil, 0, err
    }
    defer rows.Close()
    comments := make([]Comment, 0)
    for rows.Next() {
        c := Comment{}
        err = rows.Scan(&c.Id, &c.UserId, &c.LinkId, &c.ParentPath, &c.ChildrenCount,
            &c.TopParentId, &c.ParentId, &c.Deep, &c.Status, &c.Content,
            &c.CreateTime, &c.VoteUp, &c.VoteDown, &c.RedditScore)
        if err != nil {
            goku.Logger().Errorln(err.Error())
            return nil, 0, err
        }
        comments = append(comments, c)
    }

    total, err := db.Count("comment_for_user", "user_id=?", userId)
    if err != nil {
        goku.Logger().Errorln(err.Error())
    }

    return comments, total, nil
}

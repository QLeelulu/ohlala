package models

import (
    "github.com/QLeelulu/goku"
    //"github.com/QLeelulu/ohlala/golink"
    "github.com/QLeelulu/ohlala/golink/utils"
)

//收藏link
func SaveUserFavorite(f map[string]interface{}) (error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()
    _, err := db.Insert("user_favorite_link", f)

    return err
}

//删除link
func DelUserFavorite(userId int64, linkId int64) (error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()
    _, err := db.Delete("user_favorite_link", "`user_id`=? AND `link_id`=?", userId, linkId)

    return err
}

// 获取由用户收藏的link
// @page: 从1开始
func FavoriteLink_ByUser(userId int64, page, pagesize int) []Link {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    page, pagesize = utils.PageCheck(page, pagesize)

    qi := goku.SqlQueryInfo{}
    qi.Fields = "l.id, l.user_id, l.title, l.context, l.topics, l.vote_up, l.vote_down, l.view_count, l.comment_count, l.create_time"
    qi.Join = " ufl INNER JOIN `link` l ON ufl.link_id=l.id"
    qi.Where = "ufl.user_id=?"
    qi.Params = []interface{}{userId}
    qi.Limit = pagesize
    qi.Offset = pagesize * page
    qi.Order = "ufl.create_time desc"

    rows, err := db.Select("user_favorite_link", qi)
	if err != nil {
        goku.Logger().Errorln(err.Error())
		return nil
	}
    links := make([]Link, 0)
    for rows.Next() {
        link := Link{}
        err = rows.Scan(&link.Id, &link.UserId, &link.Title, &link.Context, &link.Topics,
            &link.VoteUp, &link.VoteDown, &link.ViewCount, &link.CommentCount, &link.CreateTime)
        if err != nil {
            goku.Logger().Errorln(err.Error())
            return nil
        }
        links = append(links, link)
    }

    return links
}




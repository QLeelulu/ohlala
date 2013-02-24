package models

import (
    "database/sql"
    "fmt"
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink"
    "github.com/QLeelulu/ohlala/golink/utils"
    "time"
)

/**
 * 链接推送给网站首页(更新现有数据 )
 */
func Link_for_home_update(handleTime time.Time, db *goku.MysqlDB) error {

    sql := `UPDATE tui_link_for_handle H 
		INNER JOIN tui_link_for_home T ON H.insert_time<=? AND H.data_type=2 AND T.link_id=H.link_id 
		INNER JOIN link L ON L.id=H.link_id 
		SET T.score=CASE WHEN T.data_type=2 THEN L.reddit_score -- 热门 
		WHEN T.data_type=3 THEN L.dispute_score -- 热议
		WHEN T.data_type BETWEEN 10 AND 14 THEN L.dispute_score -- 热议
		ELSE L.vote_up-L.vote_down -- 投票 
		END;`

    _, err := db.Query(sql, handleTime)

    return err
}

/**
 * 链接推送给网站首页(热门)
 */
func Link_for_home_top(handleTime time.Time, db *goku.MysqlDB) error {

    sql := `INSERT ignore INTO tui_link_for_home(link_id,create_time,data_type,score) 
		( 
		SELECT H.link_id,H.create_time,2,L.reddit_score FROM tui_link_for_handle H 
		INNER JOIN link L ON H.insert_time<=? AND L.id=H.link_id 
		); `

    _, err := db.Query(sql, handleTime)

    return err
}

/**
 * 链接推送给网站首页(热议)[3:全部时间；10:这个小时；11:今天；12:这周；13:这个月；14:今年]
 */
func Link_for_home_hot_all(handleTime time.Time, db *goku.MysqlDB) error {

    err := link_for_home_hot(3, handleTime, db)
    if err == nil {
        err = link_for_home_hot(10, handleTime, db)
    }
    if err == nil {
        err = link_for_home_hot(11, handleTime, db)
    }
    if err == nil {
        err = link_for_home_hot(12, handleTime, db)
    }
    if err == nil {
        err = link_for_home_hot(13, handleTime, db)
    }
    if err == nil {
        err = link_for_home_hot(14, handleTime, db)
    }

    return err
}

/**
 * 链接推送给网站首页(热议)[3:全部时间；10:这个小时；11:今天；12:这周；13:这个月；14:今年]
 */
func link_for_home_hot(dataType int, handleTime time.Time, db *goku.MysqlDB) error {

    var t time.Time
    switch {
    case dataType == 10:
        t = utils.ThisHour()
    case dataType == 11:
        t = utils.ThisDate()
    case dataType == 12:
        t = utils.ThisWeek()
    case dataType == 13:
        t = utils.ThisMonth()
    case dataType == 14:
        t = utils.ThisYear()
    }

    var err error
    if dataType == 3 { //3:全部时间
        sql := `INSERT ignore INTO tui_link_for_home(link_id,create_time,data_type,score) 
			( 
			SELECT H.link_id,H.create_time,?,L.dispute_score FROM tui_link_for_handle H 
			INNER JOIN link L ON H.insert_time<=? AND L.id=H.link_id 
			); `

        _, err = db.Query(sql, dataType, handleTime)
    } else {
        sql := `INSERT ignore INTO tui_link_for_home(link_id,create_time,data_type,score) 
		( 
		SELECT H.link_id,H.create_time,?,L.dispute_score FROM tui_link_for_handle H 
		INNER JOIN link L ON H.insert_time<=? AND H.create_time>=? AND L.id=H.link_id 
		); `

        _, err = db.Query(sql, dataType, handleTime, t)
    }

    return err
}

/**
 * 链接推送给网站首页(投票)[投票时间范围: 4:全部时间；5:这个小时；6:今天；7:这周；8:这个月；9:今年]
 */
func Link_for_home_vote_all(handleTime time.Time, db *goku.MysqlDB) error {

    err := link_for_home_vote(4, handleTime, db)
    if err == nil {
        err = link_for_home_vote(5, handleTime, db)
    }
    if err == nil {
        err = link_for_home_vote(6, handleTime, db)
    }
    if err == nil {
        err = link_for_home_vote(7, handleTime, db)
    }
    if err == nil {
        err = link_for_home_vote(8, handleTime, db)
    }
    if err == nil {
        err = link_for_home_vote(9, handleTime, db)
    }

    return err
}

/**
 * 链接推送给网站首页(投票)[投票时间范围: 4:全部时间；5:这个小时；6:今天；7:这周；8:这个月；9:今年]
 */
func link_for_home_vote(dataType int, handleTime time.Time, db *goku.MysqlDB) error {

    var t time.Time
    switch {
    case dataType == 5:
        t = utils.ThisHour()
    case dataType == 6:
        t = utils.ThisDate()
    case dataType == 7:
        t = utils.ThisWeek()
    case dataType == 8:
        t = utils.ThisMonth()
    case dataType == 9:
        t = utils.ThisYear()
    }

    var err error
    if dataType == 4 { //4:全部时间
        sql := `INSERT ignore INTO tui_link_for_home(link_id,create_time,data_type,score) 
			( 
			SELECT H.link_id,H.create_time,?,L.vote_up-L.vote_down FROM tui_link_for_handle H 
			INNER JOIN link L ON H.insert_time<=? AND L.id=H.link_id 
			);`

        _, err = db.Query(sql, dataType, handleTime)
    } else {
        sql := `INSERT ignore INTO tui_link_for_home(link_id,create_time,data_type,score) 
			( 
			SELECT H.link_id,H.create_time,?,L.vote_up-L.vote_down FROM tui_link_for_handle H 
			INNER JOIN link L ON H.insert_time<=? AND L.create_time>=? AND L.id=H.link_id  
			); `

        _, err = db.Query(sql, dataType, handleTime, t)
    }

    return err
}

func Del_link_for_home_all(db *goku.MysqlDB) error {

    err := del_link_for_home("data_type=2", "score DESC,link_id DESC", db)
	if err == nil {
		_, err = db.Query(`DELETE FROM tui_link_for_home WHERE ((data_type=10 OR data_type=5) AND create_time<?) OR 
							((data_type=11 OR data_type=6) AND create_time<?) OR 
							((data_type=12 OR data_type=7) AND create_time<?) OR 
							((data_type=13 OR data_type=8) AND create_time<?) OR 
			((data_type=14 OR data_type=9) AND create_time<?)`, utils.ThisHour(), utils.ThisDate(), utils.ThisWeek(), utils.ThisMonth(), utils.ThisYear())
	}
    if err == nil {
		if err == nil {
        	err = del_link_for_home("data_type IN(3,10,11,12,13,14)", "score DESC,link_id DESC", db)
		}
    }
    if err == nil {
		
        err = del_link_for_home("data_type IN(4,5,6,7,8,9)", "score DESC,link_id DESC", db)
    }

    return err
}

/** 删除`tui_link_for_home`
 * 热门, whereDataType:data_type=2    orderName:score DESC,link_id DESC
 * 热议, whereDataType:data_type IN(3,10,11,12,13,14)    orderName:score desc,link_id DESC
 * 投票, whereDataType:data_type IN(4,5,6,7,8,9)    orderName:score DESC,link_id DESC
 */
func del_link_for_home(whereDataType string, orderName string, db *goku.MysqlDB) error {

    sql := fmt.Sprintf(`SELECT data_type, tcount - %d AS del_count FROM 
		(SELECT link_id,data_type,COUNT(1) AS tcount FROM tui_link_for_home WHERE %s GROUP BY data_type) T
		WHERE T.tcount>%d;`, LinkMaxCount, whereDataType, LinkMaxCount)

    delSqlCreate := `INSERT ignore INTO tui_link_temporary_delete(id)
		( 
		SELECT link_id FROM tui_link_for_home WHERE data_type=%d ORDER BY ` + orderName + ` LIMIT %d,%d 
		); `
    delSqlDelete := `DELETE T FROM tui_link_temporary_delete D INNER JOIN tui_link_for_home T ON T.data_type=%d
		AND T.link_id=D.id; `

    var delCount int64
    var dataType int
    rows, err := db.Query(sql)
    if err == nil {
        for rows.Next() {
            rows.Scan(&dataType, &delCount)
            db.Query("DELETE FROM tui_link_temporary_delete;")
            db.Query(fmt.Sprintf(delSqlCreate, dataType, LinkMaxCount, delCount))
            db.Query(fmt.Sprintf(delSqlDelete, dataType))
        }
    }

    return err
}

// @page: 从1开始
// @orderType: 排序类型, hot:热门, hotc:热议, time:最新, vote:投票得分, ctvl:争议
// @dataType: 2:热门; 
//            3:争议[3:全部时间；10:这个小时；11:今天；12:这周；13:这个月；14:今年]; 
//            [投票时间范围: 4:全部时间；5:这个小时；6:今天；7:这周；8:这个月；9:今年]
func LinkForHome_GetByPage(orderType string, dataType, page, pagesize int) ([]Link, error) {
    if page < 1 {
        page = 1
    }
    page = page - 1
    if pagesize == 0 {
        pagesize = 20
    }

    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    var rows *sql.Rows
    var err error
    if orderType != golink.ORDER_TYPE_TIME {
        qi := goku.SqlQueryInfo{}
        qi.Fields = "l.id, l.user_id, l.title, l.context, l.topics, l.vote_up, l.vote_down, l.view_count, l.comment_count, l.create_time, l.status"
        qi.Join = " lfh INNER JOIN `link` l ON lfh.link_id=l.id"
        qi.Where = "lfh.data_type=?"
        qi.Limit = pagesize
        qi.Offset = pagesize * page
        switch orderType {
        case golink.ORDER_TYPE_HOTC: // 热议
            qi.Order = "l.comment_count desc, lfh.link_id desc"
            dataType = 2
        case golink.ORDER_TYPE_CTVL: // 争议
            qi.Order = "lfh.score DESC,lfh.link_id desc"
            dataType = 3
        case golink.ORDER_TYPE_VOTE: // 得分
            qi.Order = "lfh.score desc, lfh.link_id desc"
            dataType = 4
        default:
            qi.Order = "lfh.score desc, lfh.link_id desc"
            dataType = 2
        }

        qi.Params = []interface{}{dataType}
        rows, err = db.Select("tui_link_for_home", qi)

        if err != nil {
            goku.Logger().Errorln(err.Error())
            return nil, err
        }
        defer rows.Close()
    } else {
        qi := goku.SqlQueryInfo{}
        qi.Fields = "id, user_id, title, context, topics, vote_up, vote_down, view_count, comment_count, create_time, status"
        qi.Limit = pagesize
        qi.Offset = pagesize * page
        qi.Order = "id desc"

        rows, err = db.Select("link", qi)

        if err != nil {
            goku.Logger().Errorln(err.Error())
            return nil, err
        }
        defer rows.Close()
    }

    links := make([]Link, 0)
    for rows.Next() {
        link := Link{}
        err = rows.Scan(&link.Id, &link.UserId, &link.Title, &link.Context, &link.Topics,
            &link.VoteUp, &link.VoteDown, &link.ViewCount, &link.CommentCount, &link.CreateTime, &link.Status)
        if err != nil {
            goku.Logger().Errorln(err.Error())
            return nil, err
        }
        links = append(links, link)
    }
    return links, nil
}

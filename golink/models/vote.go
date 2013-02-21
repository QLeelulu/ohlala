package models

import (
    //"database/sql"
    "github.com/QLeelulu/goku"
    //"github.com/QLeelulu/ohlala/golink"
    "github.com/QLeelulu/ohlala/golink/utils"
    "time"
//"fmt"
)

type Vote struct {
    Id      int64
    VoteNum int64
    Success bool
    Errors  string
}

func VoteLink(linkId int64, userId int64, score int, siteRunTime string) *Vote {

    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    var update bool = false
    var vote *Vote = &Vote{0, 0, false, ""}
    checkRows, checkErr := db.Query("SELECT `user_id` FROM `link` WHERE `id` = ? LIMIT 0,1", linkId)
    if checkErr != nil {
        goku.Logger().Errorln(checkErr.Error())
        vote.Errors = "数据库错误"
    } else if checkRows.Next() == false {
        vote.Errors = "参数错误"
    } else {
        var luid int64
        checkRows.Scan(&luid)
        if luid == userId {
            vote.Errors = "不可以对自己投票"
        }
    }
    if vote.Errors != "" {
        return vote
    }

    rows, err := db.Query("SELECT score FROM `link_support_record` WHERE `link_id` = ? AND `user_id` = ? LIMIT 0,1", linkId, userId)
    if err == nil {
	
		var upVote int64
		var downVote int64
		var createTime time.Time
        if rows.Next() { //投过票的情况
            var scoreTemp int
            rows.Scan(&scoreTemp)

            //已投了支持，再投反对
            if scoreTemp == 1 && score == -1 {

                update = true
				rows, err = db.Query("SELECT vote_up,vote_down,create_time FROM `link` WHERE id=?", linkId)
				rows.Next()
				rows.Scan(&upVote, &downVote, &createTime)
				upVote = upVote-1
				downVote = downVote+1
				score := utils.LinkSortAlgorithm(createTime, upVote, downVote)
                db.Query("UPDATE `link` SET vote_up=?,vote_down=?,reddit_score=? WHERE id=?;", upVote, downVote, score, linkId)
                db.Query("UPDATE `link_support_record` SET score=-1,vote_time=NOW() WHERE `link_id` = ? AND `user_id` = ?", linkId, userId)

            } else if scoreTemp == -1 && score == 1 { //已投了反对，再投支持

                update = true
				rows, err = db.Query("SELECT vote_up,vote_down,create_time FROM `link` WHERE id=?", linkId)
				rows.Next()
				rows.Scan(&upVote, &downVote, &createTime)
				upVote = upVote+1
				downVote = downVote-1
				score := utils.LinkSortAlgorithm(createTime, upVote, downVote)
                db.Query("UPDATE `link` SET vote_up=?,vote_down=?,reddit_score=? WHERE id=?;", upVote, downVote, score, linkId)
                db.Query("UPDATE `link_support_record` SET score=1,vote_time=NOW() WHERE `link_id` = ? AND `user_id` = ?", linkId, userId)

            } else {
                vote.Errors = "您已对此投票"
            }

        } else { //没投过票的情况

            update = true
			rows, err = db.Query("SELECT vote_up,vote_down,create_time FROM `link` WHERE id=?", linkId)
			rows.Next()
			rows.Scan(&upVote, &downVote, &createTime)
            if score == 1 {
				upVote = upVote+1
				score := utils.LinkSortAlgorithm(createTime, upVote, downVote)
                db.Query("UPDATE `link` SET vote_up=?,reddit_score=? WHERE id=?;", upVote, score, linkId)

            } else {
				downVote = downVote+1
				score := utils.LinkSortAlgorithm(createTime, upVote, downVote)
                db.Query("UPDATE `link` SET vote_down=?,reddit_score=? WHERE `id`=?;", downVote, score, linkId)

            }
            db.Query("INSERT INTO `link_support_record` (link_id, user_id, score,vote_time) VALUES (?, ?, ?, NOW())", linkId, userId, score)

        }

        if update {
            vote.Id = linkId
            vote.VoteNum = upVote - downVote
            vote.Success = true
            // 存入`tui_link_for_handle` 链接处理队列表
            db.Query("INSERT ignore INTO tui_link_for_handle(link_id,create_time,user_id,insert_time,data_type) VALUES (?, ?, ?, NOW(), ?)",
                linkId, createTime, userId, 2)
        }
    }
    if err != nil {
		vote.Success = false
        goku.Logger().Errorln(err.Error())
        vote.Errors = "数据库错误"
    }

    return vote
}

func VoteComment(commentId int64, userId int64, score int, siteRunTime string) *Vote {

    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    var vote *Vote = &Vote{0, 0, false, ""}
    var updateChildrenScore bool = false

    checkRows, checkErr := db.Query("SELECT `user_id` FROM `comment` WHERE `id` = ? LIMIT 0,1", commentId) //reddit_score
    if checkErr != nil {
        goku.Logger().Errorln(checkErr.Error())
        vote.Errors = "数据库错误"
    } else if checkRows.Next() == false {
        vote.Errors = "参数错误"
    } else {
        var luid int64
        checkRows.Scan(&luid)
        if luid == userId {
            vote.Errors = "不可以对自己投票"
        }
    }
    if vote.Errors != "" {
        return vote
    }

    rows, err := db.Query("SELECT score FROM `comment_support_record` WHERE `comment_id` = ? AND `user_id` = ? LIMIT 0,1", commentId, userId)
    if err == nil {
		var upVote, downVote int64
        if rows.Next() { //投过票的情况
            var scoreTemp int
            rows.Scan(&scoreTemp)
            //已投了支持，再投反对
            if scoreTemp == 1 && score == -1 {

                updateChildrenScore = true
				rows, err = db.Query("SELECT vote_up,vote_down FROM `comment` WHERE id=?", commentId)
				rows.Next()
				rows.Scan(&upVote, &downVote)
				upVote -= 1
				downVote += 1
				score := utils.CommentSortAlgorithm(upVote, downVote)
                db.Query("UPDATE `comment` SET vote_up=?,vote_down=?,reddit_score=? WHERE id=?;", upVote, downVote, score, commentId)
                db.Query("UPDATE `comment_support_record` SET score=-1,vote_time=NOW() WHERE `comment_id` = ? AND `user_id` = ?", commentId, userId)

            } else if scoreTemp == -1 && score == 1 { //已投了反对，再投支持

                updateChildrenScore = true
				rows, err = db.Query("SELECT vote_up,vote_down FROM `comment` WHERE id=?", commentId)
				rows.Next()
				rows.Scan(&upVote, &downVote)
				upVote += 1
				downVote -= 1
				score := utils.CommentSortAlgorithm(upVote, downVote)
                db.Query("UPDATE `comment` SET vote_down=?,vote_up=?,reddit_score=? WHERE `id`=?;", downVote, upVote, score, commentId)
                db.Query("UPDATE `comment_support_record` SET score=1,vote_time=NOW() WHERE `comment_id` = ? AND `user_id` = ?", commentId, userId)

            } else {
                vote.Errors = "您已对此投票"
            }

        } else { //没投过票的情况

            updateChildrenScore = true
			rows, err = db.Query("SELECT vote_up,vote_down FROM `comment` WHERE id=?", commentId)
			rows.Next()
			rows.Scan(&upVote, &downVote)
            if score == 1 {
				upVote += 1
				score := utils.CommentSortAlgorithm(upVote, downVote)
                db.Query("UPDATE `comment` SET vote_up=?,reddit_score=? WHERE id=?;", upVote, score, commentId)

            } else {
				downVote += 1
				score := utils.CommentSortAlgorithm(upVote, downVote)
                db.Query("UPDATE `comment` SET vote_down=?,reddit_score=? WHERE `id`=?;", downVote, score, commentId)

            }
            db.Query("INSERT INTO `comment_support_record` (comment_id, user_id, score, vote_time) VALUES (?, ?, ?,NOW())", commentId, userId, score)

        }

        if updateChildrenScore {

            //rows, err = db.Query("SELECT vote_up-vote_down AS vote FROM `comment` WHERE `id` = ? LIMIT 0,1", commentId) //, reddit_score,link_id
            //if err == nil && rows.Next() {
                //var voteNum int64 = 0
                vote.Id = commentId
                //rows.Scan(&voteNum) //, &childScore, &linkId
                vote.VoteNum = upVote - downVote
                vote.Success = true
            //}
        }
    }
    if err != nil {
        goku.Logger().Errorln(err.Error())
        vote.Success = false
        vote.Errors = "数据库错误"
    }

    return vote

}

//获取link和comment的投票记录
func GetVoteRecordByUser(userId int64, page int, pagesize int) {
    /*
    	sql := `SELECT * FROM (
    (SELECT L.id, 'link' AS record_type,LSR.score AS vote_score,LSR.vote_time FROM link_support_record LSR 
    INNER JOIN link L ON LSR.link_id=L.id AND LSR.user_id=1 ORDER BY vote_time DESC LIMIT 0,200)
    UNION ALL
    (SELECT C.id,'comment' AS record_type,CSR.score AS vote_score,CSR.vote_time FROM comment_support_record CSR
    INNER JOIN comment C ON CSR.comment_id=C.id AND CSR.user_id=1 ORDER BY CSR.vote_time DESC LIMIT 0,200))T ORDER BY T.vote_time DESC LIMIT ?,?`
    */
    if page < 1 {
        page = 1
    }
    page = page - 1
    if pagesize == 0 {
        pagesize = 20
    }
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    //rows, err := db.Query(sql, pagesize * page, pagesize)

    //if err != nil {
    //goku.Logger().Errorln(err.Error())
    //return nil, err
    //}

}

package models

import (
    //"database/sql"
    "github.com/QLeelulu/goku"
    "time"
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

        if rows.Next() { //投过票的情况
            var scoreTemp int
            rows.Scan(&scoreTemp)

            //已投了支持，再投反对
            if scoreTemp == 1 && score == -1 {
                update = true
                db.Query("UPDATE `link` SET vote_up=vote_up-1,vote_down=vote_down+1,reddit_score=LOG10(ABS(vote_up-vote_down)) +  ( CASE WHEN vote_up=vote_down THEN 0 ELSE (IF(vote_up-vote_down>0, 1, -1) * TIMESTAMPDIFF(SECOND,'"+siteRunTime+"',create_time))/45000 END ) WHERE id=?;", linkId)
                db.Query("UPDATE `link_support_record` SET score=-1 WHERE `link_id` = ? AND `user_id` = ?", linkId, userId)
            } else if scoreTemp == -1 && score == 1 { //已投了反对，再投支持
                update = true
                db.Query("UPDATE `link` SET vote_down=vote_down-1,vote_up=vote_up+1,reddit_score=LOG10(ABS(vote_up-vote_down)) +  ( CASE WHEN vote_up=vote_down THEN 0 ELSE (IF(vote_up-vote_down>0, 1, -1) * TIMESTAMPDIFF(SECOND,'"+siteRunTime+"',create_time))/45000 END ) WHERE `id`=?;", linkId)
                db.Query("UPDATE `link_support_record` SET score=1 WHERE `link_id` = ? AND `user_id` = ?", linkId, userId)
            } else {
                vote.Errors = "您已对此投票"
            }

        } else { //没投过票的情况

            update = true
            if score == 1 {
                db.Query("UPDATE `link` SET vote_up=vote_up+1,reddit_score=LOG10(ABS(vote_up-vote_down)) +  ( CASE WHEN vote_up=vote_down THEN 0 ELSE (IF(vote_up-vote_down>0, 1, -1) * TIMESTAMPDIFF(SECOND,'"+siteRunTime+"',create_time))/45000 END ) WHERE id=?;", linkId)
            } else {
                db.Query("UPDATE `link` SET vote_down=vote_down+1,reddit_score=LOG10(ABS(vote_up-vote_down)) +  ( CASE WHEN vote_up=vote_down THEN 0 ELSE (IF(vote_up-vote_down>0, 1, -1) * TIMESTAMPDIFF(SECOND,'"+siteRunTime+"',create_time))/45000 END ) WHERE `id`=?;", linkId)
            }
            db.Query("INSERT INTO `link_support_record` (link_id, user_id, score) VALUES (?, ?, ?)", linkId, userId, score)

        }

        if update {
            rows, err = db.Query("SELECT vote_up-vote_down AS vote,create_time FROM `link` WHERE `id` = ? LIMIT 0,1", linkId)
            if err == nil && rows.Next() {
                var createTime time.Time
                var voteNum int64 = 0
                vote.Id = linkId
                rows.Scan(&voteNum, &createTime)
                vote.VoteNum = voteNum
                vote.Success = true

                // 存入`tui_link_for_handle` 链接处理队列表
                db.Query("INSERT ignore INTO tui_link_for_handle(link_id,create_time,user_id,insert_time,data_type) VALUES (?, ?, ?, NOW(), ?)",
                    linkId, createTime, userId, 2)
            }
        }
    }
    if err != nil {
        goku.Logger().Errorln(err.Error())
        vote.Errors = "数据库错误"
    }

    return vote
}

func VoteComment(commentId int64, userId int64, score int, siteRunTime string) *Vote { //topId int64, 

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

        if rows.Next() { //投过票的情况
            var scoreTemp int
            rows.Scan(&scoreTemp)
            //已投了支持，再投反对
            if scoreTemp == 1 && score == -1 {
                updateChildrenScore = true
                db.Query("UPDATE `comment` SET vote_up=vote_up-1,vote_down=vote_down+1,reddit_score=LOG10(ABS(vote_up-vote_down)) +  ( CASE WHEN vote_up=vote_down THEN 0 ELSE (IF(vote_up-vote_down>0, 1, -1) * TIMESTAMPDIFF(SECOND,'"+siteRunTime+"',create_time))/45000 END ) WHERE id=?;", commentId)
                db.Query("UPDATE `comment_support_record` SET score=-1 WHERE `comment_id` = ? AND `user_id` = ?", commentId, userId)
            } else if scoreTemp == -1 && score == 1 { //已投了反对，再投支持
                updateChildrenScore = true
                db.Query("UPDATE `comment` SET vote_down=vote_down-1,vote_up=vote_up+1,reddit_score=LOG10(ABS(vote_up-vote_down)) +  ( CASE WHEN vote_up=vote_down THEN 0 ELSE (IF(vote_up-vote_down>0, 1, -1) * TIMESTAMPDIFF(SECOND,'"+siteRunTime+"',create_time))/45000 END ) WHERE `id`=?;", commentId)
                db.Query("UPDATE `comment_support_record` SET score=1 WHERE `comment_id` = ? AND `user_id` = ?", commentId, userId)
            } else {
                vote.Errors = "您已对此投票"
            }

        } else { //没投过票的情况

            updateChildrenScore = true
            if score == 1 {
                db.Query("UPDATE `comment` SET vote_up=vote_up+1,reddit_score=LOG10(ABS(vote_up-vote_down)) +  ( CASE WHEN vote_up=vote_down THEN 0 ELSE (IF(vote_up-vote_down>0, 1, -1) * TIMESTAMPDIFF(SECOND,'"+siteRunTime+"',create_time))/45000 END ) WHERE id=?;", commentId)
            } else {
                db.Query("UPDATE `comment` SET vote_down=vote_down+1,reddit_score=LOG10(ABS(vote_up-vote_down)) +  ( CASE WHEN vote_up=vote_down THEN 0 ELSE (IF(vote_up-vote_down>0, 1, -1) * TIMESTAMPDIFF(SECOND,'"+siteRunTime+"',create_time))/45000 END ) WHERE `id`=?;", commentId)
            }
            db.Query("INSERT INTO `comment_support_record` (comment_id, user_id, score) VALUES (?, ?, ?)", commentId, userId, score)

        }

        if updateChildrenScore {
            //var childScore float64
            //var linkId int64
            //var oldChildScore float64

            rows, err = db.Query("SELECT vote_up-vote_down AS vote FROM `comment` WHERE `id` = ? LIMIT 0,1", commentId) //, reddit_score,link_id
            if err == nil && rows.Next() {
                var voteNum int64 = 0
                vote.Id = commentId
                rows.Scan(&voteNum) //, &childScore, &linkId
                vote.VoteNum = voteNum
                vote.Success = true
                //if topId > 0 {
                //oldData.Scan(&oldChildScore)
                //db.Query("UPDATE `comment` SET `children_reddit_score`=`children_reddit_score` - ? + ? WHERE `id` = ?", oldChildScore, childScore, topId)
                //}
                //db.Query("UPDATE `link` SET `comment_reddit_score`=`comment_reddit_score` - ? + ? WHERE `id` = ?", oldChildScore, childScore, linkId)
            }
        }
    }
    if err != nil {
        goku.Logger().Errorln(err.Error())
        vote.Errors = "数据库错误"
    }

    return vote

}

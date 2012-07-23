package models

import (
    //"database/sql"
    "github.com/QLeelulu/goku"
    "fmt"
)

type Vote struct {
    Id       int64
    VoteNum  int64
    Result   bool

}

func VoteLink(linkId int64, userId int64, score int, siteRunTime string) *Vote {

    var db *goku.MysqlDB = GetDB()
    defer db.Close()
    rows, err := db.Query("SELECT score FROM `link_support_record` WHERE `link_id` = ? AND `user_id` = ? LIMIT 0,1", linkId, userId)
    var vote *Vote = &Vote{0, 0, false}
    if err == nil {
	
	if rows.Next() { //投过票的情况
		var scoreTemp int
		rows.Scan(&scoreTemp)
fmt.Println(scoreTemp)
		//已投了支持，再投反对
		if scoreTemp == 1 && score == -1 {
		    db.Query("UPDATE `link` SET vote_up=vote_up-1,vote_down=vote_down+1,reddit_score=LOG10(ABS(vote_up-vote_down)) +  ( CASE WHEN vote_up=vote_down THEN 0 ELSE (IF(vote_up-vote_down>0, 1, -1) * TIMESTAMPDIFF(SECOND,'" + siteRunTime + "',create_time))/45000 END ) WHERE id=?;", linkId)
		    db.Query("UPDATE `link_support_record` SET score=-1 WHERE `link_id` = ? AND `user_id` = ?", linkId, userId)
		} else if scoreTemp == -1 && score == 1 { //已投了反对，再投支持
		    db.Query("UPDATE `link` SET vote_down=vote_down-1,vote_up=vote_up+1,reddit_score=LOG10(ABS(vote_up-vote_down)) +  ( CASE WHEN vote_up=vote_down THEN 0 ELSE (IF(vote_up-vote_down>0, 1, -1) * TIMESTAMPDIFF(SECOND,'" + siteRunTime + "',create_time))/45000 END ) WHERE `id`=?;", linkId)
		    db.Query("UPDATE `link_support_record` SET score=1 WHERE `link_id` = ? AND `user_id` = ?", linkId, userId)
		}
		
	} else { //没投过票的情况
		
		if score == 1 {
		    db.Query("UPDATE `link` SET vote_up=vote_up+1,reddit_score=LOG10(ABS(vote_up-vote_down)) +  ( CASE WHEN vote_up=vote_down THEN 0 ELSE (IF(vote_up-vote_down>0, 1, -1) * TIMESTAMPDIFF(SECOND,'" + siteRunTime + "',create_time))/45000 END ) WHERE id=?;", linkId)
		} else {
		    db.Query("UPDATE `link` SET vote_down=vote_down+1,reddit_score=LOG10(ABS(vote_up-vote_down)) +  ( CASE WHEN vote_up=vote_down THEN 0 ELSE (IF(vote_up-vote_down>0, 1, -1) * TIMESTAMPDIFF(SECOND,'" + siteRunTime + "',create_time))/45000 END ) WHERE `id`=?;", linkId)
		}
		db.Query("INSERT INTO `link_support_record` (link_id, user_id, score) VALUES (?, ?, ?)", linkId, userId, score)

	}

	rows, err := db.Query("SELECT vote_up-vote_down AS vote FROM `link` WHERE `id` = ? LIMIT 0,1", linkId)
	if err == nil && rows.Next() {
		var voteNum int64 = 0
		vote.Id = linkId
		rows.Scan(&voteNum)
		vote.VoteNum = voteNum
		vote.Result = true
	}
    }
    
    return vote
}

func VoteComment() {




}

package models

import (
    "github.com/QLeelulu/goku"
	"github.com/QLeelulu/ohlala/golink/utils"
    "time"
)

/**
 * 链接推送给话题(最新)
 */
func link_for_topic_later(handleTime time.Time, db *goku.MysqlDB) error {

	sql := `INSERT ignore INTO tui_link_for_topic_later(topic_id,link_id,create_time) 
		( 
		SELECT TL.topic_id,H.link_id,H.create_time FROM tui_link_for_handle H 
		INNER JOIN topic_link TL ON H.data_type=1 AND H.insert_time<=? AND H.link_id=TL.link_id 
		);`
	_, err := db.Query(sql, handleTime)

	return err
}

/**
 * 链接推送给话题(热门)
 */
func link_for_topic_top(handleTime time.Time, db *goku.MysqlDB) error {

	sql := `UPDATE tui_link_for_handle H 
		INNER JOIN tui_link_for_topic_top T ON H.data_type=2 AND H.insert_time<=? AND T.link_id=H.link_id 
		INNER JOIN link L ON L.id=H.link_id 
		SET T.reddit_score=L.reddit_score; 

		INSERT ignore INTO tui_link_for_topic_top(topic_id,link_id,create_time,reddit_score) 
		( 
		SELECT TL.topic_id,H.link_id,H.create_time,L.reddit_score FROM tui_link_for_handle H 
		INNER JOIN topic_link TL ON H.insert_time<=? AND H.link_id=TL.link_id 
		INNER JOIN link L ON L.id=H.link_id 
		); `
	_, err := db.Query(sql, handleTime)

	return err
}

/**
 * 链接推送给话题(热议)全部时间:1
 */
func link_for_topic_hop_all(handleTime time.Time, db *goku.MysqlDB) error {

	sql := `UPDATE tui_link_for_handle H 
		INNER JOIN tui_link_for_topic_hot TH ON H.data_type=2 AND H.insert_time<=? AND TH.link_id=H.link_id 
		INNER JOIN link L ON L.id=H.link_id 
		SET TH.vote_abs_score=ABS(L.vote_up-L.vote_down),TH.vote_add_score=(L.vote_up+L.vote_down); 

		INSERT ignore INTO tui_link_for_topic_hot(topic_id,link_id,create_time,vote_abs_score,vote_add_score,time_type) 
		( 
		SELECT TL.topic_id,H.link_id,H.create_time,ABS(L.vote_up-L.vote_down) AS vote_abs_score, 
		L.vote_up+L.vote_down AS vote_add_score,1 FROM tui_link_for_handle H 
		INNER JOIN topic_link TL ON H.insert_time<=? AND H.link_id=TL.link_id 
		INNER JOIN link L ON L.id=H.link_id 
		); `
	_, err := db.Query(sql, handleTime)

	return err
}
/**
 * 链接推送给话题(热议) 2:这个小时；3:今天；4:这周；5:这个月；6:今年
 */
func link_for_topic_hop_time(timeType int, handleTime time.Time, db *goku.MysqlDB) error {
	
	var t time.Time
    switch {
    case timeType == 2:
		t = utils.ThisHour()
    case timeType == 3:
		t = utils.ThisDate()
    case timeType == 4:
		t = utils.ThisWeek()
    case timeType == 5:
		t = utils.ThisMonth()
    case timeType == 6:
		t = utils.ThisYear()
    }

	sql := `INSERT ignore INTO tui_link_for_topic_hot(topic_id,link_id,create_time,vote_abs_score,vote_add_score,time_type) 
		( 
		SELECT TL.topic_id,H.link_id,H.create_time,ABS(L.vote_up-L.vote_down) AS vote_abs_score, 
		L.vote_up+L.vote_down AS vote_add_score,? AS time_type FROM tui_link_for_handle H 
		INNER JOIN topic_link TL ON H.insert_time<=? AND H.create_time>=? AND H.link_id=TL.link_id AND
		INNER JOIN link L ON L.id=H.link_id
		);`
	_, err := db.Query(sql, timeType, handleTime, t)

	return err
}

/**
 * 链接推送给话题(投票)全部时间:1
 */
func link_for_topic_vote_all(handleTime time.Time, db *goku.MysqlDB) error {

	sql := `UPDATE tui_link_for_handle H 
		INNER JOIN tui_link_for_topic_vote V ON H.data_type=2 AND V.link_id=H.link_id 
		INNER JOIN link L ON L.id=H.link_id 
		SET V.vote=(L.vote_up-L.vote_down); 

		INSERT ignore INTO tui_link_for_topic_vote(topic_id,link_id,create_time,time_type,vote) 
		( 
		SELECT TL.topic_id,H.link_id,H.create_time,1 AS time_type, 
		L.vote_up-L.vote_down AS vote FROM tui_link_for_handle H 
		INNER JOIN topic_link TL ON H.insert_time<=? AND H.link_id=TL.link_id 
		INNER JOIN link L ON L.id=H.link_id 
		);`
	_, err := db.Query(sql, handleTime)

	return err
}

/**
 * 链接推送给话题(投票) 2:这个小时；3:今天；4:这周；5:这个月；6:今年
 */
func link_for_topic_vote_time(timeType int, handleTime time.Time, db *goku.MysqlDB) error {
	
	var t time.Time
    switch {
    case timeType == 2:
		t = utils.ThisHour()
    case timeType == 3:
		t = utils.ThisDate()
    case timeType == 4:
		t = utils.ThisWeek()
    case timeType == 5:
		t = utils.ThisMonth()
    case timeType == 6:
		t = utils.ThisYear()
    }

	sql := `INSERT ignore INTO tui_link_for_topic_vote(topic_id,link_id,create_time,time_type,vote) 
		( 
		SELECT TL.topic_id,H.link_id,H.create_time,? AS time_type, 
		L.vote_up-L.vote_down AS vote FROM tui_link_for_handle H 
		INNER JOIN topic_link TL ON H.insert_time<=? AND H.create_time>=? AND H.link_id=TL.link_id 
		INNER JOIN link L ON L.id=H.link_id
		);`
	_, err := db.Query(sql, timeType, handleTime, t)

	return err
}






























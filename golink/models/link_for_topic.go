package models

import (
    "github.com/QLeelulu/goku"
	"github.com/QLeelulu/ohlala/golink/utils"
    "time"
	"fmt"
)

const (
    LinkMaxCount = 2//10000 // 队列长度
    HandleCount = 1//00 // 每次处理的数据
)





/**
 * 链接推送给话题(最新)
 */
func Link_for_topic_later(handleTime time.Time, db *goku.MysqlDB) error {

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
func Link_for_topic_top(handleTime time.Time, db *goku.MysqlDB) error {

	sql := `UPDATE tui_link_for_handle H 
		INNER JOIN tui_link_for_topic_top T ON H.data_type=2 AND H.insert_time<=? AND T.link_id=H.link_id 
		INNER JOIN link L ON L.id=H.link_id 
		SET T.reddit_score=L.reddit_score;`
	_, err := db.Query(sql, handleTime)

	if err == nil {
		sql = `INSERT ignore INTO tui_link_for_topic_top(topic_id,link_id,create_time,reddit_score) 
			( 
			SELECT TL.topic_id,H.link_id,H.create_time,L.reddit_score FROM tui_link_for_handle H 
			INNER JOIN topic_link TL ON H.insert_time<=? AND H.link_id=TL.link_id 
			INNER JOIN link L ON L.id=H.link_id 
			);`
		_, err = db.Query(sql, handleTime)
	}

	return err
}

/**
 * 链接推送给话题(热议)全部时间:1
 */
func Link_for_topic_hot_all(handleTime time.Time, db *goku.MysqlDB) error {

	sql := `UPDATE tui_link_for_handle H 
		INNER JOIN tui_link_for_topic_hot TH ON H.data_type=2 AND H.insert_time<=? AND TH.link_id=H.link_id 
		INNER JOIN link L ON L.id=H.link_id 
		SET TH.vote_abs_score=ABS(L.vote_up-L.vote_down),TH.vote_add_score=(L.vote_up+L.vote_down);`
	_, err := db.Query(sql, handleTime)

	if err == nil {
		sql = `INSERT ignore INTO tui_link_for_topic_hot(topic_id,link_id,create_time,vote_abs_score,vote_add_score,time_type) 
			( 
			SELECT TL.topic_id,H.link_id,H.create_time,ABS(L.vote_up-L.vote_down) AS vote_abs_score, 
			L.vote_up+L.vote_down AS vote_add_score,1 FROM tui_link_for_handle H 
			INNER JOIN topic_link TL ON H.insert_time<=? AND H.link_id=TL.link_id 
			INNER JOIN link L ON L.id=H.link_id 
			);`
		_, err = db.Query(sql, handleTime)
	}

	if err == nil {
		err = link_for_topic_hop_time(2, handleTime, db)
	}
	if err == nil {
		err = link_for_topic_hop_time(3, handleTime, db)
	}
	if err == nil {
		err = link_for_topic_hop_time(4, handleTime, db)
	}
	if err == nil {
		err = link_for_topic_hop_time(5, handleTime, db)
	}
	if err == nil {
		err = link_for_topic_hop_time(6, handleTime, db)
	}

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
		INNER JOIN topic_link TL ON H.insert_time<=? AND H.create_time>=? AND H.link_id=TL.link_id
		INNER JOIN link L ON L.id=H.link_id
		);`
	_, err := db.Query(sql, timeType, handleTime, t)

	return err
}

/**
 * 链接推送给话题(投票)全部时间:1
 */
func Link_for_topic_vote_all(handleTime time.Time, db *goku.MysqlDB) error {

	sql := `UPDATE tui_link_for_handle H 
		INNER JOIN tui_link_for_topic_vote V ON H.data_type=2 AND H.insert_time<=? AND V.link_id=H.link_id 
		INNER JOIN link L ON L.id=H.link_id 
		SET V.vote=(L.vote_up-L.vote_down);`
	_, err := db.Query(sql, handleTime)
	
	if err == nil {
		sql = `INSERT ignore INTO tui_link_for_topic_vote(topic_id,link_id,create_time,time_type,vote) 
			( 
			SELECT TL.topic_id,H.link_id,H.create_time,1 AS time_type, 
			L.vote_up-L.vote_down AS vote FROM tui_link_for_handle H 
			INNER JOIN topic_link TL ON H.insert_time<=? AND H.link_id=TL.link_id 
			INNER JOIN link L ON L.id=H.link_id 
			);`
		_, err = db.Query(sql, handleTime)
	}

	if err == nil {
		err = link_for_topic_vote_time(2, handleTime, db)
	}
	if err == nil {
		err = link_for_topic_vote_time(3, handleTime, db)
	}
	if err == nil {
		err = link_for_topic_vote_time(4, handleTime, db)
	}
	if err == nil {
		err = link_for_topic_vote_time(5, handleTime, db)
	}
	if err == nil {
		err = link_for_topic_vote_time(6, handleTime, db)
	}

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

func Del_link_for_topic_all(db *goku.MysqlDB) error {
	
	err := del_link_for_topic_later_top("tui_link_for_topic_top", "reddit_score DESC,link_id DESC", db)
	if err == nil {
		err = del_link_for_topic_later_top("tui_link_for_topic_later", "link_id DESC", db)
	}
	if err == nil {
		err = del_link_for_topic_hot_vote("tui_link_for_topic_hot", "vote_abs_score ASC,vote_add_score DESC,link_id DESC", db)
	}
	if err == nil {
		err = del_link_for_topic_hot_vote("tui_link_for_topic_vote", "vote DESC,link_id DESC", db)
	}

	return err
}

/**
 * 删除`tui_link_for_topic_top`最热, orderName:reddit_score DESC,link_id DESC
 * 删除`tui_link_for_topic_later`最新, orderName:link_id DESC
 */
func del_link_for_topic_later_top(tableName string, orderName string, db *goku.MysqlDB) error {
	
	sql := fmt.Sprintf(`INSERT ignore INTO tui_link_for_delete(id, time_type, del_count)
		SELECT topic_id, 0, tcount - %d FROM 
		(SELECT topic_id,COUNT(1) AS tcount FROM ` + tableName + ` GROUP BY topic_id) T
		WHERE T.tcount>%d;`, LinkMaxCount, LinkMaxCount)

	delSqlCreate := `INSERT ignore INTO tui_link_temporary_delete(id)
		( 
		SELECT link_id FROM ` + tableName + ` WHERE topic_id=%d ORDER BY ` + orderName + ` LIMIT %d,%d 
		);`
	delSqlDelete := `DELETE T FROM tui_link_temporary_delete D INNER JOIN ` + tableName + ` T ON T.topic_id=%d AND T.link_id=D.id; `
	
	iStart := 0
	var topicId int64
	var delCount int64
	
	db.Query("DELETE FROM tui_link_for_delete;")
	db.Query(sql)
	rows, err := db.Query("SELECT id, del_count FROM tui_link_for_delete LIMIT 0,?;", HandleCount)

	if err == nil {
		bWhile := rows.Next()
		bContinue := bWhile
		for bContinue && err == nil {
			for bWhile {
				rows.Scan(&topicId, &delCount)
				db.Query("DELETE FROM tui_link_temporary_delete;")
				db.Query(fmt.Sprintf(delSqlCreate, topicId, LinkMaxCount, delCount))
				db.Query(fmt.Sprintf(delSqlDelete, topicId))
				bWhile = rows.Next()
			}
			iStart += HandleCount
			rows, err = db.Query(fmt.Sprintf("SELECT id, del_count FROM tui_link_for_delete LIMIT %d,%d;", iStart, HandleCount)) 
			if err == nil {
				bWhile = rows.Next()
				bContinue = bWhile
			}
		}

	}

	return err
}

/** 
 * 删除`tui_link_for_topic_hot`热议, orderName:vote_abs_score ASC,vote_add_score DESC,link_id DESC
 * 删除`tui_link_for_topic_vote`投票, orderName:vote DESC,link_id DESC
 */
func del_link_for_topic_hot_vote(tableName string, orderName string, db *goku.MysqlDB) error {
	
	sql := fmt.Sprintf(`INSERT ignore INTO tui_link_for_delete(id, time_type, del_count)
		SELECT topic_id, time_type, tcount - %d FROM 
		(SELECT topic_id,time_type,COUNT(1) AS tcount FROM ` + tableName + ` GROUP BY topic_id,time_type) T
		WHERE T.tcount>%d;`, LinkMaxCount, LinkMaxCount)

	delSqlCreate := `INSERT ignore INTO tui_link_temporary_delete(id)
		( 
		SELECT link_id FROM ` + tableName + ` WHERE topic_id=%d AND time_type=%d ORDER BY ` + orderName + ` LIMIT %d,%d 
		); `
	delSqlDelete := `DELETE T FROM tui_link_temporary_delete D INNER JOIN ` + tableName + ` T ON T.topic_id=%d AND T.time_type=%d
		AND T.link_id=D.id; `

	
	iStart := 0
	var topicId int64
	var delCount int64
	var timeType int
	
	db.Query("DELETE FROM tui_link_for_delete;")
	db.Query(sql)
	rows, err := db.Query("SELECT id, time_type, del_count FROM tui_link_for_delete LIMIT 0,?;", HandleCount)

	if err == nil {
		bWhile := rows.Next()
		bContinue := bWhile
		for bContinue && err == nil {
			for bWhile {
				rows.Scan(&topicId, &timeType, &delCount)
				db.Query("DELETE FROM tui_link_temporary_delete;")
				db.Query(fmt.Sprintf(delSqlCreate, topicId, timeType, LinkMaxCount, delCount))
				db.Query(fmt.Sprintf(delSqlDelete, topicId, timeType))

				bWhile = rows.Next()
			}
			iStart += HandleCount
			rows, err = db.Query(fmt.Sprintf("SELECT id, time_type, del_count FROM tui_link_for_delete LIMIT %d,%d;", iStart, HandleCount)) 
			if err == nil {
				bWhile = rows.Next()
				bContinue = bWhile
			}
		}

	}

	return err
}

























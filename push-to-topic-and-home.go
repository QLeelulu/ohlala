package main

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink"
    "github.com/QLeelulu/ohlala/golink/models"
    "strconv"
    "strings"
    "time"
)

func main() {

	for {

		handleTime := time.Now()
		var db *goku.MysqlDB = GetDB()

		err := tui_link_for_topic(handleTime, db)
		if err == nil {
			err = tui_link_for_home(handleTime, db)
		}
		if err == nil {
			err = delete_tui_link_for_handle(handleTime, db)
		}
		if err != nil {
			goku.Logger().Errorln(err.Error())
		}

		db.Close()

		time.Sleep(300 * time.Second)
	}
}

/**
 * 推给话题
 */
func tui_link_for_topic(handleTime time.Time, db *goku.MysqlDB) error{

	err := models.link_for_topic_later(handleTime, db)
	if err == nil {
		err = models.link_for_topic_top(handleTime, db)
	}
	if err == nil {
		err = models.link_for_topic_hop_all(handleTime, db)
	}
	if err == nil {
		err = models.link_for_topic_vote_all(handleTime, db)
	}
	
	return err
}

/**
 * 推给话题
 */
func tui_link_for_home(handleTime time.Time, db *goku.MysqlDB) error{


	if err == nil {
		err = models.link_for_home_update(handleTime, db)
	}
	if err == nil {
		err = models.link_for_home_top(handleTime, db)
	}
	if err == nil {
		err = models.link_for_home_hot_all(handleTime, db)
	}
	if err == nil {
		err = models.link_for_home_vote_all(handleTime, db)
	}

	return err
}


/**
 * 删除tui_link_for_handle已经处理的数据
 */
func delete_tui_link_for_handle(handleTime time.Time, db *goku.MysqlDB) error {

	sql := `DELETE FROM tui_link_for_handle H WHERE H.insert_time<=?;`
	_, err := db.Query(sql, handleTime)

	return err
}




package main

import (
    "github.com/QLeelulu/goku"
    //"github.com/QLeelulu/ohlala/golink"
    "github.com/QLeelulu/ohlala/golink/models"
    //"strconv"
    //"strings"
	"fmt"
    "time"
)

func main() {
	
	delTime := time.Now()
	for {

		handleTime := time.Now()
		var db *goku.MysqlDB = models.GetDB()
//db.Debug = true

		err := tui_link_for_topic(handleTime, db)
		if err == nil {
			err = tui_link_for_home(handleTime, db)
		}
		if err == nil {
			err = tui_link_for_host(handleTime, db)
		}
		if err == nil {
			err = delete_tui_link_for_handle(handleTime, db)
		}

fmt.Println("tui wan")

		if handleTime.Sub(delTime).Seconds() >= 1800 && err == nil { // 每30分钟删除一次
			delTime = handleTime
			if err == nil {
fmt.Println("Del_link_for_home_all")
				err = models.Del_link_for_home_all(db)
			}
			if err == nil {
fmt.Println("Del_link_for_topic_all")
				err = models.Del_link_for_topic_all(db)
			}
			if err == nil {
fmt.Println("Del_link_for_host_all")
				err = models.Del_link_for_host_all(db)
			}
		}

		if err != nil {
			goku.Logger().Errorln(err.Error())
		}

		db.Close()

		time.Sleep(300 * time.Second) // 每5分钟推给话题/首页
	}
}

/**
 * 推给话题
 */
func tui_link_for_topic(handleTime time.Time, db *goku.MysqlDB) error{

	err := models.Link_for_topic_later(handleTime, db)
	if err == nil {
		err = models.Link_for_topic_top(handleTime, db)
	}
	if err == nil {
		err = models.Link_for_topic_hot_all(handleTime, db)
	}
	if err == nil {
		err = models.Link_for_topic_vote_all(handleTime, db)
	}
	
	return err
}

/**
 * 推给话题
 */
func tui_link_for_host(handleTime time.Time, db *goku.MysqlDB) error{

	err := models.Link_for_host_later(handleTime, db)
	if err == nil {
		err = models.Link_for_host_top(handleTime, db)
	}
	if err == nil {
		err = models.Link_for_host_hot_all(handleTime, db)
	}
	if err == nil {
		err = models.Link_for_host_vote_all(handleTime, db)
	}
	
	return err
}

/**
 * 推给首页
 */
func tui_link_for_home(handleTime time.Time, db *goku.MysqlDB) error{


	err := models.Link_for_home_update(handleTime, db)
	if err == nil {
		err = models.Link_for_home_top(handleTime, db)
	}
	if err == nil {
		err = models.Link_for_home_hot_all(handleTime, db)
	}
	if err == nil {
		err = models.Link_for_home_vote_all(handleTime, db)
	}

	return err
}


/**
 * 删除tui_link_for_handle已经处理的数据
 */
func delete_tui_link_for_handle(handleTime time.Time, db *goku.MysqlDB) error {

	sql := "DELETE FROM tui_link_for_handle WHERE `insert_time`<=? "
	_, err := db.Query(sql, handleTime)

	return err
}




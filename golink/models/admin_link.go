package models

import (
    //"fmt"
    "github.com/QLeelulu/goku"
    //"github.com/QLeelulu/ohlala/golink"
    //"github.com/QLeelulu/ohlala/golink/utils"
)

func Link_DelById(id int64) error {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    _, err := db.Query("UPDATE `link` SET status=2 WHERE id=?", id)
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return err
    }

    db.Query("DELETE FROM `host_link` WHERE link_id=?", id)
    db.Query("DELETE FROM `tui_link_for_host_later` WHERE link_id=?", id)
    db.Query("DELETE FROM `tui_link_for_host_top` WHERE link_id=?", id)
    db.Query("DELETE FROM `tui_link_for_host_hot` WHERE link_id=?", id)
    db.Query("DELETE FROM `tui_link_for_host_vote` WHERE link_id=?", id)
    db.Query("DELETE FROM `tui_link_for_topic_later` WHERE link_id=?", id)
    db.Query("DELETE FROM `tui_link_for_topic_top` WHERE link_id=?", id)
    db.Query("DELETE FROM `tui_link_for_topic_hot` WHERE link_id=?", id)
    db.Query("DELETE FROM `tui_link_for_topic_vote` WHERE link_id=?", id)
    db.Query("DELETE FROM `tui_link_for_home` WHERE link_id=?", id)
    db.Query("DELETE FROM `tui_link_for_handle` WHERE link_id=?", id)
    db.Query("DELETE FROM `link_for_user` WHERE link_id=?", id)
    db.Query("DELETE FROM `link_for_user_0` WHERE link_id=?", id)
    db.Query("DELETE FROM `link_for_user_1` WHERE link_id=?", id)
    db.Query("DELETE FROM `link_for_user_2` WHERE link_id=?", id)
    db.Query("DELETE FROM `link_for_user_3` WHERE link_id=?", id)
    db.Query("DELETE FROM `link_for_user_4` WHERE link_id=?", id)
    db.Query("DELETE FROM `link_for_user_5` WHERE link_id=?", id)
    db.Query("DELETE FROM `link_for_user_6` WHERE link_id=?", id)
    db.Query("DELETE FROM `link_for_user_7` WHERE link_id=?", id)
    db.Query("DELETE FROM `link_for_user_8` WHERE link_id=?", id)
    db.Query("DELETE FROM `link_for_user_9` WHERE link_id=?", id)
    db.Query("DELETE FROM `link_for_user_10` WHERE link_id=?", id)
    db.Query("DELETE FROM `link_for_user_11` WHERE link_id=?", id)
    db.Query("DELETE FROM `link_for_user_12` WHERE link_id=?", id)
    db.Query("DELETE FROM `link_for_user_13` WHERE link_id=?", id)
    db.Query("DELETE FROM `link_for_user_14` WHERE link_id=?", id)
    db.Query("DELETE FROM `link_for_user_15` WHERE link_id=?", id)
    db.Query("DELETE FROM `link_for_user_16` WHERE link_id=?", id)
    db.Query("DELETE FROM `link_for_user_17` WHERE link_id=?", id)
    db.Query("DELETE FROM `link_for_user_18` WHERE link_id=?", id)
    db.Query("DELETE FROM `link_for_user_19` WHERE link_id=?", id)
    db.Query("DELETE FROM `link_for_user_20` WHERE link_id=?", id)
    db.Query("DELETE FROM `link_for_user_21` WHERE link_id=?", id)
    db.Query("DELETE FROM `link_for_user_22` WHERE link_id=?", id)
    db.Query("DELETE FROM `link_for_user_23` WHERE link_id=?", id)
    db.Query("DELETE FROM `link_for_user_24` WHERE link_id=?", id)

    return nil
}

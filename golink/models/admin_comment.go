package models

import (
    //"fmt"
    "github.com/QLeelulu/goku"
    //"github.com/QLeelulu/ohlala/golink"
    //"github.com/QLeelulu/ohlala/golink/utils"
)


func Comment_DelById(id int64) (error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    _, err := db.Query("UPDATE `comment` SET status=2 WHERE id=?", id)
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return err
    }

    return nil
}

package models

import (
    // // "bytes"
    // "errors"
    // "fmt"
    "github.com/QLeelulu/goku"

// "github.com/QLeelulu/goku/form"
// "github.com/QLeelulu/ohlala/golink"
// "github.com/QLeelulu/ohlala/golink/utils"
// // "html/template"
// // "strings"
// "time"
)

var table_CommentForUser string = "comment_for_user"

func CommentForUser_Add(userId int64, comment Comment) error {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    m := map[string]interface{}{}

    m["user_id"] = userId
    m["comment_id"] = comment.Id
    m["link_id"] = comment.LinkId
    m["pcomment_id"] = comment.ParentId
    m["create_time"] = comment.CreateTime

    _, err := db.Insert(table_CommentForUser, m)
    if err != nil {
        goku.Logger().Errorln(err.Error())
    }
    return err
}

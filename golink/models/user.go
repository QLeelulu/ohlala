package models

import (
    "database/sql"
    "github.com/QLeelulu/goku"
    // "github.com/QLeelulu/ohlala/golink/utils"
    "strings"
    "time"
)

type User struct {
    Id                   int
    Name                 string
    Email                string
    Pwd                  string
    UserPic              string
    Description          string
    ReferenceSystem      int
    ReferenceToken       string
    ReferenceTokenSecret string
    CreateAt             time.Time
}

// 检查email地址是否存在。
// 任何出错都认为email地址存在，防止注册
func User_IsEmailExist(email string) bool {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    res, err := db.Exec("select id from `user` where `email`=? limit 1", strings.ToLower(email))
    if err != nil {
        goku.Logger().Errorln(err.Error())
        // 出错直接认为email存在
        return true
    }
    var af int64
    af, err = res.RowsAffected()
    if err != nil {
        goku.Logger().Errorln(err.Error())
        // 出错直接认为email存在
        return true
    }
    if af < 1 {
        return false
    }
    return true
}

func User_SaveMap(m map[string]interface{}) (sql.Result, error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()
    r, err := db.Insert("user", m)
    return r, err
}

func User_Update(id int, m map[string]interface{}) (sql.Result, error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()
    r, err := db.Update("user", m, "id=?", id)
    return r, err
}

func User_Delete(id int) (sql.Result, error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()
    r, err := db.Delete("user", "id=?", id)
    return r, err
}

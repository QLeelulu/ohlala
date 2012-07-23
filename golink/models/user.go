package models

import (
    "database/sql"
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink/utils"
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

    rows, err := db.Query("select id from `user` where `email_lower`=? limit 1", strings.ToLower(email))
    if err != nil {
        goku.Logger().Errorln(err.Error())
        // 出错直接认为email存在
        return true
    }
    defer rows.Close()
    if rows.Next() {
        return true
    }
    return false
}

// 检查账号密码是否正确
// 如果正确，则返回用户id
func User_CheckPwd(email, pwd string) int {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    pwd = utils.PasswordHash(pwd)
    rows, err := db.Query("select id from `user` where `email_lower`=? and pwd=? limit 1", strings.ToLower(email), pwd)
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return 0
    }
    defer rows.Close()
    if rows.Next() {
        var id int
        err = rows.Scan(&id)
        if err != nil {
            goku.Logger().Errorln(err.Error())
        } else {
            return id
        }
    }
    return 0
}

func User_SaveMap(m map[string]interface{}) (sql.Result, error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()
    m["email_lower"] = strings.ToLower(m["email"].(string))
    r, err := db.Insert("user", m)
    return r, err
}

func User_GetByTicket(ticket string) (*User, error) {
    redisClient := GetRedis()
    defer redisClient.Quit()

    id, err := redisClient.Get(ticket)
    if err != nil {
        return nil, err
    }

    if id.String() == "" {
        return nil, nil
    }

    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    var user *User = new(User)
    err = db.GetStruct(user, "id=?", id.String())
    if err != nil {
        return nil, err
    }
    if user.Id > 0 {
        return user, nil
    }
    return nil, nil
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

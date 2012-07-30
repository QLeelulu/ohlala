package models

import (
    "github.com/QLeelulu/goku"
)

/**
 * 用户的粉丝与好友
 */

// 获取用户关注的好友列表
func UserFollow_Friends(userId int64, page, pagesize int) ([]User, error) {
    if page < 1 {
        page = 1
    }
    page = page - 1
    if pagesize == 0 {
        pagesize = 20
    }
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    qi := goku.SqlQueryInfo{}
    qi.Fields = "u.id, u.name, u.email, u.user_pic"
    qi.Join = " uf INNER JOIN `user` u ON uf.follow_id=u.id"
    qi.Where = "uf.user_id=?"
    qi.Params = []interface{}{userId}
    qi.Limit = pagesize
    qi.Offset = pagesize * page
    qi.Order = "u.id desc"

    rows, err := db.Select("user_follow", qi)

    if err != nil {
        goku.Logger().Errorln(err.Error())
        return nil, err
    }
    defer rows.Close()

    users := make([]User, 0)
    for rows.Next() {
        user := User{}
        err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.UserPic)
        if err != nil {
            goku.Logger().Errorln(err.Error())
            return nil, err
        }
        users = append(users, user)
    }
    return users, nil
}

// 获取用户关注的粉丝列表
func UserFollow_Followers(userId int64, page, pagesize int) ([]User, error) {
    if page < 1 {
        page = 1
    }
    page = page - 1
    if pagesize == 0 {
        pagesize = 20
    }
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    qi := goku.SqlQueryInfo{}
    qi.Fields = "u.id, u.name, u.email, u.user_pic"
    qi.Join = " uf INNER JOIN `user` u ON uf.user_id=u.id"
    qi.Where = "uf.follow_id=?"
    qi.Params = []interface{}{userId}
    qi.Limit = pagesize
    qi.Offset = pagesize * page
    qi.Order = "u.id desc"

    rows, err := db.Select("user_follow", qi)

    if err != nil {
        goku.Logger().Errorln(err.Error())
        return nil, err
    }
    defer rows.Close()

    users := make([]User, 0)
    for rows.Next() {
        user := User{}
        err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.UserPic)
        if err != nil {
            goku.Logger().Errorln(err.Error())
            return nil, err
        }
        users = append(users, user)
    }
    return users, nil
}

package models

import (
    "fmt"

    "github.com/QLeelulu/goku"
)

// 获取用户的未读链接数

// 全部链接的最新链接的未读数
func NewestLinkUnread_All(userId, lastReadLinkId int64) (int64, error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    qi := goku.SqlQueryInfo{}
    qi.Fields = "max(id)"
    rows, err := db.Select(Table_Link, qi)
    var maxLinkId int64
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return 0, err
    }
    if rows.Next() {
        err = rows.Scan(&maxLinkId)
        if err != nil {
            goku.Logger().Errorln(err.Error())
            return 0, err
        }
    }
    return maxLinkId - lastReadLinkId - 1, nil
}

// 关注好友的最新链接的未读数
func NewestLinkUnread_Friends(userId, lastReadLinkId int64) (int64, error) {
    if userId < 1 {
        return 0, nil
    }

    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    qi := goku.SqlQueryInfo{}
    qi.Where = "`user_id`=? and `link_id`>?"
    qi.Params = []interface{}{userId, lastReadLinkId}
    qi.Fields = "count(*)"
    tableName := LinkForUser_TableName(userId)
    rows, err := db.Select(tableName, qi)
    var unreadCount int64
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return 0, err
    }
    if rows.Next() {
        err = rows.Scan(&unreadCount)
        if err != nil {
            goku.Logger().Errorln(err.Error())
            return 0, err
        }
    }
    return unreadCount, nil
}

// 更新用户已读的最新链接的最大的链接id
func NewestLinkUnread_UpdateForAll(userId, lastReadLinkId int64) error {
    if userId < 1 || lastReadLinkId < 1 {
        return nil
    }

    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    m := map[string]interface{}{
        "last_read_link_id": lastReadLinkId,
    }
    _, err := db.Update(Table_User, m, "id=?", userId)
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return err
    }
    return nil
}

// 更新用户已读的关注好友的最新链接的最大的链接id
func NewestLinkUnread_UpdateForUser(userId, lastReadLinkId int64) error {
    if userId < 1 || lastReadLinkId < 1 {
        return nil
    }

    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    m := map[string]interface{}{
        "last_read_friend_link_id": lastReadLinkId,
    }
    _, err := db.Update(Table_User, m, "id=?", userId)
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return err
    }
    return nil
}

// @userId: 用户id，如果没有用户，则传0.
func NewestLinkUnread_ToString(userId, unreadCount int64) string {
    if unreadCount < 1 {
        return ""
    } else if unreadCount > 99 {
        return "99+"
    }
    return fmt.Sprintf("%d", unreadCount)
}

package models

import (
    // _ "code.google.com/p/go-mysql-driver/mysql"
    "database/sql"
    "fmt"
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink"
    "github.com/simonz05/godis"
    _ "github.com/ziutek/mymysql/godrv"
)

const (
    Table_Link    = "link"
    Table_Comment = "comment"
)

func GetDB() *goku.MysqlDB {
    db, err := goku.OpenMysql(golink.DATABASE_Driver, golink.DATABASE_DSN)
    if err != nil {
        panic(err.Error())
    }
    return db
}

func GetRedis() *godis.Client {
    return godis.New(golink.REDIS_HOST, 0, golink.REDIS_AUTH)
}

// 加（减）表里面的统计数据
// @table: 要操作数据库表名
// @field: 要修改的字段
// @inc: 要增加或减少的值
func IncCountById(db *goku.MysqlDB, table string, id int64, field string, inc int) (sql.Result, error) {
    r, err := db.Exec(fmt.Sprintf("UPDATE `%s` SET %s=%s+? WHERE id=?;", table, field, field), inc, id)
    if err != nil {
        goku.Logger().Errorln(err.Error())
    }
    return r, err
}

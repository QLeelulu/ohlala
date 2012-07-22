package models

import (
    // _ "code.google.com/p/go-mysql-driver/mysql"
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink"
    _ "github.com/ziutek/mymysql/godrv"
)

func GetDB() *goku.MysqlDB {
    db, err := goku.OpenMysql(golink.DATABASE_Driver, golink.DATABASE_DSN)
    if err != nil {
        panic(err.Error())
    }
    return db
}

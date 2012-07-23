package models

import (
    // _ "code.google.com/p/go-mysql-driver/mysql"
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink"
    "github.com/simonz05/godis"
    _ "github.com/ziutek/mymysql/godrv"
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

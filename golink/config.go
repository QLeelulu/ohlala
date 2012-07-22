package golink

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/mustache.goku"
    "path"
    "runtime"
    "time"
)

var (
    DATABASE_Driver string = "mymysql"
    // mysql: "user:password@/dbname?charset=utf8&keepalive=1"
    // mymysql: tcp:localhost:3306*test_db/lulu/123456
    DATABASE_DSN string = "tcp:localhost:3306*link/root/112358"
)

var Config *goku.ServerConfig = &goku.ServerConfig{
    Addr:           ":8080",
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 << 20,
    //RootDir:        _, filename, _, _ := runtime.Caller(1),
    StaticPath: "static", // static content 
    ViewPath:   "views",
    Debug:      true,
}

func init() {
    // WTF, i just want to set the RootDir as current dir.
    _, filename, _, _ := runtime.Caller(1)
    Config.RootDir = path.Dir(filename)

    // template engine
    te := mustache.NewMustacheTemplateEngine()
    te.UseCache = !Config.Debug
    Config.TemplateEnginer = te

    goku.SetGlobalViewData("SiteName", "Todo - by {goku}")
}

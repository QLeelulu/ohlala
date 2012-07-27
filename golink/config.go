package golink

import (
    "github.com/QLeelulu/goku"
    // "github.com/QLeelulu/mustache.goku"
    "log"
    "os"
    "path"
    "runtime"
    "time"
)

var (
    // mysql
    DATABASE_Driver string = "mymysql"
    DATABASE_DSN    string = "tcp:localhost:3306*link/root/112358"

    // reddit time
    SITERUNTIME      string    = "2012-07-21 23:20:25"
    SITERUNTIME_TIME time.Time = time.Date(2012, time.July, 21, 23, 20, 25, 0, time.UTC)

    // redis
    REDIS_HOST string = "tcp:127.0.0.1:6379"
    REDIS_AUTH string = ""

    // errors
    ERROR_DATABASE = "数据库出错"
)

var Config *goku.ServerConfig = &goku.ServerConfig{
    Addr:           ":8080",
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 << 20,

    //RootDir:        _, filename, _, _ := runtime.Caller(1),
    StaticPath: "static", // static content 
    ViewPath:   "views",

    Logger:   log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
    LogLevel: goku.LOG_LEVEL_LOG,

    Debug: true,
}

func init() {
    // WTF, i just want to set the RootDir as current dir.
    _, filename, _, _ := runtime.Caller(1)
    Config.RootDir = path.Dir(filename)

    // // template engine
    // te := mustache.NewMustacheTemplateEngine()
    // te.UseCache = !Config.Debug
    // Config.TemplateEnginer = te

    goku.SetGlobalViewData("SiteName", "Todo - by {goku}")
}

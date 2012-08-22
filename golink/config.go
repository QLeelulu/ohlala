package golink

import (
    "flag"
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/goku/utils"
    // "github.com/QLeelulu/mustache.goku"
    "log"
    "os"
    "path"
    "runtime"
    "time"
    //"math"
)

var (
    // mysql
    DATABASE_Driver string = "mymysql"
    DATABASE_DSN    string = "tcp:localhost:3306*link/root/112358"

    // reddit time
    SITERUNTIME      string    = "2012-07-21 23:20:25"
    SITERUNTIME_TIME time.Time = time.Date(2012, time.July, 21, 23, 20, 25, 0, time.UTC)
	SCORETIMESTEMP   float64   = 45000.0

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

    loadFileConf()
}

const (
    PATH_IMAGE_UPLOAD = "/static/img/"
    PATH_IMAGE_AVATAR = PATH_IMAGE_UPLOAD + "avatar/"
    PATH_USER_AVATAR  = PATH_IMAGE_AVATAR + "user/"
    PATH_TOPIC_AVATAR = PATH_IMAGE_AVATAR + "topic/"

    // 队列KYE
    KEY_LIST_PUSH_TO_USER = "link_for_user"
)

func loadFileConf() {
    var confFile string
    // flag.StringVar(&confFile, "golink-conf", "", "golink的配置文件路径，json格式")
    // flag.Parse()

    // cmd := os.Args[1]
    flags := flag.NewFlagSet("golink-conf", flag.ContinueOnError)
    flags.StringVar(&confFile, "conf", "", "golink的配置文件路径，json格式")
    flags.Parse(os.Args[1:])

    if confFile != "" {
        conf, err := utils.LoadJsonFile(confFile)
        if err != nil {
            log.Fatalln("load conf file", confFile, "error:", err.Error())
        }
        if fc, ok := conf["DataBase"]; ok {
            dbc, ok := fc.(map[string]interface{})
            if !ok {
                log.Fatalln("conf file error: wrong DataBase Session format.")
            }
            DATABASE_Driver = dbc["DATABASE_Driver"].(string)
            DATABASE_DSN = dbc["DATABASE_DSN"].(string)
            REDIS_HOST = dbc["REDIS_HOST"].(string)
            REDIS_AUTH = dbc["REDIS_AUTH"].(string)
        }
    }
}

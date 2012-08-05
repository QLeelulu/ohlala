package main

import (
    "flag"
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/goku/utils"
    "github.com/QLeelulu/ohlala/golink"
    _ "github.com/QLeelulu/ohlala/golink/controllers" // notice this!! import controllers
    "github.com/QLeelulu/ohlala/golink/middlewares"
    "log"
    "os"
)

func main() {

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
            golink.DATABASE_Driver = dbc["DATABASE_Driver"].(string)
            golink.DATABASE_DSN = dbc["DATABASE_DSN"].(string)
            golink.REDIS_HOST = dbc["REDIS_HOST"].(string)
            golink.REDIS_AUTH = dbc["REDIS_AUTH"].(string)
        }
    }

    rt := &goku.RouteTable{Routes: golink.Routes}
    middlewares := []goku.Middlewarer{
        new(middlewares.UtilMiddleware),
    }
    s := goku.CreateServer(rt, middlewares, golink.Config)
    goku.Logger().Logln("Server start on", s.Addr)

    log.Fatal(s.ListenAndServe())
}

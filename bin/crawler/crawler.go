package main

import (
    "log"
    "path"
    "runtime"
    "time"

    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink"
    "github.com/QLeelulu/ohlala/golink/crawler"
    "github.com/QLeelulu/ohlala/golink/utils"
)

type RssConfItem struct {
    Url  string
    Name string
}

type RssConf struct {
    Items []RssConfItem
}

func main() {
    uids := []int64{10011, 10012, 10013}

    _, filename, _, _ := runtime.Caller(0)
    rssConfFile := path.Join(path.Dir(filename), "rss_urls.json")
    rssConf := RssConf{}
    err := utils.LoadJsonFile(rssConfFile, &rssConf.Items)
    if err != nil {
        log.Fatalln("load conf file", rssConfFile, "error:", err.Error())
    }

    // 这个只是为了设置Goku的Log级别，
    // 后面需要重构Goku。
    rt := &goku.RouteTable{Routes: golink.Routes}
    goku.CreateServer(rt, nil, golink.Config)

    for {
        for _, rssItem := range rssConf.Items {
            rssc := crawler.RssCrawler{}
            rssc.Name = rssItem.Name
            rssc.Url = rssItem.Url
            rssc.UserIds = uids
            rssc.Run()
        }
        time.Sleep(time.Minute * 10)
    }
}

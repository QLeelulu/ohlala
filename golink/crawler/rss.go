package crawler

import (
    "encoding/xml"
    "errors"
    "io/ioutil"
    "net"
    "net/http"
    "strings"
    "time"

    "github.com/QLeelulu/goku"
)

var timeout = time.Duration(30 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
    return net.DialTimeout(network, addr, timeout)
}

type RSS struct {
    // XMLName xml.Name `xml:"rss"`
    Channel RssChannel `xml:"channel"`
}

type RssChannel struct {
    Title       string `xml:"title"`
    Link        string `xml:"link"`
    Description string `xml:"description"`
    Items       []Item `xml:"item"`
}
type Item struct {
    Title       string `xml:"title"`
    Link        string `xml:"link"`
    Description string `xml:"description"`
}

// rss爬虫
type RssCrawler struct {
    BaseCrawler

    rss *RSS
}

func (self *RssCrawler) Run() (err error) {
    if err = self.getContent(); err != nil {
        goku.Logger().Errorln("read rss content error:", err.Error())
        return err
    }
    if self.rss == nil {
        err = errors.New("no rss content")
        goku.Logger().Errorln(err.Error())
        return err
    }
    items := self.rss.Channel.Items
    successCount := 0
    submitedCount := 0
    for i := len(items) - 1; i >= 0; i-- {
        item := items[i]
        err = self.saveLink(item.Link, item.Title)
        if err == nil {
            successCount++
        } else if strings.Index(err.Error(), "Url已经提交过") > -1 {
            submitedCount++
            if submitedCount > 4 {
                break
            }
        }
    }
    goku.Logger().Noticef("%s(%s) import %d.", self.Name, self.Url, successCount)
    return nil
}

func (self *RssCrawler) getContent() error {
    transport := http.Transport{
        Dial: dialTimeout,
    }
    client := &http.Client{
        Transport: &transport,
    }
    req, err := http.NewRequest("GET", self.Url, nil)
    if err != nil {
        return err
    }
    req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.69 Safari/537.36")

    res, err := client.Do(req)
    if err != nil {
        return err
    }
    asText, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return err
    }

    var i RSS
    err = xml.Unmarshal([]byte(asText), &i)
    if err != nil {
        return err
    }
    self.rss = &i
    return nil
}

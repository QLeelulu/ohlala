package crawler

import (
    "strings"
    // "time"
    // "fmt"

    "github.com/QLeelulu/goku"
    rss "github.com/jteeuwen/go-pkg-rss"
)

// rss爬虫
type RssCrawler struct {
    BaseCrawler
}

func (self *RssCrawler) Run() (err error) {
    err = self.PollFeed(30)
    return
}

func (self *RssCrawler) PollFeed(timeout int) error {
    feed := rss.New(timeout, true, self.chanHandler, self.itemHandler)

    // for {
    if err := feed.Fetch(self.Url, nil); err != nil {
        goku.Logger().Errorf("%s: %s", self.Url, err)
        return err
    }

    // <-time.After(time.Duration(feed.SecondsTillUpdate() * 1e9))
    // }
    return nil
}

func (self *RssCrawler) chanHandler(feed *rss.Feed, newchannels []*rss.Channel) {
    // fmt.Printf("%d new channel(s) in %s\n", len(newchannels), feed.Url)
}

func (self *RssCrawler) itemHandler(feed *rss.Feed, ch *rss.Channel, items []*rss.Item) {
    // fmt.Printf("%d new item(s) in %s\n", len(items), feed.Url)
    successCount := 0
    submitedCount := 0
    for i, l := 0, len(items); i < l; i++ {
        item := items[i]
        err := self.saveLink(item.Links[0].Href, item.Title)
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
}

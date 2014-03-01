package crawler

import (
    "github.com/sdegutis/go.assert"

    "testing"
)

func TestRssCrawler(t *testing.T) {
    rssc := RssCrawler{}
    rssc.Name = "网易头条"
    rssc.Url = "http://news.163.com/special/00011K6L/rss_newstop.xml"
    rssc.UserIds = []int64{10000, 10001, 10002, 10003}

    err := rssc.Run()
    assert.Equals(t, err, nil)

    items := rssc.rss.Channel.Items
    assert.NotEquals(t, len(items), 0)
}

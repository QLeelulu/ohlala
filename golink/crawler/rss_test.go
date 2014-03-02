package crawler

import (
    "github.com/sdegutis/go.assert"

    "testing"
)

func TestRssCrawler(t *testing.T) {
    rssc := RssCrawler{}
    rssc.Name = "网易头条"
    rssc.Url = "http://www.ruanyifeng.com/blog/atom.xml"
    rssc.UserIds = []int64{10000, 10001, 10002, 10003}

    err := rssc.Run()
    assert.Equals(t, err, nil)
}

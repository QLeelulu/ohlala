package utils

import (
    "github.com/sdegutis/go.assert"
    "strings"
    "testing"
    "time"
)

func TestSmcTimeSince(t *testing.T) {
    // (?P<name>re)
    now := time.Now()
    s := SmcTimeSince(now)
    assert.Equals(t, s, "刚刚")

    t2 := now.Add(-30 * time.Minute)
    s = SmcTimeSince(t2)
    assert.Equals(t, s, "30分钟前")

    t2 = now.Add(-10 * time.Hour)
    s = SmcTimeSince(t2)
    assert.Equals(t, s, "10小时前")

    t2 = now.Add(-24 * time.Hour)
    s = SmcTimeSince(t2)
    assert.Equals(t, strings.Contains(s, "昨天"), true)

    t2 = now.Add(-48 * time.Hour)
    s = SmcTimeSince(t2)
    assert.Equals(t, strings.Contains(s, "前天"), true)

    t2 = now.Add(-24 * 30 * 12 * time.Hour)
    s = SmcTimeSince(t2)
    assert.Equals(t, s, t2.Format("2006年1月2日 15:04"))
}

func TestMd5(t *testing.T) {
    assert.Equals(t, MD5("test-md5"), "42505649011f560d32a0a7a6427a8778")
    assert.Equals(t, MD5("ohlala-md5"), "1d8fd7200f4f9c0111d2e7fcc0a564e6")
}

package utils

import (
    "github.com/sdegutis/go.assert"
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
    assert.Equals(t, s, "刚刚")

    t2 = now.Add(-10 * time.Hour)
    s = SmcTimeSince(t2)
    assert.Equals(t, s, "刚刚")

    t2 = now.Add(-24 * time.Hour)
    s = SmcTimeSince(t2)
    assert.Equals(t, s, "刚刚")

    t2 = now.Add(-48 * time.Hour)
    s = SmcTimeSince(t2)
    assert.Equals(t, s, "刚刚")

    t2 = now.Add(-24 * 30 * 12 * time.Hour)
    s = SmcTimeSince(t2)
    assert.Equals(t, s, "刚刚")
}

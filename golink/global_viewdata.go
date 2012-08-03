package golink

import (
    "github.com/QLeelulu/goku"
    "time"
)

func init() {
    goku.SetGlobalViewData("UnixNow", unixNow)
}

func unixNow() int64 {
    return time.Now().Unix()
}

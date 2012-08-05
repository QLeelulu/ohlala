package golink

import (
    "github.com/QLeelulu/goku"
    // "html/template"
    "time"
)

func init() {
    goku.SetGlobalViewData("UnixNow", unixNow)
    goku.SetGlobalViewData("ilg", funcMap_ilg)
}

func unixNow() int64 {
    return time.Now().Unix()
}

// // First we create a FuncMap with which to register the function.
// var funcMap template.FuncMap = template.FuncMap{
//     "lg":  funcMap_lg,
//     "lge": funcMap_lge,
// }

func funcMap_ilg(a, b int) bool {
    return a > b
}

func funcMap_ilge(a, b int) bool {
    return a >= b
}

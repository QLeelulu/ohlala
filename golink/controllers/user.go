package controllers

import (
    "github.com/QLeelulu/goku"
)

var _ = goku.Controller("user").
    Get("login", func(ctx *goku.HttpContext) goku.ActionResulter {
    return ctx.View(nil)
}).
    Post("login", func(ctx *goku.HttpContext) goku.ActionResulter {
    return ctx.View(nil)
}).
    Post("reg", func(ctx *goku.HttpContext) goku.ActionResulter {
    return ctx.View(nil)
})

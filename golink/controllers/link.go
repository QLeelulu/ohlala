package controllers

import (
    "github.com/QLeelulu/goku"
)

var _ = goku.Controller("link").
    // 
    Get("index", func(ctx *goku.HttpContext) goku.ActionResulter {
    return ctx.View(nil)
}).
    Get("show", func(ctx *goku.HttpContext) goku.ActionResulter {
    return ctx.View(nil)
}).
    // 提交链接表单
    Get("submit", func(ctx *goku.HttpContext) goku.ActionResulter {
    return ctx.View(nil)
}).
    // 提交一个链接并保存到数据库
    Post("submit", func(ctx *goku.HttpContext) goku.ActionResulter {
    return ctx.View(nil)
}).
    // 添加评论
    Post("comment", func(ctx *goku.HttpContext) goku.ActionResulter {
    return ctx.View(nil)
})

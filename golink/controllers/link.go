package controllers

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink/filters"
    "github.com/QLeelulu/ohlala/golink/forms"
    "github.com/QLeelulu/ohlala/golink/models"
)

var _ = goku.Controller("link").
    // 
    Get("index", func(ctx *goku.HttpContext) goku.ActionResulter {

    return ctx.View(nil)
}).
    /**
     * 查看一个链接的评论
     */
    Get("show", func(ctx *goku.HttpContext) goku.ActionResulter {
    return ctx.View(nil)
}).

    /**
     * 提交链接的表单
     */
    Get("submit", func(ctx *goku.HttpContext) goku.ActionResulter {

    ctx.ViewData["Values"] = map[string]string{
        "title":   ctx.Get("title"),
        "context": ctx.Get("url"),
    }
    return ctx.View(nil)

}).Filters(filters.NewRequireLoginFilter()).

    /**
     * 提交一个链接并保存到数据库
     */
    Post("submit", func(ctx *goku.HttpContext) goku.ActionResulter {

    f := forms.CreateLinkSubmitForm()
    f.FillByRequest(ctx.Request)

    success, errorMsgs := models.Link_SaveForm(f, (ctx.Data["user"].(*models.User)).Id)

    if success {
        return ctx.Redirect("/")
    } else {
        ctx.ViewData["Errors"] = errorMsgs
        ctx.ViewData["Values"] = f.Values()
    }
    return ctx.View(nil)

}).Filters(filters.NewRequireLoginFilter()).

    /**
     * 添加评论
     */
    Post("comment", func(ctx *goku.HttpContext) goku.ActionResulter {

    return ctx.View(nil)

}).Filters(filters.NewRequireLoginFilter())

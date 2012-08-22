package controllers

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink/models"
    "strconv"
)

var _ = goku.Controller("discover").

    /**
     * 未登陆用户首页
     */
    Get("index", func(ctx *goku.HttpContext) goku.ActionResulter {

    ot := ctx.Get("o")
    if ot == "" {
        ot = "hot"
    }
    dt, _ := strconv.Atoi(ctx.Get("dt"))
    ctx.ViewData["Order"] = ot
    links, _ := models.LinkForHome_GetByPage(ot, dt, 1, 20)
    ctx.ViewData["Links"] = models.Link_ToVLink(links, ctx)
    ctx.ViewData["TopTab"] = "discover"
    return ctx.Render("/home/index", nil)
})

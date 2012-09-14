package controllers

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink"
    "github.com/QLeelulu/ohlala/golink/filters"
    "github.com/QLeelulu/ohlala/golink/models"
    "strconv"
)

var _ = goku.Controller("discover").

    /**
     * 未登陆用户首页
     */
    Get("index", discover_index).

    /**
     * 未登陆用户首页
     */
    Get("loadmorelink", discover_loadMoreLink).
    Filters(filters.NewAjaxFilter())

// END Controller & Action
// 

// 发现 首页
func discover_index(ctx *goku.HttpContext) goku.ActionResulter {
    ot := ctx.Get("o")
    if ot == "" {
        ot = "hot"
    }
    dt, _ := strconv.Atoi(ctx.Get("dt"))
    ctx.ViewData["Order"] = ot
    links, _ := models.LinkForHome_GetByPage(ot, dt, 1, golink.PAGE_SIZE)
    ctx.ViewData["Links"] = models.Link_ToVLink(links, ctx)
    ctx.ViewData["TopTab"] = "discover"
    ctx.ViewData["HasMoreLink"] = len(links) >= golink.PAGE_SIZE
    return ctx.Render("/home/index", nil)
}

// 加载更多link
func discover_loadMoreLink(ctx *goku.HttpContext) goku.ActionResulter {
    page, err := strconv.Atoi(ctx.Get("page"))
    success, hasmore := false, false
    errorMsgs, html := "", ""
    if err == nil && page > 1 {
        ot := ctx.Get("o")
        if ot == "" {
            ot = "hot"
        }
        dt, _ := strconv.Atoi(ctx.Get("dt"))
        links, _ := models.LinkForHome_GetByPage(ot, dt, page, golink.PAGE_SIZE)
        if links != nil && len(links) > 0 {
            ctx.ViewData["Links"] = models.Link_ToVLink(links, ctx)
            vr := ctx.RenderPartial("loadmorelink", nil)
            vr.Render(ctx, vr.Body)
            html = vr.Body.String()
            hasmore = len(links) >= golink.PAGE_SIZE
        }
        success = true
    } else {
        errorMsgs = "参数错误"
    }
    r := map[string]interface{}{
        "success": success,
        "errors":  errorMsgs,
        "html":    html,
        "hasmore": hasmore,
    }
    return ctx.Json(r)
}

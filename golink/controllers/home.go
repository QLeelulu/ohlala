package controllers

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink"
    "github.com/QLeelulu/ohlala/golink/filters"
    "github.com/QLeelulu/ohlala/golink/models"
    "strconv"
)

var _ = goku.Controller("home").
    // index
    Get("index", home_index).
    // 关于
    Get("about", home_about).
    // 免责声明
    Get("disclaimer", home_disclaimer).
    // 合作伙伴
    Get("partner", home_partner).
    // 联系信息
    Get("contact", home_contact).
    // load more
    Get("loadmorelink", home_loadMoreLink).
    Filters(filters.NewRequireLoginFilter(), filters.NewAjaxFilter())

//

func home_index(ctx *goku.HttpContext) goku.ActionResulter {
    u, ok := ctx.Data["user"]
    if !ok || u == nil {
        return ctx.Redirect("/discover")
    }
    user := u.(*models.User)
    ot := ctx.Get("o")
    if ot == "" {
        ot = "hot"
    }
    ctx.ViewData["Order"] = ot
    links, _ := models.Link_ForUser(user.Id, ot, 1, golink.PAGE_SIZE) //models.Link_GetByPage(1, 20)
    ctx.ViewData["Links"] = models.Link_ToVLink(links, ctx)
    ctx.ViewData["HasMoreLink"] = len(links) >= golink.PAGE_SIZE
    return ctx.View(nil)
}

func home_loadMoreLink(ctx *goku.HttpContext) goku.ActionResulter {
    page, err := strconv.Atoi(ctx.Get("page"))
    success, hasmore := false, false
    errorMsgs, html := "", ""
    if err == nil && page > 1 {
        user := ctx.Data["user"].(*models.User)
        ot := ctx.Get("o")
        if ot == "" {
            ot = "hot"
        }
        links, _ := models.Link_ForUser(user.Id, ot, page, golink.PAGE_SIZE)
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

func home_about(ctx *goku.HttpContext) goku.ActionResulter {
    return ctx.View(nil)
}

func home_disclaimer(ctx *goku.HttpContext) goku.ActionResulter {
    return ctx.View(nil)
}

func home_contact(ctx *goku.HttpContext) goku.ActionResulter {
    return ctx.View(nil)
}

func home_partner(ctx *goku.HttpContext) goku.ActionResulter {
    return ctx.View(nil)
}

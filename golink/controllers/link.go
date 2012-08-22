package controllers

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink/filters"
    "github.com/QLeelulu/ohlala/golink/forms"
    "github.com/QLeelulu/ohlala/golink/models"
    "strconv"
    "strings"
)

var _ = goku.Controller("link").
    // 
    Get("indexxxxx", func(ctx *goku.HttpContext) goku.ActionResulter {

    return ctx.View(nil)
}).
    /**
     * 查看一个链接的评论
     */
    Get("show", func(ctx *goku.HttpContext) goku.ActionResulter {

    linkId, _ := strconv.ParseInt(ctx.RouteData.Params["id"], 10, 64)
    link, _ := models.Link_GetById(linkId)

    if link == nil {
        ctx.ViewData["errorMsg"] = "内容不存在"
        return ctx.Render("error", nil)
    }

    vlink := models.Link_ToVLink([]models.Link{*link}, ctx)

    comments := models.Comment_SortForLink(link.Id, "hot")

    ctx.ViewData["Comments"] = comments
    return ctx.View(vlink[0])
}).

    /**
     * 提交链接的表单页面
     */
    Get("submit", func(ctx *goku.HttpContext) goku.ActionResulter {

    ctx.ViewData["Values"] = map[string]string{
        "title":   ctx.Get("t"),
        "context": ctx.Get("u"),
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
     * 提交评论并保存到数据库
     */
    Post("ajax-comment", func(ctx *goku.HttpContext) goku.ActionResulter {

    f := forms.NewCommentSubmitForm()
    f.FillByRequest(ctx.Request)

    var success bool
    var errorMsgs string
    if ctx.RouteData.Params["id"] != f.Values()["link_id"] {
        errorMsgs = "参数错误"
    } else {
        var errors []string
        success, errors = models.Comment_SaveForm(f, (ctx.Data["user"].(*models.User)).Id)
        if errors != nil {
            errorMsgs = strings.Join(errors, "\n")
        }
    }
    r := map[string]interface{}{
        "success": success,
        "errors":  errorMsgs,
    }
    return ctx.Json(r)

}).Filters(filters.NewRequireLoginFilter(), filters.NewAjaxFilter())

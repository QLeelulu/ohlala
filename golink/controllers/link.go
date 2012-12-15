package controllers

import (
    "fmt"
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink/filters"
    "github.com/QLeelulu/ohlala/golink/forms"
    "github.com/QLeelulu/ohlala/golink/models"
    "html/template"
    "strconv"
    "strings"
)

var _ = goku.Controller("link").
    /**
     * 查看某评论
     */
    Get("permacoment", link_permacoment).
    /**
     * 查看一个链接的评论
     */
    Get("show", link_show).

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

}).Filters(filters.NewRequireLoginFilter(), filters.NewAjaxFilter()).

    /**
     * 删除评论
     */
    Post("ajax-del", link_ajaxDel).Filters(filters.NewRequireLoginFilter(), filters.NewAjaxFilter())

//

// 删除link
func link_ajaxDel(ctx *goku.HttpContext) goku.ActionResulter {
    var errs string
    var ok = false

    linkId, err := strconv.ParseInt(ctx.RouteData.Params["id"], 10, 64)
    if err == nil {
        user := ctx.Data["user"].(*models.User)
        link, err := models.Link_GetById(linkId)
        if err == nil {
            // 只可以删除自己的链接
            if link.UserId == user.Id {
                err = models.Link_DelById(linkId)
                if err == nil {
                    ok = true
                }
            } else {
                errs = "不允许的操作"
            }
        }
    }

    if err != nil {
        errs = err.Error()
    }

    r := map[string]interface{}{
        "success": ok,
        "errors":  errs,
    }

    return ctx.Json(r)
}

// 查看link
func link_show(ctx *goku.HttpContext) goku.ActionResulter {
    return link_showWithComments(ctx, ctx.RouteData.Params["id"], "0")
}

func link_permacoment(ctx *goku.HttpContext) goku.ActionResulter {
    return link_showWithComments(ctx,
        ctx.RouteData.Params["id"], ctx.RouteData.Params["cid"])
}

var ORDER_NAMES map[string]string = map[string]string{
    "top":   "最佳",
    "hot":   "热议",
    "later": "最新",
    "vote":  "得分",
}

func link_showWithComments(ctx *goku.HttpContext, slinkId, scommentId string) goku.ActionResulter {

    linkId, err1 := strconv.ParseInt(slinkId, 10, 64)
    commentId, err2 := strconv.ParseInt(scommentId, 10, 64)
    if err1 != nil || err2 != nil {
        ctx.ViewData["errorMsg"] = "参数错误"
        return ctx.Render("error", nil)
    }
    link, err := models.Link_GetById(linkId)
    if err != nil {
        ctx.ViewData["errorMsg"] = "服务器开小差了 >_<!!"
        return ctx.Render("error", nil)
    }

    if link == nil {
        ctx.ViewData["errorMsg"] = "内容不存在"
        return ctx.Render("error", nil)
    }

    vlink := models.Link_ToVLink([]models.Link{*link}, ctx)
    sortType := strings.ToLower(ctx.Get("cm_order")) //"top":热门；"hot":热议；"later":最新；"vote":得分；
    if sortType == "" {
        sortType = "top"
    }
    var comments string
    if commentId > 0 {
        comments = models.GetPermalinkComment(linkId, commentId, sortType)
        ctx.ViewData["SubLinkUrl"] = fmt.Sprintf("permacoment/%d/%d/", linkId, commentId)
    } else {
        comments = models.GetSortComments("", "/", int64(0), linkId, sortType, "", false) //models.Comment_SortForLink(link.Id, "hot")
        ctx.ViewData["SubLinkUrl"] = linkId
    }

    ctx.ViewData["Comments"] = template.HTML(comments)
    ctx.ViewData["SortType"] = sortType
    ctx.ViewData["SortTypeName"] = ORDER_NAMES[sortType]

    return ctx.Render("/link/show", vlink[0])
}

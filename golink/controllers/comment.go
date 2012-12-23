package controllers

import (
    "fmt"
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink/utils"
    //"github.com/QLeelulu/goku/form"
    "github.com/QLeelulu/ohlala/golink/filters"
    "github.com/QLeelulu/ohlala/golink/models"
    "strconv"
    //"time"
    //"github.com/QLeelulu/ohlala/golink"
    //"html/template"
)

type CommentHtml struct {
    Html string
}

/**
 * 评论
 */
var _ = goku.Controller("comment").
    /**
     * 加载更多评论
     */
    Post("loadmore", comment_LoadMore).
    /**
     * 收到的评论
     */
    Get("inbox", comment_Inbox).Filters(filters.NewRequireLoginFilter())

/**
 * 加载更多评论
 */
func comment_LoadMore(ctx *goku.HttpContext) goku.ActionResulter {

    htmlObject := CommentHtml{""}
    exceptIds := ctx.Get("except_ids")
    fmt.Println("exceptIds:", exceptIds)
    parentPath := ctx.Get("parent_path")
    sortType := ctx.Get("sort_type")
    topId, err1 := strconv.ParseInt(ctx.Get("top_parent_id"), 10, 64)
    linkId, err2 := strconv.ParseInt(ctx.Get("link_id"), 10, 64)
    if err1 == nil && err2 == nil {
        htmlObject.Html = models.GetSortComments(exceptIds, parentPath, topId, linkId, sortType, "", true)
    }

    return ctx.Json(htmlObject)
}

/**
 * 收到的评论
 */
func comment_Inbox(ctx *goku.HttpContext) goku.ActionResulter {
    user := ctx.Data["user"].(*models.User)
    page, pagesize := utils.PagerParams(ctx.Request)
    comments, total, err := models.CommentForUser_GetByPage(user.Id, page, pagesize, "")
    if err != nil {
        ctx.ViewData["errorMsg"] = err.Error()
        return ctx.Render("error", nil)
    }
    ctx.ViewData["CommentList"] = comments
    ctx.ViewData["CommentCount"] = total
    ctx.ViewData["Page"] = page
    ctx.ViewData["Pagesize"] = pagesize
    err = models.Remind_Reset(user.Id, models.REMIND_COMMENT)
    if err != nil {
        goku.Logger().Errorln("Reset用户提醒信息数出错：", err.Error())
    }
    return ctx.View(nil)
}

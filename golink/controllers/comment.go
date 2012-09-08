package controllers

import (
    "fmt"
    "github.com/QLeelulu/goku"
    //"github.com/QLeelulu/goku/form"
    //"github.com/QLeelulu/ohlala/golink/filters"
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
 * 追加评论
 */
var _ = goku.Controller("comment").
    /**
     * 追加评论
     */
    Post("loadmore", func(ctx *goku.HttpContext) goku.ActionResulter {

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
})


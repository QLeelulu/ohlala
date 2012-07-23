package controllers

import (
    //"fmt"
    "github.com/QLeelulu/goku"
    //"github.com/QLeelulu/goku/form"
    "github.com/QLeelulu/ohlala/golink/models"
    "strconv"
    //"time"
    "github.com/QLeelulu/ohlala/golink"
)


/**
 * vote controller
 */
var _ = goku.Controller("vote").
    /**
     * 投票链接
     */
    Get("link", func(ctx *goku.HttpContext) goku.ActionResulter {

    vote := &models.Vote{0, 0, false}
    id, err1 := strconv.Atoi(ctx.RouteData.Params["id"])
    votetype, err2 := strconv.Atoi(ctx.RouteData.Params["votetype"])
    var score int = 1 //vote up
    if votetype == 2 { //vote down
	score = -1
    }
    var userId int64 = 1 //TODO:

    if err1 == nil && err2 == nil {
        vote = models.VoteLink(int64(id), userId, score, golink.SITERUNTIME)
    }

    return ctx.Json(vote)
}).

    /**
     * 投票评论
     */
    Get("comment", func(ctx *goku.HttpContext) goku.ActionResulter {

    id, err := strconv.Atoi(ctx.RouteData.Params["id"])
    if err == nil {
        var todo models.Todo
        todo, err = models.GetTodo(id)
        if err == nil {
            return ctx.View(todo)
        }
    }
    ctx.ViewData["errorMsg"] = err.Error()
    return ctx.Render("error", nil)
})

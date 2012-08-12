package controllers

import (
    //"fmt"
    "github.com/QLeelulu/goku"
    //"github.com/QLeelulu/goku/form"
    "github.com/QLeelulu/ohlala/golink/filters"
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
    var score int = 1  //vote up
    if votetype == 2 { //vote down
        score = -1
    }
    var userId int64 = (ctx.Data["user"].(*models.User)).Id

    if err1 == nil && err2 == nil {
        vote = models.VoteLink(int64(id), userId, score, golink.SITERUNTIME)
    }

    return ctx.Json(vote)
}).Filters(filters.NewRequireLoginFilter()).

    /**
     * 投票评论
     */
    Get("comment", func(ctx *goku.HttpContext) goku.ActionResulter {

    vote := &models.Vote{0, 0, false}
    id, err1 := strconv.Atoi(ctx.RouteData.Params["id"])
    //topId, err2 := strconv.Atoi(ctx.RouteData.Params["topid"])
    votetype, err3 := strconv.Atoi(ctx.RouteData.Params["votetype"])

    var score int = 1 //vote up
    if votetype == 2 { 
	score = -1 //vote down
    }
    var userId int64 = (ctx.Data["user"].(*models.User)).Id

    if err1 == nil && err3 == nil { //err2 == nil && 
        vote = models.VoteComment(int64(id), userId, score, golink.SITERUNTIME) //int64(topId), 
    }

    return ctx.Json(vote)

}).Filters(filters.NewRequireLoginFilter())



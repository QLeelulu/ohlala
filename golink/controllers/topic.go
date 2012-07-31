package controllers

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink/filters"
    "github.com/QLeelulu/ohlala/golink/models"
    "strconv"
)

/**
 * Controller: topic
 */
var _ = goku.Controller("topic").

    /**
     * 查看话题信息页
     */
    Get("show", func(ctx *goku.HttpContext) goku.ActionResulter {

    topicName, _ := ctx.RouteData.Params["name"]
    topic, _ := models.Topic_GetByName(topicName)

    if topic == nil {
        ctx.ViewData["errorMsg"] = "话题不存在"
        return ctx.Render("error", nil)
    }

    links, _ := models.Link_ForTopic(topic.Id, 1, 20)
    followers, _ := models.Topic_GetFollowers(topic.Id, 1, 12)

    ctx.ViewData["Links"] = links
    ctx.ViewData["Followers"] = followers
    return ctx.View(topic)

}).
    Filters(filters.NewRequireLoginFilter()).

    /**
     * 关注话题
     */
    Post("follow", func(ctx *goku.HttpContext) goku.ActionResulter {

    topicId, _ := strconv.ParseInt(ctx.RouteData.Params["id"], 10, 64)
    ok, err := models.Topic_Follow(ctx.Data["user"].(*models.User).Id, topicId)
    var errs string
    if err != nil {
        errs = err.Error()
    }
    r := map[string]interface{}{
        "success": ok,
        "errors":  errs,
    }
    return ctx.Json(r)

}).
    Filters(filters.NewRequireLoginFilter(), filters.NewAjaxFilter())

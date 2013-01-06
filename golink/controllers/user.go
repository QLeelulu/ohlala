package controllers

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink"
    "github.com/QLeelulu/ohlala/golink/filters"
    "github.com/QLeelulu/ohlala/golink/models"
    "github.com/QLeelulu/ohlala/golink/utils"
    "strconv"
)

/**
 * Controller: user
 */
var _ = goku.Controller("user").

    /**
     * 查看关注的人
     */
    Get("follows", user_Follows).
    /**
     * 查看用户的粉丝
     */
    Get("fans", user_Fans).

    /**
     * follow somebody
     */
    Post("follow", func(ctx *goku.HttpContext) goku.ActionResulter {

    followId, _ := strconv.ParseInt(ctx.RouteData.Params["id"], 10, 64)
    ok, err := models.User_Follow(ctx.Data["user"].(*models.User).Id, followId)
    var errs string
    if err != nil {
        errs = err.Error()
    }
    r := map[string]interface{}{
        "success": ok,
        "errors":  errs,
    }
    return ctx.Json(r)

}).Filters(filters.NewRequireLoginFilter(), filters.NewAjaxFilter()).

    /**
     * follow somebody
     */
    Post("unfollow", func(ctx *goku.HttpContext) goku.ActionResulter {

    followId, _ := strconv.ParseInt(ctx.RouteData.Params["id"], 10, 64)
    ok, err := models.User_UnFollow(ctx.Data["user"].(*models.User).Id, followId)
    var errs string
    if err != nil {
        errs = err.Error()
    }
    r := map[string]interface{}{
        "success": ok,
        "errors":  errs,
    }
    return ctx.Json(r)

}).Filters(filters.NewRequireLoginFilter(), filters.NewAjaxFilter()).

    /**
     * 查看用户信息页
     */
    Get("show", func(ctx *goku.HttpContext) goku.ActionResulter {

    userId, _ := strconv.ParseInt(ctx.RouteData.Params["id"], 10, 64)
    user := models.User_GetById(userId)

    if user == nil {
        ctx.ViewData["errorMsg"] = "用户不存在"
        return ctx.Render("error", nil)
    }

    links := models.Link_ByUser(user.Id, 1, golink.PAGE_SIZE)
    friends, _ := models.UserFollow_Friends(user.Id, 1, 21)
    followers, _ := models.UserFollow_Followers(user.Id, 1, 21)
    followTopics, _ := models.User_GetFollowTopics(user.Id, 1, 21, "")

    ctx.ViewData["Links"] = models.Link_ToVLink(links, ctx)
    ctx.ViewData["Friends"] = friends
    ctx.ViewData["Followers"] = followers
    ctx.ViewData["FollowTopics"] = followTopics
    ctx.ViewData["HasMoreLink"] = len(links) >= golink.PAGE_SIZE
    return ctx.View(models.User_ToVUser(user, ctx))

}).Filters(filters.NewRequireLoginFilter()).

    /**
     * 获取用户信息
     * 用于浮动层
     */
    Get("pbox-info", func(ctx *goku.HttpContext) goku.ActionResulter {

    userId, _ := strconv.ParseInt(ctx.RouteData.Params["id"], 10, 64)
    user := models.User_GetById(userId)

    if user != nil {
        return ctx.RenderPartial("pop-info", models.User_ToVUser(user, ctx))
    }
    return ctx.Html("")

}).Filters(filters.NewAjaxFilter()).

    /**
     * 新评论、新关注等提醒
     */
    Get("remind", func(ctx *goku.HttpContext) goku.ActionResulter {

    user := ctx.Data["user"].(*models.User)
    ok := false
    errs := ""
    remindInfo, err := models.Remind_ForUser(user.Id)
    if err == nil {
        ok = true
    } else {
        errs = err.Error()
    }

    r := map[string]interface{}{
        "success": ok,
        "errors":  errs,
        "remind":  remindInfo,
    }
    return ctx.Json(r)

}).Filters(filters.NewRequireLoginFilter(), filters.NewAjaxFilter()).

    /**
     * 加载更多链接
     */
    Get("loadmorelink", func(ctx *goku.HttpContext) goku.ActionResulter {

    page, pagesize := utils.PagerParams(ctx.Request)
    success, hasmore := false, false
    errorMsgs, html := "", ""
    if page > 1 {
        userId, _ := strconv.ParseInt(ctx.RouteData.Params["id"], 10, 64)
        user := models.User_GetById(userId)

        if user == nil {
            errorMsgs = "用户不存在"
        } else {
            // ot := ctx.Get("o")
            // if ot == "" {
            //     ot = "hot"
            // }
            // links, _ := models.Link_ByUser(user.Id, ot, page, golink.PAGE_SIZE)
            links := models.Link_ByUser(user.Id, page, pagesize)
            if links != nil && len(links) > 0 {
                ctx.ViewData["Links"] = models.Link_ToVLink(links, ctx)
                vr := ctx.RenderPartial("loadmorelink", nil)
                vr.Render(ctx, vr.Body)
                html = vr.Body.String()
                hasmore = len(links) >= pagesize
            }
            success = true
        }
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

}).Filters(filters.NewRequireLoginFilter(), filters.NewAjaxFilter())

//
// ==>>>>

// 查看关注的人
func user_Follows(ctx *goku.HttpContext) goku.ActionResulter {

    userId, _ := strconv.ParseInt(ctx.RouteData.Params["id"], 10, 64)
    var user *models.User
    if userId > 0 {
        user = models.User_GetById(userId)
    } else {
        if u, ok := ctx.Data["user"]; ok {
            user = u.(*models.User)
            ctx.ViewData["UserMenu"] = "um-follows"
        }
    }

    if user == nil {
        ctx.ViewData["errorMsg"] = "用户不存在"
        return ctx.Render("error", nil)
    }

    page, pagesize := utils.PagerParams(ctx.Request)
    friends, _ := models.UserFollow_Friends(user.Id, page, pagesize)

    ctx.ViewData["Friends"] = models.User_ToVUsers(friends, ctx)
    ctx.ViewData["HasMoreFriends"] = len(friends) >= pagesize
    return ctx.View(models.User_ToVUser(user, ctx))

}

// 查看粉丝
func user_Fans(ctx *goku.HttpContext) goku.ActionResulter {

    userId, _ := strconv.ParseInt(ctx.RouteData.Params["id"], 10, 64)
    var user *models.User
    if userId > 0 {
        user = models.User_GetById(userId)
    } else {
        if u, ok := ctx.Data["user"]; ok {
            user = u.(*models.User)
            ctx.ViewData["UserMenu"] = "um-fans"
        }
    }

    if user == nil {
        ctx.ViewData["errorMsg"] = "用户不存在"
        return ctx.Render("error", nil)
    }

    page, pagesize := utils.PagerParams(ctx.Request)
    followers, _ := models.UserFollow_Followers(user.Id, page, pagesize)

    ctx.ViewData["Followers"] = models.User_ToVUsers(followers, ctx)
    ctx.ViewData["HasMoreFollowers"] = len(followers) >= pagesize
    return ctx.View(models.User_ToVUser(user, ctx))

}

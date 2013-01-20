package controllers

import (
    "github.com/QLeelulu/goku"
    //"github.com/QLeelulu/goku/form"
    "errors"
    "github.com/QLeelulu/ohlala/golink"
    "github.com/QLeelulu/ohlala/golink/filters"
    "github.com/QLeelulu/ohlala/golink/models"
    "strconv"
    "strings"
    "time"
)

type FavoriteResult struct {
    Success bool   `json:"success"`
    Errors  string `json:"errors"`
}

var _ = goku.Controller("favorite").
    /**
     * 用户收藏link的首页
     */
    Get("index", func(ctx *goku.HttpContext) goku.ActionResulter {

    u, _ := ctx.Data["user"]
    user := u.(*models.User)
    links := models.FavoriteLink_ByUser(user.Id, 1, golink.PAGE_SIZE)
    ctx.ViewData["Links"] = models.Link_ToVLink(links, ctx)
    ctx.ViewData["HasMoreLink"] = len(links) >= golink.PAGE_SIZE
    ctx.ViewData["UserMenu"] = "um-favorite"

    return ctx.Render("/favorite/show", nil)

}).Filters(filters.NewRequireLoginFilter()).

    /**
     * load more
     */
    Get("loadmorelink", favorite_loadMoreLink).
    Filters(filters.NewRequireLoginFilter(), filters.NewAjaxFilter()).

    /**
     * 对收藏link的操作:opration:[add:收藏 | del:取消]
     */
    Post("opration", func(ctx *goku.HttpContext) goku.ActionResulter {

    var userId int64
    if u, ok := ctx.Data["user"]; ok && u != nil {
        userId = (u.(*models.User)).Id
    } else {
        return ctx.Json(&FavoriteResult{false, "登录已超时,请重新登录!"})
    }

    var err error
    var strOpration string = ctx.Get("opration")
    var linkId int64
    linkId, err = strconv.ParseInt(ctx.Get("linkId"), 10, 64)
    if err != nil {
        return ctx.Json(&FavoriteResult{false, "linkId: 参数错误!"})
    }

    if strOpration == "add" {
        f := map[string]interface{}{
            "user_id":     userId,
            "link_id":     linkId,
            "create_time": time.Now(),
        }
        err = models.SaveUserFavorite(f)
    } else if strOpration == "del" {
        err = models.DelUserFavorite(userId, linkId)
    } else {
        err = errors.New("opration: 参数错误")
    }

    if err != nil {
        if strings.Contains(err.Error(), "Duplicate entry") {
            err = errors.New("已经收藏过了")
        }
        return ctx.Json(&FavoriteResult{false, err.Error()})
    }

    return ctx.Json(&FavoriteResult{true, ""})
}).Filters(filters.NewRequireLoginFilter())

func favorite_loadMoreLink(ctx *goku.HttpContext) goku.ActionResulter {
    page, err := strconv.Atoi(ctx.Get("page"))
    success, hasmore := false, false
    errorMsgs, html := "", ""
    if err == nil && page > 1 {
        user := ctx.Data["user"].(*models.User)

        links := models.FavoriteLink_ByUser(user.Id, page, golink.PAGE_SIZE)
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

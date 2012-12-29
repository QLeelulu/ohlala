package controllers

import (
    "github.com/QLeelulu/goku"
    //"github.com/QLeelulu/goku/form"
    "github.com/QLeelulu/ohlala/golink/filters"
    "github.com/QLeelulu/ohlala/golink/models"
    "strconv"
    "time"
    "github.com/QLeelulu/ohlala/golink"
)

type FavoriteResult struct {
	Result        bool
	Msg           string
}



var _ = goku.Controller("favorite").
    /**
     * 用户收藏link的首页
     */
	Get("user", func(ctx *goku.HttpContext) goku.ActionResulter {
	u, ok := ctx.Data["user"]
	if !ok || u == nil {
	    return ctx.Redirect("/discover")
	}
	user := u.(*models.User)
	links := models.FavoriteLink_ByUser(user.Id , 1, golink.PAGE_SIZE)
	ctx.ViewData["Links"] = models.Link_ToVLink(links, ctx)
	ctx.ViewData["HasMoreLink"] = len(links) >= golink.PAGE_SIZE

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
		return ctx.Json(&FavoriteResult{false, "请求出错,请重试!"})
	}

	if strOpration == "add" {
		f := map[string]interface{}{
		    "user_id": userId,
		    "link_id": linkId,
		    "create_time": time.Now(),
		}
		err = models.SaveUserFavorite(f)
	} else {
		err = models.DelUserFavorite(userId, linkId)
	}

	if err != nil {
		return ctx.Json(&FavoriteResult{false, "请求出错,请重试!"})
	}

	return ctx.Json(&FavoriteResult{true, ""})
})


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
























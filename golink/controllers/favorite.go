package controllers

import (
    "strings"
    "github.com/QLeelulu/goku"
    //"github.com/QLeelulu/goku/form"
    "github.com/QLeelulu/ohlala/golink/filters"
    "github.com/QLeelulu/ohlala/golink/models"
    "github.com/QLeelulu/ohlala/golink/utils"
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

    var userId int64 = (ctx.Data["user"].(*models.User)).Id
	links := models.FavoriteLink_ByUser(userId , page, pagesize)

	return ctx.Render("/favorite/user", links)

}).Filters(filters.NewRequireLoginFilter()).
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

	var strOpration string = ctx.Get("opration")
	var linkId int64 = strconv.ParseInt(ctx.Get("linkId"), 10, 64)
	var err error
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
	} else {
		return ctx.Json(&FavoriteResult{true, ""})
	}
})



























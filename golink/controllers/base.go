package controllers

import (
    "github.com/QLeelulu/goku"
    "net/url"
)

func RequireLogin(ctx *goku.HttpContext) goku.ActionResulter {
    if u, ok := ctx.Data["user"]; !ok || u == nil {
        return ctx.Redirect("/user/login?returnurl=" + url.QueryEscape(ctx.Request.RequestURI))
    }
    return nil
}

/**
 * filters
 */
// 需要先登陆，没有登陆则跳转到登陆页
type RequireLoginFilter struct {
}

func (tf *RequireLoginFilter) OnActionExecuting(ctx *goku.HttpContext) (ar goku.ActionResulter, err error) {
    if u, ok := ctx.Data["user"]; !ok || u == nil {
        ar = ctx.Redirect("/user/login?returnurl=" + url.QueryEscape(ctx.Request.RequestURI))
    }
    return
}
func (tf *RequireLoginFilter) OnActionExecuted(ctx *goku.HttpContext) (goku.ActionResulter, error) {
    return nil, nil
}

func (tf *RequireLoginFilter) OnResultExecuting(ctx *goku.HttpContext) (goku.ActionResulter, error) {
    return nil, nil
}

func (tf *RequireLoginFilter) OnResultExecuted(ctx *goku.HttpContext) (goku.ActionResulter, error) {
    return nil, nil
}

var requireLoginFilter = new(RequireLoginFilter)

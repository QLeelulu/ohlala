package filters

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink/models"
    "net/url"
)

// 需要先登陆，没有登陆则跳转到登陆页
type RequireLoginFilter struct {
}

func (tf *RequireLoginFilter) OnActionExecuting(ctx *goku.HttpContext) (ar goku.ActionResulter, err error) {
    if u, ok := ctx.Data["user"]; !ok || u == nil {
        if ctx.IsAjax() {
            ar = ctx.Json(map[string]interface{}{
                "success":   false,
                "needLogin": true,
                "errors":    "请先登陆",
            })
        } else {
            ar = ctx.Redirect("/user/login?returnurl=" + url.QueryEscape(ctx.Request.RequestURI))
        }
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

type RequireAdminFilter struct {
    RequireLoginFilter
}

func (raf *RequireAdminFilter) OnActionExecuting(ctx *goku.HttpContext) (ar goku.ActionResulter, err error) {
    ar, err = raf.RequireLoginFilter.OnActionExecuting(ctx)
    if ar != nil || err != nil {
        return
    }
    user := ctx.Data["user"].(*models.User)
    if !user.IsAdmin() {
        if ctx.IsAjax() {
            ar = ctx.Json(map[string]interface{}{
                "success":   false,
                "needLogin": false,
                "errors":    "没有权限",
            })
        } else {
            // ctx.ViewData["errorMsg"] = "没有权限"
            // ar = ctx.Render("error", nil)
            ar = ctx.Raw("没有权限")
        }
    }
    return
}

var requireLoginFilter = new(RequireLoginFilter)
var requireAdminFilter = new(RequireAdminFilter)

func NewRequireLoginFilter() *RequireLoginFilter {
    return requireLoginFilter
}

func NewRequireAdminFilter() *RequireAdminFilter {
    return requireAdminFilter
}

package filters

import (
    "github.com/QLeelulu/goku"
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

var requireLoginFilter = new(RequireLoginFilter)

func NewRequireLoginFilter() *RequireLoginFilter {
    return requireLoginFilter
}

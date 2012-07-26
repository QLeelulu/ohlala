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

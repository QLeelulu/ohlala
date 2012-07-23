package middlewares

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink/models"
)

// 一些基本的处理
// 例如检查用户是否登陆，如果登陆则获取登陆用户信息，并添加 ctx.Data 中
type UtilMiddleware struct {
}

func (tmd *UtilMiddleware) OnBeginRequest(ctx *goku.HttpContext) (goku.ActionResulter, error) {
    return nil, nil
}

func (tmd *UtilMiddleware) OnBeginMvcHandle(ctx *goku.HttpContext) (goku.ActionResulter, error) {
    getUser(ctx)
    return nil, nil
}
func (tmd *UtilMiddleware) OnEndMvcHandle(ctx *goku.HttpContext) (goku.ActionResulter, error) {
    return nil, nil
}

func (tmd *UtilMiddleware) OnEndRequest(ctx *goku.HttpContext) (goku.ActionResulter, error) {
    return nil, nil
}

func getUser(ctx *goku.HttpContext) {
    c, err := ctx.Request.Cookie("_glut")
    if err == nil {
        user, _ := models.User_GetByTicket(c.Value)
        if user != nil {
            ctx.Data["user"] = user
            // 暂时先设置到ViewData里面吧，应该需要一个更好的办法？
            ctx.ViewData["user"] = user
        }
    }
}

package controllers

// import (
//     "github.com/QLeelulu/goku"
//     "github.com/QLeelulu/ohlala/golink/filters"
//     "github.com/QLeelulu/ohlala/golink/forms"
//     "github.com/QLeelulu/ohlala/golink/models"
//     "strconv"
// )

// var _ = goku.Controller("comment").

//     /**
//      * 提交评论并保存到数据库
//      */
//     Post("submit", func(ctx *goku.HttpContext) goku.ActionResulter {

//     f := forms.NewCommentSubmitForm()
//     f.FillByRequest(ctx.Request)

//     success, errorMsgs := models.Comment_SaveForm(f, (ctx.Data["user"].(*models.User)).Id)

//     if success {
//         return ctx.Redirect("/")
//     } else {
//         ctx.ViewData["Errors"] = errorMsgs
//         ctx.ViewData["Values"] = f.Values()
//     }
//     return ctx.View(nil)

// }).Filters(filters.NewRequireLoginFilter()).

//     /**
//      * 添加评论
//      */
//     Post("comment", func(ctx *goku.HttpContext) goku.ActionResulter {

//     return ctx.View(nil)

// }).Filters(filters.NewRequireLoginFilter())

package admin

import (
    "github.com/QLeelulu/goku"
    // "github.com/QLeelulu/ohlala/golink"
    // "github.com/QLeelulu/ohlala/golink/filters"
    "github.com/QLeelulu/ohlala/golink/models"
    // "strconv"
)

var _ = adminController.
    // index
    Get("index", admin_index)

//

func admin_index(ctx *goku.HttpContext) goku.ActionResulter {
    var db *goku.MysqlDB = models.GetDB()
    defer db.Close()

    linkCount, err := db.Count("link", "")
    if err != nil {
        ctx.ViewData["errorMsg"] = err.Error()
        return ctx.Render("error", nil)
    }
    ctx.ViewData["linkCount"] = linkCount

    userCount, err := db.Count("user", "")
    if err != nil {
        ctx.ViewData["errorMsg"] = err.Error()
        return ctx.Render("error", nil)
    }
    ctx.ViewData["userCount"] = userCount

    topicCount, err := db.Count("topic", "")
    if err != nil {
        ctx.ViewData["errorMsg"] = err.Error()
        return ctx.Render("error", nil)
    }
    ctx.ViewData["topicCount"] = topicCount

    commentCount, err := db.Count("comment", "")
    if err != nil {
        ctx.ViewData["errorMsg"] = err.Error()
        return ctx.Render("error", nil)
    }
    ctx.ViewData["commentCount"] = commentCount

    return ctx.View(nil)
}

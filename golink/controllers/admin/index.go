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

    // u, ok := ctx.Data["user"]
    // if !ok || u == nil {
    //     return ctx.Redirect("/discover")
    // }
    // user := u.(*models.User)
    // ot := ctx.Get("o")
    // if ot == "" {
    //     ot = "hot"
    // }
    // ctx.ViewData["Order"] = ot
    // links, _ := models.Link_ForUser(user.Id, ot, 1, golink.PAGE_SIZE) //models.Link_GetByPage(1, 20)
    // ctx.ViewData["Links"] = models.Link_ToVLink(links, ctx)
    // ctx.ViewData["HasMoreLink"] = len(links) >= golink.PAGE_SIZE
    // return ctx.View(nil)
}

package admin

import (
    "github.com/QLeelulu/goku"
    // "github.com/QLeelulu/ohlala/golink"
    // "github.com/QLeelulu/ohlala/golink/filters"
    // "github.com/QLeelulu/ohlala/golink/models"
    // "strconv"
)

var _ = adminController.
    // index
    Get("users", admin_users)

//

func admin_users(ctx *goku.HttpContext) goku.ActionResulter {
    return ctx.Html("admin")
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

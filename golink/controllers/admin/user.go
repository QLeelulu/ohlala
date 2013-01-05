package admin

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink/utils"
    // "github.com/QLeelulu/ohlala/golink/filters"
    "github.com/QLeelulu/ohlala/golink/models"
    // "strconv"
)

var _ = adminController.
    // index
    Get("users", admin_users)

//

func admin_users(ctx *goku.HttpContext) goku.ActionResulter {
    page, pagesize := utils.PagerParams(ctx.Request)
    users, total, err := models.User_GetList(page, pagesize, "")
    if err != nil {
        ctx.ViewData["errorMsg"] = err.Error()
        return ctx.Render("error", nil)
    }
    ctx.ViewData["UserList"] = users
    ctx.ViewData["UserCount"] = total
    ctx.ViewData["Page"] = page
    ctx.ViewData["Pagesize"] = pagesize
    ctx.ViewData["TabName"] = "users"
    return ctx.View(nil)
}

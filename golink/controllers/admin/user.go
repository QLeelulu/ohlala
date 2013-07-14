package admin

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink/utils"
    // "github.com/QLeelulu/ohlala/golink/filters"
    "github.com/QLeelulu/ohlala/golink/models"
    "strconv"
)

var _ = adminController.
    // index
    Get("users", admin_users).
    // 禁言用户等
    Post("ban_user", admin_banUser)

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

func admin_banUser(ctx *goku.HttpContext) goku.ActionResulter {
    var err error
    var errs string
    var ok = false
    var userId, status int64

    userId, err = strconv.ParseInt(ctx.Get("id"), 10, 64)
    if err == nil {
        status, err = strconv.ParseInt(ctx.Get("status"), 10, 64)
    }
    if err == nil {
        _, err = models.User_Update(userId, map[string]interface{}{"Status": status})
    }

    if err != nil {
        errs = err.Error()
    } else {
        ok = true
    }
    r := map[string]interface{}{
        "success": ok,
        "errors":  errs,
    }

    return ctx.Json(r)
}

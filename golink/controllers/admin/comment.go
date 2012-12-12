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
    Get("comments", admin_comments).
    // 删除link
    Post("del_comment", admin_del_comments)

//

func admin_comments(ctx *goku.HttpContext) goku.ActionResulter {
    page, pagesize := utils.PagerParams(ctx.Request)
    comments, total, err := models.Comment_GetByPage(page, pagesize, "")
    if err != nil {
        ctx.ViewData["errorMsg"] = err.Error()
        return ctx.Render("error", nil)
    }
    ctx.ViewData["CommentList"] = comments
    ctx.ViewData["CommentCount"] = total
    ctx.ViewData["Page"] = page
    ctx.ViewData["Pagesize"] = pagesize
    return ctx.View(nil)
}

// 删除comment
func admin_del_comments(ctx *goku.HttpContext) goku.ActionResulter {
    var errs string
    var ok = false

    id, err := strconv.ParseInt(ctx.RouteData.Params["id"], 10, 64)
    if err == nil {
        err = models.Comment_DelById(id)
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

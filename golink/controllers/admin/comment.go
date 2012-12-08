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
    Get("comments", admin_comments)

//

func admin_comments(ctx *goku.HttpContext) goku.ActionResulter {
    comments, err := models.Comment_GetByPage(1, 20, "")
    if err != nil {
        ctx.ViewData["errorMsg"] = err.Error()
        return ctx.Render("error", nil)
    }
    ctx.ViewData["CommentList"] = comments
    return ctx.View(nil)
}

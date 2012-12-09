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
    Get("comments", admin_comments).
	// 删除link
	Get("comments", admin_del_comments)

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

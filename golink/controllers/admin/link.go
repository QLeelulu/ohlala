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
    Get("links", admin_links)

//

func admin_links(ctx *goku.HttpContext) goku.ActionResulter {
    links, err := models.Link_GetByPage(1, 20, "")
    if err != nil {
        ctx.ViewData["errorMsg"] = err.Error()
        return ctx.Render("error", nil)
    }
    ctx.ViewData["LinkList"] = links
    return ctx.View(nil)
}

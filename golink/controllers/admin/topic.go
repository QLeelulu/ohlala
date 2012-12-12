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
    Get("topics", admin_topics)

//

func admin_topics(ctx *goku.HttpContext) goku.ActionResulter {
    page, pagesize := utils.PagerParams(ctx.Request)
    topics, total, err := models.Topic_GetByPage(page, pagesize, "")
    if err != nil {
        ctx.ViewData["errorMsg"] = err.Error()
        return ctx.Render("error", nil)
    }
    ctx.ViewData["TopicList"] = topics
    ctx.ViewData["TopicCount"] = total
    ctx.ViewData["Page"] = page
    ctx.ViewData["Pagesize"] = pagesize
    return ctx.View(nil)
}

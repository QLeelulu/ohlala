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
    Get("topics", admin_topics).
    /**
     * 修改话题名称
     */
    Post("topic_editname", admin_topicEditName)

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

/**
 * 修改话题名称
 */
func admin_topicEditName(ctx *goku.HttpContext) goku.ActionResulter {
    var ok = false
    var errs, name string
    topicId, err := strconv.ParseInt(ctx.Get("id"), 10, 64)
    name = ctx.Request.FormValue("name")
    if err == nil && topicId > 0 && name != "" {
        _, err = models.Topic_UpdateName(topicId, name)
        if err == nil {
            ok = true
        }
    } else if topicId < 1 || name == "" {
        errs = "参数错误"
    }

    if err != nil {
        errs = err.Error()
    }
    r := map[string]interface{}{
        "success": ok,
        "name":    name,
        "errors":  errs,
    }
    return ctx.Json(r)
}

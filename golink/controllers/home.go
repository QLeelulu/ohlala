package controllers

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink/models"
)

var _ = goku.Controller("home").
    Get("index", func(ctx *goku.HttpContext) goku.ActionResulter {
    links := models.Link_GetByPage(1, 20)
    ctx.ViewData["Links"] = links
    return ctx.View(nil)
})

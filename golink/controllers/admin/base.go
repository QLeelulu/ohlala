package admin

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink/filters"
)

var adminController *goku.ControllerBuilder = goku.Controller("_golink_admin").
    Filters(filters.NewRequireAdminFilter())

// render the view and return a *ViewResult
// it will find the view in these rules:
//      1. /{ViewPath}/{Controller}/{action}
//      2. /{ViewPath}/shared/{action}
// func adminView(ctx *goku.HttpContext, viewData interface{}) *goku.ViewResult {
//     return ctx.RenderWithLayout("", "adminLayout", viewModel)
// }

package admin

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink/filters"
)

var adminController *goku.ControllerBuilder = goku.Controller("_golink_admin").
    Filters(filters.NewRequireAdminFilter())

package golink

import (
    "github.com/QLeelulu/goku"
)

// routes
var Routes []*goku.Route = []*goku.Route{
    &goku.Route{
        Name:     "static",
        IsStatic: true,
        Pattern:  "/assets/(.*)",
    },
    &goku.Route{
        Name:       "edit",
        Pattern:    "/{controller}/{id}/{action}",
        Default:    map[string]string{"action": "edit"},
        Constraint: map[string]string{"id": "\\d+"},
    },
    &goku.Route{
        Name:    "default",
        Pattern: "/{controller}/{action}",
        Default: map[string]string{"controller": "home", "action": "index"},
    },
}

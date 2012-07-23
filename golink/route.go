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
        Default: map[string]string{"controller": "todo", "action": "index"},
    },
    &goku.Route{
        Name:       "vote",
        Pattern:    "/{controller}/{action}/{id}/{votetype}/{topid}", //1 == vote up, 2 == votet down
        Default:    map[string]string{"action": "link"},
        Constraint: map[string]string{"id": "\\d+", "topid": "\\d+", "votetype" : "\\d+"},
    },
}

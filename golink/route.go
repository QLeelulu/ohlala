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
        Name:    "topicInfo",
        Pattern: "/t/{name}/",
        Default: map[string]string{"controller": "topic", "action": "show"},
    },
    &goku.Route{
        Name:       "threeNumParam",
        Pattern:    "/{controller}/{action}/{id}/{cid}/{arg}/",
        Constraint: map[string]string{"id": "\\d+", "cid": "\\d+"},
        Default:    map[string]string{"arg": ""},
    },
    &goku.Route{
        Name:       "edit",
        Pattern:    "/{controller}/{id}/{action}/",
        Default:    map[string]string{"action": "show"},
        Constraint: map[string]string{"id": "\\d+"},
    },
    &goku.Route{
        Name:    "default",
        Pattern: "/{controller}/{action}/{arg}/",
        Default: map[string]string{"controller": "home", "action": "index", "arg": ""},
    },
}

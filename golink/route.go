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
        Name:       "twoNumParam",
        Pattern:    "/{controller}/{action}/{id}/{arg}/",
        Constraint: map[string]string{"id": "\\d+", "arg": "\\d+"},
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
    &goku.Route{
        Name:       "threeNumParam",
        Pattern:    "/{controller}/{action}/{lid}/{cid}/{arg}/",
        Constraint: map[string]string{"lid": "\\d+", "cid": "\\d+"},
    },
}

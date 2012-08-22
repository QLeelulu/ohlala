package main

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink"
    _ "github.com/QLeelulu/ohlala/golink/controllers" // notice this!! import controllers
    "github.com/QLeelulu/ohlala/golink/middlewares"
    "log"
)

func main() {

    rt := &goku.RouteTable{Routes: golink.Routes}
    middlewares := []goku.Middlewarer{
        new(middlewares.UtilMiddleware),
    }
    s := goku.CreateServer(rt, middlewares, golink.Config)
    goku.Logger().Logln("Server start on", s.Addr)

    log.Fatal(s.ListenAndServe())
}

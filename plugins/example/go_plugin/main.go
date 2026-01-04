package main

import (
	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/modules/auth"
	c "github.com/purpose168/GoAdmin/modules/config"
	"github.com/purpose168/GoAdmin/modules/db"
	"github.com/purpose168/GoAdmin/modules/service"
	"github.com/purpose168/GoAdmin/plugins"
)

type Example struct {
	*plugins.Base
}

var Plugin = &Example{
	Base: &plugins.Base{PlugName: "example"},
}

func (example *Example) InitPlugin(srv service.List) {
	example.InitBase(srv, "example")
	Plugin.App = example.initRouter(c.Prefix(), srv)
}

func (example *Example) initRouter(prefix string, srv service.List) *context.App {

	app := context.NewApp()
	route := app.Group(prefix)
	route.GET("/example", auth.Middleware(db.GetConnection(srv)), example.TestHandler)

	return app
}

func (example *Example) TestHandler(ctx *context.Context) {

}

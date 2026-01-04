package example

import (
	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/modules/auth"
	"github.com/purpose168/GoAdmin/modules/db"
	"github.com/purpose168/GoAdmin/modules/service"
)

func (e *Example) initRouter(prefix string, srv service.List) *context.App {

	app := context.NewApp()
	route := app.Group(prefix)
	route.GET("/example", auth.Middleware(db.GetConnection(srv)), e.TestHandler)

	return app
}

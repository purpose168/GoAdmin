package fasthttp

import (
	// add fasthttp adapter
	ada "github.com/purpose168/GoAdmin/adapter/fasthttp"
	// add mysql driver
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/mysql"
	// add postgresql driver
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/postgres"
	// add sqlite driver
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/sqlite"
	// add mssql driver
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/mssql"
	// add adminlte ui theme
	"github.com/purpose168/GoAdmin-themes/adminlte"

	"os"

	"github.com/buaazp/fasthttprouter"
	"github.com/purpose168/GoAdmin/engine"
	"github.com/purpose168/GoAdmin/modules/config"
	"github.com/purpose168/GoAdmin/modules/language"
	"github.com/purpose168/GoAdmin/plugins/admin"
	"github.com/purpose168/GoAdmin/plugins/admin/modules/table"
	"github.com/purpose168/GoAdmin/template"
	"github.com/purpose168/GoAdmin/template/chartjs"
	"github.com/purpose168/GoAdmin/tests/tables"
	"github.com/valyala/fasthttp"
)

func internalHandler() fasthttp.RequestHandler {
	router := fasthttprouter.New()

	eng := engine.Default()

	adminPlugin := admin.NewAdmin(tables.Generators).AddDisplayFilterXssJsFilter()
	adminPlugin.AddGenerator("user", tables.GetUserTable)

	template.AddComp(chartjs.NewChart())

	if err := eng.AddConfigFromJSON(os.Args[len(os.Args)-1]).
		AddPlugins(adminPlugin).
		Use(router); err != nil {
		panic(err)
	}

	eng.HTML("GET", "/admin", tables.GetContent)

	return func(ctx *fasthttp.RequestCtx) {
		router.Handler(ctx)
	}
}

func NewHandler(dbs config.DatabaseList, gens table.GeneratorList) fasthttp.RequestHandler {
	router := fasthttprouter.New()

	eng := engine.Default()

	template.AddComp(chartjs.NewChart())

	if err := eng.AddConfig(&config.Config{
		Databases: dbs,
		UrlPrefix: "admin",
		Store: config.Store{
			Path:   "./uploads",
			Prefix: "uploads",
		},
		Language:    language.EN,
		IndexUrl:    "/",
		Debug:       true,
		ColorScheme: adminlte.ColorschemeSkinBlack,
	}).
		AddAdapter(new(ada.Fasthttp)).
		AddGenerators(gens).
		Use(router); err != nil {
		panic(err)
	}

	eng.HTML("GET", "/admin", tables.GetContent)

	return func(ctx *fasthttp.RequestCtx) {
		router.Handler(ctx)
	}
}

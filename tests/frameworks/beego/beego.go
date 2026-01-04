package beego

import (
	// add beego adapter
	_ "github.com/purpose168/GoAdmin/adapter/beego"
	"github.com/purpose168/GoAdmin/modules/config"
	"github.com/purpose168/GoAdmin/modules/language"
	"github.com/purpose168/GoAdmin/plugins/admin/modules/table"

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

	"net/http"
	"os"

	"github.com/astaxie/beego"
	"github.com/purpose168/GoAdmin/engine"
	"github.com/purpose168/GoAdmin/plugins/admin"
	"github.com/purpose168/GoAdmin/plugins/example"
	"github.com/purpose168/GoAdmin/template"
	"github.com/purpose168/GoAdmin/template/chartjs"
	"github.com/purpose168/GoAdmin/tests/tables"
)

func internalHandler() http.Handler {

	app := beego.NewApp()

	eng := engine.Default()
	adminPlugin := admin.NewAdmin(tables.Generators)
	adminPlugin.AddGenerator("user", tables.GetUserTable)

	examplePlugin := example.NewExample()

	if err := eng.AddConfigFromJSON(os.Args[len(os.Args)-1]).
		AddPlugins(adminPlugin, examplePlugin).Use(app); err != nil {
		panic(err)
	}

	template.AddComp(chartjs.NewChart())

	eng.HTML("GET", "/admin", tables.GetContent)

	beego.BConfig.Listen.HTTPAddr = "127.0.0.1"
	beego.BConfig.Listen.HTTPPort = 9087

	return app.Handlers
}

func NewHandler(dbs config.DatabaseList, gens table.GeneratorList) http.Handler {

	app := beego.NewApp()

	eng := engine.Default()
	adminPlugin := admin.NewAdmin(gens)

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
		AddPlugins(adminPlugin).Use(app); err != nil {
		panic(err)
	}

	template.AddComp(chartjs.NewChart())

	eng.HTML("GET", "/admin", tables.GetContent)

	beego.BConfig.Listen.HTTPAddr = "127.0.0.1"
	beego.BConfig.Listen.HTTPPort = 9087

	return app.Handlers
}

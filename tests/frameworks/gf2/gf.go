package gf

import (
	// add gf adapter
	"reflect"

	"github.com/agiledragon/gomonkey"
	_ "github.com/purpose168/GoAdmin/adapter/gf2"

	// add mysql driver
	"github.com/purpose168/GoAdmin/modules/config"
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/mysql"
	"github.com/purpose168/GoAdmin/modules/language"

	// add postgresql driver
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/postgres"
	// add sqlite driver
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/sqlite"
	// add mssql driver
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/mssql"
	// add adminlte ui theme
	"github.com/GoAdminGroup/themes/adminlte"

	"net/http"
	"os"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/purpose168/GoAdmin/engine"
	"github.com/purpose168/GoAdmin/plugins/admin"
	"github.com/purpose168/GoAdmin/plugins/admin/modules/table"
	"github.com/purpose168/GoAdmin/template"
	"github.com/purpose168/GoAdmin/template/chartjs"
	"github.com/purpose168/GoAdmin/tests/tables"
)

func internalHandler() http.Handler {
	s := g.Server()

	eng := engine.Default()

	adminPlugin := admin.NewAdmin(tables.Generators).AddDisplayFilterXssJsFilter()

	template.AddComp(chartjs.NewChart())

	adminPlugin.AddGenerator("user", tables.GetUserTable)

	if err := eng.AddConfigFromJSON(os.Args[len(os.Args)-1]).
		AddPlugins(adminPlugin).
		Use(s); err != nil {
		panic(err)
	}

	eng.HTML("GET", "/admin", tables.GetContent)

	s.SetPort(8103)

	gomonkey.ApplyMethod(reflect.TypeOf(new(ghttp.Request).Session), "Close",
		func(*ghttp.Session) error {
			return nil
		})

	return s
}

func NewHandler(dbs config.DatabaseList, gens table.GeneratorList) http.Handler {

	s := g.Server(8103)

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
		AddPlugins(adminPlugin).Use(s); err != nil {
		panic(err)
	}

	template.AddComp(chartjs.NewChart())

	eng.HTML("GET", "/admin", tables.GetContent)

	return s
}

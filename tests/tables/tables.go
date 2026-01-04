package tables

import "github.com/purpose168/GoAdmin/plugins/admin/modules/table"

var Generators = map[string]table.Generator{
	"posts":    GetPostsTable,
	"authors":  GetAuthorsTable,
	"external": GetExternalTable,
}

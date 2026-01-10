// beego.go - Beego 框架适配器文件
// 包名：beego
// 作者：GoAdmin 团队
// 创建日期：2020年
// 描述：本文件提供 GoAdmin 在 Beego Web 框架下的适配器实现
//       包含两个核心函数：internalHandler 和 NewHandler
//       internalHandler 用于内部测试，NewHandler 用于通用测试
//
// 主要功能：
//   - internalHandler: 创建内部 HTTP 处理器，用于测试内置表
//   - NewHandler: 创建通用 HTTP 处理器，支持自定义配置
//
// 技术栈：
//   - Beego: Go Web 框架，提供路由和中间件功能
//   - GoAdmin: Go 后台管理框架
//   - AdminLTE: UI 主题
//   - Chart.js: 图表组件
//
// 支持的数据库：
//   - MySQL
//   - PostgreSQL
//   - SQLite
//   - MSSQL
//
// 使用场景：
//   - 集成测试
//   - 黑盒测试
//   - 框架适配验证

package beego

import (
	_ "github.com/purpose168/GoAdmin/adapter/beego"             // Beego 适配器，用于将 GoAdmin 集成到 Beego 框架
	"github.com/purpose168/GoAdmin/modules/config"              // 配置模块
	"github.com/purpose168/GoAdmin/modules/language"            // 语言模块
	"github.com/purpose168/GoAdmin/plugins/admin/modules/table" // 表生成器模块

	"github.com/purpose168/GoAdmin-themes/adminlte"               // AdminLTE UI 主题
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/mssql"    // MSSQL 数据库驱动
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/mysql"    // MySQL 数据库驱动
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/postgres" // PostgreSQL 数据库驱动
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/sqlite"   // SQLite 数据库驱动

	"net/http" // Go 标准网络包，提供 HTTP 接口
	"os"       // Go 标准操作系统包，用于读取命令行参数

	"github.com/astaxie/beego"                       // Beego Web 框架
	"github.com/purpose168/GoAdmin/engine"           // GoAdmin 核心引擎
	"github.com/purpose168/GoAdmin/plugins/admin"    // GoAdmin 管理插件
	"github.com/purpose168/GoAdmin/plugins/example"  // GoAdmin 示例插件
	"github.com/purpose168/GoAdmin/template"         // GoAdmin 模板引擎
	"github.com/purpose168/GoAdmin/template/chartjs" // Chart.js 图表组件
	"github.com/purpose168/GoAdmin/tests/tables"     // 测试数据表
)

// ==================== 内部处理器函数 ====================

// internalHandler 创建内部 HTTP 处理器
// 返回值：
//   - http.Handler: HTTP 处理器，用于处理 HTTP 请求
//
// 说明：
//
//	该函数创建一个用于测试内置表的 HTTP 处理器。
//	处理器配置：
//	- 使用默认的 GoAdmin 引擎
//	- 使用内置的表生成器
//	- 从命令行参数读取配置文件路径
//	- 添加管理插件和示例插件
//	- 使用 AdminLTE 主题
//	- 监听地址：127.0.0.1:9087
//
// 配置说明：
//   - 数据库配置：从 JSON 文件读取（命令行最后一个参数）
//   - URL 前缀：默认为 "admin"
//   - 文件存储：./uploads 目录
//   - 语言：英语（EN）
//   - 首页 URL：/
//   - 调试模式：开启
//   - 主题：AdminLTE 黑色皮肤
//
// 使用场景：
//   - 内置表的黑盒测试
//   - 集成测试
//   - 功能验证测试
//
// 注意事项：
//   - 需要通过命令行参数传递配置文件路径
//   - 端口 9087 需要未被占用
//   - 配置文件需要包含数据库连接信息
func internalHandler() http.Handler {
	// 创建新的 Beego 应用实例
	app := beego.NewApp()

	// 创建 GoAdmin 默认引擎实例
	eng := engine.Default()
	// 创建管理插件，使用内置表生成器
	adminPlugin := admin.NewAdmin(tables.Generators)
	// 添加用户表生成器
	adminPlugin.AddGenerator("user", tables.GetUserTable)

	// 创建示例插件
	examplePlugin := example.NewExample()

	// 配置引擎：从 JSON 文件读取配置，添加插件，使用 Beego 应用
	if err := eng.AddConfigFromJSON(os.Args[len(os.Args)-1]).
		AddPlugins(adminPlugin, examplePlugin).Use(app); err != nil {
		panic(err) // 配置失败时抛出异常
	}

	// 添加 Chart.js 图表组件到模板引擎
	template.AddComp(chartjs.NewChart())

	// 注册 HTML 路由：GET /admin
	eng.HTML("GET", "/admin", tables.GetContent)

	// 配置 Beego 监听地址和端口
	beego.BConfig.Listen.HTTPAddr = "127.0.0.1" // 监听地址
	beego.BConfig.Listen.HTTPPort = 9087        // 监听端口

	// 返回 Beego 应用的处理器
	return app.Handlers
}

// ==================== 通用处理器函数 ====================

// NewHandler 创建新的 HTTP 处理器
// 参数：
//   - dbs: 数据库连接配置列表
//   - gens: 表生成器列表，用于自定义表
//
// 返回值：
//   - http.Handler: HTTP 处理器，用于处理 HTTP 请求
//
// 说明：
//
//	该函数创建一个通用的 HTTP 处理器，支持自定义数据库配置和表生成器。
//	处理器配置：
//	- 使用默认的 GoAdmin 引擎
//	- 使用自定义的表生成器
//	- 使用传入的数据库配置
//	- 添加管理插件
//	- 使用 AdminLTE 主题
//	- 监听地址：127.0.0.1:9087
//
// 配置说明：
//   - 数据库配置：从参数传入
//   - URL 前缀："admin"
//   - 文件存储：./uploads 目录，前缀为 "uploads"
//   - 语言：英语（EN）
//   - 首页 URL：/
//   - 调试模式：开启
//   - 主题：AdminLTE 黑色皮肤
//
// 使用场景：
//   - 自定义表的测试
//   - 集成测试
//   - 黑盒测试
//   - 框架适配验证
//
// 使用示例：
//
//	handler := NewHandler(config.DatabaseList{
//	    "default": {
//	        Host: "127.0.0.1",
//	        Port: "3306",
//	        User: "root",
//	        Pwd: "root",
//	        Name: "go-admin-test",
//	        Driver: config.DriverMysql,
//	    },
//	}, table.GeneratorList{})
//
// 注意事项：
//   - 端口 9087 需要未被占用
//   - 数据库配置需要正确
//   - 表生成器列表可以为空（使用默认表）
func NewHandler(dbs config.DatabaseList, gens table.GeneratorList) http.Handler {
	// 创建新的 Beego 应用实例
	app := beego.NewApp()

	// 创建 GoAdmin 默认引擎实例
	eng := engine.Default()
	// 创建管理插件，使用自定义表生成器
	adminPlugin := admin.NewAdmin(gens)

	// 配置引擎：添加配置、添加插件、使用 Beego 应用
	if err := eng.AddConfig(&config.Config{
		Databases: dbs,     // 数据库配置
		UrlPrefix: "admin", // URL 前缀
		Store: config.Store{ // 文件存储配置
			Path:   "./uploads", // 文件存储路径
			Prefix: "uploads",   // URL 前缀
		},
		Language:    language.EN,                   // 语言设置为英语
		IndexUrl:    "/",                           // 首页 URL
		Debug:       true,                          // 开启调试模式
		ColorScheme: adminlte.ColorschemeSkinBlack, // 使用 AdminLTE 黑色皮肤
	}).
		AddPlugins(adminPlugin).Use(app); err != nil {
		panic(err) // 配置失败时抛出异常
	}

	// 添加 Chart.js 图表组件到模板引擎
	template.AddComp(chartjs.NewChart())

	// 注册 HTML 路由：GET /admin
	eng.HTML("GET", "/admin", tables.GetContent)

	// 配置 Beego 监听地址和端口
	beego.BConfig.Listen.HTTPAddr = "127.0.0.1" // 监听地址
	beego.BConfig.Listen.HTTPPort = 9087        // 监听端口

	// 返回 Beego 应用的处理器
	return app.Handlers
}

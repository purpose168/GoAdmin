// chi.go - Chi 框架适配器文件
// 包名：chi
// 作者：GoAdmin 团队
// 创建日期：2020年
// 描述：本文件提供 GoAdmin 在 Chi Web 框架下的适配器实现
//       包含 internalHandler 和 NewHandler 两个核心函数
//
// 主要功能：
//   - internalHandler: 创建内部 HTTP 处理器，用于测试环境
//   - NewHandler: 创建新的 HTTP 处理器，支持自定义数据库和表生成器
//
// 核心概念：
//   - Chi: Go 语言的轻量级、可组合的 HTTP 路由器
//   - GoAdmin: Go 后台管理系统框架
//   - 适配器模式: 将 GoAdmin 集成到不同的 Web 框架中
//   - 插件系统: 通过插件扩展 GoAdmin 的功能
//
// 技术栈：
//   - Chi: Web 框架
//   - GoAdmin Engine: 核心引擎
//   - Admin Plugin: 管理插件
//   - Example Plugin: 示例插件
//   - Chart.js: 图表库
//   - AdminLTE: 主题模板
//
// 数据库支持：
//   - MySQL
//   - PostgreSQL
//   - SQLite
//   - MSSQL
//
// 使用场景：
//   - 集成测试
//   - 开发环境
//   - 演示环境
//   - 框架适配验证
//
// 配置说明：
//   - URL 前缀: /admin
//   - 存储路径: ./uploads
//   - 语言: 英语
//   - 主题: AdminLTE 黑色皮肤
//
// 注意事项：
//   - 需要正确配置数据库连接
//   - 需要确保上传目录存在且有写入权限
//   - 测试环境使用 JSON 配置文件
//   - 生产环境建议使用环境变量配置

package chi

import (
	"net/http" // Go 标准网络包，提供 HTTP 客户端和服务端功能
	"os"       // Go 标准操作系统包，提供操作系统接口

	// 导入 Chi 适配器，用于将 GoAdmin 集成到 Chi 框架
	// 使用空白导入（_）触发适配器的 init() 函数
	_ "github.com/purpose168/GoAdmin/adapter/chi"
	"github.com/purpose168/GoAdmin/modules/config"              // 配置模块包，提供配置管理功能
	"github.com/purpose168/GoAdmin/modules/language"            // 语言模块包，提供多语言支持
	"github.com/purpose168/GoAdmin/plugins/admin/modules/table" // 表格模块包，提供表格生成和管理功能

	// 导入 MySQL 数据库驱动，支持 MySQL 数据库连接
	// 使用空白导入（_）触发驱动的 init() 函数
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/mysql"

	// 导入 PostgreSQL 数据库驱动，支持 PostgreSQL 数据库连接
	// 使用空白导入（_）触发驱动的 init() 函数
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/postgres"

	// 导入 SQLite 数据库驱动，支持 SQLite 数据库连接
	// 使用空白导入（_）触发驱动的 init() 函数
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/sqlite"

	// 导入 MSSQL 数据库驱动，支持 Microsoft SQL Server 数据库连接
	// 使用空白导入（_）触发驱动的 init() 函数
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/mssql"

	"github.com/purpose168/GoAdmin-themes/adminlte" // AdminLTE 主题包，提供后台管理界面主题

	"github.com/go-chi/chi"                          // Chi Web 框架包，提供 HTTP 路由器功能
	"github.com/purpose168/GoAdmin/engine"           // GoAdmin 核心引擎包，提供框架核心功能
	"github.com/purpose168/GoAdmin/plugins/admin"    // 管理插件包，提供后台管理功能
	"github.com/purpose168/GoAdmin/plugins/example"  // 示例插件包，提供示例功能
	"github.com/purpose168/GoAdmin/template"         // 模板包，提供模板渲染功能
	"github.com/purpose168/GoAdmin/template/chartjs" // Chart.js 模板包，提供图表组件
	"github.com/purpose168/GoAdmin/tests/tables"     // 测试表格包，提供测试用的表格生成器
)

// ==================== 处理器函数 ====================

// internalHandler 创建内部HTTP处理器
// 返回:
//   - http.Handler: HTTP处理器，用于处理HTTP请求
//
// 说明：
//
//	该函数创建一个用于测试环境的内部 HTTP 处理器。
//	它使用命令行参数中的最后一个参数作为 JSON 配置文件路径。
//
// 功能特性：
//   - 创建 Chi HTTP 路由器实例
//   - 初始化 GoAdmin 默认引擎
//   - 添加管理插件和示例插件
//   - 添加用户表生成器
//   - 配置 Chart.js 图表组件
//   - 设置自定义 HTML 内容
//
// 配置说明：
//   - 配置来源: 命令行最后一个参数（JSON 文件）
//   - HTML 路由: GET /admin
//
// 技术细节：
//   - 使用 chi.NewRouter() 创建 Chi 路由器实例
//   - 使用 engine.Default() 创建默认引擎实例
//   - 使用 admin.NewAdmin() 创建管理插件实例
//   - 使用 example.NewExample() 创建示例插件实例
//   - 使用 eng.AddConfigFromJSON() 从 JSON 文件加载配置
//   - 使用 eng.AddPlugins() 添加插件
//   - 使用 eng.Use() 将引擎应用到 Chi 路由器
//   - 使用 template.AddComp() 添加模板组件
//   - 使用 eng.HTML() 注册 HTML 路由
//
// 使用场景：
//   - 集成测试
//   - 单元测试
//   - 开发环境
//   - 演示环境
//
// 注意事项：
//   - 需要在命令行参数中提供 JSON 配置文件路径
//   - 配置文件必须包含数据库连接信息
//   - 测试完成后需要清理测试数据
//
// 错误处理：
//   - 如果配置加载失败，会触发 panic
//   - 如果插件添加失败，会触发 panic
//   - 如果引擎应用失败，会触发 panic
//
// 示例：
//
//	// 创建内部处理器
//	handler := internalHandler()
//
//	// 使用处理器处理请求
//	resp := httptest.NewRecorder()
//	req := httptest.NewRequest("GET", "/admin", nil)
//	handler.ServeHTTP(resp, req)
func internalHandler() http.Handler {
	// 创建 Chi HTTP 路由器实例
	r := chi.NewRouter() // 创建 Chi 路由器实例

	// 创建 GoAdmin 默认引擎实例
	eng := engine.Default() // 创建默认引擎实例

	// 创建管理插件实例，传入测试表格生成器
	adminPlugin := admin.NewAdmin(tables.Generators) // 创建管理插件实例

	// 添加用户表生成器
	adminPlugin.AddGenerator("user", tables.GetUserTable) // 添加用户表生成器

	// 创建示例插件实例
	examplePlugin := example.NewExample() // 创建示例插件实例

	// 添加 Chart.js 图表组件到模板
	template.AddComp(chartjs.NewChart()) // 添加 Chart.js 图表组件

	// 从命令行最后一个参数加载 JSON 配置，添加插件，并将引擎应用到 Chi 路由器
	if err := eng.AddConfigFromJSON(os.Args[len(os.Args)-1]). // 从 JSON 文件加载配置
									AddPlugins(adminPlugin, examplePlugin). // 添加管理插件和示例插件
									Use(r); err != nil {                    // 将引擎应用到 Chi 路由器
		panic(err) // 如果失败，触发 panic
	}

	// 注册 HTML 路由，用于显示自定义内容
	eng.HTML("GET", "/admin", tables.GetContent) // 注册 HTML 路由

	// 返回 Chi 路由器实例
	return r // 返回路由器实例
}

// NewHandler 创建新的HTTP处理器
// 参数:
//   - dbs: 数据库连接配置列表，包含数据库连接信息
//   - gens: 表生成器列表，包含所有需要管理的表的定义
//
// 返回:
//   - http.Handler: HTTP处理器，用于处理HTTP请求
//
// 说明：
//
//	该函数创建一个支持自定义配置的 HTTP 处理器。
//	它允许用户自定义数据库连接和表生成器。
//
// 功能特性：
//   - 创建 Chi HTTP 路由器实例
//   - 初始化 GoAdmin 默认引擎
//   - 添加管理插件
//   - 支持自定义数据库配置
//   - 支持自定义表生成器
//   - 配置存储路径和前缀
//   - 配置语言和主题
//   - 添加 Chart.js 图表组件
//   - 设置自定义 HTML 内容
//
// 配置说明：
//   - URL 前缀: /admin
//   - 存储路径: ./uploads
//   - 存储前缀: uploads
//   - 语言: 英语
//   - 索引URL: /
//   - 调试模式: 启用
//   - 主题: AdminLTE 黑色皮肤
//   - HTML 路由: GET /admin
//
// 技术细节：
//   - 使用 chi.NewRouter() 创建 Chi 路由器实例
//   - 使用 engine.Default() 创建默认引擎实例
//   - 使用 admin.NewAdmin() 创建管理插件实例
//   - 使用 eng.AddConfig() 添加配置
//   - 使用 eng.AddPlugins() 添加插件
//   - 使用 eng.Use() 将引擎应用到 Chi 路由器
//   - 使用 template.AddComp() 添加模板组件
//   - 使用 eng.HTML() 注册 HTML 路由
//
// 使用场景：
//   - 生产环境
//   - 开发环境
//   - 自定义配置场景
//   - 多数据库场景
//
// 注意事项：
//   - 需要确保数据库连接信息正确
//   - 需要确保上传目录存在且有写入权限
//   - 表生成器需要正确定义
//
// 错误处理：
//   - 如果配置添加失败，会触发 panic
//   - 如果插件添加失败，会触发 panic
//   - 如果引擎应用失败，会触发 panic
//
// 示例：
//
//	// 配置数据库连接
//	dbs := config.DatabaseList{
//	    {
//	        Driver: "mysql",
//	        Host:   "localhost",
//	        Port:   "3306",
//	        User:   "root",
//	        Pass:   "password",
//	        Name:   "goadmin",
//	    },
//	}
//
//	// 配置表生成器
//	gens := table.GeneratorList{
//	    "user": tables.GetUserTable,
//	}
//
//	// 创建处理器
//	handler := NewHandler(dbs, gens)
//
//	// 使用处理器处理请求
//	resp := httptest.NewRecorder()
//	req := httptest.NewRequest("GET", "/admin", nil)
//	handler.ServeHTTP(resp, req)
func NewHandler(dbs config.DatabaseList, gens table.GeneratorList) http.Handler {
	// 创建 Chi HTTP 路由器实例
	r := chi.NewRouter() // 创建 Chi 路由器实例

	// 创建 GoAdmin 默认引擎实例
	eng := engine.Default() // 创建默认引擎实例

	// 创建管理插件实例，传入表生成器列表
	adminPlugin := admin.NewAdmin(gens) // 创建管理插件实例

	// 添加 Chart.js 图表组件到模板
	template.AddComp(chartjs.NewChart()) // 添加 Chart.js 图表组件

	// 添加配置，添加插件，并将引擎应用到 Chi 路由器
	if err := eng.AddConfig(&config.Config{ // 添加配置
		Databases: dbs,     // 设置数据库连接配置
		UrlPrefix: "admin", // 设置 URL 前缀为 /admin
		Store: config.Store{ // 配置文件存储
			Path:   "./uploads", // 设置上传文件存储路径
			Prefix: "uploads",   // 设置上传文件访问前缀
		},
		Language:    language.EN,                   // 设置语言为英语
		IndexUrl:    "/",                           // 设置索引 URL
		Debug:       true,                          // 启用调试模式
		ColorScheme: adminlte.ColorschemeSkinBlack, // 设置主题为 AdminLTE 黑色皮肤
	}).
		AddPlugins(adminPlugin). // 添加管理插件
		Use(r); err != nil {     // 将引擎应用到 Chi 路由器
		panic(err) // 如果失败，触发 panic
	}

	// 注册 HTML 路由，用于显示自定义内容
	eng.HTML("GET", "/admin", tables.GetContent) // 注册 HTML 路由

	// 返回 Chi 路由器实例
	return r // 返回路由器实例
}

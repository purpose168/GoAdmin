// Package gofiber 提供 GoAdmin 在 GoFiber Web 框架下的适配器实现
//
// 本包实现了 GoAdmin 与 GoFiber Web 框架的集成，提供以下功能：
//   - internalHandler: 创建用于测试环境的内部 HTTP 请求处理器
//   - NewHandler: 创建支持自定义配置的 HTTP 请求处理器
//
// 核心概念：
//   - GoFiber: 基于 Go 语言的高性能 HTTP Web 框架（基于 FastHTTP）
//   - GoAdmin: Go 后台管理框架，提供完整的后台管理功能
//   - 适配器模式: 将 GoAdmin 集成到不同 Web 框架的设计模式
//   - 插件系统: 通过插件扩展 GoAdmin 功能的架构设计
//   - FastHTTP: 高性能的 HTTP 服务器实现
//
// 技术栈：
//   - GoFiber: Web 框架（基于 FastHTTP）
//   - GoAdmin Engine: 核心引擎
//   - Admin Plugin: 管理插件
//   - Chart.js: 图表组件
//   - AdminLTE: 后台管理界面主题
//   - FastHTTP: 高性能 HTTP 服务器
//
// 数据库支持：
//   - MySQL: 开源关系型数据库
//   - PostgreSQL: 高级开源关系型数据库
//   - SQLite: 轻量级嵌入式数据库
//   - MSSQL: Microsoft SQL Server 数据库
//
// 使用场景：
//   - 集成测试: 测试 GoAdmin 与 GoFiber 框架的集成
//   - 开发环境: 快速搭建开发环境
//   - 演示环境: 展示 GoAdmin 功能
//   - 框架适配验证: 验证 GoFiber 框架适配器的正确性
//   - 高性能场景: 利用 FastHTTP 的高性能特性
//
// 配置说明：
//   - URL 前缀: 默认为 /admin
//   - 存储路径: 默认为 ./uploads
//   - 语言: 支持多语言，默认为英语
//   - 主题: 支持 AdminLTE 主题
//   - 服务器头部: "Fiber"
//
// 注意事项：
//   - 需要正确配置数据库连接信息
//   - 需要确保上传目录存在且有写入权限
//   - JSON 配置文件需要包含完整的配置信息
//   - 环境变量配置需要正确设置
//   - GoFiber 基于 FastHTTP，与标准 HTTP 包不完全兼容
//
// 作者: GoAdmin Team
// 创建日期: 2024-01-01
// 版本: 1.0.0
package gofiber

import (
	"os" // Go 标准操作系统包，提供操作系统接口

	"github.com/gofiber/fiber/v2"                               // GoFiber Web 框架包，提供 HTTP 服务器和路由功能
	"github.com/purpose168/GoAdmin-themes/adminlte"             // AdminLTE 主题包，提供后台管理界面主题
	"github.com/purpose168/GoAdmin/engine"                      // GoAdmin 核心引擎包，提供框架核心功能
	"github.com/purpose168/GoAdmin/modules/config"              // 配置模块包，提供配置管理功能
	"github.com/purpose168/GoAdmin/modules/language"            // 语言模块包，提供多语言支持
	"github.com/purpose168/GoAdmin/plugins/admin"               // 管理插件包，提供后台管理功能
	"github.com/purpose168/GoAdmin/plugins/admin/modules/table" // 表格模块包，提供表格生成和管理功能
	"github.com/purpose168/GoAdmin/template"                    // 模板包，提供模板渲染功能
	"github.com/purpose168/GoAdmin/template/chartjs"            // Chart.js 模板包，提供图表组件
	"github.com/purpose168/GoAdmin/tests/tables"                // 测试表格包，提供测试用的表格生成器
	"github.com/valyala/fasthttp"                               // FastHTTP 包，提供高性能 HTTP 服务器功能

	ada "github.com/purpose168/GoAdmin/adapter/gofiber"           // GoFiber 适配器包，提供 GoFiber 框架适配功能
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/mssql"    // 导入 MSSQL 数据库驱动，支持 Microsoft SQL Server 数据库连接
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/mysql"    // 导入 MySQL 数据库驱动，支持 MySQL 数据库连接
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/postgres" // 导入 PostgreSQL 数据库驱动，支持 PostgreSQL 数据库连接
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/sqlite"   // 导入 SQLite 数据库驱动，支持 SQLite 数据库连接
)

// internalHandler 创建内部HTTP请求处理器
//
// 返回:
//   - fasthttp.RequestHandler: FastHTTP请求处理器，用于处理HTTP请求
//
// 说明：
//
//	该函数创建一个用于测试环境的内部 HTTP 请求处理器。
//	它使用命令行参数中的最后一个参数作为 JSON 配置文件路径。
//
// 功能特性：
//   - 创建 GoFiber 应用实例（配置服务器头部为"Fiber"）
//   - 初始化 GoAdmin 默认引擎
//   - 添加管理插件并启用 XSS 过滤
//   - 添加用户表生成器
//   - 添加 Chart.js 图表组件
//   - 从 JSON 文件加载配置
//   - 设置自定义 HTML 内容
//
// 配置说明：
//   - 配置来源: 命令行最后一个参数（JSON 文件）
//   - HTML 路由: GET /admin
//   - 服务器头部: "Fiber"
//   - XSS 过滤: 启用
//
// 技术细节：
//   - 使用 fiber.New() 创建 GoFiber 应用实例，配置服务器头部为 "Fiber"
//   - 使用 engine.Default() 创建默认引擎实例
//   - 使用 admin.NewAdmin() 创建管理插件实例
//   - 使用 AddDisplayFilterXssJsFilter() 启用 XSS 过滤
//   - 使用 adminPlugin.AddGenerator() 添加用户表生成器
//   - 使用 eng.AddConfigFromJSON() 从 JSON 文件加载配置
//   - 使用 eng.AddPlugins() 添加插件
//   - 使用 eng.Use() 将引擎应用到 GoFiber 应用
//   - 使用 template.AddComp() 添加模板组件
//   - 使用 eng.HTML() 注册 HTML 路由
//   - 使用 app.Handler() 获取 FastHTTP 请求处理器
//
// 使用场景：
//   - 集成测试
//   - 单元测试
//   - 开发环境
//   - 演示环境
//   - 高性能场景
//
// 注意事项：
//   - 需要在命令行参数中提供 JSON 配置文件路径
//   - 配置文件必须包含数据库连接信息
//   - 测试完成后需要清理测试数据
//   - GoFiber 基于 FastHTTP，与标准 HTTP 包不完全兼容
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
//	ctx := &fasthttp.RequestCtx{}
//	ctx.Request.SetRequestURI("/admin")
//	handler(ctx)
func internalHandler() fasthttp.RequestHandler {
	app := fiber.New(fiber.Config{ // 创建 GoFiber 应用实例
		ServerHeader: "Fiber", // 配置服务器头部为 "Fiber"
	})

	eng := engine.Default() // 创建默认引擎实例

	adminPlugin := admin.NewAdmin(tables.Generators).AddDisplayFilterXssJsFilter() // 创建管理插件实例并启用 XSS 过滤
	adminPlugin.AddGenerator("user", tables.GetUserTable)                          // 添加用户表生成器

	template.AddComp(chartjs.NewChart()) // 添加 Chart.js 图表组件

	if err := eng.AddConfigFromJSON(os.Args[len(os.Args)-1]). // 从 JSON 文件加载配置
									AddPlugins(adminPlugin). // 添加管理插件
									Use(app); err != nil {   // 将引擎应用到 GoFiber 应用
		panic(err) // 如果失败，触发 panic
	}

	eng.HTML("GET", "/admin", tables.GetContent) // 注册 HTML 路由

	return app.Handler() // 返回 FastHTTP 请求处理器
}

// NewHandler 创建新的HTTP请求处理器
//
// 参数:
//   - dbs: 数据库连接配置列表，包含数据库连接信息
//   - gens: 表生成器列表，包含所有需要管理的表的定义
//
// 返回:
//   - fasthttp.RequestHandler: FastHTTP请求处理器，用于处理HTTP请求
//
// 说明：
//
//	该函数创建一个支持自定义配置的 HTTP 请求处理器。
//	它允许用户自定义数据库连接和表生成器。
//
// 功能特性：
//   - 创建 GoFiber 应用实例（配置服务器头部为"Fiber"）
//   - 初始化 GoAdmin 默认引擎
//   - 添加 Chart.js 图表组件
//   - 支持自定义数据库配置
//   - 支持自定义表生成器
//   - 配置存储路径和前缀
//   - 配置语言和主题
//   - 添加 GoFiber 适配器
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
//   - 服务器头部: "Fiber"
//
// 技术细节：
//   - 使用 fiber.New() 创建 GoFiber 应用实例，配置服务器头部为 "Fiber"
//   - 使用 engine.Default() 创建默认引擎实例
//   - 使用 eng.AddConfig() 添加配置
//   - 使用 eng.AddAdapter() 添加 GoFiber 适配器
//   - 使用 eng.AddGenerators() 添加表生成器
//   - 使用 eng.Use() 将引擎应用到 GoFiber 应用
//   - 使用 template.AddComp() 添加模板组件
//   - 使用 eng.HTML() 注册 HTML 路由
//   - 使用 app.Handler() 获取 FastHTTP 请求处理器
//
// 使用场景：
//   - 生产环境
//   - 开发环境
//   - 自定义配置场景
//   - 多数据库场景
//   - 高性能场景
//
// 注意事项：
//   - 需要确保数据库连接信息正确
//   - 需要确保上传目录存在且有写入权限
//   - 表生成器需要正确定义
//   - GoFiber 基于 FastHTTP，与标准 HTTP 包不完全兼容
//
// 错误处理：
//   - 如果配置添加失败，会触发 panic
//   - 如果适配器添加失败，会触发 panic
//   - 如果表生成器添加失败，会触发 panic
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
//	ctx := &fasthttp.RequestCtx{}
//	ctx.Request.SetRequestURI("/admin")
//	handler(ctx)
func NewHandler(dbs config.DatabaseList, gens table.GeneratorList) fasthttp.RequestHandler {
	app := fiber.New(fiber.Config{ // 创建 GoFiber 应用实例
		ServerHeader: "Fiber", // 配置服务器头部为 "Fiber"
	})

	eng := engine.Default() // 创建默认引擎实例

	template.AddComp(chartjs.NewChart()) // 添加 Chart.js 图表组件

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
		AddAdapter(new(ada.Gofiber)). // 添加 GoFiber 适配器
		AddGenerators(gens).          // 添加表生成器
		Use(app); err != nil {        // 将引擎应用到 GoFiber 应用
		panic(err) // 如果失败，触发 panic
	}

	eng.HTML("GET", "/admin", tables.GetContent) // 注册 HTML 路由

	return app.Handler() // 返回 FastHTTP 请求处理器
}

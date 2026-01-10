// Package gorilla 提供 GoAdmin 在 Gorilla Mux Web 框架下的适配器实现
//
// 本包实现了 GoAdmin 与 Gorilla Mux Web 框架的集成，提供以下功能：
//   - internalHandler: 创建用于测试环境的内部 HTTP 处理器
//   - NewHandler: 创建支持自定义配置的 HTTP 处理器
//
// 核心概念：
//   - Gorilla Mux: 基于 Go 语言的高性能 HTTP 路由器
//   - GoAdmin: Go 后台管理框架，提供完整的后台管理功能
//   - 适配器模式: 将 GoAdmin 集成到不同 Web 框架的设计模式
//   - 插件系统: 通过插件扩展 GoAdmin 功能的架构设计
//   - 示例插件: 展示 GoAdmin 插件系统的示例实现
//
// 技术栈：
//   - Gorilla Mux: Web 路由器
//   - GoAdmin Engine: 核心引擎
//   - Admin Plugin: 管理插件
//   - Example Plugin: 示例插件
//   - Chart.js: 图表组件
//   - AdminLTE: 后台管理界面主题
//
// 数据库支持：
//   - MySQL: 开源关系型数据库
//   - PostgreSQL: 高级开源关系型数据库
//   - SQLite: 轻量级嵌入式数据库
//   - MSSQL: Microsoft SQL Server 数据库
//
// 使用场景：
//   - 集成测试: 测试 GoAdmin 与 Gorilla Mux 框架的集成
//   - 开发环境: 快速搭建开发环境
//   - 演示环境: 展示 GoAdmin 功能
//   - 框架适配验证: 验证 Gorilla Mux 框架适配器的正确性
//   - 插件开发: 展示如何开发自定义插件
//
// 配置说明：
//   - URL 前缀: 默认为 /admin
//   - 存储路径: 默认为 ./uploads
//   - 语言: 支持多语言，默认为英语
//   - 主题: 支持 AdminLTE 主题
//
// 注意事项：
//   - 需要正确配置数据库连接信息
//   - 需要确保上传目录存在且有写入权限
//   - JSON 配置文件需要包含完整的配置信息
//   - 环境变量配置需要正确设置
//   - 示例插件仅用于演示，生产环境可根据需要移除
//
// 作者: GoAdmin Team
// 创建日期: 2024-01-01
// 版本: 1.0.0
package gorilla

import (
	"net/http" // Go 标准 HTTP 包，提供 HTTP 客户端和服务器功能
	"os"       // Go 标准操作系统包，提供操作系统接口

	"github.com/gorilla/mux"                                    // Gorilla Mux 路由器包，提供 HTTP 路由功能
	"github.com/purpose168/GoAdmin-themes/adminlte"             // AdminLTE 主题包，提供后台管理界面主题
	"github.com/purpose168/GoAdmin/engine"                      // GoAdmin 核心引擎包，提供框架核心功能
	"github.com/purpose168/GoAdmin/modules/config"              // 配置模块包，提供配置管理功能
	"github.com/purpose168/GoAdmin/modules/language"            // 语言模块包，提供多语言支持
	"github.com/purpose168/GoAdmin/plugins/admin"               // 管理插件包，提供后台管理功能
	"github.com/purpose168/GoAdmin/plugins/admin/modules/table" // 表格模块包，提供表格生成和管理功能
	"github.com/purpose168/GoAdmin/plugins/example"             // 示例插件包，展示插件系统的使用
	"github.com/purpose168/GoAdmin/template"                    // 模板包，提供模板渲染功能
	"github.com/purpose168/GoAdmin/template/chartjs"            // Chart.js 模板包，提供图表组件
	"github.com/purpose168/GoAdmin/tests/tables"                // 测试表格包，提供测试用的表格生成器

	_ "github.com/purpose168/GoAdmin/adapter/gorilla"             // 导入 Gorilla Mux 适配器，支持 Gorilla Mux 框架集成
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/mssql"    // 导入 MSSQL 数据库驱动，支持 Microsoft SQL Server 数据库连接
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/mysql"    // 导入 MySQL 数据库驱动，支持 MySQL 数据库连接
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/postgres" // 导入 PostgreSQL 数据库驱动，支持 PostgreSQL 数据库连接
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/sqlite"   // 导入 SQLite 数据库驱动，支持 SQLite 数据库连接
)

// internalHandler 创建内部HTTP处理器
//
// 返回:
//   - http.Handler: HTTP处理器，用于处理HTTP请求
//
// 说明：
//
//	该函数创建一个用于测试环境的内部 HTTP 处理器。
//	它使用命令行参数中的最后一个参数作为 JSON 配置文件路径。
//
// 功能特性：
//   - 创建 Gorilla Mux 路由器实例
//   - 初始化 GoAdmin 默认引擎
//   - 创建示例插件
//   - 添加 Chart.js 图表组件
//   - 从 JSON 文件加载配置
//   - 添加管理插件和用户表生成器
//   - 添加示例插件
//   - 设置自定义 HTML 内容
//
// 配置说明：
//   - 配置来源: 命令行最后一个参数（JSON 文件）
//   - HTML 路由: GET /admin
//
// 技术细节：
//   - 使用 mux.NewRouter() 创建 Gorilla Mux 路由器实例
//   - 使用 engine.Default() 创建默认引擎实例
//   - 使用 example.NewExample() 创建示例插件实例
//   - 使用 admin.NewAdmin() 创建管理插件实例
//   - 使用 adminPlugin.AddGenerator() 添加用户表生成器
//   - 使用 eng.AddConfigFromJSON() 从 JSON 文件加载配置
//   - 使用 eng.AddPlugins() 添加插件（管理插件和示例插件）
//   - 使用 eng.Use() 将引擎应用到 Gorilla Mux 路由器
//   - 使用 template.AddComp() 添加模板组件
//   - 使用 eng.HTML() 注册 HTML 路由
//
// 使用场景：
//   - 集成测试
//   - 单元测试
//   - 开发环境
//   - 演示环境
//   - 插件开发
//
// 注意事项：
//   - 需要在命令行参数中提供 JSON 配置文件路径
//   - 配置文件必须包含数据库连接信息
//   - 测试完成后需要清理测试数据
//   - 示例插件仅用于演示，生产环境可根据需要移除
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
//	req, _ := http.NewRequest("GET", "/admin", nil)
//	resp := httptest.NewRecorder()
//	handler.ServeHTTP(resp, req)
func internalHandler() http.Handler {
	app := mux.NewRouter() // 创建 Gorilla Mux 路由器实例

	eng := engine.Default() // 创建默认引擎实例

	examplePlugin := example.NewExample() // 创建示例插件实例
	template.AddComp(chartjs.NewChart())  // 添加 Chart.js 图表组件

	if err := eng.AddConfigFromJSON(os.Args[len(os.Args)-1]). // 从 JSON 文件加载配置
									AddPlugins(admin.NewAdmin(tables.Generators). // 创建管理插件实例
															AddGenerator("user", tables.GetUserTable), examplePlugin). // 添加用户表生成器并添加示例插件
									Use(app); err != nil {                        // 将引擎应用到 Gorilla Mux 路由器
		panic(err) // 如果失败，触发 panic
	}

	eng.HTML("GET", "/admin", tables.GetContent) // 注册 HTML 路由

	return app // 返回路由器实例
}

// NewHandler 创建新的HTTP处理器
//
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
//   - 创建 Gorilla Mux 路由器实例
//   - 初始化 GoAdmin 默认引擎
//   - 添加 Chart.js 图表组件
//   - 支持自定义数据库配置
//   - 支持自定义表生成器
//   - 配置存储路径和前缀
//   - 配置语言和主题
//   - 添加管理插件
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
//   - 使用 mux.NewRouter() 创建 Gorilla Mux 路由器实例
//   - 使用 engine.Default() 创建默认引擎实例
//   - 使用 admin.NewAdmin() 创建管理插件实例
//   - 使用 eng.AddConfig() 添加配置
//   - 使用 eng.AddPlugins() 添加插件
//   - 使用 eng.Use() 将引擎应用到 Gorilla Mux 路由器
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
//	req, _ := http.NewRequest("GET", "/admin", nil)
//	resp := httptest.NewRecorder()
//	handler.ServeHTTP(resp, req)
func NewHandler(dbs config.DatabaseList, gens table.GeneratorList) http.Handler {
	app := mux.NewRouter() // 创建 Gorilla Mux 路由器实例

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
		AddPlugins(admin.NewAdmin(gens)). // 添加管理插件
		Use(app); err != nil {            // 将引擎应用到 Gorilla Mux 路由器
		panic(err) // 如果失败，触发 panic
	}

	eng.HTML("GET", "/admin", tables.GetContent) // 注册 HTML 路由

	return app // 返回路由器实例
}

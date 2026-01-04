// GoAdmin Buffalo 框架集成示例
// 包名称: main
// 作者: GoAdmin Team
// 创建日期: 2024
// 目的: 演示如何将 GoAdmin 后台管理系统集成到 Buffalo Web 框架中
//
// 本示例展示了以下功能:
// - Buffalo 应用实例的创建和配置
// - GoAdmin 引擎的初始化和配置
// - MySQL 数据库连接配置
// - 插件系统的使用
// - 自定义页面和路由
// - 优雅的服务关闭处理
//
// 使用说明:
// 1. 确保 MySQL 数据库已启动并创建了名为 "godmin" 的数据库
// 2. 运行此程序: go run main.go
// 3. 访问 http://127.0.0.1:9033/admin 进入管理后台
// 4. 默认登录凭证请参考 GoAdmin 文档
//
// 依赖项:
// - Buffalo Web 框架
// - GoAdmin 后台管理系统
// - MySQL 数据库驱动
// - AdminLTE 主题
// - Chart.js 图表库

package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	// 导入 Buffalo 适配器，用于将 GoAdmin 集成到 Buffalo 框架中
	// 使用下划线导入表示只执行包的 init() 函数，不直接使用包中的导出符号
	_ "github.com/purpose168/GoAdmin/adapter/buffalo"

	// 导入 MySQL 数据库驱动
	// GoAdmin 支持多种数据库驱动，此处使用 MySQL
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/mysql"

	// Buffalo Web 框架核心包
	// Buffalo 是一个全栈 Web 框架，提供了路由、中间件、模板等功能
	"github.com/gobuffalo/buffalo"

	// AdminLTE 主题包
	// AdminLTE 是一个流行的 Bootstrap 管理后台主题
	"github.com/purpose168/GoAdmin-themes/adminlte"

	// GoAdmin 引擎核心包
	// 提供了后台管理系统的核心功能，包括认证、权限、数据管理等
	"github.com/purpose168/GoAdmin/engine"

	// 数据模型示例包
	// 包含了示例数据表定义和生成器
	"github.com/purpose168/GoAdmin/examples/datamodel"

	// 配置模块包
	// 提供了配置结构和数据库连接管理
	"github.com/purpose168/GoAdmin/modules/config"

	// 语言模块包
	// 提供了多语言支持
	"github.com/purpose168/GoAdmin/modules/language"

	// 示例插件包
	// 演示如何创建和使用自定义插件
	"github.com/purpose168/GoAdmin/plugins/example"

	// 模板模块包
	// 提供了模板组件和渲染功能
	"github.com/purpose168/GoAdmin/template"

	// Chart.js 图表库集成包
	// 用于在管理后台中显示数据可视化图表
	"github.com/purpose168/GoAdmin/template/chartjs"
)

// main 函数是程序的入口点
// 执行流程:
// 1. 创建 Buffalo 应用实例
// 2. 初始化 GoAdmin 引擎
// 3. 配置数据库连接
// 4. 添加插件和生成器
// 5. 启动 Web 服务器
// 6. 监听系统信号实现优雅关闭
func main() {
	// 创建 Buffalo 应用实例
	// Buffalo.Options 结构体用于配置应用的启动参数
	// Env: 运行环境，"test" 表示测试环境
	// Addr: 监听地址和端口，格式为 "host:port"
	bu := buffalo.New(buffalo.Options{
		Env:  "test",
		Addr: "127.0.0.1:9033",
	})

	// 创建 GoAdmin 引擎实例
	// Default() 函数返回一个配置好的引擎实例，包含默认的中间件和处理器
	// 引擎负责处理管理后台的所有请求，包括认证、授权、数据管理等
	eng := engine.Default()

	// 配置 GoAdmin 引擎
	// config.Config 结构体包含了 GoAdmin 的所有配置选项
	cfg := config.Config{
		// 运行环境
		// EnvLocal 表示本地开发环境，其他选项包括 EnvProduction、EnvDevelopment
		Env: config.EnvLocal,

		// 数据库配置列表
		// GoAdmin 支持配置多个数据库连接，使用字符串键标识不同的连接
		// "default" 是默认数据库连接的键名
		Databases: config.DatabaseList{
			"default": {
				// 数据库主机地址
				Host: "127.0.0.1",

				// 数据库端口号
				// MySQL 默认端口为 3306
				Port: "3306",

				// 数据库用户名
				// 建议在生产环境中使用具有最小权限的专用账户
				User: "root",

				// 数据库密码
				// 建议在生产环境中使用环境变量或配置管理工具存储密码
				Pwd: "root",

				// 数据库名称
				// 需要在 MySQL 中预先创建该数据库
				Name: "godmin",

				// 最大空闲连接数
				// 连接池中保持的空闲连接的最大数量
				// 设置为 50 可以在高并发情况下减少连接建立的开销
				MaxIdleConns: 50,

				// 最大打开连接数
				// 连接池中允许的最大活动连接数
				// 设置为 150 可以防止数据库连接数过多导致资源耗尽
				MaxOpenConns: 150,

				// 连接最大生存时间
				// 连接在池中的最大存活时间，超过此时间将被关闭
				// 设置为 1 小时可以防止长时间持有的连接出现问题
				ConnMaxLifetime: time.Hour,

				// 数据库驱动类型
				// 使用 MySQL 驱动，GoAdmin 还支持 PostgreSQL、SQLite、MSSQL 等
				Driver: config.DriverMysql,
			},
		},

		// URL 前缀
		// 管理后台的所有路由都会以此前缀开头
		// 例如: http://127.0.0.1:9033/admin/info/user
		UrlPrefix: "admin",

		// 文件存储配置
		// 用于配置上传文件的存储路径和访问 URL
		Store: config.Store{
			// 文件存储的本地路径
			// 所有上传的文件将保存在此目录下
			Path: "./uploads",

			// URL 前缀
			// 访问上传文件时的 URL 路径前缀
			// 例如: http://127.0.0.1:9033/uploads/filename.jpg
			Prefix: "uploads",
		},

		// 界面语言
		// 设置管理后台的默认语言
		// GoAdmin 支持多种语言，包括中文、英文、日文等
		Language: language.EN,

		// 首页 URL
		// 用户登录后跳转的页面
		IndexUrl: "/",

		// 调试模式
		// 设置为 true 时会显示详细的错误信息和调试日志
		// 生产环境应设置为 false 以提高安全性和性能
		Debug: true,

		// 主题配色方案
		// AdminLTE 提供多种配色方案，如 SkinBlack、SkinBlue、SkinGreen 等
		ColorScheme: adminlte.ColorschemeSkinBlack,
	}

	// 添加 Chart.js 图表组件到模板系统
	// Chart.js 是一个流行的 JavaScript 图表库，支持多种图表类型
	// 可以在管理后台中显示数据可视化图表，如折线图、柱状图、饼图等
	template.AddComp(chartjs.NewChart())

	// 创建示例插件实例
	// 插件是 GoAdmin 的扩展机制，允许开发者添加自定义功能
	// 示例插件演示了如何创建和使用自定义插件
	examplePlugin := example.NewExample()

	// 从 Go 插件加载插件（可选功能）
	// Go 插件允许在运行时动态加载编译好的插件
	// 这对于模块化开发和插件分发非常有用
	//
	// 示例代码:
	// examplePlugin := plugins.LoadFromPlugin("../datamodel/example.so")

	// 自定义登录页面（可选功能）
	// GoAdmin 允许开发者自定义登录页面的外观和功能
	// 可以参考官方示例: https://github.com/purpose168/demo.go-admin.cn/blob/master/main.go#L39
	//
	// 示例代码:
	// template.AddComp("login", datamodel.LoginPage)

	// 从 JSON 文件加载配置（可选功能）
	// 可以将配置保存在 JSON 文件中，便于管理和版本控制
	// 适合在复杂项目中使用，可以将不同环境的配置分开管理
	//
	// 示例代码:
	// eng.AddConfigFromJSON("../datamodel/config.json")

	// 配置 GoAdmin 引擎并集成到 Buffalo 应用
	// 使用链式调用方法配置引擎，每个方法返回引擎实例以便继续调用
	//
	// AddConfig: 添加配置
	// AddGenerators: 添加数据表生成器
	// AddDisplayFilterXssJsFilter: 添加 XSS 过滤器，防止跨站脚本攻击
	// AddGenerator: 添加单个数据表生成器
	// AddPlugins: 添加插件
	// Use: 将引擎集成到 Buffalo 应用
	if err := eng.AddConfig(&cfg).
		AddGenerators(datamodel.Generators).
		AddDisplayFilterXssJsFilter().
		// 添加生成器，第一个参数是数据表的 URL 前缀
		// 示例:
		//
		// "user" => http://localhost:9033/admin/info/user
		//
		// 当访问 /admin/info/user 时，会显示用户数据表的管理界面
		AddGenerator("user", datamodel.GetUserTable).
		AddPlugins(examplePlugin).
		Use(bu); err != nil {
		// 如果配置失败，直接 panic 终止程序
		// 在生产环境中，应该记录错误日志并进行优雅的退出处理
		panic(err)
	}

	// 配置静态文件服务
	// ServeFiles 方法将本地目录映射到 URL 路径
	// 第一个参数是 URL 路径前缀
	// 第二个参数是本地文件系统目录
	// 这样用户可以通过 URL 访问上传的文件
	bu.ServeFiles("/uploads", http.Dir("./uploads"))

	// 自定义页面路由
	// eng.HTML 方法用于添加自定义 HTML 页面
	// 第一个参数是 HTTP 方法，如 "GET"、"POST"
	// 第二个参数是 URL 路径
	// 第三个参数是处理函数，返回页面的 HTML 内容
	// 这样可以创建完全自定义的页面，不受 GoAdmin 默认模板的限制
	eng.HTML("GET", "/admin", datamodel.GetContent)

	// 在新的 goroutine 中启动 Web 服务器
	// 使用 goroutine 可以让服务器在后台运行，主线程可以继续执行其他操作
	// 这是 Go 语言并发编程的典型应用场景
	// bu.Serve() 会阻塞当前 goroutine，直到服务器停止
	go func() {
		_ = bu.Serve()
	}()

	// 创建信号通道用于优雅关闭
	// make 函数创建一个带缓冲的通道，缓冲区大小为 1
	// os.Signal 类型用于接收操作系统信号
	quit := make(chan os.Signal, 1)

	// 注册要监听的系统信号
	// os.Interrupt 对应 Ctrl+C 信号
	// 当用户按下 Ctrl+C 时，信号会被发送到 quit 通道
	// signal.Notify 函数会将指定的信号转发到通道
	signal.Notify(quit, os.Interrupt)

	// 阻塞等待信号
	// <-quit 操作会阻塞，直到 quit 通道收到信号
	// 这样程序会一直运行，直到用户主动终止
	<-quit

	// 收到终止信号后，执行清理操作
	// 打印日志提示用户程序正在关闭
	log.Print("closing database connection")

	// 关闭数据库连接
	// 释放数据库资源，确保所有连接都被正确关闭
	// 这是优雅关闭的重要步骤，可以防止数据丢失或连接泄漏
	eng.MysqlConnection().Close()
}

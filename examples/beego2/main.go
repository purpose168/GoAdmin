// 包名：main
// 作者：GoAdmin 翻译团队
// 创建日期：2026-01-04
// 目的：这是一个使用 Beego2 框架的 GoAdmin 示例程序，演示如何集成 GoAdmin 管理后台到 Beego2 应用中
//
// 使用说明：
// 1. 确保已安装 MySQL 数据库并创建了名为 "godmin" 的数据库
// 2. 修改数据库连接配置（主机、端口、用户名、密码等）以匹配您的环境
// 3. 运行程序：go run main.go
// 4. 访问管理后台：http://127.0.0.1:9087/admin
//
// 功能特性：
// - 集成 Beego2 框架
// - 支持 MySQL 数据库连接
// - 使用 AdminLTE 主题
// - 包含示例插件和数据生成器
// - 支持 XSS 过滤
// - 支持文件上传功能
// - 优雅关闭处理

package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	// 导入 Beego2 适配器，用于将 GoAdmin 集成到 Beego2 框架中
	_ "github.com/purpose168/GoAdmin/adapter/beego2"
	// 导入 MySQL 驱动，用于数据库连接
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/mysql"

	// Beego2 Web 框架
	"github.com/beego/beego/v2/server/web"
	// AdminLTE 主题，用于管理后台界面
	"github.com/purpose168/GoAdmin-themes/adminlte"
	// GoAdmin 核心引擎
	"github.com/purpose168/GoAdmin/engine"
	// 数据模型示例，包含数据生成器和表结构定义
	"github.com/purpose168/GoAdmin/examples/datamodel"
	// 配置模块，用于管理应用配置
	"github.com/purpose168/GoAdmin/modules/config"
	// 语言模块，用于国际化支持
	"github.com/purpose168/GoAdmin/modules/language"
	// 示例插件，演示如何创建自定义插件
	"github.com/purpose168/GoAdmin/plugins/example"
	// 模板模块，用于页面渲染
	"github.com/purpose168/GoAdmin/template"
	// Chart.js 图表组件，用于数据可视化
	"github.com/purpose168/GoAdmin/template/chartjs"
)

// 主函数：程序入口点
// 功能：初始化 Beego2 应用，配置 GoAdmin 引擎，启动 HTTP 服务器
func main() {
	// 创建 Beego2 HTTP 服务器实例
	// web.NewHttpSever() 返回一个新的 Beego 应用对象
	应用 := web.NewHttpSever()

	// 创建 GoAdmin 默认引擎实例
	// engine.Default() 返回一个预配置的 GoAdmin 引擎，包含默认的中间件和处理器
	引擎 := engine.Default()

	// 配置 GoAdmin 引擎
	// config.Config 结构体包含了 GoAdmin 的所有配置选项
	配置 := config.Config{
		// 运行环境：本地开发环境
		// 可选值：config.EnvLocal（本地）、config.EnvProd（生产）
		Env: config.EnvLocal,
		// 数据库连接配置
		// config.DatabaseList 是一个 map 类型，可以配置多个数据库连接
		Databases: config.DatabaseList{
			"default": {
				// 数据库主机地址
				Host: "127.0.0.1",
				// 数据库端口
				Port: "3306",
				// 数据库用户名
				User: "root",
				// 数据库密码
				Pwd: "123456",
				// 数据库名称
				Name: "godmin",
				// 最大空闲连接数：连接池中保持的最大空闲连接数
				// 建议值：根据应用并发量设置，通常为 50-100
				MaxIdleConns: 50,
				// 最大打开连接数：连接池中允许的最大连接数
				// 建议值：根据数据库服务器性能设置，通常为 100-200
				MaxOpenConns: 150,
				// 连接最大生存时间：连接在连接池中存活的最长时间
				// 建议值：time.Hour（1小时）或 time.Minute*30（30分钟）
				ConnMaxLifetime: time.Hour,
				// 数据库驱动类型
				// 可选值：config.DriverMysql、config.DriverPostgresql、config.DriverSqlite 等
				Driver: config.DriverMysql,
			},
		},
		// 文件存储配置
		// 用于管理上传文件的存储路径和 URL 前缀
		Store: config.Store{
			// 文件存储路径：相对于项目根目录的路径
			Path: "./uploads",
			// URL 前缀：访问上传文件时的 URL 路径前缀
			Prefix: "uploads",
		},
		// 管理后台 URL 前缀
		// 访问管理后台的基础路径，例如：http://127.0.0.1:9087/admin
		UrlPrefix: "admin",
		// 首页 URL
		// 管理后台的首页路径
		IndexUrl: "/",
		// 调试模式
		// true：启用调试模式，显示详细错误信息
		// false：生产环境，隐藏详细错误信息
		Debug: true,
		// 语言设置
		// 可选值：language.CN（中文）、language.EN（英文）等
		Language: language.CN,
		// 主题配色方案
		// 使用 AdminLTE 主题的黑色皮肤
		// 可选值：adminlte.ColorschemeSkinBlack、adminlte.ColorschemeSkinBlue 等
		ColorScheme: adminlte.ColorschemeSkinBlack,
	}

	// 添加 Chart.js 图表组件到模板
	// chartjs.NewChart() 创建一个新的图表组件实例
	// template.AddComp() 将组件注册到模板系统中，使其可以在页面中使用
	template.AddComp(chartjs.NewChart())

	// 创建示例插件实例
	// example.NewExample() 返回一个预配置的示例插件
	示例插件 := example.NewExample()

	// 设置静态文件路径
	// 将 /uploads URL 路径映射到 uploads 目录
	// 这样可以通过 http://127.0.0.1:9087/uploads/xxx 访问上传的文件
	web.SetStaticPath("/uploads", "uploads")

	// 配置 GoAdmin 引擎并集成到 Beego2 应用中
	// 使用链式调用方法进行配置：
	// 1. AddConfig(&配置)：添加配置信息
	// 2. AddGenerators(datamodel.Generators)：添加数据生成器，用于生成 CRUD 页面
	// 3. AddDisplayFilterXssJsFilter()：添加 XSS 过滤器，防止跨站脚本攻击
	// 4. AddGenerator("user", datamodel.GetUserTable)：添加用户表生成器
	// 5. AddPlugins(示例插件)：添加自定义插件
	// 6. Use(应用)：将 GoAdmin 集成到 Beego2 应用中
	if 错误 := 引擎.AddConfig(&配置).
		AddGenerators(datamodel.Generators).
		AddDisplayFilterXssJsFilter().
		AddGenerator("user", datamodel.GetUserTable).
		AddPlugins(示例插件).
		Use(应用); 错误 != nil {
		// 如果配置失败，抛出 panic 并终止程序
		panic(错误)
	}

	// 注册 HTML 路由
	// 引擎.HTML() 方法用于注册自定义 HTML 页面路由
	// 参数说明：
	// - "GET"：HTTP 方法
	// - "/admin"：路由路径
	// - datamodel.GetContent：处理函数，返回页面内容
	引擎.HTML("GET", "/admin", datamodel.GetContent)

	// 配置 HTTP 服务器监听地址和端口
	// HTTPSAddr：HTTPS 监听地址（本例中未启用 HTTPS）
	应用.Cfg.Listen.HTTPSAddr = "127.0.0.1"
	// HTTPPort：HTTP 监听端口
	应用.Cfg.Listen.HTTPPort = 9087

	// 启动 Beego2 HTTP 服务器
	// 使用 goroutine 异步启动，不阻塞主线程
	// app.Run("") 启动服务器，空字符串表示使用配置文件中的配置
	go 应用.Run("")

	// 创建信号通道，用于优雅关闭
	// make(chan os.Signal, 1) 创建一个缓冲大小为 1 的信号通道
	退出通道 := make(chan os.Signal, 1)
	// 注册要监听的信号
	// signal.Notify() 将信号发送到退出通道
	// os.Interrupt：中断信号（Ctrl+C）
	signal.Notify(退出通道, os.Interrupt)
	// 阻塞等待信号
	// <-退出通道 会一直阻塞，直到收到信号
	<-退出通道

	// 打印关闭信息
	log.Print("closing database connection")
	// 关闭数据库连接
	// 引擎.MysqlConnection().Close() 关闭 MySQL 数据库连接
	// 优雅关闭时应该先关闭数据库连接，确保所有事务都已完成
	引擎.MysqlConnection().Close()
}

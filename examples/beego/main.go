// GoAdmin Beego 示例程序
// 本示例演示如何使用 GoAdmin 框架与 Beego Web 框架集成
// 作者: GoAdmin 团队
// 创建日期: 2024年
// 目的: 展示 GoAdmin 在 Beego 框架中的基本用法和配置
//
// 使用说明:
// 1. 确保已安装 MySQL 数据库，并创建了名为 "godmin" 的数据库
// 2. 修改数据库连接配置（Host、Port、User、Pwd）以匹配您的环境
// 3. 运行程序: go run main.go
// 4. 访问管理后台: http://127.0.0.1:9087/admin
//
// 端口占用处理:
// 如果端口 9087 被占用，请执行以下步骤:
// 1. 检查端口占用: lsof -i :9087 (Linux/Mac) 或 netstat -ano | findstr :9087 (Windows)
// 2. 终止占用进程: kill -9 <PID> (Linux/Mac) 或 taskkill /PID <PID> /F (Windows)
// 3. 修改 HTTPPort 为其他可用端口，然后重新启动服务

package main // 声明主包，Go 程序的入口

import (
	"log"       // 标准日志库，用于记录日志信息
	"os"        // 操作系统接口，用于获取系统级功能
	"os/signal" // 信号处理，用于捕获系统信号（如中断信号）
	"time"      // 时间处理库，用于设置数据库连接最大生命周期

	// GoAdmin Beego 适配器，提供 Beego 框架的集成支持
	_ "github.com/purpose168/GoAdmin/adapter/beego"
	// GoAdmin MySQL 数据库驱动，用于连接和操作 MySQL 数据库
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/mysql"

	"github.com/astaxie/beego"                         // Beego Web 框架核心包
	"github.com/purpose168/GoAdmin-themes/adminlte"    // AdminLTE 主题包，提供管理后台的 UI 主题
	"github.com/purpose168/GoAdmin/engine"             // GoAdmin 核心引擎，提供主要功能
	"github.com/purpose168/GoAdmin/examples/datamodel" // 数据模型示例，展示如何定义数据表
	"github.com/purpose168/GoAdmin/modules/config"     // 配置模块，用于管理应用配置
	"github.com/purpose168/GoAdmin/modules/language"   // 语言模块，提供多语言支持
	"github.com/purpose168/GoAdmin/plugins/example"    // 示例插件，展示插件系统的使用
	"github.com/purpose168/GoAdmin/template"           // 模板引擎，用于渲染页面
	"github.com/purpose168/GoAdmin/template/chartjs"   // Chart.js 图表组件，用于数据可视化
)

func main() {
	// 创建 Beego 应用实例
	app := beego.NewApp()

	// 创建 GoAdmin 默认引擎实例
	eng := engine.Default()

	// 配置 GoAdmin 应用
	cfg := config.Config{
		// 设置运行环境为本地开发环境
		Env: config.EnvLocal,
		// 配置数据库连接信息
		Databases: config.DatabaseList{
			"default": {
				// 数据库主机地址
				Host: "127.0.0.1",
				// 数据库端口
				Port: "3306",
				// 数据库用户名
				User: "root",
				// 数据库密码
				Pwd: "root",
				// 数据库名称
				Name: "godmin",
				// 最大空闲连接数
				MaxIdleConns: 50,
				// 最大打开连接数
				MaxOpenConns: 150,
				// 连接最大生命周期为1小时
				ConnMaxLifetime: time.Hour,
				// 使用 MySQL 驱动
				Driver: config.DriverMysql,
			},
		},
		// 配置文件存储路径
		Store: config.Store{
			// 上传文件存储路径
			Path: "./uploads",
			// URL 前缀
			Prefix: "uploads",
		},
		// 管理后台 URL 前缀
		UrlPrefix: "admin",
		// 首页 URL
		IndexUrl: "/",
		// 开启调试模式
		Debug: true,
		// 设置语言为中文
		Language: language.CN,
		// 设置主题颜色方案为黑色皮肤
		ColorScheme: adminlte.ColorschemeSkinBlack,
	}

	// 添加 Chart.js 图表组件到模板
	template.AddComp(chartjs.NewChart())

	// 创建示例插件实例
	examplePlugin := example.NewExample()

	// 配置并启动 GoAdmin 引擎
	if err := eng.AddConfig(&cfg). // 添加配置
					AddGenerators(datamodel.Generators).          // 添加数据生成器
					AddDisplayFilterXssJsFilter().                // 添加 XSS 和 JS 过滤器
					AddGenerator("user", datamodel.GetUserTable). // 添加用户表生成器
					AddPlugins(examplePlugin).                    // 添加示例插件
					Use(app); err != nil {                        // 使用 Beego 应用
		panic(err)
	}

	// 设置 HTML 路由，处理 GET 请求到 /admin 路径
	eng.HTML("GET", "/admin", datamodel.GetContent)

	// 配置 Beego 服务器监听地址
	beego.BConfig.Listen.HTTPAddr = "127.0.0.1"
	// 配置 Beego 服务器监听端口
	beego.BConfig.Listen.HTTPPort = 9087
	// 在协程中启动 Beego 应用
	go app.Run()

	// 创建一个信号通道，用于捕获中断信号
	quit := make(chan os.Signal, 1)
	// 监听系统中断信号（Ctrl+C）
	signal.Notify(quit, os.Interrupt)
	// 阻塞等待中断信号
	<-quit
	// 打印关闭数据库连接的日志
	log.Print("closing database connection")
	// 关闭 MySQL 数据库连接
	eng.MysqlConnection().Close()
}

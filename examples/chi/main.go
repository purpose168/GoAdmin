// GoAdmin Chi 框架示例程序
// 包名称: main
// 作者: GoAdmin Group
// 创建日期: 2026-01-04
// 目的: 演示如何使用 GoAdmin 框架配合 Chi 路由器构建后台管理系统
// 本示例展示了 GoAdmin 的核心功能，包括配置管理、数据表生成、插件系统和静态文件服务

package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	// 导入 Chi 适配器，用于将 GoAdmin 与 Chi 路由器集成
	// 使用空白导入符 _ 表示只执行包的 init 函数，不直接使用包中的导出标识符
	_ "github.com/purpose168/GoAdmin/adapter/chi"
	// 导入 MySQL 数据库驱动
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/mysql"

	"github.com/go-chi/chi"
	"github.com/purpose168/GoAdmin-themes/adminlte"
	"github.com/purpose168/GoAdmin/engine"
	"github.com/purpose168/GoAdmin/examples/datamodel"
	"github.com/purpose168/GoAdmin/modules/config"
	"github.com/purpose168/GoAdmin/modules/language"
	"github.com/purpose168/GoAdmin/plugins/example"
	"github.com/purpose168/GoAdmin/template"
	"github.com/purpose168/GoAdmin/template/chartjs"
)

// main 函数是程序的入口点
// 在 Go 语言中，main 包的 main 函数是程序执行的起点
// 本函数演示了如何初始化 GoAdmin 引擎、配置数据库、添加插件和启动 HTTP 服务器
func main() {
	// 创建 Chi 路由器实例
	// Chi 是一个轻量级、符合 Go 标准库风格的 HTTP 路由器
	// 它提供了路由组合、中间件支持和优雅的 API 设计
	r := chi.NewRouter()

	// 创建默认的 GoAdmin 引擎实例
	// Engine 是 GoAdmin 的核心组件，负责管理配置、数据库连接、插件和路由
	// Default() 函数返回一个预配置的引擎实例，包含默认的设置
	eng := engine.Default()

	// 配置 GoAdmin 的运行参数
	// Config 结构体包含了 GoAdmin 运行所需的所有配置信息
	cfg := config.Config{
		// Env 设置运行环境
		// EnvLocal 表示本地开发环境，其他选项包括 EnvProduction（生产环境）和 EnvDev（开发环境）
		// 不同环境会影响日志级别、错误显示等行为
		Env: config.EnvLocal,
		// Databases 配置数据库连接信息
		// DatabaseList 类型是一个 map，键是数据库连接名称，值是数据库配置
		// 这里配置了一个名为 "default" 的 MySQL 数据库连接
		Databases: config.DatabaseList{
			"default": {
				// Host 数据库服务器地址
				// 127.0.0.1 表示本地主机，也可以使用 localhost
				Host: "127.0.0.1",
				// Port 数据库服务器端口
				// MySQL 默认端口是 3306
				Port: "3306",
				// User 数据库用户名
				// root 是 MySQL 的默认超级用户
				User: "root",
				// Pwd 数据库密码
				// 在生产环境中，应该使用环境变量或配置文件存储敏感信息
				Pwd: "root",
				// Name 数据库名称
				// godmin 是示例使用的数据库名
				Name: "godmin",
				// MaxIdleConns 连接池中最大空闲连接数
				// 空闲连接是指当前未被使用的连接，保留一定数量的空闲连接可以提高性能
				// 设置为 50 表示最多保留 50 个空闲连接
				MaxIdleConns: 50,
				// MaxOpenConns 连接池中最大打开连接数
				// 打开连接包括正在使用的连接和空闲连接
				// 设置为 150 表示最多同时打开 150 个数据库连接
				MaxOpenConns: 150,
				// ConnMaxLifetime 连接的最大生命周期
				// 超过此时间的连接将被关闭并从连接池中移除
				// time.Hour 表示连接最多存活 1 小时
				ConnMaxLifetime: time.Hour,
				// Driver 数据库驱动类型
				// DriverMysql 表示使用 MySQL 数据库驱动
				Driver: config.DriverMysql,
			},
		},
		// UrlPrefix 设置后台管理系统的 URL 前缀
		// 所有后台管理相关的路由都会以此前缀开头
		// 例如：http://localhost:3333/admin/info/user
		UrlPrefix: "admin",
		// Store 配置文件存储设置
		// 用于存储用户上传的文件
		Store: config.Store{
			// Path 文件存储的本地路径
			// ./uploads 表示当前目录下的 uploads 文件夹
			Path: "./uploads",
			// Prefix 文件访问的 URL 前缀
			// uploads 表示文件将通过 /uploads 路径访问
			Prefix: "uploads",
		},
		// Language 设置后台管理系统的语言
		// language.EN 表示使用英语
		// 其他选项包括 language.ZH（中文）、language.JA（日语）等
		Language: language.EN,
		// IndexUrl 设置首页 URL
		// 用户访问后台管理系统时的默认页面
		IndexUrl: "/",
		// Debug 是否启用调试模式
		// 调试模式下会显示详细的错误信息和调试信息
		// 生产环境应该设置为 false
		Debug: true,
		// ColorScheme 设置后台管理系统的配色方案
		// adminlte.ColorschemeSkinBlack 表示使用 AdminLTE 主题的黑色皮肤
		// AdminLTE 是一个流行的 Bootstrap 后台管理模板
		ColorScheme: adminlte.ColorschemeSkinBlack,
	}

	// 添加 Chart.js 图表组件到模板系统
	// Chart.js 是一个流行的 JavaScript 图表库，用于在网页上绘制各种图表
	// template.AddComp 函数将组件注册到模板系统，使其可以在页面中使用
	// NewChart() 函数创建一个新的图表组件实例
	template.AddComp(chartjs.NewChart())

	// 创建示例插件
	// 插件是 GoAdmin 的扩展机制，允许开发者添加自定义功能
	// example.NewExample() 创建一个示例插件实例
	// 插件可以添加自定义页面、数据表、API 端点等
	examplePlugin := example.NewExample()

	// 从 Go 插件系统加载插件（已注释）
	// Go 插件系统允许在运行时动态加载编译好的 .so 文件
	// 这使得可以在不重新编译主程序的情况下扩展功能
	// 示例：
	// examplePlugin := plugins.LoadFromPlugin("../datamodel/example.so")

	// 自定义登录页面（已注释）
	// 可以通过添加自定义组件来替换默认的登录页面
	// 示例代码参考：https://github.com/GoAdminGroup/demo.go-admin.cn/blob/master/main.go#L39
	// 示例：
	// template.AddComp("login", datamodel.LoginPage)

	// 从 JSON 文件加载配置（已注释）
	// 可以将配置信息存储在 JSON 文件中，便于管理和版本控制
	// 示例：
	// eng.AddConfigFromJSON("../datamodel/config.json")

	// 配置 GoAdmin 引擎
	// 使用链式调用方法来配置引擎
	// AddConfig 添加配置信息
	// AddGenerators 添加数据表生成器
	// AddDisplayFilterXssJsFilter 添加 XSS 和 JS 过滤器，用于防止跨站脚本攻击
	// AddGenerator 添加单个数据表生成器
	// AddPlugins 添加插件
	// Use 将引擎集成到 Chi 路由器中
	if err := eng.AddConfig(&cfg).
		AddGenerators(datamodel.Generators).
		AddDisplayFilterXssJsFilter().
		// 添加生成器，第一个参数是数据表访问时的 URL 前缀
		// 示例：
		// "user" => http://localhost:9033/admin/info/user
		//
		AddGenerator("user", datamodel.GetUserTable).
		AddPlugins(examplePlugin).
		Use(r); err != nil {
		// 如果配置过程中出现错误，使用 panic 终止程序
		// panic 是 Go 语言的内置函数，用于在遇到不可恢复的错误时终止程序
		// 在生产环境中，应该使用更优雅的错误处理方式
		panic(err)
	}

	// 获取当前工作目录
	// os.Getwd() 返回当前进程的工作目录
	// 工作目录是程序启动时的目录，不是可执行文件所在的目录
	workDir, _ := os.Getwd()
	// 构建文件存储目录的完整路径
	// filepath.Join 函数用于跨平台地拼接路径，自动处理路径分隔符
	// 在 Windows 上使用反斜杠，在 Unix/Linux 上使用正斜杠
	filesDir := filepath.Join(workDir, "uploads")
	// 设置静态文件服务器
	// FileServer 函数（在下方定义）用于提供静态文件服务
	// r 是 Chi 路由器，"/uploads" 是 URL 路径，http.Dir(filesDir) 是文件系统
	FileServer(r, "/uploads", http.Dir(filesDir))

	// 自定义页面
	// eng.HTML 方法用于添加自定义 HTML 页面
	// 第一个参数是 HTTP 方法（GET、POST 等）
	// 第二个参数是 URL 路径
	// 第三个参数是处理函数，返回页面的 HTML 内容
	// 这里将 /admin 路径映射到 datamodel.GetContent 函数
	eng.HTML("GET", "/admin", datamodel.GetContent)

	// 启动 HTTP 服务器
	// 使用 goroutine 在后台启动服务器
	// goroutine 是 Go 语言的轻量级线程，可以并发执行函数
	// 使用 go 关键字启动的函数会在新的 goroutine 中执行
	// http.ListenAndServe 在指定地址启动 HTTP 服务器
	// ":3333" 表示监听所有网络接口的 3333 端口
	// r 是路由处理器，所有请求都会通过 r 进行路由
	go func() {
		_ = http.ListenAndServe(":3333", r)
	}()

	// 设置优雅关闭机制
	// 创建一个用于接收信号的通道
	// make 函数用于创建通道，第二个参数指定通道缓冲区大小
	// 通道是 Go 语言用于在 goroutine 之间通信的主要方式
	quit := make(chan os.Signal, 1)
	// 注册要接收的信号
	// signal.Notify 函数将信号发送到指定的通道
	// os.Interrupt 表示中断信号（通常是 Ctrl+C）
	// 当用户按下 Ctrl+C 时，中断信号会被发送到 quit 通道
	signal.Notify(quit, os.Interrupt)
	// 阻塞等待信号
	// <-quit 表示从通道读取数据，如果通道为空，则会阻塞等待
	// 这行代码会一直阻塞，直到收到中断信号
	<-quit
	// 打印关闭数据库连接的消息
	// log.Print 函数输出日志信息，包含时间戳
	log.Print("closing database connection")
	// 关闭数据库连接
	// MysqlConnection() 返回 MySQL 数据库连接对象
	// Close() 方法关闭数据库连接，释放资源
	eng.MysqlConnection().Close()
}

// FileServer 便捷地设置 http.FileServer 处理器，用于从 http.FileSystem 提供静态文件服务
//
// 参数说明：
//   - r: Chi 路由器实例，用于注册路由
//   - path: URL 路径前缀，所有静态文件请求都会以此前缀开头
//   - root: 文件系统，包含要提供的静态文件
//
// 功能说明：
//  1. 检查路径是否包含 URL 参数（{} 或 *），如果包含则 panic
//  2. 使用 http.StripPrefix 去除 URL 路径前缀
//  3. 如果路径不以 / 结尾，添加重定向规则
//  4. 注册通配符路由，处理所有静态文件请求
//
// 使用示例：
//
//	FileServer(r, "/static", http.Dir("./static"))
//	这将把 /static 路径映射到 ./static 目录
//
// 注意事项：
//   - 路径不能包含 URL 参数（{} 或 *）
//   - 如果路径不以 / 结尾，会自动添加重定向规则
//   - 使用通配符 * 匹配所有子路径
//
// Go 语言概念：
//   - http.FileServer: 标准库提供的文件服务器处理器
//   - http.StripPrefix: 用于去除 URL 路径前缀的中间件
//   - http.FileSystem: 文件系统接口，http.Dir 实现了该接口
func FileServer(r chi.Router, path string, root http.FileSystem) {
	// 检查路径是否包含 URL 参数
	// strings.ContainsAny 函数检查字符串是否包含指定字符集中的任意字符
	// {}* 是 URL 路径中的特殊字符，用于参数匹配和通配符
	// FileServer 不支持这些字符，因为它们会导致路由冲突
	if strings.ContainsAny(path, "{}*") {
		// 如果路径包含非法字符，使用 panic 终止程序
		// 这是因为在运行时才发现配置错误，无法继续执行
		panic("FileServer does not permit URL parameters.")
	}

	// 创建文件服务器处理器
	// http.FileServer 创建一个文件服务器处理器，从指定的文件系统提供文件
	// http.StripPrefix 是一个中间件，用于去除请求 URL 中的指定前缀
	// 例如：请求 /uploads/image.jpg，去除 /uploads 前缀后，在文件系统中查找 image.jpg
	fs := http.StripPrefix(path, http.FileServer(root))

	// 如果路径不是根路径且不以 / 结尾，添加重定向规则
	// 这是为了确保 URL 规范化，避免重定向循环
	// 例如：/uploads 重定向到 /uploads/
	if path != "/" && path[len(path)-1] != '/' {
		// 注册重定向处理器
		// http.RedirectHandler 创建一个重定向处理器
		// path+"/" 是重定向目标路径
		// http.StatusMovedPermanently 表示永久重定向（HTTP 301）
		// ServeHTTP 方法处理 HTTP 请求
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		// 在路径后添加 /
		path += "/"
	}
	// 添加通配符，匹配所有子路径
	// * 是 Chi 路由器的通配符，匹配任意字符
	// 例如：/uploads/* 匹配 /uploads/image.jpg、/uploads/docs/file.pdf 等
	path += "*"

	// 注册通配符路由处理器
	// r.Get 注册 GET 方法的路由
	// http.HandlerFunc 将普通函数转换为 http.Handler 接口
	// 匿名函数处理请求，调用 fs.ServeHTTP 提供文件服务
	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 调用文件服务器处理请求
		// w 是响应写入器，用于写入 HTTP 响应
		// r 是请求对象，包含请求的所有信息
		fs.ServeHTTP(w, r)
	}))
}

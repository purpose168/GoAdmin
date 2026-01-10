// Copyright 2019 GoAdmin Core Team. All rights reserved.
// Use of this source code is governed by a Apache-2.0 style
// license that can be found in LICENSE file.

// Package buffalo 提供 GoAdmin 与 Buffalo Web 框架的适配器实现
//
// 文件名: buffalo.go
// 包名: buffalo
// 作者: GoAdmin Core Team
// 创建日期: 2019
//
// 功能描述:
// 本包实现了 GoAdmin 管理后台与 Buffalo Web 框架的适配器，允许 GoAdmin 在 Buffalo 应用中运行
// 该适配器作为 Buffalo 框架和 GoAdmin 管理后台之间的桥梁，实现了适配器接口的所有方法
// 通过该适配器，开发者可以在 Buffalo 应用中快速集成 GoAdmin 管理后台功能
//
// 核心概念:
// - 适配器模式: 将 GoAdmin 的核心功能适配到 Buffalo 框架中
// - 上下文转换: 将 Buffalo 的 buffalo.Context 转换为 GoAdmin 的 context.Context
// - 路由集成: 将 GoAdmin 的路由注册到 Buffalo 的路由系统中
// - 中间件链: 支持 Buffalo 的中间件和 GoAdmin 的处理器链
// - 请求处理: 统一处理 HTTP 请求和响应
// - Cookie 认证: 支持 GoAdmin 的 Cookie 认证机制
// - PJAX 支持: 支持 PJAX 技术实现页面部分更新
// - 多语言: 支持多语言切换功能
//
// 技术栈:
// - Buffalo: Go 语言的 Web 框架，提供 MVC 架构、路由、中间件等功能
// - GoAdmin: Go 语言的后台管理框架，提供数据表格、表单、图表等组件
// - Go 标准库: bytes、errors、net/http、net/url、regexp、strings 等
//
// 数据库支持:
// - MySQL
// - PostgreSQL
// - SQLite
// - MSSQL
//
// 配置说明:
// - 通过 config.SetCfg(cfg) 设置 GoAdmin 配置
// - 通过 config.GetLoginUrl() 获取登录页面 URL
// - 通过 config.Url(path) 生成带前缀的 URL
// - 通过 constant.EditPKKey 获取编辑主键的常量
//
// 使用示例:
//
//	package main
//
//	import (
//		"github.com/gobuffalo/buffalo"
//		"github.com/purpose168/GoAdmin"
//		"github.com/purpose168/GoAdmin/adapter/buffalo"
//		_ "github.com/purpose168/GoAdmin/plugins/admin"
//		"github.com/purpose168/GoAdmin/plugins/admin/modules/constant"
//		_ "github.com/purpose168/GoAdmin/plugins/example"
//	)
//
//	func main() {
//		// 初始化 Buffalo 应用
//		app := buffalo.New(buffalo.Options{
//			Env:  "development",
//			Addr: "127.0.0.1:3000",
//		})
//
//		// 设置 GoAdmin 配置
//		cfg := config.Config{
//			Domain: "localhost",
//			// ... 其他配置
//		}
//		config.SetCfg(cfg)
//
//		// 使用 Buffalo 适配器
//		_ = GoAdmin.Init(buffaloAdapter(app))
//
//		// 启动 Buffalo 服务
//		app.Serve()
//	}
//
//	func buffaloAdapter(app *buffalo.App) *buffalo.Buffalo {
//		eng := buffalo.New()
//		_ = eng.SetApp(app)
//		return eng
//	}
//
// 注意事项:
// - Buffalo 适配器需要 Buffalo v0.16+ 版本
// - 路由参数格式为 :param，会自动转换为 {param} 格式
// - 表单解析的最大内存限制为 32MB
// - Cookie 认证使用默认的 Cookie 键名
// - PJAX 请求通过 X-PJAX 请求头识别
// - 语言设置通过 __ga_lang 查询参数传递
// - 支持自定义 Content-Type 响应头
// - 支持自定义响应状态码和响应体
package buffalo

import (
	"bytes"          // bytes: 字节缓冲区操作，提供字节缓冲区的读写功能
	"errors"         // errors: 错误处理，提供错误创建和处理功能
	"net/http"       // net/http: HTTP 包，提供 HTTP 客户端和服务器功能
	neturl "net/url" // net/url: URL 解析和查询，提供 URL 解析和查询参数处理功能
	"regexp"         // regexp: 正则表达式，提供正则表达式匹配和替换功能
	"strings"        // strings: 字符串操作，提供字符串处理和转换功能

	"github.com/gobuffalo/buffalo"                                 // buffalo: Buffalo Web 框架，提供 MVC 架构、路由、中间件等功能
	"github.com/purpose168/GoAdmin/adapter"                        // adapter: GoAdmin 适配器包，提供适配器接口和基础适配器
	"github.com/purpose168/GoAdmin/context"                        // context: GoAdmin 上下文包，提供请求上下文和处理器链功能
	"github.com/purpose168/GoAdmin/engine"                         // engine: GoAdmin 引擎包，提供核心引擎和适配器注册功能
	"github.com/purpose168/GoAdmin/modules/config"                 // config: GoAdmin 配置模块，提供配置管理和 URL 生成功能
	"github.com/purpose168/GoAdmin/plugins"                        // plugins: GoAdmin 插件包，提供插件接口和插件管理功能
	"github.com/purpose168/GoAdmin/plugins/admin/models"           // models: GoAdmin 管理员模型包，提供用户模型和数据访问功能
	"github.com/purpose168/GoAdmin/plugins/admin/modules/constant" // constant: GoAdmin 常量模块，提供框架常量定义（如 PjaxHeader、EditPKKey 等）
	"github.com/purpose168/GoAdmin/template/types"                 // types: GoAdmin 模板类型包，提供面板、按钮等类型定义
)

// Buffalo 结构体实现了 GoAdmin 的适配器接口
// 它作为 Buffalo 框架和 GoAdmin 管理后台之间的桥梁，实现了适配器接口的所有方法
// 嵌入了 adapter.BaseAdapter 以获得基础适配器功能（如 GetUser、GetUse、GetContent 等）
// ctx 字段存储当前的 Buffalo 上下文，用于访问请求和响应对象
// app 字段存储 Buffalo 应用实例，用于注册路由和处理器
type Buffalo struct {
	adapter.BaseAdapter                 // BaseAdapter: 基础适配器，提供通用的适配器功能
	ctx                 buffalo.Context // ctx: Buffalo 上下文，存储当前请求的上下文信息
	app                 *buffalo.App    // app: Buffalo 应用实例，用于注册路由和处理器
}

// init 函数在包导入时自动执行
// Go 语言的 init 函数会在 main 函数之前自动调用，用于包级别的初始化
// 这里使用 init 函数将 Buffalo 适配器注册到 GoAdmin 引擎中
// 这样 GoAdmin 就知道如何使用 Buffalo 框架，可以通过 engine.Register 注册多个适配器
//
// 说明:
//
//	该函数在包导入时自动执行，无需手动调用
//	通过 engine.Register(new(Buffalo)) 将 Buffalo 适配器注册到 GoAdmin 引擎中
//	GoAdmin 引擎会维护一个适配器列表，支持多个框架适配器
//	注册后，可以通过 GoAdmin.Init() 初始化 GoAdmin，并使用注册的适配器
//	适配器注册是线程安全的，可以在多个 goroutine 中并发注册
//
// 注意事项:
//
//   - init 函数在包导入时自动执行，不能手动调用
//   - 一个包可以有多个 init 函数，执行顺序不确定
//   - 适配器注册应该在 main 函数之前完成
//   - 如果注册多个适配器，GoAdmin 会根据配置选择使用哪个适配器
func init() {
	engine.Register(new(Buffalo)) // 将 Buffalo 适配器注册到 GoAdmin 引擎中
}

// User 实现了 Adapter.User 方法
// 该方法用于从当前上下文中获取用户信息，用于身份验证和权限控制
//
// 参数:
//   - ctx: 上下文接口，通常为 buffalo.Context 类型，包含当前请求的上下文信息
//
// 返回值:
//   - models.UserModel: 用户模型，包含用户信息（如用户名、角色、权限等）
//   - bool: 是否成功获取用户信息，true 表示成功，false 表示失败
//
// 功能特性:
//   - 从 Buffalo 上下文中提取用户信息
//   - 支持多种认证方式（Cookie、Session、Token 等）
//   - 支持用户权限验证
//   - 支持用户角色管理
//   - 支持用户会话管理
//
// 说明:
//
//	该方法调用了基础适配器的 GetUser 方法，传入 Buffalo 上下文和当前适配器实例
//	GetUser 方法会从上下文中提取认证信息（如 Cookie），并查询数据库获取用户信息
//	如果用户不存在或认证失败，返回空的 UserModel 和 false
//	如果用户存在且认证成功，返回完整的 UserModel 和 true
//	UserModel 包含用户的基本信息（如用户名、邮箱、角色、权限等）
//	该方法在每次请求时都会被调用，用于验证用户身份和权限
//
// 注意事项:
//
//   - 该方法需要在请求处理前调用，确保用户已认证
//   - 如果用户未认证，返回的 UserModel 为空，bool 为 false
//   - 该方法会查询数据库，可能影响性能
//   - 建议在中间件中调用该方法，统一处理用户认证
//   - 该方法支持并发调用，但需要注意数据库连接池的配置
func (bu *Buffalo) User(ctx interface{}) (models.UserModel, bool) {
	return bu.GetUser(ctx, bu) // 调用基础适配器的 GetUser 方法，传入 Buffalo 上下文和当前适配器实例
}

// Use 实现了 Adapter.Use 方法
// 该方法用于将插件注册到 Buffalo 应用中，用于扩展 GoAdmin 的功能
//
// 参数:
//   - app: 应用接口，通常为 *buffalo.App 类型，表示 Buffalo 应用实例
//   - plugs: 插件列表，包含需要注册的所有插件（如管理员插件、示例插件等）
//
// 返回值:
//   - error: 错误信息，如果注册失败则返回错误，注册成功则返回 nil
//
// 功能特性:
//   - 将 GoAdmin 插件注册到 Buffalo 应用中
//   - 支持多个插件同时注册
//   - 支持插件路由注册
//   - 支持插件中间件注册
//   - 支持插件配置管理
//
// 说明:
//
//	该方法调用了基础适配器的 GetUse 方法，传入 Buffalo 应用实例和插件列表
//	GetUse 方法会遍历插件列表，为每个插件调用 Init 方法进行初始化
//	插件初始化会注册插件的路由、中间件、配置等
//	插件路由会通过 AddHandler 方法注册到 Buffalo 应用中
//	插件中间件会添加到处理器链中，在请求处理前执行
//	插件配置会存储到 GoAdmin 的配置中，供后续使用
//	该方法通常在应用启动时调用，用于初始化所有插件
//
// 注意事项:
//
//   - 该方法需要在应用启动前调用，确保所有插件正确初始化
//   - 插件列表不能为空，至少需要注册一个插件（如管理员插件）
//   - 插件初始化顺序很重要，依赖关系需要正确处理
//   - 如果插件初始化失败，会返回错误，应用无法启动
//   - 该方法支持并发调用，但需要注意插件初始化的线程安全性
func (bu *Buffalo) Use(app interface{}, plugs []plugins.Plugin) error {
	return bu.GetUse(app, plugs, bu) // 调用基础适配器的 GetUse 方法，传入 Buffalo 应用实例、插件列表和当前适配器实例
}

// Content 实现了 Adapter.Content 方法
// 该方法用于渲染管理面板内容，将 GoAdmin 的面板渲染为 HTML 响应
//
// 参数:
//   - ctx: 上下文接口，通常为 buffalo.Context 类型，包含当前请求的上下文信息
//   - getPanelFn: 获取面板的函数，返回 types.Panel 类型的面板，用于生成管理面板内容
//   - fn: 节点处理器，用于处理上下文中的节点，支持自定义处理逻辑
//   - btns: 导航按钮列表，可变参数，用于在面板顶部显示导航按钮
//
// 返回值:
//   - 无返回值
//
// 功能特性:
//   - 渲染 GoAdmin 管理面板为 HTML 响应
//   - 支持自定义面板内容
//   - 支持自定义节点处理器
//   - 支持自定义导航按钮
//   - 支持多种面板类型（表格、表单、图表等）
//
// 说明:
//
//	该方法调用了基础适配器的 GetContent 方法，传入 Buffalo 上下文、获取面板函数、导航按钮和节点处理器
//	GetContent 方法会调用 getPanelFn 函数获取面板内容
//	面板内容包含数据表格、表单、图表等组件
//	节点处理器 fn 用于处理面板中的节点，支持自定义处理逻辑
//	导航按钮 btns 会显示在面板顶部，用于快速导航
//	渲染后的 HTML 会写入到响应中，通过 Write 方法发送给客户端
//	该方法通常在路由处理器中调用，用于渲染管理页面
//
// 注意事项:
//
//   - 该方法需要在请求处理中调用，确保上下文已正确设置
//   - getPanelFn 函数不能为空，必须返回有效的面板
//   - 节点处理器 fn 可以为空，使用默认处理逻辑
//   - 导航按钮 btns 可以为空，不显示导航按钮
//   - 该方法会修改响应头和响应体，确保在调用前未写入响应
//   - 该方法支持并发调用，但需要注意面板生成的线程安全性
func (bu *Buffalo) Content(ctx interface{}, getPanelFn types.GetPanelFn, fn context.NodeProcessor, btns ...types.Button) {
	bu.GetContent(ctx, getPanelFn, bu, btns, fn) // 调用基础适配器的 GetContent 方法，传入 Buffalo 上下文、获取面板函数、当前适配器实例、导航按钮和节点处理器
}

// HandlerFunc 定义了处理函数的类型
// 该函数接收 Buffalo 上下文，返回面板和可能的错误，用于自定义处理逻辑
//
// 参数:
//   - ctx: Buffalo 上下文，包含当前请求的上下文信息
//
// 返回值:
//   - types.Panel: 管理面板，包含面板内容（如数据表格、表单、图表等）
//   - error: 错误信息，如果处理失败则返回错误，处理成功则返回 nil
//
// 功能特性:
//   - 支持自定义处理逻辑
//   - 支持返回自定义面板
//   - 支持错误处理
//   - 支持异步处理
//
// 说明:
//
//	该类型定义了处理函数的签名，用于在 Buffalo 路由处理器中使用
//	处理函数接收 Buffalo 上下文，可以访问请求和响应对象
//	处理函数返回面板和错误，面板会被渲染为 HTML 响应
//	如果处理失败，返回错误，GoAdmin 会显示错误页面
//	处理函数可以访问数据库、调用 API、执行业务逻辑等
//	该类型通常与 Content 辅助函数一起使用，用于创建自定义处理器
//
// 使用示例:
//
//	// 定义处理函数
//	handler := func(ctx buffalo.Context) (types.Panel, error) {
//		// 自定义处理逻辑
//		panel := types.NewPanel()
//		// ... 设置面板内容
//		return panel, nil
//	}
//
//	// 使用 Content 辅助函数创建 Buffalo 处理器
//	buffaloHandler := Content(handler)
//
//	// 注册处理器到 Buffalo 应用
//	app.GET("/admin", buffaloHandler)
//
// 注意事项:
//
//   - 处理函数不能为空，必须返回有效的面板或错误
//   - 处理函数应该是线程安全的，支持并发调用
//   - 处理函数应该尽快返回，避免阻塞请求
//   - 处理函数应该正确处理错误，避免 panic
//   - 处理函数应该避免修改传入的上下文，除非必要
type HandlerFunc func(ctx buffalo.Context) (types.Panel, error)

// Content 是一个辅助函数，用于将 HandlerFunc 转换为 Buffalo 的 Handler
// 这样可以在 Buffalo 的路由中使用 GoAdmin 的处理函数，简化集成过程
//
// 参数:
//   - handler: 处理函数，接收 Buffalo 上下文并返回面板，类型为 HandlerFunc
//
// 返回值:
//   - buffalo.Handler: Buffalo 处理器函数，可以在 Buffalo 的路由中使用
//
// 功能特性:
//   - 将 HandlerFunc 转换为 Buffalo 的 Handler
//   - 简化 GoAdmin 与 Buffalo 的集成
//   - 支持在 Buffalo 路由中使用 GoAdmin 处理函数
//   - 支持自定义处理逻辑
//
// 说明:
//
//	该辅助函数接收一个 HandlerFunc 类型的处理函数，返回一个 buffalo.Handler 类型的处理器函数
//	处理器函数内部调用 GoAdmin 的 engine.Content 方法，传入 Buffalo 上下文和处理函数
//	engine.Content 方法会调用处理函数，获取面板内容，并渲染为 HTML 响应
//	处理器函数可以在 Buffalo 的路由中使用，支持 GET、POST 等 HTTP 方法
//	使用该辅助函数可以简化 GoAdmin 与 Buffalo 的集成，无需手动处理上下文转换和响应渲染
//	该辅助函数通常在应用启动时使用，用于注册管理后台的路由和处理器
//
// 使用示例:
//
//	// 定义处理函数
//	handler := func(ctx buffalo.Context) (types.Panel, error) {
//		// 自定义处理逻辑
//		panel := types.NewPanel()
//		// ... 设置面板内容
//		return panel, nil
//	}
//
//	// 使用 Content 辅助函数创建 Buffalo 处理器
//	buffaloHandler := Content(handler)
//
//	// 注册处理器到 Buffalo 应用
//	app.GET("/admin", buffaloHandler)
//
// 注意事项:
//
//   - 处理函数不能为空，必须返回有效的面板或错误
//   - 处理器函数应该在适当的 HTTP 方法中注册（如 GET、POST）
//   - 处理器函数应该尽早注册，确保在其他路由之前执行
//   - 处理器函数应该正确处理错误，避免 panic
//   - 处理器函数应该避免修改传入的上下文，除非必要
func Content(handler HandlerFunc) buffalo.Handler {
	return func(ctx buffalo.Context) error { // 返回一个 Buffalo 处理器函数
		engine.Content(ctx, func(ctx interface{}) (types.Panel, error) { // 调用 GoAdmin 的 engine.Content 方法，传入 Buffalo 上下文和处理函数
			return handler(ctx.(buffalo.Context)) // 调用处理函数，传入 Buffalo 上下文，返回面板和错误
		})
		return nil // 返回 nil，表示处理成功
	}
}

// SetApp 实现了 Adapter.SetApp 方法
// 该方法用于设置 Buffalo 应用实例到适配器中，用于后续的路由注册和请求处理
//
// 参数:
//   - app: 应用接口，必须为 *buffalo.App 类型，表示 Buffalo 应用实例
//
// 返回值:
//   - error: 错误信息，如果参数类型不正确则返回错误，设置成功则返回 nil
//
// 功能特性:
//   - 设置 Buffalo 应用实例到适配器中
//   - 支持类型检查，确保参数类型正确
//   - 支持错误处理，提供友好的错误信息
//
// 说明:
//
//	该方法接收一个应用接口，使用类型断言检查是否为 *buffalo.App 类型
//	如果类型断言失败，返回错误 "buffalo 适配器 SetApp: 参数类型错误"
//	如果类型断言成功，将 Buffalo 应用实例存储到适配器的 app 字段中
//	Buffalo 应用实例用于后续的路由注册和请求处理
//	该方法通常在适配器初始化时调用，用于设置 Buffalo 应用实例
//	该方法在 Use 方法中被调用，传入 Buffalo 应用实例
//
// 注意事项:
//
//   - app 参数必须为 *buffalo.App 类型，否则返回错误
//   - 该方法应该在 Use 方法之前调用，确保应用实例已设置
//   - 该方法只能调用一次，重复调用会覆盖之前的应用实例
//   - 该方法支持并发调用，但需要注意应用实例的线程安全性
//   - 该方法返回错误后，适配器无法正常工作
func (bu *Buffalo) SetApp(app interface{}) error {
	var (
		eng *buffalo.App // eng: Buffalo 应用实例，用于存储类型断言后的结果
		ok  bool         // ok: 类型断言是否成功，true 表示成功，false 表示失败
	)
	// 使用类型断言检查 app 是否为 *buffalo.App 类型
	// ok 为 true 表示断言成功，eng 为转换后的值
	// ok 为 false 表示断言失败，返回错误
	if eng, ok = app.(*buffalo.App); !ok {
		return errors.New("buffalo 适配器 SetApp: 参数类型错误") // 返回错误，提示参数类型不正确
	}
	bu.app = eng // 将 Buffalo 应用实例存储到适配器的 app 字段中
	return nil   // 返回 nil，表示设置成功
}

// AddHandler 实现了 Adapter.AddHandler 方法
// 该方法用于向 Buffalo 应用添加路由处理器，将 GoAdmin 的路由注册到 Buffalo 的路由系统中
//
// 参数:
//   - method: HTTP 方法，如 "GET"、"POST"、"PUT"、"DELETE" 等
//   - path: 路由路径，如 "/admin"、"/admin/users" 等
//   - handlers: 处理器链，包含要执行的中间件和处理器，类型为 context.Handlers
//
// 返回值:
//   - 无返回值
//
// 功能特性:
//   - 将 GoAdmin 的路由注册到 Buffalo 的路由系统中
//   - 支持多种 HTTP 方法（GET、POST、PUT、DELETE 等）
//   - 支持路由参数（如 :id、:name 等）
//   - 支持处理器链（中间件和处理器）
//   - 支持上下文转换（Buffalo 上下文转换为 GoAdmin 上下文）
//   - 支持响应头复制（GoAdmin 响应头复制到 Buffalo 响应中）
//   - 支持响应体写入（GoAdmin 响应体写入到 Buffalo 响应中）
//
// 说明:
//
//	该方法调用 Buffalo 应用的路由注册函数，注册路由处理器
//	路由处理器接收 Buffalo 上下文，处理 HTTP 请求和响应
//	首先使用正则表达式将 Buffalo 的路由参数格式 :param 转换为 Go 标准格式 {param}
//	reg1 匹配中间的参数，如 :id/，转换为 {id}/
//	reg2 匹配结尾的参数，如 :id$，转换为 {id}
//	然后根据方法类型获取对应的路由注册函数（GET、POST、PUT、DELETE 等）
//	路由处理器内部处理请求：
//	1. 如果路径以 / 结尾，移除末尾的 /
//	2. 创建 GoAdmin 上下文
//	3. 获取 Buffalo 的路由参数并转换为 URL 查询参数
//	4. 执行处理器链
//	5. 将 GoAdmin 响应头复制到 Buffalo 响应中
//	6. 设置响应状态码
//	7. 将响应体写入 Buffalo 响应
//	该方法通常在 Use 方法中被调用，用于注册插件的路由
//
// 注意事项:
//
//   - method 参数必须为大写，如 "GET"、"POST" 等
//   - path 参数不能为空，必须以 "/" 开头
//   - handlers 参数不能为空，至少包含一个处理器
//   - 路由参数格式为 :param，会自动转换为 {param} 格式
//   - 处理器链执行顺序为从左到右，中间件先执行，处理器后执行
//   - 该方法支持并发调用，但需要注意路由注册的线程安全性
//   - 该方法会修改响应头和响应体，确保在调用前未写入响应
func (bu *Buffalo) AddHandler(method, path string, handlers context.Handlers) {
	url := path // 保存原始路径
	// 使用正则表达式将 Buffalo 的路由参数格式 :param 转换为 Go 标准格式 {param}
	// reg1 匹配中间的参数，如 :id/，转换为 {id}/
	reg1 := regexp.MustCompile(":(.*?)/") // 正则表达式：匹配冒号后跟任意字符，直到斜杠
	// reg2 匹配结尾的参数，如 :id$，转换为 {id}
	reg2 := regexp.MustCompile(":(.*?)$")     // 正则表达式：匹配冒号后跟任意字符，直到行尾
	url = reg1.ReplaceAllString(url, "{$1}/") // 将中间的参数替换为 {param}/ 格式
	url = reg2.ReplaceAllString(url, "{$1}")  // 将结尾的参数替换为 {param} 格式

	// 根据方法类型获取对应的路由注册函数
	getHandleFunc(bu.app, strings.ToUpper(method))(url, func(c buffalo.Context) error { // 调用 getHandleFunc 获取路由注册函数，注册路由处理器

		// 如果路径以 / 结尾，移除末尾的 /
		if c.Request().URL.Path[len(c.Request().URL.Path)-1] == '/' { // 检查路径最后一个字符是否为 /
			c.Request().URL.Path = c.Request().URL.Path[:len(c.Request().URL.Path)-1] // 移除末尾的 /
		}

		// 创建 GoAdmin 上下文
		ctx := context.NewContext(c.Request()) // 创建 GoAdmin 上下文，传入 HTTP 请求

		// 获取 Buffalo 的路由参数并转换为 URL 查询参数
		params := c.Params().(neturl.Values) // 获取 Buffalo 的路由参数，类型转换为 neturl.Values

		for key, param := range params { // 遍历路由参数
			if c.Request().URL.RawQuery == "" { // 如果 URL 查询参数为空
				c.Request().URL.RawQuery += strings.ReplaceAll(key, ":", "") + "=" + param[0] // 添加第一个查询参数，格式为 ?param=value
			} else { // 如果 URL 查询参数不为空
				c.Request().URL.RawQuery += "&" + strings.ReplaceAll(key, ":", "") + "=" + param[0] // 添加后续查询参数，格式为 &param=value
			}
		}

		// 执行处理器链
		ctx.SetHandlers(handlers).Next() // 设置处理器链并执行，处理器链包含中间件和处理器
		// 将 GoAdmin 响应头复制到 Buffalo 响应中
		for key, head := range ctx.Response.Header { // 遍历 GoAdmin 响应头
			c.Response().Header().Set(key, head[0]) // 将响应头设置到 Buffalo 响应中
		}
		// 将响应体写入 Buffalo 响应
		if ctx.Response.Body != nil { // 如果 GoAdmin 响应体不为空
			buf := new(bytes.Buffer)                          // 创建字节缓冲区
			_, _ = buf.ReadFrom(ctx.Response.Body)            // 从 GoAdmin 响应体读取数据到缓冲区
			c.Response().WriteHeader(ctx.Response.StatusCode) // 设置 Buffalo 响应的状态码为 GoAdmin 响应的状态码
			_, _ = c.Response().Write(buf.Bytes())            // 将缓冲区数据写入 Buffalo 响应
		} else { // 如果 GoAdmin 响应体为空
			c.Response().WriteHeader(ctx.Response.StatusCode) // 只设置 Buffalo 响应的状态码
		}
		return nil // 返回 nil，表示处理成功
	})
}

// HandleFun 定义了 Buffalo 路由方法的函数类型
// 该类型用于表示 Buffalo 的路由注册函数，如 GET、POST 等
//
// 参数:
//   - p: 路由路径
//   - h: 处理器函数
//
// 返回值:
//   - *buffalo.RouteInfo: 路由信息指针
//
// 功能特性:
//   - 定义 Buffalo 路由注册函数的签名
//   - 支持 GET、POST、PUT、DELETE 等 HTTP 方法
//   - 支持路由信息返回
//
// 说明:
//
//	该类型定义了 Buffalo 路由注册函数的签名，用于在 getHandleFunc 函数中使用
//	Buffalo 的路由注册函数包括 GET、POST、PUT、DELETE、HEAD、OPTIONS、PATCH 等
//	该类型接收路由路径和处理器函数，返回路由信息指针
//	路由信息包含路由路径、处理器函数、HTTP 方法等信息
//	该类型通常在 getHandleFunc 函数中使用，用于根据方法名称返回对应的路由注册函数
//
// 注意事项:
//
//   - 路由路径不能为空，必须以 "/" 开头
//   - 处理器函数不能为空，必须返回有效的响应
//   - 路由信息指针可能为 nil，表示路由注册失败
type HandleFun func(p string, h buffalo.Handler) *buffalo.RouteInfo

// getHandleFunc 根据方法名称返回对应的路由注册函数
// 这是一个辅助函数，用于将字符串方法名转换为 Buffalo 的路由注册函数
//
// 参数:
//   - eng: Buffalo 应用实例
//   - method: HTTP 方法名称，如 "GET"、"POST" 等
//
// 返回值:
//   - HandleFun: 对应的路由注册函数
//
// 功能特性:
//   - 根据方法名称返回对应的路由注册函数
//   - 支持所有 HTTP 方法（GET、POST、PUT、DELETE、HEAD、OPTIONS、PATCH）
//   - 支持错误处理，提供友好的错误信息
//
// 说明:
//
//	该函数接收 Buffalo 应用实例和 HTTP 方法名称，返回对应的路由注册函数
//	支持的方法包括 GET、POST、PUT、DELETE、HEAD、OPTIONS、PATCH
//	如果方法名称不支持，panic 抛出错误 "错误的 HTTP 方法"
//	该函数通常在 AddHandler 方法中被调用，用于根据方法类型获取路由注册函数
//
// 注意事项:
//
//   - method 参数必须为大写，如 "GET"、"POST" 等
//   - method 参数必须是支持的 HTTP 方法之一
//   - 如果方法名称不支持，会 panic 抛出错误
//   - 该函数支持并发调用，但需要注意应用实例的线程安全性
func getHandleFunc(eng *buffalo.App, method string) HandleFun {
	switch method { // 根据 HTTP 方法名称选择对应的路由注册函数
	case "GET": // GET 方法
		return eng.GET // 返回 GET 路由注册函数
	case "POST": // POST 方法
		return eng.POST // 返回 POST 路由注册函数
	case "PUT": // PUT 方法
		return eng.PUT // 返回 PUT 路由注册函数
	case "DELETE": // DELETE 方法
		return eng.DELETE // 返回 DELETE 路由注册函数
	case "HEAD": // HEAD 方法
		return eng.HEAD // 返回 HEAD 路由注册函数
	case "OPTIONS": // OPTIONS 方法
		return eng.OPTIONS // 返回 OPTIONS 路由注册函数
	case "PATCH": // PATCH 方法
		return eng.PATCH // 返回 PATCH 路由注册函数
	default: // 不支持的 HTTP 方法
		panic("错误的 HTTP 方法") // panic 抛出错误，提示 HTTP 方法错误
	}
}

// Name 实现了 Adapter.Name 方法
// 该方法返回适配器的名称，用于标识适配器类型
//
// 返回值:
//   - string: 适配器名称，固定为 "buffalo"
//
// 功能特性:
//   - 返回适配器名称
//   - 用于标识适配器类型
//   - 用于日志记录和调试
//
// 说明:
//
//	该方法返回适配器的名称，固定为 "buffalo"
//	适配器名称用于标识适配器类型，可以在日志中记录，便于调试
//	适配器名称也可以用于配置文件中，指定使用哪个适配器
//	该方法通常在适配器注册时被调用，用于标识适配器类型
//
// 注意事项:
//
//   - 该方法返回的名称必须与适配器类型一致
//   - 该方法返回的名称不能为空
//   - 该方法返回的名称应该唯一，避免与其他适配器冲突
func (*Buffalo) Name() string {
	return "buffalo" // 返回适配器名称 "buffalo"
}

// SetContext 实现了 Adapter.SetContext 方法
// 该方法用于设置当前请求的上下文，用于后续的请求处理
//
// 参数:
//   - contextInterface: 上下文接口，必须为 buffalo.Context 类型，表示 Buffalo 上下文
//
// 返回值:
//   - adapter.WebFrameWork: 返回设置了上下文的新适配器实例
//
// 功能特性:
//   - 设置当前请求的上下文
//   - 支持类型检查，确保参数类型正确
//   - 支持错误处理，提供友好的错误信息
//   - 返回新的适配器实例，避免修改原适配器
//
// 说明:
//
//	该方法接收一个上下文接口，使用类型断言检查是否为 buffalo.Context 类型
//	如果类型断言失败，panic 抛出错误 "buffalo 适配器 SetContext: 参数类型错误"
//	如果类型断言成功，创建新的适配器实例，设置 Buffalo 上下文到 ctx 字段中
//	新的适配器实例用于后续的请求处理，避免修改原适配器
//	该方法通常在每个请求处理时被调用，用于设置当前请求的上下文
//	该方法在 User、Content、Redirect、SetContentType、Write、GetCookie、Lang、Path、Method、FormParam、IsPjax、Query、Request 等方法中被使用
//
// 注意事项:
//
//   - contextInterface 参数必须为 buffalo.Context 类型，否则 panic
//   - 该方法返回新的适配器实例，不会修改原适配器
//   - 该方法应该在每个请求处理时调用，确保上下文正确设置
//   - 该方法支持并发调用，但需要注意上下文的线程安全性
//   - 该方法 panic 后，请求处理会中断，返回 500 错误
func (*Buffalo) SetContext(contextInterface interface{}) adapter.WebFrameWork {
	var (
		ctx buffalo.Context // ctx: Buffalo 上下文，用于存储类型断言后的结果
		ok  bool            // ok: 类型断言是否成功，true 表示成功，false 表示失败
	)
	// 使用类型断言检查 contextInterface 是否为 buffalo.Context 类型
	// ok 为 true 表示断言成功，ctx 为转换后的值
	// ok 为 false 表示断言失败，panic 抛出错误
	if ctx, ok = contextInterface.(buffalo.Context); !ok {
		panic("buffalo 适配器 SetContext: 参数类型错误") // panic 抛出错误，提示参数类型不正确
	}
	return &Buffalo{ctx: ctx} // 返回新的适配器实例，设置 Buffalo 上下文到 ctx 字段中
}

// Redirect 实现了 Adapter.Redirect 方法
// 该方法用于重定向到登录页面，使用 HTTP 302 状态码进行临时重定向
//
// 返回值:
//   - 无返回值
//
// 功能特性:
//   - 重定向到登录页面
//   - 使用 HTTP 302 状态码进行临时重定向
//   - 支持自定义登录页面 URL
//
// 说明:
//
//	该方法调用 Buffalo 上下文的 Redirect 方法，重定向到登录页面
//	使用 HTTP 302 状态码进行临时重定向，表示资源临时移动
//	登录页面 URL 通过 config.GetLoginUrl() 获取，支持自定义登录页面 URL
//	重定向 URL 通过 config.Url() 生成，支持 URL 前缀配置
//	该方法通常在用户未认证时被调用，重定向到登录页面
//	该方法在 GoAdmin 的中间件中被调用，用于检查用户认证状态
//
// 注意事项:
//
//   - 该方法会修改响应头和响应体，确保在调用前未写入响应
//   - 该方法会终止请求处理，后续代码不会执行
//   - 该方法应该在用户未认证时调用，避免无限重定向
//   - 该方法支持并发调用，但需要注意重定向的线程安全性
func (bu *Buffalo) Redirect() {
	_ = bu.ctx.Redirect(http.StatusFound, config.Url(config.GetLoginUrl())) // 调用 Buffalo 上下文的 Redirect 方法，重定向到登录页面，忽略错误
}

// SetContentType 实现了 Adapter.SetContentType 方法
// 该方法用于设置响应的 Content-Type 头，Content-Type 由 HTMLContentType() 方法确定
//
// 返回值:
//   - 无返回值
//
// 功能特性:
//   - 设置响应的 Content-Type 头
//   - 支持自定义 Content-Type
//   - Content-Type 由 HTMLContentType() 方法确定
//
// 说明:
//
//	该方法调用 Buffalo 响应头的 Set 方法，设置 Content-Type 头
//	Content-Type 由 HTMLContentType() 方法确定，通常为 "text/html; charset=utf-8"
//	HTMLContentType() 方法是基础适配器的方法，支持自定义 Content-Type
//	该方法通常在渲染 HTML 响应前被调用，设置正确的 Content-Type
//	该方法在 GoAdmin 的渲染器中被调用，用于设置响应的 Content-Type
//
// 注意事项:
//
//   - 该方法会修改响应头，确保在调用前未写入响应
//   - 该方法应该在渲染 HTML 响应前调用，确保 Content-Type 正确
//   - 该方法支持并发调用，但需要注意响应头的线程安全性
func (bu *Buffalo) SetContentType() {
	bu.ctx.Response().Header().Set("Content-Type", bu.HTMLContentType()) // 调用 Buffalo 响应头的 Set 方法，设置 Content-Type 头
}

// Write 实现了 Adapter.Write 方法
// 该方法用于将响应体写入到响应中
//
// 参数:
//   - body: 要写入的响应体字节数组，包含 HTML、JSON 等数据
//
// 返回值:
//   - 无返回值
//
// 功能特性:
//   - 将响应体写入到响应中
//   - 支持写入任意字节数组
//   - 支持写入 HTML、JSON 等数据
//
// 说明:
//
//	该方法调用 Buffalo 响应的 WriteHeader 和 Write 方法，将响应体写入到响应中
//	首先设置响应状态码为 200（OK），然后写入响应体
//	响应体可以是 HTML、JSON 等数据，由调用方决定
//	该方法通常在渲染完成后被调用，将渲染结果写入响应
//	该方法在 GoAdmin 的渲染器中被调用，用于写入响应体
//
// 注意事项:
//
//   - 该方法会修改响应体，确保在调用前已设置响应头
//   - 该方法应该在渲染完成后调用，确保响应体完整
//   - 该方法支持并发调用，但需要注意响应体的线程安全性
//   - 该方法忽略写入错误，使用 _ 忽略返回值
func (bu *Buffalo) Write(body []byte) {
	bu.ctx.Response().WriteHeader(http.StatusOK) // 设置 Buffalo 响应的状态码为 200（OK）
	_, _ = bu.ctx.Response().Write(body)         // 调用 Buffalo 响应的 Write 方法，将响应体写入到响应中，忽略写入错误
}

// GetCookie 实现了 Adapter.GetCookie 方法
// 该方法用于从请求中获取认证 Cookie，用于用户身份验证
//
// 返回值:
//   - string: Cookie 值，包含用户认证信息
//   - error: 错误信息，如果获取失败则返回错误
//
// 功能特性:
//   - 从请求中获取认证 Cookie
//   - 支持自定义 Cookie 键名
//   - 支持错误处理
//
// 说明:
//
//	该方法调用 Buffalo 上下文的 Cookies().Get() 方法，获取认证 Cookie
//	Cookie 键名通过 CookieKey() 方法获取，支持自定义 Cookie 键名
//	CookieKey() 方法是基础适配器的方法，默认为 "admin_cookie"
//	该方法返回 Cookie 值和错误，如果 Cookie 不存在则返回错误
//	该方法通常在 User 方法中被调用，用于获取用户认证信息
//
// 注意事项:
//
//   - 该方法返回的 Cookie 值可能为空，表示用户未认证
//   - 该方法返回的错误表示 Cookie 不存在或获取失败
//   - 该方法支持并发调用，但需要注意 Cookie 的线程安全性
//   - 该方法应该在 User 方法之前调用，确保 Cookie 正确获取
func (bu *Buffalo) GetCookie() (string, error) {
	return bu.ctx.Cookies().Get(bu.CookieKey()) // 调用 Buffalo 上下文的 Cookies().Get() 方法，获取认证 Cookie
}

// Lang 实现了 Adapter.Lang 方法
// 该方法用于从 URL 查询参数中获取语言设置，用于多语言支持
//
// 返回值:
//   - string: 语言代码，如 "zh-CN"、"en-US" 等
//
// 功能特性:
//   - 从 URL 查询参数中获取语言设置
//   - 支持自定义语言参数名
//   - 支持多语言切换
//
// 说明:
//
//	该方法调用 Buffalo 请求的 URL.Query().Get() 方法，获取语言设置
//	语言参数名为 "__ga_lang"，支持自定义语言参数名
//	语言代码格式为 "zh-CN"、"en-US" 等，遵循 RFC 5646 标准
//	该方法返回语言代码，如果未设置则返回空字符串
//	该方法通常在渲染模板时被调用，用于设置语言环境
//
// 注意事项:
//
//   - 该方法返回的语言代码可能为空，表示使用默认语言
//   - 该方法支持并发调用，但需要注意查询参数的线程安全性
//   - 该方法应该在渲染模板前调用，确保语言环境正确设置
//   - 该方法返回的语言代码应该有效，否则可能导致渲染错误
func (bu *Buffalo) Lang() string {
	return bu.ctx.Request().URL.Query().Get("__ga_lang") // 调用 Buffalo 请求的 URL.Query().Get() 方法，获取语言设置
}

// Path 实现了 Adapter.Path 方法
// 该方法用于获取当前请求的路径，用于路由匹配和日志记录
//
// 返回值:
//   - string: 请求路径，如 "/admin/dashboard"、"/admin/users" 等
//
// 功能特性:
//   - 获取当前请求的路径
//   - 不包含查询参数
//   - 支持路由匹配
//
// 说明:
//
//	该方法调用 Buffalo 请求的 URL.Path 属性，获取当前请求的路径
//	路径不包含查询参数，如 "/admin/dashboard?page=1" 的路径为 "/admin/dashboard"
//	路径以 "/" 开头，表示根路径
//	该方法返回路径，如果路径为空则返回 "/"
//	该方法通常在路由匹配和日志记录时被调用
//
// 注意事项:
//
//   - 该方法返回的路径不包含查询参数
//   - 该方法返回的路径以 "/" 开头
//   - 该方法支持并发调用，但需要注意请求的线程安全性
//   - 该方法应该在请求处理中调用，确保路径正确获取
func (bu *Buffalo) Path() string {
	return bu.ctx.Request().URL.Path // 调用 Buffalo 请求的 URL.Path 属性，获取当前请求的路径
}

// Method 实现了 Adapter.Method 方法
// 该方法用于获取当前请求的 HTTP 方法，用于请求类型判断和日志记录
//
// 返回值:
//   - string: HTTP 方法，如 "GET"、"POST"、"PUT"、"DELETE" 等
//
// 功能特性:
//   - 获取当前请求的 HTTP 方法
//   - 支持所有 HTTP 方法
//   - 支持请求类型判断
//
// 说明:
//
//	该方法调用 Buffalo 请求的 Method 属性，获取当前请求的 HTTP 方法
//	HTTP 方法包括 GET、POST、PUT、DELETE、PATCH、HEAD、OPTIONS 等
//	该方法返回 HTTP 方法，为大写字母，如 "GET"、"POST" 等
//	该方法通常在请求类型判断和日志记录时被调用
//
// 注意事项:
//
//   - 该方法返回的 HTTP 方法为大写字母
//   - 该方法支持所有 HTTP 方法
//   - 该方法支持并发调用，但需要注意请求的线程安全性
//   - 该方法应该在请求处理中调用，确保 HTTP 方法正确获取
func (bu *Buffalo) Method() string {
	return bu.ctx.Request().Method // 调用 Buffalo 请求的 Method 属性，获取当前请求的 HTTP 方法
}

// FormParam 实现了 Adapter.FormParam 方法
// 该方法用于解析并获取表单参数，用于表单数据处理
//
// 返回值:
//   - neturl.Values: 表单参数的键值对集合，包含所有表单字段
//
// 功能特性:
//   - 解析并获取表单参数
//   - 支持多种表单类型（application/x-www-form-urlencoded、multipart/form-data 等）
//   - 支持文件上传
//   - 解析的最大内存限制为 32MB
//
// 说明:
//
//	该方法调用 Buffalo 请求的 ParseMultipartForm 方法，解析表单参数
//	解析的最大内存限制为 32MB，超过限制会返回错误
//	该方法返回表单参数的键值对集合，类型为 neturl.Values
//	neturl.Values 是 map[string][]string 类型，支持多值字段
//	该方法支持多种表单类型，包括 application/x-www-form-urlencoded 和 multipart/form-data
//	该方法通常在表单提交时被调用，用于获取表单数据
//
// 注意事项:
//
//   - 该方法会解析表单参数，可能影响性能
//   - 该方法的最大内存限制为 32MB，超过限制会返回错误
//   - 该方法忽略解析错误，使用 _ 忽略返回值
//   - 该方法支持并发调用，但需要注意表单解析的线程安全性
//   - 该方法应该在表单提交时调用，确保表单参数正确解析
//   - 该方法应该只调用一次，多次调用会重复解析表单
func (bu *Buffalo) FormParam() neturl.Values {
	_ = bu.ctx.Request().ParseMultipartForm(32 << 20) // 调用 Buffalo 请求的 ParseMultipartForm 方法，解析表单参数，最大内存限制为 32MB
	return bu.ctx.Request().PostForm                  // 返回表单参数的键值对集合
}

// IsPjax 实现了 Adapter.IsPjax 方法
// 该方法用于检查当前请求是否为 PJAX 请求，PJAX 是一种使用 AJAX 技术实现页面部分更新的技术
//
// 返回值:
//   - bool: 如果是 PJAX 请求则返回 true，否则返回 false
//
// 功能特性:
//   - 检查当前请求是否为 PJAX 请求
//   - 支持页面部分更新
//   - 支持无刷新导航
//
// 说明:
//
//	该方法调用 Buffalo 请求的 Header.Get() 方法，获取 PJAX 请求头
//	PJAX 请求头为 "X-PJAX"，由 constant.PjaxHeader 常量定义
//	如果 PJAX 请求头的值为 "true"，则返回 true，表示是 PJAX 请求
//	如果 PJAX 请求头的值不为 "true" 或不存在，则返回 false，表示不是 PJAX 请求
//	PJAX 是一种使用 AJAX 技术实现页面部分更新的技术，可以避免整页刷新
//	该方法通常在渲染模板时被调用，用于判断是否渲染完整页面或部分内容
//
// 注意事项:
//
//   - 该方法返回的布尔值表示是否为 PJAX 请求
//   - 该方法支持并发调用，但需要注意请求头的线程安全性
//   - 该方法应该在渲染模板前调用，确保 PJAX 状态正确判断
//   - 该方法应该与前端 PJAX 库配合使用，确保 PJAX 请求正确识别
func (bu *Buffalo) IsPjax() bool {
	return bu.ctx.Request().Header.Get(constant.PjaxHeader) == "true" // 调用 Buffalo 请求的 Header.Get() 方法，获取 PJAX 请求头，判断是否为 "true"
}

// Query 实现了 Adapter.Query 方法
// 该方法用于获取 URL 查询参数，用于查询参数处理和分页等
//
// 返回值:
//   - neturl.Values: 查询参数的键值对集合，包含所有查询参数
//
// 功能特性:
//   - 获取 URL 查询参数
//   - 支持多值字段
//   - 支持分页、排序、筛选等功能
//
// 说明:
//
//	该方法调用 Buffalo 请求的 URL.Query() 方法，获取 URL 查询参数
//	查询参数格式为 ?key1=value1&key2=value2，如 ?page=1&pageSize=10
//	该方法返回查询参数的键值对集合，类型为 neturl.Values
//	neturl.Values 是 map[string][]string 类型，支持多值字段
//	该方法支持分页、排序、筛选等功能，通过查询参数传递
//	该方法通常在数据查询时被调用，用于获取查询参数
//
// 注意事项:
//
//   - 该方法返回的查询参数可能为空，表示没有查询参数
//   - 该方法支持多值字段，如 ?id=1&id=2
//   - 该方法支持并发调用，但需要注意查询参数的线程安全性
//   - 该方法应该在数据查询时调用，确保查询参数正确获取
func (bu *Buffalo) Query() neturl.Values {
	return bu.ctx.Request().URL.Query() // 调用 Buffalo 请求的 URL.Query() 方法，获取 URL 查询参数
}

// Request 实现了 Adapter.Request 方法
// 该方法用于获取原始的 HTTP 请求对象，用于底层请求处理
//
// 返回值:
//   - *http.Request: HTTP 请求对象指针，包含请求的所有信息
//
// 功能特性:
//   - 获取原始的 HTTP 请求对象
//   - 支持底层请求处理
//   - 支持自定义请求处理逻辑
//
// 说明:
//
//	该方法调用 Buffalo 上下文的 Request() 方法，获取原始的 HTTP 请求对象
//	HTTP 请求对象包含请求的所有信息，如 URL、方法、头、体等
//	该方法返回的是原始的 HTTP 请求对象，不是 Buffalo 的上下文对象
//	该方法通常在底层请求处理时被调用，用于自定义请求处理逻辑
//
// 注意事项:
//
//   - 该方法返回的 HTTP 请求对象可能为空，表示请求未初始化
//   - 该方法返回的 HTTP 请求对象是原始对象，修改会影响请求处理
//   - 该方法支持并发调用，但需要注意请求对象的线程安全性
//   - 该方法应该在底层请求处理时调用，确保请求对象正确获取
//   - 该方法应该谨慎使用，避免修改请求对象导致意外行为
func (bu *Buffalo) Request() *http.Request {
	return bu.ctx.Request() // 调用 Buffalo 上下文的 Request() 方法，获取原始的 HTTP 请求对象
}

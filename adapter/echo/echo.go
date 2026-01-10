// Copyright 2019 GoAdmin Core Team. All rights reserved.
// Use of this source code is governed by a Apache-2.0 style
// license that can be found in LICENSE file.

// echo 包提供了 GoAdmin 与 Echo Web 框架的适配器实现
//
// 功能描述:
//
//	该适配器允许 GoAdmin 管理后台在 Echo 应用中运行
//	Echo 是一个高性能、极简的 Go Web 框架，专注于 HTTP 服务
//	该适配器实现了 GoAdmin 的适配器接口，将 Echo 的请求和响应对象转换为 GoAdmin 的上下文
//	支持 Echo 的路由系统，包括路由参数、中间件链、请求处理等
//	支持 Cookie 认证、PJAX 请求、多语言等 GoAdmin 核心功能
//
// 核心概念:
//
//	适配器模式: 将 GoAdmin 管理后台集成到 Echo 框架中，实现框架无关性
//	上下文转换: 将 Echo 的 echo.Context 转换为 GoAdmin 的上下文
//	路由集成: 将 GoAdmin 的路由注册到 Echo 的路由系统中，支持路由参数和中间件
//	中间件链: 支持 Echo 的中间件链，实现请求预处理和后处理
//	请求处理: 处理 HTTP 请求，包括表单解析、Cookie 获取、查询参数处理等
//	Cookie 认证: 支持 Cookie 认证，用于用户身份验证和会话管理
//	PJAX 支持: 支持 PJAX 请求，实现页面部分更新和无刷新导航
//	多语言: 支持多语言切换，通过 URL 查询参数 __ga_lang 指定语言
//
// 技术栈:
//
//	Echo: 高性能、极简的 Go Web 框架，提供路由、中间件、上下文等功能
//	GoAdmin: Go 后台管理框架，提供管理后台、表单、表格、权限等功能
//	Go 标准库: net/http、net/url、bytes、strings 等，提供 HTTP、URL、字节缓冲区、字符串处理等功能
//
// 数据库支持:
//
//	MySQL: 通过 go-sql-driver/mysql 驱动支持
//	PostgreSQL: 通过 lib/pq 驱动支持
//	SQLite: 通过 mattn/go-sqlite3 驱动支持
//	MSSQL: 通过 go-mssqldb 驱动支持
//
// 配置说明:
//
//	config.SetCfg: 设置 GoAdmin 的全局配置，包括数据库连接、语言、主题等
//	config.GetLoginUrl: 获取登录页面的 URL 路径，默认为 "/login"
//	config.Url: 生成完整的 URL，包含域名和路径
//	constant.EditPKKey: 编辑对象的主键参数名，默认为 "pk"
//	constant.PjaxHeader: PJAX 请求头名称，默认为 "X-PJAX"
//
// 使用示例:
//
//	package main
//
//	import (
//	    "github.com/labstack/echo/v4"
//	    "github.com/purpose168/GoAdmin/adapter/echo"
//	    "github.com/purpose168/GoAdmin/plugins/admin"
//	    "github.com/purpose168/GoAdmin/modules/config"
//	)
//
//	func main() {
//	    // 1. 初始化 Echo 应用
//	    e := echo.New()
//
//	    // 2. 设置 GoAdmin 配置
//	    cfg := config.Config{
//	        Databases: config.DatabaseList{
//	            {
//	                Host:     "127.0.0.1",
//	                Port:     "3306",
//	                User:     "root",
//	                Pwd:      "root",
//	                Name:     "goadmin",
//	                Driver:   "mysql",
//	            },
//	        },
//	        Language: "zh-CN",
//	        Theme:    "adminlte",
//	    }
//	    config.SetCfg(&cfg)
//
//	    // 3. 使用 Echo 适配器
//	    adminPlugin := admin.NewAdmin(generators.Tables)
//	    echoAdapter := echo.NewEchoAdapter()
//	    echoAdapter.Use(e, []plugins.Plugin{adminPlugin})
//
//	    // 4. 启动 Echo 服务
//	    e.Logger.Fatal(e.Start(":8080"))
//	}
//
// 注意事项:
//
//	- Echo 版本: 建议使用 Echo v4.x 或更高版本
//	- 路由参数格式: Echo 的路由参数格式为 :param，会自动转换为 URL 查询参数 ?param=value
//	- 表单解析限制: 默认最大内存限制为 32MB，可通过 ParseMultipartForm 方法调整
//	- Cookie 认证: 默认使用 Cookie 键名 "admin_cookie"，可通过 CookieKey() 方法自定义
//	- PJAX 请求: 默认使用请求头 "X-PJAX" 判断是否为 PJAX 请求
//	- 语言设置: 默认通过 URL 查询参数 __ga_lang 设置语言，如 ?__ga_lang=zh-CN
//	- Content-Type: 默认使用 HTMLContentType() 方法返回的 Content-Type，通常为 "text/html; charset=utf-8"
//	- 响应状态码: 默认使用 HTTP 200（OK）作为响应状态码，可通过 WriteHeader 方法自定义
//	- 响应体写入: 默认使用 Write 方法写入响应体，支持写入 HTML、JSON 等数据
//
// 包名：echo
// 作者：GoAdmin Core Team
// 创建日期：2019
// 目的：为 Echo 框架提供 GoAdmin 管理后台的集成支持

package echo

import (
	"bytes"    // 字节缓冲区操作，提供字节缓冲区的读写功能
	"errors"   // 错误处理，提供错误创建和处理功能
	"net/http" // HTTP 包，提供 HTTP 客户端和服务器功能
	"net/url"  // URL 解析和查询，提供 URL 解析和查询参数处理功能
	"strings"  // 字符串操作，提供字符串处理和转换功能

	"github.com/labstack/echo/v4"                                  // Echo Web 框架，提供路由、中间件、上下文等功能
	"github.com/purpose168/GoAdmin/adapter"                        // GoAdmin 适配器包，提供适配器接口和基础适配器
	"github.com/purpose168/GoAdmin/context"                        // GoAdmin 上下文包，提供请求上下文和处理器链功能
	"github.com/purpose168/GoAdmin/engine"                         // GoAdmin 引擎包，提供核心引擎和适配器注册功能
	"github.com/purpose168/GoAdmin/modules/config"                 // GoAdmin 配置模块，提供配置管理和 URL 生成功能
	"github.com/purpose168/GoAdmin/plugins"                        // GoAdmin 插件包，提供插件接口和插件管理功能
	"github.com/purpose168/GoAdmin/plugins/admin/models"           // GoAdmin 管理员模型包，提供用户模型和数据访问功能
	"github.com/purpose168/GoAdmin/plugins/admin/modules/constant" // GoAdmin 常量模块，提供框架常量定义（如 PjaxHeader、EditPKKey 等）
	"github.com/purpose168/GoAdmin/template/types"                 // GoAdmin 模板类型包，提供面板、按钮等类型定义
)

// Echo 结构体实现了 GoAdmin 的适配器接口
// 它作为 Echo 框架和 GoAdmin 管理后台之间的桥梁，实现了适配器接口的所有方法
// 嵌入了 adapter.BaseAdapter 以获得基础适配器功能（如 GetUser、GetUse、GetContent 等）
// ctx 字段存储当前的 Echo 上下文，用于访问请求和响应对象
// app 字段存储 Echo 应用实例，用于注册路由和处理器
type Echo struct {
	adapter.BaseAdapter              // 嵌入基础适配器，获得基础适配器功能
	ctx                 echo.Context // 当前 Echo 上下文，包含请求和响应对象
	app                 *echo.Echo   // Echo 应用实例，用于注册路由和处理器
}

// init 函数在包导入时自动执行
// Go 语言的 init 函数会在 main 函数之前自动调用
// 这里使用 init 函数将 Echo 适配器注册到 GoAdmin 引擎中
// 这样 GoAdmin 就知道如何使用 Echo 框架
//
// 功能特性:
//
//	自动执行: init 函数在包导入时自动执行，无需手动调用
//	适配器注册: 将 Echo 适配器注册到 GoAdmin 引擎中，使 GoAdmin 能够使用 Echo 框架
//	支持多个适配器: 可以同时注册多个适配器，GoAdmin 会根据配置选择使用哪个适配器
//	线程安全: 注册过程是线程安全的，可以在多个 goroutine 中并发调用
//
// 说明:
//
//	该函数使用 engine.Register 方法将 Echo 适配器注册到 GoAdmin 引擎中
//	Register 方法接收一个适配器接口，这里传入 new(Echo) 创建 Echo 适配器实例
//	注册后，GoAdmin 引擎会保存 Echo 适配器的引用，用于后续的路由注册和请求处理
//	该函数在包导入时自动执行，确保适配器在 main 函数执行前完成注册
//
// 注意事项:
//
//   - 该函数不能手动调用，由 Go 语言运行时自动调用
//   - 该函数的执行顺序不确定，如果有多个 init 函数，执行顺序由包导入顺序决定
//   - 该函数应该在 main 函数之前完成，确保适配器在应用启动前注册
//   - 该函数支持并发调用，但通常只执行一次
func init() {
	engine.Register(new(Echo)) // 将 Echo 适配器注册到 GoAdmin 引擎中
}

// User 实现了 Adapter.User 方法
// 该方法用于从当前上下文中获取用户信息，用于身份验证和权限控制
//
// 参数:
//   - ctx: 上下文接口，通常为 echo.Context 类型，包含请求和响应对象
//
// 返回值:
//   - models.UserModel: 用户模型，包含用户信息（如用户名、密码、角色等）
//   - bool: 是否成功获取用户信息，true 表示成功，false 表示失败
//
// 功能特性:
//
//	提取用户信息: 从上下文中提取用户信息，用于身份验证和权限控制
//	支持多种认证方式: 支持 Cookie 认证、Session 认证等多种认证方式
//	支持用户权限验证: 根据用户角色和权限验证用户是否有权访问资源
//	支持用户角色管理: 支持用户角色的分配和管理，如管理员、普通用户等
//	支持用户会话管理: 支持用户会话的创建、更新和删除，维护用户登录状态
//
// 说明:
//
//	该方法调用基础适配器的 GetUser 方法，传入上下文和适配器实例
//	GetUser 方法会从上下文中提取认证信息（如 Cookie、Session 等）
//	然后查询数据库获取用户信息，返回用户模型和是否成功的标志
//	该方法通常在每个请求处理时被调用，用于验证用户身份和权限
//
// 注意事项:
//
//   - 该方法需要在请求处理前调用，确保用户信息正确获取
//   - 该方法会查询数据库，可能影响性能，建议在中间件中调用并缓存结果
//   - 该方法返回的用户模型可能为空，表示用户未登录或认证失败
//   - 该方法支持并发调用，但需要注意数据库连接的线程安全性
//   - 该方法应该在中间件中调用，确保在处理业务逻辑前验证用户身份
func (e *Echo) User(ctx interface{}) (models.UserModel, bool) {
	return e.GetUser(ctx, e) // 调用基础适配器的 GetUser 方法，传入上下文和适配器实例
}

// Use 实现了 Adapter.Use 方法
// 该方法用于将插件注册到 Echo 应用中，用于扩展 GoAdmin 的功能
//
// 参数:
//   - app: 应用接口，必须为 *echo.Echo 类型，表示 Echo 应用实例
//   - plugs: 插件列表，包含需要注册的所有插件，如管理后台插件、认证插件等
//
// 返回值:
//   - error: 错误信息，如果注册失败则返回错误，否则返回 nil
//
// 功能特性:
//
//	注册插件: 将插件注册到 Echo 应用中，用于扩展 GoAdmin 的功能
//	支持多个插件: 可以同时注册多个插件，实现功能模块化
//	支持插件路由注册: 插件可以注册自己的路由，实现自定义功能
//	支持插件中间件注册: 插件可以注册自己的中间件，实现请求预处理和后处理
//	支持插件配置管理: 插件可以有自己的配置，通过配置文件或环境变量设置
//
// 说明:
//
//	该方法调用基础适配器的 GetUse 方法，传入应用实例、插件列表和适配器实例
//	GetUse 方法会遍历插件列表，对每个插件进行初始化和注册
//	初始化过程包括设置插件配置、注册插件路由和中间件等
//	注册完成后，插件就可以在 Echo 应用了
//	该方法通常在应用启动时调用，确保所有插件在应用启动前完成注册
//
// 注意事项:
//
//   - 该方法需要在应用启动前调用，确保所有插件在应用启动前完成注册
//   - app 参数必须为 *echo.Echo 类型，否则会返回错误
//   - plugs 参数不能为空，至少包含一个插件
//   - 插件初始化顺序很重要，某些插件可能依赖其他插件
//   - 该方法支持并发调用，但通常只调用一次
func (e *Echo) Use(app interface{}, plugs []plugins.Plugin) error {
	return e.GetUse(app, plugs, e) // 调用基础适配器的 GetUse 方法，传入应用实例、插件列表和适配器实例
}

// Content 实现了 Adapter.Content 方法
// 该方法用于渲染管理面板内容，将 GoAdmin 的面板渲染为 HTML 响应
//
// 参数:
//   - ctx: 上下文接口，通常为 echo.Context 类型，包含请求和响应对象
//   - getPanelFn: 获取面板的函数，返回 types.Panel 类型的面板
//   - fn: 节点处理器，用于处理上下文中的节点，可以为空
//   - btns: 导航按钮列表，可变参数，可以为空
//
// 返回值:
//   - 无返回值
//
// 功能特性:
//
//	渲染管理面板: 将 GoAdmin 的面板渲染为 HTML 响应，显示管理后台界面
//	支持自定义面板内容: 通过 getPanelFn 函数可以自定义面板内容
//	支持自定义节点处理器: 通过 fn 函数可以自定义节点处理逻辑
//	支持自定义导航按钮: 通过 btns 参数可以自定义导航按钮
//	支持多种面板类型: 支持表格面板、表单面板、详情面板等多种面板类型
//
// 说明:
//
//	该方法调用基础适配器的 GetContent 方法，传入上下文、获取面板函数、适配器实例、导航按钮和节点处理器
//	GetContent 方法会调用 getPanelFn 函数获取面板内容
//	然后处理面板中的节点（如表格、表单、按钮等）
//	最后显示导航按钮，并将面板渲染为 HTML 响应
//	该方法通常在请求处理时被调用，用于渲染管理后台界面
//
// 注意事项:
//
//   - 该方法需要在请求处理中调用，确保在渲染前获取正确的面板内容
//   - getPanelFn 函数不能为空，必须返回有效的面板
//   - 节点处理器 fn 可以为空，表示不处理节点
//   - 导航按钮 btns 可以为空，表示不显示导航按钮
//   - 该方法支持并发调用，但需要注意响应对象的线程安全性
func (e *Echo) Content(ctx interface{}, getPanelFn types.GetPanelFn, fn context.NodeProcessor, btns ...types.Button) {
	e.GetContent(ctx, getPanelFn, e, btns, fn) // 调用基础适配器的 GetContent 方法，传入上下文、获取面板函数、适配器实例、导航按钮和节点处理器
}

// HandlerFunc 定义了处理函数的类型
// 该函数接收 Echo 上下文，返回面板和可能的错误
//
// 参数:
//   - ctx: Echo 上下文，包含请求和响应对象
//
// 返回值:
//   - types.Panel: 管理面板，包含面板内容、标题、按钮等信息
//   - error: 错误信息，如果处理失败则返回错误，否则返回 nil
//
// 功能特性:
//
//	支持自定义处理逻辑: 可以在处理函数中实现自定义的业务逻辑
//	支持返回自定义面板: 可以返回自定义的面板，实现自定义界面
//	支持错误处理: 可以返回错误，表示处理失败
//	支持异步处理: 可以在处理函数中启动 goroutine 进行异步处理
//
// 说明:
//
//	该类型定义了处理函数的签名，用于在 Echo 路由中使用 GoAdmin 的处理函数
//	处理函数接收 Echo 上下文，可以访问请求和响应对象
//	处理函数返回面板和错误，面板会被渲染为 HTML 响应
//	处理函数可以访问数据库、调用 API、执行业务逻辑等
//	该类型通常与 Content 辅助函数一起使用，将 HandlerFunc 转换为 echo.HandlerFunc
//
// 使用示例:
//
//	func myHandler(ctx echo.Context) (types.Panel, error) {
//	    // 访问数据库
//	    users, err := db.GetAllUsers()
//	    if err != nil {
//	        return nil, err
//	    }
//
//	    // 创建面板
//	    panel := types.NewPanel("用户列表")
//	    panel.AddTable(types.NewTable(users))
//
//	    return panel, nil
//	}
//
// 注意事项:
//
//   - 处理函数不能为空，必须返回有效的面板
//   - 处理函数应该是线程安全的，可以在多个 goroutine 中并发调用
//   - 处理函数应该尽快返回，避免阻塞请求处理
//   - 处理函数应该正确处理错误，避免未捕获的异常
//   - 处理函数应该避免修改传入的上下文，确保上下文的线程安全性
type HandlerFunc func(ctx echo.Context) (types.Panel, error)

// Content 是一个辅助函数，用于将 HandlerFunc 转换为 Echo 的处理函数
// 这样可以在 Echo 的路由中使用 GoAdmin 的处理函数
//
// 参数:
//   - handler: 处理函数，接收 Echo 上下文并返回面板
//
// 返回值:
//   - echo.HandlerFunc: Echo 处理器函数，可以在 Echo 路由中使用
//
// 功能特性:
//
//	转换 HandlerFunc 为 Handler: 将 GoAdmin 的处理函数转换为 Echo 的处理器函数
//	简化集成: 简化了 GoAdmin 与 Echo 的集成过程，无需手动处理上下文转换
//	支持在 Echo 路由中使用 GoAdmin 处理函数: 可以在 Echo 路由中直接使用 GoAdmin 的处理函数
//	支持自定义处理逻辑: 可以在处理函数中实现自定义的业务逻辑
//
// 说明:
//
//	该函数接收 HandlerFunc 类型的处理函数，返回 echo.HandlerFunc 类型的处理器函数
//	处理器函数内部调用 GoAdmin 的 engine.Content 方法，传入上下文和处理函数
//	engine.Content 方法会调用处理函数，获取面板并渲染为 HTML 响应
//	该函数通常在路由注册时使用，将 GoAdmin 的处理函数注册到 Echo 路由中
//
// 使用示例:
//
//	// 定义处理函数
//	func myHandler(ctx echo.Context) (types.Panel, error) {
//	    panel := types.NewPanel("用户列表")
//	    return panel, nil
//	}
//
//	// 使用 Content 辅助函数创建 Echo 处理器
//	e := echo.New()
//	e.GET("/admin", Content(myHandler))
//
// 注意事项:
//
//   - 处理函数不能为空，必须返回有效的面板
//   - 该函数应该在适当的 HTTP 方法中注册，如 GET、POST 等
//   - 该函数应该尽早注册，确保在应用启动前完成注册
//   - 该函数应该正确处理错误，避免未捕获的异常
//   - 该函数应该避免修改传入的上下文，确保上下文的线程安全性
func Content(handler HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error { // 返回 echo.HandlerFunc 类型的处理器函数
		engine.Content(ctx, func(ctx interface{}) (types.Panel, error) { // 调用 GoAdmin 的 engine.Content 方法，传入上下文和处理函数
			return handler(ctx.(echo.Context)) // 调用处理函数，传入 Echo 上下文
		})
		return nil // 返回 nil，表示处理成功
	}
}

// SetApp 实现了 Adapter.SetApp 方法
// 该方法用于设置 Echo 应用实例到适配器中，用于后续的路由注册和请求处理
//
// 参数:
//   - app: 应用接口，必须为 *echo.Echo 类型，表示 Echo 应用实例
//
// 返回值:
//   - error: 错误信息，如果参数类型不正确则返回错误，否则返回 nil
//
// 功能特性:
//
//	设置 Echo 应用实例: 将 Echo 应用实例设置到适配器中，用于后续的路由注册和请求处理
//	支持类型检查: 使用类型断言检查参数类型，确保参数为 *echo.Echo 类型
//	支持错误处理: 如果参数类型不正确，返回错误信息
//
// 说明:
//
//	该方法使用类型断言检查 app 参数是否为 *echo.Echo 类型
//	如果类型断言成功，将 Echo 应用实例存储到适配器的 app 字段中
//	如果类型断言失败，返回错误信息，表示参数类型错误
//	该方法通常在 Use 方法中被调用，用于设置 Echo 应用实例
//	设置后，适配器就可以使用 Echo 应用实例注册路由和处理器了
//
// 注意事项:
//
//   - app 参数必须为 *echo.Echo 类型，否则会返回错误
//   - 该方法应该在 Use 方法之前调用，确保在注册插件前设置 Echo 应用实例
//   - 该方法只能调用一次，重复调用会覆盖之前的设置
//   - 该方法支持并发调用，但通常只调用一次
//   - 返回错误后适配器无法正常工作，需要检查参数类型
func (e *Echo) SetApp(app interface{}) error {
	var (
		eng *echo.Echo // Echo 应用实例
		ok  bool       // 类型断言结果
	)
	// 使用类型断言检查 app 是否为 *echo.Echo 类型
	// ok 为 true 表示断言成功，eng 为转换后的值
	if eng, ok = app.(*echo.Echo); !ok { // 如果类型断言失败
		return errors.New("echo 适配器 SetApp: 参数类型错误") // 返回错误信息
	}
	e.app = eng // 将 Echo 应用实例存储到适配器的 app 字段中
	return nil  // 返回 nil，表示设置成功
}

// AddHandler 实现了 Adapter.AddHandler 方法
// 该方法用于向 Echo 应用添加路由处理器，将 GoAdmin 的路由注册到 Echo 的路由系统中
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
//
//	注册路由: 将 GoAdmin 的路由注册到 Echo 的路由系统中
//	支持多种 HTTP 方法: 支持 GET、POST、PUT、DELETE、HEAD、OPTIONS、PATCH 等 HTTP 方法
//	支持路由参数: 支持路由参数（如 :id、:name 等），会自动转换为 URL 查询参数
//	支持处理器链: 支持处理器链（中间件和处理器），实现请求预处理和后处理
//	支持上下文转换: 将 Echo 上下文转换为 GoAdmin 上下文，实现框架无关性
//	支持响应头复制: 将 GoAdmin 响应头复制到 Echo 响应中，确保响应头正确传递
//	支持响应体写入: 将 GoAdmin 响应体写入到 Echo 响应中，确保响应体正确传递
//
// 说明:
//
//	该方法调用 Echo 应用的 Add 方法，注册路由处理器
//	路由处理器接收 Echo 上下文，处理 HTTP 请求和响应
//	首先创建 GoAdmin 上下文，传入 HTTP 请求
//	然后获取 Echo 的路由参数并转换为 URL 查询参数
//	Echo 的路由参数格式为 :param，需要转换为 ?param=value 格式
//	执行处理器链，处理器链包含中间件和处理器
//	将 GoAdmin 响应头复制到 Echo 响应中
//	将 GoAdmin 响应体写入到 Echo 响应中
//	该方法通常在 Use 方法中被调用，用于注册插件的路由
//
// 注意事项:
//
//   - method 参数必须为大写，如 "GET"、"POST" 等
//   - path 参数不能为空，必须以 "/" 开头
//   - handlers 参数不能为空，至少包含一个处理器
//   - 路由参数格式为 :param，会自动转换为 URL 查询参数 ?param=value
//   - 处理器链执行顺序为从左到右，中间件先执行，处理器后执行
//   - 该方法支持并发调用，但需要注意路由注册的线程安全性
//   - 该方法会修改响应头和响应体，确保在调用前未写入响应
func (e *Echo) AddHandler(method, path string, handlers context.Handlers) {
	e.app.Add(strings.ToUpper(method), path, func(c echo.Context) error { // 调用 Echo 应用的 Add 方法，注册路由处理器

		// 创建 GoAdmin 上下文
		ctx := context.NewContext(c.Request()) // 创建 GoAdmin 上下文，传入 HTTP 请求

		// 将 Echo 的路由参数转换为 URL 查询参数
		// Echo 的路由参数格式为 :param，需要转换为 ?param=value 格式
		for _, key := range c.ParamNames() { // 遍历 Echo 的路由参数名称
			if c.Request().URL.RawQuery == "" { // 如果 URL 查询参数为空
				c.Request().URL.RawQuery += strings.ReplaceAll(key, ":", "") + "=" + c.Param(key) // 添加第一个查询参数，格式为 ?param=value
			} else { // 如果 URL 查询参数不为空
				c.Request().URL.RawQuery += "&" + strings.ReplaceAll(key, ":", "") + "=" + c.Param(key) // 添加后续查询参数，格式为 &param=value
			}
		}

		// 执行处理器链
		ctx.SetHandlers(handlers).Next() // 设置处理器链并执行，处理器链包含中间件和处理器
		// 将 GoAdmin 响应头复制到 Echo 响应中
		for key, head := range ctx.Response.Header { // 遍历 GoAdmin 响应头
			c.Response().Header().Set(key, head[0]) // 将响应头设置到 Echo 响应中
		}
		// 将响应体写入 Echo 响应
		if ctx.Response.Body != nil { // 如果 GoAdmin 响应体不为空
			buf := new(bytes.Buffer)                            // 创建字节缓冲区
			_, _ = buf.ReadFrom(ctx.Response.Body)              // 从 GoAdmin 响应体读取数据到缓冲区
			_ = c.String(ctx.Response.StatusCode, buf.String()) // 将缓冲区数据写入 Echo 响应
		} else { // 如果 GoAdmin 响应体为空
			c.Response().WriteHeader(ctx.Response.StatusCode) // 只设置 Echo 响应的状态码
		}
		return nil // 返回 nil，表示处理成功
	})
}

// Name 实现了 Adapter.Name 方法
// 该方法返回适配器的名称，用于标识适配器类型
//
// 返回值:
//   - string: 适配器名称，固定为 "echo"
//
// 功能特性:
//
//	返回适配器名称: 返回适配器的名称，用于标识适配器类型
//	用于标识适配器类型: 可以在日志、配置文件等地方使用适配器名称
//	用于日志记录和调试: 可以在日志中记录适配器名称，便于调试
//
// 说明:
//
//	该方法返回适配器的名称，固定为 "echo"
//	适配器名称用于标识适配器类型，可以在日志中记录
//	适配器名称也可以在配置文件中使用，指定使用哪个适配器
//	该方法通常在适配器注册时被调用，用于标识适配器类型
//
// 注意事项:
//
//   - 适配器名称必须与适配器类型一致，不能随意修改
//   - 适配器名称不能为空，必须返回有效的字符串
//   - 适配器名称应该唯一，不能与其他适配器名称冲突
func (*Echo) Name() string {
	return "echo" // 返回适配器名称，固定为 "echo"
}

// SetContext 实现了 Adapter.SetContext 方法
// 该方法用于设置当前请求的上下文，用于后续的请求处理
//
// 参数:
//   - contextInterface: 上下文接口，必须为 echo.Context 类型，包含请求和响应对象
//
// 返回值:
//   - adapter.WebFrameWork: 返回设置了上下文的新适配器实例
//
// 功能特性:
//
//	设置当前请求的上下文: 将上下文设置到适配器中，用于后续的请求处理
//	支持类型检查: 使用类型断言检查参数类型，确保参数为 echo.Context 类型
//	支持错误处理: 如果参数类型不正确，会 panic
//	返回新的适配器实例: 返回一个新的适配器实例，包含设置的上下文
//
// 说明:
//
//	该方法使用类型断言检查 contextInterface 参数是否为 echo.Context 类型
//	如果类型断言成功，创建一个新的适配器实例，设置上下文到 ctx 字段中
//	如果类型断言失败，panic，表示参数类型错误
//	该方法通常在每个请求处理时被调用，用于设置当前请求的上下文
//	设置后，适配器就可以使用该上下文访问请求和响应对象了
//
// 注意事项:
//
//   - contextInterface 参数必须为 echo.Context 类型，否则会 panic
//   - 该方法返回一个新的适配器实例，不是修改当前实例
//   - 该方法应该在每个请求处理时调用，确保上下文的正确性
//   - 该方法支持并发调用，但需要注意上下文的线程安全性
//   - panic 后请求处理会中断，需要确保参数类型正确
func (*Echo) SetContext(contextInterface interface{}) adapter.WebFrameWork {
	var (
		ctx echo.Context // Echo 上下文
		ok  bool         // 类型断言结果
	)
	// 使用类型断言检查 contextInterface 是否为 echo.Context 类型
	if ctx, ok = contextInterface.(echo.Context); !ok { // 如果类型断言失败
		panic("echo 适配器 SetContext: 参数类型错误") // panic，表示参数类型错误
	}
	return &Echo{ctx: ctx} // 返回新的适配器实例，包含设置的上下文
}

// Redirect 实现了 Adapter.Redirect 方法
// 该方法用于重定向到登录页面，使用 HTTP 302 状态码进行临时重定向
//
// 返回值:
//   - 无返回值
//
// 功能特性:
//
//	重定向到登录页面: 将用户重定向到登录页面，用于用户未认证时
//	使用 HTTP 302 状态码: 使用 HTTP 302 状态码进行临时重定向
//	支持自定义登录页面 URL: 登录页面 URL 通过 config.GetLoginUrl() 获取
//
// 说明:
//
//	该方法调用 Echo 上下文的 Redirect 方法，使用 HTTP 302 状态码
//	登录页面 URL 通过 config.GetLoginUrl() 获取，默认为 "/login"
//	重定向 URL 通过 config.Url() 生成，包含域名和路径
//	该方法通常在用户未认证时被调用，将用户重定向到登录页面
//
// 注意事项:
//
//   - 该方法会修改响应头和响应体，确保在调用前未写入响应
//   - 该方法会终止请求处理，调用后不应该继续处理请求
//   - 该方法应该在用户未认证时调用，如 Cookie 失效或会话过期
//   - 该方法支持并发调用，但需要注意响应对象的线程安全性
func (e *Echo) Redirect() {
	_ = e.ctx.Redirect(http.StatusFound, config.Url(config.GetLoginUrl())) // 调用 Echo 上下文的 Redirect 方法，重定向到登录页面
}

// SetContentType 实现了 Adapter.SetContentType 方法
// 该方法用于设置响应的 Content-Type 头，Content-Type 由 HTMLContentType() 方法确定
//
// 返回值:
//   - 无返回值
//
// 功能特性:
//
//	设置响应的 Content-Type 头: 设置响应的 Content-Type 头，指定响应内容的类型
//	支持自定义 Content-Type: Content-Type 由 HTMLContentType() 方法确定
//
// 说明:
//
//	该方法调用 Echo 响应头的 Set 方法，设置 Content-Type 头
//	Content-Type 由 HTMLContentType() 方法确定，通常为 "text/html; charset=utf-8"
//	该方法通常在渲染 HTML 响应前被调用，确保浏览器正确解析响应内容
//
// 注意事项:
//
//   - 该方法会修改响应头，确保在调用前未写入响应
//   - 该方法应该在渲染 HTML 响应前调用，确保 Content-Type 正确
//   - 该方法支持并发调用，但需要注意响应对象的线程安全性
func (e *Echo) SetContentType() {
	e.ctx.Response().Header().Set("Content-Type", e.HTMLContentType()) // 设置响应的 Content-Type 头
}

// Write 实现了 Adapter.Write 方法
// 该方法用于将响应体写入到响应中
//
// 参数:
//   - body: 要写入的响应体字节数组，通常为 HTML、JSON 等数据
//
// 返回值:
//   - 无返回值
//
// 功能特性:
//
//	将响应体写入到响应中: 将响应体字节数组写入到响应中
//	支持写入任意字节数组: 可以写入 HTML、JSON、XML 等任意格式的数据
//
// 说明:
//
//	该方法调用 Echo 响应的 WriteHeader 和 Write 方法
//	首先设置响应状态码为 200（OK）
//	然后写入响应体，响应体可以是 HTML、JSON 等数据
//	该方法通常在渲染完成后被调用，将响应体写入到响应中
//
// 注意事项:
//
//   - 该方法会修改响应体，确保在调用前未写入响应
//   - 该方法应该在渲染完成后调用，确保响应体正确
//   - 该方法支持并发调用，但需要注意响应对象的线程安全性
//   - 该方法会忽略写入错误，如果写入失败不会报错
func (e *Echo) Write(body []byte) {
	e.ctx.Response().WriteHeader(http.StatusOK) // 设置响应状态码为 200（OK）
	_, _ = e.ctx.Response().Write(body)         // 写入响应体
}

// GetCookie 实现了 Adapter.GetCookie 方法
// 该方法用于从请求中获取认证 Cookie，用于用户身份验证
//
// 返回值:
//   - string: Cookie 值，如果获取成功则返回 Cookie 值，否则返回空字符串
//   - error: 错误信息，如果获取失败则返回错误，否则返回 nil
//
// 功能特性:
//
//	从请求中获取认证 Cookie: 从请求中获取认证 Cookie，用于用户身份验证
//	支持自定义 Cookie 键名: Cookie 键名通过 CookieKey() 方法获取，默认为 "admin_cookie"
//	支持错误处理: 如果获取失败，返回错误信息
//
// 说明:
//
//	该方法调用 Echo 请求的 Cookie() 方法，获取 Cookie 值
//	Cookie 键名通过 CookieKey() 方法获取，默认为 "admin_cookie"
//	返回 Cookie 值和错误，如果获取成功，错误为 nil
//	该方法通常在 User 方法中被调用，用于获取用户身份信息
//
// 注意事项:
//
//   - 返回的 Cookie 值可能为空，表示 Cookie 不存在或获取失败
//   - 返回的错误表示 Cookie 不存在或获取失败，需要检查错误
//   - 该方法支持并发调用，但需要注意请求对象的线程安全性
//   - 该方法应该在 User 方法之前调用，确保获取正确的 Cookie 值
func (e *Echo) GetCookie() (string, error) {
	cookie, err := e.ctx.Cookie(e.CookieKey()) // 获取 Cookie 值
	if err != nil {                            // 如果获取失败
		return "", err // 返回空字符串和错误
	}
	return cookie.Value, err // 返回 Cookie 值和错误
}

// Lang 实现了 Adapter.Lang 方法
// 该方法用于从 URL 查询参数中获取语言设置，用于多语言支持
//
// 返回值:
//   - string: 语言代码，如 "zh-CN"、"en-US" 等，如果未设置则返回空字符串
//
// 功能特性:
//
//	从 URL 查询参数中获取语言设置: 从 URL 查询参数中获取语言设置，用于多语言支持
//	支持自定义语言参数名: 语言参数名为 "__ga_lang"，遵循 GoAdmin 的命名规范
//	支持多语言切换: 支持多语言切换，通过 URL 查询参数 __ga_lang 指定语言
//
// 说明:
//
//	该方法调用 Echo 请求的 URL.Query().Get() 方法，获取语言参数值
//	语言参数名为 "__ga_lang"，语言代码格式为 "zh-CN"、"en-US" 等
//	语言代码遵循 RFC 5646 标准，如 "zh-CN" 表示简体中文
//	该方法通常在渲染模板前被调用，用于设置模板的语言
//
// 注意事项:
//
//   - 返回的语言代码可能为空，表示未设置语言参数
//   - 该方法支持并发调用，但需要注意请求对象的线程安全性
//   - 该方法应该在渲染模板前调用，确保语言设置正确
//   - 返回的语言代码应该有效，否则可能导致模板渲染错误
func (e *Echo) Lang() string {
	return e.ctx.Request().URL.Query().Get("__ga_lang") // 获取语言参数值
}

// Path 实现了 Adapter.Path 方法
// 该方法用于获取当前请求的路径，用于路由匹配和日志记录
//
// 返回值:
//   - string: 请求路径，如 "/admin/dashboard"，如果路径为空则返回 "/"
//
// 功能特性:
//
//	获取当前请求的路径: 获取当前请求的路径，不包含查询参数
//	支持路由匹配: 可以用于路由匹配和路径解析
//	支持日志记录: 可以在日志中记录请求路径，便于调试
//
// 说明:
//
//	该方法调用 Echo 请求的 URL.Path 属性，获取请求路径
//	路径不包含查询参数，如 "/admin/dashboard"
//	路径以 "/" 开头，如果路径为空则返回 "/"
//	该方法通常在请求处理时被调用，用于路由匹配和日志记录
//
// 注意事项:
//
//   - 返回的路径不包含查询参数，如 "?page=1&limit=10"
//   - 返回的路径以 "/" 开头，如 "/admin/dashboard"
//   - 该方法支持并发调用，但需要注意请求对象的线程安全性
//   - 该方法应该在请求处理中调用，确保路径正确
func (e *Echo) Path() string {
	return e.ctx.Request().URL.Path // 获取请求路径
}

// Method 实现了 Adapter.Method 方法
// 该方法用于获取当前请求的 HTTP 方法，用于请求类型判断和日志记录
//
// 返回值:
//   - string: HTTP 方法，如 "GET"、"POST"、"PUT"、"DELETE" 等
//
// 功能特性:
//
//	获取当前请求的 HTTP 方法: 获取当前请求的 HTTP 方法，用于请求类型判断
//	支持所有 HTTP 方法: 支持 GET、POST、PUT、DELETE、PATCH、HEAD、OPTIONS 等 HTTP 方法
//	支持请求类型判断: 可以根据 HTTP 方法判断请求类型，执行不同的处理逻辑
//	支持日志记录: 可以在日志中记录 HTTP 方法，便于调试
//
// 说明:
//
//	该方法调用 Echo 请求的 Method 属性，获取 HTTP 方法
//	HTTP 方法包括 GET、POST、PUT、DELETE、PATCH、HEAD、OPTIONS 等
//	HTTP 方法为大写字母，如 "GET"、"POST" 等
//	该方法通常在请求处理时被调用，用于请求类型判断和日志记录
//
// 注意事项:
//
//   - 返回的 HTTP 方法为大写字母，如 "GET"、"POST"
//   - 该方法支持所有 HTTP 方法，包括自定义方法
//   - 该方法支持并发调用，但需要注意请求对象的线程安全性
//   - 该方法应该在请求处理中调用，确保 HTTP 方法正确
func (e *Echo) Method() string {
	return e.ctx.Request().Method // 获取 HTTP 方法
}

// FormParam 实现了 Adapter.FormParam 方法
// 该方法用于解析并获取表单参数，用于表单数据处理
//
// 返回值:
//   - url.Values: 表单参数的键值对集合，支持多值字段
//
// 功能特性:
//
//	解析并获取表单参数: 解析并获取表单参数，用于表单数据处理
//	支持多种表单类型: 支持 application/x-www-form-urlencoded 和 multipart/form-data 表单
//	支持文件上传: 支持 multipart/form-data 表单，可以上传文件
//	解析的最大内存限制为 32MB: 解析的最大内存限制为 32MB，超过限制会写入临时文件
//
// 说明:
//
//	该方法调用 Echo 请求的 ParseMultipartForm 方法，解析表单参数
//	解析的最大内存限制为 32MB（32 << 20）
//	返回表单参数的键值对集合，类型为 url.Values
//	支持多值字段，一个键可以对应多个值
//	该方法通常在表单提交时被调用，用于获取表单数据
//
// 注意事项:
//
//   - 该方法会解析表单参数，可能影响性能，建议只调用一次
//   - 解析的最大内存限制为 32MB，超过限制会写入临时文件
//   - 该方法会忽略解析错误，如果解析失败不会报错
//   - 该方法支持并发调用，但需要注意请求对象的线程安全性
//   - 该方法应该只调用一次，重复调用会覆盖之前的解析结果
func (e *Echo) FormParam() url.Values {
	_ = e.ctx.Request().ParseMultipartForm(32 << 20) // 解析表单参数，最大内存限制为 32MB
	return e.ctx.Request().PostForm                  // 返回表单参数的键值对集合
}

// IsPjax 实现了 Adapter.IsPjax 方法
// 该方法用于检查当前请求是否为 PJAX 请求，PJAX 是一种使用 AJAX 技术实现页面部分更新的技术
//
// 返回值:
//   - bool: 如果是 PJAX 请求则返回 true，否则返回 false
//
// 功能特性:
//
//	检查当前请求是否为 PJAX 请求: 检查当前请求是否为 PJAX 请求
//	支持页面部分更新: PJAX 是一种使用 AJAX 技术实现页面部分更新的技术
//	支持无刷新导航: PJAX 支持无刷新导航，提升用户体验
//
// 说明:
//
//	该方法调用 Echo 请求的 Header.Get() 方法，获取 PJAX 请求头
//	PJAX 请求头为 "X-PJAX"，如果值为 "true" 则返回 true
//	PJAX 是一种使用 AJAX 技术实现页面部分更新的技术
//	该方法通常在渲染模板前被调用，用于判断是否为 PJAX 请求
//
// 注意事项:
//
//   - 返回的布尔值表示是否为 PJAX 请求
//   - 该方法支持并发调用，但需要注意请求对象的线程安全性
//   - 该方法应该在渲染模板前调用，确保 PJAX 请求正确处理
//   - 该方法应该与前端 PJAX 库配合使用，如 jQuery PJAX 插件
func (e *Echo) IsPjax() bool {
	return e.ctx.Request().Header.Get(constant.PjaxHeader) == "true" // 检查 PJAX 请求头是否为 "true"
}

// Query 实现了 Adapter.Query 方法
// 该方法用于获取 URL 查询参数，用于查询参数处理和分页等
//
// 返回值:
//   - url.Values: 查询参数的键值对集合，支持多值字段
//
// 功能特性:
//
//	获取 URL 查询参数: 获取 URL 查询参数，用于查询参数处理和分页等
//	支持多值字段: 支持多值字段，一个键可以对应多个值
//	支持分页、排序、筛选等功能: 支持分页、排序、筛选等功能，通过 URL 查询参数传递
//
// 说明:
//
//	该方法调用 Echo 请求的 URL.Query() 方法，获取查询参数
//	查询参数格式为 ?key1=value1&key2=value2
//	返回查询参数的键值对集合，类型为 url.Values
//	支持多值字段，一个键可以对应多个值
//	该方法通常在数据查询时被调用，用于获取查询参数
//
// 注意事项:
//
//   - 返回的查询参数可能为空，表示 URL 不包含查询参数
//   - 该方法支持多值字段，一个键可以对应多个值
//   - 该方法支持并发调用，但需要注意请求对象的线程安全性
//   - 该方法应该在数据查询时调用，确保查询参数正确
func (e *Echo) Query() url.Values {
	return e.ctx.Request().URL.Query() // 获取查询参数
}

// Request 实现了 Adapter.Request 方法
// 该方法用于获取原始的 HTTP 请求对象，用于底层请求处理
//
// 返回值:
//   - *http.Request: HTTP 请求对象指针，包含请求的所有信息
//
// 功能特性:
//
//	获取原始的 HTTP 请求对象: 获取原始的 HTTP 请求对象，用于底层请求处理
//	支持底层请求处理: 可以访问请求的所有信息，如 URL、方法、头、体等
//	支持自定义请求处理逻辑: 可以在底层请求处理中实现自定义逻辑
//
// 说明:
//
//	该方法调用 Echo 上下文的 Request() 方法，获取原始的 HTTP 请求对象
//	HTTP 请求对象包含请求的所有信息，如 URL、方法、头、体等
//	该方法通常在底层请求处理时被调用，用于自定义请求处理逻辑
//
// 注意事项:
//
//   - 返回的 HTTP 请求对象可能为空，表示请求未初始化
//   - 返回的 HTTP 请求对象是原始对象，修改会影响请求处理
//   - 该方法支持并发调用，但需要注意请求对象的线程安全性
//   - 该方法应该在底层请求处理时调用，确保请求对象正确获取
//   - 该方法应该谨慎使用，避免修改请求对象导致意外行为
func (e *Echo) Request() *http.Request {
	return e.ctx.Request() // 调用 Echo 上下文的 Request() 方法，获取原始的 HTTP 请求对象
}

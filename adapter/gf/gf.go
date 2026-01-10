// GoFrame (GF) 框架适配器包
//
// 本包实现了 GoAdmin 与 GoFrame (GF) Web 框架的集成适配器，使 GoAdmin 管理后台能够在 GoFrame 框架中运行
// GoFrame (GF) 是一个模块化、高性能、生产级的 Go 应用开发框架，本适配器实现了 GoAdmin 适配器接口的所有方法
//
// 核心概念:
//
//	适配器模式: Gf 结构体实现了 GoAdmin 的适配器接口，作为 GoFrame 框架和 GoAdmin 管理后台之间的桥梁
//	上下文转换: 将 GoFrame 请求上下文（*ghttp.Request）转换为 GoAdmin 上下文（*context.Context），实现框架无关性
//	路由集成: 将 GoAdmin 的路由注册到 GoFrame 的路由系统中，支持多种 HTTP 方法和路由参数
//	处理器链: 支持处理器链（中间件和处理器），实现请求预处理和后处理
//	请求处理: 处理 HTTP 请求和响应，包括请求参数解析、响应头设置、响应体写入等
//	Cookie 认证: 通过 Cookie 进行用户身份验证和会话管理
//	PJAX 支持: 支持 PJAX（PushState + AJAX）技术，实现页面部分更新
//	多语言: 支持多语言切换，通过 URL 查询参数 __ga_lang 指定语言
//
// 技术栈:
//
//	GoFrame (GF): 模块化、高性能、生产级的 Go 应用开发框架，提供路由、中间件、上下文等功能
//	GoAdmin: Go 后台管理框架，提供管理后台、表单、表格、权限控制等功能
//	Go 标准库: bytes、errors、net/http、net/url、regexp、strings 等
//
// 数据库支持:
//
//	MySQL: 通过 GoAdmin 的数据库驱动支持
//	PostgreSQL: 通过 GoAdmin 的数据库驱动支持
//	SQLite: 通过 GoAdmin 的数据库驱动支持
//	MSSQL: 通过 GoAdmin 的数据库驱动支持
//
// 配置说明:
//
//	使用 config.SetCfg 设置 GoAdmin 配置，包括数据库连接、语言设置等
//	使用 config.GetLoginUrl 获取登录页面 URL
//	使用 config.Url 生成完整的 URL
//	使用 constant.EditPKKey 获取编辑主键的常量
//	使用 constant.PjaxHeader 获取 PJAX 请求头的常量
//
// 使用示例:
//
//	package main
//
//	import (
//	    "github.com/gogf/gf/net/ghttp"
//	    gfadapter "github.com/purpose168/GoAdmin/adapter/gf"
//	    "github.com/purpose168/GoAdmin/modules/config"
//	    "github.com/purpose168/GoAdmin/plugins/admin"
//	)
//
//	func main() {
//	    // 初始化 GoFrame 服务器
//	    s := ghttp.GetServer()
//
//	    // 设置 GoAdmin 配置
//	    config.SetCfg(config.Config{
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
//	        UrlPrefix: "admin",
//	    })
//
//	    // 使用 GF 适配器
//	    adapter := gfadapter.New()
//	    adapter.SetApp(s)
//	    adapter.Use(s, []plugins.Plugin{admin.NewAdmin(generateTables)})
//
//	    // 启动 GoFrame 服务器
//	    s.Run()
//	}
//
// 注意事项:
//
//   - GoFrame 版本要求: 需要使用 GoFrame 1.x 版本
//   - 路由参数格式: GoFrame 的路由参数格式为 :param，适配器会自动转换为标准格式 {param}
//   - 路由绑定格式: GoFrame 使用 "METHOD:path" 格式来绑定路由，如 "GET:/admin"
//   - 表单解析: GoFrame 的 Form 字段已经包含了解析后的表单参数，无需手动解析
//   - Cookie 认证: Cookie 认证需要在配置中设置 CookieKey，默认为 "goadmin_session"
//   - PJAX 请求: PJAX 请求需要在请求头中设置 X-PJAX: true
//   - 语言设置: 语言通过 URL 查询参数 __ga_lang 指定，格式为 "zh-CN"、"en-US" 等
//   - Content-Type: 默认 Content-Type 为 "text/html; charset=utf-8"，可以通过 HTMLContentType() 方法修改
//   - 响应状态码和响应体: 响应状态码和响应体通过 Response.WriteStatus() 方法设置
//   - 请求对象转换: GoFrame 请求对象（*ghttp.Request.Request）已经是标准 HTTP 请求对象，无需转换
package gf

import (
	"bytes"    // 字节缓冲区操作，提供字节缓冲区的读写功能
	"errors"   // 错误处理，提供错误创建和处理功能
	"net/http" // HTTP 包，提供 HTTP 客户端和服务器功能
	"net/url"  // URL 解析和查询，提供 URL 解析和查询参数处理功能
	"regexp"   // 正则表达式，提供字符串匹配和替换功能
	"strings"  // 字符串操作，提供字符串处理和转换功能

	"github.com/gogf/gf/net/ghttp"                                 // GoFrame HTTP 服务器包，提供路由、中间件、上下文等功能
	"github.com/purpose168/GoAdmin/adapter"                        // GoAdmin 适配器包，提供适配器接口和基础适配器
	"github.com/purpose168/GoAdmin/context"                        // GoAdmin 上下文包，提供请求上下文和处理器链功能
	"github.com/purpose168/GoAdmin/engine"                         // GoAdmin 引擎包，提供核心引擎和适配器注册功能
	"github.com/purpose168/GoAdmin/modules/config"                 // GoAdmin 配置模块，提供配置管理和 URL 生成功能
	"github.com/purpose168/GoAdmin/modules/utils"                  // GoAdmin 工具模块，提供字符串处理等工具函数
	"github.com/purpose168/GoAdmin/plugins"                        // GoAdmin 插件包，提供插件接口和插件管理功能
	"github.com/purpose168/GoAdmin/plugins/admin/models"           // GoAdmin 管理员模型包，提供用户模型和数据访问功能
	"github.com/purpose168/GoAdmin/plugins/admin/modules/constant" // GoAdmin 常量模块，提供框架常量定义（如 PjaxHeader、EditPKKey 等）
	"github.com/purpose168/GoAdmin/template/types"                 // GoAdmin 模板类型包，提供面板、按钮等类型定义
)

// Gf 结构体实现了 GoAdmin 的适配器接口
// 它作为 GoFrame (GF) 框架和 GoAdmin 管理后台之间的桥梁，实现了适配器接口的所有方法
// 嵌入了 adapter.BaseAdapter 以获得基础适配器功能（如 GetUser、GetUse、GetContent 等）
// ctx 字段存储当前的 GoFrame 请求上下文，用于访问请求和响应对象
// app 字段存储 GoFrame 服务器实例，用于注册路由和处理器
type Gf struct {
	adapter.BaseAdapter
	ctx *ghttp.Request // GoFrame 请求上下文，用于访问请求和响应对象
	app *ghttp.Server  // GoFrame 服务器实例，用于注册路由和处理器
}

// init 函数在包导入时自动执行
// 它将 Gf 适配器注册到 GoAdmin 引擎中
// 这样 GoAdmin 就能够识别和使用 GoFrame 框架
//
// 功能特性:
//
//	自动执行: init 函数在包导入时自动执行，无需手动调用
//	适配器注册: 将 Gf 适配器注册到 GoAdmin 引擎中，使 GoAdmin 能够识别和使用 GoFrame 框架
//	支持多个适配器: GoAdmin 支持同时注册多个适配器，可以在不同框架中使用
//	线程安全: 注册操作是线程安全的，可以在多个 goroutine 中并发调用
//
// 说明:
//
//	调用 engine.Register 方法，将 Gf 适配器注册到 GoAdmin 引擎中
//	注册后，GoAdmin 就能够识别和使用 GoFrame 框架
//	该函数在包导入时自动执行，无需手动调用
//	该函数在 main 函数之前执行，确保适配器在应用启动前注册
//
// 注意事项:
//
//   - 不能手动调用 init 函数，它会在包导入时自动执行
//   - init 函数的执行顺序是不确定的，如果多个包都有 init 函数，执行顺序是不确定的
//   - init 函数应该在 main 函数之前完成执行
//   - 该函数支持并发调用，但 GoAdmin 引擎内部已经做了线程安全处理
func init() {
	engine.Register(new(Gf)) // 注册 Gf 适配器到 GoAdmin 引擎中
}

// User 实现了 Adapter.User 方法
// 该方法从请求上下文中获取当前登录的用户信息，用于身份验证和权限控制
//
// 参数:
//   - ctx: 请求上下文接口，实际类型为 *ghttp.Request
//
// 返回值:
//   - models.UserModel: 用户模型，包含用户信息（如用户名、邮箱、角色等）
//   - bool: 是否成功获取用户信息，true 表示成功，false 表示失败
//
// 功能特性:
//
//	提取用户信息: 从请求上下文中提取用户信息，用于身份验证和权限控制
//	支持多种认证方式: 支持 Cookie 认证、Session 认证等多种认证方式
//	支持用户权限验证: 支持用户权限验证，确保用户只能访问有权限的资源
//	支持用户角色管理: 支持用户角色管理，实现基于角色的访问控制（RBAC）
//	支持用户会话管理: 支持用户会话管理，实现会话超时和会话续期
//
// 说明:
//
//	调用基础适配器的 GetUser 方法，传入请求上下文和适配器实例
//	GetUser 方法会从请求上下文中提取认证信息（如 Cookie、Session 等）
//	然后查询数据库获取用户信息，返回用户模型和是否成功获取用户信息
//	该方法通常在中间件中被调用，用于验证用户身份和权限
//
// 注意事项:
//
//   - 需要在请求处理前调用，通常在中间件中调用
//   - 会查询数据库，可能会影响性能，建议在中间件中缓存用户信息
//   - 如果用户未登录，返回的用户模型为空，bool 值为 false
//   - 建议在中间件中调用，确保用户信息在整个请求处理过程中可用
//   - 该方法支持并发调用，但需要注意数据库连接池的配置
func (gf *Gf) User(ctx interface{}) (models.UserModel, bool) {
	return gf.GetUser(ctx, gf) // 调用基础适配器的 GetUser 方法
}

// Use 实现了 Adapter.Use 方法
// 该方法用于初始化并使用 GoAdmin 插件，将插件注册到 GoFrame 应用中
//
// 参数:
//   - app: Web 应用实例，实际类型为 *ghttp.Server
//   - plugs: 插件列表，包含要加载的 GoAdmin 插件（如 Admin 插件）
//
// 返回值:
//   - error: 初始化过程中的错误，成功则为 nil，失败则为错误信息
//
// 功能特性:
//
//	注册插件: 将插件注册到 GoFrame 应用中，用于扩展 GoAdmin 的功能
//	支持多个插件: 支持同时注册多个插件，实现功能组合
//	支持插件路由注册: 支持插件路由注册，将插件的路由注册到 GoFrame 的路由系统中
//	支持插件中间件注册: 支持插件中间件注册，将插件的中间件注册到 GoFrame 的处理器链中
//	支持插件配置管理: 支持插件配置管理，将插件的配置存储到 GoAdmin 的配置系统中
//
// 说明:
//
//	调用基础适配器的 GetUse 方法，传入应用实例、插件列表和适配器实例
//	GetUse 方法会遍历插件列表，初始化每个插件
//	然后调用插件的 Init 方法，传入应用实例和适配器实例
//	插件会注册自己的路由和中间件到 GoFrame 应用中
//	插件配置会存储到 GoAdmin 的配置系统中
//	该方法通常在应用启动时调用，用于初始化 GoAdmin 插件
//
// 注意事项:
//
//   - 需要在应用启动前调用，通常在 main 函数中调用
//   - app 参数必须为 *ghttp.Server 类型，否则会返回错误
//   - plugs 参数不能为空，至少包含一个插件
//   - 插件初始化顺序很重要，建议先初始化基础插件，再初始化业务插件
//   - 该方法支持并发调用，但需要注意插件初始化的线程安全性
func (gf *Gf) Use(app interface{}, plugs []plugins.Plugin) error {
	return gf.GetUse(app, plugs, gf) // 调用基础适配器的 GetUse 方法
}

// Content 实现了 Adapter.Content 方法
// 该方法用于渲染管理面板内容，将 GoAdmin 的面板渲染为 HTML 响应
//
// 参数:
//   - ctx: 请求上下文接口，实际类型为 *ghttp.Request
//   - getPanelFn: 获取面板的函数，用于生成管理面板（如表格面板、表单面板等）
//   - fn: 节点处理器，用于处理面板中的节点（如表格行、表单字段等）
//   - btns: 导航按钮列表，用于面板上的操作按钮（如新增、编辑、删除等）
//
// 返回值:
//   - 无返回值
//
// 功能特性:
//
//	渲染管理面板: 将 GoAdmin 的面板渲染为 HTML 响应，用于显示管理后台界面
//	支持自定义面板内容: 支持自定义面板内容，如表格、表单、图表等
//	支持自定义节点处理器: 支持自定义节点处理器，用于处理面板中的节点（如表格行、表单字段等）
//	支持自定义导航按钮: 支持自定义导航按钮，用于面板上的操作按钮（如新增、编辑、删除等）
//	支持多种面板类型: 支持多种面板类型，如表格面板、表单面板、详情面板等
//
// 说明:
//
//	调用基础适配器的 GetContent 方法，传入请求上下文、获取面板的函数、适配器实例、导航按钮列表和节点处理器
//	GetContent 方法会调用 getPanelFn 函数获取面板内容
//	然后处理面板中的节点，调用节点处理器 fn 处理每个节点
//	最后显示导航按钮，将面板渲染为 HTML 响应
//	该方法通常在请求处理中被调用，用于渲染管理后台界面
//
// 注意事项:
//
//   - 需要在请求处理中调用，通常在路由处理器中调用
//   - getPanelFn 函数不能为空，必须返回有效的面板
//   - 节点处理器 fn 可以为空，如果为空则不处理节点
//   - 导航按钮 btns 可以为空，如果为空则不显示导航按钮
//   - 该方法支持并发调用，但需要注意面板渲染的线程安全性
func (gf *Gf) Content(ctx interface{}, getPanelFn types.GetPanelFn, fn context.NodeProcessor, btns ...types.Button) {
	gf.GetContent(ctx, getPanelFn, gf, btns, fn) // 调用基础适配器的 GetContent 方法
}

// HandlerFunc 定义了 GoFrame 框架的处理器函数类型
// 它接收 GoFrame 请求对象并返回面板和可能的错误，用于自定义处理逻辑
//
// 参数:
//   - ctx: GoFrame 请求对象，用于访问请求和响应对象
//
// 返回值:
//   - types.Panel: 管理面板，包含面板内容（如表格、表单等）
//   - error: 处理过程中的错误，成功则为 nil，失败则为错误信息
//
// 功能特性:
//
//	支持自定义处理逻辑: 支持自定义处理逻辑，实现业务逻辑处理
//	支持返回自定义面板: 支持返回自定义面板，如表格面板、表单面板等
//	支持错误处理: 支持错误处理，返回错误信息
//	支持异步处理: 支持异步处理，使用 goroutine 实现异步任务
//
// 说明:
//
//	定义处理函数的签名，接收 GoFrame 请求对象，返回面板和错误
//	处理函数可以访问数据库、调用 API、执行业务逻辑等
//	处理函数返回的面板会被渲染为 HTML 响应
//	如果处理函数返回错误，错误会被记录到日志中
//
// 使用示例:
//
//	func myHandler(ctx *ghttp.Request) (types.Panel, error) {
//	    // 访问数据库
//	    users, err := db.QueryUsers()
//	    if err != nil {
//	        return nil, err
//	    }
//
//	    // 返回表格面板
//	    return table.NewTable(users), nil
//	}
//
//	// 使用 Content 辅助函数创建 GoFrame 处理器
//	s.BindHandler("GET:/admin", gf.Content(myHandler))
//
// 注意事项:
//
//   - 处理函数不能为空，必须返回有效的面板
//   - 处理函数应该是线程安全的，可以在多个 goroutine 中并发调用
//   - 处理函数应该尽快返回，避免阻塞请求处理
//   - 处理函数应该正确处理错误，返回错误信息
//   - 处理函数应该避免修改传入的上下文，避免影响后续处理
type HandlerFunc func(ctx *ghttp.Request) (types.Panel, error)

// Content 是一个辅助函数，用于创建 GoFrame 处理器函数
// 该处理器将 GoFrame 的请求处理与 GoAdmin 的内容渲染集成，简化集成过程
//
// 参数:
//   - handler: 处理函数，用于生成管理面板
//
// 返回值:
//   - ghttp.HandlerFunc: GoFrame 处理器函数类型
//
// 功能特性:
//
//	转换 HandlerFunc 为 HandlerFunc: 将 HandlerFunc 类型的处理函数转换为 ghttp.HandlerFunc 类型的处理器函数
//	简化集成: 简化集成过程，无需手动创建处理器
//	支持在 GoFrame 路由中使用 GoAdmin 处理函数: 支持在 GoFrame 路由中使用 GoAdmin 处理函数
//	支持自定义处理逻辑: 支持自定义处理逻辑，实现业务逻辑处理
//
// 说明:
//
//	接收 HandlerFunc 类型的处理函数，返回 ghttp.HandlerFunc 类型的处理器函数
//	处理器函数内部调用 GoAdmin 的 engine.Content 方法
//	engine.Content 方法会调用处理函数，获取面板内容
//	然后将面板渲染为 HTML 响应
//	支持 GET、POST 等 HTTP 方法
//
// 使用示例:
//
//	// 定义处理函数
//	func myHandler(ctx *ghttp.Request) (types.Panel, error) {
//	    // 访问数据库
//	    users, err := db.QueryUsers()
//	    if err != nil {
//	        return nil, err
//	    }
//
//	    // 返回表格面板
//	    return table.NewTable(users), nil
//	}
//
//	// 使用 Content 辅助函数创建 GoFrame 处理器
//	s.BindHandler("GET:/admin", gf.Content(myHandler))
//
// 注意事项:
//
//   - 处理函数不能为空，必须返回有效的面板
//   - 应该在适当的 HTTP 方法中注册处理器
//   - 应该尽早注册处理器，确保处理器能够处理请求
//   - 应该正确处理错误，返回错误信息
//   - 应该避免修改传入的上下文，避免影响后续处理
func Content(handler HandlerFunc) ghttp.HandlerFunc {
	return func(ctx *ghttp.Request) { // 返回 GoFrame 处理器函数
		// 调用 GoAdmin 引擎的内容处理方法
		// 将 GoFrame 请求对象转换为通用接口类型
		engine.Content(ctx, func(ctx interface{}) (types.Panel, error) { // 调用 engine.Content 方法
			// 类型断言：将通用接口转换回 GoFrame 请求对象
			return handler(ctx.(*ghttp.Request)) // 类型断言并调用处理函数
		})
	}
}

// SetApp 实现了 Adapter.SetApp 方法
// 该方法用于设置 GoFrame 服务器实例到适配器中，用于后续的路由注册和请求处理
//
// 参数:
//   - app: Web 应用实例，实际类型应为 *ghttp.Server
//
// 返回值:
//   - error: 如果参数类型错误则返回错误，成功则为 nil
//
// 功能特性:
//
//	设置 GoFrame 服务器实例: 将 GoFrame 服务器实例设置到适配器中，用于后续的路由注册和请求处理
//	支持类型检查: 使用类型断言检查参数类型，确保参数类型正确
//	支持错误处理: 如果参数类型错误，返回错误信息
//
// 说明:
//
//	使用类型断言检查参数类型，确保参数为 *ghttp.Server 类型
//	如果类型断言失败，返回错误信息
//	如果类型断言成功，将 GoFrame 服务器实例存储到适配器的 app 字段中
//	用于后续的路由注册和请求处理
//	该方法通常在 Use 方法中被调用，用于设置应用实例
//
// 注意事项:
//
//   - app 参数必须为 *ghttp.Server 类型，否则会返回错误
//   - 应该在 Use 方法之前调用，确保应用实例在插件初始化前设置
//   - 只能调用一次，多次调用会覆盖之前的应用实例
//   - 该方法支持并发调用，但需要注意应用实例的线程安全性
//   - 如果返回错误，适配器无法正常工作，需要检查参数类型
func (gf *Gf) SetApp(app interface{}) error {
	var (
		eng *ghttp.Server // GoFrame 服务器实例
		ok  bool          // 类型断言结果
	)

	// 类型断言：验证传入的参数是否为 *ghttp.Server 类型
	// GoFrame 使用断言来确保类型安全，这是 Go 语言中处理接口的常见模式
	if eng, ok = app.(*ghttp.Server); !ok { // 类型断言
		return errors.New("gf 适配器 SetApp: 参数类型错误") // 返回错误
	}

	gf.app = eng // 保存 GoFrame 服务器实例
	return nil   // 返回 nil，表示设置成功
}

// AddHandler 实现了 Adapter.AddHandler 方法
// 该方法用于向 GoFrame 应用添加路由处理器，将 GoAdmin 的路由注册到 GoFrame 的路由系统中
// GoFrame 的路由参数格式为 :param，需要转换为标准 URL 查询参数格式 ?param=value
// GoFrame 使用 "METHOD:path" 格式来绑定路由，如 "GET:/admin"
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
//	注册路由: 将 GoAdmin 的路由注册到 GoFrame 的路由系统中
//	支持多种 HTTP 方法: 支持 GET、POST、PUT、DELETE、HEAD、OPTIONS、PATCH 等 HTTP 方法
//	支持路由参数: 支持路由参数（如 :id、:name 等），会自动转换为 URL 查询参数
//	支持处理器链: 支持处理器链（中间件和处理器），实现请求预处理和后处理
//	支持上下文转换: 将 GoFrame 上下文转换为 GoAdmin 上下文，实现框架无关性
//	支持响应头复制: 将 GoAdmin 响应头复制到 GoFrame 响应中，确保响应头正确传递
//	支持响应体写入: 将 GoAdmin 响应体写入到 GoFrame 响应中，确保响应体正确传递
//
// 说明:
//
//	使用 GoFrame 的 BindHandler 方法绑定路由处理器，GoFrame 使用 "METHOD:path" 格式来绑定路由
//	路由处理器接收 GoFrame 请求对象，处理 HTTP 请求和响应
//	首先创建 GoAdmin 上下文，传入 GoFrame 的标准 HTTP 请求对象
//	然后使用正则表达式提取路由参数，GoFrame 的路由参数格式为 :param
//	将路由参数转换为 URL 查询参数，例如：/user/:id → /user?id=123
//	执行处理器链，处理器链包含中间件和处理器
//	将 GoAdmin 响应头复制到 GoFrame 响应中
//	如果响应体不为空，使用 GoFrame 的 WriteStatus 方法写入响应状态码和响应体
//	否则只写入响应状态码
//	该方法通常在 Use 方法中被调用，用于注册插件的路由
//
// 注意事项:
//
//   - method 参数必须为大写，如 "GET"、"POST" 等
//   - path 参数不能为空，必须以 "/" 开头
//   - handlers 参数不能为空，至少包含一个处理器
//   - 路由参数格式为 :param，会自动转换为 ?param=value 格式
//   - 处理器链执行顺序为从左到右，中间件先执行，处理器后执行
//   - 该方法支持并发调用，但需要注意路由注册的线程安全性
//   - 该方法会修改响应头和响应体，确保在调用前未写入响应
func (gf *Gf) AddHandler(method, path string, handlers context.Handlers) {
	// 使用 GoFrame 的 BindHandler 方法绑定路由处理器
	// GoFrame 使用 "METHOD:path" 格式来绑定路由
	gf.app.BindHandler(strings.ToUpper(method)+":"+path, func(c *ghttp.Request) { // 绑定路由处理器

		// 创建 GoAdmin 上下文，传入 GF 的标准 HTTP 请求对象
		ctx := context.NewContext(c.Request) // 创建 GoAdmin 上下文

		// 复制路径用于参数提取
		newPath := path // 复制路径

		// 使用正则表达式匹配路由参数
		// GoFrame 的路由参数格式为 :param，例如 /user/:id
		// 需要提取这些参数并转换为 URL 查询参数
		reg1 := regexp.MustCompile(":(.*?)/") // 匹配路径中间的参数，如 /user/:id/
		reg2 := regexp.MustCompile(":(.*?)$") // 匹配路径末尾的参数，如 /user/:id

		// 查找所有路由参数
		params := reg1.FindAllString(newPath, -1)                   // 查找所有中间参数
		newPath = reg1.ReplaceAllString(newPath, "")                // 移除中间参数
		params = append(params, reg2.FindAllString(newPath, -1)...) // 查找所有末尾参数

		// 将路由参数转换为 URL 查询参数
		// 例如：/user/:id → /user?id=123
		for _, param := range params { // 遍历所有路由参数
			// 移除参数名中的冒号和斜杠
			p := utils.ReplaceAll(param, ":", "", "/", "") // 移除冒号和斜杠

			// 构建查询字符串
			// 如果是第一个参数，直接添加；后续参数用 & 连接
			// 使用 GF 的 GetRequestString 方法获取路由参数值
			if c.Request.URL.RawQuery == "" { // 如果查询字符串为空
				c.Request.URL.RawQuery += p + "=" + c.GetRequestString(p) // 添加第一个参数
			} else { // 如果查询字符串不为空
				c.Request.URL.RawQuery += "&" + p + "=" + c.GetRequestString(p) // 添加后续参数
			}
		}

		// 设置处理器链并执行
		// GoAdmin 使用处理器链模式，每个处理器可以处理请求并传递给下一个
		ctx.SetHandlers(handlers).Next() // 设置处理器链并执行

		// 将 GoAdmin 响应头复制到 GF 响应中
		for key, head := range ctx.Response.Header { // 遍历 GoAdmin 响应头
			c.Response.Header().Add(key, head[0]) // 将响应头设置到 GoFrame 响应中
		}

		// 如果响应体不为空，则写入响应
		if ctx.Response.Body != nil { // 如果 GoAdmin 响应体不为空
			buf := new(bytes.Buffer)               // 创建字节缓冲区
			_, _ = buf.ReadFrom(ctx.Response.Body) // 从 GoAdmin 响应体读取数据到缓冲区
			// 使用 GF 的 WriteStatus 方法写入响应状态码和响应体
			c.Response.WriteStatus(ctx.Response.StatusCode, buf.Bytes()) // 写入响应状态码和响应体
		} else { // 如果 GoAdmin 响应体为空
			// 如果响应体为空，只写入状态码
			c.Response.WriteStatus(ctx.Response.StatusCode) // 只写入响应状态码
		}
	})
}

// Name 实现了 Adapter.Name 方法
// 返回适配器的名称，用于标识不同的框架适配器
//
// 返回值:
//   - string: 适配器名称，固定为 "gf"
//
// 功能特性:
//
//	返回适配器名称: 返回适配器的名称，用于标识不同的框架适配器
//	用于标识适配器类型: 用于标识适配器类型，便于日志记录和调试
//	用于日志记录和调试: 用于日志记录和调试，便于追踪适配器的使用情况
//
// 说明:
//
//	返回适配器的名称，固定为 "gf"
//	用于标识适配器类型，便于日志记录和调试
//	可以在日志中记录适配器名称，便于追踪适配器的使用情况
//	可以在配置文件中指定适配器名称，便于配置管理
//
// 注意事项:
//
//   - 适配器名称必须与适配器类型一致，不能随意修改
//   - 适配器名称不能为空，必须返回有效的字符串
//   - 适配器名称应该唯一，不能与其他适配器名称冲突
func (*Gf) Name() string {
	return "gf" // 返回适配器名称
}

// SetContext 实现了 Adapter.SetContext 方法
// 该方法用于设置当前请求的上下文到适配器中，用于后续的请求处理
//
// 参数:
//   - contextInterface: 请求上下文接口，实际类型应为 *ghttp.Request
//
// 返回值:
//   - adapter.WebFrameWork: 返回新的适配器实例，包含设置的上下文
//
// 功能特性:
//
//	设置当前请求的上下文: 设置当前请求的上下文到适配器中，用于后续的请求处理
//	支持类型检查: 使用类型断言检查参数类型，确保参数类型正确
//	支持错误处理: 如果参数类型错误，使用 panic 终止程序
//	返回新的适配器实例: 返回新的适配器实例，包含设置的上下文
//
// 说明:
//
//	使用类型断言检查参数类型，确保参数为 *ghttp.Request 类型
//	如果类型断言失败，使用 panic 终止程序，这是 Go 语言中处理严重错误的常见方式
//	如果类型断言成功，创建新的 Gf 适配器实例，设置 GoFrame 请求上下文到 ctx 字段中
//	返回新的适配器实例，用于后续的请求处理
//	该方法通常在请求处理中被调用，用于设置请求上下文
//
// 注意事项:
//
//   - contextInterface 参数必须为 *ghttp.Request 类型，否则会 panic
//   - 返回新的适配器实例，不是修改当前适配器实例
//   - 应该在每个请求处理时调用，确保每个请求都有独立的上下文
//   - 该方法支持并发调用，但需要注意上下文的线程安全性
//   - 如果 panic，请求处理会中断，需要确保参数类型正确
func (*Gf) SetContext(contextInterface interface{}) adapter.WebFrameWork {
	var (
		ctx *ghttp.Request // GoFrame 请求上下文
		ok  bool           // 类型断言结果
	)

	// 类型断言：验证传入的参数是否为 *ghttp.Request 类型
	// 如果类型不匹配，使用 panic 终止程序，这是 Go 语言中处理严重错误的常见方式
	if ctx, ok = contextInterface.(*ghttp.Request); !ok { // 类型断言
		panic("gf 适配器 SetContext: 参数类型错误") // 终止程序
	}

	// 返回新的 Gf 适配器实例，包含设置的上下文
	return &Gf{ctx: ctx} // 返回新的适配器实例
}

// Redirect 实现了 Adapter.Redirect 方法
// 该方法用于重定向到登录页面
// 当用户未登录或会话过期时，GoAdmin 会调用此方法将用户重定向到登录页面
//
// 返回值:
//   - 无返回值
//
// 功能特性:
//
//	重定向到登录页面: 将用户重定向到登录页面，用于用户未登录或会话过期的情况
//	使用 HTTP 302 状态码: 使用 HTTP 302 状态码进行临时重定向
//	支持自定义登录页面 URL: 支持自定义登录页面 URL，通过 config.GetLoginUrl() 获取
//
// 说明:
//
//	调用 GoFrame 请求的 Response.RedirectTo 方法，使用 HTTP 302 状态码
//	登录页面 URL 通过 config.GetLoginUrl() 获取
//	重定向 URL 通过 config.Url() 生成，确保 URL 格式正确
//	该方法通常在中间件中被调用，用于重定向未登录的用户
//
// 注意事项:
//
//   - 会修改响应头和响应体，确保在调用前未写入响应
//   - 会终止请求处理，调用后不应该再写入响应
//   - 应该在用户未认证时调用，确保用户被重定向到登录页面
//   - 该方法支持并发调用，但需要注意重定向的线程安全性
func (gf *Gf) Redirect() {
	gf.ctx.Response.RedirectTo(config.Url(config.GetLoginUrl())) // 重定向到登录页面
}

// SetContentType 实现了 Adapter.SetContentType 方法
// 该方法用于设置响应的内容类型
// GoAdmin 默认使用 HTML 内容类型，确保浏览器正确渲染页面
//
// 返回值:
//   - 无返回值
//
// 功能特性:
//
//	设置响应的内容类型: 设置响应的 Content-Type 头，确保浏览器正确渲染页面
//	支持自定义 Content-Type: 支持自定义 Content-Type，通过 HTMLContentType() 方法获取
//
// 说明:
//
//	调用 GoFrame 响应头的 Add 方法，设置 Content-Type 头
//	Content-Type 由 HTMLContentType() 方法确定，通常为 "text/html; charset=utf-8"
//	该方法通常在渲染 HTML 响应前调用，确保浏览器正确渲染页面
//
// 注意事项:
//
//   - 会修改响应头，确保在调用前未写入响应
//   - 应该在渲染 HTML 响应前调用，确保浏览器正确渲染页面
//   - 该方法支持并发调用，但需要注意响应头的线程安全性
func (gf *Gf) SetContentType() {
	gf.ctx.Response.Header().Add("Content-Type", gf.HTMLContentType()) // 设置 Content-Type 头
}

// Write 实现了 Adapter.Write 方法
// 该方法用于写入响应体
//
// 参数:
//   - body: 要写入的响应体字节数组
//
// 返回值:
//   - 无返回值
//
// 功能特性:
//
//	将响应体写入到响应中: 将响应体写入到响应中，用于返回 HTML、JSON 等数据
//	支持写入任意字节数组: 支持写入任意字节数组，包括 HTML、JSON、XML 等
//
// 说明:
//
//	调用 GoFrame 响应的 WriteStatus 方法，结束请求并返回响应
//	响应体可以是 HTML、JSON 等数据
//	响应状态码为 http.StatusOK（200）
//	该方法通常在渲染完成后调用，用于返回响应
//
// 注意事项:
//
//   - 会修改响应体，确保在调用前未写入响应
//   - 应该在渲染完成后调用，确保响应体完整
//   - 该方法支持并发调用，但需要注意响应体的线程安全性
//   - 忽略写入错误，如果写入失败不会返回错误
func (gf *Gf) Write(body []byte) {
	gf.ctx.Response.WriteStatus(http.StatusOK, body) // 写入响应状态码和响应体
}

// GetCookie 实现了 Adapter.GetCookie 方法
// 该方法用于获取指定名称的 Cookie 值
// Cookie 用于存储用户会话信息，如登录凭证
//
// 返回值:
//   - string: Cookie 的值
//   - error: 获取 Cookie 时的错误，总是返回 nil
//
// 功能特性:
//
//	从请求中获取认证 Cookie: 从请求中获取认证 Cookie，用于用户身份验证
//	支持自定义 Cookie 键名: 支持自定义 Cookie 键名，通过 CookieKey() 方法获取
//	总是返回 nil 错误: 总是返回 nil 错误，简化错误处理
//
// 说明:
//
//	调用 GoFrame 请求的 Cookie.Get 方法，获取指定名称的 Cookie 值
//	Cookie 键名通过 CookieKey() 方法获取，默认为 "goadmin_session"
//	返回 Cookie 值和 nil 错误
//	该方法通常在中间件中被调用，用于获取用户会话信息
//
// 注意事项:
//
//   - 返回的 Cookie 值可能为空，如果 Cookie 不存在
//   - 返回的错误总是为 nil，简化错误处理
//   - 该方法支持并发调用，但需要注意 Cookie 的线程安全性
//   - 应该在 User 方法之前调用，确保用户会话信息可用
func (gf *Gf) GetCookie() (string, error) {
	return gf.ctx.Cookie.Get(gf.CookieKey()), nil // 获取 Cookie 值
}

// Lang 实现了 Adapter.Lang 方法
// 该方法用于从查询参数中获取语言设置
// GoAdmin 支持多语言，语言通过 URL 查询参数 __ga_lang 指定
//
// 返回值:
//   - string: 语言代码，如 "zh-CN"、"en-US" 等
//
// 功能特性:
//
//	从 URL 查询参数中获取语言设置: 从 URL 查询参数中获取语言设置，用于多语言支持
//	支持自定义语言参数名: 支持自定义语言参数名，默认为 "__ga_lang"
//	支持多语言切换: 支持多语言切换，通过修改 URL 查询参数实现
//
// 说明:
//
//	调用 GoFrame 请求的 URL.Query().Get 方法，获取指定名称的查询参数值
//	语言参数名为 "__ga_lang"，格式为 "zh-CN"、"en-US" 等
//	返回语言代码，如果查询参数不存在则返回空字符串
//	该方法通常在渲染模板前调用，用于设置模板语言
//
// 注意事项:
//
//   - 返回的语言代码可能为空，如果查询参数不存在
//   - 返回的语言代码应该有效，否则可能导致模板渲染错误
//   - 该方法支持并发调用，但需要注意查询参数的线程安全性
//   - 应该在渲染模板前调用，确保模板语言正确
func (gf *Gf) Lang() string {
	return gf.ctx.Request.URL.Query().Get("__ga_lang") // 获取语言代码
}

// Path 实现了 Adapter.Path 方法
// 该方法用于获取当前请求的路径
//
// 返回值:
//   - string: 请求路径，如 "/admin/info"
//
// 功能特性:
//
//	获取当前请求的路径: 获取当前请求的路径，用于路由匹配和日志记录
//	不包含查询参数: 路径不包含查询参数，只包含路径部分
//	支持路由匹配: 支持路由匹配，用于确定当前请求的路由
//
// 说明:
//
//	调用 GoFrame 请求的 URL.Path 方法，获取请求路径
//	路径不包含查询参数，只包含路径部分
//	路径以 "/" 开头，如 "/admin/info"
//	该方法通常在请求处理中被调用，用于路由匹配和日志记录
//
// 注意事项:
//
//   - 返回的路径不包含查询参数，如果需要查询参数请使用 Query 方法
//   - 返回的路径以 "/" 开头，格式为 "/path/to/resource"
//   - 该方法支持并发调用，但需要注意路径的线程安全性
//   - 应该在请求处理中调用，确保路径正确
func (gf *Gf) Path() string {
	return gf.ctx.URL.Path // 获取请求路径
}

// Method 实现了 Adapter.Method 方法
// 该方法用于获取当前请求的 HTTP 方法
//
// 返回值:
//   - string: HTTP 方法，如 "GET"、"POST"、"PUT"、"DELETE" 等
//
// 功能特性:
//
//	获取当前请求的 HTTP 方法: 获取当前请求的 HTTP 方法，用于请求类型判断和日志记录
//	支持所有 HTTP 方法: 支持所有 HTTP 方法，包括 GET、POST、PUT、DELETE、PATCH、HEAD、OPTIONS 等
//	支持请求类型判断: 支持请求类型判断，用于区分不同的请求类型
//
// 说明:
//
//	调用 GoFrame 请求的 Method 方法，获取 HTTP 方法
//	HTTP 方法包括 GET、POST、PUT、DELETE、PATCH、HEAD、OPTIONS 等
//	返回的 HTTP 方法为大写字母，如 "GET"、"POST" 等
//	该方法通常在请求处理中被调用，用于请求类型判断和日志记录
//
// 注意事项:
//
//   - 返回的 HTTP 方法为大写字母，格式为 "GET"、"POST" 等
//   - 返回的 HTTP 方法应该是有效的 HTTP 方法，否则可能导致请求处理错误
//   - 该方法支持并发调用，但需要注意 HTTP 方法的线程安全性
//   - 应该在请求处理中调用，确保 HTTP 方法正确
func (gf *Gf) Method() string {
	return gf.ctx.Method // 获取 HTTP 方法
}

// FormParam 实现了 Adapter.FormParam 方法
// 该方法用于获取表单参数
// 表单参数通常来自 POST 请求的请求体
// GoFrame 的 Form 字段已经包含了解析后的表单参数
//
// 返回值:
//   - url.Values: 表单参数的键值对集合
//
// 功能特性:
//
//	获取表单参数: 获取表单参数，用于表单数据处理
//	支持多种表单类型: 支持多种表单类型，包括 application/x-www-form-urlencoded 和 multipart/form-data
//	支持文件上传: 支持文件上传，通过 multipart/form-data 格式
//	支持多值字段: 支持多值字段，一个键可以对应多个值
//	自动解析: GoFrame 的 Form 字段已经包含了解析后的表单参数，无需手动解析
//
// 说明:
//
//	返回 GoFrame 请求的 Form 字段，该字段已经包含了解析后的表单参数
//	表单参数包括 POST 请求的请求体中的数据
//	支持多值字段，一个键可以对应多个值
//	该方法通常在请求处理中被调用，用于获取表单数据
//
// 注意事项:
//
//   - 返回的表单参数可能为空，如果请求体不包含表单数据
//   - 返回的表单参数支持多值字段，一个键可以对应多个值
//   - 该方法支持并发调用，但需要注意表单参数的线程安全性
//   - GoFrame 的 Form 字段已经自动解析，无需手动调用 ParseForm 或 ParseMultipartForm
//   - 应该在数据查询时调用，确保表单参数正确
func (gf *Gf) FormParam() url.Values {
	return gf.ctx.Form // 返回表单参数
}

// IsPjax 实现了 Adapter.IsPjax 方法
// 该方法用于判断当前请求是否为 PJAX 请求
// PJAX (PushState + AJAX) 是一种技术，允许在不刷新整个页面的情况下更新页面内容
// GoAdmin 使用 PJAX 来提供更流畅的用户体验
//
// 返回值:
//   - bool: 如果是 PJAX 请求返回 true，否则返回 false
//
// 功能特性:
//
//	检查当前请求是否为 PJAX 请求: 检查当前请求是否为 PJAX 请求，用于页面部分更新
//	支持页面部分更新: 支持页面部分更新，提供更流畅的用户体验
//	支持无刷新导航: 支持无刷新导航，避免页面完全刷新
//
// 说明:
//
//	调用 GoFrame 请求的 Header.Get 方法，获取指定名称的请求头值
//	PJAX 请求头为 "X-PJAX"，如果值为 "true" 则返回 true
//	返回的布尔值表示是否为 PJAX 请求
//	该方法通常在渲染模板前调用，用于判断是否需要渲染完整页面
//
// 注意事项:
//
//   - 返回的布尔值表示是否为 PJAX 请求
//   - 如果是 PJAX 请求，应该只渲染部分页面内容
//   - 该方法支持并发调用，但需要注意请求头的线程安全性
//   - 应该在渲染模板前调用，确保页面渲染正确
//   - 应该与前端 PJAX 库配合使用，确保 PJAX 功能正常
func (gf *Gf) IsPjax() bool {
	return gf.ctx.Header.Get(constant.PjaxHeader) == "true" // 检查是否为 PJAX 请求
}

// Query 实现了 Adapter.Query 方法
// 该方法用于获取 URL 查询参数
// 查询参数是 URL 中 ? 后面的键值对
//
// 返回值:
//   - url.Values: 查询参数的键值对集合
//
// 功能特性:
//
//	获取 URL 查询参数: 获取 URL 查询参数，用于查询参数处理和分页等
//	支持多值字段: 支持多值字段，一个键可以对应多个值
//	支持分页、排序、筛选等功能: 支持分页、排序、筛选等功能，通过查询参数传递
//
// 说明:
//
//	调用 GoFrame 请求的 URL.Query 方法，获取查询参数
//	查询参数是 URL 中 ? 后面的键值对，如 ?page=1&size=10
//	返回查询参数的键值对集合，可以通过键获取对应的值
//	支持多值字段，一个键可以对应多个值
//	该方法通常在请求处理中被调用，用于获取查询参数
//
// 注意事项:
//
//   - 返回的查询参数可能为空，如果 URL 不包含查询参数
//   - 返回的查询参数支持多值字段，一个键可以对应多个值
//   - 该方法支持并发调用，但需要注意查询参数的线程安全性
//   - 应该在数据查询时调用，确保查询参数正确
func (gf *Gf) Query() url.Values {
	return gf.ctx.Request.URL.Query() // 获取查询参数
}

// Request 实现了 Adapter.Request 方法
// 该方法用于获取原始的 HTTP 请求对象
//
// 返回值:
//   - *http.Request: 标准 HTTP 请求对象指针
//
// 功能特性:
//
//	获取原始的 HTTP 请求对象: 获取原始的 HTTP 请求对象，用于底层请求处理
//	支持底层请求处理: 支持底层请求处理，用于自定义请求处理逻辑
//	支持自定义请求处理逻辑: 支持自定义请求处理逻辑，实现特殊的请求处理需求
//
// 说明:
//
//	返回 GoFrame 请求的 Request 字段，该字段已经是标准 HTTP 请求对象
//	HTTP 请求对象包含请求的所有信息，如请求方法、URL、请求头、请求体等
//	该方法通常在底层请求处理中被调用，用于访问原始请求对象
//
// 注意事项:
//
//   - 返回的 HTTP 请求对象可能为空，如果请求上下文未正确初始化
//   - 返回的 HTTP 请求对象是 GoFrame 请求的 Request 字段，不是转换后的对象
//   - 该方法支持并发调用，但需要注意 HTTP 请求对象的线程安全性
//   - 应该在底层请求处理时调用，确保请求对象正确
//   - 应该谨慎使用，避免修改请求对象导致意外的行为
func (gf *Gf) Request() *http.Request {
	return gf.ctx.Request // 返回原始的 HTTP 请求对象
}

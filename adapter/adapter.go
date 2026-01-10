// Copyright 2019 GoAdmin Core Team. All rights reserved.
// Use of this source code is governed by a Apache-2.0 style
// license that can be found in the LICENSE file.
//
// 版权所有 2019 GoAdmin 核心团队。保留所有权利。
// 本源代码的使用受 Apache-2.0 风格许可证约束，
// 该许可证可在 LICENSE 文件中找到。

package adapter

import (
	"bytes"    // 字节缓冲区操作，用于存储渲染后的HTML内容
	"fmt"      // 格式化I/O，用于字符串格式化和错误消息生成
	"net/http" // HTTP客户端和服务器，提供HTTP请求和响应处理
	"net/url"  // URL解析和查询，用于处理URL参数和表单数据

	"github.com/purpose168/GoAdmin/context"              // 上下文管理，提供请求和响应的抽象
	"github.com/purpose168/GoAdmin/modules/auth"         // 认证模块，处理用户身份验证和权限管理
	"github.com/purpose168/GoAdmin/modules/config"       // 配置模块，提供配置读取和URL处理功能
	"github.com/purpose168/GoAdmin/modules/constant"     // 常量定义，包含框架使用的各种常量
	"github.com/purpose168/GoAdmin/modules/db"           // 数据库模块，提供数据库连接和操作接口
	"github.com/purpose168/GoAdmin/modules/errors"       // 错误处理，提供框架特定的错误类型和消息
	"github.com/purpose168/GoAdmin/modules/logger"       // 日志模块，提供日志记录功能
	"github.com/purpose168/GoAdmin/modules/menu"         // 菜单模块，提供菜单生成和管理功能
	"github.com/purpose168/GoAdmin/plugins"              // 插件接口，定义插件的基本结构和功能
	"github.com/purpose168/GoAdmin/plugins/admin/models" // 管理模型，提供用户模型等数据结构
	"github.com/purpose168/GoAdmin/template"             // 模板引擎，提供HTML模板渲染功能
	"github.com/purpose168/GoAdmin/template/types"       // 模板类型，定义面板、按钮等模板组件
)

// WebFrameWork 接口定义了Web框架适配器需要实现的方法
// 该接口是GoAdmin与不同Web框架（如Gin、Echo、Chi等）之间的桥梁
// 通过实现这个接口，GoAdmin可以适配到各种Go Web框架
//
// 设计模式说明：
// - 这是一个适配器模式(Adapter Pattern)的应用
// - 允许GoAdmin核心功能与不同的Web框架解耦
// - 每个Web框架都需要实现这个接口来提供适配器
//
// 使用示例：
//
//	// Gin框架适配器示例
//	type GinAdapter struct {
//	    BaseAdapter
//	    engine *gin.Engine
//	    ctx    *gin.Context
//	}
//
//	func (g *GinAdapter) Name() string {
//	    return "gin"
//	}
//
//	// 实现其他接口方法...
type WebFrameWork interface {
	// Name 返回适配器的名称
	// 用于标识当前使用的是哪个Web框架
	// 返回值示例: "gin", "echo", "chi", "fiber" 等
	Name() string

	// Use 将GoAdmin插件注册到Web框架中
	// 该方法会遍历所有插件，并将它们的路由处理器添加到Web框架中
	//
	// 参数说明：
	//   - app: Web框架的实例（如 *gin.Engine, *echo.Echo 等）
	//   - plugins: 插件列表，每个插件提供一组路由和处理函数
	//
	// 返回值：
	//   - error: 如果设置应用或添加处理器失败，返回错误信息
	//
	// 注意事项：
	//   - 插件的路由会根据插件的Prefix()方法添加前缀
	//   - 每个插件可能包含多个路由和处理函数
	//   - 该方法通常在初始化时调用一次
	Use(app interface{}, plugins []plugins.Plugin) error

	// Content 渲染并返回管理面板的内容
	// 这是GoAdmin的核心方法，负责生成整个管理页面的HTML
	//
	// 参数说明：
	//   - ctx: Web框架的上下文对象（如 *gin.Context）
	//   - fn: 获取面板内容的函数，返回 types.Panel
	//   - fn2: 节点处理器，用于处理面板的回调函数
	//   - navButtons: 导航按钮列表，显示在页面顶部
	//
	// 工作流程：
	//   1. 从Cookie中获取用户信息
	//   2. 验证用户身份和权限
	//   3. 调用getPanelFn获取面板内容
	//   4. 渲染完整的HTML页面（包括菜单、面板、按钮等）
	//   5. 将HTML写入响应
	//
	// 注意事项：
	//   - 如果用户未登录或权限不足，会重定向到登录页
	//   - 支持PJAX（部分页面加载）请求
	//   - 会根据IsPjax()决定渲染完整页面还是仅面板内容
	Content(ctx interface{}, fn types.GetPanelFn, fn2 context.NodeProcessor, navButtons ...types.Button)

	// User 从上下文中获取当前登录的用户信息
	//
	// 参数说明：
	//   - ctx: Web框架的上下文对象
	//
	// 返回值：
	//   - models.UserModel: 用户模型，包含用户详细信息
	//   - bool: 是否成功获取用户，true表示成功，false表示失败
	//
	// 注意事项：
	//   - 该方法通过Cookie验证用户身份
	//   - 如果Cookie无效或用户不存在，返回空模型和false
	//   - 获取用户后会释放数据库连接
	User(ctx interface{}) (models.UserModel, bool)

	// AddHandler 添加路由处理器到Web框架
	//
	// 参数说明：
	//   - method: HTTP方法（GET, POST, PUT, DELETE 等）
	//   - path: 路由路径（如 "/admin/user"）
	//   - handlers: 处理器链，按顺序执行
	//
	// 使用场景：
	//   - 注册GoAdmin的管理路由
	//   - 注册插件的API路由
	//   - 注册自定义的路由处理器
	//
	// 注意事项：
	//   - handlers是处理器链，会按顺序执行
	//   - 中间件应该在handlers链的前面
	//   - 路径可以包含参数（如 "/user/:id"）
	AddHandler(method, path string, handlers context.Handlers)

	// DisableLog 禁用日志记录
	// 用于在不需要日志记录的场景下关闭日志输出
	//
	// 注意事项：
	//   - 该方法通常在开发或测试环境中使用
	//   - 不同适配器的实现可能不同
	//   - 某些适配器可能不支持此功能
	DisableLog()

	// Static 配置静态文件服务
	// 用于提供CSS、JS、图片等静态资源
	//
	// 参数说明：
	//   - prefix: URL前缀（如 "/static"）
	//   - path: 文件系统路径（如 "./assets/static"）
	//
	// 使用示例：
	//   Static("/static", "./assets/static")
	//   // 访问 http://localhost:8080/static/style.css
	//   // 会返回 ./assets/static/style.css 文件
	//
	// 注意事项：
	//   - path必须是绝对路径或相对于工作目录的路径
	//   - 如果目录不存在，可能会返回404错误
	Static(prefix, path string)

	// Run 启动Web服务器
	// 开始监听HTTP请求
	//
	// 返回值：
	//   - error: 如果启动失败，返回错误信息
	//
	// 注意事项：
	//   - 该方法会阻塞当前goroutine
	//   - 通常在main函数的最后调用
	//   - 端口配置通常来自配置文件
	Run() error

	// SetApp 设置Web框架的应用实例
	// 用于将适配器与Web框架实例关联
	//
	// 参数说明：
	//   - app: Web框架的实例（如 *gin.Engine, *echo.Echo 等）
	//
	// 返回值：
	//   - error: 如果设置失败，返回错误信息
	//
	// 注意事项：
	//   - 该方法通常在Use()方法中被调用
	//   - 只需要调用一次
	SetApp(app interface{}) error

	// SetConnection 设置数据库连接
	// 用于将数据库连接注入到适配器中
	//
	// 参数说明：
	//   - conn: 数据库连接对象
	//
	// 注意事项：
	//   - 该连接会被用于所有数据库操作
	//   - 连接通常在初始化时设置
	//   - BaseAdapter已经提供了默认实现
	SetConnection(conn db.Connection)

	// GetConnection 获取数据库连接
	// 返回当前适配器使用的数据库连接
	//
	// 返回值：
	//   - db.Connection: 数据库连接对象
	//
	// 注意事项：
	//   - BaseAdapter已经提供了默认实现
	//   - 返回的连接应该被正确管理（释放或复用）
	GetConnection() db.Connection

	// SetContext 设置当前请求的上下文
	// 用于在适配器中保存请求上下文
	//
	// 参数说明：
	//   - ctx: Web框架的上下文对象
	//
	// 返回值：
	//   - WebFrameWork: 返回适配器自身，支持链式调用
	//
	// 使用示例：
	//   adapter.SetContext(ctx).GetCookie()
	//
	// 注意事项：
	//   - 每个请求都应该调用此方法设置上下文
	//   - 返回适配器自身，支持链式调用
	SetContext(ctx interface{}) WebFrameWork

	// GetCookie 从请求中获取Cookie值
	// 用于获取认证Cookie或其他Cookie
	//
	// 返回值：
	//   - string: Cookie的值
	//   - error: 如果获取失败，返回错误信息
	//
	// 注意事项：
	//   - 默认获取名为"admin_cookie"的Cookie
	//   - 如果Cookie不存在，返回空字符串和错误
	GetCookie() (string, error)

	// Lang 获取当前请求的语言设置
	// 用于国际化支持
	//
	// 返回值：
	//   - string: 语言代码（如 "zh-CN", "en-US"）
	//
	// 注意事项：
	//   - 语言通常从请求头或Cookie中获取
	//   - 默认语言通常在配置文件中设置
	Lang() string

	// Path 获取当前请求的路径
	// 不包含查询参数
	//
	// 返回值：
	//   - string: 请求路径（如 "/admin/user/list"）
	//
	// 使用场景：
	//   - 权限验证
	//   - 菜单激活状态判断
	//   - 路由匹配
	Path() string

	// Method 获取当前请求的HTTP方法
	//
	// 返回值：
	//   - string: HTTP方法（GET, POST, PUT, DELETE 等）
	//
	// 使用场景：
	//   - 权限验证（区分读取和修改操作）
	//   - 路由匹配
	Method() string

	// Request 获取原始的HTTP请求对象
	// 返回标准库的http.Request
	//
	// 返回值：
	//   - *http.Request: HTTP请求对象
	//
	// 使用场景：
	//   - 访问请求头、请求体等原始数据
	//   - 与不依赖特定框架的代码交互
	Request() *http.Request

	// FormParam 获取表单参数
	// 解析并返回POST/PUT请求的表单数据
	//
	// 返回值：
	//   - url.Values: 表单参数的键值对
	//
	// 使用示例：
	//   params := adapter.FormParam()
	//   name := params.Get("name")
	//
	// 注意事项：
	//   - 包含Content-Type为application/x-www-form-urlencoded的数据
	//   - 对于multipart/form-data，需要特殊处理
	FormParam() url.Values

	// Query 获取URL查询参数
	// 解析并返回URL中的查询字符串
	//
	// 返回值：
	//   - url.Values: 查询参数的键值对
	//
	// 使用示例：
	//   query := adapter.Query()
	//   page := query.Get("page")
	//
	// 注意事项：
	//   - 返回?后面的参数
	//   - 参数值会自动进行URL解码
	Query() url.Values

	// IsPjax 判断当前请求是否为PJAX请求
	// PJAX是一种技术，允许只更新页面的一部分
	//
	// 返回值：
	//   - bool: true表示是PJAX请求，false表示普通请求
	//
	// 注意事项：
	//   - PJAX请求通常包含X-PJAX请求头
	//   - PJAX请求只返回面板内容，不返回完整页面
	IsPjax() bool

	// Redirect 重定向到登录页
	// 当用户未登录或权限不足时调用
	//
	// 注意事项：
	//   - 通常重定向到配置文件中指定的登录页
	//   - 会保存原始请求URL，登录后可以跳转回来
	Redirect()

	// SetContentType 设置响应的Content-Type
	// 通常设置为"text/html; charset=utf-8"
	//
	// 注意事项：
	//   - 在写入响应体之前调用
	//   - BaseAdapter已经提供了默认实现
	SetContentType()

	// Write 写入响应体
	// 将数据写入HTTP响应
	//
	// 参数说明：
	//   - body: 要写入的字节数据
	//
	// 注意事项：
	//   - 只能调用一次，多次调用会覆盖之前的内容
	//   - 应该在SetContentType()之后调用
	Write(body []byte)

	// CookieKey 获取认证Cookie的键名
	// 返回用于存储用户认证信息的Cookie名称
	//
	// 返回值：
	//   - string: Cookie键名，默认为"admin_cookie"
	//
	// 注意事项：
	//   - BaseAdapter已经提供了默认实现
	//   - 可以在配置文件中自定义
	CookieKey() string

	// HTMLContentType 返回HTML响应的Content-Type
	//
	// 返回值：
	//   - string: Content-Type值，默认为"text/html; charset=utf-8"
	//
	// 注意事项：
	//   - BaseAdapter已经提供了默认实现
	//   - 指定UTF-8编码以支持中文
	HTMLContentType() string
}

// BaseAdapter 是Web框架适配器的基础实现
// 提供了WebFrameWork接口的通用方法，可以被各个Web框架适配器嵌入使用
//
// 设计模式说明：
// - 这是一个模板方法模式(Template Method Pattern)的应用
// - BaseAdapter提供了通用的方法实现
// - 具体的Web框架适配器只需要实现框架特定的方法
// - 减少了重复代码，提高了代码复用性
//
// 使用示例：
//
//	type GinAdapter struct {
//	    BaseAdapter
//	    engine *gin.Engine
//	    ctx    *gin.Context
//	}
//
//	// GinAdapter只需要实现Gin特定的方法
//	// 其他方法可以直接使用BaseAdapter的实现
type BaseAdapter struct {
	// db 数据库连接
	// 用于存储和管理数据库连接，供所有数据库操作使用
	// 该连接在初始化时通过SetConnection方法设置
	db db.Connection
}

// SetConnection 设置数据库连接
// 将数据库连接注入到适配器中，供后续使用
//
// 参数说明：
//   - conn: 数据库连接对象，实现了db.Connection接口
//
// 使用场景：
//   - 在初始化GoAdmin时设置数据库连接
//   - 在需要切换数据库连接时调用
//
// 注意事项：
//   - 该连接会被所有数据库操作共享
//   - 应该确保连接的有效性
//   - 通常在应用启动时调用一次
func (base *BaseAdapter) SetConnection(conn db.Connection) {
	base.db = conn
}

// GetConnection 获取数据库连接
// 返回当前适配器使用的数据库连接
//
// 返回值：
//   - db.Connection: 数据库连接对象
//
// 使用场景：
//   - 在需要执行数据库查询时获取连接
//   - 在验证用户身份时访问数据库
//   - 在获取菜单数据时查询数据库
//
// 注意事项：
//   - 返回的连接应该被正确管理
//   - 不要关闭返回的连接，它可能被其他地方使用
//   - 如果连接未设置，返回nil
func (base *BaseAdapter) GetConnection() db.Connection {
	return base.db
}

// HTMLContentType 返回HTML响应的Content-Type
// 设置为"text/html; charset=utf-8"以支持中文显示
//
// 返回值：
//   - string: Content-Type值，固定为"text/html; charset=utf-8"
//
// 使用场景：
//   - 在设置响应头时使用
//   - 确保浏览器正确解析HTML和中文内容
//
// 注意事项：
//   - UTF-8编码支持所有Unicode字符，包括中文
//   - 该方法被Content()方法调用
//   - 子类可以覆盖此方法以返回不同的Content-Type
func (*BaseAdapter) HTMLContentType() string {
	return "text/html; charset=utf-8"
}

// CookieKey 获取认证Cookie的键名
// 返回用于存储用户认证信息的Cookie名称
//
// 返回值：
//   - string: Cookie键名，默认为"admin_cookie"
//
// 使用场景：
//   - 在设置Cookie时使用
//   - 在获取Cookie时使用
//   - 在验证用户身份时使用
//
// 注意事项：
//   - 默认值来自auth.DefaultCookieKey常量
//   - 可以在配置文件中自定义
//   - 确保前后端使用相同的Cookie键名
func (*BaseAdapter) CookieKey() string {
	return auth.DefaultCookieKey
}

// GetUser 从上下文中获取当前登录的用户信息
// 通过Cookie验证用户身份并返回用户模型
//
// 参数说明：
//   - ctx: Web框架的上下文对象（如 *gin.Context）
//   - wf: WebFrameWork接口实例，用于调用其他方法
//
// 返回值：
//   - models.UserModel: 用户模型，包含用户详细信息
//   - bool: 是否成功获取用户，true表示成功，false表示失败
//
// 工作流程：
//  1. 设置上下文到适配器
//  2. 从请求中获取Cookie
//  3. 使用Cookie验证用户身份
//  4. 释放数据库连接
//  5. 返回用户模型和验证结果
//
// 注意事项：
//   - 如果Cookie获取失败或无效，返回空模型和false
//   - 获取用户后会释放数据库连接，避免连接泄漏
//   - 该方法被User()方法调用
func (*BaseAdapter) GetUser(ctx interface{}, wf WebFrameWork) (models.UserModel, bool) {
	// 从请求中获取Cookie
	cookie, err := wf.SetContext(ctx).GetCookie()

	// 如果获取Cookie失败，返回空模型和false
	if err != nil {
		return models.UserModel{}, false
	}

	// 使用Cookie验证用户身份并获取用户信息
	// auth.GetCurUser会查询数据库验证Cookie的有效性
	user, exist := auth.GetCurUser(cookie, wf.GetConnection())

	// 释放数据库连接，避免连接泄漏
	return user.ReleaseConn(), exist
}

// GetUse 将GoAdmin插件注册到Web框架中
// 该方法会遍历所有插件，并将它们的路由处理器添加到Web框架中
//
// 参数说明：
//   - app: Web框架的实例（如 *gin.Engine, *echo.Echo 等）
//   - plugin: 插件列表，每个插件提供一组路由和处理函数
//   - wf: WebFrameWork接口实例，用于调用其他方法
//
// 返回值：
//   - error: 如果设置应用或添加处理器失败，返回错误信息
//
// 工作流程：
//  1. 调用SetApp设置Web框架实例
//  2. 遍历所有插件
//  3. 对每个插件，遍历其所有路由
//  4. 如果插件有前缀，将前缀添加到路由路径
//  5. 调用AddHandler注册路由处理器
//
// 使用示例：
//
//	err := adapter.GetUse(ginEngine, []plugins.Plugin{adminPlugin}, adapter)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// 注意事项：
//   - 插件的路由会根据插件的Prefix()方法添加前缀
//   - 每个插件可能包含多个路由和处理函数
//   - 该方法通常在初始化时调用一次
//   - 如果插件的前缀为空，则不添加前缀
func (*BaseAdapter) GetUse(app interface{}, plugin []plugins.Plugin, wf WebFrameWork) error {
	// 设置Web框架实例到适配器
	if err := wf.SetApp(app); err != nil {
		return err
	}

	// 遍历所有插件
	for _, plug := range plugin {
		// 遍历插件的所有路由
		for path, handlers := range plug.GetHandler() {
			// 如果插件有前缀，将前缀添加到路由路径
			if plug.Prefix() == "" {
				// 没有前缀，直接注册路由
				wf.AddHandler(path.Method, path.URL, handlers)
			} else {
				// 有前缀，将前缀添加到路由路径
				// config.Url()会处理URL格式
				wf.AddHandler(path.Method, config.Url("/"+plug.Prefix()+path.URL), handlers)
			}
		}
	}

	return nil
}

// Run 启动Web服务器
// 这是一个占位方法，具体实现由各个Web框架适配器提供
//
// 返回值：
//   - 无，直接panic
//
// 注意事项：
//   - 该方法必须被子类覆盖
//   - 如果直接调用，会触发panic
//   - 每个Web框架适配器都应该实现自己的Run方法
func (*BaseAdapter) Run() error { panic("not implement") }

// DisableLog 禁用日志记录
// 这是一个占位方法，具体实现由各个Web框架适配器提供
//
// 注意事项：
//   - 该方法必须被子类覆盖
//   - 如果直接调用，会触发panic
//   - 每个Web框架适配器都应该实现自己的DisableLog方法
func (*BaseAdapter) DisableLog() { panic("not implement") }

// Static 配置静态文件服务
// 这是一个占位方法，具体实现由各个Web框架适配器提供
//
// 参数说明：
//   - _: URL前缀（占位符）
//   - _: 文件系统路径（占位符）
//
// 注意事项：
//   - 该方法必须被子类覆盖
//   - 如果直接调用，会触发panic
//   - 每个Web框架适配器都应该实现自己的Static方法
func (*BaseAdapter) Static(_, _ string) { panic("not implement") }

// GetContent 渲染并返回管理面板的内容
// 这是GoAdmin的核心方法，负责生成整个管理页面的HTML
//
// 参数说明：
//   - ctx: Web框架的上下文对象（如 *gin.Context）
//   - getPanelFn: 获取面板内容的函数，返回 types.Panel
//   - wf: WebFrameWork接口实例，用于调用其他方法
//   - navButtons: 导航按钮列表，显示在页面顶部
//   - fn: 节点处理器，用于处理面板的回调函数
//
// 工作流程：
//  1. 设置上下文到适配器
//  2. 从请求中获取Cookie
//  3. 验证用户身份（如果Cookie无效或为空，重定向到登录页）
//  4. 使用Cookie验证用户身份并获取用户信息
//  5. 如果验证失败，重定向到登录页
//  6. 检查用户权限
//  7. 如果权限不足，显示403错误页面
//  8. 如果权限足够，调用getPanelFn获取面板内容
//  9. 如果获取面板失败，显示错误页面
//  10. 执行面板的回调函数
//  11. 获取模板（根据是否为PJAX请求）
//  12. 渲染完整的HTML页面
//  13. 设置Content-Type
//  14. 将HTML写入响应
//
// 使用示例：
//
//	adapter.GetContent(ctx, func(ctx interface{}) (types.Panel, error) {
//	    return types.Panel{
//	        Title: "用户列表",
//	        Content: template.HTML("<table>...</table>"),
//	    }, nil
//	}, adapter, navButtons, func(callbacks ...context.Node) {
//	    // 处理回调
//	})
//
// 注意事项：
//   - 如果用户未登录或权限不足，会重定向到登录页
//   - 支持PJAX（部分页面加载）请求
//   - 会根据IsPjax()决定渲染完整页面还是仅面板内容
//   - 模板渲染失败会记录错误日志
//   - 面板内容会根据环境（生产/开发）进行不同的处理
func (base *BaseAdapter) GetContent(ctx interface{}, getPanelFn types.GetPanelFn, wf WebFrameWork,
	navButtons types.Buttons, fn context.NodeProcessor) {

	// 设置上下文到适配器，并获取Cookie
	var (
		newBase          = wf.SetContext(ctx)
		cookie, hasError = newBase.GetCookie()
	)

	// 如果获取Cookie失败或Cookie为空，重定向到登录页
	if hasError != nil || cookie == "" {
		newBase.Redirect()
		return
	}

	// 使用Cookie验证用户身份并获取用户信息
	user, authSuccess := auth.GetCurUser(cookie, wf.GetConnection())

	// 如果验证失败，重定向到登录页
	if !authSuccess {
		newBase.Redirect()
		return
	}

	// 准备获取面板内容
	var (
		panel types.Panel
		err   error
	)

	// 创建GoAdmin上下文
	gctx := context.NewContext(newBase.Request())

	// 检查用户权限
	// auth.CheckPermissions会检查用户是否有权访问当前路径和方法
	if !auth.CheckPermissions(user, newBase.Path(), newBase.Method(), newBase.FormParam()) {
		// 权限不足，显示403错误页面
		panel = template.WarningPanel(gctx, errors.NoPermission, template.NoPermission403Page)
	} else {
		// 权限足够，调用getPanelFn获取面板内容
		panel, err = getPanelFn(ctx)
		// 如果获取面板失败，显示错误页面
		if err != nil {
			panel = template.WarningPanel(gctx, err.Error())
		}
	}

	// 执行面板的回调函数
	// 这些回调可能包含一些需要在渲染前执行的逻辑
	fn(panel.Callbacks...)

	// 获取模板（根据是否为PJAX请求）
	// PJAX请求只返回面板内容，不返回完整页面
	tmpl, tmplName := template.Default(gctx).GetTemplate(newBase.IsPjax())

	// 创建缓冲区用于存储渲染后的HTML
	buf := new(bytes.Buffer)

	// 渲染模板
	// types.NewPage创建页面参数，包含用户、菜单、面板等信息
	hasError = tmpl.ExecuteTemplate(buf, tmplName, types.NewPage(gctx, &types.NewPageParam{
		User:         user,                                                                                                                // 当前用户
		Menu:         menu.GetGlobalMenu(user, wf.GetConnection(), newBase.Lang()).SetActiveClass(config.URLRemovePrefix(newBase.Path())), // 全局菜单，并设置当前激活的菜单项
		Panel:        panel.GetContent(config.IsProductionEnvironment()),                                                                  // 面板内容，根据环境处理
		Assets:       template.GetComponentAssetImportHTML(gctx),                                                                          // 资源文件（CSS、JS等）
		Buttons:      navButtons.CheckPermission(user),                                                                                    // 导航按钮，根据用户权限过滤
		TmplHeadHTML: template.Default(gctx).GetHeadHTML(),                                                                                // 模板头部HTML
		TmplFootJS:   template.Default(gctx).GetFootJS(),                                                                                  // 模板底部JS
		Iframe:       newBase.Query().Get(constant.IframeKey) == "true",                                                                   // 是否在iframe中显示
	}))

	// 如果模板渲染失败，记录错误日志
	if hasError != nil {
		logger.Error(fmt.Sprintf("error: %s adapter content, ", newBase.Name()), hasError)
	}

	// 设置响应的Content-Type
	newBase.SetContentType()

	// 将渲染后的HTML写入响应
	newBase.Write(buf.Bytes())
}

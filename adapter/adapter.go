// Copyright 2019 GoAdmin Core Team. All rights reserved.
// Use of this source code is governed by a Apache-2.0 style
// license that can be found in the LICENSE file.

// adapter包提供了Web框架与GoAdmin框架之间的适配器接口
// 该包定义了适配器需要实现的标准接口，使GoAdmin能够与不同的Web框架集成
//
// 主要功能：
//   - 定义WebFrameWork接口，规范适配器的行为
//   - 提供BaseAdapter基础适配器，包含通用功能
//   - 实现用户认证、权限检查、内容渲染等核心逻辑
//   - 支持多种Web框架（如gin、chi、echo、net/http等）
//
// 架构设计：
//   - WebFrameWork接口：定义适配器必须实现的方法
//   - BaseAdapter结构体：提供基础功能，各具体适配器可嵌入使用
//   - 适配器模式：将不同Web框架的API统一为GoAdmin可用的接口
//
// 核心概念：
//   - 适配器（Adapter）：连接Web框架和GoAdmin的桥梁
//   - 上下文（Context）：封装HTTP请求和响应
//   - 插件（Plugin）：可插拔的功能模块
//   - 面板（Panel）：页面内容的数据结构
//
// 使用场景：
//   - 为新的Web框架创建适配器
//   - 理解GoAdmin与Web框架的交互方式
//   - 扩展或修改适配器功能
//
// 注意事项：
//   - 所有适配器必须实现WebFrameWork接口
//   - 适配器应保持框架无关性
//   - 正确处理错误和异常情况
//
// 作者: GoAdmin Core Team
// 创建日期: 2019-01-01
// 版本: 1.0.0
package adapter

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/modules/auth"
	"github.com/purpose168/GoAdmin/modules/config"
	"github.com/purpose168/GoAdmin/modules/constant"
	"github.com/purpose168/GoAdmin/modules/db"
	"github.com/purpose168/GoAdmin/modules/errors"
	"github.com/purpose168/GoAdmin/modules/logger"
	"github.com/purpose168/GoAdmin/modules/menu"
	"github.com/purpose168/GoAdmin/plugins"
	"github.com/purpose168/GoAdmin/plugins/admin/models"
	"github.com/purpose168/GoAdmin/template"
	"github.com/purpose168/GoAdmin/template/types"
)

// WebFrameWork接口是Web框架与GoAdmin之间的适配器接口
// 所有Web框架适配器都必须实现此接口，以使GoAdmin能够与不同的Web框架无缝集成
//
// 接口设计原理：
//   - 使用接口抽象不同Web框架的API差异
//   - 提供统一的方法签名，使GoAdmin代码框架无关
//   - 支持路由注册、上下文处理、用户认证等核心功能
//
// 实现要求：
//  1. 必须实现所有定义的方法
//  2. 保持方法行为的一致性
//  3. 正确处理错误和异常
//  4. 遵循Go语言的接口设计最佳实践
//
// 已实现的适配器：
//   - gin: adapter/gin/gin.go
//   - chi: adapter/chi/chi.go
//   - echo: adapter/echo/echo.go
//   - net/http: adapter/nethttp/nethttp.go
//   - buffalo: adapter/buffalo/buffalo.go
//   - beego: adapter/beego/beego.go
//   - iris: adapter/iris/iris.go
//   - fiber: adapter/fiber/fiber.go
type WebFrameWork interface {
	// Name返回Web框架的名称
	//
	// 返回值：
	//   - string: 框架名称，如"gin"、"chi"、"echo"等
	//
	// 使用场景：
	//   - 识别当前使用的适配器类型
	//   - 日志记录和调试
	//   - 条件判断和分支处理
	Name() string

	// Use方法将插件注入到Web框架引擎中
	//
	// 参数说明：
	//   - app: Web框架引擎实例（如*gin.Engine、*chi.Mux等）
	//   - plugins: 插件列表，包含需要注册的所有插件
	//
	// 返回值：
	//   - error: 注册过程中的错误，成功返回nil
	//
	// 工作原理：
	//   1. 设置Web框架引擎到适配器
	//   2. 遍历所有插件
	//   3. 为每个插件注册路由和处理器
	//   4. 处理插件前缀（如有）
	//
	// 注意事项：
	//   - app参数类型必须与适配器匹配
	//   - 插件的路由会被自动注册到框架
	//   - 前缀会自动添加到插件路由前
	Use(app interface{}, plugins []plugins.Plugin) error

	// Content方法将面板HTML响应添加到Web框架上下文中
	//
	// 参数说明：
	//   - ctx: Web框架上下文，包含HTTP请求和响应
	//   - fn: 获取面板的函数，用于生成面板内容
	//   - fn2: 节点处理器，用于处理面板中的节点
	//   - navButtons: 导航按钮列表，用于在面板上显示操作按钮
	//
	// 工作流程：
	//   1. 从上下文中获取用户信息
	//   2. 验证用户权限
	//   3. 调用getPanelFn生成面板内容
	//   4. 渲染HTML模板
	//   5. 写入响应
	//
	// 使用场景：
	//   - 渲染管理后台页面
	//   - 显示数据表格和表单
	//   - 处理用户交互操作
	Content(ctx interface{}, fn types.GetPanelFn, fn2 context.NodeProcessor, navButtons ...types.Button)

	// User从给定的Web框架上下文中获取认证用户模型
	//
	// 参数说明：
	//   - ctx: Web框架上下文，包含HTTP请求和响应
	//
	// 返回值：
	//   - models.UserModel: 用户模型，包含用户详细信息
	//   - bool: 是否成功获取用户信息，true表示成功，false表示失败
	//
	// 工作原理：
	//   1. 从上下文中提取cookie
	//   2. 使用cookie验证用户身份
	//   3. 从数据库获取用户信息
	//   4. 返回用户模型
	//
	// 使用示例：
	//   user, ok := adapter.User(ctx)
	//   if !ok {
	//       // 处理未登录情况
	//   }
	User(ctx interface{}) (models.UserModel, bool)

	// AddHandler将GoAdmin的路由和处理器注入到Web框架中
	//
	// 参数说明：
	//   - method: HTTP方法，如GET、POST、PUT、DELETE等
	//   - path: 路由路径，支持路径参数（如:/user/:id）
	//   - handlers: 处理器链，按顺序执行
	//
	// 工作原理：
	//   1. 将路由注册到Web框架
	//   2. 绑定处理器链到路由
	//   3. 处理路径参数（如需要）
	//
	// 使用示例：
	//   AddHandler("GET", "/users", handlers)
	//   AddHandler("POST", "/users", handlers)
	AddHandler(method, path string, handlers context.Handlers)

	// DisableLog禁用日志记录
	//
	// 使用场景：
	//   - 在测试环境中禁用日志
	//   - 减少日志输出
	//   - 性能优化
	DisableLog()

	// Static注册静态文件路由
	//
	// 参数说明：
	//   - prefix: URL前缀，如"/static"
	//   - path: 文件系统路径，如"./assets"
	//
	// 使用示例：
	//   Static("/static", "./assets")
	//   访问: http://localhost:8080/static/logo.png
	Static(prefix, path string)

	// Run启动Web服务器
	//
	// 返回值：
	//   - error: 启动过程中的错误
	//
	// 注意事项：
	//   - 此方法会阻塞，直到服务器停止
	//   - 需要先配置好所有路由和中间件
	Run() error

	// 辅助函数
	// ================================

	// SetApp设置Web框架引擎实例到适配器
	//
	// 参数说明：
	//   - app: Web框架引擎实例
	//
	// 返回值：
	//   - error: 设置过程中的错误
	//
	// 注意事项：
	//   - 必须在注册路由之前调用
	SetApp(app interface{}) error

	// SetConnection设置数据库连接
	//
	// 参数说明：
	//   - conn: 数据库连接对象
	//
	// 使用场景：
	//   - 配置数据库连接
	//   - 初始化适配器
	SetConnection(conn db.Connection)

	// GetConnection获取数据库连接
	//
	// 返回值：
	//   - db.Connection: 数据库连接对象
	//
	// 使用场景：
	//   - 执行数据库查询
	//   - 用户认证
	//   - 权限验证
	GetConnection() db.Connection

	// SetContext设置上下文对象到适配器
	//
	// 参数说明：
	//   - ctx: Web框架上下文
	//
	// 返回值：
	//   - WebFrameWork: 返回适配器实例
	//
	// 使用场景：
	//   - 在请求处理过程中设置上下文
	//   - 链式调用适配器方法
	SetContext(ctx interface{}) WebFrameWork

	// GetCookie从请求中获取指定名称的cookie值
	//
	// 返回值：
	//   - string: cookie的值
	//   - error: 获取失败时的错误信息
	//
	// 使用场景：
	//   - 获取会话token
	//   - 读取用户偏好设置
	//   - 验证用户身份
	GetCookie() (string, error)

	// Lang从查询参数中获取语言设置
	//
	// 返回值：
	//   - string: 语言代码，如"zh-CN"、"en-US"等
	//
	// 使用场景：
	//   - 国际化支持
	//   - 根据用户语言显示不同内容
	Lang() string

	// Path获取请求的路径部分
	//
	// 返回值：
	//   - string: 请求路径，如"/admin/users"
	//
	// 使用场景：
	//   - 路由匹配
	//   - 权限验证
	//   - 日志记录
	Path() string

	// Method获取HTTP请求方法
	//
	// 返回值：
	//   - string: HTTP方法，如GET、POST、PUT、DELETE等
	//
	// 使用场景：
	//   - 区分不同类型的请求
	//   - 权限控制
	//   - 日志记录
	Method() string

	// Request获取原始HTTP请求对象
	//
	// 返回值：
	//   - *http.Request: 原始HTTP请求对象
	//
	// 使用场景：
	//   - 访问请求的所有属性
	//   - 获取请求头、Cookie等
	Request() *http.Request

	// FormParam获取表单参数
	//
	// 返回值：
	//   - url.Values: 表单参数的键值对集合
	//
	// 使用场景：
	//   - 处理表单提交
	//   - 文件上传
	//   - POST请求参数获取
	FormParam() url.Values

	// Query获取URL查询参数
	//
	// 返回值：
	//   - url.Values: 查询参数的键值对集合
	//
	// 使用场景：
	//   - 获取GET请求参数
	//   - 分页参数
	//   - 搜索和过滤参数
	Query() url.Values

	// IsPjax判断是否为PJAX请求
	//
	// 返回值：
	//   - bool: true表示是PJAX请求，false表示不是
	//
	// PJAX说明：
	//   PJAX（pushState + AJAX）允许在不刷新整个页面的情况下更新部分内容
	IsPjax() bool

	// Redirect重定向到登录页面
	//
	// 使用场景：
	//   - 用户未登录时重定向
	//   - 会话过期时重新认证
	Redirect()

	// SetContentType设置响应的Content-Type头部
	//
	// 注意事项：
	//   - 必须在写入响应体之前调用
	SetContentType()

	// Write写入响应体数据
	//
	// 参数说明：
	//   - body: 要写入的字节数据
	Write(body []byte)

	// CookieKey返回cookie的键名
	//
	// 返回值：
	//   - string: cookie键名
	CookieKey() string

	// HTMLContentType返回HTML的Content-Type值
	//
	// 返回值：
	//   - string: Content-Type值，通常为"text/html; charset=utf-8"
	HTMLContentType() string
}

// BaseAdapter是基础适配器结构体，包含一些通用的辅助函数
// 各具体适配器可以嵌入此结构体以复用这些功能
//
// 设计模式：
//   - 使用结构体组合实现代码复用
//   - 提供默认实现，具体适配器可选择性覆盖
//   - 遵循Go语言的组合优于继承原则
//
// 字段说明：
//   - db: 数据库连接对象，用于数据库操作
//
// 使用示例：
//
//	type GinAdapter struct {
//	    adapter.BaseAdapter
//	    engine *gin.Engine
//	}
type BaseAdapter struct {
	db db.Connection
}

// SetConnection设置数据库连接
//
// 参数说明：
//   - conn: 数据库连接对象
//
// 使用场景：
//   - 初始化适配器时配置数据库连接
//   - 切换数据库连接
//   - 多数据库支持
func (base *BaseAdapter) SetConnection(conn db.Connection) {
	base.db = conn
}

// GetConnection获取数据库连接
//
// 返回值：
//   - db.Connection: 数据库连接对象
//
// 使用场景：
//   - 执行数据库查询
//   - 用户认证
//   - 权限验证
func (base *BaseAdapter) GetConnection() db.Connection {
	return base.db
}

// HTMLContentType返回默认的HTML Content-Type头部值
//
// 返回值：
//   - string: Content-Type值，固定返回"text/html; charset=utf-8"
//
// 使用场景：
//   - 设置响应的Content-Type头部
//   - 确保正确的字符编码
func (*BaseAdapter) HTMLContentType() string {
	return "text/html; charset=utf-8"
}

// CookieKey返回cookie的键名
//
// 返回值：
//   - string: cookie键名，从auth包获取默认值
//
// 使用场景：
//   - 获取认证cookie的键名
//   - 设置和读取用户会话
func (*BaseAdapter) CookieKey() string {
	return auth.DefaultCookieKey
}

// GetUser是从上下文中获取认证用户模型的辅助函数
//
// 参数说明：
//   - ctx: Web框架上下文，包含HTTP请求和响应
//   - wf: WebFrameWork适配器实例
//
// 返回值：
//   - models.UserModel: 用户模型，包含用户详细信息
//   - bool: 是否成功获取用户信息，true表示成功，false表示失败
//
// 工作原理：
//  1. 使用SetContext设置上下文
//  2. 从上下文中获取cookie
//  3. 使用cookie验证用户身份
//  4. 从数据库获取用户信息
//  5. 释放数据库连接并返回用户模型
//
// 使用示例：
//
//	user, ok := adapter.GetUser(ctx, wf)
//	if !ok {
//	    // 处理未登录情况
//	}
func (*BaseAdapter) GetUser(ctx interface{}, wf WebFrameWork) (models.UserModel, bool) {
	cookie, err := wf.SetContext(ctx).GetCookie()

	if err != nil {
		return models.UserModel{}, false
	}

	user, exist := auth.GetCurUser(cookie, wf.GetConnection())
	return user.ReleaseConn(), exist
}

// GetUse是将插件添加到Web框架的辅助函数
//
// 参数说明：
//   - app: Web框架引擎实例
//   - plugin: 插件列表，包含需要注册的所有插件
//   - wf: WebFrameWork适配器实例
//
// 返回值：
//   - error: 注册过程中的错误，成功返回nil
//
// 工作原理：
//  1. 调用SetApp设置Web框架引擎
//  2. 遍历所有插件
//  3. 获取每个插件的路由和处理器
//  4. 如果插件有前缀，将前缀添加到路由前
//  5. 调用AddHandler注册路由
//
// 注意事项：
//   - 插件前缀会自动添加到路由前
//   - 所有插件的路由都会被注册
//   - app参数类型必须与适配器匹配
func (*BaseAdapter) GetUse(app interface{}, plugin []plugins.Plugin, wf WebFrameWork) error {
	if err := wf.SetApp(app); err != nil {
		return err
	}

	for _, plug := range plugin {
		for path, handlers := range plug.GetHandler() {
			if plug.Prefix() == "" {
				wf.AddHandler(path.Method, path.URL, handlers)
			} else {
				wf.AddHandler(path.Method, config.Url("/"+plug.Prefix()+path.URL), handlers)
			}
		}
	}

	return nil
}

// Run是BaseAdapter的默认实现，会触发panic
//
// 注意事项：
//   - 具体适配器必须覆盖此方法
//   - 此方法仅作为占位符，不应被调用
func (*BaseAdapter) Run() error { panic("not implement") }

// DisableLog是BaseAdapter的默认实现，会触发panic
//
// 注意事项：
//   - 具体适配器必须覆盖此方法
//   - 此方法仅作为占位符，不应被调用
func (*BaseAdapter) DisableLog() { panic("not implement") }

// Static是BaseAdapter的默认实现，会触发panic
//
// 注意事项：
//   - 具体适配器必须覆盖此方法
//   - 此方法仅作为占位符，不应被调用
func (*BaseAdapter) Static(_, _ string) { panic("not implement") }

// GetContent是adapter.Content的辅助函数，用于生成和渲染管理面板内容
//
// 参数说明：
//   - ctx: Web框架上下文，包含HTTP请求和响应
//   - getPanelFn: 获取面板的函数，用于生成面板内容
//   - wf: WebFrameWork适配器实例
//   - navButtons: 导航按钮列表，用于在面板上显示操作按钮
//   - fn: 节点处理器，用于处理面板中的节点
//
// 工作流程：
//  1. 从上下文中获取cookie并验证用户身份
//  2. 检查用户权限
//  3. 生成面板内容
//  4. 渲染HTML模板
//  5. 写入响应
//
// 错误处理：
//   - Cookie获取失败或为空：重定向到登录页
//   - 用户认证失败：重定向到登录页
//   - 权限不足：显示403错误页面
//   - 面板生成失败：显示错误信息
func (base *BaseAdapter) GetContent(ctx interface{}, getPanelFn types.GetPanelFn, wf WebFrameWork,
	navButtons types.Buttons, fn context.NodeProcessor) {

	// 步骤1: 初始化上下文并获取认证cookie
	// SetContext将Web框架的上下文设置到适配器中，便于后续操作
	// GetCookie从请求中提取认证cookie，用于验证用户身份
	var (
		newBase          = wf.SetContext(ctx)
		cookie, hasError = newBase.GetCookie()
	)

	// 步骤2: 验证cookie是否存在
	// 如果获取cookie失败或cookie为空，说明用户未登录，重定向到登录页
	if hasError != nil || cookie == "" {
		newBase.Redirect()
		return
	}

	// 步骤3: 使用cookie从数据库获取用户信息
	// GetCurUser通过cookie查询用户表，返回用户模型和认证状态
	user, authSuccess := auth.GetCurUser(cookie, wf.GetConnection())

	// 步骤4: 检查用户认证是否成功
	// 如果认证失败，重定向到登录页
	if !authSuccess {
		newBase.Redirect()
		return
	}

	// 步骤5: 准备面板和错误变量
	var (
		panel types.Panel
		err   error
	)

	// 步骤6: 创建GoAdmin上下文
	// NewContext封装原始HTTP请求，提供便捷的请求处理方法
	gctx := context.NewContext(newBase.Request())

	// 步骤7: 检查用户权限
	// CheckPermissions验证用户是否有权访问当前路径和方法
	// 参数包括: 用户模型、请求路径、HTTP方法、表单参数
	if !auth.CheckPermissions(user, newBase.Path(), newBase.Method(), newBase.FormParam()) {
		// 权限不足，显示403错误页面
		panel = template.WarningPanel(gctx, errors.NoPermission, template.NoPermission403Page)
	} else {
		// 权限验证通过，调用getPanelFn生成面板内容
		panel, err = getPanelFn(ctx)
		// 如果生成面板时出错，显示错误信息面板
		if err != nil {
			panel = template.WarningPanel(gctx, err.Error())
		}
	}

	// 步骤8: 执行面板回调函数
	// Callbacks是面板中定义的回调函数列表，用于处理异步操作
	fn(panel.Callbacks...)

	// 步骤9: 获取模板
	// 根据是否为PJAX请求选择不同的模板
	// PJAX请求使用简化模板，普通请求使用完整模板
	tmpl, tmplName := template.Default(gctx).GetTemplate(newBase.IsPjax())

	// 步骤10: 准备模板渲染缓冲区
	// 使用bytes.Buffer避免频繁的内存分配
	buf := new(bytes.Buffer)

	// 步骤11: 执行模板渲染
	// NewPage创建页面数据结构，包含以下关键信息:
	//   - User: 当前登录用户信息
	//   - Menu: 全局菜单，根据用户权限动态生成，并设置当前激活项
	//   - Panel: 面板内容，根据环境决定是否压缩
	//   - Assets: 前端资源（CSS/JS）的导入HTML
	//   - Buttons: 导航按钮，根据用户权限过滤
	//   - TmplHeadHTML: 自定义头部HTML
	//   - TmplFootJS: 自定义底部JavaScript
	//   - Iframe: 判断是否在iframe中显示
	hasError = tmpl.ExecuteTemplate(buf, tmplName, types.NewPage(gctx, &types.NewPageParam{
		User:         user,
		Menu:         menu.GetGlobalMenu(user, wf.GetConnection(), newBase.Lang()).SetActiveClass(config.URLRemovePrefix(newBase.Path())),
		Panel:        panel.GetContent(config.IsProductionEnvironment()),
		Assets:       template.GetComponentAssetImportHTML(gctx),
		Buttons:      navButtons.CheckPermission(user),
		TmplHeadHTML: template.Default(gctx).GetHeadHTML(),
		TmplFootJS:   template.Default(gctx).GetFootJS(),
		Iframe:       newBase.Query().Get(constant.IframeKey) == "true",
	}))

	// 步骤12: 处理模板渲染错误
	// 如果渲染失败，记录错误日志
	if hasError != nil {
		logger.Error(fmt.Sprintf("error: %s adapter content, ", newBase.Name()), hasError)
	}

	// 步骤13: 设置响应Content-Type头部
	// 必须在写入响应体之前调用
	newBase.SetContentType()

	// 步骤14: 将渲染后的HTML写入响应
	newBase.Write(buf.Bytes())
}

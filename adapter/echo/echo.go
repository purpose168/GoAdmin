// Copyright 2019 GoAdmin Core Team. All rights reserved.
// Use of this source code is governed by a Apache-2.0 style
// license that can be found in LICENSE file.

// echo 包提供了 GoAdmin 与 Echo Web 框架的适配器实现
// 该适配器允许 GoAdmin 管理后台在 Echo 应用中运行
// Echo 是一个高性能、极简的 Go Web 框架
// 包名：echo
// 作者：GoAdmin Core Team
// 创建日期：2019
// 目的：为 Echo 框架提供 GoAdmin 管理后台的集成支持

package echo

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/purpose168/GoAdmin/adapter"
	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/engine"
	"github.com/purpose168/GoAdmin/modules/config"
	"github.com/purpose168/GoAdmin/plugins"
	"github.com/purpose168/GoAdmin/plugins/admin/models"
	"github.com/purpose168/GoAdmin/plugins/admin/modules/constant"
	"github.com/purpose168/GoAdmin/template/types"
)

// Echo 结构体实现了 GoAdmin 的适配器接口
// 它作为 Echo 框架和 GoAdmin 管理后台之间的桥梁
// 嵌入了 adapter.BaseAdapter 以获得基础适配器功能
// ctx 字段存储当前的 Echo 上下文
// app 字段存储 Echo 应用实例
type Echo struct {
	adapter.BaseAdapter
	ctx echo.Context
	app *echo.Echo
}

// init 函数在包导入时自动执行
// Go 语言的 init 函数会在 main 函数之前自动调用
// 这里使用 init 函数将 Echo 适配器注册到 GoAdmin 引擎中
// 这样 GoAdmin 就知道如何使用 Echo 框架
func init() {
	engine.Register(new(Echo))
}

// User 实现了 Adapter.User 方法
// 该方法用于从当前上下文中获取用户信息
// 参数：
//   - ctx: 上下文接口，通常为 echo.Context 类型
//
// 返回值：
//   - models.UserModel: 用户模型，包含用户信息
//   - bool: 是否成功获取用户信息，true 表示成功
func (e *Echo) User(ctx interface{}) (models.UserModel, bool) {
	return e.GetUser(ctx, e)
}

// Use 实现了 Adapter.Use 方法
// 该方法用于将插件注册到 Echo 应用中
// 参数：
//   - app: 应用接口，通常为 *echo.Echo 类型
//   - plugs: 插件列表，包含需要注册的所有插件
//
// 返回值：
//   - error: 错误信息，如果注册失败则返回错误
func (e *Echo) Use(app interface{}, plugs []plugins.Plugin) error {
	return e.GetUse(app, plugs, e)
}

// Content 实现了 Adapter.Content 方法
// 该方法用于渲染管理面板内容
// 参数：
//   - ctx: 上下文接口，通常为 echo.Context 类型
//   - getPanelFn: 获取面板的函数，返回 types.Panel 类型的面板
//   - fn: 节点处理器，用于处理上下文中的节点
//   - btns: 导航按钮列表，可变参数
func (e *Echo) Content(ctx interface{}, getPanelFn types.GetPanelFn, fn context.NodeProcessor, btns ...types.Button) {
	e.GetContent(ctx, getPanelFn, e, btns, fn)
}

// HandlerFunc 定义了处理函数的类型
// 该函数接收 Echo 上下文，返回面板和可能的错误
// 参数：
//   - ctx: Echo 上下文
//
// 返回值：
//   - types.Panel: 管理面板
//   - error: 错误信息，如果处理失败则返回错误
type HandlerFunc func(ctx echo.Context) (types.Panel, error)

// Content 是一个辅助函数，用于将 HandlerFunc 转换为 Echo 的处理函数
// 这样可以在 Echo 的路由中使用 GoAdmin 的处理函数
// 参数：
//   - handler: 处理函数，接收 Echo 上下文并返回面板
//
// 返回值：
//   - echo.HandlerFunc: Echo 处理器函数
func Content(handler HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		engine.Content(ctx, func(ctx interface{}) (types.Panel, error) {
			return handler(ctx.(echo.Context))
		})
		return nil
	}
}

// SetApp 实现了 Adapter.SetApp 方法
// 该方法用于设置 Echo 应用实例到适配器中
// 参数：
//   - app: 应用接口，必须为 *echo.Echo 类型
//
// 返回值：
//   - error: 错误信息，如果参数类型不正确则返回错误
func (e *Echo) SetApp(app interface{}) error {
	var (
		eng *echo.Echo
		ok  bool
	)
	// 使用类型断言检查 app 是否为 *echo.Echo 类型
	// ok 为 true 表示断言成功，eng 为转换后的值
	if eng, ok = app.(*echo.Echo); !ok {
		return errors.New("echo 适配器 SetApp: 参数类型错误")
	}
	e.app = eng
	return nil
}

// AddHandler 实现了 Adapter.AddHandler 方法
// 该方法用于向 Echo 应用添加路由处理器
// 参数：
//   - method: HTTP 方法，如 "GET"、"POST" 等
//   - path: 路由路径，如 "/admin" 等
//   - handlers: 处理器链，包含要执行的中间件和处理器
func (e *Echo) AddHandler(method, path string, handlers context.Handlers) {
	e.app.Add(strings.ToUpper(method), path, func(c echo.Context) error {
		// 创建 GoAdmin 上下文
		ctx := context.NewContext(c.Request())

		// 将 Echo 的路由参数转换为 URL 查询参数
		// Echo 的路由参数格式为 :param，需要转换为 ?param=value 格式
		for _, key := range c.ParamNames() {
			if c.Request().URL.RawQuery == "" {
				c.Request().URL.RawQuery += strings.ReplaceAll(key, ":", "") + "=" + c.Param(key)
			} else {
				c.Request().URL.RawQuery += "&" + strings.ReplaceAll(key, ":", "") + "=" + c.Param(key)
			}
		}

		// 执行处理器链
		ctx.SetHandlers(handlers).Next()
		// 将 GoAdmin 响应头复制到 Echo 响应中
		for key, head := range ctx.Response.Header {
			c.Response().Header().Set(key, head[0])
		}
		// 将响应体写入 Echo 响应
		if ctx.Response.Body != nil {
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(ctx.Response.Body)
			_ = c.String(ctx.Response.StatusCode, buf.String())
		} else {
			c.Response().WriteHeader(ctx.Response.StatusCode)
		}
		return nil
	})
}

// Name 实现了 Adapter.Name 方法
// 该方法返回适配器的名称
// 返回值：
//   - string: 适配器名称，固定为 "echo"
func (*Echo) Name() string {
	return "echo"
}

// SetContext 实现了 Adapter.SetContext 方法
// 该方法用于设置当前请求的上下文
// 参数：
//   - contextInterface: 上下文接口，必须为 echo.Context 类型
//
// 返回值：
//   - adapter.WebFrameWork: 返回设置了上下文的新适配器实例
func (*Echo) SetContext(contextInterface interface{}) adapter.WebFrameWork {
	var (
		ctx echo.Context
		ok  bool
	)
	// 使用类型断言检查 contextInterface 是否为 echo.Context 类型
	if ctx, ok = contextInterface.(echo.Context); !ok {
		panic("echo 适配器 SetContext: 参数类型错误")
	}
	return &Echo{ctx: ctx}
}

// Redirect 实现了 Adapter.Redirect 方法
// 该方法用于重定向到登录页面
// 使用 HTTP 302 状态码进行临时重定向
func (e *Echo) Redirect() {
	_ = e.ctx.Redirect(http.StatusFound, config.Url(config.GetLoginUrl()))
}

// SetContentType 实现了 Adapter.SetContentType 方法
// 该方法用于设置响应的 Content-Type 头
// Content-Type 由 HTMLContentType() 方法确定
func (e *Echo) SetContentType() {
	e.ctx.Response().Header().Set("Content-Type", e.HTMLContentType())
}

// Write 实现了 Adapter.Write 方法
// 该方法用于将响应体写入到响应中
// 参数：
//   - body: 要写入的响应体字节数组
func (e *Echo) Write(body []byte) {
	e.ctx.Response().WriteHeader(http.StatusOK)
	_, _ = e.ctx.Response().Write(body)
}

// GetCookie 实现了 Adapter.GetCookie 方法
// 该方法用于从请求中获取认证 Cookie
// 返回值：
//   - string: Cookie 值
//   - error: 错误信息，如果获取失败则返回错误
func (e *Echo) GetCookie() (string, error) {
	cookie, err := e.ctx.Cookie(e.CookieKey())
	if err != nil {
		return "", err
	}
	return cookie.Value, err
}

// Lang 实现了 Adapter.Lang 方法
// 该方法用于从 URL 查询参数中获取语言设置
// 返回值：
//   - string: 语言代码，如 "zh-CN"、"en-US" 等
func (e *Echo) Lang() string {
	return e.ctx.Request().URL.Query().Get("__ga_lang")
}

// Path 实现了 Adapter.Path 方法
// 该方法用于获取当前请求的路径
// 返回值：
//   - string: 请求路径，如 "/admin/dashboard"
func (e *Echo) Path() string {
	return e.ctx.Request().URL.Path
}

// Method 实现了 Adapter.Method 方法
// 该方法用于获取当前请求的 HTTP 方法
// 返回值：
//   - string: HTTP 方法，如 "GET"、"POST"、"PUT"、"DELETE" 等
func (e *Echo) Method() string {
	return e.ctx.Request().Method
}

// FormParam 实现了 Adapter.FormParam 方法
// 该方法用于解析并获取表单参数
// 解析的最大内存限制为 32MB
// 返回值：
//   - url.Values: 表单参数的键值对集合
func (e *Echo) FormParam() url.Values {
	_ = e.ctx.Request().ParseMultipartForm(32 << 20)
	return e.ctx.Request().PostForm
}

// IsPjax 实现了 Adapter.IsPjax 方法
// 该方法用于检查当前请求是否为 PJAX 请求
// PJAX 是一种使用 AJAX 技术实现页面部分更新的技术
// 返回值：
//   - bool: 如果是 PJAX 请求则返回 true，否则返回 false
func (e *Echo) IsPjax() bool {
	return e.ctx.Request().Header.Get(constant.PjaxHeader) == "true"
}

// Query 实现了 Adapter.Query 方法
// 该方法用于获取 URL 查询参数
// 返回值：
//   - url.Values: 查询参数的键值对集合
func (e *Echo) Query() url.Values {
	return e.ctx.Request().URL.Query()
}

// Request 实现了 Adapter.Request 方法
// 该方法用于获取原始的 HTTP 请求对象
// 返回值：
//   - *http.Request: HTTP 请求对象指针
func (e *Echo) Request() *http.Request {
	return e.ctx.Request()
}

// Copyright 2019 GoAdmin Core Team. All rights reserved.
// Use of this source code is governed by a Apache-2.0 style
// license that can be found in LICENSE file.

// chi 包提供了 GoAdmin 与 Chi Web 框架的适配器实现
// 该适配器允许 GoAdmin 管理后台在 Chi 应用中运行
// Chi 是一个轻量级、可组合的 HTTP 路由器
// 包名：chi
// 作者：GoAdmin Core Team
// 创建日期：2019
// 目的：为 Chi 框架提供 GoAdmin 管理后台的集成支持

package chi

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/go-chi/chi"
	"github.com/purpose168/GoAdmin/adapter"
	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/engine"
	cfg "github.com/purpose168/GoAdmin/modules/config"
	"github.com/purpose168/GoAdmin/plugins"
	"github.com/purpose168/GoAdmin/plugins/admin/models"
	"github.com/purpose168/GoAdmin/plugins/admin/modules/constant"
	"github.com/purpose168/GoAdmin/template/types"
)

// Chi 结构体实现了 GoAdmin 的适配器接口
// 它作为 Chi 框架和 GoAdmin 管理后台之间的桥梁
// 嵌入了 adapter.BaseAdapter 以获得基础适配器功能
// ctx 字段存储当前的 Chi 上下文
// app 字段存储 Chi 路由器实例
type Chi struct {
	adapter.BaseAdapter
	ctx Context
	app *chi.Mux
}

// init 函数在包导入时自动执行
// Go 语言的 init 函数会在 main 函数之前自动调用
// 这里使用 init 函数将 Chi 适配器注册到 GoAdmin 引擎中
// 这样 GoAdmin 就知道如何使用 Chi 框架
func init() {
	engine.Register(new(Chi))
}

// User 实现了 Adapter.User 方法
// 该方法用于从当前上下文中获取用户信息
// 参数：
//   - ctx: 上下文接口，通常为 Context 类型
//
// 返回值：
//   - models.UserModel: 用户模型，包含用户信息
//   - bool: 是否成功获取用户信息，true 表示成功
func (ch *Chi) User(ctx interface{}) (models.UserModel, bool) {
	return ch.GetUser(ctx, ch)
}

// Use 实现了 Adapter.Use 方法
// 该方法用于将插件注册到 Chi 应用中
// 参数：
//   - app: 应用接口，通常为 *chi.Mux 类型
//   - plugs: 插件列表，包含需要注册的所有插件
//
// 返回值：
//   - error: 错误信息，如果注册失败则返回错误
func (ch *Chi) Use(app interface{}, plugs []plugins.Plugin) error {
	return ch.GetUse(app, plugs, ch)
}

// Content 实现了 Adapter.Content 方法
// 该方法用于渲染管理面板内容
// 参数：
//   - ctx: 上下文接口，通常为 Context 类型
//   - getPanelFn: 获取面板的函数，返回 types.Panel 类型的面板
//   - fn: 节点处理器，用于处理上下文中的节点
//   - btns: 导航按钮列表，可变参数
func (ch *Chi) Content(ctx interface{}, getPanelFn types.GetPanelFn, fn context.NodeProcessor, btns ...types.Button) {
	ch.GetContent(ctx, getPanelFn, ch, btns, fn)
}

// HandlerFunc 定义了处理函数的类型
// 该函数接收 Chi 上下文，返回面板和可能的错误
// 参数：
//   - ctx: Chi 上下文
//
// 返回值：
//   - types.Panel: 管理面板
//   - error: 错误信息，如果处理失败则返回错误
type HandlerFunc func(ctx Context) (types.Panel, error)

// Content 是一个辅助函数，用于将 HandlerFunc 转换为 HTTP 处理函数
// 这样可以在 Chi 的路由中使用 GoAdmin 的处理函数
// 参数：
//   - handler: 处理函数，接收 Chi 上下文并返回面板
//
// 返回值：
//   - http.HandlerFunc: HTTP 处理器函数
func Content(handler HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := Context{
			Request:  request,
			Response: writer,
		}
		engine.Content(ctx, func(ctx interface{}) (types.Panel, error) {
			return handler(ctx.(Context))
		})
	}
}

// SetApp 实现了 Adapter.SetApp 方法
// 该方法用于设置 Chi 路由器实例到适配器中
// 参数：
//   - app: 应用接口，必须为 *chi.Mux 类型
//
// 返回值：
//   - error: 错误信息，如果参数类型不正确则返回错误
func (ch *Chi) SetApp(app interface{}) error {
	var (
		eng *chi.Mux
		ok  bool
	)
	// 使用类型断言检查 app 是否为 *chi.Mux 类型
	// ok 为 true 表示断言成功，eng 为转换后的值
	if eng, ok = app.(*chi.Mux); !ok {
		return errors.New("chi 适配器 SetApp: 参数类型错误")
	}
	ch.app = eng
	return nil
}

// AddHandler 实现了 Adapter.AddHandler 方法
// 该方法用于向 Chi 路由器添加路由处理器
// 参数：
//   - method: HTTP 方法，如 "GET"、"POST" 等
//   - path: 路由路径，如 "/admin" 等
//   - handlers: 处理器链，包含要执行的中间件和处理器
func (ch *Chi) AddHandler(method, path string, handlers context.Handlers) {
	url := path
	// 使用正则表达式将 Chi 的路由参数格式 :param 转换为 Go 标准格式 {param}
	// reg1 匹配中间的参数，如 :id/，转换为 {id}/
	reg1 := regexp.MustCompile(":(.*?)/")
	// reg2 匹配结尾的参数，如 :id$，转换为 {id}
	reg2 := regexp.MustCompile(":(.*?)$")
	url = reg1.ReplaceAllString(url, "{$1}/")
	url = reg2.ReplaceAllString(url, "{$1}")

	// 如果路径以 // 开头，移除第一个 /
	if len(url) > 1 && url[0] == '/' && url[1] == '/' {
		url = url[1:]
	}

	// 根据方法类型获取对应的路由注册函数
	getHandleFunc(ch.app, strings.ToUpper(method))(url, func(w http.ResponseWriter, r *http.Request) {

		// 如果路径以 / 结尾，移除末尾的 /
		if r.URL.Path[len(r.URL.Path)-1] == '/' {
			r.URL.Path = r.URL.Path[:len(r.URL.Path)-1]
		}

		// 创建 GoAdmin 上下文
		ctx := context.NewContext(r)

		// 获取 Chi 的路由参数并转换为 URL 查询参数
		params := chi.RouteContext(r.Context()).URLParams

		for i := 0; i < len(params.Values); i++ {
			if r.URL.RawQuery == "" {
				r.URL.RawQuery += strings.ReplaceAll(params.Keys[i], ":", "") + "=" + params.Values[i]
			} else {
				r.URL.RawQuery += "&" + strings.ReplaceAll(params.Keys[i], ":", "") + "=" + params.Values[i]
			}
		}

		// 执行处理器链
		ctx.SetHandlers(handlers).Next()
		// 将 GoAdmin 响应头复制到 Chi 响应中
		for key, head := range ctx.Response.Header {
			w.Header().Set(key, head[0])
		}
		// 将响应体写入 Chi 响应
		if ctx.Response.Body != nil {
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(ctx.Response.Body)
			w.WriteHeader(ctx.Response.StatusCode)
			_, _ = w.Write(buf.Bytes())
		} else {
			w.WriteHeader(ctx.Response.StatusCode)
		}
	})
}

// HandleFun 定义了 Chi 路由方法的函数类型
// 该类型用于表示 Chi 的路由注册函数，如 Get、Post 等
// 参数：
//   - pattern: 路由模式
//   - handlerFn: 处理器函数
type HandleFun func(pattern string, handlerFn http.HandlerFunc)

// getHandleFunc 根据方法名称返回对应的路由注册函数
// 这是一个辅助函数，用于将字符串方法名转换为 Chi 的路由注册函数
// 参数：
//   - eng: Chi 路由器实例
//   - method: HTTP 方法名称，如 "GET"、"POST" 等
//
// 返回值：
//   - HandleFun: 对应的路由注册函数
func getHandleFunc(eng *chi.Mux, method string) HandleFun {
	switch method {
	case "GET":
		return eng.Get
	case "POST":
		return eng.Post
	case "PUT":
		return eng.Put
	case "DELETE":
		return eng.Delete
	case "HEAD":
		return eng.Head
	case "OPTIONS":
		return eng.Options
	case "PATCH":
		return eng.Patch
	default:
		panic("错误的 HTTP 方法")
	}
}

// Context 封装了 Chi 的 Request 和 Response 对象
// 该结构体用于在适配器中传递请求和响应信息
type Context struct {
	Request  *http.Request
	Response http.ResponseWriter
}

// SetContext 实现了 Adapter.SetContext 方法
// 该方法用于设置当前请求的上下文
// 参数：
//   - contextInterface: 上下文接口，必须为 Context 类型
//
// 返回值：
//   - adapter.WebFrameWork: 返回设置了上下文的新适配器实例
func (*Chi) SetContext(contextInterface interface{}) adapter.WebFrameWork {
	var (
		ctx Context
		ok  bool
	)
	// 使用类型断言检查 contextInterface 是否为 Context 类型
	if ctx, ok = contextInterface.(Context); !ok {
		panic("chi 适配器 SetContext: 参数类型错误")
	}
	return &Chi{ctx: ctx}
}

// Name 实现了 Adapter.Name 方法
// 该方法返回适配器的名称
// 返回值：
//   - string: 适配器名称，固定为 "chi"
func (*Chi) Name() string {
	return "chi"
}

// Redirect 实现了 Adapter.Redirect 方法
// 该方法用于重定向到登录页面
// 使用 HTTP 302 状态码进行临时重定向
func (ch *Chi) Redirect() {
	http.Redirect(ch.ctx.Response, ch.ctx.Request, cfg.Url(cfg.GetLoginUrl()), http.StatusFound)
}

// SetContentType 实现了 Adapter.SetContentType 方法
// 该方法用于设置响应的 Content-Type 头
// Content-Type 由 HTMLContentType() 方法确定
func (ch *Chi) SetContentType() {
	ch.ctx.Response.Header().Set("Content-Type", ch.HTMLContentType())
}

// Write 实现了 Adapter.Write 方法
// 该方法用于将响应体写入到响应中
// 参数：
//   - body: 要写入的响应体字节数组
func (ch *Chi) Write(body []byte) {
	ch.ctx.Response.WriteHeader(http.StatusOK)
	_, _ = ch.ctx.Response.Write(body)
}

// GetCookie 实现了 Adapter.GetCookie 方法
// 该方法用于从请求中获取认证 Cookie
// 返回值：
//   - string: Cookie 值
//   - error: 错误信息，如果获取失败则返回错误
func (ch *Chi) GetCookie() (string, error) {
	cookie, err := ch.ctx.Request.Cookie(ch.CookieKey())
	if err != nil {
		return "", err
	}
	return cookie.Value, err
}

// Lang 实现了 Adapter.Lang 方法
// 该方法用于从 URL 查询参数中获取语言设置
// 返回值：
//   - string: 语言代码，如 "zh-CN"、"en-US" 等
func (ch *Chi) Lang() string {
	return ch.ctx.Request.URL.Query().Get("__ga_lang")
}

// Path 实现了 Adapter.Path 方法
// 该方法用于获取当前请求的路径
// 返回值：
//   - string: 请求路径，如 "/admin/dashboard"
func (ch *Chi) Path() string {
	return ch.ctx.Request.URL.Path
}

// Method 实现了 Adapter.Method 方法
// 该方法用于获取当前请求的 HTTP 方法
// 返回值：
//   - string: HTTP 方法，如 "GET"、"POST"、"PUT"、"DELETE" 等
func (ch *Chi) Method() string {
	return ch.ctx.Request.Method
}

// FormParam 实现了 Adapter.FormParam 方法
// 该方法用于解析并获取表单参数
// 解析的最大内存限制为 32MB
// 返回值：
//   - url.Values: 表单参数的键值对集合
func (ch *Chi) FormParam() url.Values {
	_ = ch.ctx.Request.ParseMultipartForm(32 << 20)
	return ch.ctx.Request.PostForm
}

// IsPjax 实现了 Adapter.IsPjax 方法
// 该方法用于检查当前请求是否为 PJAX 请求
// PJAX 是一种使用 AJAX 技术实现页面部分更新的技术
// 返回值：
//   - bool: 如果是 PJAX 请求则返回 true，否则返回 false
func (ch *Chi) IsPjax() bool {
	return ch.ctx.Request.Header.Get(constant.PjaxHeader) == "true"
}

// Query 实现了 Adapter.Query 方法
// 该方法用于获取 URL 查询参数
// 返回值：
//   - url.Values: 查询参数的键值对集合
func (ch *Chi) Query() url.Values {
	return ch.ctx.Request.URL.Query()
}

// Request 实现了 Adapter.Request 方法
// 该方法用于获取原始的 HTTP 请求对象
// 返回值：
//   - *http.Request: HTTP 请求对象指针
func (ch *Chi) Request() *http.Request {
	return ch.ctx.Request
}

// Copyright 2019 GoAdmin Core Team. All rights reserved.
// Use of this source code is governed by a Apache-2.0 style
// license that can be found in LICENSE file.

// fasthttp 包提供了 GoAdmin 与 Fasthttp Web 框架的适配器实现
// 该适配器允许 GoAdmin 管理后台在 Fasthttp 应用中运行
// Fasthttp 是一个高性能的 Go HTTP 服务器和客户端库
// 包名：fasthttp
// 作者：GoAdmin Core Team
// 创建日期：2019
// 目的：为 Fasthttp 框架提供 GoAdmin 管理后台的集成支持

package fasthttp

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/buaazp/fasthttprouter"
	"github.com/purpose168/GoAdmin/adapter"
	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/engine"
	"github.com/purpose168/GoAdmin/modules/config"
	"github.com/purpose168/GoAdmin/plugins"
	"github.com/purpose168/GoAdmin/plugins/admin/models"
	"github.com/purpose168/GoAdmin/plugins/admin/modules/constant"
	"github.com/purpose168/GoAdmin/template/types"
	"github.com/valyala/fasthttp"
)

// Fasthttp 结构体实现了 GoAdmin 的适配器接口
// 它作为 Fasthttp 框架和 GoAdmin 管理后台之间的桥梁
// 嵌入了 adapter.BaseAdapter 以获得基础适配器功能
// ctx 字段存储当前的 Fasthttp 请求上下文
// app 字段存储 Fasthttp 路由器实例
type Fasthttp struct {
	adapter.BaseAdapter
	ctx *fasthttp.RequestCtx
	app *fasthttprouter.Router
}

// init 函数在包导入时自动执行
// Go 语言的 init 函数会在 main 函数之前自动调用
// 这里使用 init 函数将 Fasthttp 适配器注册到 GoAdmin 引擎中
// 这样 GoAdmin 就知道如何使用 Fasthttp 框架
func init() {
	engine.Register(new(Fasthttp))
}

// User 实现了 Adapter.User 方法
// 该方法用于从当前上下文中获取用户信息
// 参数：
//   - ctx: 上下文接口，通常为 *fasthttp.RequestCtx 类型
//
// 返回值：
//   - models.UserModel: 用户模型，包含用户信息
//   - bool: 是否成功获取用户信息，true 表示成功
func (fast *Fasthttp) User(ctx interface{}) (models.UserModel, bool) {
	return fast.GetUser(ctx, fast)
}

// Use 实现了 Adapter.Use 方法
// 该方法用于将插件注册到 Fasthttp 应用中
// 参数：
//   - app: 应用接口，通常为 *fasthttprouter.Router 类型
//   - plugs: 插件列表，包含需要注册的所有插件
//
// 返回值：
//   - error: 错误信息，如果注册失败则返回错误
func (fast *Fasthttp) Use(app interface{}, plugs []plugins.Plugin) error {
	return fast.GetUse(app, plugs, fast)
}

// Content 实现了 Adapter.Content 方法
// 该方法用于渲染管理面板内容
// 参数：
//   - ctx: 上下文接口，通常为 *fasthttp.RequestCtx 类型
//   - getPanelFn: 获取面板的函数，返回 types.Panel 类型的面板
//   - fn: 节点处理器，用于处理上下文中的节点
//   - btns: 导航按钮列表，可变参数
func (fast *Fasthttp) Content(ctx interface{}, getPanelFn types.GetPanelFn, fn context.NodeProcessor, btns ...types.Button) {
	fast.GetContent(ctx, getPanelFn, fast, btns, fn)
}

// HandlerFunc 定义了处理函数的类型
// 该函数接收 Fasthttp 请求上下文，返回面板和可能的错误
// 参数：
//   - ctx: Fasthttp 请求上下文指针
//
// 返回值：
//   - types.Panel: 管理面板
//   - error: 错误信息，如果处理失败则返回错误
type HandlerFunc func(ctx *fasthttp.RequestCtx) (types.Panel, error)

// Content 是一个辅助函数，用于将 HandlerFunc 转换为 Fasthttp 的请求处理器
// 这样可以在 Fasthttp 的路由中使用 GoAdmin 的处理函数
// 参数：
//   - handler: 处理函数，接收 Fasthttp 请求上下文并返回面板
//
// 返回值：
//   - fasthttp.RequestHandler: Fasthttp 请求处理器函数
func Content(handler HandlerFunc) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		engine.Content(ctx, func(ctx interface{}) (types.Panel, error) {
			return handler(ctx.(*fasthttp.RequestCtx))
		})
	}
}

// SetApp 实现了 Adapter.SetApp 方法
// 该方法用于设置 Fasthttp 路由器实例到适配器中
// 参数：
//   - app: 应用接口，必须为 *fasthttprouter.Router 类型
//
// 返回值：
//   - error: 错误信息，如果参数类型不正确则返回错误
func (fast *Fasthttp) SetApp(app interface{}) error {
	var (
		eng *fasthttprouter.Router
		ok  bool
	)
	// 使用类型断言检查 app 是否为 *fasthttprouter.Router 类型
	// ok 为 true 表示断言成功，eng 为转换后的值
	if eng, ok = app.(*fasthttprouter.Router); !ok {
		return errors.New("fasthttp 适配器 SetApp: 参数类型错误")
	}

	fast.app = eng
	return nil
}

// AddHandler 实现了 Adapter.AddHandler 方法
// 该方法用于向 Fasthttp 路由器添加路由处理器
// 参数：
//   - method: HTTP 方法，如 "GET"、"POST" 等
//   - path: 路由路径，如 "/admin" 等
//   - handlers: 处理器链，包含要执行的中间件和处理器
func (fast *Fasthttp) AddHandler(method, path string, handlers context.Handlers) {
	fast.app.Handle(strings.ToUpper(method), path, func(c *fasthttp.RequestCtx) {
		// 将 Fasthttp 请求上下文转换为标准 HTTP 请求
		httpreq := convertCtx(c)
		// 创建 GoAdmin 上下文
		ctx := context.NewContext(httpreq)

		// 获取 Fasthttp 的路由参数并转换为 URL 查询参数
		var params = make(map[string]string)
		c.VisitUserValues(func(i []byte, i2 interface{}) {
			if value, ok := i2.(string); ok {
				params[string(i)] = value
			}
		})

		for key, value := range params {
			if httpreq.URL.RawQuery == "" {
				httpreq.URL.RawQuery += strings.ReplaceAll(key, ":", "") + "=" + value
			} else {
				httpreq.URL.RawQuery += "&" + strings.ReplaceAll(key, ":", "") + "=" + value
			}
		}

		// 执行处理器链
		ctx.SetHandlers(handlers).Next()
		// 将 GoAdmin 响应头复制到 Fasthttp 响应中
		for key, head := range ctx.Response.Header {
			c.Response.Header.Set(key, head[0])
		}
		// 将响应体写入 Fasthttp 响应
		if ctx.Response.Body != nil {
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(ctx.Response.Body)
			_, _ = c.WriteString(buf.String())
		}
		c.Response.SetStatusCode(ctx.Response.StatusCode)
	})
}

// convertCtx 将 Fasthttp 请求上下文转换为标准 HTTP 请求
// 这是一个辅助函数，用于在 GoAdmin 和 Fasthttp 之间转换请求对象
// 参数：
//   - ctx: Fasthttp 请求上下文
//
// 返回值：
//   - *http.Request: 标准 HTTP 请求对象指针
func convertCtx(ctx *fasthttp.RequestCtx) *http.Request {
	var r http.Request

	body := ctx.PostBody()
	r.Method = string(ctx.Method())
	r.Proto = "HTTP/1.1"
	r.ProtoMajor = 1
	r.ProtoMinor = 1
	r.RequestURI = string(ctx.RequestURI())
	r.ContentLength = int64(len(body))
	r.Host = string(ctx.Host())
	r.RemoteAddr = ctx.RemoteAddr().String()

	// 复制请求头
	hdr := make(http.Header)
	ctx.Request.Header.VisitAll(func(k, v []byte) {
		sk := string(k)
		sv := string(v)
		switch sk {
		case "Transfer-Encoding":
			r.TransferEncoding = append(r.TransferEncoding, sv)
		default:
			hdr.Set(sk, sv)
		}
	})
	r.Header = hdr
	r.Body = &netHTTPBody{body}

	// 解析请求 URI
	rURL, err := url.ParseRequestURI(r.RequestURI)
	if err != nil {
		ctx.Logger().Printf("fasthttp 适配器 convertCtx: 无法解析请求URI %q: %s", r.RequestURI, err)
		ctx.Error("fasthttp 适配器 convertCtx: 内部服务器错误", fasthttp.StatusInternalServerError)
		return &r
	}
	r.URL = rURL
	return &r
}

// netHTTPBody 实现了 io.Reader 接口，用于包装 Fasthttp 的请求体
// 该结构体允许将 Fasthttp 的字节切片转换为可读取的流
type netHTTPBody struct {
	b []byte
}

// Read 实现了 io.Reader 接口的 Read 方法
// 该方法从包装的字节切片中读取数据
// 参数：
//   - p: 读取缓冲区
//
// 返回值：
//   - int: 实际读取的字节数
//   - error: 错误信息，如果读取失败则返回错误
func (r *netHTTPBody) Read(p []byte) (int, error) {
	if len(r.b) == 0 {
		return 0, io.EOF
	}
	n := copy(p, r.b)
	r.b = r.b[n:]
	return n, nil
}

// Close 实现了 io.Closer 接口的 Close 方法
// 该方法清空包装的字节切片
// 返回值：
//   - error: 总是返回 nil
func (r *netHTTPBody) Close() error {
	r.b = r.b[:0]
	return nil
}

// Name 实现了 Adapter.Name 方法
// 该方法返回适配器的名称
// 返回值：
//   - string: 适配器名称，固定为 "fasthttp"
func (*Fasthttp) Name() string {
	return "fasthttp"
}

// SetContext 实现了 Adapter.SetContext 方法
// 该方法用于设置当前请求的上下文
// 参数：
//   - contextInterface: 上下文接口，必须为 *fasthttp.RequestCtx 类型
//
// 返回值：
//   - adapter.WebFrameWork: 返回设置了上下文的新适配器实例
func (*Fasthttp) SetContext(contextInterface interface{}) adapter.WebFrameWork {
	var (
		ctx *fasthttp.RequestCtx
		ok  bool
	)
	// 使用类型断言检查 contextInterface 是否为 *fasthttp.RequestCtx 类型
	if ctx, ok = contextInterface.(*fasthttp.RequestCtx); !ok {
		panic("fasthttp 适配器 SetContext: 参数类型错误")
	}
	return &Fasthttp{ctx: ctx}
}

// Redirect 实现了 Adapter.Redirect 方法
// 该方法用于重定向到登录页面
// 使用 HTTP 302 状态码进行临时重定向
func (fast *Fasthttp) Redirect() {
	fast.ctx.Redirect(config.Url(config.GetLoginUrl()), http.StatusFound)
}

// SetContentType 实现了 Adapter.SetContentType 方法
// 该方法用于设置响应的 Content-Type 头
// Content-Type 由 HTMLContentType() 方法确定
func (fast *Fasthttp) SetContentType() {
	fast.ctx.Response.Header.Set("Content-Type", fast.HTMLContentType())
}

// Write 实现了 Adapter.Write 方法
// 该方法用于将响应体写入到响应中
// 参数：
//   - body: 要写入的响应体字节数组
func (fast *Fasthttp) Write(body []byte) {
	_, _ = fast.ctx.Write(body)
}

// GetCookie 实现了 Adapter.GetCookie 方法
// 该方法用于从请求中获取认证 Cookie
// 返回值：
//   - string: Cookie 值
//   - error: 错误信息，当前实现总是返回 nil
func (fast *Fasthttp) GetCookie() (string, error) {
	return string(fast.ctx.Request.Header.Cookie(fast.CookieKey())), nil
}

// Lang 实现了 Adapter.Lang 方法
// 该方法用于从 URL 查询参数中获取语言设置
// 返回值：
//   - string: 语言代码，如 "zh-CN"、"en-US" 等
func (fast *Fasthttp) Lang() string {
	return string(fast.ctx.Request.URI().QueryArgs().Peek("__ga_lang"))
}

// Path 实现了 Adapter.Path 方法
// 该方法用于获取当前请求的路径
// 返回值：
//   - string: 请求路径，如 "/admin/dashboard"
func (fast *Fasthttp) Path() string {
	return string(fast.ctx.Path())
}

// Method 实现了 Adapter.Method 方法
// 该方法用于获取当前请求的 HTTP 方法
// 返回值：
//   - string: HTTP 方法，如 "GET"、"POST"、"PUT"、"DELETE" 等
func (fast *Fasthttp) Method() string {
	return string(fast.ctx.Method())
}

// FormParam 实现了 Adapter.FormParam 方法
// 该方法用于解析并获取表单参数
// 返回值：
//   - url.Values: 表单参数的键值对集合
func (fast *Fasthttp) FormParam() url.Values {
	f, _ := fast.ctx.MultipartForm()
	if f != nil {
		return f.Value
	}
	return url.Values{}
}

// IsPjax 实现了 Adapter.IsPjax 方法
// 该方法用于检查当前请求是否为 PJAX 请求
// PJAX 是一种使用 AJAX 技术实现页面部分更新的技术
// 返回值：
//   - bool: 如果是 PJAX 请求则返回 true，否则返回 false
func (fast *Fasthttp) IsPjax() bool {
	return string(fast.ctx.Request.Header.Peek(constant.PjaxHeader)) == "true"
}

// Query 实现了 Adapter.Query 方法
// 该方法用于获取 URL 查询参数
// 返回值：
//   - url.Values: 查询参数的键值对集合
func (fast *Fasthttp) Query() url.Values {
	queryStr := fast.ctx.URI().QueryString()
	queryObj, err := url.Parse(string(queryStr))

	if err != nil {
		return url.Values{}
	}

	return queryObj.Query()
}

// Request 实现了 Adapter.Request 方法
// 该方法用于获取原始的 HTTP 请求对象
// 返回值：
//   - *http.Request: HTTP 请求对象指针
func (fast *Fasthttp) Request() *http.Request {
	return convertCtx(fast.ctx)
}

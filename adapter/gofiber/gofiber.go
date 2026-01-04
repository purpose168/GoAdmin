// 版权所有 2019 GoAdmin 核心团队。保留所有权利。
// 本源代码的使用受 Apache-2.0 风格许可证约束
// 该许可证可在 LICENSE 文件中找到

package gofiber

import (
	"errors"   // 错误处理包
	"io"       // 输入输出接口
	"net/http" // HTTP客户端和服务器实现
	"net/url"  // URL解析和查询
	"strings"  // 字符串操作

	"github.com/gofiber/fiber/v2"
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

// Gofiber 结构体实现了 GoAdmin 的适配器接口
// 它作为 GoFiber 框架和 GoAdmin 管理后台之间的桥梁
// 嵌入了 adapter.BaseAdapter 以获得基础适配器功能
// ctx 字段存储当前的 GoFiber 请求上下文
// app 字段存储 GoFiber 应用实例
type Gofiber struct {
	adapter.BaseAdapter
	ctx *fiber.Ctx
	app *fiber.App
}

// init 函数在包导入时自动执行
// 它将 Gofiber 适配器注册到 GoAdmin 引擎中
// 这样 GoAdmin 就能够识别和使用 GoFiber 框架
func init() {
	engine.Register(new(Gofiber))
}

// User 实现了 Adapter.User 方法
// 该方法从请求上下文中获取当前登录的用户信息
// 参数：
//   - ctx: 请求上下文接口，实际类型为 *fiber.Ctx
//
// 返回值：
//   - models.UserModel: 用户模型，包含用户信息
//   - bool: 是否成功获取用户信息
func (gof *Gofiber) User(ctx interface{}) (models.UserModel, bool) {
	return gof.GetUser(ctx, gof)
}

// Use 实现了 Adapter.Use 方法
// 该方法用于初始化并使用 GoAdmin 插件
// 参数：
//   - app: Web 应用实例，实际类型为 *fiber.App
//   - plugs: 插件列表，包含要加载的 GoAdmin 插件
//
// 返回值：
//   - error: 初始化过程中的错误，成功则为 nil
func (gof *Gofiber) Use(app interface{}, plugs []plugins.Plugin) error {
	return gof.GetUse(app, plugs, gof)
}

// Content 实现了 Adapter.Content 方法
// 该方法用于渲染管理面板内容
// 参数：
//   - ctx: 请求上下文接口
//   - getPanelFn: 获取面板的函数，用于生成管理面板
//   - fn: 节点处理器，用于处理面板中的节点
//   - btns: 按钮列表，用于面板上的操作按钮
func (gof *Gofiber) Content(ctx interface{}, getPanelFn types.GetPanelFn, fn context.NodeProcessor, btns ...types.Button) {
	gof.GetContent(ctx, getPanelFn, gof, btns, fn)
}

// SetApp 实现了 Adapter.SetApp 方法
// 该方法用于设置 GoFiber 应用实例到适配器中
// 参数：
//   - app: Web 应用实例，实际类型应为 *fiber.App
//
// 返回值：
//   - error: 如果参数类型错误则返回错误，成功则为 nil
func (gof *Gofiber) SetApp(app interface{}) error {
	var (
		eng *fiber.App
		ok  bool
	)

	// 类型断言：验证传入的参数是否为 *fiber.App 类型
	// GoFiber 使用断言来确保类型安全，这是 Go 语言中处理接口的常见模式
	if eng, ok = app.(*fiber.App); !ok {
		return errors.New("gofiber 适配器 SetApp: 参数类型错误")
	}

	// 保存 GoFiber 应用实例
	gof.app = eng
	return nil
}

// AddHandler 实现了 Adapter.AddHandler 方法
// 该方法用于向 GoFiber 应用添加路由处理器
// GoFiber 基于 fasthttp，需要将 fasthttp 请求转换为标准 HTTP 请求
// 参数：
//   - method: HTTP 方法，如 "GET"、"POST" 等
//   - path: 路由路径，如 "/admin" 等
//   - handlers: 处理器链，包含要执行的中间件和处理器
func (gof *Gofiber) AddHandler(method, path string, handlers context.Handlers) {
	// 使用 GoFiber 的 Add 方法注册路由处理器
	gof.app.Add(strings.ToUpper(method), path, func(c *fiber.Ctx) error {
		// 将 GoFiber 的 fasthttp 请求上下文转换为标准 HTTP 请求
		httpreq := convertCtx(c.Context())
		// 创建 GoAdmin 上下文
		ctx := context.NewContext(httpreq)

		// 将 GoFiber 的路由参数转换为 URL 查询参数
		// GoFiber 已经将路由参数解析到 c.Route().Params 中
		// 例如：/user/:id → 需要转换为：/user?id=123
		for _, key := range c.Route().Params {
			// 移除参数名中的冒号
			// 如果是第一个参数，直接添加；后续参数用 & 连接
			if httpreq.URL.RawQuery == "" {
				httpreq.URL.RawQuery += strings.ReplaceAll(key, ":", "") + "=" + c.Params(key)
			} else {
				httpreq.URL.RawQuery += "&" + strings.ReplaceAll(key, ":", "") + "=" + c.Params(key)
			}
		}

		// 设置处理器链并执行
		// GoAdmin 使用处理器链模式，每个处理器可以处理请求并传递给下一个
		ctx.SetHandlers(handlers).Next()

		// 将 GoAdmin 响应头复制到 GoFiber 响应中
		for key, head := range ctx.Response.Header {
			c.Set(key, head[0])
		}

		// 使用 GoFiber 的 SendStream 方法发送响应
		// GoFiber 的响应体已经是一个 io.Reader，可以直接发送
		return c.Status(ctx.Response.StatusCode).SendStream(ctx.Response.Body)
	})
}

// convertCtx 将 fasthttp 请求上下文转换为标准 HTTP 请求
// 这是一个辅助函数，用于在 GoAdmin 和 GoFiber 之间转换请求对象
// GoFiber 基于 fasthttp，而 GoAdmin 使用标准 net/http 包
// 参数：
//   - ctx: fasthttp 请求上下文
//
// 返回值：
//   - *http.Request: 标准 HTTP 请求对象指针
func convertCtx(ctx *fasthttp.RequestCtx) *http.Request {
	var r http.Request

	// 获取请求体
	body := ctx.PostBody()

	// 设置 HTTP 方法
	r.Method = string(ctx.Method())

	// 设置 HTTP 协议版本
	r.Proto = "HTTP/1.1"
	r.ProtoMajor = 1
	r.ProtoMinor = 1

	// 设置请求 URI
	r.RequestURI = string(ctx.RequestURI())

	// 设置内容长度
	r.ContentLength = int64(len(body))

	// 设置主机名
	r.Host = string(ctx.Host())

	// 设置远程地址
	r.RemoteAddr = ctx.RemoteAddr().String()

	// 转换请求头
	hdr := make(http.Header)
	ctx.Request.Header.VisitAll(func(k, v []byte) {
		sk := string(k)
		sv := string(v)
		switch sk {
		case "Transfer-Encoding":
			// 特殊处理 Transfer-Encoding 头
			r.TransferEncoding = append(r.TransferEncoding, sv)
		default:
			// 普通请求头直接设置
			hdr.Set(sk, sv)
		}
	})
	r.Header = hdr

	// 设置请求体
	r.Body = &netHTTPBody{body}

	// 解析请求 URI
	rURL, err := url.ParseRequestURI(r.RequestURI)
	if err != nil {
		// 如果解析失败，记录错误并返回内部服务器错误
		ctx.Logger().Printf("无法解析请求 URI %q: %s", r.RequestURI, err)
		ctx.Error("内部服务器错误", fasthttp.StatusInternalServerError)
		return &r
	}
	r.URL = rURL
	return &r
}

// netHTTPBody 实现了 io.ReadCloser 接口
// 它包装了 fasthttp 的请求体字节数组，使其兼容标准 HTTP 请求体
type netHTTPBody struct {
	b []byte
}

// Read 实现了 io.Reader 接口的 Read 方法
// 该方法用于读取请求体数据
// 参数：
//   - p: 用于存储读取数据的字节切片
//
// 返回值：
//   - int: 实际读取的字节数
//   - error: 读取错误，如果读取完成则返回 io.EOF
func (r *netHTTPBody) Read(p []byte) (int, error) {
	if len(r.b) == 0 {
		return 0, io.EOF
	}
	n := copy(p, r.b)
	r.b = r.b[n:]
	return n, nil
}

// Close 实现了 io.Closer 接口的 Close 方法
// 该方法用于关闭请求体，释放资源
// 返回值：
//   - error: 关闭错误，总是返回 nil
func (r *netHTTPBody) Close() error {
	r.b = r.b[:0]
	return nil
}

// Name 实现了 Adapter.Name 方法
// 返回适配器的名称，用于标识不同的框架适配器
// 返回值：
//   - string: 适配器名称 "gofiber"
func (*Gofiber) Name() string {
	return "gofiber"
}

// SetContext 实现了 Adapter.SetContext 方法
// 该方法用于设置当前请求的上下文到适配器中
// 参数：
//   - contextInterface: 请求上下文接口，实际类型应为 *fiber.Ctx
//
// 返回值：
//   - adapter.WebFrameWork: 返回新的适配器实例，包含设置的上下文
func (*Gofiber) SetContext(contextInterface interface{}) adapter.WebFrameWork {
	var (
		ctx *fiber.Ctx
		ok  bool
	)

	// 类型断言：验证传入的参数是否为 *fiber.Ctx 类型
	// 如果类型不匹配，使用 panic 终止程序，这是 Go 语言中处理严重错误的常见方式
	if ctx, ok = contextInterface.(*fiber.Ctx); !ok {
		panic("gofiber 适配器 SetContext: 参数类型错误")
	}

	// 返回新的 Gofiber 适配器实例，包含设置的上下文
	return &Gofiber{ctx: ctx}
}

// Redirect 实现了 Adapter.Redirect 方法
// 该方法用于重定向到登录页面
// 当用户未登录或会话过期时，GoAdmin 会调用此方法将用户重定向到登录页面
// 使用 http.StatusFound (302) 状态码进行临时重定向
func (gof *Gofiber) Redirect() {
	_ = gof.ctx.Redirect(config.Url(config.GetLoginUrl()), http.StatusFound)
}

// SetContentType 实现了 Adapter.SetContentType 方法
// 该方法用于设置响应的内容类型
// GoAdmin 默认使用 HTML 内容类型，确保浏览器正确渲染页面
func (gof *Gofiber) SetContentType() {
	gof.ctx.Response().Header.Set("Content-Type", gof.HTMLContentType())
}

// Write 实现了 Adapter.Write 方法
// 该方法用于写入响应体
// 参数：
//   - body: 要写入的响应体字节数组
func (gof *Gofiber) Write(body []byte) {
	_, _ = gof.ctx.Write(body)
}

// GetCookie 实现了 Adapter.GetCookie 方法
// 该方法用于获取指定名称的 Cookie 值
// Cookie 用于存储用户会话信息，如登录凭证
// 返回值：
//   - string: Cookie 的值
//   - error: 获取 Cookie 时的错误（GoFiber 的 Cookies 方法不返回错误）
func (gof *Gofiber) GetCookie() (string, error) {
	return string(gof.ctx.Cookies(gof.CookieKey())), nil
}

// Lang 实现了 Adapter.Lang 方法
// 该方法用于从查询参数中获取语言设置
// GoAdmin 支持多语言，语言通过 URL 查询参数 __ga_lang 指定
// 返回值：
//   - string: 语言代码，如 "zh-CN"、"en-US" 等
func (gof *Gofiber) Lang() string {
	return string(gof.ctx.Request().URI().QueryArgs().Peek("__ga_lang"))
}

// Path 实现了 Adapter.Path 方法
// 该方法用于获取当前请求的路径
// 返回值：
//   - string: 请求路径，如 "/admin/info"
func (gof *Gofiber) Path() string {
	return string(gof.ctx.Path())
}

// Method 实现了 Adapter.Method 方法
// 该方法用于获取当前请求的 HTTP 方法
// 返回值：
//   - string: HTTP 方法，如 "GET"、"POST"、"PUT"、"DELETE" 等
func (gof *Gofiber) Method() string {
	return string(gof.ctx.Method())
}

// FormParam 实现了 Adapter.FormParam 方法
// 该方法用于获取表单参数
// 表单参数通常来自 POST 请求的请求体
// GoFiber 的 MultipartForm 方法会自动解析 multipart/form-data 格式的表单数据
// 返回值：
//   - url.Values: 表单参数的键值对集合
func (gof *Gofiber) FormParam() url.Values {
	f, _ := gof.ctx.MultipartForm()
	if f != nil {
		return f.Value
	}
	return url.Values{}
}

// IsPjax 实现了 Adapter.IsPjax 方法
// 该方法用于判断当前请求是否为 PJAX 请求
// PJAX (PushState + AJAX) 是一种技术，允许在不刷新整个页面的情况下更新页面内容
// GoAdmin 使用 PJAX 来提供更流畅的用户体验
// 返回值：
//   - bool: 如果是 PJAX 请求返回 true，否则返回 false
func (gof *Gofiber) IsPjax() bool {
	return string(gof.ctx.Request().Header.Peek(constant.PjaxHeader)) == "true"
}

// Query 实现了 Adapter.Query 方法
// 该方法用于获取 URL 查询参数
// 查询参数是 URL 中 ? 后面的键值对
// GoFiber 使用 QueryString() 方法获取查询字符串，然后解析为 url.Values
// 返回值：
//   - url.Values: 查询参数的键值对集合
func (gof *Gofiber) Query() url.Values {
	queryStr := gof.ctx.Context().QueryArgs().QueryString()
	queryObj, err := url.Parse(string(queryStr))

	if err != nil {
		return url.Values{}
	}

	return queryObj.Query()
}

// Request 实现了 Adapter.Request 方法
// 该方法用于获取原始的 HTTP 请求对象
// 需要将 GoFiber 的 fasthttp 请求上下文转换为标准 HTTP 请求
// 返回值：
//   - *http.Request: 标准 HTTP 请求对象指针
func (gof *Gofiber) Request() *http.Request {
	return convertCtx(gof.ctx.Context())
}

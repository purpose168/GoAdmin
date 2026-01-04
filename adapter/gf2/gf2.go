package gf2

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/purpose168/GoAdmin/adapter"
	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/engine"
	"github.com/purpose168/GoAdmin/modules/config"
	"github.com/purpose168/GoAdmin/modules/constant"
	"github.com/purpose168/GoAdmin/modules/utils"
	"github.com/purpose168/GoAdmin/plugins"
	"github.com/purpose168/GoAdmin/plugins/admin/models"
	"github.com/purpose168/GoAdmin/template/types"
)

// GF2 结构体实现了 GoAdmin 的适配器接口
// 它作为 GoFrame v2 (GF2) 框架和 GoAdmin 管理后台之间的桥梁
// 嵌入了 adapter.BaseAdapter 以获得基础适配器功能
// ctx 字段存储当前的 GF2 请求上下文
// app 字段存储 GF2 服务器实例
type GF2 struct {
	adapter.BaseAdapter
	ctx *ghttp.Request
	app *ghttp.Server
}

// init 函数在包导入时自动执行
// 它将 GF2 适配器注册到 GoAdmin 引擎中
// 这样 GoAdmin 就能够识别和使用 GoFrame v2 框架
func init() {
	engine.Register(new(GF2))
}

// Name 实现了 Adapter.Name 方法
// 返回适配器的名称，用于标识不同的框架适配器
// 返回值：
//   - string: 适配器名称 "gf2"
func (*GF2) Name() string {
	return "gf2"
}

// Use 实现了 Adapter.Use 方法
// 该方法用于初始化并使用 GoAdmin 插件
// 参数：
//   - app: Web 应用实例，实际类型为 *ghttp.Server
//   - plugins: 插件列表，包含要加载的 GoAdmin 插件
//
// 返回值：
//   - error: 初始化过程中的错误，成功则为 nil
func (gf2 *GF2) Use(app interface{}, plugins []plugins.Plugin) error {
	return gf2.GetUse(app, plugins, gf2)
}

// Content 实现了 Adapter.Content 方法
// 该方法用于渲染管理面板内容
// 参数：
//   - ctx: 请求上下文接口
//   - getPanelFn: 获取面板的函数，用于生成管理面板
//   - fn: 节点处理器，用于处理面板中的节点
//   - btns: 按钮列表，用于面板上的操作按钮
func (gf2 *GF2) Content(ctx interface{}, getPanelFn types.GetPanelFn, fn context.NodeProcessor, btns ...types.Button) {
	gf2.GetContent(ctx, getPanelFn, gf2, btns, fn)
}

// User 实现了 Adapter.User 方法
// 该方法从请求上下文中获取当前登录的用户信息
// 参数：
//   - ctx: 请求上下文接口，实际类型为 *ghttp.Request
//
// 返回值：
//   - models.UserModel: 用户模型，包含用户信息
//   - bool: 是否成功获取用户信息
func (gf2 *GF2) User(ctx interface{}) (models.UserModel, bool) {
	return gf2.GetUser(ctx, gf2)
}

// AddHandler 实现了 Adapter.AddHandler 方法
// 该方法用于向 GoFrame v2 应用添加路由处理器
// GoFrame v2 的路由参数格式为 :param，需要转换为标准 URL 查询参数格式 ?param=value
// 参数：
//   - method: HTTP 方法，如 "GET"、"POST" 等
//   - path: 路由路径，如 "/admin" 等
//   - handlers: 处理器链，包含要执行的中间件和处理器
func (gf2 *GF2) AddHandler(method, path string, handlers context.Handlers) {
	// 使用 GoFrame v2 的 BindHandler 方法绑定路由处理器
	// GoFrame v2 使用 "METHOD:path" 格式来绑定路由
	gf2.app.BindHandler(strings.ToUpper(method)+":"+path, func(c *ghttp.Request) {
		// 创建 GoAdmin 上下文，传入 GF2 的标准 HTTP 请求对象
		ctx := context.NewContext(c.Request)

		// 复制路径用于参数提取
		newPath := path

		// 使用正则表达式匹配路由参数
		// GoFrame v2 的路由参数格式为 :param，例如 /user/:id
		// 需要提取这些参数并转换为 URL 查询参数
		reg1 := regexp.MustCompile(":(.*?)/") // 匹配路径中间的参数，如 /user/:id/
		reg2 := regexp.MustCompile(":(.*?)$") // 匹配路径末尾的参数，如 /user/:id

		// 查找所有路由参数
		params := reg1.FindAllString(newPath, -1)
		newPath = reg1.ReplaceAllString(newPath, "")
		params = append(params, reg2.FindAllString(newPath, -1)...)

		// 将路由参数转换为 URL 查询参数
		// 例如：/user/:id → /user?id=123
		for _, param := range params {
			// 移除参数名中的冒号和斜杠
			p := utils.ReplaceAll(param, ":", "", "/", "")

			// 构建查询字符串
			// 如果是第一个参数，直接添加；后续参数用 & 连接
			// 使用 GF2 的 GetRequest 方法获取路由参数值，并调用 String() 方法转换为字符串
			// 注意：GoFrame v2 的 GetRequest 返回的是 gvar.Var 类型，需要调用 String() 方法
			if c.Request.URL.RawQuery == "" {
				c.Request.URL.RawQuery += p + "=" + c.GetRequest(p).String()
			} else {
				c.Request.URL.RawQuery += "&" + p + "=" + c.GetRequest(p).String()
			}
		}

		// 设置处理器链并执行
		// GoAdmin 使用处理器链模式，每个处理器可以处理请求并传递给下一个
		ctx.SetHandlers(handlers).Next()

		// 将 GoAdmin 响应头复制到 GF2 响应中
		for key, head := range ctx.Response.Header {
			c.Response.Header().Add(key, head[0])
		}

		// 如果响应体不为空，则写入响应
		if ctx.Response.Body != nil {
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(ctx.Response.Body)
			// 使用 GF2 的 WriteStatus 方法写入响应状态码和响应体
			c.Response.WriteStatus(ctx.Response.StatusCode, buf.Bytes())
		} else {
			// 如果响应体为空，只写入状态码
			c.Response.WriteStatus(ctx.Response.StatusCode)
		}
	})
}

// SetApp 实现了 Adapter.SetApp 方法
// 该方法用于设置 GoFrame v2 服务器实例到适配器中
// 参数：
//   - app: Web 应用实例，实际类型应为 *ghttp.Server
//
// 返回值：
//   - error: 如果参数类型错误则返回错误，成功则为 nil
func (gf2 *GF2) SetApp(app interface{}) error {
	var (
		eng *ghttp.Server
		ok  bool
	)

	// 类型断言：验证传入的参数是否为 *ghttp.Server 类型
	// GoFrame v2 使用断言来确保类型安全，这是 Go 语言中处理接口的常见模式
	if eng, ok = app.(*ghttp.Server); !ok {
		return errors.New("gf2 适配器 SetApp: 参数类型错误")
	}

	// 保存 GoFrame v2 服务器实例
	gf2.app = eng
	return nil
}

// SetContext 实现了 Adapter.SetContext 方法
// 该方法用于设置当前请求的上下文到适配器中
// 参数：
//   - contextInterface: 请求上下文接口，实际类型应为 *ghttp.Request
//
// 返回值：
//   - adapter.WebFrameWork: 返回新的适配器实例，包含设置的上下文
func (*GF2) SetContext(contextInterface interface{}) adapter.WebFrameWork {
	var (
		ctx *ghttp.Request
		ok  bool
	)

	// 类型断言：验证传入的参数是否为 *ghttp.Request 类型
	// 如果类型不匹配，使用 panic 终止程序，这是 Go 语言中处理严重错误的常见方式
	if ctx, ok = contextInterface.(*ghttp.Request); !ok {
		panic("gf2 适配器 SetContext: 参数类型错误")
	}

	// 返回新的 GF2 适配器实例，包含设置的上下文
	return &GF2{ctx: ctx}
}

// GetCookie 实现了 Adapter.GetCookie 方法
// 该方法用于获取指定名称的 Cookie 值
// Cookie 用于存储用户会话信息，如登录凭证
// 返回值：
//   - string: Cookie 的值
//   - error: 获取 Cookie 时的错误（GoFrame v2 的 Get 方法返回 gvar.Var 类型，需要调用 String() 方法）
func (gf2 *GF2) GetCookie() (string, error) {
	return gf2.ctx.Cookie.Get(gf2.CookieKey()).String(), nil
}

// Lang 实现了 Adapter.Lang 方法
// 该方法用于从查询参数中获取语言设置
// GoAdmin 支持多语言，语言通过 URL 查询参数 __ga_lang 指定
// 返回值：
//   - string: 语言代码，如 "zh-CN"、"en-US" 等
func (gf2 *GF2) Lang() string {
	return gf2.ctx.Request.URL.Query().Get("__ga_lang")
}

// Path 实现了 Adapter.Path 方法
// 该方法用于获取当前请求的路径
// 返回值：
//   - string: 请求路径，如 "/admin/info"
func (gf2 *GF2) Path() string {
	return gf2.ctx.URL.Path
}

// Method 实现了 Adapter.Method 方法
// 该方法用于获取当前请求的 HTTP 方法
// 返回值：
//   - string: HTTP 方法，如 "GET"、"POST"、"PUT"、"DELETE" 等
func (gf2 *GF2) Method() string {
	return gf2.ctx.Method
}

// FormParam 实现了 Adapter.FormParam 方法
// 该方法用于获取表单参数
// 表单参数通常来自 POST 请求的请求体
// GoFrame v2 的 Form 字段已经包含了解析后的表单参数
// 返回值：
//   - url.Values: 表单参数的键值对集合
func (gf2 *GF2) FormParam() url.Values {
	return gf2.ctx.Form
}

// Query 实现了 Adapter.Query 方法
// 该方法用于获取 URL 查询参数
// 查询参数是 URL 中 ? 后面的键值对
// 返回值：
//   - url.Values: 查询参数的键值对集合
func (gf2 *GF2) Query() url.Values {
	return gf2.ctx.Request.URL.Query()
}

// IsPjax 实现了 Adapter.IsPjax 方法
// 该方法用于判断当前请求是否为 PJAX 请求
// PJAX (PushState + AJAX) 是一种技术，允许在不刷新整个页面的情况下更新页面内容
// GoAdmin 使用 PJAX 来提供更流畅的用户体验
// 返回值：
//   - bool: 如果是 PJAX 请求返回 true，否则返回 false
func (gf2 *GF2) IsPjax() bool {
	return gf2.ctx.Header.Get(constant.PjaxHeader) == "true"
}

// Redirect 实现了 Adapter.Redirect 方法
// 该方法用于重定向到登录页面
// 当用户未登录或会话过期时，GoAdmin 会调用此方法将用户重定向到登录页面
func (gf2 *GF2) Redirect() {
	gf2.ctx.Response.RedirectTo(config.Url(config.GetLoginUrl()))
}

// SetContentType 实现了 Adapter.SetContentType 方法
// 该方法用于设置响应的内容类型
// GoAdmin 默认使用 HTML 内容类型，确保浏览器正确渲染页面
func (gf2 *GF2) SetContentType() {
	gf2.ctx.Response.Header().Add("Content-Type", gf2.HTMLContentType())
}

// Write 实现了 Adapter.Write 方法
// 该方法用于写入响应体
// 参数：
//   - body: 要写入的响应体字节数组
func (gf2 *GF2) Write(body []byte) {
	gf2.ctx.Response.WriteStatus(http.StatusOK, body)
}

// Request 实现了 Adapter.Request 方法
// 该方法用于获取原始的 HTTP 请求对象
// 返回值：
//   - *http.Request: 标准 HTTP 请求对象指针
func (gf2 *GF2) Request() *http.Request {
	return gf2.ctx.Request
}

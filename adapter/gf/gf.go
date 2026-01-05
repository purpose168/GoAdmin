// 版权所有 2019 GoAdmin 核心团队。保留所有权利。
// 本源代码的使用受 Apache-2.0 风格许可证约束
// 该许可证可在 LICENSE 文件中找到

// gf 包提供了 GoAdmin 与 GoFrame (GF) 框架的适配器实现
//
// 该适配器作为 GoAdmin 管理后台与 GoFrame 框架之间的桥梁，使得 GoAdmin
// 可以在 GoFrame 应用中无缝运行。它实现了 adapter.WebFrameWork 接口，
// 提供了路由处理、请求上下文管理、响应写入等核心功能。
//
// 主要功能：
//   - 路由注册：将 GoAdmin 的路由注册到 GoFrame 框架中
//   - 请求处理：处理 HTTP 请求并传递给 GoAdmin 引擎
//   - 响应写入：将 GoAdmin 的响应写入到 GoFrame 的响应对象中
//   - 上下文管理：管理请求上下文和用户会话
//   - 参数提取：从路由参数和查询参数中提取数据
//
// 使用示例：
//
//	import (
//	    "github.com/gogf/gf/net/ghttp"
//	    gfadapter "github.com/purpose168/GoAdmin/adapter/gf"
//	    "github.com/purpose168/GoAdmin/plugins/admin"
//	)
//
//	func main() {
//	    s := ghttp.GetServer()
//
//	    // 使用 GF 适配器初始化 GoAdmin
//	    admin.SetAdapter(gfadapter.New())
//
//	    // 添加路由
//	    admin.AddHandler("GET", "/admin", func(ctx *ghttp.Request) {
//	        // 处理逻辑
//	    })
//
//	    s.Run()
//	}
//
// 注意事项：
//   - GoFrame 的路由参数格式为 :param，适配器会自动转换为 URL 查询参数
//   - 适配器支持 PJAX 请求，提供更流畅的用户体验
//   - 需要确保传入的参数类型正确，否则会返回错误或 panic
package gf

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gogf/gf/net/ghttp"
	"github.com/purpose168/GoAdmin/adapter"
	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/engine"
	"github.com/purpose168/GoAdmin/modules/config"
	"github.com/purpose168/GoAdmin/modules/utils"
	"github.com/purpose168/GoAdmin/plugins"
	"github.com/purpose168/GoAdmin/plugins/admin/models"
	"github.com/purpose168/GoAdmin/plugins/admin/modules/constant"
	"github.com/purpose168/GoAdmin/template/types"
)

// Gf 结构体实现了 GoAdmin 的适配器接口
// 它作为 GoFrame (GF) 框架和 GoAdmin 管理后台之间的桥梁
// 嵌入了 adapter.BaseAdapter 以获得基础适配器功能
// ctx 字段存储当前的 GF 请求上下文
// app 字段存储 GF 服务器实例
type Gf struct {
	adapter.BaseAdapter
	ctx *ghttp.Request
	app *ghttp.Server
}

// init 函数在包导入时自动执行
// 它将 Gf 适配器注册到 GoAdmin 引擎中
// 这样 GoAdmin 就能够识别和使用 GoFrame 框架
func init() {
	engine.Register(new(Gf))
}

// User 实现了 Adapter.User 方法
// 该方法从请求上下文中获取当前登录的用户信息
// 参数：
//   - ctx: 请求上下文接口，实际类型为 *ghttp.Request
//
// 返回值：
//   - models.UserModel: 用户模型，包含用户信息
//   - bool: 是否成功获取用户信息
func (gf *Gf) User(ctx interface{}) (models.UserModel, bool) {
	return gf.GetUser(ctx, gf)
}

// Use 实现了 Adapter.Use 方法
// 该方法用于初始化并使用 GoAdmin 插件
// 参数：
//   - app: Web 应用实例，实际类型为 *ghttp.Server
//   - plugs: 插件列表，包含要加载的 GoAdmin 插件
//
// 返回值：
//   - error: 初始化过程中的错误，成功则为 nil
func (gf *Gf) Use(app interface{}, plugs []plugins.Plugin) error {
	return gf.GetUse(app, plugs, gf)
}

// Content 实现了 Adapter.Content 方法
// 该方法用于渲染管理面板内容
// 参数：
//   - ctx: 请求上下文接口
//   - getPanelFn: 获取面板的函数，用于生成管理面板
//   - fn: 节点处理器，用于处理面板中的节点
//   - btns: 按钮列表，用于面板上的操作按钮
func (gf *Gf) Content(ctx interface{}, getPanelFn types.GetPanelFn, fn context.NodeProcessor, btns ...types.Button) {
	gf.GetContent(ctx, getPanelFn, gf, btns, fn)
}

// HandlerFunc 定义了 GoFrame 框架的处理器函数类型
// 它接收 GF 请求对象并返回面板和可能的错误
// 参数：
//   - ctx: GoFrame 请求对象
//
// 返回值：
//   - types.Panel: 管理面板
//   - error: 处理过程中的错误
type HandlerFunc func(ctx *ghttp.Request) (types.Panel, error)

// Content 是一个辅助函数，用于创建 GF 处理器函数
// 该处理器将 GF 的请求处理与 GoAdmin 的内容渲染集成
// 参数：
//   - handler: 处理器函数，用于生成管理面板
//
// 返回值：
//   - ghttp.HandlerFunc: GF 处理器函数类型
func Content(handler HandlerFunc) ghttp.HandlerFunc {
	return func(ctx *ghttp.Request) {
		// 调用 GoAdmin 引擎的内容处理方法
		// 将 GF 请求对象转换为通用接口类型
		engine.Content(ctx, func(ctx interface{}) (types.Panel, error) {
			// 类型断言：将通用接口转换回 GF 请求对象
			return handler(ctx.(*ghttp.Request))
		})
	}
}

// SetApp 实现了 Adapter.SetApp 方法
// 该方法用于设置 GoFrame 服务器实例到适配器中
// 参数：
//   - app: Web 应用实例，实际类型应为 *ghttp.Server
//
// 返回值：
//   - error: 如果参数类型错误则返回错误，成功则为 nil
func (gf *Gf) SetApp(app interface{}) error {
	var (
		eng *ghttp.Server
		ok  bool
	)

	// 类型断言：验证传入的参数是否为 *ghttp.Server 类型
	// GoFrame 使用断言来确保类型安全，这是 Go 语言中处理接口的常见模式
	if eng, ok = app.(*ghttp.Server); !ok {
		return errors.New("gf 适配器 SetApp: 参数类型错误")
	}

	// 保存 GoFrame 服务器实例
	gf.app = eng
	return nil
}

// AddHandler 实现了 Adapter.AddHandler 方法
// 该方法用于向 GoFrame 应用添加路由处理器
// GoFrame 的路由参数格式为 :param，需要转换为标准 URL 查询参数格式 ?param=value
// 参数：
//   - method: HTTP 方法，如 "GET"、"POST" 等
//   - path: 路由路径，如 "/admin" 等
//   - handlers: 处理器链，包含要执行的中间件和处理器
func (gf *Gf) AddHandler(method, path string, handlers context.Handlers) {
	// 使用 GoFrame 的 BindHandler 方法绑定路由处理器
	// GoFrame 使用 "METHOD:path" 格式来绑定路由
	gf.app.BindHandler(strings.ToUpper(method)+":"+path, func(c *ghttp.Request) {
		// 创建 GoAdmin 上下文，传入 GF 的标准 HTTP 请求对象
		ctx := context.NewContext(c.Request)

		// 复制路径用于参数提取
		newPath := path

		// 使用正则表达式匹配路由参数
		// GoFrame 的路由参数格式为 :param，例如 /user/:id
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
			// 使用 GF 的 GetRequestString 方法获取路由参数值
			if c.Request.URL.RawQuery == "" {
				c.Request.URL.RawQuery += p + "=" + c.GetRequestString(p)
			} else {
				c.Request.URL.RawQuery += "&" + p + "=" + c.GetRequestString(p)
			}
		}

		// 设置处理器链并执行
		// GoAdmin 使用处理器链模式，每个处理器可以处理请求并传递给下一个
		ctx.SetHandlers(handlers).Next()

		// 将 GoAdmin 响应头复制到 GF 响应中
		for key, head := range ctx.Response.Header {
			c.Response.Header().Add(key, head[0])
		}

		// 如果响应体不为空，则写入响应
		if ctx.Response.Body != nil {
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(ctx.Response.Body)
			// 使用 GF 的 WriteStatus 方法写入响应状态码和响应体
			c.Response.WriteStatus(ctx.Response.StatusCode, buf.Bytes())
		} else {
			// 如果响应体为空，只写入状态码
			c.Response.WriteStatus(ctx.Response.StatusCode)
		}
	})
}

// Name 实现了 Adapter.Name 方法
// 返回适配器的名称，用于标识不同的框架适配器
// 返回值：
//   - string: 适配器名称 "gf"
func (*Gf) Name() string {
	return "gf"
}

// SetContext 实现了 Adapter.SetContext 方法
// 该方法用于设置当前请求的上下文到适配器中
// 参数：
//   - contextInterface: 请求上下文接口，实际类型应为 *ghttp.Request
//
// 返回值：
//   - adapter.WebFrameWork: 返回新的适配器实例，包含设置的上下文
func (*Gf) SetContext(contextInterface interface{}) adapter.WebFrameWork {
	var (
		ctx *ghttp.Request
		ok  bool
	)

	// 类型断言：验证传入的参数是否为 *ghttp.Request 类型
	// 如果类型不匹配，使用 panic 终止程序，这是 Go 语言中处理严重错误的常见方式
	if ctx, ok = contextInterface.(*ghttp.Request); !ok {
		panic("gf 适配器 SetContext: 参数类型错误")
	}

	// 返回新的 Gf 适配器实例，包含设置的上下文
	return &Gf{ctx: ctx}
}

// Redirect 实现了 Adapter.Redirect 方法
// 该方法用于重定向到登录页面
// 当用户未登录或会话过期时，GoAdmin 会调用此方法将用户重定向到登录页面
func (gf *Gf) Redirect() {
	gf.ctx.Response.RedirectTo(config.Url(config.GetLoginUrl()))
}

// SetContentType 实现了 Adapter.SetContentType 方法
// 该方法用于设置响应的内容类型
// GoAdmin 默认使用 HTML 内容类型，确保浏览器正确渲染页面
func (gf *Gf) SetContentType() {
	gf.ctx.Response.Header().Add("Content-Type", gf.HTMLContentType())
}

// Write 实现了 Adapter.Write 方法
// 该方法用于写入响应体
// 参数：
//   - body: 要写入的响应体字节数组
func (gf *Gf) Write(body []byte) {
	gf.ctx.Response.WriteStatus(http.StatusOK, body)
}

// GetCookie 实现了 Adapter.GetCookie 方法
// 该方法用于获取指定名称的 Cookie 值
// Cookie 用于存储用户会话信息，如登录凭证
// 返回值：
//   - string: Cookie 的值
//   - error: 获取 Cookie 时的错误（GoFrame 的 Get 方法不返回错误，因此这里总是返回 nil）
func (gf *Gf) GetCookie() (string, error) {
	return gf.ctx.Cookie.Get(gf.CookieKey()), nil
}

// Lang 实现了 Adapter.Lang 方法
// 该方法用于从查询参数中获取语言设置
// GoAdmin 支持多语言，语言通过 URL 查询参数 __ga_lang 指定
// 返回值：
//   - string: 语言代码，如 "zh-CN"、"en-US" 等
func (gf *Gf) Lang() string {
	return gf.ctx.Request.URL.Query().Get("__ga_lang")
}

// Path 实现了 Adapter.Path 方法
// 该方法用于获取当前请求的路径
// 返回值：
//   - string: 请求路径，如 "/admin/info"
func (gf *Gf) Path() string {
	return gf.ctx.URL.Path
}

// Method 实现了 Adapter.Method 方法
// 该方法用于获取当前请求的 HTTP 方法
// 返回值：
//   - string: HTTP 方法，如 "GET"、"POST"、"PUT"、"DELETE" 等
func (gf *Gf) Method() string {
	return gf.ctx.Method
}

// FormParam 实现了 Adapter.FormParam 方法
// 该方法用于获取表单参数
// 表单参数通常来自 POST 请求的请求体
// GoFrame 的 Form 字段已经包含了解析后的表单参数
// 返回值：
//   - url.Values: 表单参数的键值对集合
func (gf *Gf) FormParam() url.Values {
	return gf.ctx.Form
}

// IsPjax 实现了 Adapter.IsPjax 方法
// 该方法用于判断当前请求是否为 PJAX 请求
// PJAX (PushState + AJAX) 是一种技术，允许在不刷新整个页面的情况下更新页面内容
// GoAdmin 使用 PJAX 来提供更流畅的用户体验
// 返回值：
//   - bool: 如果是 PJAX 请求返回 true，否则返回 false
func (gf *Gf) IsPjax() bool {
	return gf.ctx.Header.Get(constant.PjaxHeader) == "true"
}

// Query 实现了 Adapter.Query 方法
// 该方法用于获取 URL 查询参数
// 查询参数是 URL 中 ? 后面的键值对
// 返回值：
//   - url.Values: 查询参数的键值对集合
func (gf *Gf) Query() url.Values {
	return gf.ctx.Request.URL.Query()
}

// Request 实现了 Adapter.Request 方法
// 该方法用于获取原始的 HTTP 请求对象
// 返回值：
//   - *http.Request: 标准 HTTP 请求对象指针
func (gf *Gf) Request() *http.Request {
	return gf.ctx.Request
}

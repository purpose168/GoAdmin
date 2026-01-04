/***
# 文件名称: ../../adapter/gear/gear.go
# 作者: eavesmy
# 邮箱: eavesmy@gmail.com
# 创建时间: 2021年06月03日 星期四 19时05分06秒
# 描述: Gear 框架的 GoAdmin 适配器实现
# Gear 是一个轻量级的 Go Web 框架，本文件实现了 GoAdmin 与 Gear 框架的集成
# 主要功能包括：
#   - 请求上下文转换
#   - 路由参数处理
#   - HTTP 响应管理
#   - Cookie 和会话管理
#   - PJAX 支持
***/

package gear

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/purpose168/GoAdmin/adapter"
	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/engine"
	"github.com/purpose168/GoAdmin/modules/config"
	"github.com/purpose168/GoAdmin/modules/utils"
	"github.com/purpose168/GoAdmin/plugins"
	"github.com/purpose168/GoAdmin/plugins/admin/models"
	"github.com/purpose168/GoAdmin/plugins/admin/modules/constant"
	"github.com/purpose168/GoAdmin/template/types"
	"github.com/teambition/gear"
)

// Gear 结构体实现了 GoAdmin 的适配器接口
// 它作为 Gear 框架和 GoAdmin 管理后台之间的桥梁
// 嵌入了 adapter.BaseAdapter 以获得基础适配器功能
// ctx 字段存储当前的 Gear 请求上下文
// app 字段存储 Gear 应用实例
// router 字段存储 Gear 路由器实例，用于注册路由
type Gear struct {
	adapter.BaseAdapter
	ctx    *gear.Context
	app    *gear.App
	router *gear.Router
}

// init 函数在包导入时自动执行
// 它将 Gear 适配器注册到 GoAdmin 引擎中
// 这样 GoAdmin 就能够识别和使用 Gear 框架
func init() {
	engine.Register(new(Gear))
}

// User 实现了 Adapter.User 方法
// 该方法从请求上下文中获取当前登录的用户信息
// 参数：
//   - ctx: 请求上下文接口，实际类型为 *gear.Context
//
// 返回值：
//   - models.UserModel: 用户模型，包含用户信息
//   - bool: 是否成功获取用户信息
func (gears *Gear) User(ctx interface{}) (models.UserModel, bool) {
	return gears.GetUser(ctx, gears)
}

// Use 实现了 Adapter.Use 方法
// 该方法用于初始化并使用 GoAdmin 插件
// 参数：
//   - app: Web 应用实例，实际类型为 *gear.App
//   - plugs: 插件列表，包含要加载的 GoAdmin 插件
//
// 返回值：
//   - error: 初始化过程中的错误，成功则为 nil
func (gears *Gear) Use(app interface{}, plugs []plugins.Plugin) error {
	return gears.GetUse(app, plugs, gears)
}

// Content 实现了 Adapter.Content 方法
// 该方法用于渲染管理面板内容
// 参数：
//   - ctx: 请求上下文接口
//   - getPanelFn: 获取面板的函数，用于生成管理面板
//   - fn: 节点处理器，用于处理面板中的节点
//   - btns: 按钮列表，用于面板上的操作按钮
func (gears *Gear) Content(ctx interface{}, getPanelFn types.GetPanelFn, fn context.NodeProcessor, btns ...types.Button) {
	gears.GetContent(ctx, getPanelFn, gears, btns, fn)
}

// HandlerFunc 定义了 Gear 框架的处理器函数类型
// 它接收 Gear 上下文并返回面板和可能的错误
// 参数：
//   - ctx: Gear 请求上下文
//
// 返回值：
//   - types.Panel: 管理面板
//   - error: 处理过程中的错误
type HandlerFunc func(ctx *gear.Context) (types.Panel, error)

// Content 是一个辅助函数，用于创建 Gear 中间件
// 该中间件将 Gear 的请求处理与 GoAdmin 的内容渲染集成
// 参数：
//   - handler: 处理器函数，用于生成管理面板
//
// 返回值：
//   - gear.Middleware: Gear 中间件函数
func Content(handler HandlerFunc) gear.Middleware {
	return func(ctx *gear.Context) error {
		// 调用 GoAdmin 引擎的内容处理方法
		// 将 Gear 上下文转换为通用接口类型
		engine.Content(ctx, func(ctx interface{}) (types.Panel, error) {
			// 类型断言：将通用接口转换回 Gear 上下文
			return handler(ctx.(*gear.Context))
		})
		return nil
	}
}

// SetApp 实现了 Adapter.SetApp 方法
// 该方法用于设置 Gear 应用实例到适配器中
// 参数：
//   - app: Web 应用实例，实际类型应为 *gear.App
//
// 返回值：
//   - error: 如果参数类型错误则返回错误，成功则为 nil
func (gears *Gear) SetApp(app interface{}) error {
	// 初始化 Gear 路由器
	gears.router = gear.NewRouter()

	var (
		eng *gear.App
		ok  bool
	)

	// 类型断言：验证传入的参数是否为 *gear.App 类型
	// Gear 使用断言来确保类型安全，这是 Go 语言中处理接口的常见模式
	if eng, ok = app.(*gear.App); !ok {
		return errors.New("gear 适配器 SetApp: 参数类型错误")
	}

	// 保存 Gear 应用实例
	gears.app = eng
	return nil
}

// AddHandler 实现了 Adapter.AddHandler 方法
// 该方法用于向 Gear 应用添加路由处理器
// Gear 的路由参数格式为 :param，需要转换为标准 URL 查询参数格式 ?param=value
// 参数：
//   - method: HTTP 方法，如 "GET"、"POST" 等
//   - path: 路由路径，如 "/admin" 等
//   - handlers: 处理器链，包含要执行的中间件和处理器
func (gears *Gear) AddHandler(method, path string, handlers context.Handlers) {

	// 如果路由器未初始化，则创建新的路由器
	if gears.router == nil {
		gears.router = gear.NewRouter()
	}

	// 使用 Gear 路由器注册路由处理器
	// Gear 使用链式处理器模式，每个处理器可以处理请求并传递给下一个
	gears.router.Handle(strings.ToUpper(method), path, func(c *gear.Context) error {

		// 创建 GoAdmin 上下文，传入 Gear 的请求对象
		ctx := context.NewContext(c.Req)

		// 复制路径用于参数提取
		newPath := path

		// 使用正则表达式匹配路由参数
		// Gear 的路由参数格式为 :param，例如 /user/:id
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
			if c.Req.URL.RawQuery == "" {
				c.Req.URL.RawQuery += p + "=" + c.Param(p)
			} else {
				c.Req.URL.RawQuery += "&" + p + "=" + c.Param(p)
			}
		}

		// 设置处理器链并执行
		// GoAdmin 使用处理器链模式，每个处理器可以处理请求并传递给下一个
		ctx.SetHandlers(handlers).Next()

		// 将 GoAdmin 响应头复制到 Gear 响应中
		for key, head := range ctx.Response.Header {
			c.Res.Header().Add(key, head[0])
		}

		// 如果响应体不为空，则写入响应
		if ctx.Response.Body != nil {
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(ctx.Response.Body)

			// 使用 Gear 的 End 方法结束请求并返回响应
			return c.End(ctx.Response.StatusCode, buf.Bytes())
		}

		// 设置响应状态码
		c.Status(ctx.Response.StatusCode)
		return nil
	})

	// 将路由器作为中间件添加到 Gear 应用中
	// 这样 Gear 应用就能够处理 GoAdmin 的路由
	gears.app.UseHandler(gears.router)
}

// Name 实现了 Adapter.Name 方法
// 返回适配器的名称，用于标识不同的框架适配器
// 返回值：
//   - string: 适配器名称 "gear"
func (*Gear) Name() string {
	return "gear"
}

// SetContext 实现了 Adapter.SetContext 方法
// 该方法用于设置当前请求的上下文到适配器中
// 参数：
//   - contextInterface: 请求上下文接口，实际类型应为 *gear.Context
//
// 返回值：
//   - adapter.WebFrameWork: 返回新的适配器实例，包含设置的上下文
func (*Gear) SetContext(contextInterface interface{}) adapter.WebFrameWork {
	var (
		ctx *gear.Context
		ok  bool
	)

	// 类型断言：验证传入的参数是否为 *gear.Context 类型
	// 如果类型不匹配，使用 panic 终止程序，这是 Go 语言中处理严重错误的常见方式
	if ctx, ok = contextInterface.(*gear.Context); !ok {
		panic("gear 适配器 SetContext: 参数类型错误")
	}

	// 返回新的 Gear 适配器实例，包含设置的上下文
	return &Gear{ctx: ctx}
}

// Redirect 实现了 Adapter.Redirect 方法
// 该方法用于重定向到登录页面
// 当用户未登录或会话过期时，GoAdmin 会调用此方法将用户重定向到登录页面
func (gears *Gear) Redirect() {
	gears.ctx.Redirect(config.Url(config.GetLoginUrl()))
}

// SetContentType 实现了 Adapter.SetContentType 方法
// 该方法用于设置响应的内容类型
// GoAdmin 默认使用 HTML 内容类型，确保浏览器正确渲染页面
func (gears *Gear) SetContentType() {
	gears.ctx.Res.Header().Set("Content-Type", gears.HTMLContentType())
}

// Write 实现了 Adapter.Write 方法
// 该方法用于写入响应体
// 参数：
//   - body: 要写入的响应体字节数组
func (gears *Gear) Write(body []byte) {
	gears.ctx.End(http.StatusOK, body)
}

// GetCookie 实现了 Adapter.GetCookie 方法
// 该方法用于获取指定名称的 Cookie 值
// Cookie 用于存储用户会话信息，如登录凭证
// 返回值：
//   - string: Cookie 的值
//   - error: 获取 Cookie 时的错误
func (gears *Gear) GetCookie() (string, error) {
	return gears.ctx.Cookies.Get(gears.CookieKey())
}

// Lang 实现了 Adapter.Lang 方法
// 该方法用于从查询参数中获取语言设置
// GoAdmin 支持多语言，语言通过 URL 查询参数 __ga_lang 指定
// 返回值：
//   - string: 语言代码，如 "zh-CN"、"en-US" 等
func (gears *Gear) Lang() string {
	return gears.ctx.Req.URL.Query().Get("__ga_lang")
}

// Path 实现了 Adapter.Path 方法
// 该方法用于获取当前请求的路径
// 返回值：
//   - string: 请求路径，如 "/admin/info"
func (gears *Gear) Path() string {
	return gears.ctx.Req.URL.Path
}

// Method 实现了 Adapter.Method 方法
// 该方法用于获取当前请求的 HTTP 方法
// 返回值：
//   - string: HTTP 方法，如 "GET"、"POST"、"PUT"、"DELETE" 等
func (gears *Gear) Method() string {
	return gears.ctx.Req.Method
}

// FormParam 实现了 Adapter.FormParam 方法
// 该方法用于获取表单参数
// 表单参数通常来自 POST 请求的请求体
// ParseMultipartForm 用于解析 multipart/form-data 格式的表单数据
// 32 << 20 表示最大内存限制为 32MB
// 返回值：
//   - url.Values: 表单参数的键值对集合
func (gears *Gear) FormParam() url.Values {
	_ = gears.ctx.Req.ParseMultipartForm(32 << 20)
	return gears.ctx.Req.PostForm
}

// IsPjax 实现了 Adapter.IsPjax 方法
// 该方法用于判断当前请求是否为 PJAX 请求
// PJAX (PushState + AJAX) 是一种技术，允许在不刷新整个页面的情况下更新页面内容
// GoAdmin 使用 PJAX 来提供更流畅的用户体验
// 返回值：
//   - bool: 如果是 PJAX 请求返回 true，否则返回 false
func (gears *Gear) IsPjax() bool {
	return gears.ctx.Req.Header.Get(constant.PjaxHeader) == "true"
}

// Query 实现了 Adapter.Query 方法
// 该方法用于获取 URL 查询参数
// 查询参数是 URL 中 ? 后面的键值对
// 返回值：
//   - url.Values: 查询参数的键值对集合
func (gears *Gear) Query() url.Values {
	return gears.ctx.Req.URL.Query()
}

// Request 实现了 Adapter.Request 方法
// 该方法用于获取原始的 HTTP 请求对象
// 返回值：
//   - *http.Request: 标准 HTTP 请求对象指针
func (gears *Gear) Request() *http.Request {
	return gears.ctx.Req
}

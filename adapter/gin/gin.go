// 版权所有 2019 GoAdmin 核心团队。保留所有权利。
// 本源代码的使用受 Apache-2.0 风格许可证约束
// 该许可证可在 LICENSE 文件中找到

package gin

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/purpose168/GoAdmin/adapter"
	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/engine"
	"github.com/purpose168/GoAdmin/modules/config"
	"github.com/purpose168/GoAdmin/plugins"
	"github.com/purpose168/GoAdmin/plugins/admin/models"
	"github.com/purpose168/GoAdmin/plugins/admin/modules/constant"
	"github.com/purpose168/GoAdmin/template/types"
)

// Gin 结构体实现了 GoAdmin 的适配器接口
// 它作为 Gin 框架和 GoAdmin 管理后台之间的桥梁
// 嵌入了 adapter.BaseAdapter 以获得基础适配器功能
// ctx 字段存储当前的 Gin 请求上下文
// app 字段存储 Gin 引擎实例
type Gin struct {
	adapter.BaseAdapter
	ctx *gin.Context
	app *gin.Engine
}

// init 函数在包导入时自动执行
// 它将 Gin 适配器注册到 GoAdmin 引擎中
// 这样 GoAdmin 就能够识别和使用 Gin 框架
func init() {
	engine.Register(new(Gin))
}

// User 实现了 Adapter.User 方法
// 该方法从请求上下文中获取当前登录的用户信息
// 参数：
//   - ctx: 请求上下文接口，实际类型为 *gin.Context
//
// 返回值：
//   - models.UserModel: 用户模型，包含用户信息
//   - bool: 是否成功获取用户信息
func (gins *Gin) User(ctx interface{}) (models.UserModel, bool) {
	return gins.GetUser(ctx, gins)
}

// Use 实现了 Adapter.Use 方法
// 该方法用于初始化并使用 GoAdmin 插件
// 参数：
//   - app: Web 应用实例，实际类型为 *gin.Engine
//   - plugs: 插件列表，包含要加载的 GoAdmin 插件
//
// 返回值：
//   - error: 初始化过程中的错误，成功则为 nil
func (gins *Gin) Use(app interface{}, plugs []plugins.Plugin) error {
	return gins.GetUse(app, plugs, gins)
}

// Content 实现了 Adapter.Content 方法
// 该方法用于渲染管理面板内容
// 参数：
//   - ctx: 请求上下文接口
//   - getPanelFn: 获取面板的函数，用于生成管理面板
//   - fn: 节点处理器，用于处理面板中的节点
//   - btns: 按钮列表，用于面板上的操作按钮
func (gins *Gin) Content(ctx interface{}, getPanelFn types.GetPanelFn, fn context.NodeProcessor, btns ...types.Button) {
	gins.GetContent(ctx, getPanelFn, gins, btns, fn)
}

// HandlerFunc 定义了 Gin 框架的处理器函数类型
// 它接收 Gin 上下文并返回面板和可能的错误
// 参数：
//   - ctx: Gin 请求上下文
//
// 返回值：
//   - types.Panel: 管理面板
//   - error: 处理过程中的错误
type HandlerFunc func(ctx *gin.Context) (types.Panel, error)

// Content 是一个辅助函数，用于创建 Gin 处理器函数
// 该处理器将 Gin 的请求处理与 GoAdmin 的内容渲染集成
// 参数：
//   - handler: 处理器函数，用于生成管理面板
//
// 返回值：
//   - gin.HandlerFunc: Gin 处理器函数类型
func Content(handler HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 调用 GoAdmin 引擎的内容处理方法
		// 将 Gin 上下文转换为通用接口类型
		engine.Content(ctx, func(ctx interface{}) (types.Panel, error) {
			// 类型断言：将通用接口转换回 Gin 上下文
			return handler(ctx.(*gin.Context))
		})
	}
}

// SetApp 实现了 Adapter.SetApp 方法
// 该方法用于设置 Gin 引擎实例到适配器中
// 参数：
//   - app: Web 应用实例，实际类型应为 *gin.Engine
//
// 返回值：
//   - error: 如果参数类型错误则返回错误，成功则为 nil
func (gins *Gin) SetApp(app interface{}) error {
	var (
		eng *gin.Engine
		ok  bool
	)

	// 类型断言：验证传入的参数是否为 *gin.Engine 类型
	// Gin 使用断言来确保类型安全，这是 Go 语言中处理接口的常见模式
	if eng, ok = app.(*gin.Engine); !ok {
		return errors.New("gin 适配器 SetApp: 参数类型错误")
	}

	// 保存 Gin 引擎实例
	gins.app = eng
	return nil
}

// AddHandler 实现了 Adapter.AddHandler 方法
// 该方法用于向 Gin 应用添加路由处理器
// Gin 的路由参数已经解析在 c.Params 中，需要转换为标准 URL 查询参数格式 ?param=value
// 参数：
//   - method: HTTP 方法，如 "GET"、"POST" 等
//   - path: 路由路径，如 "/admin" 等
//   - handlers: 处理器链，包含要执行的中间件和处理器
func (gins *Gin) AddHandler(method, path string, handlers context.Handlers) {
	// 使用 Gin 的 Handle 方法注册路由处理器
	gins.app.Handle(strings.ToUpper(method), path, func(c *gin.Context) {
		// 创建 GoAdmin 上下文，传入 Gin 的标准 HTTP 请求对象
		ctx := context.NewContext(c.Request)

		// 将 Gin 的路由参数转换为 URL 查询参数
		// Gin 已经将路由参数解析到 c.Params 中，每个参数包含 Key 和 Value
		// 例如：/user/:id → c.Params = [{Key: ":id", Value: "123"}]
		// 需要转换为：/user?id=123
		for _, param := range c.Params {
			// 移除参数名中的冒号
			// 如果是第一个参数，直接添加；后续参数用 & 连接
			if c.Request.URL.RawQuery == "" {
				c.Request.URL.RawQuery += strings.ReplaceAll(param.Key, ":", "") + "=" + param.Value
			} else {
				c.Request.URL.RawQuery += "&" + strings.ReplaceAll(param.Key, ":", "") + "=" + param.Value
			}
		}

		// 设置处理器链并执行
		// GoAdmin 使用处理器链模式，每个处理器可以处理请求并传递给下一个
		ctx.SetHandlers(handlers).Next()

		// 将 GoAdmin 响应头复制到 Gin 响应中
		for key, head := range ctx.Response.Header {
			c.Header(key, head[0])
		}

		// 如果响应体不为空，则写入响应
		if ctx.Response.Body != nil {
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(ctx.Response.Body)
			// 使用 Gin 的 String 方法写入响应状态码和响应体
			c.String(ctx.Response.StatusCode, buf.String())
		} else {
			// 如果响应体为空，只写入状态码
			c.Status(ctx.Response.StatusCode)
		}
	})
}

// Name 实现了 Adapter.Name 方法
// 返回适配器的名称，用于标识不同的框架适配器
// 返回值：
//   - string: 适配器名称 "gin"
func (*Gin) Name() string {
	return "gin"
}

// SetContext 实现了 Adapter.SetContext 方法
// 该方法用于设置当前请求的上下文到适配器中
// 参数：
//   - contextInterface: 请求上下文接口，实际类型应为 *gin.Context
//
// 返回值：
//   - adapter.WebFrameWork: 返回新的适配器实例，包含设置的上下文
func (*Gin) SetContext(contextInterface interface{}) adapter.WebFrameWork {
	var (
		ctx *gin.Context
		ok  bool
	)

	// 类型断言：验证传入的参数是否为 *gin.Context 类型
	// 如果类型不匹配，使用 panic 终止程序，这是 Go 语言中处理严重错误的常见方式
	if ctx, ok = contextInterface.(*gin.Context); !ok {
		panic("gin 适配器 SetContext: 参数类型错误")
	}

	// 返回新的 Gin 适配器实例，包含设置的上下文
	return &Gin{ctx: ctx}
}

// Redirect 实现了 Adapter.Redirect 方法
// 该方法用于重定向到登录页面
// 当用户未登录或会话过期时，GoAdmin 会调用此方法将用户重定向到登录页面
// 使用 http.StatusFound (302) 状态码进行临时重定向
// Abort() 方法用于停止后续的处理器执行
func (gins *Gin) Redirect() {
	gins.ctx.Redirect(http.StatusFound, config.Url(config.GetLoginUrl()))
	gins.ctx.Abort()
}

// SetContentType 实现了 Adapter.SetContentType 方法
// 该方法用于设置响应的内容类型
// Gin 框架会自动设置内容类型，因此这里不需要额外操作
// GoAdmin 默认使用 HTML 内容类型，确保浏览器正确渲染页面
func (*Gin) SetContentType() {}

// Write 实现了 Adapter.Write 方法
// 该方法用于写入响应体
// 参数：
//   - body: 要写入的响应体字节数组
func (gins *Gin) Write(body []byte) {
	gins.ctx.Data(http.StatusOK, gins.HTMLContentType(), body)
}

// GetCookie 实现了 Adapter.GetCookie 方法
// 该方法用于获取指定名称的 Cookie 值
// Cookie 用于存储用户会话信息，如登录凭证
// 返回值：
//   - string: Cookie 的值
//   - error: 获取 Cookie 时的错误
func (gins *Gin) GetCookie() (string, error) {
	return gins.ctx.Cookie(gins.CookieKey())
}

// Lang 实现了 Adapter.Lang 方法
// 该方法用于从查询参数中获取语言设置
// GoAdmin 支持多语言，语言通过 URL 查询参数 __ga_lang 指定
// 返回值：
//   - string: 语言代码，如 "zh-CN"、"en-US" 等
func (gins *Gin) Lang() string {
	return gins.ctx.Request.URL.Query().Get("__ga_lang")
}

// Path 实现了 Adapter.Path 方法
// 该方法用于获取当前请求的路径
// 返回值：
//   - string: 请求路径，如 "/admin/info"
func (gins *Gin) Path() string {
	return gins.ctx.Request.URL.Path
}

// Method 实现了 Adapter.Method 方法
// 该方法用于获取当前请求的 HTTP 方法
// 返回值：
//   - string: HTTP 方法，如 "GET"、"POST"、"PUT"、"DELETE" 等
func (gins *Gin) Method() string {
	return gins.ctx.Request.Method
}

// FormParam 实现了 Adapter.FormParam 方法
// 该方法用于获取表单参数
// 表单参数通常来自 POST 请求的请求体
// ParseMultipartForm 用于解析 multipart/form-data 格式的表单数据
// 32 << 20 表示最大内存限制为 32MB
// 返回值：
//   - url.Values: 表单参数的键值对集合
func (gins *Gin) FormParam() url.Values {
	_ = gins.ctx.Request.ParseMultipartForm(32 << 20)
	return gins.ctx.Request.PostForm
}

// IsPjax 实现了 Adapter.IsPjax 方法
// 该方法用于判断当前请求是否为 PJAX 请求
// PJAX (PushState + AJAX) 是一种技术，允许在不刷新整个页面的情况下更新页面内容
// GoAdmin 使用 PJAX 来提供更流畅的用户体验
// 返回值：
//   - bool: 如果是 PJAX 请求返回 true，否则返回 false
func (gins *Gin) IsPjax() bool {
	return gins.ctx.Request.Header.Get(constant.PjaxHeader) == "true"
}

// Query 实现了 Adapter.Query 方法
// 该方法用于获取 URL 查询参数
// 查询参数是 URL 中 ? 后面的键值对
// 返回值：
//   - url.Values: 查询参数的键值对集合
func (gins *Gin) Query() url.Values {
	return gins.ctx.Request.URL.Query()
}

// Request 实现了 Adapter.Request 方法
// 该方法用于获取原始的 HTTP 请求对象
// 返回值：
//   - *http.Request: 标准 HTTP 请求对象指针
func (gins *Gin) Request() *http.Request {
	return gins.ctx.Request
}

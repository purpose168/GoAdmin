// beego2 包提供了 GoAdmin 与 Beego v2 Web 框架的适配器实现
// 该适配器允许 GoAdmin 管理后台在 Beego v2 应用中运行
// 包名：beego2
// 作者：GoAdmin Core Team
// 创建日期：2019
// 目的：为 Beego v2 框架提供 GoAdmin 管理后台的集成支持

package beego2

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/purpose168/GoAdmin/adapter"
	gctx "github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/engine"
	"github.com/purpose168/GoAdmin/modules/config"
	"github.com/purpose168/GoAdmin/modules/constant"
	"github.com/purpose168/GoAdmin/plugins"
	"github.com/purpose168/GoAdmin/plugins/admin/models"
	"github.com/purpose168/GoAdmin/template/types"
)

// Beego2 结构体实现了 GoAdmin 的适配器接口
// 它作为 Beego v2 框架和 GoAdmin 管理后台之间的桥梁
// 嵌入了 adapter.BaseAdapter 以获得基础适配器功能
// ctx 字段存储当前的 Beego v2 上下文
// app 字段存储 Beego v2 应用实例
type Beego2 struct {
	adapter.BaseAdapter
	ctx *context.Context
	app *web.HttpServer
}

// init 函数在包导入时自动执行
// Go 语言的 init 函数会在 main 函数之前自动调用
// 这里使用 init 函数将 Beego2 适配器注册到 GoAdmin 引擎中
// 这样 GoAdmin 就知道如何使用 Beego v2 框架
func init() {
	engine.Register(new(Beego2))
}

// Name 实现了 Adapter.Name 方法
// 该方法返回适配器的名称
// 返回值：
//   - string: 适配器名称，固定为 "beego2"
func (*Beego2) Name() string {
	return "beego2"
}

// Use 实现了 Adapter.Use 方法
// 该方法用于将插件注册到 Beego v2 应用中
// 参数：
//   - app: 应用接口，通常为 *web.HttpServer 类型
//   - plugins: 插件列表，包含需要注册的所有插件
//
// 返回值：
//   - error: 错误信息，如果注册失败则返回错误
func (bee2 *Beego2) Use(app interface{}, plugins []plugins.Plugin) error {
	return bee2.GetUse(app, plugins, bee2)
}

// Content 实现了 Adapter.Content 方法
// 该方法用于渲染管理面板内容
// 参数：
//   - ctx: 上下文接口，通常为 *context.Context 类型
//   - getPanelFn: 获取面板的函数，返回 types.Panel 类型的面板
//   - fn: 节点处理器，用于处理上下文中的节点
//   - navButtons: 导航按钮列表，可变参数
func (bee2 *Beego2) Content(ctx interface{}, getPanelFn types.GetPanelFn, fn gctx.NodeProcessor, navButtons ...types.Button) {
	bee2.GetContent(ctx, getPanelFn, bee2, navButtons, fn)
}

// User 实现了 Adapter.User 方法
// 该方法用于从当前上下文中获取用户信息
// 参数：
//   - ctx: 上下文接口，通常为 *context.Context 类型
//
// 返回值：
//   - models.UserModel: 用户模型，包含用户信息
//   - bool: 是否成功获取用户信息，true 表示成功
func (bee2 *Beego2) User(ctx interface{}) (models.UserModel, bool) {
	return bee2.GetUser(ctx, bee2)
}

// AddHandler 实现了 Adapter.AddHandler 方法
// 该方法用于向 Beego v2 应用添加路由处理器
// 参数：
//   - method: HTTP 方法，如 "GET"、"POST" 等
//   - path: 路由路径，如 "/admin" 等
//   - handlers: 处理器链，包含要执行的中间件和处理器
func (bee2 *Beego2) AddHandler(method, path string, handlers gctx.Handlers) {
	bee2.app.Handlers.AddMethod(method, path, func(c *context.Context) {
		// 将 Beego v2 路由参数转换为 URL 查询参数
		// Beego v2 的路由参数格式为 :param，需要转换为 ?param=value 格式
		for key, value := range c.Input.Params() {
			if c.Request.URL.RawQuery == "" {
				c.Request.URL.RawQuery += strings.ReplaceAll(key, ":", "") + "=" + value
			} else {
				c.Request.URL.RawQuery += "&" + strings.ReplaceAll(key, ":", "") + "=" + value
			}
		}
		// 创建 GoAdmin 上下文并执行处理器链
		ctx := gctx.NewContext(c.Request)
		ctx.SetHandlers(handlers).Next()
		// 将 GoAdmin 响应头复制到 Beego v2 响应中
		for key, head := range ctx.Response.Header {
			c.ResponseWriter.Header().Add(key, head[0])
		}
		// 设置响应状态码
		c.ResponseWriter.WriteHeader(ctx.Response.StatusCode)
		// 将响应体写入 Beego v2 响应
		if ctx.Response.Body != nil {
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(ctx.Response.Body)
			c.WriteString(buf.String())
		}
	})
}

// SetApp 实现了 Adapter.SetApp 方法
// 该方法用于设置 Beego v2 应用实例到适配器中
// 参数：
//   - app: 应用接口，必须为 *web.HttpServer 类型
//
// 返回值：
//   - error: 错误信息，如果参数类型不正确则返回错误
func (bee2 *Beego2) SetApp(app interface{}) error {
	var (
		eng *web.HttpServer
		ok  bool
	)
	// 使用类型断言检查 app 是否为 *web.HttpServer 类型
	// ok 为 true 表示断言成功，eng 为转换后的值
	if eng, ok = app.(*web.HttpServer); !ok {
		return errors.New("beego2 adapter SetApp: wrong parameter")
	}
	bee2.app = eng
	return nil
}

// SetContext 实现了 Adapter.SetContext 方法
// 该方法用于设置当前请求的上下文
// 参数：
//   - contextInterface: 上下文接口，必须为 *context.Context 类型
//
// 返回值：
//   - adapter.WebFrameWork: 返回设置了上下文的新适配器实例
func (*Beego2) SetContext(contextInterface interface{}) adapter.WebFrameWork {
	var (
		ctx *context.Context
		ok  bool
	)
	// 使用类型断言检查 contextInterface 是否为 *context.Context 类型
	if ctx, ok = contextInterface.(*context.Context); !ok {
		panic("beego2 adapter SetContext: wrong parameter")
	}
	return &Beego2{ctx: ctx}
}

// GetCookie 实现了 Adapter.GetCookie 方法
// 该方法用于从请求中获取认证 Cookie
// 返回值：
//   - string: Cookie 值
//   - error: 错误信息，当前实现总是返回 nil
func (bee2 *Beego2) GetCookie() (string, error) {
	return bee2.ctx.GetCookie(bee2.CookieKey()), nil
}

// Lang 实现了 Adapter.Lang 方法
// 该方法用于从 URL 查询参数中获取语言设置
// 返回值：
//   - string: 语言代码，如 "zh-CN"、"en-US" 等
func (bee2 *Beego2) Lang() string {
	return bee2.ctx.Request.URL.Query().Get("__ga_lang")
}

// Path 实现了 Adapter.Path 方法
// 该方法用于获取当前请求的路径
// 返回值：
//   - string: 请求路径，如 "/admin/dashboard"
func (bee2 *Beego2) Path() string {
	return bee2.ctx.Request.URL.Path
}

// Method 实现了 Adapter.Method 方法
// 该方法用于获取当前请求的 HTTP 方法
// 返回值：
//   - string: HTTP 方法，如 "GET"、"POST"、"PUT"、"DELETE" 等
func (bee2 *Beego2) Method() string {
	return bee2.ctx.Request.Method
}

// FormParam 实现了 Adapter.FormParam 方法
// 该方法用于解析并获取表单参数
// 解析的最大内存限制为 32MB
// 返回值：
//   - url.Values: 表单参数的键值对集合
func (bee2 *Beego2) FormParam() url.Values {
	_ = bee2.ctx.Request.ParseMultipartForm(32 << 20)
	return bee2.ctx.Request.PostForm
}

// Query 实现了 Adapter.Query 方法
// 该方法用于获取 URL 查询参数
// 返回值：
//   - url.Values: 查询参数的键值对集合
func (bee2 *Beego2) Query() url.Values {
	return bee2.ctx.Request.URL.Query()
}

// IsPjax 实现了 Adapter.IsPjax 方法
// 该方法用于检查当前请求是否为 PJAX 请求
// PJAX 是一种使用 AJAX 技术实现页面部分更新的技术
// 返回值：
//   - bool: 如果是 PJAX 请求则返回 true，否则返回 false
func (bee2 *Beego2) IsPjax() bool {
	return bee2.ctx.Request.Header.Get(constant.PjaxHeader) == "true"
}

// Redirect 实现了 Adapter.Redirect 方法
// 该方法用于重定向到登录页面
// 使用 HTTP 302 状态码进行临时重定向
func (bee2 *Beego2) Redirect() {
	bee2.ctx.Redirect(http.StatusFound, config.Url(config.GetLoginUrl()))
}

// SetContentType 实现了 Adapter.SetContentType 方法
// 该方法用于设置响应的 Content-Type 头
// Content-Type 由 HTMLContentType() 方法确定
func (bee2 *Beego2) SetContentType() {
	bee2.ctx.ResponseWriter.Header().Set("Content-Type", bee2.HTMLContentType())
}

// Write 实现了 Adapter.Write 方法
// 该方法用于将响应体写入到响应中
// 参数：
//   - body: 要写入的响应体字节数组
func (bee2 *Beego2) Write(body []byte) {
	_, _ = bee2.ctx.ResponseWriter.Write(body)
}

// Request 实现了 Adapter.Request 方法
// 该方法用于获取原始的 HTTP 请求对象
// 返回值：
//   - *http.Request: HTTP 请求对象指针
func (bee2 *Beego2) Request() *http.Request {
	return bee2.ctx.Request
}

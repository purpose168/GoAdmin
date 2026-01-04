// nethttp包提供了GoAdmin框架与Go标准库net/http之间的适配器实现
// 该适配器允许GoAdmin框架与标准库的http.ServeMux无缝集成
//
// 主要功能：
//   - 提供标准库http.ServeMux的适配器
//   - 支持路由注册和参数提取
//   - 实现请求/响应的上下文封装
//   - 提供用户认证和会话管理接口
//
// 使用场景：
//   - 当项目使用标准库net/http作为HTTP服务器时
//   - 需要轻量级的HTTP框架集成
//   - 不想引入第三方路由库（如gin、chi等）
//
// 注意事项：
//   - 路由参数格式从:param转换为{param}以符合标准库规范
//   - 支持GET、POST、PUT、DELETE等HTTP方法
//   - 自动处理路径参数并注入到查询参数中
//
// 作者: purpose168
// 创建日期: 2020-01-01
// 版本: 1.0.0
package nethttp

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/purpose168/GoAdmin/adapter"
	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/engine"
	cfg "github.com/purpose168/GoAdmin/modules/config"
	"github.com/purpose168/GoAdmin/modules/constant"
	"github.com/purpose168/GoAdmin/plugins"
	"github.com/purpose168/GoAdmin/plugins/admin/models"
	"github.com/purpose168/GoAdmin/template/types"
)

// NetHTTP结构体是GoAdmin框架与net/http标准库之间的适配器
// 它实现了adapter.Adapter接口，将标准库的http.ServeMux封装为框架可用的形式
//
// 字段说明：
//   - BaseAdapter: 基础适配器，提供通用功能
//   - ctx: 上下文对象，封装了HTTP请求和响应
//   - app: http.ServeMux实例，用于路由注册和处理
type NetHTTP struct {
	adapter.BaseAdapter
	ctx Context
	app *http.ServeMux
}

// init函数在包加载时自动执行
// 将NetHTTP适配器注册到引擎中，使其可以被框架识别和使用
func init() {
	engine.Register(new(NetHTTP))
}

// User实现Adapter.User接口方法
// 从上下文中获取当前登录的用户信息
//
// 参数说明：
//   - ctx: 上下文接口，包含HTTP请求和响应信息
//
// 返回值：
//   - models.UserModel: 用户模型，包含用户详细信息
//   - bool: 是否成功获取用户信息，true表示成功，false表示失败
//
// 使用示例：
//
//	user, ok := adapter.User(ctx)
//	if !ok {
//	    // 处理未登录情况
//	}
func (nh *NetHTTP) User(ctx interface{}) (models.UserModel, bool) {
	return nh.GetUser(ctx, nh)
}

// Use实现Adapter.Use接口方法
// 将插件注册到HTTP应用中
//
// 参数说明：
//   - app: HTTP应用实例，通常是*http.ServeMux类型
//   - plugs: 插件列表，包含需要注册的所有插件
//
// 返回值：
//   - error: 注册过程中的错误，成功返回nil
//
// 注意事项：
//   - app参数必须是*http.ServeMux类型，否则会返回错误
//   - 所有插件都会被注册到同一个ServeMux实例中
func (nh *NetHTTP) Use(app interface{}, plugs []plugins.Plugin) error {
	return nh.GetUse(app, plugs, nh)
}

// Content实现Adapter.Content接口方法
// 生成并返回管理面板的内容
//
// 参数说明：
//   - ctx: 上下文接口，包含HTTP请求和响应信息
//   - getPanelFn: 获取面板的函数，用于生成面板内容
//   - fn: 节点处理器，用于处理面板中的节点
//   - btns: 按钮列表，用于在面板上显示操作按钮
//
// 使用场景：
//   - 渲染管理后台页面
//   - 显示数据表格和表单
//   - 处理用户交互操作
func (nh *NetHTTP) Content(ctx interface{}, getPanelFn types.GetPanelFn, fn context.NodeProcessor, btns ...types.Button) {
	nh.GetContent(ctx, getPanelFn, nh, btns, fn)
}

// HandlerFunc是处理HTTP请求的函数类型
// 它接收Context参数并返回Panel和可能的错误
//
// 参数说明：
//   - ctx: 上下文对象，包含HTTP请求和响应信息
//
// 返回值：
//   - types.Panel: 面板对象，包含要显示的内容
//   - error: 处理过程中的错误，成功返回nil
//
// 使用示例：
//
//	handler := func(ctx Context) (types.Panel, error) {
//	    // 处理请求并返回面板
//	    return panel, nil
//	}
type HandlerFunc func(ctx Context) (types.Panel, error)

// Content函数将HandlerFunc转换为标准的http.HandlerFunc
// 这样可以将GoAdmin的面板处理函数注册到标准库的路由中
//
// 参数说明：
//   - handler: HandlerFunc类型的处理函数
//
// 返回值：
//   - http.HandlerFunc: 标准库的HTTP处理函数
//
// 工作原理：
//  1. 创建Context对象封装请求和响应
//  2. 调用engine.Content处理请求
//  3. 将Context转换为interface{}类型传递给handler
//
// 使用示例：
//
//	mux.HandleFunc("/admin", Content(adminHandler))
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

// SetApp实现Adapter.SetApp接口方法
// 设置HTTP应用实例到适配器中
//
// 参数说明：
//   - app: HTTP应用实例，必须是*http.ServeMux类型
//
// 返回值：
//   - error: 设置过程中的错误，成功返回nil
//
// 注意事项：
//   - 如果app不是*http.ServeMux类型，会返回错误
//   - 该方法必须在注册路由之前调用
func (nh *NetHTTP) SetApp(app interface{}) error {
	var (
		eng *http.ServeMux
		ok  bool
	)
	if eng, ok = app.(*http.ServeMux); !ok {
		return errors.New("net/http adapter SetApp: wrong parameter")
	}
	nh.app = eng
	return nil
}

// AddHandler实现Adapter.AddHandler接口方法
// 添加路由处理器到ServeMux中
//
// 参数说明：
//   - method: HTTP方法，如GET、POST、PUT、DELETE等
//   - path: 路由路径，支持路径参数（如:/user/:id）
//   - handlers: 处理器链，按顺序执行
//
// 工作原理：
//  1. 将路径参数从:param格式转换为{param}格式（标准库规范）
//  2. 移除路径开头的双斜杠
//  3. 创建方法+路径的模式字符串（如"GET /users/{id}"）
//  4. 注册处理器到ServeMux
//  5. 在处理器中提取路径参数并注入到查询参数中
//
// 路径参数处理：
//   - 原始路径: /users/:id/posts/:post_id
//   - 转换后: /users/{id}/posts/{post_id}
//   - 模式: GET /users/{id}/posts/{post_id}
//
// 使用示例：
//
//	AddHandler("GET", "/users/:id", handlers)
//	AddHandler("POST", "/users", handlers)
func (nh *NetHTTP) AddHandler(method, path string, handlers context.Handlers) {
	// 步骤1: 将路径参数从:param格式转换为{param}格式以符合标准库规范
	// 例如: /users/:id/posts/:post_id -> /users/{id}/posts/{post_id}
	url := path
	reg1 := regexp.MustCompile(":(.*?)/")     // 匹配路径中间的参数 (如 :id/)
	reg2 := regexp.MustCompile(":(.*?)$")     // 匹配路径末尾的参数 (如 :id)
	url = reg1.ReplaceAllString(url, "{$1}/") // 替换中间参数为{param}/
	url = reg2.ReplaceAllString(url, "{$1}")  // 替换末尾参数为{param}

	// 步骤2: 移除路径开头的双斜杠，避免路由匹配问题
	if len(url) > 1 && url[0] == '/' && url[1] == '/' {
		url = url[1:]
	}

	// 步骤3: 创建标准库ServeMux所需的路由模式字符串
	// 格式: "METHOD /path/{param}"
	pattern := fmt.Sprintf("%s %s", strings.ToUpper(method), url)

	// 步骤4: 注册处理器到ServeMux
	nh.app.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		// 步骤4.1: 移除URL路径末尾的斜杠，确保路径匹配一致性
		if r.URL.Path[len(r.URL.Path)-1] == '/' {
			r.URL.Path = r.URL.Path[:len(r.URL.Path)-1]
		}

		// 步骤4.2: 创建框架上下文对象，封装HTTP请求
		ctx := context.NewContext(r)

		// 步骤4.3: 从URL路径中提取路径参数并注入到查询参数中
		// 这样框架可以通过Query()方法统一获取参数
		params := getPathParams(url, r.URL.Path)
		for key, value := range params {
			if r.URL.RawQuery == "" {
				// 第一个参数，直接添加
				r.URL.RawQuery += strings.ReplaceAll(key, ":", "") + "=" + value
			} else {
				// 后续参数，使用&连接
				r.URL.RawQuery += "&" + strings.ReplaceAll(key, ":", "") + "=" + value
			}
		}

		// 步骤4.4: 执行处理器链
		ctx.SetHandlers(handlers).Next()

		// 步骤4.5: 将框架响应头复制到HTTP响应头
		for key, head := range ctx.Response.Header {
			w.Header().Set(key, head[0])
		}

		// 步骤4.6: 写入响应体
		if ctx.Response.Body != nil {
			// 有响应体: 读取并写入
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(ctx.Response.Body)
			w.WriteHeader(ctx.Response.StatusCode)
			_, _ = w.Write(buf.Bytes())
		} else {
			// 无响应体: 仅写入状态码
			w.WriteHeader(ctx.Response.StatusCode)
		}
	})
}

// HandleFun是路由方法的函数类型
// 用于在ServeMux中注册路由处理器
//
// 参数说明：
//   - pattern: 路由模式字符串，如"GET /users/{id}"
//   - handlerFn: HTTP处理函数
//
// 使用示例：
//
//	handleFun("GET /users", handlerFunc)
type HandleFun func(pattern string, handlerFn http.HandlerFunc)

// Context结构体封装了HTTP请求和响应对象
// 它作为适配器与框架之间的桥梁，提供对HTTP请求和响应的访问
//
// 字段说明：
//   - Request: HTTP请求对象，包含请求的所有信息
//   - Response: HTTP响应写入器，用于写入响应数据
//
// 使用场景：
//   - 在处理器中访问请求参数、头部、cookie等
//   - 写入响应数据、设置响应头
//   - 处理文件上传、表单提交等
type Context struct {
	Request  *http.Request
	Response http.ResponseWriter
}

// SetContext实现Adapter.SetContext接口方法
// 设置上下文对象到适配器中
//
// 参数说明：
//   - contextInterface: 上下文接口，必须是Context类型
//
// 返回值：
//   - adapter.WebFrameWork: 返回适配器实例
//
// 注意事项：
//   - 如果contextInterface不是Context类型，会触发panic
//   - 该方法用于在请求处理过程中设置上下文
func (*NetHTTP) SetContext(contextInterface interface{}) adapter.WebFrameWork {
	var (
		ctx Context
		ok  bool
	)
	if ctx, ok = contextInterface.(Context); !ok {
		panic("net/http adapter SetContext: wrong parameter")
	}
	return &NetHTTP{ctx: ctx}
}

// Name实现Adapter.Name接口方法
// 返回适配器的名称
//
// 返回值：
//   - string: 适配器名称，固定返回"net/http"
//
// 使用场景：
//   - 识别当前使用的适配器类型
//   - 日志记录和调试
func (*NetHTTP) Name() string {
	return "net/http"
}

// Redirect实现Adapter.Redirect接口方法
// 重定向到登录页面
//
// 工作原理：
//  1. 从配置中获取登录URL
//  2. 使用HTTP 302状态码进行重定向
//  3. 重定向到登录页面
//
// 使用场景：
//   - 用户未登录时重定向到登录页
//   - 会话过期时重新认证
func (nh *NetHTTP) Redirect() {
	http.Redirect(nh.ctx.Response, nh.ctx.Request, cfg.Url(cfg.GetLoginUrl()), http.StatusFound)
}

// SetContentType实现Adapter.SetContentType接口方法
// 设置响应的Content-Type头部
//
// 工作原理：
//  1. 调用HTMLContentType()方法获取内容类型
//  2. 设置到响应头中
//
// 注意事项：
//   - 通常设置为"text/html; charset=utf-8"
//   - 必须在写入响应体之前调用
func (nh *NetHTTP) SetContentType() {
	nh.ctx.Response.Header().Set("Content-Type", nh.HTMLContentType())
}

// Write实现Adapter.Write接口方法
// 写入响应体数据
//
// 参数说明：
//   - body: 要写入的字节数据
//
// 工作原理：
//  1. 设置HTTP状态码为200 OK
//  2. 将body数据写入响应
//
// 使用示例：
//
//	Write([]byte("Hello, World!"))
func (nh *NetHTTP) Write(body []byte) {
	nh.ctx.Response.WriteHeader(http.StatusOK)
	_, _ = nh.ctx.Response.Write(body)
}

// GetCookie实现Adapter.GetCookie接口方法
// 从请求中获取指定名称的cookie值
//
// 返回值：
//   - string: cookie的值
//   - error: 获取失败时的错误信息
//
// 使用场景：
//   - 获取会话token
//   - 读取用户偏好设置
//   - 验证用户身份
func (nh *NetHTTP) GetCookie() (string, error) {
	cookie, err := nh.ctx.Request.Cookie(nh.CookieKey())
	if err != nil {
		return "", err
	}
	return cookie.Value, err
}

// Lang实现Adapter.Lang接口方法
// 从查询参数中获取语言设置
//
// 返回值：
//   - string: 语言代码，如"zh-CN"、"en-US"等
//
// 使用场景：
//   - 国际化支持
//   - 根据用户语言显示不同内容
//
// 注意事项：
//   - 通过查询参数__ga_lang传递
//   - 示例: /admin?__ga_lang=zh-CN
func (nh *NetHTTP) Lang() string {
	return nh.ctx.Request.URL.Query().Get("__ga_lang")
}

// Path实现Adapter.Path接口方法
// 获取请求的路径部分
//
// 返回值：
//   - string: 请求路径，如"/admin/users"
//
// 使用场景：
//   - 路由匹配
//   - 权限验证
//   - 日志记录
func (nh *NetHTTP) Path() string {
	return nh.ctx.Request.URL.Path
}

// Method实现Adapter.Method接口方法
// 获取HTTP请求方法
//
// 返回值：
//   - string: HTTP方法，如GET、POST、PUT、DELETE等
//
// 使用场景：
//   - 区分不同类型的请求
//   - 权限控制（如只允许GET请求）
//   - 日志记录和监控
func (nh *NetHTTP) Method() string {
	return nh.ctx.Request.Method
}

// FormParam实现Adapter.FormParam接口方法
// 获取表单参数
//
// 返回值：
//   - url.Values: 表单参数的键值对集合
//
// 工作原理：
//  1. 解析multipart表单数据（最大32MB）
//  2. 返回PostForm字段中的表单参数
//
// 使用场景：
//   - 处理表单提交
//   - 文件上传
//   - POST请求参数获取
//
// 注意事项：
//   - 支持application/x-www-form-urlencoded和multipart/form-data
//   - 自动解析表单数据
func (nh *NetHTTP) FormParam() url.Values {
	_ = nh.ctx.Request.ParseMultipartForm(32 << 20)
	return nh.ctx.Request.PostForm
}

// IsPjax实现Adapter.IsPjax接口方法
// 判断是否为PJAX请求
//
// 返回值：
//   - bool: true表示是PJAX请求，false表示不是
//
// PJAX说明：
//
//	PJAX（pushState + AJAX）是一种技术，允许在不刷新整个页面的情况下更新部分内容
//	通过X-PJAX头部标识
//
// 使用场景：
//   - 实现无刷新页面更新
//   - 提升用户体验
//   - 减少服务器负载
func (nh *NetHTTP) IsPjax() bool {
	return nh.ctx.Request.Header.Get(constant.PjaxHeader) == "true"
}

// Query实现Adapter.Query接口方法
// 获取URL查询参数
//
// 返回值：
//   - url.Values: 查询参数的键值对集合
//
// 使用场景：
//   - 获取GET请求参数
//   - 分页参数（page、limit）
//   - 搜索和过滤参数
//
// 使用示例：
//
//	query := Query()
//	page := query.Get("page")
//	keyword := query.Get("keyword")
func (nh *NetHTTP) Query() url.Values {
	return nh.ctx.Request.URL.Query()
}

// Request实现Adapter.Request接口方法
// 获取原始HTTP请求对象
//
// 返回值：
//   - *http.Request: 原始HTTP请求对象
//
// 使用场景：
//   - 访问请求的所有属性
//   - 获取请求头、Cookie等
//   - 传递给其他需要请求对象的函数
func (nh *NetHTTP) Request() *http.Request {
	return nh.ctx.Request
}

// getPathParams根据给定的路由模式从URL中提取路径参数
//
// 参数说明：
//   - pattern: 路由模式字符串，如"/users/{id}/posts/{post_id}"
//   - url: 实际请求的URL路径，如"/users/123/posts/456"
//
// 返回值：
//   - map[string]string: 参数名到参数值的映射，如{"id": "123", "post_id": "456"}
//
// 工作原理：
//  1. 将路由模式中的{param}转换为正则表达式命名捕获组
//  2. 编译正则表达式
//  3. 使用正则表达式匹配URL路径
//  4. 提取命名捕获组的值
//
// 转换示例：
//
//	输入模式: "/users/{id}/posts/{post_id}"
//	正则表达式: "^/users/(?P<id>\w+)/posts/(?P<post_id>\w+)$"
//	输入URL: "/users/123/posts/456"
//	提取结果: {"id": "123", "post_id": "456"}
//
// 正则表达式说明：
//   - ^: 匹配字符串开头
//   - (?P<name>\w+): 命名捕获组，匹配一个或多个单词字符
//   - $: 匹配字符串结尾
//
// 注意事项：
//   - 参数名只能包含字母、数字和下划线
//   - 参数值只能包含单词字符（字母、数字、下划线）
//   - 如果URL不匹配模式，返回nil
//
// 使用示例：
//
//	params := getPathParams("/users/{id}", "/users/123")
//	id := params["id"]  // id = "123"
func getPathParams(pattern, url string) map[string]string {
	params := make(map[string]string)

	placeholderRegex := regexp.MustCompile(`\{(\w+)\}`)
	regexPattern := "^" + placeholderRegex.ReplaceAllStringFunc(pattern, func(s string) string {
		return `(?P<` + s[1:len(s)-1] + `>\w+)`
	}) + "$"

	regex := regexp.MustCompile(regexPattern)

	match := regex.FindStringSubmatch(url)
	if match == nil {
		return nil
	}

	for i, name := range regex.SubexpNames() {
		if i != 0 && name != "" {
			params[name] = match[i]
		}
	}

	return params
}

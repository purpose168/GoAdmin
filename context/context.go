// Copyright 2019 GoAdmin Core Team. All rights reserved.
// Use of this source code is governed by a Apache-2.0 style
// license that can be found in the LICENSE file.

// context包提供了HTTP请求和响应的上下文管理
// 该包是GoAdmin框架的核心组件之一，提供了轻量级的请求/响应处理机制
//
// 主要功能：
//   - Context结构体：封装HTTP请求和响应，提供便捷的访问方法
//   - App结构体：路由管理器，支持路由注册和中间件
//   - RouterGroup结构体：路由分组，支持前缀和中间件
//   - 路径处理：提供路径标准化和连接功能
//
// 设计理念：
//   - 简化Web框架上下文，提供框架无关的接口
//   - 支持中间件模式，实现请求处理链
//   - 提供丰富的辅助方法，简化常见操作
//   - 支持路由参数和通配符
//
// 核心组件：
//   - Context：请求上下文，包含Request、Response和UserValue
//   - App：应用实例，管理路由和处理器
//   - RouterGroup：路由分组，支持嵌套和前缀
//   - Handler：处理器函数类型
//   - Handlers：处理器链类型
//
// 使用场景：
//   - 插件开发：在插件中使用Context处理请求和响应
//   - 适配器开发：将Web框架的上下文转换为Context
//   - 路由管理：使用App和RouterGroup注册路由
//   - 中间件开发：实现请求拦截和处理
//
// 注意事项：
//   - Context是轻量级的，不包含Web框架特定的功能
//   - 适配器负责将Context转换为Web框架的上下文
//   - 中间件必须调用Next()才能继续处理链
//   - 路由参数使用:__前缀标识
//
// 作者: GoAdmin Core Team
// 创建日期: 2019-01-01
// 版本: 1.0.0
package context

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/purpose168/GoAdmin/modules/constant"
)

const abortIndex int8 = math.MaxInt8 / 2

// Context结构体是Web框架上下文的简化版本
// 它是GoAdmin框架的核心组件，用于在插件中自定义请求和响应处理
// 适配器负责将Context转换为Web框架的上下文
//
// 字段说明：
//   - Request: HTTP请求对象，来自net/http包
//   - Response: HTTP响应对象，来自net/http包
//   - UserValue: 用户自定义的键值对存储，用于在处理器链中传递数据
//   - index: 当前处理器在处理器链中的索引，用于中间件控制
//   - handlers: 处理器链，包含所有要执行的处理器
//
// 使用场景：
//   - 在插件中访问请求和响应
//   - 在中间件中控制请求处理流程
//   - 在处理器链中传递数据
//
// 注意事项：
//   - Request和Response属于net/http包，不是Web框架特定的类型
//   - UserValue是线程安全的，每个请求都有独立的副本
//   - index用于Abort()和Next()方法控制处理流程
type Context struct {
	Request   *http.Request
	Response  *http.Response
	UserValue map[string]interface{}
	index     int8
	handlers  Handlers
}

// Path结构体用于请求和响应的匹配
//
// 字段说明：
//   - URL: 原始注册的URL路径
//   - Method: HTTP方法（GET、POST、PUT、DELETE等）
//
// 使用场景：
//   - 作为路由表的键
//   - 匹配请求路径和方法
//   - 存储路由信息
type Path struct {
	URL    string
	Method string
}

// RouterMap是路由器的映射类型
// 键为路由名称，值为Router对象
type RouterMap map[string]Router

// Get从RouterMap中获取指定名称的路由器
//
// 参数说明：
//   - name: 路由器名称
//
// 返回值：
//   - Router: 路由器对象，如果不存在则返回零值
func (r RouterMap) Get(name string) Router {
	return r[name]
}

// Router结构体表示一个路由器
//
// 字段说明：
//   - Methods: 支持的HTTP方法列表
//   - Patten: 路由模式，可能包含参数占位符
//
// 使用场景：
//   - 定义路由的HTTP方法
//   - 存储路由模式
//   - 生成URL
type Router struct {
	Methods []string
	Patten  string
}

// Method返回路由器的第一个HTTP方法
//
// 返回值：
//   - string: 第一个HTTP方法
//
// 使用场景：
//   - 获取路由的主要方法
//   - 用于路由匹配
func (r Router) Method() string {
	return r.Methods[0]
}

// GetURL根据给定的参数值生成完整的URL
//
// 参数说明：
//   - value: 参数键值对，格式为[key1, value1, key2, value2, ...]
//
// 返回值：
//   - string: 替换参数后的完整URL
//
// 工作原理：
//   - 遍历参数对
//   - 将路由模式中的:__key替换为对应的value
//
// 使用示例：
//
//	router := Router{Patten: "/user/:__id/info/:__type"}
//	url := router.GetURL("id", "123", "type", "detail")
//	// url = "/user/123/info/detail"
func (r Router) GetURL(value ...string) string {
	u := r.Patten
	for i := 0; i < len(value); i += 2 {
		u = strings.ReplaceAll(u, ":__"+value[i], value[i+1])
	}
	return u
}

// NodeProcessor是节点处理器函数类型
// 用于处理面板中的节点
type NodeProcessor func(...Node)

// Node结构体表示一个路由节点
//
// 字段说明：
//   - Path: 路由路径
//   - Method: HTTP方法
//   - Handlers: 处理器链
//   - Value: 节点的自定义值
type Node struct {
	Path     string
	Method   string
	Handlers []Handler
	Value    map[string]interface{}
}

// SetUserValue设置用户上下文的值
//
// 参数说明：
//   - key: 键名
//   - value: 键值
//
// 使用场景：
//   - 在中间件中设置数据供后续处理器使用
//   - 在处理器链中传递数据
//   - 存储请求级别的数据
//
// 使用示例：
//
//	ctx.SetUserValue("userID", "123")
//	userID := ctx.GetUserValue("userID")
func (ctx *Context) SetUserValue(key string, value interface{}) {
	ctx.UserValue[key] = value
}

// GetUserValue获取指定键的值
//
// 参数说明：
//   - key: 键名
//
// 返回值：
//   - interface{}: 键对应的值，如果不存在则返回nil
//
// 使用场景：
//   - 从中间件中获取之前设置的数据
//   - 在处理器链中传递数据
func (ctx *Context) GetUserValue(key string) interface{} {
	return ctx.UserValue[key]
}

// Path返回请求的URL路径
//
// 返回值：
//   - string: URL路径，如"/admin/users"
//
// 使用场景：
//   - 路由匹配
//   - 权限验证
//   - 日志记录
func (ctx *Context) Path() string {
	return ctx.Request.URL.Path
}

// Abort中止上下文的处理
//
// 工作原理：
//   - 将index设置为abortIndex（最大int8值的一半）
//   - 这样Next()方法会立即停止执行后续处理器
//
// 使用场景：
//   - 在中间件中拒绝请求
//   - 权限验证失败时停止处理
//   - 错误处理时提前返回
//
// 使用示例：
//
//	if !hasPermission {
//	    ctx.Abort()
//	    return
//	}
func (ctx *Context) Abort() {
	ctx.index = abortIndex
}

// Next应该在中间件内部使用
//
// 工作原理：
//   - 递增index
//   - 执行下一个处理器
//   - 直到所有处理器执行完毕或被Abort()
//
// 使用场景：
//   - 在中间件中传递控制权给下一个处理器
//   - 实现请求处理链
//
// 注意事项：
//   - 必须在中间件中调用，否则请求处理会停止
//   - 调用Next()后，后续中间件会在Next()返回后继续执行
//
// 使用示例：
//
//	func middleware(ctx *context.Context) {
//	    ctx.SetUserValue("start", time.Now())
//	    ctx.Next()
//	    duration := time.Since(ctx.GetUserValue("start").(time.Time))
//	}
func (ctx *Context) Next() {
	ctx.index++
	for s := int8(len(ctx.handlers)); ctx.index < s; ctx.index++ {
		ctx.handlers[ctx.index](ctx)
	}
}

// SetHandlers设置Context的处理器链
//
// 参数说明：
//   - handlers: 处理器链
//
// 返回值：
//   - *Context: 返回Context本身，支持链式调用
//
// 使用场景：
//   - 在路由匹配后设置处理器链
//   - 初始化请求处理
func (ctx *Context) SetHandlers(handlers Handlers) *Context {
	ctx.handlers = handlers
	return ctx
}

// Method返回请求的HTTP方法
//
// 返回值：
//   - string: HTTP方法，如GET、POST、PUT、DELETE等
//
// 使用场景：
//   - 区分不同类型的请求
//   - 权限控制
//   - 日志记录
func (ctx *Context) Method() string {
	return ctx.Request.Method
}

// NewContext在适配器中使用，返回一个包含请求、UserValue和默认Response的Context
//
// 参数说明：
//   - req: HTTP请求对象
//
// 返回值：
//   - *Context: 新创建的Context对象
//
// 工作原理：
//   - 创建Context结构体
//   - 初始化UserValue为空map
//   - 初始化Response为默认状态（200 OK）
//   - 初始化index为-1
//
// 使用场景：
//   - 适配器将Web框架的请求转换为Context
//   - 创建新的请求上下文
func NewContext(req *http.Request) *Context {

	return &Context{
		Request:   req,
		UserValue: make(map[string]interface{}),
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
		},
		index: -1,
	}
}

const (
	HeaderContentType = "Content-Type"

	HeaderLastModified    = "Last-Modified"
	HeaderIfModifiedSince = "If-Modified-Since"
	HeaderCacheControl    = "Cache-Control"
	HeaderETag            = "ETag"

	HeaderContentDisposition = "Content-Disposition"
	HeaderContentLength      = "Content-Length"
	HeaderContentEncoding    = "Content-Encoding"

	GzipHeaderValue      = "gzip"
	HeaderAcceptEncoding = "Accept-Encoding"
	HeaderVary           = "Vary"

	ThemeKey = "__ga_theme"
)

func (ctx *Context) BindJSON(data interface{}) error {
	if ctx.Request.Body != nil {
		b, err := io.ReadAll(ctx.Request.Body)
		if err == nil {
			return json.Unmarshal(b, data)
		}
		return err
	}
	return errors.New("empty request body")
}

func (ctx *Context) MustBindJSON(data interface{}) {
	if ctx.Request.Body != nil {
		b, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(b, data)
		if err != nil {
			panic(err)
		}
	}
	panic("empty request body")
}

// Write save the given status code, headers and body string into the response.
func (ctx *Context) Write(code int, header map[string]string, Body string) {
	ctx.Response.StatusCode = code
	for key, head := range header {
		ctx.AddHeader(key, head)
	}
	ctx.Response.Body = io.NopCloser(strings.NewReader(Body))
}

// JSON serializes the given struct as JSON into the response body.
// It also sets the Content-Type as "application/json".
func (ctx *Context) JSON(code int, Body map[string]interface{}) {
	ctx.Response.StatusCode = code
	ctx.SetContentType("application/json")
	BodyStr, err := json.Marshal(Body)
	if err != nil {
		panic(err)
	}
	ctx.Response.Body = io.NopCloser(bytes.NewReader(BodyStr))
}

// DataWithHeaders save the given status code, headers and body data into the response.
func (ctx *Context) DataWithHeaders(code int, header map[string]string, data []byte) {
	ctx.Response.StatusCode = code
	for key, head := range header {
		ctx.AddHeader(key, head)
	}
	ctx.Response.Body = io.NopCloser(bytes.NewBuffer(data))
}

// Data writes some data into the body stream and updates the HTTP code.
func (ctx *Context) Data(code int, contentType string, data []byte) {
	ctx.Response.StatusCode = code
	ctx.SetContentType(contentType)
	ctx.Response.Body = io.NopCloser(bytes.NewBuffer(data))
}

// Redirect add redirect url to header.
func (ctx *Context) Redirect(path string) {
	ctx.Response.StatusCode = http.StatusFound
	ctx.SetContentType("text/html; charset=utf-8")
	ctx.AddHeader("Location", path)
}

// HTML output html response.
func (ctx *Context) HTML(code int, body string) {
	ctx.SetContentType("text/html; charset=utf-8")
	ctx.SetStatusCode(code)
	ctx.WriteString(body)
}

// HTMLByte output html response.
func (ctx *Context) HTMLByte(code int, body []byte) {
	ctx.SetContentType("text/html; charset=utf-8")
	ctx.SetStatusCode(code)
	ctx.Response.Body = io.NopCloser(bytes.NewBuffer(body))
}

// WriteString save the given body string into the response.
func (ctx *Context) WriteString(body string) {
	ctx.Response.Body = io.NopCloser(strings.NewReader(body))
}

// SetStatusCode save the given status code into the response.
func (ctx *Context) SetStatusCode(code int) {
	ctx.Response.StatusCode = code
}

// SetContentType save the given content type header into the response header.
func (ctx *Context) SetContentType(contentType string) {
	ctx.AddHeader(HeaderContentType, contentType)
}

func (ctx *Context) SetLastModified(modtime time.Time) {
	if !IsZeroTime(modtime) {
		ctx.AddHeader(HeaderLastModified, modtime.UTC().Format(http.TimeFormat)) // or modtime.UTC()?
	}
}

var unixEpochTime = time.Unix(0, 0)

// IsZeroTime reports whether t is obviously unspecified (either zero or Unix()=0).
func IsZeroTime(t time.Time) bool {
	return t.IsZero() || t.Equal(unixEpochTime)
}

// ParseTime parses a time header (such as the Date: header),
// trying each forth formats
// that are allowed by HTTP/1.1:
// time.RFC850, and time.ANSIC.
var ParseTime = func(text string) (t time.Time, err error) {
	t, err = time.Parse(http.TimeFormat, text)
	if err != nil {
		return http.ParseTime(text)
	}

	return
}

func (ctx *Context) WriteNotModified() {
	// RFC 7232 section 4.1:
	// a sender SHOULD NOT generate representation metadata other than the
	// above listed fields unless said metadata exists for the purpose of
	// guiding cache updates (e.g.," Last-Modified" might be useful if the
	// response does not have an ETag field).
	delete(ctx.Response.Header, HeaderContentType)
	delete(ctx.Response.Header, HeaderContentLength)
	if ctx.Headers(HeaderETag) != "" {
		delete(ctx.Response.Header, HeaderLastModified)
	}
	ctx.SetStatusCode(http.StatusNotModified)
}

func (ctx *Context) CheckIfModifiedSince(modtime time.Time) (bool, error) {
	if method := ctx.Method(); method != http.MethodGet && method != http.MethodHead {
		return false, errors.New("skip: method")
	}
	ims := ctx.Headers(HeaderIfModifiedSince)
	if ims == "" || IsZeroTime(modtime) {
		return false, errors.New("skip: zero time")
	}
	t, err := ParseTime(ims)
	if err != nil {
		return false, errors.New("skip: " + err.Error())
	}
	// sub-second precision, so
	// use mtime < t+1s instead of mtime <= t to check for unmodified.
	if modtime.UTC().Before(t.Add(1 * time.Second)) {
		return false, nil
	}
	return true, nil
}

// LocalIP return the request client ip.
func (ctx *Context) LocalIP() string {
	xForwardedFor := ctx.Request.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(ctx.Request.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(ctx.Request.RemoteAddr)); err == nil {
		return ip
	}

	return "127.0.0.1"
}

// SetCookie save the given cookie obj into the response Set-Cookie header.
func (ctx *Context) SetCookie(cookie *http.Cookie) {
	if v := cookie.String(); v != "" {
		ctx.AddHeader("Set-Cookie", v)
	}
}

// Query get the query parameter of url.
func (ctx *Context) Query(key string) string {
	return ctx.Request.URL.Query().Get(key)
}

// QueryAll get the query parameters of url.
func (ctx *Context) QueryAll(key string) []string {
	return ctx.Request.URL.Query()[key]
}

// QueryDefault get the query parameter of url. If it is empty, return the default.
func (ctx *Context) QueryDefault(key, def string) string {
	value := ctx.Query(key)
	if value == "" {
		return def
	}
	return value
}

// Lang get the query parameter of url with given key __ga_lang.
func (ctx *Context) Lang() string {
	return ctx.Query("__ga_lang")
}

// Theme get the request theme with given key __ga_theme.
func (ctx *Context) Theme() string {
	queryTheme := ctx.Query(ThemeKey)
	if queryTheme != "" {
		return queryTheme
	}
	cookieTheme := ctx.Cookie(ThemeKey)
	if cookieTheme != "" {
		return cookieTheme
	}
	return ctx.RefererQuery(ThemeKey)
}

// Headers get the value of request headers key.
func (ctx *Context) Headers(key string) string {
	return ctx.Request.Header.Get(key)
}

// Referer get the url string of request header Referer.
func (ctx *Context) Referer() string {
	return ctx.Headers("Referer")
}

// RefererURL get the url.URL object of request header Referer.
func (ctx *Context) RefererURL() *url.URL {
	ref := ctx.Headers("Referer")
	if ref == "" {
		return nil
	}
	u, err := url.Parse(ref)
	if err != nil {
		return nil
	}
	return u
}

// RefererQuery retrieve the value of given key from url.URL object of request header Referer.
func (ctx *Context) RefererQuery(key string) string {
	if u := ctx.RefererURL(); u != nil {
		return u.Query().Get(key)
	}
	return ""
}

// FormValue get the value of request form key.
func (ctx *Context) FormValue(key string) string {
	return ctx.Request.FormValue(key)
}

// PostForm get the values of request form.
func (ctx *Context) PostForm() url.Values {
	_ = ctx.Request.ParseMultipartForm(32 << 20)
	return ctx.Request.PostForm
}

func (ctx *Context) WantHTML() bool {
	return ctx.Method() == "GET" && strings.Contains(ctx.Headers("Accept"), "html")
}

func (ctx *Context) WantJSON() bool {
	return strings.Contains(ctx.Headers("Accept"), "json")
}

// AddHeader adds the key, value pair to the header.
func (ctx *Context) AddHeader(key, value string) {
	ctx.Response.Header.Add(key, value)
}

// PjaxUrl add pjax url header.
func (ctx *Context) PjaxUrl(url string) {
	ctx.Response.Header.Add(constant.PjaxUrlHeader, url)
}

// IsPjax check request is pjax or not.
func (ctx *Context) IsPjax() bool {
	return ctx.Headers(constant.PjaxHeader) == "true"
}

// IsIframe check request is iframe or not.
func (ctx *Context) IsIframe() bool {
	return ctx.Query(constant.IframeKey) == "true" || ctx.Headers(constant.IframeKey) == "true"
}

// SetHeader set the key, value pair to the header.
func (ctx *Context) SetHeader(key, value string) {
	ctx.Response.Header.Set(key, value)
}

func (ctx *Context) GetContentType() string {
	return ctx.Request.Header.Get("Content-Type")
}

func (ctx *Context) Cookie(name string) string {
	for _, ck := range ctx.Request.Cookies() {
		if ck.Name == name {
			return ck.Value
		}
	}
	return ""
}

// User return the current login user.
func (ctx *Context) User() interface{} {
	return ctx.UserValue["user"]
}

// ServeContent serves content, headers are autoset
// receives three parameters, it's low-level function, instead you can use .ServeFile(string,bool)/SendFile(string,string)
//
// You can define your own "Content-Type" header also, after this function call
// Doesn't implements resuming (by range), use ctx.SendFile instead
func (ctx *Context) ServeContent(content io.ReadSeeker, filename string, modtime time.Time, gzipCompression bool) error {
	if modified, err := ctx.CheckIfModifiedSince(modtime); !modified && err == nil {
		ctx.WriteNotModified()
		return nil
	}

	if ctx.GetContentType() == "" {
		ctx.SetContentType(filename)
	}

	buf, _ := io.ReadAll(content)
	ctx.Response.Body = io.NopCloser(bytes.NewBuffer(buf))
	return nil
}

// ServeFile serves a view file, to send a file ( zip for example) to the client you should use the SendFile(serverfilename,clientfilename)
func (ctx *Context) ServeFile(filename string, gzipCompression bool) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("%d", http.StatusNotFound)
	}
	defer func() {
		_ = f.Close()
	}()
	fi, _ := f.Stat()
	if fi.IsDir() {
		return ctx.ServeFile(path.Join(filename, "index.html"), gzipCompression)
	}

	return ctx.ServeContent(f, fi.Name(), fi.ModTime(), gzipCompression)
}

type HandlerMap map[Path]Handlers

// App结构体是包的核心结构体
// App作为插件实体的成员，包含请求和对应的处理器
// Prefix是URL前缀，MiddlewareList用于控制流程
//
// 字段说明：
//   - Requests: 请求路径列表
//   - Handlers: 处理器映射，键为Path，值为处理器链
//   - Middlewares: 中间件列表，用于控制请求流程
//   - Prefix: URL前缀，用于路由分组
//   - Routers: 路由器映射，键为路由名称
//   - routeIndex: 路由索引，用于跟踪当前路由
//   - routeANY: 是否为ANY路由（匹配所有HTTP方法）
//
// 使用场景：
//   - 插件中定义路由
//   - 管理中间件
//   - 分组路由
//   - 生成路由URL
type App struct {
	Requests    []Path
	Handlers    HandlerMap
	Middlewares Handlers
	Prefix      string

	Routers    RouterMap
	routeIndex int
	routeANY   bool
}

// NewApp返回一个空的App实例
//
// 返回值：
//   - *App: 新创建的App实例
//
// 工作原理：
//   - 初始化Requests为空切片
//   - 初始化Handlers为空映射
//   - 设置Prefix为"/"
//   - 初始化Middlewares为空切片
//   - 设置routeIndex为-1
//   - 初始化Routers为空映射
//
// 使用场景：
//   - 创建新的应用实例
//   - 初始化插件路由
func NewApp() *App {
	return &App{
		Requests:    make([]Path, 0),
		Handlers:    make(HandlerMap),
		Prefix:      "/",
		Middlewares: make([]Handler, 0),
		routeIndex:  -1,
		Routers:     make(RouterMap),
	}
}

// Handler定义了中间件使用的处理器函数类型
//
// 参数说明：
//   - ctx: 上下文对象
//
// 使用场景：
//   - 定义路由处理器
//   - 定义中间件
//   - 定义错误处理器
type Handler func(ctx *Context)

// Handlers是Handler的数组类型
//
// 使用场景：
//   - 处理器链
//   - 中间件列表
type Handlers []Handler

// AppendReqAndResp将请求信息和处理器存储到app中
// 支持路由参数。路由参数将被识别为通配符存储到Path结构的RegUrl中
//
// 路由参数示例：
//
//	/user/:id      => /user/(.*)
//	/user/:id/info => /user/(.*?)/info
//
// # RegUrl将用于识别传入的路径并查找处理器
//
// 参数说明：
//   - url: 路由路径
//   - method: HTTP方法
//   - handler: 处理器链
//
// 工作原理：
//   - 将路径添加到Requests列表
//   - 递增routeIndex
//   - 将处理器添加到Handlers映射
//   - 处理器包含所有中间件和最终处理器
//
// 使用场景：
//   - 注册路由
//   - 添加处理器
func (app *App) AppendReqAndResp(url, method string, handler []Handler) {

	app.Requests = append(app.Requests, Path{
		URL:    join(app.Prefix, url),
		Method: method,
	})
	app.routeIndex++

	app.Handlers[Path{
		URL:    join(app.Prefix, url),
		Method: method,
	}] = append(app.Middlewares, handler...)
}

// Find是findPath的公共辅助方法
//
// 参数说明：
//   - url: 请求URL
//   - method: HTTP方法
//
// 返回值：
//   - []Handler: 处理器链，如果不存在则返回nil
//
// 使用场景：
//   - 查找路由处理器
//   - 路由匹配
func (app *App) Find(url, method string) []Handler {
	app.routeANY = false
	return app.Handlers[Path{URL: url, Method: method}]
}

// POST是app.AppendReqAndResp(url, "post", handler)的快捷方法
//
// 参数说明：
//   - url: 路由路径
//   - handler: 处理器链
//
// 返回值：
//   - *App: 返回App本身，支持链式调用
//
// 使用示例：
//
//	app.POST("/users", handler1, handler2)
func (app *App) POST(url string, handler ...Handler) *App {
	app.routeANY = false
	app.AppendReqAndResp(url, "post", handler)
	return app
}

// GET是app.AppendReqAndResp(url, "get", handler)的快捷方法
//
// 参数说明：
//   - url: 路由路径
//   - handler: 处理器链
//
// 返回值：
//   - *App: 返回App本身，支持链式调用
//
// 使用示例：
//
//	app.GET("/users", handler1, handler2)
func (app *App) GET(url string, handler ...Handler) *App {
	app.routeANY = false
	app.AppendReqAndResp(url, "get", handler)
	return app
}

// DELETE是app.AppendReqAndResp(url, "delete", handler)的快捷方法
//
// 参数说明：
//   - url: 路由路径
//   - handler: 处理器链
//
// 返回值：
//   - *App: 返回App本身，支持链式调用
//
// 使用示例：
//
//	app.DELETE("/users/:id", handler)
func (app *App) DELETE(url string, handler ...Handler) *App {
	app.routeANY = false
	app.AppendReqAndResp(url, "delete", handler)
	return app
}

// PUT是app.AppendReqAndResp(url, "put", handler)的快捷方法
//
// 参数说明：
//   - url: 路由路径
//   - handler: 处理器链
//
// 返回值：
//   - *App: 返回App本身，支持链式调用
//
// 使用示例：
//
//	app.PUT("/users/:id", handler)
func (app *App) PUT(url string, handler ...Handler) *App {
	app.routeANY = false
	app.AppendReqAndResp(url, "put", handler)
	return app
}

// OPTIONS是app.AppendReqAndResp(url, "options", handler)的快捷方法
//
// 参数说明：
//   - url: 路由路径
//   - handler: 处理器链
//
// 返回值：
//   - *App: 返回App本身，支持链式调用
//
// 使用示例：
//
//	app.OPTIONS("/users", handler)
func (app *App) OPTIONS(url string, handler ...Handler) *App {
	app.routeANY = false
	app.AppendReqAndResp(url, "options", handler)
	return app
}

// HEAD是app.AppendReqAndResp(url, "head", handler)的快捷方法
//
// 参数说明：
//   - url: 路由路径
//   - handler: 处理器链
//
// 返回值：
//   - *App: 返回App本身，支持链式调用
//
// 使用示例：
//
//	app.HEAD("/users", handler)
func (app *App) HEAD(url string, handler ...Handler) *App {
	app.routeANY = false
	app.AppendReqAndResp(url, "head", handler)
	return app
}

// ANY注册一个匹配所有HTTP方法的路由
// 包括：GET、POST、PUT、HEAD、OPTIONS、DELETE
//
// 参数说明：
//   - url: 路由路径
//   - handler: 处理器链
//
// 返回值：
//   - *App: 返回App本身，支持链式调用
//
// 使用场景：
//   - 需要处理所有HTTP方法的路由
//   - 简化路由注册
//
// 使用示例：
//
//	app.ANY("/api", handler)
func (app *App) ANY(url string, handler ...Handler) *App {
	app.routeANY = true
	app.AppendReqAndResp(url, "post", handler)
	app.AppendReqAndResp(url, "get", handler)
	app.AppendReqAndResp(url, "delete", handler)
	app.AppendReqAndResp(url, "put", handler)
	app.AppendReqAndResp(url, "options", handler)
	app.AppendReqAndResp(url, "head", handler)
	return app
}

// Name为路由命名
//
// 参数说明：
//   - name: 路由名称
//
// 工作原理：
//   - 如果是ANY路由，Methods包含所有HTTP方法
//   - 否则Methods只包含当前路由的方法
//   - Patten为当前路由的URL
//
// 使用场景：
//   - 为路由命名
//   - 生成路由URL
func (app *App) Name(name string) {
	if app.routeANY {
		app.Routers[name] = Router{
			Methods: []string{"POST", "GET", "DELETE", "PUT", "OPTIONS", "HEAD"},
			Patten:  app.Requests[app.routeIndex].URL,
		}
	} else {
		app.Routers[name] = Router{
			Methods: []string{app.Requests[app.routeIndex].Method},
			Patten:  app.Requests[app.routeIndex].URL,
		}
	}
}

// Group为App添加中间件和前缀
//
// 参数说明：
//   - prefix: URL前缀
//   - middleware: 中间件列表
//
// 返回值：
//   - *RouterGroup: 新的路由分组
//
// 工作原理：
//   - 创建新的RouterGroup
//   - 继承App的所有中间件
//   - 添加新的中间件
//   - 设置前缀
//
// 使用场景：
//   - 路由分组
//   - 嵌套路由
//   - 共享中间件
//
// 使用示例：
//
//	api := app.Group("/api", authMiddleware)
//	api.GET("/users", handler)
//	// 注册为 /api/users
func (app *App) Group(prefix string, middleware ...Handler) *RouterGroup {
	return &RouterGroup{
		app:         app,
		Middlewares: append(app.Middlewares, middleware...),
		Prefix:      slash(prefix),
	}
}

// RouterGroup是路由分组结构体
//
// 字段说明：
//   - app: 所属的App实例
//   - Middlewares: 中间件列表
//   - Prefix: URL前缀
//
// 使用场景：
//   - 路由分组
//   - 嵌套路由
//   - 共享中间件
type RouterGroup struct {
	app         *App
	Middlewares Handlers
	Prefix      string
}

// AppendReqAndResp存储请求信息和处理器到app
// 支持路由参数。路由参数将被识别为通配符并存储到Path结构的RegUrl中。例如：
//
//	/user/:id      => /user/(.*)
//	/user/:id/info => /user/(.*?)/info
//
// # RegUrl将用于识别传入的路径并查找处理器
//
// 参数说明：
//   - url: 路由路径
//   - method: HTTP方法
//   - handler: 处理器链
//
// 工作原理：
//   - 将URL和Method添加到app.Requests
//   - 复制RouterGroup的中间件
//   - 将中间件和处理器合并后存储到app.Handlers
func (g *RouterGroup) AppendReqAndResp(url, method string, handler []Handler) {

	g.app.Requests = append(g.app.Requests, Path{
		URL:    join(g.Prefix, url),
		Method: method,
	})
	g.app.routeIndex++

	var h = make([]Handler, len(g.Middlewares))
	copy(h, g.Middlewares)

	g.app.Handlers[Path{
		URL:    join(g.Prefix, url),
		Method: method,
	}] = append(h, handler...)
}

// POST是g.AppendReqAndResp(url, "post", handler)的快捷方法
//
// 参数说明：
//   - url: 路由路径
//   - handler: 处理器链
//
// 返回值：
//   - *RouterGroup: 返回RouterGroup本身，支持链式调用
//
// 使用示例：
//
//	group.POST("/users", handler)
func (g *RouterGroup) POST(url string, handler ...Handler) *RouterGroup {
	g.app.routeANY = false
	g.AppendReqAndResp(url, "post", handler)
	return g
}

// GET是g.AppendReqAndResp(url, "get", handler)的快捷方法
//
// 参数说明：
//   - url: 路由路径
//   - handler: 处理器链
//
// 返回值：
//   - *RouterGroup: 返回RouterGroup本身，支持链式调用
//
// 使用示例：
//
//	group.GET("/users", handler)
func (g *RouterGroup) GET(url string, handler ...Handler) *RouterGroup {
	g.app.routeANY = false
	g.AppendReqAndResp(url, "get", handler)
	return g
}

// DELETE是g.AppendReqAndResp(url, "delete", handler)的快捷方法
//
// 参数说明：
//   - url: 路由路径
//   - handler: 处理器链
//
// 返回值：
//   - *RouterGroup: 返回RouterGroup本身，支持链式调用
//
// 使用示例：
//
//	group.DELETE("/users/:id", handler)
func (g *RouterGroup) DELETE(url string, handler ...Handler) *RouterGroup {
	g.app.routeANY = false
	g.AppendReqAndResp(url, "delete", handler)
	return g
}

// PUT是g.AppendReqAndResp(url, "put", handler)的快捷方法
//
// 参数说明：
//   - url: 路由路径
//   - handler: 处理器链
//
// 返回值：
//   - *RouterGroup: 返回RouterGroup本身，支持链式调用
//
// 使用示例：
//
//	group.PUT("/users/:id", handler)
func (g *RouterGroup) PUT(url string, handler ...Handler) *RouterGroup {
	g.app.routeANY = false
	g.AppendReqAndResp(url, "put", handler)
	return g
}

// OPTIONS是g.AppendReqAndResp(url, "options", handler)的快捷方法
//
// 参数说明：
//   - url: 路由路径
//   - handler: 处理器链
//
// 返回值：
//   - *RouterGroup: 返回RouterGroup本身，支持链式调用
//
// 使用示例：
//
//	group.OPTIONS("/users", handler)
func (g *RouterGroup) OPTIONS(url string, handler ...Handler) *RouterGroup {
	g.app.routeANY = false
	g.AppendReqAndResp(url, "options", handler)
	return g
}

// HEAD是g.AppendReqAndResp(url, "head", handler)的快捷方法
//
// 参数说明：
//   - url: 路由路径
//   - handler: 处理器链
//
// 返回值：
//   - *RouterGroup: 返回RouterGroup本身，支持链式调用
//
// 使用示例：
//
//	group.HEAD("/users", handler)
func (g *RouterGroup) HEAD(url string, handler ...Handler) *RouterGroup {
	g.app.routeANY = false
	g.AppendReqAndResp(url, "head", handler)
	return g
}

// ANY注册一个匹配所有HTTP方法的路由
// 包括：GET、POST、PUT、HEAD、OPTIONS、DELETE
//
// 参数说明：
//   - url: 路由路径
//   - handler: 处理器链
//
// 返回值：
//   - *RouterGroup: 返回RouterGroup本身，支持链式调用
//
// 使用场景：
//   - 需要处理所有HTTP方法的路由
//   - 简化路由注册
//
// 使用示例：
//
//	group.ANY("/api", handler)
func (g *RouterGroup) ANY(url string, handler ...Handler) *RouterGroup {
	g.app.routeANY = true
	g.AppendReqAndResp(url, "post", handler)
	g.AppendReqAndResp(url, "get", handler)
	g.AppendReqAndResp(url, "delete", handler)
	g.AppendReqAndResp(url, "put", handler)
	g.AppendReqAndResp(url, "options", handler)
	g.AppendReqAndResp(url, "head", handler)
	return g
}

// Name为路由命名
//
// 参数说明：
//   - name: 路由名称
//
// 工作原理：
//   - 调用app.Name方法
//   - 将路由名称与当前路由关联
//
// 使用场景：
//   - 为路由命名
//   - 生成路由URL
func (g *RouterGroup) Name(name string) {
	g.app.Name(name)
}

// Group为RouterGroup添加中间件和前缀
//
// 参数说明：
//   - prefix: URL前缀
//   - middleware: 中间件列表
//
// 返回值：
//   - *RouterGroup: 新的RouterGroup
//
// 工作原理：
//   - 创建新的RouterGroup
//   - 继承当前RouterGroup的所有中间件
//   - 添加新的中间件
//   - 连接当前前缀和新前缀
//
// 使用场景：
//   - 嵌套路由分组
//   - 共享中间件
//   - 多级路由组织
//
// 使用示例：
//
//	v1 := api.Group("/v1", middleware1)
//	v2 := v1.Group("/v2", middleware2)
//	v2.GET("/users", handler)
//	// 注册为 /api/v1/v2/users
func (g *RouterGroup) Group(prefix string, middleware ...Handler) *RouterGroup {
	return &RouterGroup{
		app:         g.app,
		Middlewares: append(g.Middlewares, middleware...),
		Prefix:      join(slash(g.Prefix), slash(prefix)),
	}
}

// slash修复格式错误的路径
//
// 参数说明：
//   - prefix: 路径前缀
//
// 返回值：
//   - string: 规范化后的路径
//
// 转换规则：
//
//	""      => "/"
//	"abc/"  => "/abc"
//	"/abc/" => "/abc"
//	"/abc"  => "/abc"
//	"/"     => "/"
//
// 工作原理：
//   - 去除首尾空格
//   - 确保以/开头
//   - 去除末尾的/
//   - 特殊处理根路径"/"
func slash(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" || prefix == "/" {
		return "/"
	}
	if prefix[0] != '/' {
		if prefix[len(prefix)-1] == '/' {
			return "/" + prefix[:len(prefix)-1]
		}
		return "/" + prefix
	}
	if prefix[len(prefix)-1] == '/' {
		return prefix[:len(prefix)-1]
	}
	return prefix
}

// join连接两个路径
//
// 参数说明：
//   - prefix: 前缀路径
//   - suffix: 后缀路径
//
// 返回值：
//   - string: 连接后的路径
//
// 工作原理：
//   - 如果prefix是"/"，直接返回suffix
//   - 如果suffix是"/"，直接返回prefix
//   - 否则直接拼接两个字符串
func join(prefix, suffix string) string {
	if prefix == "/" {
		return suffix
	}
	if suffix == "/" {
		return prefix
	}
	return prefix + suffix
}

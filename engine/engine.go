// 版权所有 2019 GoAdmin 核心团队。保留所有权利。
// 本源代码的使用受 Apache-2.0 风格许可证管辖
// 该许可证可在 LICENSE 文件中找到。

// 包 engine 提供 GoAdmin 框架的核心引擎功能
//
// Engine 是 GoAdmin 的核心组件，负责协调插件、适配器和各种服务。
// 它提供了配置管理、数据库连接、路由注册、页面渲染等核心功能。
//
// 主要组件：
//   - Engine 结构体：核心引擎，管理插件列表、适配器、服务和导航按钮
//   - 插件系统：支持动态添加和管理插件
//   - 适配器系统：连接不同 Web 框架与 GoAdmin
//   - 服务系统：提供认证、配置、UI 等核心服务
//
// 核心功能：
//   - 配置管理：支持从 JSON、YAML、INI 文件加载配置
//   - 数据库连接：支持 MySQL、PostgreSQL、SQLite、MSSQL、OceanBase
//   - 路由注册：提供 HTML、Data、HTMLFile 等路由注册方法
//   - 页面渲染：支持面板组件的渲染和展示
//   - 中间件支持：提供认证、错误处理、日志记录等中间件
//
// 使用示例：
//
//	// 创建引擎实例
//	eng := engine.Default()
//
//	// 添加配置
//	eng.AddConfig(&config.Config{
//	    Databases: config.DatabaseList{
//	        Default: "default",
//	        ...
//	    },
//	})
//
//	// 添加插件
//	eng.AddPlugins(admin.NewAdmin())
//
//	// 启用适配器
//	eng.Use(router)
//
// 设计理念：
//   - 插件化架构：通过插件系统扩展功能
//   - 适配器模式：支持多种 Web 框架
//   - 服务化设计：核心功能通过服务提供
//   - 中间件链：支持灵活的请求处理流程
package engine

import (
	"bytes"
	"encoding/json"
	errors2 "errors"
	"fmt"
	template2 "html/template"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/purpose168/GoAdmin/modules/language"
	"github.com/purpose168/GoAdmin/template/icon"
	"github.com/purpose168/GoAdmin/template/types/action"

	"github.com/purpose168/GoAdmin/adapter"
	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/modules/auth"
	"github.com/purpose168/GoAdmin/modules/config"
	"github.com/purpose168/GoAdmin/modules/db"
	"github.com/purpose168/GoAdmin/modules/errors"
	"github.com/purpose168/GoAdmin/modules/logger"
	"github.com/purpose168/GoAdmin/modules/menu"
	"github.com/purpose168/GoAdmin/modules/service"
	"github.com/purpose168/GoAdmin/modules/system"
	"github.com/purpose168/GoAdmin/modules/ui"
	"github.com/purpose168/GoAdmin/plugins"
	"github.com/purpose168/GoAdmin/plugins/admin"
	"github.com/purpose168/GoAdmin/plugins/admin/models"
	"github.com/purpose168/GoAdmin/plugins/admin/modules/response"
	"github.com/purpose168/GoAdmin/plugins/admin/modules/table"
	"github.com/purpose168/GoAdmin/template"
	"github.com/purpose168/GoAdmin/template/types"
)

// Engine 是 GoAdmin 的核心组件
//
// 它包含两个主要属性：
//   - PluginList：插件数组，包含所有注册的插件
//   - Adapter：Web框架上下文和GoAdmin上下文的适配器
//
// 适配器和插件的关系：
//
//	适配器使用插件，插件包含路由器和控制器方法，
//	将这些内容注入到框架实体中，使其正常工作。
//
// 其他属性：
//   - Services：服务列表，提供认证、配置、UI等核心服务
//   - NavButtons：导航栏按钮，显示在页面顶部导航栏
//   - config：配置对象，存储全局配置信息
//   - announceLock：同步锁，用于确保公告只打印一次
type Engine struct {
	PluginList   plugins.Plugins
	Adapter      adapter.WebFrameWork
	Services     service.List
	NavButtons   *types.Buttons
	config       *config.Config
	announceLock sync.Once
}

// Default 返回默认的引擎实例
//
// 创建并初始化一个Engine实例，设置默认适配器、服务列表和导航按钮
//
// 返回值：
//   - *Engine：默认的引擎实例
//
// 实现细节：
//   - 设置默认适配器为defaultAdapter
//   - 获取服务列表
//   - 初始化导航按钮
func Default() *Engine {
	engine = &Engine{
		Adapter:    defaultAdapter,
		Services:   service.GetServices(),
		NavButtons: new(types.Buttons),
	}
	return engine
}

// Use 启用适配器
//
// 参数：
//   - router：Web框架的路由器实例
//
// 返回值：
//   - error：如果适配器为空则返回错误，否则返回nil
//
// 工作流程：
//  1. 检查适配器是否为空，为空则触发panic
//  2. 添加CSRF令牌服务
//  3. 初始化站点设置
//  4. 初始化跳转导航按钮
//  5. 初始化插件
//  6. 打印初始化成功消息
//  7. 调用适配器的Use方法，将插件列表注入到框架中
func (eng *Engine) Use(router interface{}) error {
	if eng.Adapter == nil {
		emptyAdapterPanic()
	}

	eng.Services.Add(auth.InitCSRFTokenSrv(eng.DefaultConnection()))
	eng.initSiteSetting()
	eng.initJumpNavButtons()
	eng.initPlugins()

	printInitMsg(language.Get("initialize success"))

	return eng.Adapter.Use(router, eng.PluginList)
}

// AddPlugins 添加插件
//
// 参数：
//   - plugs：要添加的插件列表
//
// 返回值：
//   - *Engine：引擎实例，支持链式调用
//
// 工作原理：
//
//	遍历所有传入的插件，将它们添加到引擎的PluginList中
//	如果传入的插件列表为空，则直接返回引擎实例
func (eng *Engine) AddPlugins(plugs ...plugins.Plugin) *Engine {

	if len(plugs) == 0 {
		return eng
	}

	for _, plug := range plugs {
		eng.PluginList = eng.PluginList.Add(plug)
	}

	return eng
}

// AddPluginList 添加插件列表
//
// 参数：
//   - plugs：要添加的插件列表
//
// 返回值：
//   - *Engine：引擎实例，支持链式调用
//
// 工作原理：
//
//	遍历所有传入的插件，将它们添加到引擎的PluginList中
//	如果传入的插件列表为空，则直接返回引擎实例
func (eng *Engine) AddPluginList(plugs plugins.Plugins) *Engine {

	if len(plugs) == 0 {
		return eng
	}

	for _, plug := range plugs {
		eng.PluginList = eng.PluginList.Add(plug)
	}

	return eng
}

// FindPluginByName 根据名称查找注册的插件
//
// 参数：
//   - name：要查找的插件名称
//
// 返回值：
//   - plugins.Plugin：找到的插件实例
//   - bool：是否找到插件
//
// 工作原理：
//
//	遍历所有注册的插件，比较插件名称与传入的名称
//	如果找到匹配的插件，则返回该插件和true，否则返回nil和false
func (eng *Engine) FindPluginByName(name string) (plugins.Plugin, bool) {
	for _, plug := range eng.PluginList {
		if plug.Name() == name {
			return plug, true
		}
	}
	return nil, false
}

// AddAuthService 使用给定的回调函数自定义认证逻辑
//
// 参数：
//   - processor：认证处理器函数
//
// 返回值：
//   - *Engine：引擎实例，支持链式调用
//
// 工作原理：
//
//	创建一个新的认证服务，并将其添加到引擎的服务列表中
//	服务名称为"auth"
func (eng *Engine) AddAuthService(processor auth.Processor) *Engine {
	eng.Services.Add("auth", auth.NewService(processor))
	return eng
}

// ============================
// 配置 API
// ============================

// announce 打印初始化公告
//
// 工作原理：
//   - 如果配置的Debug模式为true，只打印一次公告
//   - 使用sync.Once确保公告只打印一次
//   - 打印"goadmin is now running"等信息
//
// 返回值：
//   - *Engine：引擎实例，支持链式调用
func (eng *Engine) announce() *Engine {
	if eng.config.Debug {
		eng.announceLock.Do(func() {
			fmt.Printf(language.Get("goadmin is now running. \nrunning in \"debug\" mode. switch to \"release\" mode in production.\n\n"))
		})
	}
	return eng
}

// AddConfig 设置全局配置
//
// 参数：
//   - cfg：配置对象指针
//
// 返回值：
//   - *Engine：引擎实例，支持链式调用
//
// 工作流程：
//  1. 设置配置
//  2. 打印初始化公告
//  3. 初始化数据库连接
func (eng *Engine) AddConfig(cfg *config.Config) *Engine {
	return eng.setConfig(cfg).announce().initDatabase()
}

// setConfig 设置引擎的配置
//
// 参数：
//   - cfg：配置对象指针
//
// 返回值：
//   - *Engine：引擎实例，支持链式调用
//
// 工作原理：
//  1. 初始化配置
//  2. 检查主题和GoAdmin版本的兼容性
//  3. 如果版本不兼容，触发panic
//
// 兼容性检查：
//   - 检查GoAdmin版本是否与主题兼容
//   - 检查主题版本是否与GoAdmin兼容
func (eng *Engine) setConfig(cfg *config.Config) *Engine {
	eng.config = config.Initialize(cfg)
	sysCheck, themeCheck := template.CheckRequirements()
	if !sysCheck {
		logger.Panicf(language.Get("wrong goadmin version, theme %s required goadmin version are %s"),
			eng.config.Theme, strings.Join(template.Default().GetRequirements(), ","))
	}
	if !themeCheck {
		logger.Panicf(language.Get("wrong theme version, goadmin %s required version of theme %s is %s"),
			system.Version(), eng.config.Theme, strings.Join(system.RequireThemeVersion()[eng.config.Theme], ","))
	}
	return eng
}

// AddConfigFromJSON 从JSON文件设置全局配置
//
// 参数：
//   - path：JSON文件路径
//
// 返回值：
//   - *Engine：引擎实例，支持链式调用
//
// 工作原理：
//  1. 从JSON文件读取配置
//  2. 设置配置
//  3. 打印初始化公告
//  4. 初始化数据库连接
func (eng *Engine) AddConfigFromJSON(path string) *Engine {
	cfg := config.ReadFromJson(path)
	return eng.setConfig(&cfg).announce().initDatabase()
}

// AddConfigFromYAML 从YAML文件设置全局配置
//
// 参数：
//   - path：YAML文件路径
//
// 返回值：
//   - *Engine：引擎实例，支持链式调用
//
// 工作原理：
//  1. 从YAML文件读取配置
//  2. 设置配置
//  3. 打印初始化公告
//  4. 初始化数据库连接
func (eng *Engine) AddConfigFromYAML(path string) *Engine {
	cfg := config.ReadFromYaml(path)
	return eng.setConfig(&cfg).announce().initDatabase()
}

// AddConfigFromINI 从INI文件设置全局配置
//
// 参数：
//   - path：INI文件路径
//
// 返回值：
//   - *Engine：引擎实例，支持链式调用
//
// 工作原理：
//  1. 从INI文件读取配置
//  2. 设置配置
//  3. 打印初始化公告
//  4. 初始化数据库连接
func (eng *Engine) AddConfigFromINI(path string) *Engine {
	cfg := config.ReadFromINI(path)
	return eng.setConfig(&cfg).announce().initDatabase()
}

// InitDatabase 初始化所有数据库连接
//
// 返回值：
//   - *Engine：引擎实例，支持链式调用
//
// 工作流程：
//  1. 打印初始化数据库连接消息
//  2. 遍历所有数据库配置，按驱动类型分组
//  3. 为每个驱动类型初始化数据库连接
//  4. 检查默认适配器是否为空
//  5. 设置默认适配器的连接
//  6. 设置引擎适配器的连接
func (eng *Engine) initDatabase() *Engine {
	printInitMsg(language.Get("initialize database connections"))
	for driver, databaseCfg := range eng.config.Databases.GroupByDriver() {
		eng.Services.Add(driver, db.GetConnectionByDriver(driver).InitDB(databaseCfg))
	}
	if defaultAdapter == nil {
		emptyAdapterPanic()
	}
	defaultConnection := db.GetConnection(eng.Services)
	defaultAdapter.SetConnection(defaultConnection)
	eng.Adapter.SetConnection(defaultConnection)
	return eng
}

// AddAdapter 添加引擎的适配器
//
// 参数：
//   - ada：Web框架适配器
//
// 返回值：
//   - *Engine：引擎实例，支持链式调用
//
// 工作原理：
//  1. 设置引擎的适配器
//  2. 设置默认适配器
func (eng *Engine) AddAdapter(ada adapter.WebFrameWork) *Engine {
	eng.Adapter = ada
	defaultAdapter = ada
	return eng
}

// defaultAdapter 是引擎的默认适配器
var defaultAdapter adapter.WebFrameWork

var engine *Engine

// navButtons 是导航栏中的默认按钮
var navButtons = new(types.Buttons)

func emptyAdapterPanic() {
	logger.Panic(language.Get("adapter is nil, import the default adapter or use addadapter method add the adapter"))
}

// Register 设置引擎的默认适配器
func Register(ada adapter.WebFrameWork) {
	if ada == nil {
		emptyAdapterPanic()
	}
	defaultAdapter = ada
}

// User 调用 defaultAdapter 的 User 方法
func User(ctx interface{}) (models.UserModel, bool) {
	return defaultAdapter.User(ctx)
}

// User 调用引擎适配器的 User 方法
func (eng *Engine) User(ctx interface{}) (models.UserModel, bool) {
	return eng.Adapter.User(ctx)
}

// ============================
// 数据库连接 API
// ============================

// DB 返回给定驱动的数据库连接
//
// 参数：
//   - driver：数据库驱动名称（mysql、postgresql、sqlite等）
//
// 返回值：
//   - db.Connection：数据库连接实例
//
// 工作原理：
//
//	从服务列表中获取指定驱动的数据库连接服务
//	然后从服务中获取实际的数据库连接
func (eng *Engine) DB(driver string) db.Connection {
	return db.GetConnectionFromService(eng.Services.Get(driver))
}

// DefaultConnection 返回默认数据库连接
//
// 返回值：
//   - db.Connection：默认数据库连接实例
//
// 工作原理：
//  1. 获取配置中默认数据库的驱动名称
//  2. 调用DB函数获取该驱动的数据库连接
func (eng *Engine) DefaultConnection() db.Connection {
	return eng.DB(eng.config.Databases.GetDefault().Driver)
}

// MysqlConnection 返回MySQL数据库连接
//
// 返回值：
//   - db.Connection：MySQL数据库连接实例
//
// 工作原理：
//
//	从服务列表中获取MySQL驱动的数据库连接
func (eng *Engine) MysqlConnection() db.Connection {
	return db.GetConnectionFromService(eng.Services.Get(db.DriverMysql))
}

// MssqlConnection 返回MSSQL数据库连接
//
// 返回值：
//   - db.Connection：MSSQL数据库连接实例
//
// 工作原理：
//
//	从服务列表中获取MSSQL驱动的数据库连接
func (eng *Engine) MssqlConnection() db.Connection {
	return db.GetConnectionFromService(eng.Services.Get(db.DriverMssql))
}

// PostgresqlConnection 返回PostgreSQL数据库连接
//
// 返回值：
//   - db.Connection：PostgreSQL数据库连接实例
//
// 工作原理：
//
//	从服务列表中获取PostgreSQL驱动的数据库连接
func (eng *Engine) PostgresqlConnection() db.Connection {
	return db.GetConnectionFromService(eng.Services.Get(db.DriverPostgresql))
}

// SqliteConnection 返回SQLite数据库连接
//
// 返回值：
//   - db.Connection：SQLite数据库连接实例
//
// 工作原理：
//
//	从服务列表中获取SQLite驱动的数据库连接
func (eng *Engine) SqliteConnection() db.Connection {
	return db.GetConnectionFromService(eng.Services.Get(db.DriverSqlite))
}

// OceanBaseConnection 返回OceanBase数据库连接
//
// 返回值：
//   - db.Connection：OceanBase数据库连接实例
//
// 工作原理：
//
//	从服务列表中获取OceanBase驱动的数据库连接
func (eng *Engine) OceanBaseConnection() db.Connection {
	return db.GetConnectionFromService(eng.Services.Get(db.DriverOceanBase))
}

// ConnectionSetter 是一个用于配置数据库连接的函数类型
//
// 参数：
//   - db.Connection：数据库连接实例
//
// 使用场景：
//
//	用于在获取数据库连接后对连接进行进一步配置
//	可以在连接上设置超时、连接池大小等参数
type ConnectionSetter func(db.Connection)

// ResolveConnection 解析指定驱动的数据库连接
//
// 参数：
//   - setter：连接配置函数
//   - driver：数据库驱动名称
//
// 返回值：
//   - *Engine：引擎实例，支持链式调用
//
// 工作原理：
//  1. 获取指定驱动的数据库连接
//  2. 调用setter函数对连接进行配置
//  3. 返回引擎实例，支持链式调用
func (eng *Engine) ResolveConnection(setter ConnectionSetter, driver string) *Engine {
	setter(eng.DB(driver))
	return eng
}

// ResolveMysqlConnection 解析MySQL数据库连接
//
// 参数：
//   - setter：连接配置函数
//
// 返回值：
//   - *Engine：引擎实例，支持链式调用
//
// 工作原理：
//
//	调用ResolveConnection函数，指定MySQL驱动
func (eng *Engine) ResolveMysqlConnection(setter ConnectionSetter) *Engine {
	eng.ResolveConnection(setter, db.DriverMysql)
	return eng
}

// ResolveMssqlConnection 解析MSSQL数据库连接
//
// 参数：
//   - setter：连接配置函数
//
// 返回值：
//   - *Engine：引擎实例，支持链式调用
//
// 工作原理：
//
//	调用ResolveConnection函数，指定MSSQL驱动
func (eng *Engine) ResolveMssqlConnection(setter ConnectionSetter) *Engine {
	eng.ResolveConnection(setter, db.DriverMssql)
	return eng
}

// ResolveSqliteConnection 解析SQLite数据库连接
//
// 参数：
//   - setter：连接配置函数
//
// 返回值：
//   - *Engine：引擎实例，支持链式调用
//
// 工作原理：
//
//	调用ResolveConnection函数，指定SQLite驱动
func (eng *Engine) ResolveSqliteConnection(setter ConnectionSetter) *Engine {
	eng.ResolveConnection(setter, db.DriverSqlite)
	return eng
}

// ResolvePostgresqlConnection 解析PostgreSQL数据库连接
//
// 参数：
//   - setter：连接配置函数
//
// 返回值：
//   - *Engine：引擎实例，支持链式调用
//
// 工作原理：
//
//	调用ResolveConnection函数，指定PostgreSQL驱动
func (eng *Engine) ResolvePostgresqlConnection(setter ConnectionSetter) *Engine {
	eng.ResolveConnection(setter, db.DriverPostgresql)
	return eng
}

// Setter 是一个用于配置引擎的函数类型
//
// 参数：
//   - *Engine：引擎实例
//
// 使用场景：
//
//	用于在克隆引擎实例时对新实例进行配置
type Setter func(*Engine)

// Clone 复制一个新的引擎实例
//
// 参数：
//   - e：用于接收复制结果的引擎实例指针
//
// 返回值：
//   - *Engine：复制后的引擎实例
//
// 注意：
//
//	该方法目前的实现只是将引擎实例赋值给e，
//	并没有真正创建一个新的副本。
func (eng *Engine) Clone(e *Engine) *Engine {
	e = eng
	return eng
}

// ClonedBySetter 通过设置器函数复制一个新的引擎实例
//
// 参数：
//   - setter：引擎配置函数
//
// 返回值：
//   - *Engine：配置后的引擎实例
//
// 工作原理：
//  1. 调用setter函数对引擎实例进行配置
//  2. 返回配置后的引擎实例
//
// 注意：
//
//	该方法目前的实现只是对原引擎实例进行配置，
//	并没有真正创建一个新的副本。
func (eng *Engine) ClonedBySetter(setter Setter) *Engine {
	setter(eng)
	return eng
}

func (eng *Engine) deferHandler(conn db.Connection) context.Handler {
	return func(ctx *context.Context) {
		defer func(ctx *context.Context) {
			if user, ok := ctx.UserValue["user"].(models.UserModel); ok {
				var input []byte
				form := ctx.Request.MultipartForm
				if form != nil {
					input, _ = json.Marshal((*form).Value)
				}

				models.OperationLog().SetConn(conn).New(user.Id, ctx.Path(), ctx.Method(), ctx.LocalIP(), string(input))
			}

			if err := recover(); err != nil {
				logger.Error(err)
				logger.Error(string(debug.Stack()))

				var (
					errMsg string
					ok     bool
					e      error
				)

				if errMsg, ok = err.(string); !ok {
					if e, ok = err.(error); ok {
						errMsg = e.Error()
					}
				}

				if errMsg == "" {
					errMsg = "系统错误"
				}

				if ctx.WantJSON() {
					response.Error(ctx, errMsg)
					return
				}

				eng.errorPanelHTML(ctx, new(bytes.Buffer), errors2.New(errMsg))
			}
		}(ctx)
		ctx.Next()
	}
}

// wrapWithAuthMiddleware 将认证中间件包装到给定的处理器中
func (eng *Engine) wrapWithAuthMiddleware(handler context.Handler) context.Handlers {
	conn := db.GetConnection(eng.Services)
	return []context.Handler{eng.deferHandler(conn), response.OffLineHandler, auth.Middleware(conn), handler}
}

// wrap 将处理器包装到中间件链中（不包含认证中间件）
func (eng *Engine) wrap(handler context.Handler) context.Handlers {
	conn := db.GetConnection(eng.Services)
	return []context.Handler{eng.deferHandler(conn), response.OffLineHandler, handler}
}

// ============================
// 初始化方法
// ============================

// AddNavButtons 添加导航按钮
func (eng *Engine) AddNavButtons(title template2.HTML, icon string, action types.Action) *Engine {
	btn := types.GetNavButton(title, icon, action)
	*eng.NavButtons = append(*eng.NavButtons, btn)
	return eng
}

// AddNavButtonsRaw 添加导航按钮（原始按钮对象）
func (eng *Engine) AddNavButtonsRaw(btns ...types.Button) *Engine {
	*eng.NavButtons = append(*eng.NavButtons, btns...)
	return eng
}

type navJumpButtonParam struct {
	Exist      bool
	Icon       string
	BtnName    string
	URL        string
	Title      string
	TitleScore string
}

func (eng *Engine) addJumpNavButton(param navJumpButtonParam) *Engine {
	if param.Exist {
		*eng.NavButtons = (*eng.NavButtons).AddNavButton(param.Icon, param.BtnName,
			action.JumpInNewTab(config.Url(param.URL),
				language.GetWithScope(param.Title, param.TitleScore)))
	}
	return eng
}

func printInitMsg(msg string) {
	logger.Info(msg)
}

// initJumpNavButtons 初始化跳转导航按钮
//
// 工作流程：
//   - 打印初始化导航按钮消息
//   - 遍历所有导航按钮参数，添加到导航栏
//   - 将导航按钮设置为全局导航按钮
//   - 将导航按钮服务添加到服务列表中
func (eng *Engine) initJumpNavButtons() {
	printInitMsg(language.Get("initialize navigation buttons"))
	for _, param := range eng.initNavJumpButtonParams() {
		eng.addJumpNavButton(param)
	}
	navButtons = eng.NavButtons
	eng.Services.Add(ui.ServiceKey, ui.NewService(eng.NavButtons))
}

// initPlugins 初始化插件
//
// 工作流程：
//   - 打印初始化插件消息
//   - 添加admin插件和获取的插件列表
//   - 创建插件生成器列表
//   - 遍历所有插件（除了admin插件）
//   - 初始化每个插件
//   - 如果插件需要安装，添加插件生成器
//   - 合并所有插件的生成器
func (eng *Engine) initPlugins() {

	printInitMsg(language.Get("initialize plugins"))

	eng.AddPlugins(admin.NewAdmin()).AddPluginList(plugins.Get())

	var plugGenerators = make(table.GeneratorList)

	for i := range eng.PluginList {
		if eng.PluginList[i].Name() != "admin" {
			printInitMsg("--> " + eng.PluginList[i].Name())
			eng.PluginList[i].InitPlugin(eng.Services)
			if !eng.PluginList[i].GetInfo().SkipInstallation {
				eng.AddGenerator("plugin_"+eng.PluginList[i].Name(), eng.PluginList[i].GetSettingPage())
			}
			plugGenerators = plugGenerators.Combine(eng.PluginList[i].GetGenerators())
		}
	}
	adm := eng.AdminPlugin().AddGenerators(plugGenerators)
	adm.InitPlugin(eng.Services)
	plugins.Add(adm)
}

// initNavJumpButtonParams 初始化导航栏跳转按钮参数
//
// 返回值：
//   - []navJumpButtonParam: 导航栏按钮参数列表
//
// 工作原理：
//   - 返回所有默认导航栏按钮的参数
//   - 根据配置决定是否显示各个按钮
//   - 包含：站点设置、代码生成工具、站点信息、插件管理
//
// 使用场景：
//   - 内部方法，由initJumpNavButtons调用
//   - 定义默认导航栏按钮
func (eng *Engine) initNavJumpButtonParams() []navJumpButtonParam {
	return []navJumpButtonParam{
		{
			Exist:      !eng.config.HideConfigCenterEntrance,
			Icon:       icon.Gear,
			BtnName:    types.NavBtnSiteName,
			URL:        "/info/site/edit",
			Title:      "site setting",
			TitleScore: "config",
		}, {
			Exist:      !eng.config.HideToolEntrance && eng.config.IsNotProductionEnvironment(),
			Icon:       icon.Wrench,
			BtnName:    types.NavBtnToolName,
			URL:        "/info/generate/new",
			Title:      "code generate tool",
			TitleScore: "tool",
		}, {
			Exist:      !eng.config.HideAppInfoEntrance,
			Icon:       icon.Info,
			BtnName:    types.NavBtnInfoName,
			URL:        "/application/info",
			Title:      "site info",
			TitleScore: "system",
		}, {
			Exist:      !eng.config.HidePluginEntrance,
			Icon:       icon.Th,
			BtnName:    types.NavBtnPlugName,
			URL:        "/plugins",
			Title:      "plugins",
			TitleScore: "plugin",
		},
	}
}

// initSiteSetting 初始化站点设置
//
// 工作原理：
//   - 从数据库加载站点配置
//   - 使用配置初始化Site模型
//   - 将配置更新到数据库
//   - 将配置服务添加到服务列表
//   - 初始化错误处理
//
// 使用场景：
//   - 内部方法，由Use调用
//   - 初始化站点配置
func (eng *Engine) initSiteSetting() {

	printInitMsg(language.Get("initialize configuration"))

	err := eng.config.Update(models.Site().
		SetConn(eng.DefaultConnection()).
		Init(eng.config.ToMap()).
		AllToMap())
	if err != nil {
		logger.Panic(err)
	}
	eng.Services.Add("config", config.SrvWithConfig(eng.config))

	errors.Init()
}

// ============================
// HTML 内容渲染 API
// ============================

// Content 调用Engine适配器的Content方法
//
// 参数说明：
//   - ctx: Web框架上下文
//   - panel: 面板生成函数
//
// 工作原理：
//   - 检查适配器是否为空，为空则panic
//   - 调用适配器的Content方法
//   - 传递admin插件的AddOperationFn
//   - 传递导航栏按钮
//
// 使用场景：
//   - 渲染GoAdmin页面
//   - 显示管理面板
//
// 使用示例：
//
//	eng.Content(c, func(ctx *context.Context) types.Panel {
//	    return components.GetTable()
//	})
func (eng *Engine) Content(ctx interface{}, panel types.GetPanelFn) {
	if eng.Adapter == nil {
		emptyAdapterPanic()
	}
	eng.Adapter.Content(ctx, panel, eng.AdminPlugin().GetAddOperationFn(), *eng.NavButtons...)
}

// Content 调用defaultAdapter的Content方法
//
// 参数说明：
//   - ctx: Web框架上下文
//   - panel: 面板生成函数
//
// 工作原理：
//   - 检查默认适配器是否为空，为空则panic
//   - 调用默认适配器的Content方法
//   - 传递admin插件的AddOperationFn
//   - 传递全局导航栏按钮
//
// 使用场景：
//   - 渲染GoAdmin页面
//   - 显示管理面板
//
// 使用示例：
//
//	Content(c, func(ctx *context.Context) types.Panel {
//	    return components.GetTable()
//	})
func Content(ctx interface{}, panel types.GetPanelFn) {
	if defaultAdapter == nil {
		emptyAdapterPanic()
	}
	defaultAdapter.Content(ctx, panel, engine.AdminPlugin().GetAddOperationFn(), *navButtons...)
}

// Data 将路由和对应的处理器注入到Web框架
//
// 参数说明：
//   - method: HTTP方法（GET、POST等）
//   - url: 路由路径
//   - handler: 处理器函数
//   - noAuth: 是否需要认证，默认需要认证
//
// 工作原理：
//   - 如果noAuth为true，使用wrap包装处理器（不需要认证）
//   - 否则使用wrapWithAuthMiddleware包装处理器（需要认证）
//   - 将处理器添加到适配器
//
// 使用场景：
//   - 注册API路由
//   - 注册自定义处理器
//
// 使用示例：
//
//	eng.Data("GET", "/api/data", func(ctx *context.Context) {
//	    ctx.JSON(200, map[string]interface{}{"data": "ok"})
//	})
//	eng.Data("POST", "/api/data", handler, true) // 不需要认证
func (eng *Engine) Data(method, url string, handler context.Handler, noAuth ...bool) {
	if len(noAuth) > 0 && noAuth[0] {
		eng.Adapter.AddHandler(method, url, eng.wrap(handler))
	} else {
		eng.Adapter.AddHandler(method, url, eng.wrapWithAuthMiddleware(handler))
	}
}

// HTML将路由和对应的处理器注入到Web框架，处理器由给定函数包装
//
// 参数说明：
//   - method: HTTP方法（GET、POST等）
//   - url: 路由路径
//   - fn: 面板信息生成函数
//   - noAuth: 是否需要认证，默认需要认证
//
// 工作原理：
//   - 创建处理器函数
//   - 调用fn获取面板信息
//   - 如果出错则显示警告面板
//   - 执行面板的回调函数
//   - 获取模板和用户信息
//   - 渲染页面并返回HTML
//   - 根据noAuth参数决定是否添加认证中间件
//
// 使用场景：
//   - 注册HTML页面路由
//   - 渲染管理面板
//
// 使用示例：
//
//	eng.HTML("GET", "/dashboard", func(ctx *context.Context) (types.Panel, error) {
//	    return components.GetDashboard(), nil
//	})
func (eng *Engine) HTML(method, url string, fn types.GetPanelInfoFn, noAuth ...bool) {

	var handler = func(ctx *context.Context) {
		panel, err := fn(ctx)
		if err != nil {
			panel = template.WarningPanel(ctx, err.Error())
		}

		eng.AdminPlugin().GetAddOperationFn()(panel.Callbacks...)

		var (
			tmpl, tmplName = template.Default(ctx).GetTemplate(ctx.IsPjax())

			user = auth.Auth(ctx)
			buf  = new(bytes.Buffer)
		)

		hasError := tmpl.ExecuteTemplate(buf, tmplName, types.NewPage(ctx, &types.NewPageParam{
			User:         user,
			Menu:         menu.GetGlobalMenu(user, eng.Adapter.GetConnection(), ctx.Lang()).SetActiveClass(config.URLRemovePrefix(ctx.Path())),
			Panel:        panel.GetContent(eng.config.IsProductionEnvironment()),
			Assets:       template.GetComponentAssetImportHTML(ctx),
			Buttons:      eng.NavButtons.CheckPermission(user),
			TmplHeadHTML: template.Default(ctx).GetHeadHTML(),
			TmplFootJS:   template.Default(ctx).GetFootJS(),
			Iframe:       ctx.IsIframe(),
		}))

		if hasError != nil {
			logger.Error(fmt.Sprintf("错误：%s 适配器内容，", eng.Adapter.Name()), hasError)
		}

		ctx.HTMLByte(http.StatusOK, buf.Bytes())
	}

	if len(noAuth) > 0 && noAuth[0] {
		eng.Adapter.AddHandler(method, url, eng.wrap(handler))
	} else {
		eng.Adapter.AddHandler(method, url, eng.wrapWithAuthMiddleware(handler))
	}
}

// HTMLFile将路由和对应的处理器注入到Web框架，处理器返回给定HTML文件路径的面板内容
//
// 参数说明：
//   - method: HTTP方法（GET、POST等）
//   - url: 路由路径
//   - path: HTML文件路径
//   - data: 模板数据
//   - noAuth: 是否需要认证，默认需要认证
//
// 工作原理：
//   - 创建处理器函数
//   - 解析HTML文件模板
//   - 执行模板并将结果写入缓冲区
//   - 如果出错则显示错误面板
//   - 获取模板和用户信息
//   - 渲染页面并返回HTML
//   - 根据noAuth参数决定是否添加认证中间件
//
// 使用场景：
//   - 注册HTML文件路由
//   - 渲染自定义HTML页面
//
// 使用示例：
//
//	eng.HTMLFile("GET", "/custom", "views/custom.html", map[string]interface{}{
//	    "title": "自定义页面",
//	})
func (eng *Engine) HTMLFile(method, url, path string, data map[string]interface{}, noAuth ...bool) {

	var handler = func(ctx *context.Context) {

		cbuf := new(bytes.Buffer)

		t, err := template2.ParseFiles(path)
		if err != nil {
			eng.errorPanelHTML(ctx, cbuf, err)
			return
		} else if err := t.Execute(cbuf, data); err != nil {
			eng.errorPanelHTML(ctx, cbuf, err)
			return
		}

		var (
			tmpl, tmplName = template.Default(ctx).GetTemplate(ctx.IsPjax())

			user = auth.Auth(ctx)
			buf  = new(bytes.Buffer)
		)

		hasError := tmpl.ExecuteTemplate(buf, tmplName, types.NewPage(ctx, &types.NewPageParam{
			User: user,
			Menu: menu.GetGlobalMenu(user, eng.Adapter.GetConnection(), ctx.Lang()).SetActiveClass(eng.config.URLRemovePrefix(ctx.Path())),
			Panel: types.Panel{
				Content: template.HTML(cbuf.String()),
			},
			Assets:       template.GetComponentAssetImportHTML(ctx),
			Buttons:      eng.NavButtons.CheckPermission(user),
			TmplHeadHTML: template.Default(ctx).GetHeadHTML(),
			TmplFootJS:   template.Default(ctx).GetFootJS(),
			Iframe:       ctx.IsIframe(),
		}))

		if hasError != nil {
			logger.Error(fmt.Sprintf("错误：%s 适配器内容，", eng.Adapter.Name()), hasError)
		}

		ctx.HTMLByte(http.StatusOK, buf.Bytes())
	}

	if len(noAuth) > 0 && noAuth[0] {
		eng.Adapter.AddHandler(method, url, eng.wrap(handler))
	} else {
		eng.Adapter.AddHandler(method, url, eng.wrapWithAuthMiddleware(handler))
	}
}

// HTMLFiles将路由和对应的处理器注入到Web框架，处理器返回给定HTML文件路径的面板内容
//
// 参数说明：
//   - method: HTTP方法（GET、POST等）
//   - url: 路由路径
//   - data: 模板数据
//   - files: HTML文件路径列表
//
// 工作原理：
//   - 使用htmlFilesHandler创建处理器
//   - 添加认证中间件
//   - 将处理器添加到适配器
//
// 使用场景：
//   - 注册多个HTML文件路由
//   - 渲染自定义HTML页面
//
// 使用示例：
//
//	eng.HTMLFiles("GET", "/custom", map[string]interface{}{
//	    "title": "自定义页面",
//	}, "views/header.html", "views/content.html", "views/footer.html")
func (eng *Engine) HTMLFiles(method, url string, data map[string]interface{}, files ...string) {
	eng.Adapter.AddHandler(method, url, eng.wrapWithAuthMiddleware(eng.htmlFilesHandler(data, files...)))
}

// HTMLFilesNoAuth将路由和对应的处理器注入到Web框架，处理器返回给定HTML文件路径的面板内容，不需要认证
//
// 参数说明：
//   - method: HTTP方法（GET、POST等）
//   - url: 路由路径
//   - data: 模板数据
//   - files: HTML文件路径列表
//
// 工作原理：
//   - 使用htmlFilesHandler创建处理器
//   - 不添加认证中间件
//   - 将处理器添加到适配器
//
// 使用场景：
//   - 注册公开访问的HTML文件路由
//   - 渲染自定义HTML页面
//
// 使用示例：
//
//	eng.HTMLFilesNoAuth("GET", "/public", map[string]interface{}{
//	    "title": "公开页面",
//	}, "views/public.html")
func (eng *Engine) HTMLFilesNoAuth(method, url string, data map[string]interface{}, files ...string) {
	eng.Adapter.AddHandler(method, url, eng.wrap(eng.htmlFilesHandler(data, files...)))
}

// htmlFilesHandler 创建处理器，返回给定HTML文件路径的面板内容
//
// 参数说明：
//   - data: 模板数据
//   - files: HTML文件路径列表
//
// 返回值：
//   - context.Handler: 处理器函数
//
// 工作原理：
//   - 解析HTML文件模板
//   - 执行模板并将结果写入缓冲区
//   - 如果出错则显示错误面板
//   - 获取模板和用户信息
//   - 渲染页面并返回HTML
//
// 使用场景：
//   - 内部方法，由HTMLFiles和HTMLFilesNoAuth调用
//   - 创建HTML文件处理器
func (eng *Engine) htmlFilesHandler(data map[string]interface{}, files ...string) context.Handler {
	return func(ctx *context.Context) {

		cbuf := new(bytes.Buffer)

		t, err := template2.ParseFiles(files...)
		if err != nil {
			eng.errorPanelHTML(ctx, cbuf, err)
			return
		} else if err := t.Execute(cbuf, data); err != nil {
			eng.errorPanelHTML(ctx, cbuf, err)
			return
		}

		var (
			tmpl, tmplName = template.Default(ctx).GetTemplate(ctx.IsPjax())

			user = auth.Auth(ctx)
			buf  = new(bytes.Buffer)
		)

		hasError := tmpl.ExecuteTemplate(buf, tmplName, types.NewPage(ctx, &types.NewPageParam{
			User: user,
			Menu: menu.GetGlobalMenu(user, eng.Adapter.GetConnection(), ctx.Lang()).SetActiveClass(eng.config.URLRemovePrefix(ctx.Path())),
			Panel: types.Panel{
				Content: template.HTML(cbuf.String()),
			},
			Assets:       template.GetComponentAssetImportHTML(ctx),
			Buttons:      eng.NavButtons.CheckPermission(user),
			TmplHeadHTML: template.Default(ctx).GetHeadHTML(),
			TmplFootJS:   template.Default(ctx).GetFootJS(),
			Iframe:       ctx.IsIframe(),
		}))

		if hasError != nil {
			logger.Error(fmt.Sprintf("错误：%s 适配器内容，", eng.Adapter.Name()), hasError)
		}

		ctx.HTMLByte(http.StatusOK, buf.Bytes())
	}
}

// errorPanelHTML 将错误面板HTML添加到上下文响应
//
// 参数说明：
//   - ctx: 上下文对象
//   - buf: 缓冲区
//   - err: 错误对象
//
// 工作原理：
//   - 获取模板和用户信息
//   - 渲染错误面板
//   - 返回HTML响应
//
// 使用场景：
//   - 内部方法，用于显示错误页面
//   - 错误处理
func (eng *Engine) errorPanelHTML(ctx *context.Context, buf *bytes.Buffer, err error) {
	user := auth.Auth(ctx)
	tmpl, tmplName := template.Default(ctx).GetTemplate(ctx.IsPjax())

	hasError := tmpl.ExecuteTemplate(buf, tmplName, types.NewPage(ctx, &types.NewPageParam{
		User:         user,
		Menu:         menu.GetGlobalMenu(user, eng.Adapter.GetConnection(), ctx.Lang()).SetActiveClass(eng.config.URLRemovePrefix(ctx.Path())),
		Panel:        template.WarningPanel(ctx, err.Error()).GetContent(eng.config.IsProductionEnvironment()),
		Assets:       template.GetComponentAssetImportHTML(ctx),
		Buttons:      (*eng.NavButtons).CheckPermission(user),
		TmplHeadHTML: template.Default(ctx).GetHeadHTML(),
		TmplFootJS:   template.Default(ctx).GetFootJS(),
		Iframe:       ctx.IsIframe(),
	}))

	if hasError != nil {
		logger.Error(fmt.Sprintf("错误：%s 适配器内容，", eng.Adapter.Name()), hasError)
	}

	ctx.HTMLByte(http.StatusOK, buf.Bytes())
}

// ============================
// Admin 插件 API
// ============================

// AddGenerators 添加admin生成器
//
// 参数说明：
//   - list: 生成器列表
//
// 返回值：
//   - *Engine: 返回Engine本身，支持链式调用
//
// 工作原理：
//   - 查找admin插件
//   - 如果存在则添加生成器
//   - 如果不存在则创建新的admin插件并添加生成器
//
// 使用场景：
//   - 添加表格模型生成器
//   - 扩展admin功能
//
// 使用示例：
//
//	eng.AddGenerators(table.GeneratorList{
//	    table.GetGeneratorsForModel(&User{}),
//	})
func (eng *Engine) AddGenerators(list ...table.GeneratorList) *Engine {
	plug, exist := eng.FindPluginByName("admin")
	if exist {
		plug.(*admin.Admin).AddGenerators(list...)
		return eng
	}
	eng.PluginList = append(eng.PluginList, admin.NewAdmin(list...))
	return eng
}

// AdminPlugin 获取admin插件，如果不存在则创建一个
//
// 返回值：
//   - *admin.Admin: admin插件实例
//
// 工作原理：
//   - 查找admin插件
//   - 如果存在则返回
//   - 如果不存在则创建新的admin插件并添加到插件列表
//
// 使用场景：
//   - 获取admin插件
//   - 配置admin功能
//
// 使用示例：
//
//	adm := eng.AdminPlugin()
//	adm.SetCaptcha(map[string]string{"driver": "recaptcha"})
func (eng *Engine) AdminPlugin() *admin.Admin {
	plug, exist := eng.FindPluginByName("admin")
	if exist {
		return plug.(*admin.Admin)
	}
	adm := admin.NewAdmin()
	eng.PluginList = append(eng.PluginList, adm)
	return adm
}

// SetCaptcha 设置验证码配置
//
// 参数说明：
//   - captcha: 验证码配置
//
// 返回值：
//   - *Engine: 返回Engine本身，支持链式调用
//
// 工作原理：
//   - 获取admin插件
//   - 设置验证码配置
//
// 使用场景：
//   - 配置验证码
//   - 防止机器人攻击
//
// 使用示例：
//
//	eng.SetCaptcha(map[string]string{
//	    "driver": "recaptcha",
//	    "site_key": "xxx",
//	    "secret_key": "yyy",
//	})
func (eng *Engine) SetCaptcha(captcha map[string]string) *Engine {
	eng.AdminPlugin().SetCaptcha(captcha)
	return eng
}

// SetCaptchaDriver 使用驱动设置验证码配置
//
// 参数说明：
//   - driver: 验证码驱动名称
//
// 返回值：
//   - *Engine: 返回Engine本身，支持链式调用
//
// 工作原理：
//   - 获取admin插件
//   - 设置验证码驱动
//
// 使用场景：
//   - 配置验证码驱动
//   - 快速设置验证码
//
// 使用示例：
//
//	eng.SetCaptchaDriver("recaptcha")
func (eng *Engine) SetCaptchaDriver(driver string) *Engine {
	eng.AdminPlugin().SetCaptcha(map[string]string{"driver": driver})
	return eng
}

// AddGenerator 添加表格模型生成器
//
// 参数说明：
//   - key: 生成器键名
//   - g: 生成器对象
//
// 返回值：
//   - *Engine: 返回Engine本身，支持链式调用
//
// 工作原理：
//   - 获取admin插件
//   - 添加生成器到admin插件
//
// 使用场景：
//   - 添加单个表格生成器
//   - 管理数据表
//
// 使用示例：
//
//	eng.AddGenerator("user", table.GetGeneratorsForModel(&User{}))
func (eng *Engine) AddGenerator(key string, g table.Generator) *Engine {
	eng.AdminPlugin().AddGenerator(key, g)
	return eng
}

// AddGlobalDisplayProcessFn 调用types.AddGlobalDisplayProcessFn
//
// 参数说明：
//   - f: 字段过滤函数
//
// 返回值：
//   - *Engine: 返回Engine本身，支持链式调用
//
// 工作原理：
//   - 添加全局显示处理函数
//   - 用于处理字段显示
//
// 使用场景：
//   - 自定义字段显示逻辑
//   - 全局字段处理
//
// 使用示例：
//
//	eng.AddGlobalDisplayProcessFn(func(value types.FieldModel) interface{} {
//	    return strings.ToUpper(value.Value)
//	})
func (eng *Engine) AddGlobalDisplayProcessFn(f types.FieldFilterFn) *Engine {
	types.AddGlobalDisplayProcessFn(f)
	return eng
}

// AddDisplayFilterLimit 调用types.AddDisplayFilterLimit
//
// 参数说明：
//   - limit: 显示长度限制
//
// 返回值：
//   - *Engine: 返回Engine本身，支持链式调用
//
// 工作原理：
//   - 设置字段显示长度限制
//   - 超过限制的文本将被截断
//
// 使用场景：
//   - 限制字段显示长度
//   - 防止长文本破坏布局
//
// 使用示例：
//
//	eng.AddDisplayFilterLimit(50)
func (eng *Engine) AddDisplayFilterLimit(limit int) *Engine {
	types.AddLimit(limit)
	return eng
}

// AddDisplayFilterTrimSpace 调用types.AddDisplayFilterTrimSpace
//
// 返回值：
//   - *Engine: 返回Engine本身，支持链式调用
//
// 工作原理：
//   - 添加字段显示过滤函数
//   - 自动去除字段值的首尾空格
//
// 使用场景：
//   - 去除字段值空格
//   - 数据清洗
//
// 使用示例：
//
//	eng.AddDisplayFilterTrimSpace()
func (eng *Engine) AddDisplayFilterTrimSpace() *Engine {
	types.AddTrimSpace()
	return eng
}

// AddDisplayFilterSubstr 调用types.AddDisplayFilterSubstr
//
// 参数说明：
//   - start: 起始位置
//   - end: 结束位置
//
// 返回值：
//   - *Engine: 返回Engine本身，支持链式调用
//
// 工作原理：
//   - 添加字段显示过滤函数
//   - 截取字段值的指定范围
//
// 使用场景：
//   - 截取字段值
//   - 部分显示
//
// 使用示例：
//
//	eng.AddDisplayFilterSubstr(0, 10)
func (eng *Engine) AddDisplayFilterSubstr(start int, end int) *Engine {
	types.AddSubstr(start, end)
	return eng
}

// AddDisplayFilterToTitle 调用types.AddDisplayFilterToTitle
//
// 返回值：
//   - *Engine: 返回Engine本身，支持链式调用
//
// 工作原理：
//   - 添加字段显示过滤函数
//   - 将字段值转换为标题格式（首字母大写）
//
// 使用场景：
//   - 标题格式化
//   - 字段值美化
//
// 使用示例：
//
//	eng.AddDisplayFilterToTitle()
func (eng *Engine) AddDisplayFilterToTitle() *Engine {
	types.AddToTitle()
	return eng
}

// AddDisplayFilterToUpper 调用types.AddDisplayFilterToUpper
//
// 返回值：
//   - *Engine: 返回Engine本身，支持链式调用
//
// 工作原理：
//   - 添加字段显示过滤函数
//   - 将字段值转换为大写
//
// 使用场景：
//   - 大写转换
//   - 字段值格式化
//
// 使用示例：
//
//	eng.AddDisplayFilterToUpper()
func (eng *Engine) AddDisplayFilterToUpper() *Engine {
	types.AddToUpper()
	return eng
}

// AddDisplayFilterToLower 调用types.AddDisplayFilterToLower
//
// 返回值：
//   - *Engine: 返回Engine本身，支持链式调用
//
// 工作原理：
//   - 添加字段显示过滤函数
//   - 将字段值转换为小写
//
// 使用场景：
//   - 小写转换
//   - 字段值格式化
//
// 使用示例：
//
//	eng.AddDisplayFilterToLower()
func (eng *Engine) AddDisplayFilterToLower() *Engine {
	types.AddToLower()
	return eng
}

// AddDisplayFilterXssFilter 调用types.AddDisplayFilterXssFilter
//
// 返回值：
//   - *Engine: 返回Engine本身，支持链式调用
//
// 工作原理：
//   - 添加字段显示过滤函数
//   - 过滤XSS攻击代码
//
// 使用场景：
//   - XSS防护
//   - 安全过滤
//
// 使用示例：
//
//	eng.AddDisplayFilterXssFilter()
func (eng *Engine) AddDisplayFilterXssFilter() *Engine {
	types.AddXssFilter()
	return eng
}

// AddDisplayFilterXssJsFilter 调用types.AddDisplayFilterXssJsFilter
//
// 返回值：
//   - *Engine: 返回Engine本身，支持链式调用
//
// 工作原理：
//   - 添加字段显示过滤函数
//   - 过滤JavaScript XSS攻击代码
//
// 使用场景：
//   - JavaScript XSS防护
//   - 安全过滤
//
// 使用示例：
//
//	eng.AddDisplayFilterXssJsFilter()
func (eng *Engine) AddDisplayFilterXssJsFilter() *Engine {
	types.AddXssJsFilter()
	return eng
}

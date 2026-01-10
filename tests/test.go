// test.go - 测试框架核心文件
// 包名：tests
// 作者：GoAdmin 团队
// 创建日期：2020年
// 描述：本文件提供 GoAdmin 测试框架的核心功能，包括数据清理、测试套件执行、
//       错误处理等。支持多种数据库（MySQL、PostgreSQL、SQLite、MSSQL）和
//       多种 Web 框架（Gin、Echo、Chi、FastHTTP 等）的测试
//
// 主要功能：
//   - Cleaner: 清理测试数据并插入初始测试数据
//   - BlackBoxTestSuitOfBuiltInTables: 执行内置表的黑盒测试
//   - BlackBoxTestSuit: 通用黑盒测试套件框架
//   - checkErr: 错误检查辅助函数
//
// 支持的数据库：
//   - MySQL: 使用 DriverMysql
//   - PostgreSQL: 使用 DriverPostgresql
//   - SQLite: 使用 DriverSqlite
//   - MSSQL: 使用 DriverMssql
//
// 支持的 Web 框架：
//   - Gin
//   - Echo
//   - Chi
//   - FastHTTP
//   - Beego
//   - Buffalo
//   - Gear
//   - GoFrame (GF)
//   - Iris
//   - Gorilla Mux
//   - Net/http

package tests

import (
	"net/http" // Go 标准网络包，提供 HTTP 客户端和服务端功能
	"strings"  // Go 标准字符串包，提供字符串操作功能
	"testing"  // Go 标准测试包

	"github.com/gavv/httpexpect"                                // HTTP 测试库，用于编写 HTTP 断言
	"github.com/purpose168/GoAdmin/modules/config"              // 配置模块
	"github.com/purpose168/GoAdmin/modules/db"                  // 数据库模块
	"github.com/purpose168/GoAdmin/modules/db/dialect"          // 数据库方言模块
	"github.com/purpose168/GoAdmin/plugins/admin/modules/table" // 表生成器模块
	"github.com/purpose168/GoAdmin/tests/common"                // 测试公共模块
	"github.com/purpose168/GoAdmin/tests/frameworks/fasthttp"   // FastHTTP 框架适配器
	fasthttp2 "github.com/valyala/fasthttp"                     // FastHTTP 库（使用别名避免冲突）
)

// ==================== 数据清理函数 ====================

// Cleaner 清理测试数据并插入初始测试数据
// 参数：config - 数据库配置列表
// 说明：
//
//	该函数执行以下操作：
//	1. 安全检查：验证数据库名称或 DSN 包含 "test" 字符串，防止误操作生产数据库
//	2. 清理数据：删除所有测试表中的数据
//	3. 重置自增 ID：根据数据库类型重置自增序列
//	4. 插入初始数据：插入测试所需的初始用户、角色、权限、菜单等数据
//
// 测试表列表：
//   - goadmin_users: 用户表
//   - goadmin_user_permissions: 用户权限关联表
//   - goadmin_session: 会话表
//   - goadmin_roles: 角色表
//   - goadmin_role_users: 角色用户关联表
//   - goadmin_role_permissions: 角色权限关联表
//   - goadmin_role_menu: 角色菜单关联表
//   - goadmin_permissions: 权限表
//   - goadmin_operation_log: 操作日志表
//   - goadmin_menu: 菜单表
//
// 初始数据：
//   - 用户：admin（管理员）、operator（操作员）
//   - 角色：Administrator（管理员）、Operator（操作员）
//   - 权限：All permission（所有权限）、Dashboard（仪表板）
//   - 菜单：Dashboard、Admin、Users、Permission、Menu、Operation log、User、test2 menu
//
// 数据库支持：
//   - MySQL: 使用 ALTER TABLE ... AUTO_INCREMENT = 1 重置自增
//   - PostgreSQL: 使用 ALTER SEQUENCE ... RESTART WITH 1 重置序列
//   - SQLite: 使用 update sqlite_sequence set seq = 0 重置序列
//   - MSSQL: 使用 DBCC CHECKIDENT (... RESEED, 0) 重置标识
//
// 注意事项：
//   - 如果数据库名称或 DSN 不包含 "test"，函数会 panic
//   - 该函数会修改数据库内容，确保使用测试数据库
//   - 密码使用 bcrypt 加密，admin 密码为 "admin"，operator 密码为 "operator"
func Cleaner(config config.DatabaseList) {
	// 构建检查语句，用于验证是否为测试数据库
	checkStatement := ""

	// 根据数据库驱动类型和配置构建检查语句
	if config.GetDefault().Driver != "sqlite" { // 非 SQLite 数据库
		if config.GetDefault().Dsn == "" { // 如果 DSN 为空，使用数据库名称
			checkStatement = config.GetDefault().Name
		} else { // 否则使用 DSN
			checkStatement = config.GetDefault().Dsn
		}
	} else { // SQLite 数据库
		if config.GetDefault().Dsn == "" { // 如果 DSN 为空，使用文件路径
			checkStatement = config.GetDefault().File
		} else { // 否则使用 DSN
			checkStatement = config.GetDefault().Dsn
		}
	}

	// 安全检查：确保数据库名称或 DSN 包含 "test" 字符串
	// 防止误操作生产数据库
	if !strings.Contains(checkStatement, "test") {
		panic("wrong database") // 如果不包含 "test"，抛出异常
	}

	// 定义所有需要清理的测试表
	var allTables = [...]string{
		"goadmin_users",            // 用户表
		"goadmin_user_permissions", // 用户权限关联表
		"goadmin_session",          // 会话表
		"goadmin_roles",            // 角色表
		"goadmin_role_users",       // 角色用户关联表
		"goadmin_role_permissions", // 角色权限关联表
		"goadmin_role_menu",        // 角色菜单关联表
		"goadmin_permissions",      // 权限表
		"goadmin_operation_log",    // 操作日志表
		"goadmin_menu",             // 菜单表
	}

	// 定义需要重置自增 ID 的表
	var autoIncrementTable = [...]string{
		"goadmin_menu",        // 菜单表
		"goadmin_permissions", // 权限表
		"goadmin_roles",       // 角色表
		"goadmin_users",       // 用户表
	}

	// 定义需要插入的初始测试数据
	var insertData = map[string][]dialect.H{
		// 用户数据：admin（管理员）和 operator（操作员）
		"goadmin_users": {
			{"username": "admin", "name": "admin", "password": "$2a$10$TEDU/aUxLkr2wCxGxI62/.yOtzrzfv426DLLdyha9H2GpWRggB0di", "remember_token": "tlNcBVK9AvfYH7WEnwB1RKvocJu8FfRy4um3DJtwdHuJy0dwFsLOgAc0xUfh"},      // admin 用户，密码为 "admin"
			{"username": "operator", "name": "operator", "password": "$2a$10$rVqOzHjN2MdlEprRflb1eGP0oZXuSrbJLOmJagFsCd81YZm0bsh.", "remember_token": "tlNcBVK9AvfYH7WEnwB1RKvocJu8FfRy4um3DJtwdHuJy0dwFsLOgAc0xUfh"}, // operator 用户，密码为 "operator"
		},
		// 角色数据：Administrator（管理员）和 Operator（操作员）
		"goadmin_roles": {
			{"name": "Administrator", "slug": "administrator"}, // 管理员角色
			{"name": "Operator", "slug": "operator"},           // 操作员角色
		},
		// 权限数据：所有权限和仪表板权限
		"goadmin_permissions": {
			{"name": "All permission", "slug": "*", "http_method": "", "http_path": "*"},                       // 所有权限
			{"name": "Dashboard", "slug": "dashboard", "http_method": "GET,PUT,POST,DELETE", "http_path": "/"}, // 仪表板权限
		},
		// 菜单数据：系统菜单结构
		"goadmin_menu": {
			{"parent_id": 0, "type": 1, "order": 2, "title": "Admin", "icon": "fa-tasks", "uri": ""},                       // 管理菜单
			{"parent_id": 1, "type": 1, "order": 2, "title": "Users", "icon": "fa-users", "uri": "/info/manager"},          // 用户管理菜单
			{"parent_id": 0, "type": 1, "order": 3, "title": "test2 menu", "icon": "fa-angellist", "uri": "/example/test"}, // 测试菜单
			{"parent_id": 1, "type": 1, "order": 4, "title": "Permission", "icon": "fa-ban", "uri": "/info/permission"},    // 权限管理菜单
			{"parent_id": 1, "type": 1, "order": 5, "title": "Menu", "icon": "fa-bars", "uri": "/menu"},                    // 菜单管理菜单
			{"parent_id": 1, "type": 1, "order": 6, "title": "Operation log", "icon": "fa-history", "uri": "/info/op"},     // 操作日志菜单
			{"parent_id": 0, "type": 1, "order": 1, "title": "Dashboard", "icon": "fa-bar-chart", "uri": "/"},              // 仪表板菜单
			{"parent_id": 0, "type": 1, "order": 7, "title": "User", "icon": "fa-users", "uri": "/info/user"},              // 用户信息菜单
		},
		// 角色用户关联数据
		"goadmin_role_users": {
			{"user_id": 1, "role_id": 1}, // admin 用户关联 Administrator 角色
			{"user_id": 2, "role_id": 2}, // operator 用户关联 Operator 角色
		},
		// 用户权限关联数据
		"goadmin_user_permissions": {
			{"user_id": 1, "permission_id": 1}, // admin 用户拥有所有权限
			{"user_id": 2, "permission_id": 2}, // operator 用户拥有仪表板权限
		},
		// 角色权限关联数据
		"goadmin_role_permissions": {
			{"role_id": 1, "permission_id": 1}, // Administrator 角色拥有所有权限
			{"role_id": 1, "permission_id": 2}, // Administrator 角色拥有仪表板权限
			{"role_id": 2, "permission_id": 2}, // Operator 角色拥有仪表板权限
		},
		// 角色菜单关联数据
		"goadmin_role_menu": {
			{"role_id": 1, "menu_id": 1}, // Administrator 角色可以访问 Admin 菜单
			{"role_id": 1, "menu_id": 7}, // Administrator 角色可以访问 Dashboard 菜单
			{"role_id": 2, "menu_id": 7}, // Operator 角色可以访问 Dashboard 菜单
			{"role_id": 1, "menu_id": 8}, // Administrator 角色可以访问 User 菜单
			{"role_id": 2, "menu_id": 8}, // Operator 角色可以访问 User 菜单
			{"role_id": 1, "menu_id": 3}, // Administrator 角色可以访问 test2 menu
		},
	}

	// 根据配置初始化数据库连接
	conn := db.GetConnectionByDriver(config.GetDefault().Driver).InitDB(config)

	// 清理数据：删除所有测试表中的数据
	for _, t := range allTables {
		_ = db.WithDriver(conn).Table(t).Delete() // 删除表中的所有数据
	}

	// 重置自增 ID：根据数据库类型执行不同的重置命令
	switch config.GetDefault().Driver {
	case db.DriverMysql: // MySQL 数据库
		for _, t := range autoIncrementTable {
			checkErr(conn.Exec(`ALTER TABLE ` + t + ` AUTO_INCREMENT = 1`)) // 重置自增 ID 为 1
		}
	case db.DriverMssql: // MSSQL 数据库
		for _, t := range autoIncrementTable {
			checkErr(conn.Exec(`DBCC CHECKIDENT (` + t + `, RESEED, 0)`)) // 重置标识种子为 0
		}
	case db.DriverPostgresql: // PostgreSQL 数据库
		for _, t := range autoIncrementTable {
			checkErr(conn.Exec(`ALTER SEQUENCE ` + t + `_myid_seq RESTART WITH  1`)) // 重置序列从 1 开始
		}
	case db.DriverSqlite: // SQLite 数据库
		for _, t := range autoIncrementTable {
			checkErr(conn.Exec(`update sqlite_sequence set seq = 0 where name = '` + t + `'`)) // 重置序列为 0
		}
	}

	// 插入初始数据：遍历所有表并插入测试数据
	for t, data := range insertData {
		for _, d := range data {
			checkErr(db.WithDriver(conn).Table(t).Insert(d)) // 插入数据并检查错误
		}
	}
}

// ==================== 测试套件函数 ====================

// BlackBoxTestSuitOfBuiltInTables 内置表黑盒测试套件
// 参数：
//   - t: 测试对象，用于测试断言和报告
//   - fn: 处理器生成函数，用于创建 HTTP 处理器
//   - config: 数据库配置列表
//   - isFasthttp: 可变参数，是否使用 FastHTTP 框架（默认为 false）
//
// 说明：
//
//	该函数是内置表黑盒测试的便捷包装器，调用通用的 BlackBoxTestSuit 函数。
//	使用 Cleaner 函数清理和初始化测试数据，使用 common.Test 函数执行测试。
//
// 测试内容：
//   - 认证测试（登录、登出）
//   - 权限管理测试（增删改查）
//   - 角色管理测试（增删改查、权限分配）
//   - 管理员管理测试（增删改查、角色分配）
//   - 菜单管理测试（增删改查、层级管理）
//   - 操作日志测试（记录查询）
//
// 使用示例：
//
//	func TestMyTest(t *testing.T) {
//	    BlackBoxTestSuitOfBuiltInTables(t, gin.NewHandler, myConfig)
//	}
func BlackBoxTestSuitOfBuiltInTables(t *testing.T, fn HandlerGenFn, config config.DatabaseList, isFasthttp ...bool) {
	// 调用通用黑盒测试套件，传入处理器生成函数、数据库配置、
	// 数据清理函数（Cleaner）和测试函数（common.Test）
	BlackBoxTestSuit(t, fn, config, nil, Cleaner, common.Test, isFasthttp...)
}

// ==================== 辅助函数 ====================

// checkErr 错误检查函数
// 参数：
//   - _: 忽略的返回值（通常用于忽略成功时的返回值）
//   - err: 错误对象
//
// 说明：
//
//	该函数检查错误是否存在，如果存在则 panic。
//	用于简化错误处理代码，避免重复的 if err != nil 检查。
//
// 使用示例：
//
//	checkErr(db.Exec("INSERT INTO ..."))
//	checkErr(http.Get(url))
func checkErr(_ interface{}, err error) {
	if err != nil { // 如果错误不为 nil
		panic(err) // 抛出异常
	}
}

// ==================== 通用测试套件函数 ====================

// BlackBoxTestSuit 通用黑盒测试套件
// 参数：
//   - t: 测试对象，用于测试断言和报告
//   - fn: 处理器生成函数，用于创建 HTTP 处理器
//   - config: 数据库配置列表
//   - gens: 表生成器列表（可选，用于自定义表）
//   - cleaner: 数据清理函数，用于清理和初始化测试数据
//   - tester: 测试函数，用于执行具体的测试逻辑
//   - isFasthttp: 可变参数，是否使用 FastHTTP 框架（默认为 false）
//
// 说明：
//
//	该函数提供通用的黑盒测试框架，支持多种 Web 框架和数据库。
//	执行流程：
//	1. 调用 cleaner 函数清理和初始化测试数据
//	2. 根据 isFasthttp 参数选择测试框架（FastHTTP 或标准 HTTP）
//	3. 创建 httpexpect 客户端并执行测试
//
// 支持的框架：
//   - 标准框架（Gin、Echo、Chi 等）：使用 httpexpect.NewBinder
//   - FastHTTP 框架：使用 httpexpect.NewFastBinder
//
// 使用场景：
//   - 编写自定义测试套件
//   - 测试自定义表和功能
//   - 集成第三方框架
//
// 使用示例：
//
//	func TestMyCustomTest(t *testing.T) {
//	    BlackBoxTestSuit(t, gin.NewHandler, myConfig, myGenerators, myCleaner, myTester)
//	}
func BlackBoxTestSuit(t *testing.T, fn HandlerGenFn,
	config config.DatabaseList,
	gens table.GeneratorList,
	cleaner DataCleaner,
	tester Tester, isFasthttp ...bool) {
	// 清理数据：调用数据清理函数
	cleaner(config)

	// 执行测试：根据 isFasthttp 参数选择测试框架
	if len(isFasthttp) > 0 && isFasthttp[0] { // 如果使用 FastHTTP 框架
		tester(httpexpect.WithConfig(httpexpect.Config{
			Client: &http.Client{
				Transport: httpexpect.NewFastBinder(fasthttp.NewHandler(config, gens)), // 使用 FastHTTP 绑定器
				Jar:       httpexpect.NewJar(),                                         // 创建 Cookie Jar
			},
			Reporter: httpexpect.NewAssertReporter(t), // 创建断言报告器
		}))
	} else { // 使用标准 HTTP 框架
		tester(httpexpect.WithConfig(httpexpect.Config{
			Client: &http.Client{
				Transport: httpexpect.NewBinder(fn(config, gens)), // 使用标准 HTTP 绑定器
				Jar:       httpexpect.NewJar(),                    // 创建 Cookie Jar
			},
			Reporter: httpexpect.NewAssertReporter(t), // 创建断言报告器
		}))
	}
}

// ==================== 类型定义 ====================

// Tester 测试函数类型
// 描述：
//
//	定义了测试函数的类型签名，用于封装具体的测试逻辑
//
// 参数：
//   - e: httpexpect.Expect 对象，用于 HTTP 请求和断言
//
// 使用示例：
//
//	func myTest(e *httpexpect.Expect) {
//	    e.GET("/admin").Expect().Status(http.StatusOK)
//	}
type Tester func(e *httpexpect.Expect)

// DataCleaner 数据清理函数类型
// 描述：
//
//	定义了数据清理函数的类型签名，用于清理和初始化测试数据
//
// 参数：
//   - config: 数据库配置列表
//
// 使用示例：
//
//	func myCleaner(config config.DatabaseList) {
//	    // 清理和初始化测试数据
//	}
type DataCleaner func(config config.DatabaseList)

// HandlerGenFn 处理器生成函数类型
// 描述：
//
//	定义了 HTTP 处理器生成函数的类型签名，用于创建 HTTP 处理器
//
// 参数：
//   - config: 数据库配置列表
//   - gens: 表生成器列表
//
// 返回值：
//   - http.Handler: HTTP 处理器
//
// 使用示例：
//
//	func myHandler(config config.DatabaseList, gens table.GeneratorList) http.Handler {
//	    r := gin.Default()
//	    // 配置路由
//	    return r
//	}
type HandlerGenFn func(config config.DatabaseList, gens table.GeneratorList) http.Handler

// FasthttpHandlerGenFn FastHTTP 处理器生成函数类型
// 描述：
//
//	定义了 FastHTTP 处理器生成函数的类型签名，用于创建 FastHTTP 请求处理器
//
// 参数：
//   - config: 数据库配置列表
//   - gens: 表生成器列表
//
// 返回值：
//   - fasthttp2.RequestHandler: FastHTTP 请求处理器
//
// 使用示例：
//
//	func myFastHandler(config config.DatabaseList, gens table.GeneratorList) fasthttp.RequestHandler {
//	    return func(ctx *fasthttp.RequestCtx) {
//	        // 处理请求
//	    }
//	}
type FasthttpHandlerGenFn func(config config.DatabaseList, gens table.GeneratorList) fasthttp2.RequestHandler

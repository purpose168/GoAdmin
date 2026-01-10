// common.go - 通用测试文件
// 包名：common
// 作者：GoAdmin 团队
// 创建日期：2020年
// 描述：本文件包含 GoAdmin 管理插件的测试套件，使用 httpexpect 库进行 HTTP 请求测试
//       提供了完整的测试框架，包括认证、权限、角色、管理员、API、菜单、操作日志、外部数据源和表格测试

package common

import (
	"fmt"    // 格式化输出包，提供字符串格式化功能
	"regexp" // 正则表达式包，提供字符串匹配和提取功能

	"github.com/gavv/httpexpect"                               // httpexpect 包，提供 HTTP 请求测试功能
	"github.com/mgutz/ansi"                                    // ansi 包，提供终端颜色输出功能
	"github.com/purpose168/GoAdmin/plugins/admin/modules/form" // 表单模块，提供表单处理功能
)

// reg 正则表达式，用于从HTML中提取token值
//
// 说明:
//
//	该正则表达式用于从 HTML 表单中提取 CSRF token 的值
//	CSRF token 是一种安全机制，用于防止跨站请求伪造攻击
//	正则表达式匹配格式: <input type="hidden" name="__goadmin_token__" value='token值'>
var reg, _ = regexp.Compile("<input type=\"hidden\" name=\"" + form.TokenKey + "\" value='(.*?)'>")

// ExtraTest 执行GoAdmin管理插件的扩展测试
// 包含权限、角色、管理员、API、菜单、操作日志、外部数据源和正常表格测试
//
// 参数:
//   - e: HTTP期望对象，用于发送HTTP请求和验证响应
//
// 功能特性:
//   - 执行认证测试并获取会话 cookie
//   - 测试权限管理功能
//   - 测试角色管理功能
//   - 测试管理员管理功能
//   - 测试 API 接口功能
//   - 测试菜单管理功能
//   - 测试操作日志功能
//   - 测试外部数据源功能
//   - 测试正常表格功能
//
// 测试流程:
//  1. 执行认证测试，获取会话 cookie
//  2. 执行权限管理测试
//  3. 执行角色管理测试
//  4. 执行管理员管理测试
//  5. 执行 API 接口测试
//  6. 执行菜单管理测试
//  7. 执行操作日志测试
//  8. 执行外部数据源测试
//  9. 执行正常表格测试
//
// 说明:
//
//	该函数是扩展测试套件，包含了 GoAdmin 管理插件的所有功能测试
//	测试使用黑盒测试方法，通过 HTTP 请求验证各个功能是否正常工作
//	所有测试都使用 httpexpect 库进行 HTTP 请求和响应验证
//	认证测试返回的 cookie 用于后续所有测试的身份验证
func ExtraTest(e *httpexpect.Expect) {
	fmt.Println()                                                // 打印空行
	fmt.Println("============================================")  // 打印分隔线
	printlnWithColor("Basic Function Black-Box Testing", "blue") // 使用蓝色打印 "Basic Function Black-Box Testing"
	fmt.Println("============================================")  // 打印分隔线
	fmt.Println()                                                // 打印空行

	cookie := authTest(e)       // 执行认证测试并获取 cookie
	permissionTest(e, cookie)   // 权限检查测试
	roleTest(e, cookie)         // 角色检查测试
	managerTest(e, cookie)      // 管理员检查测试
	apiTest(e, cookie)          // API 检查测试
	menuTest(e, cookie)         // 菜单检查测试
	operationLogTest(e, cookie) // 操作日志检查测试
	externalTest(e, cookie)     // 从外部数据源获取数据检查测试
	normalTest(e, cookie)       // 正常表格测试
}

// Test 执行GoAdmin管理插件的基本功能测试
// 包含权限、角色、管理员、菜单和操作日志测试
//
// 参数:
//   - e: HTTP期望对象，用于发送HTTP请求和验证响应
//
// 功能特性:
//   - 执行认证测试并获取会话 cookie
//   - 测试权限管理功能
//   - 测试角色管理功能
//   - 测试管理员管理功能
//   - 测试菜单管理功能
//   - 测试操作日志功能
//
// 测试流程:
//  1. 执行认证测试，获取会话 cookie
//  2. 执行权限管理测试
//  3. 执行角色管理测试
//  4. 执行管理员管理测试
//  5. 执行菜单管理测试
//  6. 执行操作日志测试
//
// 说明:
//
//	该函数是基本测试套件，包含了 GoAdmin 管理插件的核心功能测试
//	测试使用黑盒测试方法，通过 HTTP 请求验证各个功能是否正常工作
//	所有测试都使用 httpexpect 库进行 HTTP 请求和响应验证
//	认证测试返回的 cookie 用于后续所有测试的身份验证
//	与 ExtraTest 的区别是不包含 API 测试、外部数据源测试和正常表格测试
func Test(e *httpexpect.Expect) {
	fmt.Println()                                                // 打印空行
	fmt.Println("============================================")  // 打印分隔线
	printlnWithColor("Basic Function Black-Box Testing", "blue") // 使用蓝色打印 "Basic Function Black-Box Testing"
	fmt.Println("============================================")  // 打印分隔线
	fmt.Println()                                                // 打印空行

	cookie := authTest(e)       // 执行认证测试并获取 cookie
	permissionTest(e, cookie)   // 权限检查测试
	roleTest(e, cookie)         // 角色检查测试
	managerTest(e, cookie)      // 管理员检查测试
	menuTest(e, cookie)         // 菜单检查测试
	operationLogTest(e, cookie) // 操作日志检查测试
}

// printlnWithColor 使用指定颜色打印消息
//
// 参数:
//   - msg: 要打印的消息
//   - color: 颜色名称（如"blue"、"green"等）
//
// 功能:
//   - 使用 ansi 库将消息以指定颜色输出到控制台
//
// 使用场景:
//   - 打印测试标题和重要信息
//   - 区分不同类型的输出信息
//   - 提高测试输出的可读性
//
// 说明:
//
//	该函数用于在测试输出中添加颜色，使测试结果更易于阅读
//	支持的颜色包括: black, red, green, yellow, blue, magenta, cyan, white 等
//	使用 ANSI 转义序列实现终端颜色输出
func printlnWithColor(msg string, color string) {
	fmt.Println(ansi.Color(msg, color)) // 使用 ansi 库将消息以指定颜色打印到控制台
}

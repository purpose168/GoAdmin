// auth.go - 认证测试文件
// 包名：common
// 作者：GoAdmin 团队
// 创建日期：2020年
// 描述：本文件包含认证相关的测试用例，使用 httpexpect 库进行 HTTP 请求测试
//       涵盖登录页面显示、空密码登录、正常登录、登出、会话管理、并发登录限制等核心功能测试

package common

import (
	"fmt"      // 格式化输出包，提供字符串格式化功能
	"net/http" // HTTP 包，提供 HTTP 客户端和服务器功能

	"github.com/gavv/httpexpect"                   // httpexpect 包，提供 HTTP 请求测试功能
	"github.com/purpose168/GoAdmin/modules/auth"   // 认证模块，提供认证和会话管理功能
	"github.com/purpose168/GoAdmin/modules/config" // 配置模块，提供配置管理功能
)

// authTest 执行认证测试
// 测试登录、登出、会话管理等认证相关功能
//
// 参数:
//   - e: HTTP期望对象，用于发送HTTP请求和验证响应
//
// 返回:
//   - *http.Cookie: 登录成功后的会话 cookie，用于后续请求的身份验证
//
// 功能特性:
//   - 测试登录页面显示，验证页面是否正常加载
//   - 测试空密码登录，验证错误处理
//   - 测试正常登录，验证登录是否成功并获取会话 cookie
//   - 测试未登录状态下的登出，验证错误处理
//   - 测试已登录状态下的登出，验证登出是否成功
//   - 测试再次登录，验证是否可以重新登录
//   - 测试限制同一用户同时登录，验证并发登录限制
//   - 测试 cookie 失效的情况，验证会话过期处理
//   - 测试登录成功，验证使用有效 cookie 是否可以访问受保护资源
//
// 测试流程:
//  1. 测试登录页面显示，验证 HTTP 状态码为 200
//  2. 测试空密码登录，验证 HTTP 状态码为 400
//  3. 测试正常登录，验证 HTTP 状态码为 200，获取会话 cookie
//  4. 测试未登录状态下的登出，验证 HTTP 状态码为 200
//  5. 测试已登录状态下的登出，验证 HTTP 状态码为 200
//  6. 测试再次登录，验证 HTTP 状态码为 200，获取新的会话 cookie
//  7. 测试限制同一用户同时登录，验证 HTTP 状态码为 200，获取新的会话 cookie
//  8. 测试 cookie 失效的情况，验证 HTTP 状态码为 200，响应体包含 "login"
//  9. 测试登录成功，验证 HTTP 状态码为 200，响应体包含 "Dashboard"
//
// 说明:
//
//	该函数使用 httpexpect 库进行 HTTP 请求测试，验证认证系统的各个功能是否正常工作
//	认证系统支持会话管理，通过 cookie 维护用户登录状态
//	认证系统支持并发登录限制，防止同一账号在多个设备同时登录
//	认证系统支持会话过期，cookie 失效后需要重新登录
//	认证系统支持登出功能，清除会话 cookie
//	认证系统使用默认 cookie 键名（auth.DefaultCookieKey）
//	认证系统支持自定义登录 URL（通过 config.GetLoginUrl() 获取）
func authTest(e *httpexpect.Expect) *http.Cookie {
	printlnWithColor("Auth", "blue")            // 使用蓝色打印 "Auth"
	fmt.Println("============================") // 打印分隔线

	// ==================== 登录页面显示测试 ====================
	// 测试登录页面显示，验证页面是否正常加载
	printlnWithColor("login: show", "green") // 使用绿色打印 "login: show"
	e.GET(config.Url(config.GetLoginUrl())). // 发送 GET 请求到登录页面
							Expect().   // 获取响应期望对象
							Status(200) // 验证 HTTP 状态码为 200

	// ==================== 测试空密码登录 ====================
	// 测试空密码登录，验证错误处理
	printlnWithColor("login: empty password", "green") // 使用绿色打印 "login: empty password"
	e.POST(config.Url("/signin")).                     // 发送 POST 请求到登录接口
								WithJSON(map[string]string{ // 添加 JSON 格式的表单数据
			"username": "admin", // 用户名为 "admin"
			"password": "",      // 密码为空
		}).Expect().Status(400) // 验证 HTTP 状态码为 400（错误）

	// ==================== 登录测试 ====================
	// 测试正常登录，验证登录是否成功并获取会话 cookie
	printlnWithColor("login", "green")      // 使用绿色打印 "login"
	sesID := e.POST(config.Url("/signin")). // 发送 POST 请求到登录接口
						WithForm(map[string]string{ // 添加表单数据
			"username": "admin", // 用户名为 "admin"
			"password": "admin", // 密码为 "admin"
		}).Expect().Status(200).       // 验证 HTTP 状态码为 200
		Cookie(auth.DefaultCookieKey). // 获取默认 cookie 键的值
		Raw()                          // 获取 cookie 的原始值

	// ==================== 未登录状态下登出测试 ====================
	// 测试未登录状态下的登出，验证错误处理
	printlnWithColor("logout: without login", "green") // 使用绿色打印 "logout: without login"
	e.GET(config.Url("/logout")).                      // 发送 GET 请求到登出接口
								Expect().   // 获取响应期望对象
								Status(200) // 验证 HTTP 状态码为 200

	// ==================== 登出测试 ====================
	// 测试已登录状态下的登出，验证登出是否成功
	printlnWithColor("logout", "green") // 使用绿色打印 "logout"
	e.GET(config.Url("/logout")).       // 发送 GET 请求到登出接口
						WithCookie(auth.DefaultCookieKey, sesID.Value). // 携带会话 cookie
						Expect().                                       // 获取响应期望对象
						Status(200)                                     // 验证 HTTP 状态码为 200

	// ==================== 再次登录测试 ====================
	// 测试再次登录，验证是否可以重新登录
	printlnWithColor("login again", "green")  // 使用绿色打印 "login again"
	cookie1 := e.POST(config.Url("/signin")). // 发送 POST 请求到登录接口
							WithForm(map[string]string{ // 添加表单数据
			"username": "admin", // 用户名为 "admin"
			"password": "admin", // 密码为 "admin"
		}).Expect().Status(200).       // 验证 HTTP 状态码为 200
		Cookie(auth.DefaultCookieKey). // 获取默认 cookie 键的值
		Raw()                          // 获取 cookie 的原始值

	// ==================== 测试限制同一用户同时登录 ====================
	// 测试限制同一用户同时登录，验证并发登录限制
	printlnWithColor("login again：restrict users from logging in at the same time", "green") // 使用绿色打印 "login again：restrict users from logging in at the same time"
	cookie2 := e.POST(config.Url("/signin")).                                                // 发送 POST 请求到登录接口
													WithForm(map[string]string{ // 添加表单数据
			"username": "admin", // 用户名为 "admin"
			"password": "admin", // 密码为 "admin"
		}).Expect().Status(200).       // 验证 HTTP 状态码为 200
		Cookie(auth.DefaultCookieKey). // 获取默认 cookie 键的值
		Raw()                          // 获取 cookie 的原始值

	// ==================== 登录成功测试 ====================
	// 测试 cookie 失效的情况，验证会话过期处理
	printlnWithColor("cookie failure", "green") // 使用绿色打印 "cookie failure"
	e.GET(config.Url("/")).                     // 发送 GET 请求到首页
							WithCookie(auth.DefaultCookieKey, cookie1.Value). // 携带第一个 cookie（已失效）
							Expect().                                         // 获取响应期望对象
							Status(200).                                      // 验证 HTTP 状态码为 200
							Body().Contains("login")                          // 验证响应体包含 "login"（重定向到登录页面）

	// ==================== 测试登录成功 ====================
	// 测试登录成功，验证使用有效 cookie 是否可以访问受保护资源
	printlnWithColor("login success", "green") // 使用绿色打印 "login success"
	e.GET(config.Url("/")).                    // 发送 GET 请求到首页
							WithCookie(auth.DefaultCookieKey, cookie2.Value). // 携带第二个 cookie（有效）
							Expect().                                         // 获取响应期望对象
							Status(200).                                      // 验证 HTTP 状态码为 200
							Body().Contains("Dashboard")                      // 验证响应体包含 "Dashboard"（成功访问受保护资源）

	return cookie2 // 返回有效的会话 cookie
}

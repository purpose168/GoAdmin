// manager.go - 管理员管理测试文件
// 包名：common
// 作者：GoAdmin 团队
// 创建日期：2020年
// 描述：本文件包含管理员管理的测试用例，使用 httpexpect 库进行 HTTP 请求测试
//       涵盖管理员列表显示、编辑管理员、新建管理员、管理员登录等核心功能测试

package common

import (
	"fmt"      // 格式化输出包，提供字符串格式化功能
	"net/http" // HTTP 包，提供 HTTP 客户端和服务器功能

	"github.com/gavv/httpexpect"                                   // httpexpect 包，提供 HTTP 请求测试功能
	"github.com/purpose168/GoAdmin/modules/config"                 // 配置模块，提供配置管理功能
	"github.com/purpose168/GoAdmin/modules/errors"                 // 错误模块，提供错误定义
	"github.com/purpose168/GoAdmin/plugins/admin/modules/constant" // 常量模块，提供常量定义
	"github.com/purpose168/GoAdmin/plugins/admin/modules/form"     // 表单模块，提供表单处理功能
)

// managerTest 执行管理员管理测试
// 测试管理员的显示、编辑、新建、登录等功能
//
// 参数:
//   - e: HTTP期望对象，用于发送HTTP请求和验证响应
//   - sesID: 会话ID cookie，用于身份验证
//
// 功能特性:
//   - 显示管理员列表，验证是否包含管理员信息
//   - 编辑管理员（使用错误的 token），验证错误处理
//   - 显示编辑管理员表单，提取表单 token
//   - 编辑管理员（使用正确的 token），验证编辑是否成功
//   - 显示新建管理员表单，提取表单 token
//   - 新建管理员，验证创建是否成功
//   - 测试管理员登录（错误密码），验证错误处理
//   - 测试管理员登录（正确密码），验证登录是否成功
//
// 测试流程:
//  1. 显示管理员列表，验证包含管理员信息
//  2. 编辑管理员（使用错误的 token），验证返回错误
//  3. 显示编辑管理员表单（ID=1），提取 token
//  4. 编辑管理员（ID=1），验证编辑成功
//  5. 显示新建管理员表单，提取 token
//  6. 新建 tester 管理员，验证创建成功
//  7. 测试 tester 管理员登录（错误密码），验证返回错误
//  8. 测试 tester 管理员登录（正确密码），验证登录成功
//
// 说明:
//
//	该函数使用 httpexpect 库进行 HTTP 请求测试，验证管理员管理的各个功能是否正常工作
//	所有请求都携带会话 ID cookie 进行身份验证
//	表单提交使用 multipart/form-data 格式
//	使用正则表达式从表单中提取 CSRF token
//	管理员支持角色权限控制（通过 role_id[] 字段指定角色）
//	管理员支持权限分配（通过 permission_id[] 字段指定权限）
//	管理员支持头像管理（通过 avatar__delete_flag 字段控制头像删除）
//	管理员登录支持密码验证（通过 password 和 password_again 字段）
func managerTest(e *httpexpect.Expect, sesID *http.Cookie) {
	fmt.Println()                               // 打印空行
	printlnWithColor("Manager", "blue")         // 使用蓝色打印 "Manager"
	fmt.Println("============================") // 打印分隔线

	// ==================== 显示管理员列表测试 ====================
	// 测试显示管理员列表，验证是否包含管理员信息
	printlnWithColor("show", "green")   // 使用绿色打印 "show"
	e.GET(config.Url("/info/manager")). // 发送 GET 请求到管理员列表页面
						WithCookie(sesID.Name, sesID.Value).                        // 携带会话 ID cookie
						Expect().                                                   // 获取响应期望对象
						Status(200).                                                // 验证 HTTP 状态码为 200
						Body().Contains("Managers").Contains("admin").Contains("1") // 验证响应体包含 "Managers"、"admin" 和 "1"

	// ==================== 编辑管理员测试（错误的token） ====================
	// 测试编辑管理员（使用错误的 token），验证错误处理
	printlnWithColor("edit", "green")    // 使用绿色打印 "edit"
	e.POST(config.Url("/edit/manager")). // 发送 POST 请求到编辑管理员接口
						WithCookie(sesID.Name, sesID.Value). // 携带会话 ID cookie
						WithMultipart().                     // 使用 multipart/form-data 格式
						WithForm(map[string]interface{}{     // 添加表单字段
			"username":        "admin",                                                                       // 用户名为 "admin"
			"name":            "admin1",                                                                      // 管理员名称为 "admin1"
			"password":        "admin",                                                                       // 密码为 "admin"
			"password_again":  "admin",                                                                       // 确认密码为 "admin"
			"role_id[]":       1,                                                                             // 角色 ID 为 1
			"permission_id[]": 1,                                                                             // 权限 ID 为 1
			form.PreviousKey:  config.Url("/info/manager?__page=1&__pageSize=10&__sort=id&__sort_type=desc"), // 上一个页面 URL
			"id":              "1",                                                                           // 管理员 ID 为 1
			form.TokenKey:     "123",                                                                         // CSRF token（错误的 token）
		}).Expect().Status(200).Body().Contains(errors.EditFailWrongToken) // 验证 HTTP 状态码为 200，响应体包含错误信息

	// ==================== 不带ID的显示表单测试（已注释） ====================
	// 测试不带 ID 显示编辑表单，验证是否返回错误
	//printlnWithColor("show form: without id", "green") // 使用绿色打印 "show form: without id"
	//e.GET(config.Url("/info/manager/edit")). // 发送 GET 请求到编辑管理员表单页面（不带 ID）
	//	WithCookie(sesID.Name, sesID.Value). // 携带会话 ID cookie
	//	Expect().Status(200).Body().Contains(errors.WrongID) // 验证响应体包含错误信息

	// ==================== 显示编辑表单测试 ====================
	// 测试显示编辑管理员表单（ID=1），并从表单中提取 CSRF token
	printlnWithColor("show form", "green")               // 使用绿色打印 "show form"
	formBody := e.GET(config.Url("/info/manager/edit")). // 发送 GET 请求到编辑管理员表单页面
								WithQuery(constant.EditPKKey, "1").  // 添加查询参数，指定编辑的管理员 ID 为 1
								WithCookie(sesID.Name, sesID.Value). // 携带会话 ID cookie
								Expect().Status(200).Body()          // 验证 HTTP 状态码为 200，获取响应体
	token := reg.FindStringSubmatch(formBody.Raw()) // 使用正则表达式从表单中提取 CSRF token

	// ==================== 编辑管理员测试（正确的token） ====================
	// 测试编辑管理员（使用正确的 token），验证编辑是否成功
	printlnWithColor("edit form", "green")      // 使用绿色打印 "edit form"
	res := e.POST(config.Url("/edit/manager")). // 发送 POST 请求到编辑管理员接口
							WithCookie(sesID.Name, sesID.Value). // 携带会话 ID cookie
							WithMultipart().                     // 使用 multipart/form-data 格式
							WithForm(map[string]interface{}{     // 添加表单字段
			"username":            "admin",                                                                       // 用户名为 "admin"
			"name":                "admin1",                                                                      // 管理员名称为 "admin1"
			"password":            "admin",                                                                       // 密码为 "admin"
			"password_again":      "admin",                                                                       // 确认密码为 "admin"
			"avatar__delete_flag": "0",                                                                           // 头像删除标志为 0（不删除头像）
			"role_id[]":           1,                                                                             // 角色 ID 为 1
			"permission_id[]":     1,                                                                             // 权限 ID 为 1
			form.PreviousKey:      config.Url("/info/manager?__page=1&__pageSize=10&__sort=id&__sort_type=desc"), // 上一个页面 URL
			"id":                  "1",                                                                           // 管理员 ID 为 1
			form.TokenKey:         token[1],                                                                      // CSRF token（正确的 token）
		}).Expect().Status(200) // 验证 HTTP 状态码为 200
	res.Header("X-Pjax-Url").Contains(config.Url("/info/")) // 验证响应头 X-Pjax-Url 包含 "/info/"
	res.Body().Contains("admin1")                           // 验证响应体包含 "admin1"

	// ==================== 显示新建表单测试 ====================
	// 测试显示新建管理员表单，并从表单中提取 CSRF token
	printlnWithColor("show new form", "green")         // 使用绿色打印 "show new form"
	formBody = e.GET(config.Url("/info/manager/new")). // 发送 GET 请求到新建管理员表单页面
								WithCookie(sesID.Name, sesID.Value). // 携带会话 ID cookie
								Expect().Status(200).Body()          // 验证 HTTP 状态码为 200，获取响应体
	token = reg.FindStringSubmatch(formBody.Raw()) // 使用正则表达式从表单中提取 CSRF token

	// ==================== 新建管理员测试 ====================
	// 测试新建管理员，验证创建是否成功
	printlnWithColor("new manager test", "green") // 使用绿色打印 "new manager test"
	res = e.POST(config.Url("/new/manager")).     // 发送 POST 请求到新建管理员接口
							WithCookie(sesID.Name, sesID.Value). // 携带会话 ID cookie
							WithMultipart().                     // 使用 multipart/form-data 格式
							WithForm(map[string]interface{}{     // 添加表单字段
			"username":            "tester",                                                                      // 用户名为 "tester"
			"name":                "tester",                                                                      // 管理员名称为 "tester"
			"password":            "tester",                                                                      // 密码为 "tester"
			"password_again":      "tester",                                                                      // 确认密码为 "tester"
			"avatar__delete_flag": "0",                                                                           // 头像删除标志为 0（不删除头像）
			"role_id[]":           1,                                                                             // 角色 ID 为 1
			"permission_id[]":     1,                                                                             // 权限 ID 为 1
			form.PreviousKey:      config.Url("/info/manager?__page=1&__pageSize=10&__sort=id&__sort_type=desc"), // 上一个页面 URL
			"id":                  "1",                                                                           // 管理员 ID 为 1
			form.TokenKey:         token[1],                                                                      // CSRF token
		}).Expect().Status(200) // 验证 HTTP 状态码为 200
	res.Header("X-Pjax-Url").Contains(config.Url("/info/")) // 验证响应头 X-Pjax-Url 包含 "/info/"
	res.Body().Contains("tester")                           // 验证响应体包含 "tester"

	// ==================== 测试管理员登录：错误密码 ====================
	// 测试测试员使用错误密码登录，验证错误处理
	printlnWithColor("tester login: wrong password", "green") // 使用绿色打印 "tester login: wrong password"
	e.POST(config.Url("/signin")).                            // 发送 POST 请求到登录接口
									WithForm(map[string]string{ // 添加表单字段
			"username": "tester", // 用户名为 "tester"
			"password": "admin",  // 密码为 "admin"（错误的密码）
		}).Expect().Status(400) // 验证 HTTP 状态码为 400（错误）

	// ==================== 测试管理员登录成功 ====================
	// 测试测试员使用正确密码登录，验证登录是否成功
	printlnWithColor("tester login success", "green") // 使用绿色打印 "tester login success"
	e.POST(config.Url("/signin")).                    // 发送 POST 请求到登录接口
								WithForm(map[string]string{ // 添加表单字段
			"username": "tester", // 用户名为 "tester"
			"password": "tester", // 密码为 "tester"（正确的密码）
		}).Expect().Status(200).JSON().Equal(map[string]interface{}{ // 验证 HTTP 状态码为 200，解析 JSON 响应
		"code": 200, // 验证 JSON 字段 "code" 等于 200
		"data": map[string]interface{}{ // 验证 JSON 字段 "data"
			"url": "/" + config.GetUrlPrefix(), // 验证 URL 字段包含 URL 前缀
		},
		"msg": "ok", // 验证 JSON 字段 "msg" 等于 "ok"
	})
}

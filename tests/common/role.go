// role.go - 角色管理测试文件
// 包名：common
// 作者：GoAdmin 团队
// 创建日期：2020年
// 描述：本文件包含角色管理的测试用例，使用 httpexpect 库进行 HTTP 请求测试
//       涵盖角色列表显示、新建角色、编辑角色、删除角色等核心功能测试

package common

import (
	"fmt"      // 格式化输出包，提供字符串格式化功能
	"net/http" // HTTP 包，提供 HTTP 客户端和服务器功能

	"github.com/gavv/httpexpect"                                   // httpexpect 包，提供 HTTP 请求测试功能
	"github.com/purpose168/GoAdmin/modules/config"                 // 配置模块，提供配置管理功能
	"github.com/purpose168/GoAdmin/plugins/admin/modules/constant" // 常量模块，提供常量定义
	"github.com/purpose168/GoAdmin/plugins/admin/modules/form"     // 表单模块，提供表单处理功能
)

// roleTest 执行角色管理测试
// 测试角色的显示、新建、编辑、删除等功能
//
// 参数:
//   - e: HTTP期望对象，用于发送HTTP请求和验证响应
//   - sesID: 会话ID cookie，用于身份验证
//
// 功能特性:
//   - 显示角色列表，验证是否包含 Administrator 和 Operator 角色
//   - 显示新建角色表单，提取表单 token
//   - 新建角色，验证角色是否创建成功
//   - 显示编辑角色表单，提取表单 token
//   - 编辑角色，验证角色是否更新成功
//   - 新建 tester2 角色
//   - 删除 tester2 角色，验证删除是否成功
//
// 测试流程:
//  1. 显示角色列表，验证包含 Administrator 和 Operator
//  2. 显示新建角色表单，提取 token
//  3. 新建 tester 角色，验证创建成功
//  4. 显示编辑角色表单（ID=3），提取 token
//  5. 编辑角色（ID=3），验证更新成功
//  6. 显示新建角色表单，提取 token
//  7. 新建 tester2 角色
//  8. 删除 tester2 角色，验证删除成功
//
// 说明:
//
//	该函数使用 httpexpect 库进行 HTTP 请求测试，验证角色管理的各个功能是否正常工作
//	所有请求都携带会话 ID cookie 进行身份验证
//	表单提交使用 multipart/form-data 格式
//	使用正则表达式从表单中提取 CSRF token
func roleTest(e *httpexpect.Expect, sesID *http.Cookie) {
	fmt.Println()                               // 打印空行
	printlnWithColor("Role", "blue")            // 使用蓝色打印 "Role"
	fmt.Println("============================") // 打印分隔线

	// ==================== 显示角色列表测试 ====================
	// 测试显示角色列表，验证是否包含 Administrator 和 Operator 角色
	printlnWithColor("show", "green") // 使用绿色打印 "show"
	e.GET(config.Url("/info/roles")). // 发送 GET 请求到角色列表页面
						WithCookie(sesID.Name, sesID.Value).                  // 携带会话 ID cookie
						Expect().                                             // 获取响应期望对象
						Status(200).                                          // 验证 HTTP 状态码为 200
						Body().Contains("Administrator").Contains("Operator") // 验证响应体包含 "Administrator" 和 "Operator"

	// ==================== 显示新建表单测试 ====================
	// 测试显示新建角色表单，并从表单中提取 CSRF token
	printlnWithColor("show new form", "green")        // 使用绿色打印 "show new form"
	formBody := e.GET(config.Url("/info/roles/new")). // 发送 GET 请求到新建角色表单页面
								WithCookie(sesID.Name, sesID.Value). // 携带会话 ID cookie
								Expect().Status(200).Body()          // 验证 HTTP 状态码为 200，获取响应体
	token := reg.FindStringSubmatch(formBody.Raw()) // 使用正则表达式从表单中提取 CSRF token

	// ==================== 新建角色测试 ====================
	// 测试新建角色，验证角色是否创建成功
	printlnWithColor("new roles test", "green") // 使用绿色打印 "new roles test"
	res := e.POST(config.Url("/new/roles")).    // 发送 POST 请求到新建角色接口
							WithCookie(sesID.Name, sesID.Value).   // 携带会话 ID cookie
							WithMultipart().                       // 使用 multipart/form-data 格式
							WithFormField("permission_id[]", "3"). // 添加权限 ID 字段（权限 ID 为 3）
							WithForm(map[string]interface{}{       // 添加表单字段
			"name":           "tester",                                                                    // 角色名称为 "tester"
			"slug":           "tester",                                                                    // 角色标识为 "tester"
			form.PreviousKey: config.Url("/info/roles?__page=1&__pageSize=10&__sort=id&__sort_type=desc"), // 上一个页面 URL
			form.TokenKey:    token[1],                                                                    // CSRF token
		}).Expect().Status(200) // 验证 HTTP 状态码为 200
	res.Header("X-Pjax-Url").Contains(config.Url("/info/")) // 验证响应头 X-Pjax-Url 包含 "/info/"
	res.Body().Contains("tester")                           // 验证响应体包含 "tester"

	// ==================== 不带ID的显示表单测试（已注释） ====================
	// 测试不带 ID 显示编辑表单，验证是否返回错误
	//printlnWithColor("show form: without id", "green") // 使用绿色打印 "show form: without id"
	//e.GET(config.Url("/info/roles/edit")). // 发送 GET 请求到编辑角色表单页面（不带 ID）
	//	WithCookie(sesID.Name, sesID.Value). // 携带会话 ID cookie
	//	Expect().Status(200).Body().Contains(errors.WrongID) // 验证响应体包含错误信息

	// ==================== 显示编辑表单测试 ====================
	// 测试显示编辑角色表单（ID=3），并从表单中提取 CSRF token
	printlnWithColor("show form", "green")            // 使用绿色打印 "show form"
	formBody = e.GET(config.Url("/info/roles/edit")). // 发送 GET 请求到编辑角色表单页面
								WithQuery(constant.EditPKKey, "3").  // 添加查询参数，指定编辑的角色 ID 为 3
								WithCookie(sesID.Name, sesID.Value). // 携带会话 ID cookie
								Expect().Status(200).Body()          // 验证 HTTP 状态码为 200，获取响应体
	token = reg.FindStringSubmatch(formBody.Raw()) // 使用正则表达式从表单中提取 CSRF token

	// ==================== 编辑角色测试 ====================
	// 测试编辑角色（ID=3），验证角色是否更新成功
	printlnWithColor("edit form", "green")   // 使用绿色打印 "edit form"
	res = e.POST(config.Url("/edit/roles")). // 发送 POST 请求到编辑角色接口
							WithCookie(sesID.Name, sesID.Value).   // 携带会话 ID cookie
							WithMultipart().                       // 使用 multipart/form-data 格式
							WithFormField("permission_id[]", "3"). // 添加权限 ID 字段（权限 ID 为 3）
							WithFormField("permission_id[]", "2"). // 添加权限 ID 字段（权限 ID 为 2）
							WithForm(map[string]interface{}{       // 添加表单字段
			"name":           "tester",                                                                    // 角色名称为 "tester"
			"slug":           "tester",                                                                    // 角色标识为 "tester"
			form.PreviousKey: config.Url("/info/roles?__page=1&__pageSize=10&__sort=id&__sort_type=desc"), // 上一个页面 URL
			form.TokenKey:    token[1],                                                                    // CSRF token
			"id":             "3",                                                                         // 角色 ID 为 3
		}).Expect().Status(200) // 验证 HTTP 状态码为 200
	res.Header("X-Pjax-Url").Contains(config.Url("/info/")) // 验证响应头 X-Pjax-Url 包含 "/info/"
	res.Body().Contains("tester")                           // 验证响应体包含 "tester"

	// ==================== 显示新建表单测试 ====================
	// 测试显示新建角色表单，并从表单中提取 CSRF token
	printlnWithColor("show new form", "green")       // 使用绿色打印 "show new form"
	formBody = e.GET(config.Url("/info/roles/new")). // 发送 GET 请求到新建角色表单页面
								WithCookie(sesID.Name, sesID.Value). // 携带会话 ID cookie
								Expect().Status(200).Body()          // 验证 HTTP 状态码为 200，获取响应体
	token = reg.FindStringSubmatch(formBody.Raw()) // 使用正则表达式从表单中提取 CSRF token

	// ==================== 新建 tester2 角色测试 ====================
	// 测试新建 tester2 角色
	printlnWithColor("new tester2", "green") // 使用绿色打印 "new tester2"
	e.POST(config.Url("/new/roles")).        // 发送 POST 请求到新建角色接口
							WithCookie(sesID.Name, sesID.Value).   // 携带会话 ID cookie
							WithMultipart().                       // 使用 multipart/form-data 格式
							WithFormField("permission_id[]", "3"). // 添加权限 ID 字段（权限 ID 为 3）
							WithForm(map[string]interface{}{       // 添加表单字段
			"name":           "tester2",                                                                   // 角色名称为 "tester2"
			"slug":           "tester2",                                                                   // 角色标识为 "tester2"
			form.PreviousKey: config.Url("/info/roles?__page=1&__pageSize=10&__sort=id&__sort_type=desc"), // 上一个页面 URL
			form.TokenKey:    token[1],                                                                    // CSRF token
		}).Expect().Status(200) // 验证 HTTP 状态码为 200

	// ==================== 删除 tester2 角色测试 ====================
	// 测试删除 tester2 角色，验证删除是否成功
	printlnWithColor("delete roles tester2", "green") // 使用绿色打印 "delete roles tester2"
	e.POST(config.Url("/delete/roles")).              // 发送 POST 请求到删除角色接口
								WithCookie(sesID.Name, sesID.Value).  // 携带会话 ID cookie
								WithMultipart().                      // 使用 multipart/form-data 格式
								WithFormField("id", "3").             // 添加 ID 字段（角色 ID 为 3）
								Expect().Status(200).JSON().Object(). // 验证 HTTP 状态码为 200，解析 JSON 响应
								ValueEqual("code", 200).              // 验证 JSON 字段 "code" 等于 200
								ValueEqual("msg", "ok")               // 验证 JSON 字段 "msg" 等于 "ok"
}

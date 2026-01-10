// menu.go - 菜单管理测试文件
// 包名：common
// 作者：GoAdmin 团队
// 创建日期：2020年
// 描述：本文件包含菜单管理的测试用例，使用 httpexpect 库进行 HTTP 请求测试
//       涵盖菜单显示、新建菜单、编辑菜单、删除菜单等核心功能测试

package common

import (
	"fmt"      // 格式化输出包，提供字符串格式化功能
	"net/http" // HTTP 包，提供 HTTP 客户端和服务器功能

	"github.com/gavv/httpexpect"                                   // httpexpect 包，提供 HTTP 请求测试功能
	"github.com/purpose168/GoAdmin/modules/config"                 // 配置模块，提供配置管理功能
	"github.com/purpose168/GoAdmin/modules/language"               // 语言模块，提供多语言支持
	"github.com/purpose168/GoAdmin/plugins/admin/modules/constant" // 常量模块，提供常量定义
	"github.com/purpose168/GoAdmin/plugins/admin/modules/form"     // 表单模块，提供表单处理功能
)

// menuTest 执行菜单管理测试
// 测试菜单的显示、新建、编辑、删除等功能
//
// 参数:
//   - e: HTTP期望对象，用于发送HTTP请求和验证响应
//   - sesID: 会话ID cookie，用于身份验证
//
// 功能特性:
//   - 显示菜单列表，验证是否包含菜单管理标题
//   - 显示新建菜单表单，提取表单 token
//   - 新建菜单，验证菜单是否创建成功
//   - 显示编辑菜单表单，提取表单 token
//   - 编辑菜单，验证菜单是否更新成功
//   - 新建 tester2 菜单
//   - 删除 tester2 菜单，验证删除是否成功
//
// 测试流程:
//  1. 显示菜单列表，验证包含菜单管理标题
//  2. 显示新建菜单表单，提取 token
//  3. 新建 test menu，验证创建成功
//  4. 显示编辑菜单表单（ID=3），提取 token
//  5. 编辑菜单（ID=3），验证更新成功
//  6. 新建 test2 menu
//  7. 删除 test2 menu，验证删除成功
//
// 说明:
//
//	该函数使用 httpexpect 库进行 HTTP 请求测试，验证菜单管理的各个功能是否正常工作
//	所有请求都携带会话 ID cookie 进行身份验证
//	表单提交使用 multipart/form-data 格式
//	使用正则表达式从表单中提取 CSRF token
//	菜单支持层级结构（通过 parent_id 指定父级菜单）
//	菜单支持角色权限控制（通过 roles[] 字段指定可访问的角色）
//	菜单支持图标配置（使用 Font Awesome 图标）
//	菜单支持多语言显示（通过 language.Get 获取本地化文本）
func menuTest(e *httpexpect.Expect, sesID *http.Cookie) {
	fmt.Println()                               // 打印空行
	printlnWithColor("Menu", "blue")            // 使用蓝色打印 "Menu"
	fmt.Println("============================") // 打印分隔线

	// ==================== 显示菜单测试 ====================
	// 测试显示菜单列表，验证是否包含菜单管理标题
	printlnWithColor("show", "green")       // 使用绿色打印 "show"
	formBody := e.GET(config.Url("/menu")). // 发送 GET 请求到菜单管理页面
						WithCookie(sesID.Name, sesID.Value).          // 携带会话 ID cookie
						Expect().                                     // 获取响应期望对象
						Status(200).                                  // 验证 HTTP 状态码为 200
						Body().Contains(language.Get("menus manage")) // 验证响应体包含菜单管理标题（使用多语言）
	token := reg.FindStringSubmatch(formBody.Raw()) // 使用正则表达式从表单中提取 CSRF token

	// ==================== 新建菜单测试 ====================
	// 测试新建菜单，验证菜单是否创建成功
	printlnWithColor("new menu test", "green") // 使用绿色打印 "new menu test"
	res := e.POST(config.Url("/menu/new")).    // 发送 POST 请求到新建菜单接口
							WithCookie(sesID.Name, sesID.Value). // 携带会话 ID cookie
							WithMultipart().                     // 使用 multipart/form-data 格式
							WithFormField("roles[]", "1").       // 添加角色字段（角色 ID 为 1，表示该菜单对角色 1 可见）
							WithForm(map[string]interface{}{     // 添加表单字段
			"parent_id":      0,               // 父级菜单 ID 为 0，表示顶级菜单
			"title":          "test menu",     // 菜单标题为 "test menu"
			"header":         "",              // 菜单头部为空
			"icon":           "fa-angellist",  // 菜单图标为 Font Awesome 的 fa-angellist 图标
			"uri":            "/example/test", // 菜单 URI 为 "/example/test"
			form.PreviousKey: "/admin/menu",   // 上一个页面 URL 为 "/admin/menu"
			form.TokenKey:    token[1],        // CSRF token
		}).Expect().Status(200) // 验证 HTTP 状态码为 200
	res.Header("X-Pjax-Url").Contains(config.Url("/menu"))     // 验证响应头 X-Pjax-Url 包含 "/menu"
	res.Body().Contains("test menu").Contains("/example/test") // 验证响应体包含 "test menu" 和 "/example/test"

	// ==================== 不带ID的显示表单测试（已注释） ====================
	// 测试不带 ID 显示编辑表单，验证是否返回错误
	//printlnWithColor("show form: without id", "green") // 使用绿色打印 "show form: without id"
	//e.GET(config.Url("/menu/edit/show")). // 发送 GET 请求到编辑菜单表单页面（不带 ID）
	//	WithCookie(sesID.Name, sesID.Value). // 携带会话 ID cookie
	//	Expect().Status(200).Body().Contains(errors.WrongID) // 验证响应体包含错误信息

	// ==================== 显示编辑表单测试 ====================
	// 测试显示编辑菜单表单（ID=3），并从表单中提取 CSRF token
	printlnWithColor("show form", "green")           // 使用绿色打印 "show form"
	formBody = e.GET(config.Url("/menu/edit/show")). // 发送 GET 请求到编辑菜单表单页面
								WithQuery(constant.EditPKKey, "3").  // 添加查询参数，指定编辑的菜单 ID 为 3
								WithCookie(sesID.Name, sesID.Value). // 携带会话 ID cookie
								Expect().Status(200).Body()          // 验证 HTTP 状态码为 200，获取响应体
	token = reg.FindStringSubmatch(formBody.Raw()) // 使用正则表达式从表单中提取 CSRF token

	// ==================== 编辑菜单测试 ====================
	// 测试编辑菜单（ID=3），验证菜单是否更新成功
	printlnWithColor("edit form", "green")  // 使用绿色打印 "edit form"
	res = e.POST(config.Url("/menu/edit")). // 发送 POST 请求到编辑菜单接口
						WithCookie(sesID.Name, sesID.Value). // 携带会话 ID cookie
						WithMultipart().                     // 使用 multipart/form-data 格式
						WithFormField("roles[]", "1").       // 添加角色字段（角色 ID 为 1，表示该菜单对角色 1 可见）
						WithForm(map[string]interface{}{     // 添加表单字段
			"parent_id":      0,               // 父级菜单 ID 为 0，表示顶级菜单
			"title":          "test2 menu",    // 菜单标题为 "test2 menu"
			"header":         "",              // 菜单头部为空
			"icon":           "fa-angellist",  // 菜单图标为 Font Awesome 的 fa-angellist 图标
			"uri":            "/example/test", // 菜单 URI 为 "/example/test"
			form.PreviousKey: "/admin/menu",   // 上一个页面 URL 为 "/admin/menu"
			form.TokenKey:    token[1],        // CSRF token
			"id":             "3",             // 菜单 ID 为 3
		}).Expect().Status(200) // 验证 HTTP 状态码为 200
	res.Header("X-Pjax-Url").Contains(config.Url("/menu"))      // 验证响应头 X-Pjax-Url 包含 "/menu"
	res.Body().Contains("test2 menu").Contains("/example/test") // 验证响应体包含 "test2 menu" 和 "/example/test"

	// ==================== 新建 tester2 菜单测试 ====================
	// 测试新建 tester2 菜单
	printlnWithColor("new tester2", "green") // 使用绿色打印 "new tester2"
	e.POST(config.Url("/menu/new")).         // 发送 POST 请求到新建菜单接口
							WithCookie(sesID.Name, sesID.Value). // 携带会话 ID cookie
							WithMultipart().                     // 使用 multipart/form-data 格式
							WithFormField("roles[]", "1").       // 添加角色字段（角色 ID 为 1，表示该菜单对角色 1 可见）
							WithForm(map[string]interface{}{     // 添加表单字段
			"parent_id":      0,                // 父级菜单 ID 为 0，表示顶级菜单
			"title":          "test2 menu",     // 菜单标题为 "test2 menu"
			"header":         "",               // 菜单头部为空
			"icon":           "fa-angellist",   // 菜单图标为 Font Awesome 的 fa-angellist 图标
			"uri":            "/example/test2", // 菜单 URI 为 "/example/test2"
			form.PreviousKey: "/admin/menu",    // 上一个页面 URL 为 "/admin/menu"
			form.TokenKey:    token[1],         // CSRF token
		}).Expect().Status(200) // 验证 HTTP 状态码为 200

	// ==================== 删除 tester2 菜单测试 ====================
	// 测试删除 tester2 菜单，验证删除是否成功
	printlnWithColor("delete menu tester2", "green") // 使用绿色打印 "delete menu tester2"
	e.POST(config.Url("/menu/delete")).              // 发送 POST 请求到删除菜单接口
								WithQuery("id", "9").                 // 添加查询参数，指定删除的菜单 ID 为 9
								WithCookie(sesID.Name, sesID.Value).  // 携带会话 ID cookie
								WithMultipart().                      // 使用 multipart/form-data 格式
								Expect().Status(200).JSON().Object(). // 验证 HTTP 状态码为 200，解析 JSON 响应
								ValueEqual("code", 200).              // 验证 JSON 字段 "code" 等于 200
								ValueEqual("msg", "delete succeed")   // 验证 JSON 字段 "msg" 等于 "delete succeed"
}

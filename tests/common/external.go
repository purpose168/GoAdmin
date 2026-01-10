// external.go - 外部数据源测试文件
// 包名：common
// 作者：GoAdmin 团队
// 创建日期：2020年
// 描述：本文件包含外部数据源的测试用例，使用 httpexpect 库进行 HTTP 请求测试
//       涵盖外部数据源信息显示、编辑表单显示、新建表单显示等核心功能测试

package common

import (
	"fmt"      // 格式化输出包，提供字符串格式化功能
	"net/http" // HTTP 包，提供 HTTP 客户端和服务器功能

	"github.com/gavv/httpexpect"                                   // httpexpect 包，提供 HTTP 请求测试功能
	"github.com/purpose168/GoAdmin/modules/config"                 // 配置模块，提供配置管理功能
	"github.com/purpose168/GoAdmin/plugins/admin/modules/constant" // 常量模块，提供常量定义
)

// externalTest 执行外部数据源测试
// 测试从外部数据源获取信息、显示表单等功能
//
// 参数:
//   - e: HTTP期望对象，用于发送HTTP请求和验证响应
//   - sesID: 会话ID cookie，用于身份验证
//
// 功能特性:
//   - 显示外部数据源信息，验证是否包含预期内容
//   - 显示编辑外部数据源表单，验证表单是否正确显示
//   - 显示新建外部数据源表单，验证表单是否正确显示
//
// 测试流程:
//  1. 显示外部数据源信息，验证包含 "External"、"this is a title" 和 "10"
//  2. 显示编辑外部数据源表单（ID=10），验证表单正确显示
//  3. 显示新建外部数据源表单，验证表单正确显示
//
// 说明:
//
//	该函数使用 httpexpect 库进行 HTTP 请求测试，验证外部数据源管理的各个功能是否正常工作
//	所有请求都携带会话 ID cookie 进行身份验证
//	外部数据源是指从外部 API 或数据源获取数据，而不是直接从数据库查询
//	外部数据源支持自定义数据获取逻辑，可以集成第三方数据源
//	外部数据源支持编辑和新建操作，但数据实际存储在外部系统中
func externalTest(e *httpexpect.Expect, sesID *http.Cookie) {
	fmt.Println()                               // 打印空行
	printlnWithColor("External", "blue")        // 使用蓝色打印 "External"
	fmt.Println("============================") // 打印分隔线

	// ==================== 显示外部数据源信息测试 ====================
	// 测试显示外部数据源信息，验证是否包含预期内容
	printlnWithColor("show", "green")    // 使用绿色打印 "show"
	e.GET(config.Url("/info/external")). // 发送 GET 请求到外部数据源信息页面
						WithCookie(sesID.Name, sesID.Value).                                   // 携带会话 ID cookie
						Expect().                                                              // 获取响应期望对象
						Status(200).                                                           // 验证 HTTP 状态码为 200
						Body().Contains("External").Contains("this is a title").Contains("10") // 验证响应体包含 "External"、"this is a title" 和 "10"

	// ==================== 不带ID的显示表单测试（已注释） ====================
	// 测试不带 ID 显示编辑表单，验证是否返回错误
	//printlnWithColor("show form: without id", "green") // 使用绿色打印 "show form: without id"
	//e.GET(config.Url("/info/external/edit")). // 发送 GET 请求到编辑外部数据源表单页面（不带 ID）
	//	WithCookie(sesID.Name, sesID.Value). // 携带会话 ID cookie
	//	Expect().Status(200).Body().Contains(errors.WrongID) // 验证响应体包含错误信息

	// ==================== 显示编辑表单测试 ====================
	// 测试显示编辑外部数据源表单（ID=10），验证表单是否正确显示
	printlnWithColor("show form", "green")    // 使用绿色打印 "show form"
	e.GET(config.Url("/info/external/edit")). // 发送 GET 请求到编辑外部数据源表单页面
							WithQuery(constant.EditPKKey, "10"). // 添加查询参数，指定编辑的外部数据源 ID 为 10
							WithCookie(sesID.Name, sesID.Value). // 携带会话 ID cookie
							Expect().Status(200).Body()          // 验证 HTTP 状态码为 200，获取响应体

	// ==================== 显示新建表单测试 ====================
	// 测试显示新建外部数据源表单，验证表单是否正确显示
	printlnWithColor("show new form", "green") // 使用绿色打印 "show new form"
	e.GET(config.Url("/info/external/new")).   // 发送 GET 请求到新建外部数据源表单页面
							WithCookie(sesID.Name, sesID.Value). // 携带会话 ID cookie
							Expect().Status(200).Body()          // 验证 HTTP 状态码为 200，获取响应体
}

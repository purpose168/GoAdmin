// Package common 提供GoAdmin框架的通用测试函数
//
// 文件名: api.go
// 包名: common
// 作者: GoAdmin Team
// 创建日期: 2026-01-10
//
// 功能描述:
// 本文件提供了API接口测试功能，用于验证GoAdmin框架的API接口是否正常工作
// 主要测试管理员相关的API接口，包括列表查询、表单获取等功能
//
// 核心概念:
// - API测试: 使用httpexpect库对GoAdmin框架的API接口进行自动化测试
// - JSON响应: API接口返回JSON格式的响应数据，包含状态码和业务数据
// - HTTP请求: 使用GET/POST等HTTP方法访问API接口
// - 请求头: 设置Accept等请求头，指定响应数据格式
// - 查询参数: 通过URL查询参数传递请求参数（如编辑ID）
// - Cookie认证: 通过会话cookie进行身份验证
//
// 技术栈:
// - httpexpect: HTTP请求测试库，提供链式API进行HTTP请求和响应验证
// - config: GoAdmin配置模块，提供URL生成和配置管理功能
// - constant: GoAdmin常量模块，提供框架常量定义
//
// 测试范围:
// - 管理员列表API: 测试获取管理员列表数据
// - 更新表单API: 测试获取编辑管理员的表单数据
// - 创建表单API: 测试获取新建管理员的表单数据
//
// 使用场景:
// - 集成测试: 在集成测试中验证API接口功能
// - 回归测试: 在代码变更后验证API接口是否正常
// - CI/CD: 在持续集成流程中自动运行API测试
//
// 注意事项:
// - API接口需要身份验证，必须携带有效的会话cookie
// - API接口返回JSON格式数据，包含code、data、msg等字段
// - code=200表示请求成功，code=400表示请求参数错误
// - 编辑表单需要提供编辑对象的ID（通过constant.EditPKKey参数）
// - Accept请求头设置为"application/json, text/plain, */*"，表示接受JSON格式响应
package common

import (
	"fmt"      // fmt: 格式化输出包，提供字符串格式化功能
	"net/http" // net/http: HTTP包，提供HTTP客户端和服务器功能

	"github.com/gavv/httpexpect"                                   // httpexpect: HTTP请求测试库，提供链式API进行HTTP请求和响应验证
	"github.com/purpose168/GoAdmin/modules/config"                 // config: GoAdmin配置模块，提供URL生成和配置管理功能
	"github.com/purpose168/GoAdmin/plugins/admin/modules/constant" // constant: GoAdmin常量模块，提供框架常量定义（如EditPKKey）
)

// apiTest 执行API测试
// 测试管理员相关的API接口，验证API接口的响应数据和状态码是否正确
//
// 参数:
//   - e: HTTP期望对象，用于发送HTTP请求和验证响应
//   - sesID: 会话ID cookie，用于身份验证
//
// 返回值:
//   - 无返回值
//
// 功能特性:
//   - 测试管理员列表API，验证能否正确获取管理员列表数据
//   - 测试更新表单API，验证能否正确获取编辑管理员的表单数据
//   - 测试创建表单API，验证能否正确获取新建管理员的表单数据
//   - 验证API接口的HTTP状态码（200表示成功，400表示参数错误）
//   - 验证API接口的JSON响应格式（code字段表示请求状态）
//   - 支持JSON格式响应，便于前端集成和数据处理
//   - 支持查询参数传递（如编辑ID）
//   - 支持请求头设置（如Accept指定响应格式）
//
// 测试流程:
//  1. 测试管理员列表API，验证返回状态码为200，code为200
//  2. 测试更新表单API（ID=1），验证返回状态码为200，code为200
//  3. 测试创建表单API，验证返回状态码为200，code为200
//
// 说明:
//
//	该函数使用 httpexpect 库进行 HTTP 请求测试，验证 API 接口的各个功能是否正常工作
//	所有请求都携带会话 ID cookie 进行身份验证
//	所有请求都设置 Accept 请求头为 "application/json, text/plain, */*"，表示接受 JSON 格式响应
//	API 接口返回 JSON 格式数据，包含 code、data、msg 等字段
//	code=200 表示请求成功，code=400 表示请求参数错误
//	编辑表单 API 需要提供编辑对象的 ID（通过 constant.EditPKKey 参数）
//	API 接口遵循 RESTful 风格，使用 GET 方法获取数据
//	API 接口支持分页、排序、筛选等功能（通过 URL 查询参数）
//	API 接口支持跨域访问（通过 CORS 头）
//	API 接口支持数据验证（通过参数校验）
//	API 接口支持错误处理（通过错误码和错误消息）
//
// 示例:
//
//	// 创建 HTTP 期望对象
//	e := httpexpect.Default(t, "http://localhost:9033")
//
//	// 登录获取会话 cookie
//	sesID := loginTest(e)
//
//	// 执行 API 测试
//	apiTest(e, sesID)
//
// 注意事项:
//
//   - API 接口需要身份验证，必须携带有效的会话 cookie
//   - API 接口返回 JSON 格式数据，包含 code、data、msg 等字段
//   - code=200 表示请求成功，code=400 表示请求参数错误
//   - 编辑表单需要提供编辑对象的 ID（通过 constant.EditPKKey 参数）
//   - Accept 请求头设置为 "application/json, text/plain, */*"，表示接受 JSON 格式响应
//   - API 接口遵循 RESTful 风格，使用 GET 方法获取数据
//   - API 接口支持分页、排序、筛选等功能（通过 URL 查询参数）
//   - API 接口支持跨域访问（通过 CORS 头）
//   - API 接口支持数据验证（通过参数校验）
//   - API 接口支持错误处理（通过错误码和错误消息）
func apiTest(e *httpexpect.Expect, sesID *http.Cookie) {

	fmt.Println()                               // 打印空行
	printlnWithColor("Api", "blue")             // 使用蓝色打印 "Api"
	fmt.Println("============================") // 打印分隔线

	// ==================== 测试管理员列表API ====================
	// 测试获取管理员列表数据，验证API接口是否正常工作
	printlnWithColor("show", "green")       // 使用绿色打印 "show"
	e.GET(config.Url("/api/list/manager")). // 发送 GET 请求到管理员列表API
						WithHeader("Accept", "application/json, text/plain, */*"). // 设置 Accept 请求头，接受 JSON 格式响应
						WithCookie(sesID.Name, sesID.Value).                       // 携带会话 ID cookie 进行身份验证
						Expect().                                                  // 获取响应期望对象
						Status(200).                                               // 验证 HTTP 状态码为 200（请求成功）
						JSON().                                                    // 解析 JSON 响应
						Object().                                                  // 获取 JSON 对象
						ValueEqual("code", 200)                                    // 验证 JSON 字段 "code" 等于 200（业务成功）

	// ==================== 测试更新表单API（已注释） ====================
	// 测试不带ID的更新表单API，验证参数错误处理（已注释）
	// 说明: 编辑表单必须提供编辑对象的ID，否则返回400错误
	//printlnWithColor("update form without id", "green") // 使用绿色打印 "update form without id"
	//e.GET(config.Url("/api/edit/form/manager")). // 发送 GET 请求到编辑表单API（不带ID）
	//	WithHeader("Accept", "application/json, text/plain, */*"). // 设置 Accept 请求头
	//	WithCookie(sesID.Name, sesID.Value). // 携带会话 ID cookie
	//	Expect(). // 获取响应期望对象
	//	Status(400). // 验证 HTTP 状态码为 400（参数错误）
	//	JSON(). // 解析 JSON 响应
	//	Object(). // 获取 JSON 对象
	//	ValueEqual("code", 400) // 验证 JSON 字段 "code" 等于 400（参数错误）

	// ==================== 测试更新表单API ====================
	// 测试获取编辑管理员的表单数据，验证API接口是否正常工作
	printlnWithColor("update form", "green")     // 使用绿色打印 "update form"
	e.GET(config.Url("/api/edit/form/manager")). // 发送 GET 请求到编辑表单API
							WithHeader("Accept", "application/json, text/plain, */*"). // 设置 Accept 请求头，接受 JSON 格式响应
							WithQuery(constant.EditPKKey, "1").                        // 添加查询参数，指定编辑的管理员 ID 为 1（constant.EditPKKey 是编辑主键的常量）
							WithCookie(sesID.Name, sesID.Value).                       // 携带会话 ID cookie 进行身份验证
							Expect().                                                  // 获取响应期望对象
							Status(200).                                               // 验证 HTTP 状态码为 200（请求成功）
							JSON().                                                    // 解析 JSON 响应
							Object().                                                  // 获取 JSON 对象
							ValueEqual("code", 200)                                    // 验证 JSON 字段 "code" 等于 200（业务成功）

	// ==================== 测试创建表单API ====================
	// 测试获取新建管理员的表单数据，验证API接口是否正常工作
	printlnWithColor("create form", "green")       // 使用绿色打印 "create form"
	e.GET(config.Url("/api/create/form/manager")). // 发送 GET 请求到创建表单API
							WithHeader("Accept", "application/json, text/plain, */*"). // 设置 Accept 请求头，接受 JSON 格式响应
							WithCookie(sesID.Name, sesID.Value).                       // 携带会话 ID cookie 进行身份验证
							Expect().                                                  // 获取响应期望对象
							Status(200).                                               // 验证 HTTP 状态码为 200（请求成功）
							JSON().                                                    // 解析 JSON 响应
							Object().                                                  // 获取 JSON 对象
							ValueEqual("code", 200)                                    // 验证 JSON 字段 "code" 等于 200（业务成功）
}

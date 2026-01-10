package common

import (
	"fmt"
	"net/http"

	"github.com/gavv/httpexpect"
	"github.com/purpose168/GoAdmin/modules/config"
)

// normalTest 执行正常表格测试
// 测试用户表的显示、导出、显示表单等功能
// 参数:
//   - e: HTTP期望对象，用于发送HTTP请求和验证响应
//   - sesID: 会话ID cookie，用于身份验证
func normalTest(e *httpexpect.Expect, sesID *http.Cookie) {

	fmt.Println()
	printlnWithColor("Normal Table", "blue")
	fmt.Println("============================")

	// 显示用户列表测试

	// 测试显示用户列表
	printlnWithColor("show", "green")
	e.GET(config.Url("/info/user")).
		WithCookie(sesID.Name, sesID.Value).
		Expect().
		Status(200).
		Body().Contains("Users")

	// 导出测试

	// 测试导出用户数据
	printlnWithColor("export test", "green")
	e.POST(config.Url("/export/user")).
		WithCookie(sesID.Name, sesID.Value).
		WithMultipart().
		WithFormField("id", "1").
		Expect().Status(200)

	// 不带ID的显示表单测试（已注释）

	//printlnWithColor("show form: without id", "green")
	//e.GET(config.Url("/info/user/edit")).
	//	WithCookie(sesID.Name, sesID.Value).
	//	Expect().Status(200).Body().Contains(errors.WrongID)

	// 显示编辑表单测试（已注释）

	//printlnWithColor("show form", "green")
	//e.GET(config.Url("/info/user/edit")).
	//	WithQuery(constant.EditPKKey, "362").
	//	WithCookie(sesID.Name, sesID.Value).
	//	Expect().Status(200).Body()

	// 显示新建表单测试

	// 测试显示新建用户表单
	printlnWithColor("show new form", "green")
	e.GET(config.Url("/info/user/new")).
		WithCookie(sesID.Name, sesID.Value).
		Expect().Status(200).Body()
}

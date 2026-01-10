package common

import (
	"fmt"
	"net/http"

	"github.com/gavv/httpexpect"
	"github.com/purpose168/GoAdmin/modules/config"
	"github.com/purpose168/GoAdmin/modules/language"
)

// operationLogTest 执行操作日志测试
// 测试操作日志的显示功能
// 参数:
//   - e: HTTP期望对象，用于发送HTTP请求和验证响应
//   - sesID: 会话ID cookie，用于身份验证
func operationLogTest(e *httpexpect.Expect, sesID *http.Cookie) {

	fmt.Println()
	printlnWithColor("Operation Log", "blue")
	fmt.Println("============================")

	// 显示操作日志测试

	// 测试显示操作日志
	printlnWithColor("show", "green")
	e.GET(config.Url("/info/op")).
		WithCookie(sesID.Name, sesID.Value).
		Expect().
		Status(200).
		Body().Contains(language.Get("operation log"))
}

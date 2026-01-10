package gofiber

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/purpose168/GoAdmin/tests/common"
)

// TestGofiber 执行 GoFiber 框架测试
// 参数:
//   - t: 测试对象，用于报告测试结果
//
// 功能:
//   - 创建HTTP客户端配置
//   - 使用内部处理器作为传输层（FastHTTP绑定器）
//   - 创建新的Cookie Jar
//   - 执行额外的通用测试
func TestGofiber(t *testing.T) {
	common.ExtraTest(httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewFastBinder(internalHandler()), // 使用内部处理器作为传输层
			Jar:       httpexpect.NewJar(),                         // 创建新的Cookie Jar
		},
		Reporter: httpexpect.NewAssertReporter(t), // 创建断言报告器
	}))
}

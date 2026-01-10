// beego_test.go - Beego 框架测试文件
// 包名：beego
// 作者：GoAdmin 团队
// 创建日期：2020年
// 描述：本文件提供 GoAdmin 在 Beego Web 框架下的集成测试
//       测试使用 httpexpect 库进行 HTTP 请求和断言，验证 GoAdmin
//       管理后台在 Beego 框架下的功能完整性
//
// 测试内容：
//   - 认证功能测试（登录、登出）
//   - 权限管理测试（增删改查）
//   - 角色管理测试（增删改查、权限分配）
//   - 管理员管理测试（增删改查、角色分配）
//   - 菜单管理测试（增删改查、层级管理）
//   - 操作日志测试（记录查询）
//   - 外部数据源测试
//   - 正常表格测试
//
// 测试框架：
//   - httpexpect: HTTP 测试库，用于编写 HTTP 断言
//   - Beego: Go Web 框架，提供路由和中间件功能
//
// 运行方式：
//   在项目根目录执行：go test ./tests/frameworks/beego -v
//   或运行特定测试：go test ./tests/frameworks/beego -run TestBeego -v
//
// 注意事项：
//   - 测试前需要确保数据库服务已启动
//   - 测试数据库 go-admin-test 需要预先创建
//   - 测试会修改数据库内容，建议使用独立的测试数据库

package beego

import (
	"net/http" // Go 标准网络包，提供 HTTP 客户端和服务端功能
	"testing"  // Go 标准测试包

	"github.com/gavv/httpexpect"                 // HTTP 测试库，用于编写 HTTP 断言
	"github.com/purpose168/GoAdmin/tests/common" // 测试公共模块
)

// ==================== 测试函数定义 ====================

// TestBeego 执行 Beego 框架测试
// 参数：
//   - t: 测试对象，用于测试断言和报告
//
// 说明：
//
//	该函数执行 GoAdmin 在 Beego 框架下的完整集成测试。
//	测试流程：
//	1. 创建 httpexpect 客户端配置
//	2. 使用 internalHandler() 创建的 Beego 处理器作为传输层
//	3. 创建新的 Cookie Jar 用于管理会话
//	4. 创建断言报告器用于报告测试结果
//	5. 调用 common.ExtraTest 执行额外的通用测试
//
// 测试内容：
//   - 认证测试（登录、登出、会话管理）
//   - 权限管理测试（创建、读取、更新、删除权限）
//   - 角色管理测试（创建、读取、更新、删除角色，分配权限）
//   - 管理员管理测试（创建、读取、更新、删除管理员，分配角色）
//   - 菜单管理测试（创建、读取、更新、删除菜单，管理层级）
//   - 操作日志测试（查询操作记录）
//   - 外部数据源测试（从外部数据源获取数据）
//   - 正常表格测试（标准 CRUD 操作）
//
// 技术细节：
//   - 使用 httpexpect.NewBinder 将 Beego 处理器绑定到 httpexpect
//   - 使用 httpexpect.NewJar 创建 Cookie Jar，自动管理会话
//   - 使用 httpexpect.NewAssertReporter 创建断言报告器，输出测试结果
//
// 使用场景：
//   - CI/CD 流水线中的自动化测试
//   - 代码提交前的回归测试
//   - 版本发布前的集成测试
//   - Beego 框架适配验证
//
// 运行方式：
//
//	在项目根目录执行：go test ./tests/frameworks/beego -v
//	或运行特定测试：go test ./tests/frameworks/beego -run TestBeego -v
//
// 注意事项：
//   - 确保 Beego 框架已正确安装和配置
//   - 测试前需要启动数据库服务
//   - 测试数据库需要预先创建
//   - 测试会创建和删除测试数据
func TestBeego(t *testing.T) {
	// 调用通用测试函数，传入 httpexpect 配置
	common.ExtraTest(httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			// 使用内部处理器作为传输层，将 Beego 处理器绑定到 httpexpect
			Transport: httpexpect.NewBinder(internalHandler()), // 使用内部处理器作为传输层
			// 创建新的 Cookie Jar，用于管理会话和 Cookie
			Jar: httpexpect.NewJar(), // 创建新的 Cookie Jar
		},
		// 创建断言报告器，用于报告测试结果
		Reporter: httpexpect.NewAssertReporter(t), // 创建断言报告器
	}))
}

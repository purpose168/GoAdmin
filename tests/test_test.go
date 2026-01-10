// test_test.go - 测试套件入口文件
// 包名：tests
// 作者：GoAdmin 团队
// 创建日期：2020年
// 描述：本文件是 GoAdmin 测试套件的入口文件，用于执行内置表的黑盒测试
//       测试框架支持多种 Web 框架（Gin、Echo、Chi 等），测试内容涵盖
//       权限管理、角色管理、管理员管理、菜单管理、操作日志等核心功能

package tests

import (
	"testing" // Go 标准测试包
	"time"    // Go 标准时间包

	"github.com/purpose168/GoAdmin/modules/config"       // 配置模块
	"github.com/purpose168/GoAdmin/tests/frameworks/gin" // Gin 框架测试适配器
)

// ==================== 测试函数定义 ====================

// TestBlackBoxTestSuitOfBuiltInTables 内置表黑盒测试套件
// 参数：t - 测试对象，用于测试断言和报告
// 说明：
//
//	该函数执行 GoAdmin 内置表的完整黑盒测试，包括：
//	- 权限（Permission）表测试
//	- 角色（Role）表测试
//	- 管理员（Manager）表测试
//	- 菜单（Menu）表测试
//	- 操作日志（Operation Log）表测试
//
// 测试框架：
//
//	使用 Gin Web 框架作为测试环境，通过 gin.NewHandler 创建处理器
//
// 数据库配置：
//
//	使用 MySQL 数据库，连接信息如下：
//	- 主机：127.0.0.1
//	- 端口：3306
//	- 用户：root
//	- 密码：root
//	- 数据库名：go-admin-test
//	- 最大空闲连接数：50
//	- 最大打开连接数：150
//	- 连接最大生命周期：1小时
//	- 连接最大空闲时间：0（无限制）
//
// 测试内容：
//  1. 认证测试（登录、登出）
//  2. 权限管理测试（增删改查）
//  3. 角色管理测试（增删改查、权限分配）
//  4. 管理员管理测试（增删改查、角色分配）
//  5. 菜单管理测试（增删改查、层级管理）
//  6. 操作日志测试（记录查询）
//
// 使用场景：
//   - CI/CD 流水线中的自动化测试
//   - 代码提交前的回归测试
//   - 版本发布前的集成测试
//   - 功能验证测试
//
// 注意事项：
//   - 测试前需要确保 MySQL 服务已启动
//   - 测试数据库 go-admin-test 需要预先创建
//   - 测试会修改数据库内容，建议使用独立的测试数据库
//   - 测试过程中会创建和删除测试数据
//
// 运行方式：
//
//	在项目根目录执行：go test ./tests -v
//	或运行特定测试：go test ./tests -run TestBlackBoxTestSuitOfBuiltInTables -v
func TestBlackBoxTestSuitOfBuiltInTables(t *testing.T) {
	// 调用黑盒测试套件，传入 Gin 框架处理器和数据库配置
	BlackBoxTestSuitOfBuiltInTables(t, gin.NewHandler, config.DatabaseList{
		"default": { // 默认数据库连接配置
			Host:            "127.0.0.1",        // 数据库主机地址
			Port:            "3306",             // 数据库端口号
			User:            "root",             // 数据库用户名
			Pwd:             "root",             // 数据库密码
			Name:            "go-admin-test",    // 数据库名称
			MaxIdleConns:    50,                 // 最大空闲连接数
			MaxOpenConns:    150,                // 最大打开连接数
			ConnMaxLifetime: time.Hour,          // 连接最大生命周期（1小时）
			ConnMaxIdleTime: 0,                  // 连接最大空闲时间（0表示无限制）
			Driver:          config.DriverMysql, // 数据库驱动类型（MySQL）
		},
	})
}

// Package web 提供用户验收测试（UAT）的测试框架
// 包描述:
//
//	该包提供了完整的用户验收测试框架，用于自动化测试Web应用的UI交互
//	基于agouti库和ChromeDriver实现浏览器自动化操作
//
// 主要功能:
//   - 提供测试套件框架（UserAcceptanceTestSuit）
//   - 管理浏览器驱动和页面生命周期
//   - 支持本地和远程（无头模式）测试
//   - 提供测试辅助工具和颜色输出
//
// 核心组件:
//   - Testers: 测试函数类型，定义测试逻辑
//   - ServerStarter: 服务器启动函数类型，用于启动测试服务器
//   - UserAcceptanceTestSuit: 主测试套件函数
//
// 使用场景:
//   - Web应用的端到端（E2E）测试
//   - 用户验收测试（UAT）
//   - UI自动化测试
//   - 回归测试
//
// 技术栈:
//   - agouti: Go语言的Web驱动库
//   - ChromeDriver: Chrome浏览器的WebDriver
//   - ansi: 终端颜色输出库
//   - GoAdmin: Go后台管理框架
//
// 依赖说明:
//   - 需要安装Chrome浏览器和对应版本的ChromeDriver
//   - ChromeDriver版本必须与Chrome版本匹配
//   - 使用以下链接获取最新版本: https://googlechromelabs.github.io/chrome-for-testing/
//
// 作者: GoAdmin Team
// 创建日期: 2020-01-01
package web

import (
	"fmt"
	"testing"
	"time"

	"github.com/mgutz/ansi"

	_ "github.com/purpose168/GoAdmin-themes/adminlte"
	_ "github.com/purpose168/GoAdmin/adapter/gin"
	_ "github.com/purpose168/GoAdmin/modules/db/drivers/mysql"

	"github.com/sclevine/agouti"
)

// Testers 测试函数类型
// 描述:
//
//	定义了测试函数的类型签名，用于封装具体的测试逻辑
//
// 参数:
//   - t: testing.T对象，用于测试断言和报告
//   - page: Page对象，用于浏览器页面操作
//
// 使用示例:
//
//	func myTest(t *testing.T, page *Page) {
//	    page.NavigateTo("/admin")
//	    page.Contain("欢迎")
//	}
//
// 说明:
//
//	该类型用于UserAcceptanceTestSuit函数，允许传入自定义的测试逻辑
type Testers func(t *testing.T, page *Page)

// ServerStarter 服务器启动函数类型
// 描述:
//
//	定义了服务器启动函数的类型签名，用于启动测试服务器
//
// 参数:
//   - quit: 退出通道，用于通知服务器停止
//
// 使用示例:
//
//	func startServer(quit chan struct{}) {
//	    go func() {
//	        <-quit
//	        // 停止服务器
//	    }()
//	    // 启动服务器
//	}
//
// 说明:
//
//	该类型用于UserAcceptanceTestSuit函数，允许传入自定义的服务器启动逻辑
type ServerStarter func(quit chan struct{})

// UserAcceptanceTestSuit 用户验收测试套件
// 描述:
//
//	提供完整的用户验收测试框架，自动化执行UI测试
//	管理浏览器驱动、页面生命周期和测试执行流程
//
// 参数:
//   - t: testing.T对象，用于测试断言和报告
//   - testers: Testers类型，包含具体的测试逻辑
//   - serverStarter: ServerStarter类型，用于启动测试服务器
//   - local: bool类型，是否在本地运行（true=显示浏览器窗口，false=无头模式）
//   - options: 可变参数，Chrome浏览器启动选项
//
// 功能:
//  1. 启动测试服务器（通过serverStarter）
//  2. 配置Chrome浏览器选项（默认或自定义）
//  3. 启动ChromeDriver和浏览器
//  4. 创建新页面
//  5. 执行测试逻辑（通过testers）
//  6. 清理资源（关闭页面、停止驱动、停止服务器）
//
// 默认Chrome选项说明:
//   - --user-agent: 设置用户代理字符串（模拟Mac Chrome浏览器）
//   - --window-size: 设置浏览器窗口大小（1500x900像素）
//   - --incognito: 使用无痕模式
//   - --blink-settings=imagesEnabled=true: 启用图片加载
//   - --no-default-browser-check: 禁用默认浏览器检查
//   - --ignore-ssl-errors=true: 忽略SSL错误
//   - --ssl-protocol=any: 允许任意SSL协议
//   - --no-sandbox: 禁用沙箱（用于CI/CD环境）
//   - --disable-breakpad: 禁用崩溃报告
//   - --disable-gpu: 禁用GPU加速
//   - --disable-logging: 禁用日志
//   - --no-zygote: 禁用zygote进程
//   - --allow-running-insecure-content: 允许运行不安全内容
//   - --headless: 无头模式（仅当local=false时添加）
//
// 使用场景:
//   - 本地开发时进行UI测试（local=true）
//   - CI/CD流水线中进行自动化测试（local=false）
//   - 回归测试
//   - 用户验收测试（UAT）
//
// 重要说明:
//   - 确保ChromeDriver版本与Chrome浏览器版本匹配
//   - 使用以下链接获取最新版本: https://googlechromelabs.github.io/chrome-for-testing/
//   - 测试失败时会panic，错误信息会包含详细原因
//
// 使用示例:
//
//	func TestLogin(t *testing.T) {
//	    UserAcceptanceTestSuit(t, myTesters, startServer, false)
//	}
//
// 错误处理:
//   - 启动驱动失败: panic "failed to start driver"
//   - 打开页面失败: panic "failed to open page"
//   - 关闭页面失败: 打印错误信息但不中断
//   - 停止驱动失败: 打印错误信息但不中断
func UserAcceptanceTestSuit(t *testing.T, testers Testers, serverStarter ServerStarter, local bool, options ...string) {
	// 创建退出通道，用于通知服务器停止
	var quit = make(chan struct{})
	// 启动测试服务器
	go serverStarter(quit)

	// 如果没有提供自定义选项，使用默认的Chrome浏览器选项
	if len(options) == 0 {
		options = []string{
			// 设置用户代理字符串，模拟Mac Chrome浏览器
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36",
			// 设置浏览器窗口大小为1500x900像素
			"--window-size=1500,900",
			// 使用无痕模式
			"--incognito",
			// 启用图片加载
			"--blink-settings=imagesEnabled=true",
			// 禁用默认浏览器检查
			"--no-default-browser-check",
			// 忽略SSL错误
			"--ignore-ssl-errors=true",
			// 允许任意SSL协议
			"--ssl-protocol=any",
			// 禁用沙箱（用于CI/CD环境）
			"--no-sandbox",
			// 禁用崩溃报告
			"--disable-breakpad",
			// 禁用GPU加速
			"--disable-gpu",
			// 禁用日志
			"--disable-logging",
			// 禁用zygote进程
			"--no-zygote",
			// 允许运行不安全内容
			"--allow-running-insecure-content",
		}
		// 如果不是本地运行，添加无头模式选项
		if !local {
			options = append(options, "--headless")
		}
	}

	// 创建ChromeDriver，配置浏览器选项和功能
	driver := agouti.ChromeDriver(
		// 设置Chrome命令行参数
		agouti.ChromeOptions("args", options),
		// 设置Chrome功能
		agouti.Desired(
			agouti.Capabilities{
				// 设置日志首选项，记录所有性能日志
				"loggingPrefs": map[string]string{
					"performance": "ALL",
				},
				// 接受SSL证书
				"acceptSslCerts": true,
				// 接受不安全证书
				"acceptInsecureCerts": true,
			},
		))
	// 启动ChromeDriver
	err := driver.Start()
	if err != nil {
		// 启动失败时panic
		panic("failed to start driver, error: " + err.Error())
	}

	// 创建新页面
	page, err := driver.NewPage()
	if err != nil {
		// 创建页面失败时panic
		panic("failed to open page, error: " + err.Error())
	}

	// 打印测试开始信息
	fmt.Println()
	fmt.Println("============================================")
	// 使用蓝色打印测试标题
	printlnWithColor("User Acceptance Testing", "blue")
	fmt.Println("============================================")
	fmt.Println()

	// 执行测试逻辑
	testers(t, &Page{T: t, Page: page, Driver: driver, Quit: quit})

	// 等待2秒，确保所有操作完成
	wait(2)

	// 如果不是本地运行，清理资源
	if !local {
		// 关闭页面窗口
		err = page.CloseWindow()
		if err != nil {
			// 关闭失败时打印错误信息
			fmt.Println("failed to close page, error: ", err)
		}

		// 销毁页面对象
		err = page.Destroy()
		if err != nil {
			// 销毁失败时打印错误信息
			fmt.Println("failed to destroy page, error: ", err)
		}

		// 停止ChromeDriver
		err = driver.Stop()
		if err != nil {
			// 停止失败时打印错误信息
			fmt.Println("failed to stop driver, error: ", err)
		}
	}

	// 发送退出信号，通知服务器停止
	quit <- struct{}{}
}

// printlnWithColor 使用指定颜色打印消息
// 参数:
//   - msg: 要打印的消息内容
//   - color: 颜色名称（如"blue"、"green"等）
//
// 功能:
//   - 使用ansi库将消息以指定颜色输出到控制台
//
// 使用场景:
//   - 打印测试标题和重要信息
//   - 区分不同类型的输出信息
//   - 提高测试输出的可读性
//
// 说明:
//
//	支持的颜色包括: black, red, green, yellow, blue, magenta, cyan, white等
func printlnWithColor(msg string, color string) {
	fmt.Println(ansi.Color(msg, color))
}

// printPart 打印测试部分信息
// 参数:
//   - part: 部分名称或描述
//
// 功能:
//   - 使用蓝色打印带"> "前缀的部分信息
//
// 使用场景:
//   - 标记测试的不同阶段
//   - 显示测试执行的进度
//   - 区分不同的测试模块
//
// 说明:
//
//	内部调用printlnWithColor，固定使用蓝色
func printPart(part string) {
	printlnWithColor("> "+part, colorBlue)
}

// wait 等待指定的时间
// 参数:
//   - t: 等待的秒数
//
// 功能:
//   - 暂停执行指定的秒数
//
// 使用场景:
//   - 等待页面加载完成
//   - 等待动画效果完成
//   - 等待异步操作完成
//   - 测试之间的延迟
//
// 说明:
//
//	使用time.Sleep实现，参数单位为秒
func wait(t int) {
	time.Sleep(time.Duration(t) * time.Second)
}

// basePath 测试服务器的基础URL
// 描述:
//
//	定义测试服务器的基础地址，所有测试URL都基于此地址
//
// 说明:
//
//	默认端口为9033，可根据实际测试环境修改
const basePath = "http://localhost:9033"

// url 构建完整的测试URL
// 参数:
//   - suffix: URL后缀（如"/info"、"/edit"等）
//
// 返回值: 完整的URL字符串
// 功能:
//   - 将基础路径、"/admin"和后缀拼接成完整URL
//   - 处理后缀为"/"的特殊情况
//
// 使用场景:
//   - 构建测试页面的URL
//   - 导航到不同的管理页面
//
// 说明:
//
//	如果suffix为"/"，则将其转换为空字符串
//	最终URL格式: basePath + "/admin" + suffix
//
// 示例:
//
//	url("/info") -> "http://localhost:9033/admin/info"
//	url("/") -> "http://localhost:9033/admin"
func url(suffix string) string {
	if suffix == "/" {
		suffix = ""
	}
	return basePath + "/admin" + suffix
}

// colorBlue 蓝色常量
// 描述:
//
//	用于终端输出的蓝色标识符
//
// 使用场景:
//   - 打印测试标题
//   - 打印测试部分信息
//   - 标记重要信息
const colorBlue = "blue"

// colorGreen 绿色常量
// 描述:
//
//	用于终端输出的绿色标识符
//
// 使用场景:
//   - 打印成功信息
//   - 打印通过的测试
//   - 标记正面结果
const colorGreen = "green"

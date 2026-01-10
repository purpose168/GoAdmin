// Package web 提供 Web 页面测试的辅助工具
//
// 本包基于 agouti 库提供了 Web 页面操作的封装方法，用于编写端到端（E2E）测试，简化浏览器自动化操作。
//
// 主要功能：
//   - 页面元素定位和操作（点击、输入、获取文本等）
//   - 页面内容验证（包含、不包含、CSS 样式等）
//   - 鼠标操作（移动、点击等）
//   - 等待和导航控制
//
// 使用场景：
//   - Web 应用的 UI 自动化测试
//   - 浏览器兼容性测试
//   - 交互功能验证测试
//
// 技术栈：
//   - agouti: Go 语言的 Web 驱动库
//   - testify: Go 语言的测试断言库
//
// 核心概念：
//   - Page: 页面结构体，封装了 agouti.Page 和测试相关的信息
//   - 元素定位: 使用 XPath 或选择器定位页面元素
//   - 元素操作: 点击、输入、获取文本、获取属性等
//   - 页面验证: 验证页面内容、CSS 样式、元素状态等
//   - 等待控制: 等待页面加载、动画效果、异步操作等
//   - 导航控制: 导航到指定页面
//   - 错误处理: 捕获 panic 并打印堆栈信息
//   - 资源清理: 销毁页面对象，停止浏览器驱动
//
// 元素定位方式：
//   - XPath: 使用 XPath 表达式定位元素
//   - Selection: 使用 agouti.Selection 对象定位元素
//
// 元素操作类型：
//   - 点击: Click、ClickS
//   - 输入: Fill
//   - 获取文本: Text
//   - 获取属性: Attr、Value
//   - 获取 CSS: Css、CssS
//
// 页面验证类型：
//   - 内容验证: Contain、NoContain
//   - CSS 验证: Css、CssS、Display、Nondisplay
//   - 文本验证: Text
//   - 属性验证: Attr、Value
//
// 鼠标操作类型：
//   - 移动: MoveMouseBy
//
// 等待控制：
//   - Wait: 等待指定的时间
//   - Click: 点击后自动等待（默认 1 秒）
//   - NavigateTo: 导航后自动等待（2 秒）
//
// 导航控制：
//   - NavigateTo: 导航到指定路径
//
// 错误处理：
//   - Destroy: 捕获 panic 并打印堆栈信息
//   - 断言失败: 使用 assert 包进行断言，失败时标记测试失败
//
// 资源清理：
//   - Destroy: 销毁页面对象，停止浏览器驱动
//   - Quit: 发送退出信号
//
// 注意事项：
//   - 需要确保浏览器驱动已正确安装
//   - 需要确保页面元素已正确加载
//   - 需要确保 XPath 表达式正确
//   - 需要确保等待时间足够
//   - 需要确保测试环境稳定
//
// 作者: GoAdmin Team
// 创建日期: 2020-01-01
// 版本: 1.0.0
package web

import (
	"fmt"           // 格式化输出包，提供字符串格式化功能
	"runtime/debug" // 运行时调试包，提供堆栈信息打印功能
	"strings"       // 字符串处理包，提供字符串分割、连接等功能
	"testing"       // 测试包，提供测试框架和断言功能
	"time"          // 时间包，提供时间相关功能

	"github.com/sclevine/agouti"         // agouti 包，提供 Web 驱动功能
	"github.com/stretchr/testify/assert" // testify 包，提供测试断言功能
)

// Page 页面结构体
//
// 描述：
//
//	封装了 agouti.Page 和测试相关的信息，提供便捷的页面操作方法。
//
// 字段说明：
//   - Page: agouti.Page 对象，用于浏览器页面操作
//   - T: testing.T 对象，用于测试断言和报告
//   - Driver: agouti.WebDriver 对象，用于浏览器驱动管理
//   - Quit: 退出通道，用于通知测试结束
//
// 使用示例：
//
//	page := &Page{
//	    Page: driver.NewPage(),
//	    T: t,
//	    Driver: driver,
//	    Quit: make(chan struct{}),
//	}
//
// 注意事项：
//   - 需要确保 Page 对象已正确初始化
//   - 需要确保 T 对象已正确传入
//   - 需要确保 Driver 对象已正确初始化
//   - 需要确保 Quit 通道已正确创建
type Page struct {
	*agouti.Page                   // agouti.Page 对象，用于浏览器页面操作
	T            *testing.T        // testing.T 对象，用于测试断言和报告
	Driver       *agouti.WebDriver // agouti.WebDriver 对象，用于浏览器驱动管理
	Quit         chan struct{}     // 退出通道，用于通知测试结束
}

// Destroy 销毁页面并清理资源
//
// 功能：
//   - 捕获 panic 并打印堆栈信息
//   - 销毁页面对象
//   - 停止浏览器驱动
//   - 标记测试失败
//   - 发送退出信号
//
// 使用场景：
//   - 测试结束时清理资源
//   - 发生 panic 时的错误处理
//
// 说明：
//
//	该方法使用 defer 调用，确保在任何情况下都能正确清理资源。
//	如果发生 panic，会打印堆栈信息，销毁页面对象，停止浏览器驱动，
//	标记测试失败，并发送退出信号。
//
// 注意事项：
//   - 该方法应该使用 defer 调用
//   - 该方法会捕获 panic 并处理
//   - 该方法会标记测试失败
//   - 该方法会发送退出信号
func (page *Page) Destroy() {
	if r := recover(); r != nil { // 捕获 panic
		debug.PrintStack()               // 打印堆栈信息
		fmt.Println("Recovered in f", r) // 打印 panic 信息
		_ = page.Page.Destroy()          // 销毁页面对象
		_ = page.Driver.Stop()           // 停止浏览器驱动
		page.T.Fail()                    // 标记测试失败
		page.Quit <- struct{}{}          // 发送退出信号
	}
}

// Wait 等待指定的时间
//
// 参数:
//   - t: 等待的秒数
//
// 功能：
//   - 暂停执行指定的秒数
//
// 使用场景：
//   - 等待页面加载完成
//   - 等待动画效果完成
//   - 等待异步操作完成
//
// 说明：
//
//	使用 time.Sleep 实现，参数单位为秒。
//
// 注意事项：
//   - 等待时间不宜过长，影响测试效率
//   - 等待时间不宜过短，可能导致操作失败
func (page *Page) Wait(t int) {
	time.Sleep(time.Duration(t) * time.Second) // 暂停执行指定的秒数
}

// Contain 验证页面 HTML 内容是否包含指定字符串
//
// 参数:
//   - s: 要查找的字符串
//
// 功能：
//   - 获取页面 HTML 内容
//   - 验证是否包含指定字符串
//
// 使用场景：
//   - 验证页面是否显示了特定内容
//   - 验证操作是否成功
//   - 验证错误消息是否显示
//
// 断言：
//   - HTML 获取不能出错
//   - 必须包含指定字符串
//
// 注意事项：
//   - 字符串区分大小写
//   - 字符串需要转义特殊字符
//   - 需要确保页面已加载完成
func (page *Page) Contain(s string) {
	content, err := page.HTML()                              // 获取页面 HTML 内容
	assert.Equal(page.T, err, nil)                           // 断言 HTML 获取不能出错
	assert.Equal(page.T, strings.Contains(content, s), true) // 断言必须包含指定字符串
}

// NoContain 验证页面 HTML 内容是否不包含指定字符串
//
// 参数:
//   - s: 要查找的字符串
//
// 功能：
//   - 获取页面 HTML 内容
//   - 验证是否不包含指定字符串
//
// 使用场景：
//   - 验证页面是否不显示特定内容
//   - 验证操作后内容是否消失
//   - 验证错误消息是否隐藏
//
// 断言：
//   - HTML 获取不能出错
//   - 必须不包含指定字符串
//
// 注意事项：
//   - 字符串区分大小写
//   - 字符串需要转义特殊字符
//   - 需要确保页面已加载完成
func (page *Page) NoContain(s string) {
	content, err := page.HTML()                               // 获取页面 HTML 内容
	assert.Equal(page.T, err, nil)                            // 断言 HTML 获取不能出错
	assert.Equal(page.T, strings.Contains(content, s), false) // 断言必须不包含指定字符串
}

// Css 验证指定 XPath 元素的 CSS 样式属性值
//
// 参数:
//   - xpath: 元素的 XPath 表达式
//   - css: CSS 属性名（如 display、color 等）
//   - res: 期望的 CSS 属性值
//
// 功能：
//   - 通过 XPath 定位元素
//   - 获取元素的指定 CSS 属性值
//   - 验证 CSS 属性值是否匹配期望值
//
// 使用场景：
//   - 验证元素的显示状态（display: block/none）
//   - 验证元素的颜色、字体等样式
//   - 验证元素的布局属性
//
// 断言：
//   - CSS 属性获取不能出错
//   - CSS 属性值必须等于期望值
//
// 注意事项：
//   - XPath 表达式需要正确
//   - CSS 属性名需要正确
//   - 需要确保元素已加载完成
//   - 需要确保元素可见
func (page *Page) Css(xpath, css, res string) {
	style, err := page.FindByXPath(xpath).CSS(css) // 获取元素的指定 CSS 属性值
	assert.Equal(page.T, err, nil)                 // 断言 CSS 属性获取不能出错
	assert.Equal(page.T, style, res)               // 断言 CSS 属性值必须等于期望值
}

// CssS 验证指定选择器元素的 CSS 样式属性值
//
// 参数:
//   - s: agouti.Selection 对象，表示元素选择器
//   - css: CSS 属性名（如 display、color 等）
//   - res: 期望的 CSS 属性值
//
// 功能：
//   - 通过选择器定位元素
//   - 获取元素的指定 CSS 属性值
//   - 验证 CSS 属性值是否匹配期望值
//
// 使用场景：
//   - 验证元素的显示状态（display: block/none）
//   - 验证元素的颜色、字体等样式
//   - 验证元素的布局属性
//
// 说明：
//
//	与 Css 方法的区别是使用 Selection 对象而不是 XPath。
//
// 断言：
//   - CSS 属性获取不能出错
//   - CSS 属性值必须等于期望值
//
// 注意事项：
//   - Selection 对象需要正确
//   - CSS 属性名需要正确
//   - 需要确保元素已加载完成
//   - 需要确保元素可见
func (page *Page) CssS(s *agouti.Selection, css, res string) {
	style, err := s.CSS(css)         // 获取元素的指定 CSS 属性值
	assert.Equal(page.T, err, nil)   // 断言 CSS 属性获取不能出错
	assert.Equal(page.T, style, res) // 断言 CSS 属性值必须等于期望值
}

// Text 验证指定 XPath 元素的文本内容
//
// 参数:
//   - xpath: 元素的 XPath 表达式
//   - text: 期望的文本内容
//
// 功能：
//   - 通过 XPath 定位元素
//   - 获取元素的文本内容
//   - 验证文本内容是否匹配期望值
//
// 使用场景：
//   - 验证按钮文本是否正确
//   - 验证标题、标签等文本内容
//   - 验证错误提示信息
//
// 断言：
//   - 文本获取不能出错
//   - 文本内容必须等于期望值
//
// 注意事项：
//   - XPath 表达式需要正确
//   - 文本内容区分大小写
//   - 需要确保元素已加载完成
//   - 需要确保元素可见
func (page *Page) Text(xpath, text string) {
	mli1, err := page.FindByXPath(xpath).Text() // 获取元素的文本内容
	assert.Equal(page.T, err, nil)              // 断言文本获取不能出错
	assert.Equal(page.T, mli1, text)            // 断言文本内容必须等于期望值
}

// MoveMouseBy 移动鼠标到指定偏移量位置
//
// 参数:
//   - xOffset: X 轴偏移量（像素）
//   - yOffset: Y 轴偏移量（像素）
//
// 功能：
//   - 将鼠标从当前位置移动到指定偏移量的位置
//
// 使用场景：
//   - 测试悬停效果
//   - 测试拖拽操作
//   - 测试鼠标交互效果
//
// 断言：
//   - 鼠标移动操作不能出错
//
// 注意事项：
//   - 偏移量可以是正数或负数
//   - 需要确保鼠标位置有效
//   - 需要确保页面已加载完成
func (page *Page) MoveMouseBy(xOffset, yOffset int) {
	assert.Equal(page.T, page.Page.MoveMouseBy(xOffset, yOffset), nil) // 断言鼠标移动操作不能出错
}

// Display 验证指定 XPath 元素的 display 属性为 block
//
// 参数:
//   - xpath: 元素的 XPath 表达式
//
// 功能：
//   - 验证元素的 display CSS 属性值为 block
//   - 即验证元素是否可见（块级显示）
//
// 使用场景：
//   - 验证元素是否显示
//   - 验证操作后元素是否出现
//   - 验证展开/折叠效果
//
// 说明：
//
//	内部调用 Css 方法验证 display 属性。
//
// 注意事项：
//   - XPath 表达式需要正确
//   - 需要确保元素已加载完成
//   - 需要确保元素可见
func (page *Page) Display(xpath string) {
	page.Css(xpath, "display", "block") // 验证 display 属性为 block
}

// Nondisplay 验证指定 XPath 元素的 display 属性为 none
//
// 参数:
//   - xpath: 元素的 XPath 表达式
//
// 功能：
//   - 验证元素的 display CSS 属性值为 none
//   - 即验证元素是否隐藏
//
// 使用场景：
//   - 验证元素是否隐藏
//   - 验证操作后元素是否消失
//   - 验证展开/折叠效果
//
// 说明：
//
//	内部调用 Css 方法验证 display 属性。
//
// 注意事项：
//   - XPath 表达式需要正确
//   - 需要确保元素已加载完成
func (page *Page) Nondisplay(xpath string) {
	page.Css(xpath, "display", "none") // 验证 display 属性为 none
}

// Value 验证指定 XPath 元素的 value 属性值
//
// 参数:
//   - xpath: 元素的 XPath 表达式
//   - value: 期望的 value 属性值
//
// 功能：
//   - 通过 XPath 定位元素
//   - 获取元素的 value 属性值
//   - 验证 value 属性值是否匹配期望值
//
// 使用场景：
//   - 验证输入框的值
//   - 验证选择框的选中值
//   - 验证隐藏字段的值
//
// 断言：
//   - value 属性获取不能出错
//   - value 属性值必须等于期望值
//
// 注意事项：
//   - XPath 表达式需要正确
//   - value 属性值区分大小写
//   - 需要确保元素已加载完成
func (page *Page) Value(xpath, value string) {
	val, err := page.FindByXPath(xpath).Attribute("value") // 获取元素的 value 属性值
	assert.Equal(page.T, err, nil)                         // 断言 value 属性获取不能出错
	assert.Equal(page.T, val, value)                       // 断言 value 属性值必须等于期望值
}

// Click 点击指定 XPath 元素
//
// 参数:
//   - xpath: 元素的 XPath 表达式
//   - intervals: 可选参数，点击后等待的秒数（默认为 1 秒）
//
// 功能：
//   - 通过 XPath 定位元素
//   - 点击该元素
//   - 等待指定时间（默认 1 秒）
//
// 使用场景：
//   - 点击按钮
//   - 点击链接
//   - 点击复选框、单选框
//   - 点击菜单项
//
// 说明：
//
//	intervals 为可变参数，如果不提供则默认等待 1 秒。
//
// 断言：
//   - 点击操作不能出错
//
// 注意事项：
//   - XPath 表达式需要正确
//   - 需要确保元素已加载完成
//   - 需要确保元素可见
//   - 需要确保元素可点击
func (page *Page) Click(xpath string, intervals ...int) {
	assert.Equal(page.T, page.FindByXPath(xpath).Click(), nil) // 断言点击操作不能出错
	interval := 1                                              // 默认等待 1 秒
	if len(intervals) > 0 {                                    // 如果提供了等待时间
		interval = intervals[0] // 使用提供的等待时间
	}
	page.Wait(interval) // 等待指定时间
}

// ClickS 点击指定选择器元素
//
// 参数:
//   - s: agouti.Selection 对象，表示元素选择器
//   - intervals: 可选参数，点击后等待的秒数（默认为 1 秒）
//
// 功能：
//   - 通过选择器定位元素
//   - 点击该元素
//   - 等待指定时间（默认 1 秒）
//
// 使用场景：
//   - 点击按钮
//   - 点击链接
//   - 点击复选框、单选框
//   - 点击菜单项
//
// 说明：
//
//	与 Click 方法的区别是使用 Selection 对象而不是 XPath。
//	intervals 为可变参数，如果不提供则默认等待 1 秒。
//
// 断言：
//   - 点击操作不能出错
//
// 注意事项：
//   - Selection 对象需要正确
//   - 需要确保元素已加载完成
//   - 需要确保元素可见
//   - 需要确保元素可点击
func (page *Page) ClickS(s *agouti.Selection, intervals ...int) {
	assert.Equal(page.T, s.Click(), nil) // 断言点击操作不能出错
	interval := 1                        // 默认等待 1 秒
	if len(intervals) > 0 {              // 如果提供了等待时间
		interval = intervals[0] // 使用提供的等待时间
	}
	page.Wait(interval) // 等待指定时间
}

// Attr 验证指定选择器元素的属性值
//
// 参数:
//   - s: agouti.Selection 对象，表示元素选择器
//   - attr: 属性名（如 value、href、class 等）
//   - res: 期望的属性值
//
// 功能：
//   - 通过选择器定位元素
//   - 获取元素的指定属性值
//   - 验证属性值是否匹配期望值
//
// 使用场景：
//   - 验证链接的 href 属性
//   - 验证图片的 src 属性
//   - 验证元素的 class 属性
//   - 验证自定义属性
//
// 断言：
//   - 属性获取不能出错
//   - 属性值必须等于期望值
//
// 注意事项：
//   - Selection 对象需要正确
//   - 属性名需要正确
//   - 属性值区分大小写
//   - 需要确保元素已加载完成
func (page *Page) Attr(s *agouti.Selection, attr, res string) {
	style, err := s.Attribute(attr)  // 获取元素的指定属性值
	assert.Equal(page.T, err, nil)   // 断言属性获取不能出错
	assert.Equal(page.T, style, res) // 断言属性值必须等于期望值
}

// Fill 向指定 XPath 元素填充内容
//
// 参数:
//   - xpath: 元素的 XPath 表达式
//   - content: 要填充的内容
//
// 功能：
//   - 通过 XPath 定位元素
//   - 向元素填充指定的内容
//
// 使用场景：
//   - 向输入框输入文本
//   - 向文本域输入多行文本
//   - 向可编辑元素输入内容
//
// 说明：
//
//	通常用于 input、textarea 等可输入元素。
//
// 断言：
//   - 填充操作不能出错
//
// 注意事项：
//   - XPath 表达式需要正确
//   - 需要确保元素已加载完成
//   - 需要确保元素可见
//   - 需要确保元素可编辑
func (page *Page) Fill(xpath, content string) {
	assert.Equal(page.T, page.FindByXPath(xpath).Fill(content), nil) // 断言填充操作不能出错
}

// NavigateTo 导航到指定路径
//
// 参数:
//   - path: 要导航到的 URL 路径
//
// 功能：
//   - 导航到指定的 URL 路径
//   - 等待 2 秒让页面加载完成
//
// 使用场景：
//   - 打开指定页面
//   - 页面跳转测试
//   - 测试页面导航功能
//
// 说明：
//
//	导航后会自动等待 2 秒，确保页面加载完成。
//
// 断言：
//   - 导航操作不能出错
//
// 注意事项：
//   - URL 路径需要正确
//   - 需要确保页面可访问
//   - 需要确保网络连接正常
func (page *Page) NavigateTo(path string) {
	assert.Equal(page.T, page.Navigate(path), nil) // 断言导航操作不能出错
	page.Wait(2)                                   // 等待 2 秒让页面加载完成
}

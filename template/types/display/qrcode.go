package display

import (
	"html/template"

	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/template/types"
)

// Qrcode 二维码显示生成器
// 用于将字段值转换为二维码显示
// 继承自 BaseDisplayFnGenerator，提供基础的显示函数生成能力
// 用户可以点击二维码图标查看对应的二维码图片
type Qrcode struct {
	types.BaseDisplayFnGenerator
}

// init 包初始化函数
// 在包加载时自动执行，将 Qrcode 类型注册到显示函数生成器注册表中
// 注册键名为 "qrcode"，可以通过该键名创建 Qrcode 实例
func init() {
	types.RegisterDisplayFnGenerator("qrcode", new(Qrcode))
}

// Get 获取字段过滤函数
// 根据传入的参数生成一个字段过滤函数，用于将字段值转换为二维码显示
//
// 参数：
//   - ctx: 上下文对象，包含请求相关的上下文信息
//   - args: 可变参数，当前实现不使用任何参数
//
// 返回值：
//   - FieldFilterFn: 字段过滤函数，接收 FieldModel 对象并返回转换后的 HTML
//     HTML 包含一个二维码图标和原始字段值
//
// 实现原理：
//  1. 使用 qrserver.com 的公共 API 生成二维码图片
//  2. 二维码图片 URL 格式：https://api.qrserver.com/v1/create-qr-code/?size=150x150&data={字段值}
//  3. 使用 Bootstrap 的 popover 组件显示二维码图片
//  4. 点击二维码图标时，弹出包含二维码图片的提示框
//
// 使用示例：
//
//	// 示例：将字段值显示为二维码
//	display.Qrcode{}.Get(ctx)
//
// 注意事项：
//   - 使用外部 API 生成二维码，需要网络连接
//   - 二维码图片大小固定为 150x150 像素
//   - 需要引入 Font Awesome 图标库（fa-qrcode 图标）
//   - 需要引入 Bootstrap 的 popover 组件
//   - 需要配合 JS() 方法返回的 JavaScript 代码使用
//   - 字段值会自动进行 URL 编码
func (q *Qrcode) Get(ctx *context.Context, args ...interface{}) types.FieldFilterFn {
	return func(value types.FieldModel) interface{} {
		// 构建二维码图片 URL
		// 使用 qrserver.com 的公共 API
		// size=150x150: 二维码图片大小
		// data={value.Value}: 二维码包含的数据内容
		src := `https://api.qrserver.com/v1/create-qr-code/?size=150x150&amp;data=` + value.Value

		// 返回包含二维码图标的 HTML
		// 使用 template.HTML 类型，避免 HTML 转义
		// HTML 结构：
		//   - <a> 标签：二维码图标链接
		//     - class="grid-column-qrcode": 用于 JavaScript 选择器
		//     - data-content: 包含二维码图片 HTML
		//     - data-toggle="popover": 启用 Bootstrap popover
		//     - tabindex="0": 使元素可聚焦
		//   - <i> 标签：Font Awesome 的二维码图标
		//   - value.Value：原始字段值显示
		return template.HTML(`
<a href="javascript:void(0);" class="grid-column-qrcode text-muted" 
	data-content="<img src='` + src + `' 
style='height:150px;width:150px;'/>" data-toggle="popover" tabindex="0" data-original-title="" title="">
<i class="fa fa-qrcode"></i>
</a>&nbsp;` + value.Value + `
`)
	}
}

// JS 获取 JavaScript 代码
// 返回实现二维码显示功能的 JavaScript 代码
// 该代码为所有二维码图标添加 Bootstrap popover 功能
//
// 返回值：
//   - template.HTML: JavaScript 代码，用于初始化 popover 组件
//
// 实现原理：
//  1. 使用 jQuery 选择器找到所有带有 .grid-column-qrcode 类的元素
//  2. 调用 Bootstrap 的 popover() 方法初始化提示框
//  3. 配置选项：
//     - html: true: 允许在提示框中使用 HTML 内容
//     - container: 'body': 将提示框添加到 body 元素中，避免样式问题
//     - trigger: 'focus': 当元素获得焦点时显示提示框
//
// 注意事项：
//   - 需要引入 jQuery 库
//   - 需要引入 Bootstrap 的 popover 组件
//   - popover 的内容来自 HTML 中的 data-content 属性
//   - 点击图标后，图标会获得焦点，触发 popover 显示
//   - 点击其他地方会关闭 popover
func (q *Qrcode) JS() template.HTML {
	return template.HTML(`
$('.grid-column-qrcode').popover({
	html: true,
	container: 'body',
	trigger: 'focus'
});
`)
}

package display

import (
	"html/template"

	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/template/types"
)

// Copyable 可复制显示生成器
// 用于将字段值转换为带有复制按钮的可视化显示
// 继承自 BaseDisplayFnGenerator，提供基础的显示函数生成能力
// 用户可以点击复制按钮将字段值复制到剪贴板
type Copyable struct {
	types.BaseDisplayFnGenerator
}

// init 包初始化函数
// 在包加载时自动执行，将 Copyable 类型注册到显示函数生成器注册表中
// 注册键名为 "copyable"，可以通过该键名创建 Copyable 实例
func init() {
	types.RegisterDisplayFnGenerator("copyable", new(Copyable))
}

// Get 获取字段过滤函数
// 根据传入的参数生成一个字段过滤函数，用于将字段值转换为带有复制按钮的显示
//
// 参数：
//   - ctx: 上下文对象，包含请求相关的上下文信息
//   - args: 可变参数，当前实现不使用任何参数
//
// 返回值：
//   - FieldFilterFn: 字段过滤函数，接收 FieldModel 对象并返回转换后的 HTML
//     HTML 包含一个复制按钮和原始字段值
//
// 使用示例：
//
//	// 示例：将字段值显示为可复制格式
//	display.Copyable{}.Get(ctx)
//
// 注意事项：
//   - 生成的HTML包含一个复制按钮（使用 Font Awesome 的 fa-copy 图标）
//   - 复制按钮点击后会调用 JS 方法中定义的复制逻辑
//   - 复制成功后会显示 "Copied!" 提示
//   - 需要配合 JS() 方法返回的 JavaScript 代码使用
//   - 需要引入 Font Awesome 图标库和 Bootstrap 的 tooltip 组件
func (c *Copyable) Get(ctx *context.Context, args ...interface{}) types.FieldFilterFn {
	return func(value types.FieldModel) interface{} {
		// 返回包含复制按钮和字段值的 HTML
		// 使用 template.HTML 类型，避免 HTML 转义
		// HTML 结构：
		//   - <a> 标签：复制按钮，包含 data-content 属性存储要复制的内容
		//   - <i> 标签：Font Awesome 的复制图标
		//   - &nbsp;：空格分隔符
		//   - value.Value：原始字段值
		return template.HTML(`
<a href="javascript:void(0);" class="grid-column-copyable text-muted" data-content="` + value.Value + `"
title="Copied!" data-placement="bottom">
<i class="fa fa-copy"></i>
</a>&nbsp;` + value.Value + `
`)
	}
}

// JS 获取 JavaScript 代码
// 返回实现复制功能的 JavaScript 代码
// 该代码通过事件委托的方式为所有复制按钮添加点击事件处理
//
// 返回值：
//   - template.HTML: JavaScript 代码，用于实现复制功能
//
// 实现原理：
//  1. 使用 jQuery 的事件委托机制监听 body 上的点击事件
//  2. 当点击带有 .grid-column-copyable 类的元素时触发
//  3. 从元素的 data-content 属性中获取要复制的内容
//  4. 创建一个临时的 input 元素，将内容写入并选中
//  5. 调用 document.execCommand("copy") 执行复制操作
//  6. 移除临时元素
//  7. 显示 tooltip 提示用户复制成功
//
// 注意事项：
//   - 需要引入 jQuery 库
//   - 需要引入 Bootstrap 的 tooltip 组件
//   - document.execCommand("copy") 是传统的复制方法，现代浏览器支持
//   - 对于更现代的复制方案，可以使用 Clipboard API
//   - 事件委托方式可以动态添加的复制按钮也能正常工作
func (c *Copyable) JS() template.HTML {
	return template.HTML(`
$('body').on('click','.grid-column-copyable',(function (e) {
	// 从点击元素的 data-content 属性中获取要复制的内容
	var content = $(this).data('content');

	// 创建一个临时的 input 元素用于复制操作
	var temp = $('<input>');

	// 将临时 input 元素添加到 body 中
	$("body").append(temp);

	// 设置 input 的值为要复制的内容，并选中所有文本
	temp.val(content).select();

	// 执行浏览器原生的复制命令
	// 将选中的文本复制到剪贴板
	document.execCommand("copy");

	// 移除临时 input 元素
	temp.remove();

	// 显示 Bootstrap tooltip 提示用户复制成功
	// tooltip 的标题在 HTML 中设置为 "Copied!"
	$(this).tooltip('show');
}))
`)
}

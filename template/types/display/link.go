package display

import (
	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/template"
	"github.com/purpose168/GoAdmin/template/types"
)

// Link 链接显示生成器
// 用于将字段值转换为可点击的链接显示
// 继承自 BaseDisplayFnGenerator，提供基础的显示函数生成能力
// 支持自定义URL前缀和是否在新标签页中打开
type Link struct {
	types.BaseDisplayFnGenerator
}

// init 包初始化函数
// 在包加载时自动执行，将 Link 类型注册到显示函数生成器注册表中
// 注册键名为 "link"，可以通过该键名创建 Link 实例
func init() {
	types.RegisterDisplayFnGenerator("link", new(Link))
}

// Get 获取字段过滤函数
// 根据传入的参数生成一个字段过滤函数，用于将字段值转换为链接显示
//
// 参数：
//   - ctx: 上下文对象，包含请求相关的上下文信息
//   - args: 可变参数，包含以下内容：
//   - args[0]（可选）: string 类型，URL 前缀
//     如果不提供，直接使用字段值作为链接URL
//     如果提供，将前缀拼接到字段值前面
//   - args[1]（可选）: []bool 类型，是否在新标签页中打开
//     如果提供，使用第一个元素的值（true 表示在新标签页打开）
//
// 返回值：
//   - FieldFilterFn: 字段过滤函数，接收 FieldModel 对象并返回转换后的链接 HTML
//
// 使用示例：
//
//	// 示例1：直接使用字段值作为链接
//	display.Link{}.Get(ctx)
//
//	// 示例2：添加 URL 前缀
//	display.Link{}.Get(ctx, "/user/profile/")
//	// 如果字段值为 "123"，生成的链接为 "/user/profile/123"
//
//	// 示例3：在新标签页中打开
//	display.Link{}.Get(ctx, "/article/", []bool{true})
//
// 注意事项：
//   - 参数 args[0] 如果提供，必须是 string 类型，否则会引发运行时 panic
//   - 参数 args[1] 如果提供，必须是 []bool 类型，否则会引发运行时 panic
//   - 使用 template.Default(ctx).Link() 构建链接组件
//   - URL 前缀和字段值之间会自动拼接，不会添加额外的斜杠
//   - OpenInNewTab() 方法会在链接中添加 target="_blank" 属性
func (l *Link) Get(ctx *context.Context, args ...interface{}) types.FieldFilterFn {
	// 初始化 URL 前缀为空字符串
	prefix := ""

	// 初始化是否在新标签页中打开为 false
	openInNewTabs := false

	// 如果提供了第一个参数，使用该参数作为 URL 前缀
	if len(args) > 0 {
		prefix = args[0].(string)
	}

	// 如果提供了第二个参数，检查是否在新标签页中打开
	if len(args) > 1 {
		// 使用类型断言检查参数是否为 []bool 类型
		if openInNewTabsArr, ok := args[1].([]bool); ok {
			// 如果类型正确，使用第一个元素的值
			openInNewTabs = openInNewTabsArr[0]
		}
	}

	return func(value types.FieldModel) interface{} {
		// 如果需要在新标签页中打开，添加 target="_blank" 属性
		if openInNewTabs {
			return template.Default(ctx).Link().SetURL(prefix + value.Value).OpenInNewTab().GetContent()
		} else {
			// 否则在当前标签页中打开
			return template.Default(ctx).Link().SetURL(prefix + value.Value).GetContent()
		}
	}
}

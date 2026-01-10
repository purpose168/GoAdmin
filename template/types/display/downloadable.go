package display

import (
	"html/template"

	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/template/types"
)

// Downloadable 可下载显示生成器
// 用于将字段值转换为可下载链接的可视化显示
// 继承自 BaseDisplayFnGenerator，提供基础的显示函数生成能力
// 用户可以点击下载按钮下载对应的文件
type Downloadable struct {
	types.BaseDisplayFnGenerator
}

// init 包初始化函数
// 在包加载时自动执行，将 Downloadable 类型注册到显示函数生成器注册表中
// 注册键名为 "downloadable"，可以通过该键名创建 Downloadable 实例
func init() {
	types.RegisterDisplayFnGenerator("downloadable", new(Downloadable))
}

// Get 获取字段过滤函数
// 根据传入的参数生成一个字段过滤函数，用于将字段值转换为可下载链接
//
// 参数：
//   - ctx: 上下文对象，包含请求相关的上下文信息
//   - args: 可变参数，包含以下内容：
//   - args[0]: []string 类型，URL 前缀数组
//   - 如果为空数组，直接使用字段值作为下载链接
//   - 如果包含元素，将第一个元素作为 URL 前缀，拼接到字段值前面
//
// 返回值：
//   - FieldFilterFn: 字段过滤函数，接收 FieldModel 对象并返回转换后的 HTML
//     HTML 包含一个下载链接和下载图标
//
// 使用示例：
//
//	// 示例1：直接使用字段值作为下载链接
//	display.Downloadable{}.Get(ctx, []string{})
//
//	// 示例2：添加 URL 前缀
//	display.Downloadable{}.Get(ctx, []string{"/uploads/"})
//	// 如果字段值为 "document.pdf"，生成的链接为 "/uploads/document.pdf"
//
//	// 示例3：使用完整的 CDN 前缀
//	display.Downloadable{}.Get(ctx, []string{"https://cdn.example.com/files/"})
//
// 注意事项：
//   - 参数 args[0] 必须是 []string 类型，否则会引发运行时 panic
//   - 生成的HTML包含一个下载链接（使用 Font Awesome 的 fa-download 图标）
//   - download 属性指定下载时的文件名
//   - target="_blank" 在新标签页中打开链接
//   - 需要引入 Font Awesome 图标库
//   - URL 前缀和字段值之间会自动拼接，不会添加额外的斜杠
func (d *Downloadable) Get(ctx *context.Context, args ...interface{}) types.FieldFilterFn {
	return func(value types.FieldModel) interface{} {
		// 从参数中提取 URL 前缀数组
		param := args[0].([]string)

		// 默认使用字段值作为下载链接
		u := value.Value

		// 如果提供了 URL 前缀，将前缀拼接到字段值前面
		if len(param) > 0 {
			u = param[0] + u
		}

		// 返回包含下载链接的 HTML
		// 使用 template.HTML 类型，避免 HTML 转义
		// HTML 结构：
		//   - <a> 标签：下载链接
		//     - href: 下载地址
		//     - download: 下载时的文件名（使用字段值）
		//     - target="_blank": 在新标签页中打开
		//     - class="text-muted": 文本颜色样式
		//   - <i> 标签：Font Awesome 的下载图标
		//   - value.Value：文件名显示
		return template.HTML(`
<a href="` + u + `" download="` + value.Value + `" target="_blank" class="text-muted">
	<i class="fa fa-download"></i> ` + value.Value + `
</a>
`)
	}
}

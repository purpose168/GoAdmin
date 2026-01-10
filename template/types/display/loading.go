package display

import (
	"html/template"

	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/template/types"
)

// Loading 加载状态显示生成器
// 用于将字段值转换为加载状态的图标显示
// 继承自 BaseDisplayFnGenerator，提供基础的显示函数生成能力
// 当字段值匹配指定值时，显示旋转的刷新图标
type Loading struct {
	types.BaseDisplayFnGenerator
}

// init 包初始化函数
// 在包加载时自动执行，将 Loading 类型注册到显示函数生成器注册表中
// 注册键名为 "loading"，可以通过该键名创建 Loading 实例
func init() {
	types.RegisterDisplayFnGenerator("loading", new(Loading))
}

// Get 获取字段过滤函数
// 根据传入的参数生成一个字段过滤函数，用于将字段值转换为加载状态显示
//
// 参数：
//   - ctx: 上下文对象，包含请求相关的上下文信息
//   - args: 可变参数，必须包含以下内容：
//   - args[0]: []string 类型，需要显示加载状态的字段值列表
//     当字段值等于列表中的任何一个值时，显示加载图标
//
// 返回值：
//   - FieldFilterFn: 字段过滤函数，接收 FieldModel 对象并返回转换后的 HTML
//     如果字段值匹配，返回旋转的刷新图标 HTML
//     如果字段值不匹配，返回原始字段值
//
// 使用示例：
//
//	// 示例1：当字段值为 "processing" 时显示加载图标
//	display.Loading{}.Get(ctx, []string{"processing"})
//
//	// 示例2：当字段值为 "loading" 或 "pending" 时显示加载图标
//	display.Loading{}.Get(ctx, []string{"loading", "pending"})
//
//	// 示例3：当字段值为 "1" 或 "true" 时显示加载图标
//	display.Loading{}.Get(ctx, []string{"1", "true"})
//
// 注意事项：
//   - 参数 args[0] 必须是 []string 类型，否则会引发运行时 panic
//   - 加载图标使用 Font Awesome 的 fa-refresh 图标，并带有 fa-spin 旋转动画
//   - 加载图标使用 text-primary 类显示为蓝色
//   - 如果字段值不匹配列表中的任何值，只显示原始字段值
//   - 需要引入 Font Awesome 图标库和 CSS 动画
func (l *Loading) Get(ctx *context.Context, args ...interface{}) types.FieldFilterFn {
	return func(value types.FieldModel) interface{} {
		// 从参数中提取需要显示加载状态的字段值列表
		param := args[0].([]string)

		// 遍历字段值列表，检查当前字段值是否匹配
		for i := 0; i < len(param); i++ {
			// 如果找到匹配的字段值，返回加载图标 HTML
			if value.Value == param[i] {
				// 使用 template.HTML 类型，避免 HTML 转义
				// 图标说明：
				//   - fa-refresh: Font Awesome 的刷新图标
				//   - fa-spin: 旋转动画类
				//   - text-primary: Bootstrap 的主要颜色类（蓝色）
				return template.HTML(`<i class="fa fa-refresh fa-spin text-primary"></i>`)
			}
		}

		// 如果字段值不匹配列表中的任何值，返回原始字段值
		return value.Value
	}
}

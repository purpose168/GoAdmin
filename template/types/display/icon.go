package display

import (
	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/template/icon"
	"github.com/purpose168/GoAdmin/template/types"
)

// Icon 图标显示生成器
// 用于将字段值转换为对应的图标显示
// 继承自 BaseDisplayFnGenerator，提供基础的显示函数生成能力
// 支持根据不同的字段值显示不同的图标
type Icon struct {
	types.BaseDisplayFnGenerator
}

// init 包初始化函数
// 在包加载时自动执行，将 Icon 类型注册到显示函数生成器注册表中
// 注册键名为 "icon"，可以通过该键名创建 Icon 实例
func init() {
	types.RegisterDisplayFnGenerator("icon", new(Icon))
}

// Get 获取字段过滤函数
// 根据传入的参数生成一个字段过滤函数，用于将字段值转换为图标显示
//
// 参数：
//   - ctx: 上下文对象，包含请求相关的上下文信息
//   - args: 可变参数，包含以下内容：
//   - args[0]: map[string]string 类型，字段值到图标类的映射
//     键是字段值，值是对应的图标类名（如 "fa-check", "fa-times", "fa-warning" 等）
//   - args[1]（可选）: string 类型，默认图标类名
//     当字段值不匹配 args[0] 中的任何键时，使用此图标
//
// 返回值：
//   - FieldFilterFn: 字段过滤函数，接收 FieldModel 对象并返回转换后的图标 HTML
//
// 使用示例：
//
//	// 示例1：根据状态显示不同图标
//	display.Icon{}.Get(ctx, map[string]string{
//	    "success": "fa-check",      // 成功状态显示勾选图标
//	    "error":   "fa-times",      // 错误状态显示叉号图标
//	    "warning": "fa-exclamation", // 警告状态显示感叹号图标
//	})
//
//	// 示例2：指定默认图标
//	display.Icon{}.Get(ctx, map[string]string{
//	    "yes": "fa-check",
//	    "no":  "fa-times",
//	}, "fa-question") // 其他值显示问号图标
//
// 注意事项：
//   - 参数 args[0] 必须是 map[string]string 类型，否则会引发运行时 panic
//   - 图标类名基于 Font Awesome 图标库，需要引入 Font Awesome CSS
//   - 如果没有匹配的图标且没有默认图标，只显示原始字段值
//   - icon.Icon 函数会生成完整的图标 HTML 结构
func (i *Icon) Get(ctx *context.Context, args ...interface{}) types.FieldFilterFn {
	return func(value types.FieldModel) interface{} {
		// 从参数中提取字段值到图标类的映射
		icons := args[0].(map[string]string)

		// 初始化默认图标为空字符串
		defaultIcon := ""

		// 如果提供了第二个参数，使用该参数作为默认图标
		if len(args) > 1 {
			defaultIcon = args[1].(string)
		}

		// 遍历图标映射，查找当前字段值对应的图标
		for k, iconClass := range icons {
			// 如果找到匹配的字段值，生成对应的图标 HTML
			if k == value.Value {
				return icon.Icon(iconClass)
			}
		}

		// 如果没有找到匹配的字段值，但设置了默认图标，使用默认图标
		if defaultIcon != "" {
			return icon.Icon(defaultIcon)
		}

		// 如果没有匹配的图标且没有默认图标，只返回原始字段值
		return value.Value
	}
}

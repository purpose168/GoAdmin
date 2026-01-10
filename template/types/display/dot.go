package display

import (
	"html/template"

	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/template/types"
)

// Dot 点状标记显示生成器
// 用于将字段值转换为带有彩色圆点的可视化显示
// 继承自 BaseDisplayFnGenerator，提供基础的显示函数生成能力
// 支持根据不同的字段值显示不同颜色的圆点标记
type Dot struct {
	types.BaseDisplayFnGenerator
}

// init 包初始化函数
// 在包加载时自动执行，将 Dot 类型注册到显示函数生成器注册表中
// 注册键名为 "dot"，可以通过该键名创建 Dot 实例
func init() {
	types.RegisterDisplayFnGenerator("dot", new(Dot))
}

// Get 获取字段过滤函数
// 根据传入的参数生成一个字段过滤函数，用于将字段值转换为带有彩色圆点的显示
//
// 参数：
//   - ctx: 上下文对象，包含请求相关的上下文信息
//   - args: 可变参数，包含以下内容：
//   - args[0]: map[string]types.FieldDotColor 类型，字段值到颜色的映射
//     键是字段值，值是对应的颜色标识（如 "success", "danger", "warning", "info" 等）
//   - args[1]（可选）: types.FieldDotColor 类型，默认颜色
//     当字段值不匹配 args[0] 中的任何键时，使用此颜色
//
// 返回值：
//   - FieldFilterFn: 字段过滤函数，接收 FieldModel 对象并返回转换后的 HTML
//     HTML 包含一个彩色圆点和原始字段值
//
// 使用示例：
//
//	// 示例1：根据状态显示不同颜色的圆点
//	display.Dot{}.Get(ctx, map[string]types.FieldDotColor{
//	    "active":  "success",  // 活跃状态显示绿色圆点
//	    "inactive": "danger",  // 非活跃状态显示红色圆点
//	    "pending": "warning",  // 待处理状态显示黄色圆点
//	})
//
//	// 示例2：指定默认颜色
//	display.Dot{}.Get(ctx, map[string]types.FieldDotColor{
//	    "yes": "success",
//	    "no":  "danger",
//	}, "info") // 其他值显示蓝色圆点
//
// 注意事项：
//   - 参数 args[0] 必须是 map[string]types.FieldDotColor 类型，否则会引发运行时 panic
//   - 圆点样式使用 Bootstrap 的 label 类，需要引入 Bootstrap CSS
//   - 圆点大小为 8x8 像素，圆形显示
//   - 圆点和字段值之间有两个空格分隔
//   - 如果没有匹配的颜色且没有默认颜色，只显示原始字段值
func (d *Dot) Get(ctx *context.Context, args ...interface{}) types.FieldFilterFn {
	return func(value types.FieldModel) interface{} {
		// 从参数中提取字段值到颜色的映射
		icons := args[0].(map[string]types.FieldDotColor)

		// 初始化默认颜色为空字符串
		defaultDot := types.FieldDotColor("")

		// 如果提供了第二个参数，使用该参数作为默认颜色
		if len(args) > 1 {
			defaultDot = args[1].(types.FieldDotColor)
		}

		// 遍历颜色映射，查找当前字段值对应的颜色
		for k, style := range icons {
			// 如果找到匹配的字段值，生成带有对应颜色圆点的 HTML
			if k == value.Value {
				// 使用 template.HTML 类型，避免 HTML 转义
				// 圆点样式：
				//   - width/height: 8px - 圆点大小
				//   - padding: 0 - 无内边距
				//   - border-radius: 50% - 圆形
				//   - display: inline-block - 行内块级元素
				return template.HTML(`<span class="label-`+style+`"
					style="width: 8px;height: 8px;padding: 0;border-radius: 50%;display: inline-block;">
					</span>&nbsp;&nbsp;`) +
					template.HTML(value.Value)
			}
		}

		// 如果没有找到匹配的字段值，但设置了默认颜色，使用默认颜色
		if defaultDot != "" {
			return template.HTML(`<span class="label-`+defaultDot+`"
					style="width: 8px;height: 8px;padding: 0;border-radius: 50%;display: inline-block;">
					</span>&nbsp;&nbsp;`) +
				template.HTML(value.Value)
		}

		// 如果没有匹配的颜色且没有默认颜色，只返回原始字段值
		return value.Value
	}
}

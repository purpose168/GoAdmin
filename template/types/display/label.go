package display

import (
	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/template"
	"github.com/purpose168/GoAdmin/template/types"
)

// Label 标签显示生成器
// 用于将字段值转换为标签样式的显示
// 继承自 BaseDisplayFnGenerator，提供基础的显示函数生成能力
// 支持自定义标签颜色和类型
type Label struct {
	types.BaseDisplayFnGenerator
}

// init 包初始化函数
// 在包加载时自动执行，将 Label 类型注册到显示函数生成器注册表中
// 注册键名为 "label"，可以通过该键名创建 Label 实例
func init() {
	types.RegisterDisplayFnGenerator("label", new(Label))
}

// Get 获取字段过滤函数
// 根据传入的参数生成一个字段过滤函数，用于将字段值转换为标签显示
//
// 参数：
//   - ctx: 上下文对象，包含请求相关的上下文信息
//   - args: 可变参数，包含以下内容：
//   - args[0]: []types.FieldLabelParam 类型，标签参数数组
//   - 如果为空数组，使用默认样式（success 类型）
//   - 如果包含1个元素，使用该元素的 Color 和 Type 属性
//
// 返回值：
//   - FieldFilterFn: 字段过滤函数，接收 FieldModel 对象并返回转换后的标签 HTML
//
// 使用示例：
//
//	// 示例1：使用默认样式（success 类型）
//	display.Label{}.Get(ctx, []types.FieldLabelParam{})
//
//	// 示例2：自定义标签颜色和类型
//	display.Label{}.Get(ctx, []types.FieldLabelParam{
//	    {
//	        Color: "#ff0000",  // 自定义颜色
//	        Type:  "danger",   // Bootstrap 标签类型
//	    },
//	})
//
// 注意事项：
//   - 参数 args[0] 必须是 []types.FieldLabelParam 类型，否则会引发运行时 panic
//   - 使用 template.Default(ctx).Label() 构建标签组件
//   - 默认类型为 "success"，显示为绿色标签
//   - 标签类型支持 Bootstrap 的 label 类型（success、info、warning、danger、primary）
//   - 如果参数数组长度大于1，返回空字符串
func (label *Label) Get(ctx *context.Context, args ...interface{}) types.FieldFilterFn {
	return func(value types.FieldModel) interface{} {
		// 从参数中提取标签参数数组
		params := args[0].([]types.FieldLabelParam)

		// 如果参数数组为空，使用默认样式（success 类型）
		if len(params) == 0 {
			return template.Default(ctx).Label().
				SetContent(template.HTML(value.Value)).
				SetType("success").
				GetContent()
		} else if len(params) == 1 {
			// 如果参数数组包含1个元素，使用该元素的颜色和类型
			return template.Default(ctx).Label().
				SetContent(template.HTML(value.Value)).
				SetColor(params[0].Color).
				SetType(params[0].Type).
				GetContent()
		}

		// 如果参数数组长度大于1，返回空字符串
		return ""
	}
}

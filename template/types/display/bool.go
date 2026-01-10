package display

import (
	"strings"

	"github.com/GoAdminGroup/html"
	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/template/icon"
	"github.com/purpose168/GoAdmin/template/types"
)

// Bool 布尔值显示生成器
// 用于将布尔值或类似布尔值的字段转换为可视化的图标显示
// 继承自 BaseDisplayFnGenerator，提供基础的显示函数生成能力
type Bool struct {
	types.BaseDisplayFnGenerator
}

// init 包初始化函数
// 在包加载时自动执行，将 Bool 类型注册到显示函数生成器注册表中
// 注册键名为 "bool"，可以通过该键名创建 Bool 实例
func init() {
	types.RegisterDisplayFnGenerator("bool", new(Bool))
}

// Get 获取字段过滤函数
// 根据传入的参数生成一个字段过滤函数，用于将字段值转换为布尔图标显示
//
// 参数：
//   - ctx: 上下文对象，包含请求相关的上下文信息
//   - args: 可变参数，支持以下格式：
//   - 无参数：将 "0" 或 "false" 显示为红色叉号，其他值显示为绿色勾号
//   - 1个参数：当字段值等于该参数时显示绿色勾号，否则显示红色叉号
//   - 2个参数：第一个参数表示"通过"的值，第二个参数表示"失败"的值
//
// 返回值：
//   - FieldFilterFn: 字段过滤函数，接收 FieldModel 对象并返回转换后的显示内容
//
// 使用示例：
//
//	// 示例1：默认模式，0/false显示为叉号，其他显示为勾号
//	display.Bool{}.Get(ctx)
//
//	// 示例2：指定通过值，字段值为"active"时显示勾号
//	display.Bool{}.Get(ctx, "active")
//
//	// 示例3：指定通过和失败值，"yes"显示勾号，"no"显示叉号
//	display.Bool{}.Get(ctx, "yes", "no")
//
// 注意事项：
//   - 参数 args[0] 必须是 []string 类型，否则会引发运行时 panic
//   - 当字段值不匹配任何条件时，返回空字符串
//   - 使用了 html.Style 来设置图标颜色，绿色表示通过，红色表示失败
func (b *Bool) Get(ctx *context.Context, args ...interface{}) types.FieldFilterFn {
	return func(value types.FieldModel) interface{} {
		// 从参数中提取字符串数组
		params := args[0].([]string)

		// 创建通过状态的图标：绿色勾号
		// 使用 icon.Check 勾选图标，并设置颜色为绿色
		pass := icon.IconWithStyle(icon.Check, html.Style{
			"color": "green",
		})

		// 创建失败状态的图标：红色叉号
		// 使用 icon.Remove 移除图标，并设置颜色为红色
		fail := icon.IconWithStyle(icon.Remove, html.Style{
			"color": "red",
		})

		// 模式1：无参数，使用默认布尔值判断
		// 当字段值为 "0" 或 "false"（不区分大小写）时显示失败图标
		// 其他值显示通过图标
		if len(params) == 0 {
			if value.Value == "0" || strings.ToLower(value.Value) == "false" {
				return fail
			}
			return pass
		} else if len(params) == 1 {
			// 模式2：1个参数，指定通过值
			// 当字段值等于参数值时显示通过图标，否则显示失败图标
			if value.Value == params[0] {
				return pass
			}
			return fail
		} else {
			// 模式3：2个参数，分别指定通过和失败值
			// params[0] 表示通过值，params[1] 表示失败值
			// 当字段值匹配通过值时显示通过图标
			if value.Value == params[0] {
				return pass
			}
			// 当字段值匹配失败值时显示失败图标
			if value.Value == params[1] {
				return fail
			}
		}

		// 字段值不匹配任何条件时，返回空字符串
		// 这种情况通常发生在模式3中，字段值既不是通过值也不是失败值
		return ""
	}
}

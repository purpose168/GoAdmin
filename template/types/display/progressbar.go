package display

import (
	"fmt"
	"html/template"
	"strconv"

	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/template/types"
)

// ProgressBar 进度条显示生成器
// 用于将字段值转换为进度条的可视化显示
// 继承自 BaseDisplayFnGenerator，提供基础的显示函数生成能力
// 支持自定义进度条样式、大小和最大值
type ProgressBar struct {
	types.BaseDisplayFnGenerator
}

// init 包初始化函数
// 在包加载时自动执行，将 ProgressBar 类型注册到显示函数生成器注册表中
// 注册键名为 "progressbar"，可以通过该键名创建 ProgressBar 实例
func init() {
	types.RegisterDisplayFnGenerator("progressbar", new(ProgressBar))
}

// Get 获取字段过滤函数
// 根据传入的参数生成一个字段过滤函数，用于将字段值转换为进度条显示
//
// 参数：
//   - ctx: 上下文对象，包含请求相关的上下文信息
//   - args: 可变参数，必须包含以下内容：
//   - args[0]: []types.FieldProgressBarData 类型，进度条参数数组
//   - 如果为空数组，使用默认样式（primary 类型、sm 大小、最大值 100）
//   - 如果包含1个元素，使用该元素的 Style、Size 和 Max 属性
//   - Style: 进度条样式（primary、success、info、warning、danger）
//   - Size: 进度条大小（sm、md、lg）
//   - Max: 进度条最大值（默认为 100）
//
// 返回值：
//   - FieldFilterFn: 字段过滤函数，接收 FieldModel 对象并返回转换后的进度条 HTML
//
// 使用示例：
//
//	// 示例1：使用默认样式
//	display.ProgressBar{}.Get(ctx, []types.FieldProgressBarData{})
//
//	// 示例2：自定义样式和大小
//	display.ProgressBar{}.Get(ctx, []types.FieldProgressBarData{
//	    {
//	        Style: "success",  // 绿色进度条
//	        Size:  "lg",      // 大尺寸
//	        Max:   200,      // 最大值为 200
//	    },
//	})
//
//	// 示例3：自定义最大值
//	display.ProgressBar{}.Get(ctx, []types.FieldProgressBarData{
//	    {
//	        Style: "warning",
//	        Max:   1000,
//	    },
//	})
//
// 注意事项：
//   - 参数 args[0] 必须是 []types.FieldProgressBarData 类型，否则会引发运行时 panic
//   - 字段值必须是表示进度的数字字符串，否则解析会失败
//   - 进度条基于 Bootstrap 的 progress 组件，需要引入 Bootstrap CSS
//   - 进度条显示百分比和可视化进度条两部分
//   - 如果字段值解析失败，会显示 0% 的进度条
func (p *ProgressBar) Get(ctx *context.Context, args ...interface{}) types.FieldFilterFn {
	return func(value types.FieldModel) interface{} {
		// 从参数中提取进度条参数数组
		param := args[0].([]types.FieldProgressBarData)

		// 初始化默认样式为 primary（蓝色）
		style := "primary"
		// 初始化默认大小为 sm（小尺寸）
		size := "sm"
		// 初始化默认最大值为 100
		max := 100

		// 如果提供了参数，使用参数中的自定义值
		if len(param) > 0 {
			// 如果指定了样式，使用自定义样式
			if param[0].Style != "" {
				style = param[0].Style
			}
			// 如果指定了大小，使用自定义大小
			if param[0].Size != "" {
				size = param[0].Size
			}
			// 如果指定了最大值，使用自定义最大值
			if param[0].Max != 0 {
				max = param[0].Max
			}
		}

		// 将字段值（字符串）解析为整数
		// 字段值应该是表示进度的数字字符串
		base, _ := strconv.Atoi(value.Value)

		// 计算进度百分比
		// 公式：(当前值 / 最大值) * 100
		// 使用 fmt.Sprintf 格式化为整数，不保留小数
		per := fmt.Sprintf("%.0f", float32(base)/float32(max)*100)

		// 返回进度条 HTML
		// 使用 template.HTML 类型，避免 HTML 转义
		// HTML 结构：
		//   - row: Bootstrap 行容器
		//   - col-sm-3: 显示百分比文本的列
		//   - progress: Bootstrap 进度条容器
		//   - progress-bar: 实际的进度条元素
		return template.HTML(`
<div class="row" style="min-width: 100px;">
	<span class="col-sm-3" style="color:#777;width: 60px">` + per + `%</span>
	<div class="progress progress-` + size + ` col-sm-9" style="padding-left: 0;width: 100px;margin-left: -13px;">
		<div class="progress-bar progress-bar-` + style + `" role="progressbar" aria-valuenow="1" 
			aria-valuemin="0" aria-valuemax="` + strconv.Itoa(max) + `" style="width: ` + per + `%">
		</div>
	</div>
</div>`)
	}
}

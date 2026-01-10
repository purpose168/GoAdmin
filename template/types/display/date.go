package display

import (
	"strconv"
	"time"

	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/template/types"
)

// Date 日期显示生成器
// 用于将时间戳字段转换为格式化的日期字符串显示
// 继承自 BaseDisplayFnGenerator，提供基础的显示函数生成能力
// 支持自定义日期格式，使用 Go 标准库的 time.Format 方法
type Date struct {
	types.BaseDisplayFnGenerator
}

// init 包初始化函数
// 在包加载时自动执行，将 Date 类型注册到显示函数生成器注册表中
// 注册键名为 "date"，可以通过该键名创建 Date 实例
func init() {
	types.RegisterDisplayFnGenerator("date", new(Date))
}

// Get 获取字段过滤函数
// 根据传入的参数生成一个字段过滤函数，用于将时间戳转换为格式化的日期字符串
//
// 参数：
//   - ctx: 上下文对象，包含请求相关的上下文信息
//   - args: 可变参数，必须包含以下内容：
//   - args[0]: string 类型，日期格式字符串
//     使用 Go 标准库的日期格式化语法
//     常用格式示例：
//   - "2006-01-02"：年-月-日
//   - "2006-01-02 15:04:05"：年-月-日 时:分:秒
//   - "2006/01/02"：年/月/日
//   - "15:04"：时:分
//
// 返回值：
//   - FieldFilterFn: 字段过滤函数，接收 FieldModel 对象并返回格式化后的日期字符串
//
// 使用示例：
//
//	// 示例1：显示为 "2024-01-06" 格式
//	display.Date{}.Get(ctx, "2006-01-02")
//
//	// 示例2：显示为 "2024-01-06 14:30:45" 格式
//	display.Date{}.Get(ctx, "2006-01-02 15:04:05")
//
//	// 示例3：显示为 "2024/01/02" 格式
//	display.Date{}.Get(ctx, "2006/01/02")
//
// 注意事项：
//   - 参数 args[0] 必须是 string 类型，否则会引发运行时 panic
//   - 字段值必须是 Unix 时间戳（秒级），否则解析会失败
//   - Go 的日期格式化使用特定的参考时间 "Mon Jan 2 15:04:05 MST 2006"
//     对应的数字为：01=月, 02=日, 03=时(12小时制), 04=分, 05=秒, 06=年, 15=时(24小时制)
//   - 如果时间戳解析失败，会返回零值时间（1970-01-01 00:00:00 UTC）
func (d *Date) Get(ctx *context.Context, args ...interface{}) types.FieldFilterFn {
	return func(value types.FieldModel) interface{} {
		// 从参数中提取日期格式字符串
		format := args[0].(string)

		// 将字段值（字符串）解析为 64 位整数
		// 字段值应该是 Unix 时间戳（秒级）
		ts, _ := strconv.ParseInt(value.Value, 10, 64)

		// 将 Unix 时间戳转换为 time.Time 对象
		// time.Unix 的第二个参数是纳秒部分，这里设置为 0
		tm := time.Unix(ts, 0)

		// 使用指定的格式格式化时间并返回
		return tm.Format(format)
	}
}

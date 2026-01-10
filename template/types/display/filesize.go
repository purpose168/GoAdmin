package display

import (
	"strconv"

	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/modules/utils"
	"github.com/purpose168/GoAdmin/template/types"
)

// FileSize 文件大小显示生成器
// 用于将字节大小的字段值转换为人类可读的文件大小格式
// 继承自 BaseDisplayFnGenerator，提供基础的显示函数生成能力
// 支持自动转换为 B、KB、MB、GB、TB 等单位
type FileSize struct {
	types.BaseDisplayFnGenerator
}

// init 包初始化函数
// 在包加载时自动执行，将 FileSize 类型注册到显示函数生成器注册表中
// 注册键名为 "filesize"，可以通过该键名创建 FileSize 实例
func init() {
	types.RegisterDisplayFnGenerator("filesize", new(FileSize))
}

// Get 获取字段过滤函数
// 根据传入的参数生成一个字段过滤函数，用于将字节大小转换为人类可读格式
//
// 参数：
//   - ctx: 上下文对象，包含请求相关的上下文信息
//   - args: 可变参数，当前实现不使用任何参数
//
// 返回值：
//   - FieldFilterFn: 字段过滤函数，接收 FieldModel 对象并返回转换后的文件大小字符串
//     返回格式示例：
//   - "0 B"
//   - "1.5 KB"
//   - "2.3 MB"
//   - "1.2 GB"
//   - "500 TB"
//
// 使用示例：
//
//	// 示例：将字节大小转换为可读格式
//	display.FileSize{}.Get(ctx)
//
// 转换规则：
//   - 0 - 1023 字节：显示为 B（字节）
//   - 1024 - 1048575 字节：显示为 KB（千字节）
//   - 1048576 - 1073741823 字节：显示为 MB（兆字节）
//   - 1073741824 - 1099511627775 字节：显示为 GB（吉字节）
//   - 1099511627776 字节及以上：显示为 TB（太字节）
//
// 注意事项：
//   - 字段值必须是表示字节大小的字符串，否则解析会失败
//   - 使用 utils.FileSize 函数进行转换，该函数自动选择合适的单位
//   - 转换结果保留一位小数，便于阅读
//   - 如果字段值解析失败，会返回 "0 B"
//   - 单位换算基于 1024 进制（1 KB = 1024 B）
func (f *FileSize) Get(ctx *context.Context, args ...interface{}) types.FieldFilterFn {
	return func(value types.FieldModel) interface{} {
		// 将字段值（字符串）解析为 64 位无符号整数
		// 字段值应该是表示字节大小的数字字符串
		size, _ := strconv.ParseUint(value.Value, 10, 64)

		// 调用 utils.FileSize 函数将字节大小转换为人类可读格式
		// 该函数会自动选择合适的单位（B、KB、MB、GB、TB）
		return utils.FileSize(size)
	}
}

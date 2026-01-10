package display

import (
	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/template"
	"github.com/purpose168/GoAdmin/template/types"
)

// Image 图片显示生成器
// 用于将字段值转换为图片显示
// 继承自 BaseDisplayFnGenerator，提供基础的显示函数生成能力
// 支持自定义图片尺寸和URL前缀
type Image struct {
	types.BaseDisplayFnGenerator
}

// init 包初始化函数
// 在包加载时自动执行，将 Image 类型注册到显示函数生成器注册表中
// 注册键名为 "image"，可以通过该键名创建 Image 实例
func init() {
	types.RegisterDisplayFnGenerator("image", new(Image))
}

// Get 获取字段过滤函数
// 根据传入的参数生成一个字段过滤函数，用于将字段值转换为图片显示
//
// 参数：
//   - ctx: 上下文对象，包含请求相关的上下文信息
//   - args: 可变参数，必须包含以下内容：
//   - args[0]: string 类型，图片宽度（如 "100px", "50%"）
//   - args[1]: string 类型，图片高度（如 "100px", "50%"）
//   - args[2]: []string 类型，URL 前缀数组
//   - 如果为空数组，直接使用字段值作为图片URL
//   - 如果包含元素，将第一个元素作为 URL 前缀，拼接到字段值前面
//
// 返回值：
//   - FieldFilterFn: 字段过滤函数，接收 FieldModel 对象并返回转换后的图片 HTML
//
// 使用示例：
//
//	// 示例1：直接使用字段值作为图片URL
//	display.Image{}.Get(ctx, "100px", "100px", []string{})
//
//	// 示例2：添加 URL 前缀
//	display.Image{}.Get(ctx, "200px", "200px", []string{"/uploads/"})
//	// 如果字段值为 "photo.jpg"，生成的图片URL为 "/uploads/photo.jpg"
//
//	// 示例3：使用百分比尺寸
//	display.Image{}.Get(ctx, "50%", "auto", []string{"https://cdn.example.com/images/"})
//
// 注意事项：
//   - 参数 args[0]、args[1] 必须是 string 类型，否则会引发运行时 panic
//   - 参数 args[2] 必须是 []string 类型，否则会引发运行时 panic
//   - 如果字段值为空字符串，返回空字符串
//   - 使用 template.Default(ctx).Image() 构建图片组件
//   - URL 前缀和字段值之间会自动拼接，不会添加额外的斜杠
func (image *Image) Get(ctx *context.Context, args ...interface{}) types.FieldFilterFn {
	// 从参数中提取 URL 前缀数组
	param := args[2].([]string)

	return func(value types.FieldModel) interface{} {
		// 如果字段值为空，返回空字符串
		if value.Value == "" {
			return ""
		}

		// 如果提供了 URL 前缀，将前缀拼接到字段值前面
		if len(param) > 0 {
			// 使用模板构建图片组件，设置宽度、高度和带前缀的图片源
			return template.Default(ctx).Image().SetWidth(args[0].(string)).SetHeight(args[1].(string)).
				SetSrc(template.HTML(param[0] + value.Value)).GetContent()
		} else {
			// 使用模板构建图片组件，设置宽度、高度和原始字段值作为图片源
			return template.Default(ctx).Image().SetWidth(args[0].(string)).SetHeight(args[1].(string)).
				SetSrc(template.HTML(value.Value)).GetContent()
		}
	}
}

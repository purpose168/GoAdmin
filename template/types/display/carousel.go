package display

import (
	"html/template"
	"strconv"

	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/template/types"
)

// Carousel 轮播图显示生成器
// 用于将字段值转换为可视化的轮播图组件显示
// 继承自 BaseDisplayFnGenerator，提供基础的显示函数生成能力
// 基于 Bootstrap 的 carousel 组件实现，支持图片轮播展示
type Carousel struct {
	types.BaseDisplayFnGenerator
}

// init 包初始化函数
// 在包加载时自动执行，将 Carousel 类型注册到显示函数生成器注册表中
// 注册键名为 "carousel"，可以通过该键名创建 Carousel 实例
func init() {
	types.RegisterDisplayFnGenerator("carousel", new(Carousel))
}

// Get 获取字段过滤函数
// 根据传入的参数生成一个字段过滤函数，用于将字段值转换为轮播图显示
//
// 参数：
//   - ctx: 上下文对象，包含请求相关的上下文信息
//   - args: 可变参数，必须包含以下内容：
//   - args[0]: FieldGetImgArrFn 类型，用于从字段值中提取图片数组的函数
//   - args[1]: []int 类型，指定轮播图的宽度和高度
//   - 如果为空数组，使用默认宽度 300px 和高度 200px
//   - 如果包含1个元素，设置宽度，高度使用默认值 200px
//   - 如果包含2个元素，分别设置宽度和高度
//
// 返回值：
//   - FieldFilterFn: 字段过滤函数，接收 FieldModel 对象并返回转换后的轮播图 HTML
//
// 使用示例：
//
//	// 示例1：使用默认尺寸（300x200）
//	display.Carousel{}.Get(ctx, getImages, []int{})
//
//	// 示例2：自定义宽度为 500px，高度使用默认值
//	display.Carousel{}.Get(ctx, getImages, []int{500})
//
//	// 示例3：自定义宽度和高度为 800x600
//	display.Carousel{}.Get(ctx, getImages, []int{800, 600})
//
// 注意事项：
//   - 参数 args[0] 必须是 FieldGetImgArrFn 类型，否则会引发运行时 panic
//   - 参数 args[1] 必须是 []int 类型，否则会引发运行时 panic
//   - 生成的轮播图基于 Bootstrap 框架，需要引入 Bootstrap CSS 和 JS
//   - 每个轮播图的 ID 会根据字段 ID 自动生成，确保唯一性
//   - 图片路径由 FieldGetImgArrFn 函数提供，可以是相对路径或绝对路径
func (c *Carousel) Get(ctx *context.Context, args ...interface{}) types.FieldFilterFn {
	return func(value types.FieldModel) interface{} {
		// 从参数中提取图片数组获取函数
		// 该函数负责从字段值中解析出图片URL数组
		fn := args[0].(types.FieldGetImgArrFn)

		// 从参数中提取尺寸数组
		// size[0] 表示宽度，size[1] 表示高度
		size := args[1].([]int)

		// 设置默认宽度为 300px
		width := "300"
		// 设置默认高度为 200px
		height := "200"

		// 如果提供了宽度参数，使用自定义宽度
		if len(size) > 0 {
			width = strconv.Itoa(size[0])
		}

		// 如果提供了高度参数，使用自定义高度
		if len(size) > 1 {
			height = strconv.Itoa(size[1])
		}

		// 调用图片数组获取函数，从字段值中提取图片URL列表
		images := fn(value.Value)

		// 初始化轮播图指示器（底部的小圆点）HTML字符串
		indicators := ""
		// 初始化轮播图内容项 HTML字符串
		items := ""
		// 初始化活动状态标记，用于标记当前显示的图片
		active := ""

		// 遍历所有图片，生成轮播图的各个组件
		for i, img := range images {
			// 生成指示器列表项
			// data-target: 指向轮播图容器的ID
			// data-slide-to: 指示器对应的图片索引
			indicators += `<li data-target="#carousel-value-` + value.ID + `" data-slide-to="` +
				strconv.Itoa(i) + `" class=""></li>`

			// 第一张图片设置为活动状态（active）
			if i == 0 {
				active = " active"
			} else {
				// 其他图片不设置活动状态
				active = ""
			}

			// 生成轮播图内容项
			// item: Bootstrap 的轮播图项类
			// active: 标记当前显示的图片
			// img: 图片标签，包含图片URL和样式
			// carousel-caption: 图片说明区域（当前为空）
			items += `<div class="item` + active + `">
            <img src="` + img + `" alt=""
style="max-width:` + width + `px;max-height:` + height + `px;display: block;margin-left: auto;margin-right: auto;" />
            <div class="carousel-caption"></div>
        </div>`
		}

		// 返回完整的轮播图 HTML
		// 使用 template.HTML 类型，避免 HTML 转义
		// 轮播图组件包含：
		//   - 指示器（底部小圆点）
		//   - 轮播内容（图片列表）
		//   - 左右控制按钮（切换上一张/下一张）
		return template.HTML(`
<div id="carousel-value-` + value.ID + `" class="carousel slide" data-ride="carousel" width="` + width + `" height="` + height + `"
style="padding: 5px;border: 1px solid #f4f4f4;background-color:white;width:` + width + `px;">
    <ol class="carousel-indicators">
		` + indicators + `
    </ol>
    <div class="carousel-inner">
       ` + items + `
    </div>
    <a class="left carousel-control" href="#carousel-value-` + value.ID + `" data-slide="prev">
        <span class="fa fa-angle-left"></span>
    </a>
    <a class="right carousel-control" href="#carousel-value-` + value.ID + `" data-slide="next">
        <span class="fa fa-angle-right"></span>
    </a>
</div>
`)
	}
}

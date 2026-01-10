// Package tables 提供仪表板页面内容的生成功能
//
// 本包实现了 GoAdmin 管理后台的仪表板界面构建和管理，提供以下功能：
//   - GetContent: 获取仪表板页面的完整内容
//
// 核心概念：
//   - 仪表板: 管理后台的主页面，展示关键指标和数据概览
//   - UI组件: 各种可复用的界面元素，如信息框、表格、图表等
//   - 响应式布局: 根据屏幕尺寸自动调整的布局系统
//   - 模板系统: GoAdmin 的模板渲染系统
//   - 图表系统: 基于 Chart.js 的数据可视化组件
//   - 组件库: AdminLTE 主题提供的各种 UI 组件
//
// 技术栈：
//   - GoAdmin Template: 模板系统，提供组件渲染功能
//   - AdminLTE Components: UI 组件库，提供各种界面组件
//   - Chart.js: 图表库，提供数据可视化功能
//   - Bootstrap: 前端框架，提供基础样式和布局
//   - HTML Template: Go 标准模板引擎
//
// UI组件类型：
//   - Info Box: 信息框，显示关键业务指标
//   - Table: 表格，展示结构化数据
//   - Product List: 产品列表，展示产品信息
//   - Chart: 图表，展示数据趋势
//   - Small Box: 小方框，快速统计数据
//   - Pie Chart: 饼图，展示占比数据
//   - Tabs: 标签页，多标签内容展示
//   - Popup: 弹窗，模态对话框
//   - Progress Group: 进度条组，展示任务进度
//   - Description: 描述组件，展示数据变化
//
// 使用场景：
//   - 数据可视化: 在仪表板中展示关键数据和指标
//   - 业务监控: 实时监控业务运行状态
//   - 数据分析: 通过图表分析数据趋势
//   - 用户管理: 展示用户信息和统计数据
//   - 销售分析: 展示销售数据和趋势
//   - 系统监控: 监控系统性能和资源使用
//
// 配置说明：
//   - 响应式断点: XS(超小屏幕)、SM(小屏幕)、MD(中等屏幕)、LG(大屏幕)
//   - 颜色主题: 支持多种颜色主题（aqua、red、green、yellow、blue、danger、warning等）
//   - 图标库: 支持 Ion Icons 和 Font Awesome 图标
//   - 图表类型: 支持折线图、饼图、柱状图等
//
// 注意事项：
//   - 需要正确引入 AdminLTE 主题和 Chart.js 库
//   - 组件的 ID 必须唯一，避免冲突
//   - 响应式布局需要考虑不同屏幕尺寸的显示效果
//   - 图表数据需要正确格式化
//   - HTML 内容需要进行转义，防止 XSS 攻击
//
// 作者: GoAdmin Team
// 创建日期: 2019-01-01
// 版本: 1.0.0
package tables

import (
	"html/template" // Go 标准模板包，提供 HTML 模板渲染功能

	"github.com/purpose168/GoAdmin-themes/adminlte/components/chart_legend"   // 图例组件包，提供图表图例功能
	"github.com/purpose168/GoAdmin-themes/adminlte/components/description"    // 描述组件包，提供数据描述功能
	"github.com/purpose168/GoAdmin-themes/adminlte/components/infobox"        // 信息框组件包，提供信息框功能
	"github.com/purpose168/GoAdmin-themes/adminlte/components/productlist"    // 产品列表组件包，提供产品列表功能
	"github.com/purpose168/GoAdmin-themes/adminlte/components/progress_group" // 进度条组组件包，提供进度条功能
	"github.com/purpose168/GoAdmin-themes/adminlte/components/smallbox"       // 小方框组件包，提供小方框功能
	"github.com/purpose168/GoAdmin/context"                                   // 上下文包，提供请求上下文管理功能
	tmpl "github.com/purpose168/GoAdmin/template"                             // 模板包，提供模板渲染功能
	"github.com/purpose168/GoAdmin/template/chartjs"                          // Chart.js 模板包，提供图表组件功能
	"github.com/purpose168/GoAdmin/template/icon"                             // 图标包，提供图标定义功能
	"github.com/purpose168/GoAdmin/template/types"                            // 类型包，提供类型定义功能
)

// GetContent 获取仪表板页面的内容
//
// 参数:
//   - ctx: 上下文对象，包含请求的上下文信息
//
// 返回:
//   - types.Panel: 包含页面内容的面板对象
//   - error: 错误信息，如果生成过程中出现错误则返回
//
// 说明：
//
//	该函数构建并返回一个完整的仪表板页面，包含多个 UI 组件：
//	- 信息框（Info Box）：显示关键指标
//	- 表格（Table）：展示最新订单
//	- 产品列表（Product List）：显示最近添加的产品
//	- 图表（Chart）：销售趋势和进度信息
//	- 小方框（Small Box）：快速统计数据
//	- 饼图（Pie Chart）：浏览器使用情况
//	- 标签页（Tabs）：多标签内容展示
//	- 弹窗（Popup）：模态对话框
//
// 功能特性：
//   - 创建信息框组件，显示 CPU 流量、点赞数、销售量、新会员数
//   - 创建表格组件，展示最新订单信息
//   - 创建产品列表组件，展示最近添加的产品
//   - 创建折线图组件，展示销售趋势
//   - 创建进度条组组件，展示目标完成情况
//   - 创建描述组件，展示数据变化
//   - 创建小方框组件，快速统计数据
//   - 创建饼图组件，展示浏览器使用情况
//   - 创建标签页组件，展示多标签内容
//   - 创建弹窗组件，提供模态对话框
//   - 响应式布局，适配不同屏幕尺寸
//
// 页面结构：
//   - 第一行：四个信息框（CPU 流量、点赞数、销售量、新会员）
//   - 第二行：表格（最新订单）+ 产品列表（最近添加的产品）
//   - 第三行：月度回顾报告（折线图 + 进度条 + 描述组件）
//   - 第四行：四个小方框（新用户统计）
//   - 第五行：浏览器使用情况（饼图 + 图例）+ 标签页 + 弹窗
//
// 技术细节：
//   - 使用 tmpl.Default() 获取默认模板组件集合
//   - 使用 components.Col() 获取列组件，用于布局控制
//   - 使用 components.Row() 获取行组件，用于组合列
//   - 使用 components.Box() 获取盒子组件，用于包裹内容
//   - 使用 components.Table() 获取表格组件
//   - 使用 components.Tabs() 获取标签页组件
//   - 使用 components.Popup() 获取弹窗组件
//   - 使用 infobox.New() 创建信息框
//   - 使用 productlist.New() 创建产品列表
//   - 使用 chartjs.Line() 创建折线图
//   - 使用 chartjs.Pie() 创建饼图
//   - 使用 progress_group.New() 创建进度条组
//   - 使用 description.New() 创建描述组件
//   - 使用 smallbox.New() 创建小方框
//   - 使用 chart_legend.New() 创建图例
//   - 使用 types.Size() 设置列大小
//   - 使用 template.HTML() 将字符串转换为 HTML
//
// 使用场景：
//   - 管理后台: 作为 GoAdmin 管理后台的主页面
//   - 数据监控: 实时监控关键业务指标
//   - 数据分析: 通过图表分析数据趋势
//   - 业务概览: 快速了解业务运行状态
//   - 用户管理: 展示用户信息和统计数据
//
// 注意事项：
//   - 需要正确引入 AdminLTE 主题和 Chart.js 库
//   - 组件的 ID 必须唯一，避免冲突
//   - 响应式布局需要考虑不同屏幕尺寸的显示效果
//   - 图表数据需要正确格式化
//   - HTML 内容需要进行转义，防止 XSS 攻击
//   - 图标需要使用支持的图标库（Ion Icons 或 Font Awesome）
//
// 错误处理：
//   - 如果组件创建失败，会在运行时返回错误
//   - 如果模板渲染失败，会在运行时返回错误
//
// 示例：
//
//	// 创建上下文
//	ctx := context.NewContext(r)
//
//	// 获取仪表板内容
//	panel, err := GetContent(ctx)
//	if err != nil {
//	    // 处理错误
//	    log.Fatal(err)
//	}
//
//	// 使用面板内容
//	content := panel.Content
//	title := panel.Title
//	description := panel.Description
func GetContent(ctx *context.Context) (types.Panel, error) {

	components := tmpl.Default(ctx) // 获取默认的模板组件集合
	colComp := components.Col()     // 获取列组件，用于布局控制

	/**************************
	 * Info Box - 信息框部分
	 * 用于显示关键业务指标和统计数据
	/**************************/

	// 创建第一个信息框：显示 CPU 流量
	// SetText: 设置显示的文本
	// SetColor: 设置颜色主题（aqua-青色）
	// SetNumber: 设置显示的数值
	// SetIcon: 设置图标（使用 ion 图标库）
	infobox1 := infobox.New().
		SetText("CPU流量").
		SetColor("aqua").
		SetNumber("100").
		SetIcon("ion-ios-gear-outline").
		GetContent()

	// 创建第二个信息框：显示点赞数
	// 使用 Google Plus 图标，显示美元金额
	infobox2 := infobox.New().
		SetText("点赞数").
		SetColor("red").
		SetNumber("1030.00<small>$</small>").
		SetIcon(icon.GooglePlus).
		GetContent()

	// 创建第三个信息框：显示销售量
	// 使用购物车图标，绿色主题
	infobox3 := infobox.New().
		SetText("销售量").
		SetColor("green").
		SetNumber("760").
		SetIcon("ion-ios-cart-outline").
		GetContent()

	// 创建第四个信息框：显示新会员数
	// 使用人员图标，黄色主题
	// 注：支持 SVG 图标
	infobox4 := infobox.New().
		SetText("新会员").
		SetColor("yellow").
		SetNumber("2,349").
		SetIcon("ion-ios-people-outline").
		GetContent()

	// 设置列的大小
	// Size 参数：MD(中等屏幕)占 6 列，SM(小屏幕)占 3 列，XS(超小屏幕)占 12 列
	// 这种响应式设计确保在不同屏幕尺寸下都能良好显示
	var size = types.Size(6, 3, 0).XS(12)                                                                   // 设置列大小：中等屏幕 6 列，小屏幕 3 列，超小屏幕 12 列
	infoboxCol1 := colComp.SetSize(size).SetContent(infobox1).GetContent()                                  // 将信息框 1 放入列中
	infoboxCol2 := colComp.SetSize(size).SetContent(infobox2).GetContent()                                  // 将信息框 2 放入列中
	infoboxCol3 := colComp.SetSize(size).SetContent(infobox3).GetContent()                                  // 将信息框 3 放入列中
	infoboxCol4 := colComp.SetSize(size).SetContent(infobox4).GetContent()                                  // 将信息框 4 放入列中
	row1 := components.Row().SetContent(infoboxCol1 + infoboxCol2 + infoboxCol3 + infoboxCol4).GetContent() // 将四个信息框列组合成一行

	/**************************
	 * Box - 盒子和表格部分
	 * 用于展示最新订单信息
	/**************************/

	// 创建表格组件
	// SetInfoList: 设置表格数据列表
	// 每个 map 代表一行数据，key 为列名，value 为单元格内容
	table := components.Table().SetInfoList([]map[string]types.InfoItem{
		{
			"订单ID": {Content: "OR9842"},
			"商品":   {Content: "使命召唤IV"},
			"状态":   {Content: "已发货"},
			"热度":   {Content: "90%"},
		}, {
			"订单ID": {Content: "OR9842"},
			"商品":   {Content: "使命召唤IV"},
			"状态":   {Content: "已发货"},
			"热度":   {Content: "90%"},
		}, {
			"订单ID": {Content: "OR9842"},
			"商品":   {Content: "使命召唤IV"},
			"状态":   {Content: "已发货"},
			"热度":   {Content: "90%"},
		}, {
			"订单ID": {Content: "OR9842"},
			"商品":   {Content: "使命召唤IV"},
			"状态":   {Content: "已发货"},
			"热度":   {Content: "90%"},
		},
	}).SetThead(types.Thead{ // 设置表头
		{Head: "订单ID"},
		{Head: "商品"},
		{Head: "状态"},
		{Head: "热度"},
	}).GetContent()

	// 创建盒子组件包裹表格
	// WithHeadBorder: 添加头部边框
	// SetHeader: 设置标题
	// SetHeadColor: 设置头部背景色
	// SetBody: 设置主体内容
	// SetFooter: 设置底部内容（包含操作按钮）
	boxInfo := components.Box().
		WithHeadBorder().
		SetHeader("最新订单").
		SetHeadColor("#f7f7f7").
		SetBody(table).
		SetFooter(`<div class="clearfix"><a href="javascript:void(0)" class="btn btn-sm btn-info btn-flat pull-left">处理订单</a><a href="javascript:void(0)" class="btn btn-sm btn-default btn-flat pull-right">查看所有新订单</a> </div>`).
		GetContent()

	// 将表格盒子放入列中，中等屏幕占 8 列
	tableCol := colComp.SetSize(types.SizeMD(8)).SetContent(row1 + boxInfo).GetContent() // 将第一行（信息框）和表格盒子放入列中

	/**************************
	 * Product List - 产品列表部分
	 * 用于展示最近添加的产品
	/**************************/

	// 创建产品列表组件
	// SetData: 设置产品数据
	// 每个产品包含：标题、是否有表格、标签类型、标签、描述
	productList := productlist.New().SetData([]map[string]string{
		{
			"title":       "GoAdmin",
			"has_tabel":   "true",
			"labeltype":   "warning",
			"label":       "免费",
			"description": `帮助你构建数据可视化系统的框架`,
		}, {
			"title":       "GoAdmin",
			"has_tabel":   "true",
			"labeltype":   "warning",
			"label":       "免费",
			"description": `帮助你构建数据可视化系统的框架`,
		}, {
			"title":       "GoAdmin",
			"has_tabel":   "true",
			"labeltype":   "warning",
			"label":       "免费",
			"description": `帮助你构建数据可视化系统的框架`,
		}, {
			"title":       "GoAdmin",
			"has_tabel":   "true",
			"labeltype":   "warning",
			"label":       "免费",
			"description": `帮助你构建数据可视化系统的框架`,
		},
	}).GetContent()

	// 创建警告主题的盒子包裹产品列表
	// SetTheme: 设置主题颜色（warning-警告色）
	boxWarning := components.Box().SetTheme("warning").WithHeadBorder().SetHeader("最近添加的产品").
		SetBody(productList).
		SetFooter(`<a href="javascript:void(0)" class="uppercase">查看所有产品</a>`).
		GetContent()

	// 将产品列表盒子放入列中，中等屏幕占 4 列
	newsCol := colComp.SetSize(types.SizeMD(4)).SetContent(boxWarning).GetContent() // 将产品列表盒子放入列中

	// 将表格列和产品列组合成一行
	row5 := components.Row().SetContent(tableCol + newsCol).GetContent() // 将表格列和产品列组合成一行

	/**************************
	 * Box - 图表和进度条部分
	 * 用于展示月度报告和目标完成情况
	/**************************/

	// 创建折线图组件
	// SetID: 设置图表 ID，用于在页面中唯一标识
	// SetHeight: 设置图表高度
	// SetTitle: 设置图表标题
	// SetLabels: 设置 X 轴标签
	// AddDataSet: 添加数据集
	// DSData: 设置数据集的数值
	// DSFill: 设置是否填充区域
	// DSBorderColor: 设置边框颜色
	// DSLineTension: 设置线条张力（曲线平滑度）
	line := chartjs.Line() // 创建折线图

	lineChart := line.
		SetID("salechart").
		SetHeight(180).
		SetTitle("销售：2019年1月1日 - 2019年7月30日").
		SetLabels([]string{"一月", "二月", "三月", "四月", "五月", "六月", "七月"}).
		AddDataSet("电子产品").
		DSData([]float64{65, 59, 80, 81, 56, 55, 40}).
		DSFill(false).
		DSBorderColor("rgb(210, 214, 222)").
		DSLineTension(0.1).
		AddDataSet("数字商品").
		DSData([]float64{28, 48, 40, 19, 86, 27, 90}).
		DSFill(false).
		DSBorderColor("rgba(60,141,188,1)").
		DSLineTension(0.1).
		GetContent()

	// 创建进度组的标题
	title := `<p class="text-center"><strong>目标完成情况</strong></p>` // 进度组标题

	// 创建第一个进度组：添加商品到购物车
	// SetTitle: 设置进度条标题
	// SetColor: 设置进度条颜色
	// SetDenominator: 设置分母（总数）
	// SetMolecular: 设置分子（当前值）
	// SetPercent: 设置百分比
	progressGroup := progress_group.New().
		SetTitle("添加商品到购物车").
		SetColor("#76b2d4").
		SetDenominator(200).
		SetMolecular(160).
		SetPercent(80).
		GetContent()

	// 创建第二个进度组：完成购买
	progressGroup1 := progress_group.New().
		SetTitle("完成购买").
		SetColor("#f17c6e").
		SetDenominator(400).
		SetMolecular(310).
		SetPercent(80).
		GetContent()

	// 创建第三个进度组：访问高级页面
	progressGroup2 := progress_group.New().
		SetTitle("访问高级页面").
		SetColor("#ace0ae").
		SetDenominator(800).
		SetMolecular(490).
		SetPercent(80).
		GetContent()

	// 创建第四个进度组：发送咨询
	progressGroup3 := progress_group.New().
		SetTitle("发送咨询").
		SetColor("#fdd698").
		SetDenominator(500).
		SetMolecular(250).
		SetPercent(50).
		GetContent()

	// 将折线图放入列中，中等屏幕占 8 列
	boxInternalCol1 := colComp.SetContent(lineChart).SetSize(types.SizeMD(8)).GetContent() // 将折线图放入列中
	// 将进度组放入列中，中等屏幕占 4 列
	boxInternalCol2 := colComp.
		SetContent(template.HTML(title) + progressGroup + progressGroup1 + progressGroup2 + progressGroup3). // 将标题和四个进度组组合
		SetSize(types.SizeMD(4)).
		GetContent()

	// 将折线图列和进度组列组合成一行
	boxInternalRow := components.Row().SetContent(boxInternalCol1 + boxInternalCol2).GetContent() // 将折线图列和进度组列组合成一行

	// 创建描述组件 1：总收入
	// SetPercent: 设置百分比变化
	// SetNumber: 设置数值
	// SetTitle: 设置标题
	// SetArrow: 设置箭头方向（up-上升，down-下降）
	// SetColor: 设置颜色（green-绿色表示增长，red-红色表示下降）
	// SetBorder: 设置边框位置
	description1 := description.New().
		SetPercent("17").
		SetNumber("¥140,100").
		SetTitle("总收入").
		SetArrow("up").
		SetColor("green").
		SetBorder("right").
		GetContent()

	// 创建描述组件 2
	description2 := description.New().
		SetPercent("2").
		SetNumber("440,560").
		SetTitle("总收入").
		SetArrow("down").
		SetColor("red").
		SetBorder("right").
		GetContent()

	// 创建描述组件 3
	description3 := description.New().
		SetPercent("12").
		SetNumber("¥140,050").
		SetTitle("总收入").
		SetArrow("up").
		SetColor("green").
		SetBorder("right").
		GetContent()

	// 创建描述组件 4
	description4 := description.New().
		SetPercent("1").
		SetNumber("30943").
		SetTitle("总收入").
		SetArrow("up").
		SetColor("green").
		GetContent()

	// 设置描述组件的列大小：超小屏幕占 6 列，小屏幕占 3 列
	size2 := types.SizeXS(6).SM(3)                                                  // 设置列大小：超小屏幕 6 列，小屏幕 3 列
	boxInternalCol3 := colComp.SetContent(description1).SetSize(size2).GetContent() // 将描述组件 1 放入列中
	boxInternalCol4 := colComp.SetContent(description2).SetSize(size2).GetContent() // 将描述组件 2 放入列中
	boxInternalCol5 := colComp.SetContent(description3).SetSize(size2).GetContent() // 将描述组件 3 放入列中
	boxInternalCol6 := colComp.SetContent(description4).SetSize(size2).GetContent() // 将描述组件 4 放入列中

	// 将四个描述组件列组合成一行
	boxInternalRow2 := components.Row().SetContent(boxInternalCol3 + boxInternalCol4 + boxInternalCol5 + boxInternalCol6).GetContent() // 将四个描述组件列组合成一行

	// 创建盒子组件包裹图表和进度条
	box := components.Box().WithHeadBorder().SetHeader("月度回顾报告").
		SetBody(boxInternalRow).
		SetFooter(boxInternalRow2).
		GetContent()

	// 将盒子放入列中，中等屏幕占 12 列（全宽）
	boxcol := colComp.SetContent(box).SetSize(types.SizeMD(12)).GetContent() // 将盒子放入列中
	// 将盒子列组合成一行
	row2 := components.Row().SetContent(boxcol).GetContent() // 将盒子列组合成一行

	/**************************
	 * Small Box - 小方框部分
	 * 用于快速展示关键统计数据
	/**************************/

	// 创建第一个小方框：蓝色主题，显示新用户数
	smallbox1 := smallbox.New().SetColor("blue").SetIcon("ion-ios-gear-outline").SetUrl("/").SetTitle("新用户").SetValue("345￥").GetContent()
	// 创建第二个小方框：黄色主题，显示百分比
	smallbox2 := smallbox.New().SetColor("yellow").SetIcon("ion-ios-cart-outline").SetUrl("/").SetTitle("新用户").SetValue("80%").GetContent()
	// 创建第三个小方框：红色主题，显示金额
	smallbox3 := smallbox.New().SetColor("red").SetIcon("fa-user").SetUrl("/").SetTitle("新用户").SetValue("645￥").GetContent()
	// 创建第四个小方框：绿色主题，显示金额
	smallbox4 := smallbox.New().SetColor("green").SetIcon("ion-ios-cart-outline").SetUrl("/").SetTitle("新用户").SetValue("889￥").GetContent()

	// 将四个小方框放入列中
	col1 := colComp.SetSize(size).SetContent(smallbox1).GetContent() // 将小方框 1 放入列中
	col2 := colComp.SetSize(size).SetContent(smallbox2).GetContent() // 将小方框 2 放入列中
	col3 := colComp.SetSize(size).SetContent(smallbox3).GetContent() // 将小方框 3 放入列中
	col4 := colComp.SetSize(size).SetContent(smallbox4).GetContent() // 将小方框 4 放入列中

	// 将四个小方框列组合成一行
	row3 := components.Row().SetContent(col1 + col2 + col3 + col4).GetContent() // 将四个小方框列组合成一行

	/**************************
	 * Pie Chart - 饼图部分
	 * 用于展示浏览器使用情况
	/**************************/

	// 创建饼图组件
	// SetHeight: 设置图表高度
	// SetLabels: 设置标签（浏览器名称）
	// SetID: 设置图表 ID
	// AddDataSet: 添加数据集
	// DSData: 设置数据值
	// DSBackgroundColor: 设置背景颜色
	pie := chartjs.Pie().
		SetHeight(170).
		SetLabels([]string{"导航器", "Opera", "Safari", "火狐", "IE", "Chrome"}).
		SetID("pieChart").
		AddDataSet("Chrome").
		DSData([]float64{100, 300, 600, 400, 500, 700}).
		DSBackgroundColor([]chartjs.Color{
			"rgb(255, 205, 86)", "rgb(54, 162, 235)", "rgb(255, 99, 132)", "rgb(255, 205, 86)", "rgb(54, 162, 235)", "rgb(255, 99, 132)",
		}).
		GetContent()

	// 创建图例组件
	// SetData: 设置图例数据，包含标签和颜色
	legend := chart_legend.New().SetData([]map[string]string{
		{
			"label": " Chrome",
			"color": "red",
		}, {
			"label": " IE",
			"color": "Green",
		}, {
			"label": " 火狐",
			"color": "yellow",
		}, {
			"label": " Safari",
			"color": "blue",
		}, {
			"label": " Opera",
			"color": "light-blue",
		}, {
			"label": " 导航器",
			"color": "gray",
		},
	}).GetContent()

	// 创建危险主题的盒子包裹饼图和图例
	boxDanger := components.Box().SetTheme("danger").WithHeadBorder().SetHeader("浏览器使用情况").
		SetBody(components.Row().
			SetContent(colComp.SetSize(types.SizeMD(8)).
				SetContent(pie).
				GetContent() + colComp.SetSize(types.SizeMD(4)).
				SetContent(legend).
				GetContent()).GetContent()).
		SetFooter(`<p class="text-center"><a href="javascript:void(0)" class="uppercase">查看所有用户</a></p>`).
		GetContent()

	/**************************
	 * Tabs and Popup - 标签页和弹窗部分
	 * 用于展示多标签内容和模态对话框
	/**************************/

	// 创建标签页组件
	// SetData: 设置标签页数据
	// 每个标签页包含：标题和内容
	tabs := components.Tabs().SetData([]map[string]template.HTML{
		{
			"title": "标签1",
			"content": template.HTML(`<b>如何使用：</b>

                <p>与原始Bootstrap标签完全相同，只是应该使用自定义包装器.nav-tabs-custom来实现此样式。</p>
                一种美妙的宁静占据了我的整个灵魂，
                就像这些我全心全意享受的春天的甜美早晨。
                我独自一人，在这个地方感受到存在的魅力，
                这个地方是为像我这样的灵魂的幸福而创造的。我是如此快乐，
                我亲爱的朋友，如此沉浸在单纯的宁静存在的精致感觉中，
                以至于我忽视了我的天赋。在目前这一刻，我无法画出一笔；
                然而我觉得我从未比现在更是一个伟大的艺术家。`),
		}, {
			"title": "标签2",
			"content": template.HTML(`
                欧洲语言属于同一个家族。它们的独立存在是一个神话。
                对于科学、音乐、体育等，欧洲使用相同的词汇。这些语言仅在
                它们的语法、发音和最常见的单词上有所不同。每个人都意识到为什么
                一种新的通用语言是可取的：人们可以拒绝支付昂贵的翻译费用。为了
                实现这一点，有必要有统一的语法、发音和更常见的
                单词。如果几种语言融合，结果语言的语法将更简单
                和规律，比个别语言的语法更简单。
              `),
		}, {
			"title": "标签3",
			"content": template.HTML(`
                Lorem Ipsum只是印刷和排版行业的虚拟文本。
                自1500年代以来，Lorem Ipsum一直是该行业的标准虚拟文本，
                当时一位不知名的印刷商拿了一个字型盘并将其打乱以制作字体样本书。
                它不仅存活了五个世纪，而且跨越到电子排版，
                基本上保持不变。它在1960年代随着包含Lorem Ipsum段落的Letraset
                表的发布而流行起来，最近还有桌面发布软件
                如Aldus PageMaker包括Lorem Ipsum的版本。
              `),
		},
	}).GetContent()

	// 创建测试按钮，用于触发弹窗
	buttonTest := `<button type="button" class="btn btn-primary" data-toggle="modal" data-target="#exampleModal" data-whatever="@mdo">为@mdo打开模态框</button>` // 测试按钮

	// 创建弹窗表单内容
	popupForm := `<form>
          <div class="form-group">
            <label for="recipient-name" class="col-form-label">收件人：</label>
            <input type="text" class="form-control" id="recipient-name">
          </div>
          <div class="form-group">
            <label for="message-text" class="col-form-label">消息：</label>
            <textarea class="form-control" id="message-text"></textarea>
          </div>
        </form>` // 弹窗表单

	// 创建弹窗组件
	// SetID: 设置弹窗 ID
	// SetFooter: 设置底部按钮
	// SetTitle: 设置标题
	// SetBody: 设置主体内容
	popup := components.Popup().SetID("exampleModal").
		SetFooter("保存更改").
		SetTitle("这是一个弹窗").
		SetBody(template.HTML(popupForm)).
		GetContent()

	// 将标签页和按钮放入列中，中等屏幕占 8 列
	col5 := colComp.SetSize(types.SizeMD(8)).SetContent(tabs + template.HTML(buttonTest)).GetContent() // 将标签页和按钮放入列中
	// 将饼图盒子和弹窗放入列中，中等屏幕占 4 列
	col6 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger + popup).GetContent() // 将饼图盒子和弹窗放入列中

	// 将两列组合成一行
	row4 := components.Row().SetContent(col5 + col6).GetContent() // 将两列组合成一行

	// 返回面板对象
	// Content: 页面内容（所有行的组合）
	// Title: 页面标题
	// Description: 页面描述
	return types.Panel{
		Content:     row3 + row2 + row5 + row4, // 页面内容：小方框行 + 月度回顾报告行 + 表格和产品列表行 + 标签页和弹窗行
		Title:       "仪表板",                     // 页面标题
		Description: "仪表板示例",                   // 页面描述
	}, nil // 返回 nil 表示没有错误
}

package datamodel

import (
	"html/template"

	"github.com/purpose168/GoAdmin-themes/adminlte/components/chart_legend"
	"github.com/purpose168/GoAdmin-themes/adminlte/components/description"
	"github.com/purpose168/GoAdmin-themes/adminlte/components/infobox"
	"github.com/purpose168/GoAdmin-themes/adminlte/components/productlist"
	"github.com/purpose168/GoAdmin-themes/adminlte/components/progress_group"
	"github.com/purpose168/GoAdmin-themes/adminlte/components/smallbox"
	"github.com/purpose168/GoAdmin/context"
	tmpl "github.com/purpose168/GoAdmin/template"
	"github.com/purpose168/GoAdmin/template/chartjs"
	"github.com/purpose168/GoAdmin/template/icon"
	"github.com/purpose168/GoAdmin/template/types"
)

// GetContent 返回首页的内容面板
// 该函数创建并配置一个包含多种UI组件的仪表板页面，用于展示数据可视化和管理信息
// 参数 ctx：上下文对象，用于传递请求上下文信息
// 返回值：配置完成的面板对象和错误信息
func GetContent(ctx *context.Context) (types.Panel, error) {

	// 获取默认模板组件集合
	// components 包含了各种可用的UI组件，如表格、盒子、按钮等
	components := tmpl.Default(ctx)

	// 获取列组件
	// colComp 用于创建和管理页面布局中的列
	colComp := components.Col()

	/**************************
	 * 信息框 (Info Box)
	/**************************/

	// 创建第一个信息框：显示CPU流量
	// SetText：设置显示的文本标签
	// SetColor：设置主题颜色（aqua=青色）
	// SetNumber：显示的数值
	// SetIcon：设置图标（使用Ionicons图标库）
	infobox1 := infobox.New().
		SetText("CPU流量").
		SetColor("aqua").
		SetNumber("100").
		SetIcon("ion-ios-gear-outline").
		GetContent()

	// 创建第二个信息框：显示点赞数
	// SetNumber支持HTML标签，可以添加货币符号等
	// icon.GooglePlus：使用预定义的Google Plus图标
	infobox2 := infobox.New().
		SetText("点赞").
		SetColor("red").
		SetNumber("1030.00<small>￥</small>").
		SetIcon(icon.GooglePlus).
		GetContent()

	// 创建第三个信息框：显示销售数量
	// SetColor("green")：设置绿色主题
	infobox3 := infobox.New().
		SetText("销售").
		SetColor("green").
		SetNumber("760").
		SetIcon("ion-ios-cart-outline").
		GetContent()

	// 创建第四个信息框：显示新成员数量
	// SetColor("yellow")：设置黄色主题
	// SetNumber支持带逗号的数字格式
	infobox4 := infobox.New().
		SetText("新成员").
		SetColor("yellow").
		SetNumber("2,349").
		SetIcon("ion-ios-people-outline"). // 支持SVG图标
		GetContent()

	// 设置列的大小
	// Size(6, 3, 0)：分别表示在中等屏幕(md)、小屏幕(sm)和超小屏幕(xs)上的列宽
	// XS(12)：在超小屏幕上占满整行（12列）
	var size = types.Size(6, 3, 0).XS(12)

	// 为每个信息框创建列并设置内容
	infoboxCol1 := colComp.SetSize(size).SetContent(infobox1).GetContent()
	infoboxCol2 := colComp.SetSize(size).SetContent(infobox2).GetContent()
	infoboxCol3 := colComp.SetSize(size).SetContent(infobox3).GetContent()
	infoboxCol4 := colComp.SetSize(size).SetContent(infobox4).GetContent()

	// 将四个信息框列组合成一行
	row1 := components.Row().SetContent(infoboxCol1 + infoboxCol2 + infoboxCol3 + infoboxCol4).GetContent()

	/**************************
	 * 盒子组件 (Box)
	/**************************/

	// 创建表格组件，显示最新订单信息
	// SetInfoList：设置表格的数据内容，每行是一个map，key是列名，value是InfoItem结构
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
	}).SetThead(types.Thead{
		{Head: "订单ID"},
		{Head: "商品"},
		{Head: "状态"},
		{Head: "热度"},
	}).GetContent()

	// 创建盒子组件，包含表格
	// WithHeadBorder()：添加头部边框
	// SetHeader：设置盒子标题
	// SetHeadColor：设置头部背景颜色
	// SetBody：设置盒子主体内容
	// SetFooter：设置底部内容，支持HTML
	boxInfo := components.Box().
		WithHeadBorder().
		SetHeader("最新订单").
		SetHeadColor("#f7f7f7").
		SetBody(table).
		SetFooter(`<div class="clearfix"><a href="javascript:void(0)" class="btn btn-sm btn-info btn-flat pull-left">处理订单</a><a href="javascript:void(0)" class="btn btn-sm btn-default btn-flat pull-right">查看所有新订单</a> </div>`).
		GetContent()

	// 创建包含信息框行和订单盒子的列
	tableCol := colComp.SetSize(types.SizeMD(8)).SetContent(row1 + boxInfo).GetContent()

	/**************************
	 * 产品列表 (Product List)
	/**************************/

	// 创建产品列表组件
	// SetData：设置产品数据，每个产品包含图片、标题、标签、描述等信息
	productList := productlist.New().SetData([]map[string]string{
		{
			"img":         "http://adminlte.io/themes/AdminLTE/dist/img/default-50x50.gif",
			"title":       "GoAdmin",
			"has_tabel":   "true",
			"labeltype":   "warning",
			"label":       "免费",
			"description": `一个帮助您构建数据可视化系统的框架`,
		}, {
			"img":         "http://adminlte.io/themes/AdminLTE/dist/img/default-50x50.gif",
			"title":       "GoAdmin",
			"has_tabel":   "true",
			"labeltype":   "warning",
			"label":       "免费",
			"description": `一个帮助您构建数据可视化系统的框架`,
		}, {
			"img":         "http://adminlte.io/themes/AdminLTE/dist/img/default-50x50.gif",
			"title":       "GoAdmin",
			"has_tabel":   "true",
			"labeltype":   "warning",
			"label":       "免费",
			"description": `一个帮助您构建数据可视化系统的框架`,
		}, {
			"img":         "http://adminlte.io/themes/AdminLTE/dist/img/default-50x50.gif",
			"title":       "GoAdmin",
			"has_tabel":   "true",
			"labeltype":   "warning",
			"label":       "免费",
			"description": `一个帮助您构建数据可视化系统的框架`,
		},
	}).GetContent()

	// 创建警告主题的盒子，显示最近添加的产品
	// SetTheme("warning")：设置主题为警告色（黄色）
	boxWarning := components.Box().SetTheme("warning").WithHeadBorder().SetHeader("最近添加的产品").
		SetBody(productList).
		SetFooter(`<a href="javascript:void(0)" class="uppercase">查看所有产品</a>`).
		GetContent()

	// 创建包含产品列表盒子的列
	newsCol := colComp.SetSize(types.SizeMD(4)).SetContent(boxWarning).GetContent()

	// 将表格列和产品列表列组合成一行
	row5 := components.Row().SetContent(tableCol + newsCol).GetContent()

	/**************************
	 * 盒子组件 (Box) - 包含图表和进度条
	/**************************/

	// 创建折线图组件
	line := chartjs.Line()

	// 配置折线图
	// SetID：设置图表的DOM元素ID
	// SetHeight：设置图表高度（像素）
	// SetTitle：设置图表标题
	// SetLabels：设置X轴标签
	// AddDataSet：添加数据集
	// DSData：设置数据集的数据
	// DSFill：设置是否填充区域
	// DSBorderColor：设置线条颜色
	// DSLineTension：设置线条张力（0为直线，1为平滑曲线）
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

	// 创建进度条组的标题
	title := `<p class="text-center"><strong>目标完成度</strong></p>`

	// 创建第一个进度条：添加商品到购物车
	// SetTitle：设置进度条标题
	// SetColor：设置进度条颜色
	// SetDenominator：设置分母（总数）
	// SetMolecular：设置分子（当前值）
	// SetPercent：设置百分比
	progressGroup := progress_group.New().
		SetTitle("添加商品到购物车").
		SetColor("#76b2d4").
		SetDenominator(200).
		SetMolecular(160).
		SetPercent(80).
		GetContent()

	// 创建第二个进度条：完成购买
	progressGroup1 := progress_group.New().
		SetTitle("完成购买").
		SetColor("#f17c6e").
		SetDenominator(400).
		SetMolecular(310).
		SetPercent(80).
		GetContent()

	// 创建第三个进度条：访问高级页面
	progressGroup2 := progress_group.New().
		SetTitle("访问高级页面").
		SetColor("#ace0ae").
		SetDenominator(800).
		SetMolecular(490).
		SetPercent(80).
		GetContent()

	// 创建第四个进度条：发送咨询
	progressGroup3 := progress_group.New().
		SetTitle("发送咨询").
		SetColor("#fdd698").
		SetDenominator(500).
		SetMolecular(250).
		SetPercent(50).
		GetContent()

	// 创建包含折线图的列（占8列）
	boxInternalCol1 := colComp.SetContent(lineChart).SetSize(types.SizeMD(8)).GetContent()

	// 创建包含进度条组的列（占4列）
	boxInternalCol2 := colComp.
		SetContent(template.HTML(title) + progressGroup + progressGroup1 + progressGroup2 + progressGroup3).
		SetSize(types.SizeMD(4)).
		GetContent()

	// 将折线图列和进度条列组合成一行
	boxInternalRow := components.Row().SetContent(boxInternalCol1 + boxInternalCol2).GetContent()

	// 创建描述组件1：显示总收入
	// SetPercent：设置百分比变化
	// SetNumber：显示的数值
	// SetTitle：设置标题
	// SetArrow：设置箭头方向（up=上升，down=下降）
	// SetColor：设置颜色主题
	// SetBorder：设置边框位置
	description1 := description.New().
		SetPercent("17").
		SetNumber("¥140,100").
		SetTitle("总收入").
		SetArrow("up").
		SetColor("green").
		SetBorder("right").
		GetContent()

	// 创建描述组件2
	description2 := description.New().
		SetPercent("2").
		SetNumber("440,560").
		SetTitle("总收入").
		SetArrow("down").
		SetColor("red").
		SetBorder("right").
		GetContent()

	// 创建描述组件3
	description3 := description.New().
		SetPercent("12").
		SetNumber("¥140,050").
		SetTitle("总收入").
		SetArrow("up").
		SetColor("green").
		SetBorder("right").
		GetContent()

	// 创建描述组件4（无右边框）
	description4 := description.New().
		SetPercent("1").
		SetNumber("30943").
		SetTitle("总收入").
		SetArrow("up").
		SetColor("green").
		GetContent()

	// 设置描述组件列的大小
	// SizeXS(6).SM(3)：在超小屏幕上占6列，在小屏幕上占3列
	size2 := types.SizeXS(6).SM(3)

	// 为每个描述组件创建列
	boxInternalCol3 := colComp.SetContent(description1).SetSize(size2).GetContent()
	boxInternalCol4 := colComp.SetContent(description2).SetSize(size2).GetContent()
	boxInternalCol5 := colComp.SetContent(description3).SetSize(size2).GetContent()
	boxInternalCol6 := colComp.SetContent(description4).SetSize(size2).GetContent()

	// 将四个描述组件列组合成一行
	boxInternalRow2 := components.Row().SetContent(boxInternalCol3 + boxInternalCol4 + boxInternalCol5 + boxInternalCol6).GetContent()

	// 创建月度回顾报告盒子
	// SetBody：设置主体内容（图表和进度条）
	// SetFooter：设置底部内容（描述组件）
	box := components.Box().WithHeadBorder().SetHeader("月度回顾报告").
		SetBody(boxInternalRow).
		SetFooter(boxInternalRow2).
		GetContent()

	// 创建包含月度报告盒子的列
	boxcol := colComp.SetContent(box).SetSize(types.SizeMD(12)).GetContent()

	// 将月度报告列组合成一行
	row2 := components.Row().SetContent(boxcol).GetContent()

	/**************************
	 * 小盒子 (Small Box)
	/**************************/

	// 创建四个小盒子组件
	// SetColor：设置主题颜色
	// SetIcon：设置图标
	// SetUrl：设置点击跳转的URL
	// SetTitle：设置标题
	// SetValue：设置显示的数值
	smallbox1 := smallbox.New().SetColor("blue").SetIcon("ion-ios-gear-outline").SetUrl("/").SetTitle("新用户").SetValue("345￥").GetContent()
	smallbox2 := smallbox.New().SetColor("yellow").SetIcon("ion-ios-cart-outline").SetUrl("/").SetTitle("新用户").SetValue("80%").GetContent()
	smallbox3 := smallbox.New().SetColor("red").SetIcon("fa-user").SetUrl("/").SetTitle("新用户").SetValue("645￥").GetContent()
	smallbox4 := smallbox.New().SetColor("green").SetIcon("ion-ios-cart-outline").SetUrl("/").SetTitle("新用户").SetValue("889￥").GetContent()

	// 为每个小盒子创建列
	col1 := colComp.SetSize(size).SetContent(smallbox1).GetContent()
	col2 := colComp.SetSize(size).SetContent(smallbox2).GetContent()
	col3 := colComp.SetSize(size).SetContent(smallbox3).GetContent()
	col4 := colComp.SetSize(size).SetContent(smallbox4).GetContent()

	// 将四个小盒子列组合成一行
	row3 := components.Row().SetContent(col1 + col2 + col3 + col4).GetContent()

	/**************************
	 * 饼图 (Pie Chart)
	/**************************/

	// 创建饼图组件
	// SetHeight：设置饼图高度
	// SetLabels：设置各扇区的标签
	// SetID：设置图表ID
	// AddDataSet：添加数据集
	// DSData：设置各扇区的数值
	// DSBackgroundColor：设置各扇区的背景色
	pie := chartjs.Pie().
		SetHeight(170).
		SetLabels([]string{"Navigator", "Opera", "Safari", "FireFox", "IE", "Chrome"}).
		SetID("pieChart").
		AddDataSet("Chrome").
		DSData([]float64{100, 300, 600, 400, 500, 700}).
		DSBackgroundColor([]chartjs.Color{
			"rgb(255, 205, 86)", "rgb(54, 162, 235)", "rgb(255, 99, 132)", "rgb(255, 205, 86)", "rgb(54, 162, 235)", "rgb(255, 99, 132)",
		}).
		GetContent()

	// 创建图表图例组件
	// SetData：设置图例数据，每个图例包含标签和颜色
	legend := chart_legend.New().SetData([]map[string]string{
		{
			"label": " Chrome",
			"color": "red",
		}, {
			"label": " IE",
			"color": "Green",
		}, {
			"label": " FireFox",
			"color": "yellow",
		}, {
			"label": " Safari",
			"color": "blue",
		}, {
			"label": " Opera",
			"color": "light-blue",
		}, {
			"label": " Navigator",
			"color": "gray",
		},
	}).GetContent()

	// 创建浏览器使用情况盒子（危险主题）
	// SetTheme("danger")：设置主题为危险色（红色）
	boxDanger := components.Box().SetTheme("danger").WithHeadBorder().SetHeader("浏览器使用情况").
		SetBody(components.Row().
			SetContent(colComp.SetSize(types.SizeMD(8)).
				SetContent(pie).
				GetContent() + colComp.SetSize(types.SizeMD(4)).
				SetContent(legend).
				GetContent()).GetContent()).
		SetFooter(`<p class="text-center"><a href="javascript:void(0)" class="uppercase">查看所有用户</a></p>`).
		GetContent()

	// 创建标签页组件
	// SetData：设置标签页数据，每个标签页包含标题和内容
	tabs := components.Tabs().SetData([]map[string]template.HTML{
		{
			"title": "标签页1",
			"content": template.HTML(`<b>使用方法：</b>

                <p>与原始的bootstrap标签页完全相同，只是您应该使用
                  自定义包装器 <code>.nav-tabs-custom</code> 来实现这种样式。</p>
                一种美妙的宁静占据了我的整个灵魂，
                就像这些我全心全意享受的春天的甜美早晨。
                我独自一人，在这个地方感受到了存在的魅力，
                这个地方是为像我这样的灵魂的幸福而创造的。我是如此快乐，
                我亲爱的朋友，如此沉浸在单纯的宁静存在的精致感觉中，
                以至于我忽视了才华。在目前这一刻，我无法画出一笔；
                然而我觉得我从未像现在这样成为一个伟大的艺术家。`),
		}, {
			"title": "标签页2",
			"content": template.HTML(`
                欧洲语言是同一个家族的成员。它们各自的存在是一个神话。
                对于科学、音乐、体育等，欧洲使用相同的词汇。这些语言仅在
                语法、发音和最常用的单词上有所不同。每个人都意识到为什么
                一种新的通用语言是可取的：人们可以拒绝支付昂贵的翻译费用。为了
                实现这一点，有必要有统一的语法、发音和更多常见的
                单词。如果几种语言合并，结果语言的语法比
                个别语言的语法更简单和规则。
              `),
		}, {
			"title": "标签页3",
			"content": template.HTML(`
                Lorem Ipsum只是印刷和排版行业的虚拟文本。
                自1500年代以来，Lorem Ipsum一直是该行业的标准虚拟文本，
                当时一位不知名的印刷商拿了一个字盘并将其打乱以制作一个字体样本书。
                它不仅存活了五个世纪，而且跨越到电子排版，
                基本上保持不变。它在1960年代随着包含Lorem Ipsum段落的Letraset
                表的发布而流行起来，最近随着桌面出版软件
                如Aldus PageMaker（包括Lorem Ipsum的版本）而流行起来。
              `),
		},
	}).GetContent()

	// 创建测试按钮，用于打开弹窗
	// data-toggle="modal"：Bootstrap模态框触发器
	// data-target：指定要打开的模态框ID
	buttonTest := `<button type="button" class="btn btn-primary" data-toggle="modal" data-target="#exampleModal" data-whatever="@mdo">为@mdo打开模态框</button>`

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
        </form>`

	// 创建弹窗组件
	// SetID：设置弹窗ID（与按钮的data-target对应）
	// SetFooter：设置底部按钮文字
	// SetTitle：设置弹窗标题
	// SetBody：设置弹窗主体内容
	popup := components.Popup().SetID("exampleModal").
		SetFooter("保存更改").
		SetTitle("这是一个弹窗").
		SetBody(template.HTML(popupForm)).
		GetContent()

	// 创建包含标签页和按钮的列（占8列）
	col5 := colComp.SetSize(types.SizeMD(8)).SetContent(tabs + template.HTML(buttonTest)).GetContent()

	// 创建包含浏览器使用情况盒子和弹窗的列（占4列）
	col6 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger + popup).GetContent()

	// 将两列组合成一行
	row4 := components.Row().SetContent(col5 + col6).GetContent()

	// 返回最终的面板配置
	// Content：页面的主要内容（所有行的组合）
	// Title：页面标题
	// Description：页面描述
	return types.Panel{
		Content:     row3 + row2 + row5 + row4,
		Title:       "仪表板",
		Description: "仪表板示例",
	}, nil
}

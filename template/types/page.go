// Copyright 2019 GoAdmin Core Team. All rights reserved.
// Use of this source code is governed by a Apache-2.0 style
// license that can be found in the LICENSE file.

package types

import (
	"bytes"
	"fmt"
	"html/template"
	"strconv"
	textTmpl "text/template"

	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/modules/config"
	"github.com/purpose168/GoAdmin/modules/menu"
	"github.com/purpose168/GoAdmin/modules/system"
	"github.com/purpose168/GoAdmin/modules/utils"
	"github.com/purpose168/GoAdmin/plugins/admin/models"
)

// Attribute 是模板的组件接口。模板的每个组件都应该实现它
type Attribute struct {
	TemplateList map[string]string // 模板列表
	Separation   bool              // 是否分隔
}

// Page 在模板中用作顶级变量
type Page struct {
	// User 是登录用户
	User models.UserModel

	// Menu 是模板的左侧菜单
	Menu menu.Menu

	// Panel 是模板的主要内容
	Panel Panel

	// System 包含一些系统信息
	System SystemInfo

	// UrlPrefix 是URL的前缀
	UrlPrefix string

	// Title 是网页的标题
	Title string

	// Logo 是模板的logo
	Logo template.HTML

	// MiniLogo 是模板的缩小版logo
	MiniLogo template.HTML

	// ColorScheme 是模板的颜色方案
	ColorScheme string

	// IndexUrl 是站点的主页URL
	IndexUrl string

	// AssetUrl 是资源的CDN链接
	CdnUrl string

	// Custom html in the tag head -> 自定义head标签内的HTML
	CustomHeadHtml template.HTML

	// Custom html after body -> 自定义body后的HTML
	CustomFootHtml template.HTML

	TmplHeadHTML template.HTML // 模板头部HTML
	TmplFootJS   template.HTML // 模板底部JS

	// Components assets -> 组件资源
	AssetsList template.HTML

	// Footer info -> 页脚信息
	FooterInfo template.HTML

	// Load as Iframe or not -> 是否以iframe加载
	Iframe bool

	// Whether update menu or not -> 是否更新菜单
	UpdateMenu bool

	// Top Nav Buttons -> 顶部导航按钮
	navButtons     Buttons
	NavButtonsHTML template.HTML
}

// NewPageParam 是创建新页面的参数结构体
type NewPageParam struct {
	User           models.UserModel // 用户模型
	Menu           *menu.Menu       // 菜单
	UpdateMenu     bool             // 是否更新菜单
	Panel          Panel            // 面板
	Logo           template.HTML    // Logo
	Assets         template.HTML    // 资源
	Buttons        Buttons          // 按钮
	Iframe         bool             // 是否以iframe加载
	TmplHeadHTML   template.HTML    // 模板头部HTML
	TmplFootJS     template.HTML    // 模板底部JS
	NavButtonsHTML template.HTML    // 导航按钮HTML
	NavButtonsJS   template.HTML    // 导航按钮JS
}

// NavButtonsAndJS 获取导航按钮和JS
// 参数:
//   - ctx: 上下文对象
//
// 返回: 导航按钮HTML和底部HTML
func (param *NewPageParam) NavButtonsAndJS(ctx *context.Context) (template.HTML, template.HTML) {
	navBtnFooter := template.HTML("")
	navBtn := template.HTML("")
	btnJS := template.JS("")

	for _, btn := range param.Buttons {
		if btn.IsType(ButtonTypeNavDropDown) {
			content, js := btn.Content(ctx)
			navBtn += content
			btnJS += js
			for _, item := range btn.(*NavDropDownButton).Items {
				navBtnFooter += item.GetAction().FooterContent(ctx)
				_, js := item.Content(ctx)
				btnJS += js
			}
		} else {
			navBtnFooter += btn.GetAction().FooterContent(ctx)
			content, js := btn.Content(ctx)
			navBtn += content
			btnJS += js
		}
	}

	return template.HTML(ParseTableDataTmpl(navBtn)),
		navBtnFooter + template.HTML(ParseTableDataTmpl(`<script>`+btnJS+`</script>`))
}

// NewPage 创建新页面
// 参数:
//   - ctx: 上下文对象
//   - param: 页面参数
//
// 返回: 页面对象
func NewPage(ctx *context.Context, param *NewPageParam) *Page {

	if param.NavButtonsHTML == template.HTML("") {
		param.NavButtonsHTML, param.NavButtonsJS = param.NavButtonsAndJS(ctx)
	}

	logo := param.Logo
	if logo == template.HTML("") {
		logo = config.GetLogo()
	}

	return &Page{
		User:       param.User,
		Menu:       *param.Menu,
		Panel:      param.Panel,
		UpdateMenu: param.UpdateMenu,
		System: SystemInfo{
			Version: system.Version(),
			Theme:   config.GetTheme(),
		},
		UrlPrefix:      config.AssertPrefix(),
		Title:          config.GetTitle(),
		Logo:           logo,
		MiniLogo:       config.GetMiniLogo(),
		ColorScheme:    config.GetColorScheme(),
		IndexUrl:       config.GetIndexURL(),
		CdnUrl:         config.GetAssetUrl(),
		CustomHeadHtml: config.GetCustomHeadHtml(),
		CustomFootHtml: config.GetCustomFootHtml() + param.NavButtonsJS,
		FooterInfo:     config.GetFooterInfo(),
		AssetsList:     param.Assets,
		navButtons:     param.Buttons,
		Iframe:         param.Iframe,
		NavButtonsHTML: param.NavButtonsHTML,
		TmplHeadHTML:   param.TmplHeadHTML,
		TmplFootJS:     param.TmplFootJS,
	}
}

// AddButton 添加按钮
// 参数:
//   - ctx: 上下文对象
//   - title: 按钮标题
//   - icon: 图标
//   - action: 操作对象
//
// 返回: 更新后的页面对象
func (page *Page) AddButton(ctx *context.Context, title template.HTML, icon string, action Action) *Page {
	page.navButtons = append(page.navButtons, GetNavButton(title, icon, action))
	page.CustomFootHtml += action.FooterContent(ctx)
	return page
}

// NewPagePanel 创建新页面面板
// 参数:
//   - panel: 面板对象
//
// 返回: 页面对象
func NewPagePanel(panel Panel) *Page {
	return &Page{
		Panel: panel,
		System: SystemInfo{
			Version: system.Version(),
		},
	}
}

// SystemInfo 包含系统的基本信息
type SystemInfo struct {
	Version string // 版本号
	Theme   string // 主题
}

// TableRowData 是表格行数据结构体
type TableRowData struct {
	Id    template.HTML       // 行ID
	Ids   template.HTML       // 行ID列表
	Value map[string]InfoItem // 值映射
}

// ParseTableDataTmpl 解析表格数据模板
// 参数:
//   - content: 模板内容
//
// 返回: 解析后的字符串
func ParseTableDataTmpl(content interface{}) string {
	var (
		c  string
		ok bool
	)
	if c, ok = content.(string); !ok {
		if cc, ok := content.(template.HTML); ok {
			c = string(cc)
		} else {
			c = string(content.(template.JS))
		}
	}
	t := template.New("row_data_tmpl")
	t, _ = t.Parse(c)
	buf := new(bytes.Buffer)
	_ = t.Execute(buf, TableRowData{Ids: `typeof(selectedRows)==="function" ? selectedRows().join() : ""`})
	return buf.String()
}

// ParseTableDataTmplWithID 解析带ID的表格数据模板
// 参数:
//   - id: 行ID
//   - content: 模板内容
//   - value: 可选的值映射
//
// 返回: 解析后的字符串
func ParseTableDataTmplWithID(id template.HTML, content string, value ...map[string]InfoItem) string {
	t := textTmpl.New("row_data_tmpl")
	t, _ = t.Parse(content)
	buf := new(bytes.Buffer)
	v := make(map[string]InfoItem)
	if len(value) > 0 {
		v = value[0]
	}
	_ = t.Execute(buf, TableRowData{
		Id:    id,
		Ids:   `typeof(selectedRows)==="function" ? selectedRows().join() : ""`,
		Value: v,
	})
	return buf.String()
}

// Panel 包含模板的主要内容，用作pjax
type Panel struct {
	Title       template.HTML // 标题
	Description template.HTML // 描述
	Content     template.HTML // 内容

	CSS template.CSS // 样式
	JS  template.JS  // JavaScript
	Url string       // URL

	// Whether to toggle the sidebar -> 是否切换侧边栏
	MiniSidebar bool

	// Auto refresh page switch -> 自动刷新页面开关
	AutoRefresh bool
	// Refresh page intervals, the unit is second -> 刷新页面间隔，单位为秒
	RefreshInterval []int

	Callbacks Callbacks // 回调列表
}

// Component 是组件接口
type Component interface {
	GetContent() template.HTML // 获取内容
	GetJS() template.JS        // 获取JavaScript
	GetCSS() template.CSS      // 获取样式
	GetCallbacks() Callbacks   // 获取回调
}

// AddComponent 添加组件
// 参数:
//   - comp: 组件对象
//
// 返回: 更新后的面板
func (p Panel) AddComponent(comp Component) Panel {
	p.JS += comp.GetJS()
	p.CSS += comp.GetCSS()
	p.Content += comp.GetContent()
	p.Callbacks = append(p.Callbacks, comp.GetCallbacks()...)
	return p
}

// AddJS 添加JavaScript
// 参数:
//   - js: JavaScript代码
//
// 返回: 更新后的面板
func (p Panel) AddJS(js template.JS) Panel {
	p.JS += js
	return p
}

// GetContent 获取面板内容
// 参数:
//   - params: 可选参数，第一个参数表示是否压缩，第二个参数表示是否启用动画
//
// 返回: 更新后的面板
func (p Panel) GetContent(params ...bool) Panel {

	compress := false

	if len(params) > 0 {
		compress = params[0]
	}

	var (
		animation, style, remove = template.HTML(""), template.HTML(""), template.HTML("")
		ani                      = config.GetAnimation()
	)

	if ani.Type != "" && (len(params) < 2 || params[1]) {
		animation = template.HTML(` class='pjax-container-content animated ` + ani.Type + `'`)
		if ani.Delay != 0 {
			style = template.HTML(fmt.Sprintf(`animation-delay: %fs;-webkit-animation-delay: %fs;`, ani.Delay, ani.Delay))
		}
		if ani.Duration != 0 {
			style = template.HTML(fmt.Sprintf(`animation-duration: %fs;-webkit-animation-duration: %fs;`, ani.Duration, ani.Duration))
		}
		if style != "" {
			style = ` style="` + style + `"`
		}
		remove = template.HTML(`<script>
		$('.pjax-container-content .modal.fade').on('show.bs.modal', function (event) {
            // 修复Animate.css
			$('.pjax-container-content').removeClass('` + ani.Type + `');
        });
		</script>`)
	}

	p.Content = `<div` + animation + style + ">" + p.Content + "</div>" + remove
	if p.MiniSidebar {
		p.Content += `<script>$("body").addClass("sidebar-collapse")</script>`
	}
	if p.AutoRefresh {
		refreshTime := 60
		if len(p.RefreshInterval) > 0 {
			refreshTime = p.RefreshInterval[0]
		}

		p.Content += `<script>
window.setTimeout(function(){
	$.pjax.reload('#pjax-container');	
}, ` + template.HTML(strconv.Itoa(refreshTime*1000)) + `);
</script>`
	}

	if compress {
		utils.CompressedContent(&p.Content)
	}

	return p
}

// GetPanelFn 是获取面板的函数类型
type GetPanelFn func(ctx interface{}) (Panel, error)

// GetPanelInfoFn 是获取面板信息的函数类型
type GetPanelInfoFn func(ctx *context.Context) (Panel, error)

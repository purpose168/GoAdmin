// Copyright 2019 GoAdmin Core Team. All rights reserved.
// Use of this source code is governed by a Apache-2.0 style
// license that can be found in the LICENSE file.

package template

import (
	"bytes"
	"errors"
	"html/template"
	"path"
	"plugin"
	"strconv"
	"strings"
	"sync"

	"github.com/purpose168/GoAdmin/context"
	c "github.com/purpose168/GoAdmin/modules/config"
	errors2 "github.com/purpose168/GoAdmin/modules/errors"
	"github.com/purpose168/GoAdmin/modules/language"
	"github.com/purpose168/GoAdmin/modules/logger"
	"github.com/purpose168/GoAdmin/modules/menu"
	"github.com/purpose168/GoAdmin/modules/system"
	"github.com/purpose168/GoAdmin/modules/utils"
	"github.com/purpose168/GoAdmin/plugins/admin/models"
	"github.com/purpose168/GoAdmin/template/login"
	"github.com/purpose168/GoAdmin/template/types"
	"golang.org/x/text/cases"
	textLang "golang.org/x/text/language"
)

// Template 是包含UI组件方法的接口
// 它将在插件中用于自定义UI
type Template interface {
	Name() string

	// 组件

	// 布局
	Col() types.ColAttribute
	Row() types.RowAttribute

	// 表单和表格
	Form() types.FormAttribute
	Table() types.TableAttribute
	DataTable() types.DataTableAttribute

	TreeView() types.TreeViewAttribute
	Tree() types.TreeAttribute
	Tabs() types.TabsAttribute
	Alert() types.AlertAttribute
	Link() types.LinkAttribute

	Paginator() types.PaginatorAttribute
	Popup() types.PopupAttribute
	Box() types.BoxAttribute

	Label() types.LabelAttribute
	Image() types.ImgAttribute

	Button() types.ButtonAttribute

	// 构建器方法
	GetTmplList() map[string]string
	GetAssetList() []string
	GetAssetImportHTML(exceptComponents ...string) template.HTML
	GetAsset(string) ([]byte, error)
	GetTemplate(bool) (*template.Template, string)
	GetVersion() string
	GetRequirements() []string
	GetHeadHTML() template.HTML
	GetFootJS() template.HTML
	Get404HTML() template.HTML
	Get500HTML() template.HTML
	Get403HTML() template.HTML
}

// PageType 页面类型枚举
type PageType uint8

const (
	NormalPage          PageType = iota // 正常页面
	Missing404Page                      // 404页面
	Error500Page                        // 500错误页面
	NoPermission403Page                 // 403无权限页面
)

// GetPageTypeFromPageError 根据页面错误获取页面类型
// 参数:
//   - err: 页面错误对象
//
// 返回: 对应的页面类型
func GetPageTypeFromPageError(err errors2.PageError) PageType {
	if err == nil {
		return NormalPage
	} else if err == errors2.PageError403 {
		return NoPermission403Page
	} else if err == errors2.PageError404 {
		return Missing404Page
	} else {
		return Error500Page
	}
}

// 组件类型常量定义
const (
	CompCol       = "col"       // 列组件
	CompRow       = "row"       // 行组件
	CompForm      = "form"      // 表单组件
	CompTable     = "table"     // 表格组件
	CompDataTable = "datatable" // 数据表格组件
	CompTree      = "tree"      // 树组件
	CompTreeView  = "treeview"  // 树视图组件
	CompTabs      = "tabs"      // 标签页组件
	CompAlert     = "alert"     // 警告组件
	CompLink      = "link"      // 链接组件
	CompPaginator = "paginator" // 分页组件
	CompPopup     = "popup"     // 弹窗组件
	CompBox       = "box"       // 盒子组件
	CompLabel     = "label"     // 标签组件
	CompImage     = "image"     // 图片组件
	CompButton    = "button"    // 按钮组件
)

// HTML 将字符串转换为HTML类型
// 参数:
//   - s: HTML字符串
//
// 返回: HTML类型
func HTML(s string) template.HTML {
	return template.HTML(s)
}

// CSS 将字符串转换为CSS类型
// 参数:
//   - s: CSS字符串
//
// 返回: CSS类型
func CSS(s string) template.CSS {
	return template.CSS(s)
}

// JS 将字符串转换为JS类型
// 参数:
//   - s: JavaScript字符串
//
// 返回: JS类型
func JS(s string) template.JS {
	return template.JS(s)
}

// templateMap 包含已注册的模板
var templateMap = make(map[string]Template)

// Get 根据主题名称获取模板接口
// 如果找不到名称，会触发panic
// 参数:
//   - ctx: 上下文对象
//   - theme: 主题名称
//
// 返回: 模板接口
func Get(ctx *context.Context, theme string) Template {
	if ctx != nil {
		queryTheme := ctx.Theme()
		if queryTheme != "" {
			if temp, ok := templateMap[queryTheme]; ok {
				return temp
			}
		}
	}
	if temp, ok := templateMap[theme]; ok {
		return temp
	}
	panic("wrong theme name")
}

// Default 获取使用全局配置设置的主题名称的默认模板
// 如果找不到名称，会触发panic
// 参数:
//   - ctx: 可选的上下文对象
//
// 返回: 模板接口
func Default(ctx ...*context.Context) Template {
	if len(ctx) > 0 && ctx[0] != nil {
		queryTheme := ctx[0].Theme()
		if queryTheme != "" {
			if temp, ok := templateMap[queryTheme]; ok {
				return temp
			}
		}
	}
	if temp, ok := templateMap[c.GetTheme()]; ok {
		return temp
	}
	panic("wrong theme name")
}

var (
	templateMu sync.Mutex
	compMu     sync.Mutex
)

// Add 通过提供的主题名称使模板可用
// 如果使用相同的名称调用两次Add或模板为nil，会触发panic
// 参数:
//   - name: 主题名称
//   - temp: 模板接口
func Add(name string, temp Template) {
	templateMu.Lock()
	defer templateMu.Unlock()
	if temp == nil {
		panic("template is nil")
	}
	if _, dup := templateMap[name]; dup {
		panic("add template twice " + name)
	}
	templateMap[name] = temp
}

// CheckRequirements 检查主题和GoAdmin的相互依赖限制
// 第一个返回参数表示GoAdmin版本是否满足所使用主题的要求
// 第二个返回参数表示所使用主题的版本是否满足GoAdmin的要求
// 返回:
//   - bool: GoAdmin版本是否满足主题要求
//   - bool: 主题版本是否满足GoAdmin要求
func CheckRequirements() (bool, bool) {
	if !CheckThemeRequirements() {
		return false, true
	}
	// 不在默认官方主题中的主题将被忽略
	if !utils.InArray(DefaultThemeNames, Default().Name()) {
		return true, true
	}
	return true, VersionCompare(Default().GetVersion(), system.RequireThemeVersion()[Default().Name()])
}

// CheckThemeRequirements 检查主题要求
// 返回: GoAdmin版本是否满足主题要求
func CheckThemeRequirements() bool {
	return VersionCompare(system.Version(), Default().GetRequirements())
}

// VersionCompare 比较版本号
// 参数:
//   - toCompare: 要比较的版本号
//   - versions: 版本号列表
//
// 返回: 是否匹配或满足版本要求
func VersionCompare(toCompare string, versions []string) bool {
	for _, v := range versions {
		if v == toCompare || utils.CompareVersion(v, toCompare) {
			return true
		}
	}
	return false
}

// GetPageContentFromPageType 根据页面类型获取页面内容
// 参数:
//   - ctx: 上下文对象
//   - title: 页面标题
//   - desc: 页面描述
//   - msg: 消息内容
//   - pt: 页面类型
//
// 返回: 标题HTML、描述HTML、内容HTML
func GetPageContentFromPageType(ctx *context.Context, title, desc, msg string, pt PageType) (template.HTML, template.HTML, template.HTML) {
	if c.GetDebug() {
		return template.HTML(title), template.HTML(desc), Default(ctx).Alert().SetTitle(errors2.MsgWithIcon).Warning(msg)
	}

	if pt == Missing404Page {
		if c.GetCustom404HTML() != template.HTML("") {
			return "", "", c.GetCustom404HTML()
		} else {
			return "", "", Default(ctx).Get404HTML()
		}
	} else if pt == NoPermission403Page {
		if c.GetCustom404HTML() != template.HTML("") {
			return "", "", c.GetCustom403HTML()
		} else {
			return "", "", Default(ctx).Get403HTML()
		}
	} else {
		if c.GetCustom500HTML() != template.HTML("") {
			return "", "", c.GetCustom500HTML()
		} else {
			return "", "", Default(ctx).Get500HTML()
		}
	}
}

// DefaultThemeNames 默认主题名称列表
var DefaultThemeNames = []string{"sword", "adminlte"}

// Themes 获取所有已注册的主题名称
// 返回: 主题名称列表
func Themes() []string {
	names := make([]string, len(templateMap))
	i := 0
	for k := range templateMap {
		names[i] = k
		i++
	}
	return names
}

// AddFromPlugin 从插件添加模板
// 参数:
//   - name: 模板名称
//   - mod: 插件模块路径
func AddFromPlugin(name string, mod string) {

	plug, err := plugin.Open(mod)
	if err != nil {
		logger.Error("AddFromPlugin err", err)
		panic(err)
	}

	tempPlugin, err := plug.Lookup(cases.Title(textLang.Und).String(name))
	if err != nil {
		logger.Error("AddFromPlugin err", err)
		panic(err)
	}

	var temp Template
	temp, ok := tempPlugin.(Template)
	if !ok {
		logger.Error("AddFromPlugin err: unexpected type from module symbol")
		panic(errors.New("AddFromPlugin err: unexpected type from module symbol"))
	}

	Add(name, temp)
}

// Component 是表示UI组件的接口
type Component interface {
	// GetTemplate 返回一个 *template.Template 和一个给定的键
	GetTemplate() (*template.Template, string)

	// GetAssetList 返回组件中使用的资源URL后缀
	// 示例:
	//
	// {{.UrlPrefix}}/assets/login/css/bootstrap.min.css => login/css/bootstrap.min.css
	//
	// 参见:
	// https://github.com/purpose168/GoAdmin/blob/master/template/login/theme1.tmpl#L32
	// https://github.com/purpose168/GoAdmin/blob/master/template/login/list.go
	GetAssetList() []string

	// GetAsset 根据相应的URL后缀返回资源内容
	// 建议使用go-bindata工具生成资源内容
	//
	// 参见: http://github.com/jteeuwen/go-bindata
	GetAsset(string) ([]byte, error)

	GetContent() template.HTML

	IsAPage() bool

	GetName() string

	GetJS() template.JS
	GetCSS() template.CSS
	GetCallbacks() types.Callbacks
}

// compMap 组件映射表
var compMap = map[string]Component{
	"login": login.GetLoginComponent(),
}

// GetComp 根据注册名称获取组件
// 如果找不到名称，会触发panic
// 参数:
//   - name: 组件名称
//
// 返回: 组件接口
func GetComp(name string) Component {
	if comp, ok := compMap[name]; ok {
		return comp
	}
	panic("wrong component name")
}

// GetComponentAsset 获取所有组件的资源列表
// 返回: 资源URL后缀列表
func GetComponentAsset() []string {
	assets := make([]string, 0)
	for _, comp := range compMap {
		assets = append(assets, comp.GetAssetList()...)
	}
	return assets
}

// GetComponentAssetWithinPage 获取页面内组件的资源列表（不包括页面组件）
// 返回: 资源URL后缀列表
func GetComponentAssetWithinPage() []string {
	assets := make([]string, 0)
	for _, comp := range compMap {
		if !comp.IsAPage() {
			assets = append(assets, comp.GetAssetList()...)
		}
	}
	return assets
}

// GetComponentAssetImportHTML 获取组件资源导入HTML
// 参数:
//   - ctx: 上下文对象
//
// 返回: 资源导入HTML
func GetComponentAssetImportHTML(ctx *context.Context) (res template.HTML) {
	res = Default(ctx).GetAssetImportHTML(c.GetExcludeThemeComponents()...)
	assets := GetComponentAssetWithinPage()
	for i := 0; i < len(assets); i++ {
		res += getHTMLFromAssetUrl(assets[i])
	}
	return
}

// getHTMLFromAssetUrl 根据资源URL获取对应的HTML标签
// 参数:
//   - s: 资源URL后缀
//
// 返回: HTML标签
func getHTMLFromAssetUrl(s string) template.HTML {
	switch path.Ext(s) {
	case ".css":
		return template.HTML(`<link rel="stylesheet" href="` + c.GetAssetUrl() + c.Url("/assets"+s) + `">`)
	case ".js":
		return template.HTML(`<script src="` + c.GetAssetUrl() + c.Url("/assets"+s) + `"></script>`)
	default:
		return ""
	}
}

// GetAsset 根据路径获取资源内容
// 参数:
//   - path: 资源路径
//
// 返回: 资源内容和错误
func GetAsset(path string) ([]byte, error) {
	for _, comp := range compMap {
		res, err := comp.GetAsset(path)
		if err == nil {
			return res, nil
		}
	}
	return nil, errors.New(path + " not found")
}

// AddComp 通过提供的名称使组件可用
// 如果使用相同的名称调用两次Add或组件为nil，会触发panic
// 参数:
//   - comp: 组件接口
func AddComp(comp Component) {
	compMu.Lock()
	defer compMu.Unlock()
	if comp == nil {
		panic("component is nil")
	}
	if _, dup := compMap[comp.GetName()]; dup {
		panic("add component twice " + comp.GetName())
	}
	compMap[comp.GetName()] = comp
}

// AddLoginComp 添加指定的登录组件
// 参数:
//   - comp: 组件接口
func AddLoginComp(comp Component) {
	compMu.Lock()
	defer compMu.Unlock()
	compMap["login"] = comp
}

// SetComp 通过提供的名称使组件可用
// 如果键对应的值为空或组件为nil，会触发panic
// 参数:
//   - name: 组件名称
//   - comp: 组件接口
func SetComp(name string, comp Component) {
	compMu.Lock()
	defer compMu.Unlock()
	if comp == nil {
		panic("component is nil")
	}
	if _, dup := compMap[name]; dup {
		compMap[name] = comp
	}
}

// ExecuteParam 执行参数结构体
type ExecuteParam struct {
	User       models.UserModel   // 用户模型
	Tmpl       *template.Template // 模板对象
	TmplName   string             // 模板名称
	IsPjax     bool               // 是否为PJAX请求
	Panel      types.Panel        // 面板
	Logo       template.HTML      // Logo HTML
	Config     *c.Config          // 配置对象
	Menu       *menu.Menu         // 菜单
	Animation  bool               // 是否启用动画
	Buttons    types.Buttons      // 按钮
	NoCompress bool               // 是否不压缩
	Iframe     bool               // 是否在iframe中
}

// updateNavAndLogoJS 更新导航和Logo的JavaScript
// 参数:
//   - logo: Logo HTML
//
// 返回: JavaScript代码
func updateNavAndLogoJS(logo template.HTML) template.JS {
	if logo == template.HTML("") {
		return ""
	}
	return `$(function () {
	$(".logo-lg").html("` + template.JS(logo) + `");
});`
}

// updateNavJS 更新导航的JavaScript
// 参数:
//   - isPjax: 是否为PJAX请求
//
// 返回: JavaScript代码
func updateNavJS(isPjax bool) template.JS {
	if !isPjax {
		return ""
	}
	return `$(function () {
	let lis = $(".user-menu .dropdown-menu li");
	for (i = 0; i < lis.length - 2; i++) {
		$(lis[i]).remove();
	}
	$(".user-menu .dropdown-menu").prepend($("#navbar-nav-custom").html());
});`
}

// ExecuteOptions 执行选项结构体
type ExecuteOptions struct {
	Animation         bool                           // 是否启用动画
	NoCompress        bool                           // 是否不压缩
	HideSideBar       bool                           // 是否隐藏侧边栏
	HideHeader        bool                           // 是否隐藏头部
	UpdateMenu        bool                           // 是否更新菜单
	NavDropDownButton []*types.NavDropDownItemButton // 导航下拉按钮
}

// GetExecuteOptions 获取执行选项
// 参数:
//   - options: 执行选项列表
//
// 返回: 执行选项
func GetExecuteOptions(options []ExecuteOptions) ExecuteOptions {
	if len(options) == 0 {
		return ExecuteOptions{Animation: true}
	}
	return options[0]
}

// Execute 执行模板渲染
// 参数:
//   - ctx: 上下文对象
//   - param: 执行参数
//
// 返回: 渲染后的缓冲区
func Execute(ctx *context.Context, param *ExecuteParam) *bytes.Buffer {

	buf := new(bytes.Buffer)
	err := param.Tmpl.ExecuteTemplate(buf, param.TmplName,
		types.NewPage(ctx, &types.NewPageParam{
			User:       param.User,
			Menu:       param.Menu,
			Assets:     GetComponentAssetImportHTML(ctx),
			Buttons:    param.Buttons,
			Iframe:     param.Iframe,
			UpdateMenu: param.IsPjax,
			Panel: param.Panel.
				GetContent(append([]bool{param.Config.IsProductionEnvironment() && !param.NoCompress},
					param.Animation)...).AddJS(param.Menu.GetUpdateJS(param.IsPjax)).
				AddJS(updateNavAndLogoJS(param.Logo)).AddJS(updateNavJS(param.IsPjax)),
			TmplHeadHTML: Default(ctx).GetHeadHTML(),
			TmplFootJS:   Default(ctx).GetFootJS(),
			Logo:         param.Logo,
		}))
	if err != nil {
		logger.Error("template execute error", err)
	}
	return buf
}

// WarningPanel 创建警告面板
// 参数:
//   - ctx: 上下文对象
//   - msg: 消息内容
//   - pts: 可选的页面类型
//
// 返回: 面板对象
func WarningPanel(ctx *context.Context, msg string, pts ...PageType) types.Panel {
	pt := Error500Page
	if len(pts) > 0 {
		pt = pts[0]
	}
	pageTitle, description, content := GetPageContentFromPageType(ctx, msg, msg, msg, pt)
	return types.Panel{
		Content:     content,
		Description: description,
		Title:       pageTitle,
	}
}

// WarningPanelWithDescAndTitle 创建带描述和标题的警告面板
// 参数:
//   - ctx: 上下文对象
//   - msg: 消息内容
//   - desc: 描述
//   - title: 标题
//   - pts: 可选的页面类型
//
// 返回: 面板对象
func WarningPanelWithDescAndTitle(ctx *context.Context, msg, desc, title string, pts ...PageType) types.Panel {
	pt := Error500Page
	if len(pts) > 0 {
		pt = pts[0]
	}
	pageTitle, description, content := GetPageContentFromPageType(ctx, msg, desc, title, pt)
	return types.Panel{
		Content:     content,
		Description: description,
		Title:       pageTitle,
	}
}

// DefaultFuncMap 默认模板函数映射
var DefaultFuncMap = template.FuncMap{
	"lang":     language.Get,         // 获取语言翻译
	"langHtml": language.GetFromHtml, // 从HTML获取语言翻译
	"link": func(cdnUrl, prefixUrl, assetsUrl string) string {
		// 生成链接URL
		if cdnUrl == "" {
			return prefixUrl + assetsUrl
		}
		return cdnUrl + assetsUrl
	},
	"isLinkUrl": func(s string) bool {
		// 判断是否为链接URL
		return (len(s) > 7 && s[:7] == "http://") || (len(s) > 8 && s[:8] == "https://")
	},
	"render": func(s, old, repl template.HTML) template.HTML {
		// 渲染HTML，替换字符串
		return template.HTML(strings.ReplaceAll(string(s), string(old), string(repl)))
	},
	"renderJS": func(s template.JS, old, repl template.HTML) template.JS {
		// 渲染JavaScript，替换字符串
		return template.JS(strings.ReplaceAll(string(s), string(old), string(repl)))
	},
	"divide": func(a, b int) int {
		// 整数除法
		return a / b
	},
	"renderRowDataHTML": func(id, content template.HTML, value ...map[string]types.InfoItem) template.HTML {
		// 渲染行数据HTML
		return template.HTML(types.ParseTableDataTmplWithID(id, string(content), value...))
	},
	"renderRowDataJS": func(id template.HTML, content template.JS, value ...map[string]types.InfoItem) template.JS {
		// 渲染行数据JavaScript
		return template.JS(types.ParseTableDataTmplWithID(id, string(content), value...))
	},
	"attr": func(s template.HTML) template.HTMLAttr {
		// 转换为HTML属性
		return template.HTMLAttr(s)
	},
	"js": func(s interface{}) template.JS {
		// 转换为JavaScript
		if ss, ok := s.(string); ok {
			return template.JS(ss)
		}
		if ss, ok := s.(template.HTML); ok {
			return template.JS(ss)
		}
		return ""
	},
	"changeValue": func(f types.FormField, index int) types.FormField {
		// 更新表单字段的值
		if len(f.ValueArr) > 0 {
			f.Value = template.HTML(f.ValueArr[index])
		}
		if len(f.OptionsArr) > 0 {
			f.Options = f.OptionsArr[index]
		}
		if f.FormType.IsSelect() {
			f.FieldClass += "_" + strconv.Itoa(index)
		}
		return f
	},
}

// BaseComponent 基础组件结构体
type BaseComponent struct {
	Name      string          // 组件名称
	HTMLData  string          // HTML数据
	CSS       template.CSS    // CSS样式
	JS        template.JS     // JavaScript代码
	Callbacks types.Callbacks // 回调函数
}

// IsAPage 判断是否为页面组件
// 返回: 是否为页面组件
func (b *BaseComponent) IsAPage() bool { return false }

// GetName 获取组件名称
// 返回: 组件名称
func (b *BaseComponent) GetName() string { return b.Name }

// GetAssetList 获取资源列表
// 返回: 资源URL后缀列表
func (b *BaseComponent) GetAssetList() []string { return make([]string, 0) }

// GetAsset 根据名称获取资源
// 参数:
//   - name: 资源名称
//
// 返回: 资源内容和错误
func (b *BaseComponent) GetAsset(name string) ([]byte, error) { return nil, nil }

// GetJS 获取JavaScript代码
// 返回: JavaScript代码
func (b *BaseComponent) GetJS() template.JS { return b.JS }

// GetCSS 获取CSS样式
// 返回: CSS样式
func (b *BaseComponent) GetCSS() template.CSS { return b.CSS }

// GetCallbacks 获取回调函数
// 返回: 回调函数列表
func (b *BaseComponent) GetCallbacks() types.Callbacks { return b.Callbacks }

// BindActionTo 将动作绑定到组件
// 参数:
//   - ctx: 上下文对象
//   - action: 动作对象
//   - id: 按钮ID
func (b *BaseComponent) BindActionTo(ctx *context.Context, action types.Action, id string) {
	action.SetBtnId(id)
	b.JS += action.Js()
	b.HTMLData += string(action.ExtContent(ctx))
	b.Callbacks = append(b.Callbacks, action.GetCallbacks())
}

// GetContentWithData 使用数据获取内容
// 参数:
//   - obj: 数据对象
//
// 返回: HTML内容
func (b *BaseComponent) GetContentWithData(obj interface{}) template.HTML {
	buffer := new(bytes.Buffer)
	tmpl, defineName := b.GetTemplate()
	err := tmpl.ExecuteTemplate(buffer, defineName, obj)
	if err != nil {
		logger.Error(b.Name+" GetContent error:", err)
	}
	return template.HTML(buffer.String())
}

// GetTemplate 获取模板对象
// 返回: 模板对象和模板名称
func (b *BaseComponent) GetTemplate() (*template.Template, string) {
	tmpl, err := template.New(b.Name).
		Funcs(DefaultFuncMap).
		Parse(b.HTMLData)

	if err != nil {
		logger.Error(b.Name+" GetTemplate Error: ", err)
	}

	return tmpl, b.Name
}

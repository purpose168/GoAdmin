package types

import (
	"html/template"
	"net/url"

	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/modules/language"
	"github.com/purpose168/GoAdmin/modules/utils"
	"github.com/purpose168/GoAdmin/plugins/admin/models"
)

// Button 接口定义了按钮的基本行为
type Button interface {
	// Content 返回按钮的HTML内容和JavaScript代码
	Content(ctx *context.Context) (template.HTML, template.JS)
	// GetAction 获取按钮的动作
	GetAction() Action
	// URL 返回按钮的URL
	URL() string
	// METHOD 返回按钮的HTTP方法
	METHOD() string
	// ID 返回按钮的ID
	ID() string
	// Type 返回按钮的类型
	Type() string
	// GetName 返回按钮的名称
	GetName() string
	// SetName 设置按钮的名称
	SetName(name string)
	// IsType 检查按钮是否为指定类型
	IsType(t string) bool
}

// BaseButton 是按钮的基础结构体
type BaseButton struct {
	Id, Url, Method, Name, TypeName string        // 按钮ID、URL、方法、名称和类型名称
	Title                           template.HTML // 按钮标题
	Action                          Action        // 按钮动作
}

// Content 返回空内容
func (b *BaseButton) Content() (template.HTML, template.JS) { return "", "" }

// GetAction 返回按钮动作
func (b *BaseButton) GetAction() Action { return b.Action }

// ID 返回按钮ID
func (b *BaseButton) ID() string { return b.Id }

// URL 返回按钮URL
func (b *BaseButton) URL() string { return b.Url }

// Type 返回按钮类型
func (b *BaseButton) Type() string { return b.TypeName }

// IsType 检查按钮是否为指定类型
func (b *BaseButton) IsType(t string) bool { return b.TypeName == t }

// METHOD 返回HTTP方法
func (b *BaseButton) METHOD() string { return b.Method }

// GetName 返回按钮名称
func (b *BaseButton) GetName() string { return b.Name }

// SetName 设置按钮名称
func (b *BaseButton) SetName(name string) { b.Name = name }

// DefaultButton 是默认按钮结构体
type DefaultButton struct {
	*BaseButton
	Color     template.HTML // 按钮背景色
	TextColor template.HTML // 按钮文字颜色
	Icon      string        // 按钮图标
	Direction template.HTML // 按钮方向
	Group     bool          // 是否为按钮组
}

// GetDefaultButton 创建默认按钮
func GetDefaultButton(title template.HTML, icon string, action Action, colors ...template.HTML) *DefaultButton {
	return defaultButton(title, "right", icon, action, false, colors...)
}

// GetDefaultButtonGroup 创建默认按钮组
func GetDefaultButtonGroup(title template.HTML, icon string, action Action, colors ...template.HTML) *DefaultButton {
	return defaultButton(title, "right", icon, action, true, colors...)
}

// defaultButton 创建默认按钮的内部函数
func defaultButton(title, direction template.HTML, icon string, action Action, group bool, colors ...template.HTML) *DefaultButton {
	id := btnUUID()
	action.SetBtnId("." + id)

	var color, textColor template.HTML
	if len(colors) > 0 {
		color = colors[0]
	}
	if len(colors) > 1 {
		textColor = colors[1]
	}
	node := action.GetCallbacks()
	return &DefaultButton{
		BaseButton: &BaseButton{
			Id:     id,
			Title:  title,
			Action: action,
			Url:    node.Path,
			Method: node.Method,
		},
		Group:     group,
		Color:     color,
		TextColor: textColor,
		Icon:      icon,
		Direction: direction,
	}
}

// GetColumnButton 创建列按钮
func GetColumnButton(title template.HTML, icon string, action Action, colors ...template.HTML) *DefaultButton {
	return defaultButton(title, "", icon, action, true, colors...)
}

// Content 生成按钮的HTML内容和JavaScript代码
func (b *DefaultButton) Content(ctx *context.Context) (template.HTML, template.JS) {

	color := template.HTML("")
	if b.Color != template.HTML("") {
		color = template.HTML(`background-color:`) + b.Color + template.HTML(`;`)
	}
	textColor := template.HTML("")
	if b.TextColor != template.HTML("") {
		textColor = template.HTML(`color:`) + b.TextColor + template.HTML(`;`)
	}

	style := template.HTML("")
	addColor := color + textColor

	if addColor != template.HTML("") {
		style = template.HTML(`style="`) + addColor + template.HTML(`"`)
	}

	h := template.HTML("")
	if b.Group {
		h += `<div class="btn-group pull-` + b.Direction + `" style="margin-right: 10px">`
	}

	h += `<a ` + style + ` class="` + template.HTML(b.Id) + ` btn btn-sm btn-default ` + b.Action.BtnClass() + `" ` + b.Action.BtnAttribute() + `>
                    <i class="fa ` + template.HTML(b.Icon) + `"></i>&nbsp;&nbsp;` + b.Title + `
                </a>`
	if b.Group {
		h += `</div>`
	}
	return h + b.Action.ExtContent(ctx), b.Action.Js()
}

// ActionButton 是动作按钮结构体
type ActionButton struct {
	*BaseButton
}

// GetActionButton 创建动作按钮
func GetActionButton(title template.HTML, action Action, ids ...string) *ActionButton {

	id := ""
	if len(ids) > 0 {
		id = ids[0]
	} else {
		id = "action-info-btn-" + utils.Uuid(10)
	}

	action.SetBtnId("." + id)
	node := action.GetCallbacks()

	return &ActionButton{
		BaseButton: &BaseButton{
			Id:     id,
			Title:  title,
			Action: action,
			Url:    node.Path,
			Method: node.Method,
		},
	}
}

// Content 生成动作按钮的HTML内容和JavaScript代码
func (b *ActionButton) Content(ctx *context.Context) (template.HTML, template.JS) {
	h := template.HTML(`<li style="cursor: pointer;"><a data-id="{{.Id}}" class="`+template.HTML(b.Id)+` `+
		b.Action.BtnClass()+`" `+b.Action.BtnAttribute()+`>`+b.Title+`</a></li>`) + b.Action.ExtContent(ctx)
	return h, b.Action.Js()
}

// ActionIconButton 是动作图标按钮结构体
type ActionIconButton struct {
	Icon template.HTML // 图标
	*BaseButton
}

// GetActionIconButton 创建动作图标按钮
func GetActionIconButton(icon string, action Action, ids ...string) *ActionIconButton {

	id := ""
	if len(ids) > 0 {
		id = ids[0]
	} else {
		id = "action-info-btn-" + utils.Uuid(10)
	}

	action.SetBtnId("." + id)
	node := action.GetCallbacks()

	return &ActionIconButton{
		Icon: template.HTML(icon),
		BaseButton: &BaseButton{
			Id:     id,
			Action: action,
			Url:    node.Path,
			Method: node.Method,
		},
	}
}

// Content 生成动作图标按钮的HTML内容和JavaScript代码
func (b *ActionIconButton) Content(ctx *context.Context) (template.HTML, template.JS) {
	h := template.HTML(`<a data-id="{{.Id}}" class="`+template.HTML(b.Id)+` `+
		b.Action.BtnClass()+`" `+b.Action.BtnAttribute()+`><i class="fa `+b.Icon+`" style="font-size: 16px;"></i></a>`) + b.Action.ExtContent(ctx)
	return h, b.Action.Js()
}

// Buttons 是按钮切片类型
type Buttons []Button

// Add 添加按钮
func (b Buttons) Add(btn Button) Buttons {
	return append(b, btn)
}

// Content 生成所有按钮的HTML内容和JavaScript代码
func (b Buttons) Content(ctx *context.Context) (template.HTML, template.JS) {
	h := template.HTML("")
	j := template.JS("")

	for _, btn := range b {
		hh, jj := btn.Content(ctx)
		h += hh
		j += jj
	}
	return h, j
}

// Copy 复制按钮切片
func (b Buttons) Copy() Buttons {
	var c = make(Buttons, len(b))
	copy(c, b)
	return c
}

// FooterContent 生成所有按钮的页脚内容
func (b Buttons) FooterContent(ctx *context.Context) template.HTML {
	footer := template.HTML("")

	for _, btn := range b {
		footer += btn.GetAction().FooterContent(ctx)
	}
	return footer
}

// CheckPermission 检查用户是否有权限访问按钮
func (b Buttons) CheckPermission(user models.UserModel) Buttons {
	btns := make(Buttons, 0)
	for _, btn := range b {
		if btn.IsType(ButtonTypeNavDropDown) {
			items := make([]Button, 0)
			for _, navItem := range btn.(*NavDropDownButton).Items {
				if user.CheckPermissionByUrlMethod(btn.URL(), btn.METHOD(), url.Values{}) {
					items = append(items, navItem)
				}
			}
			if len(items) > 0 {
				btns = append(btns, btn)
			}
		} else if user.CheckPermissionByUrlMethod(btn.URL(), btn.METHOD(), url.Values{}) {
			btns = append(btns, btn)
		}
	}
	return btns
}

// CheckPermissionWhenURLAndMethodNotEmpty 当URL和方法不为空时检查权限
func (b Buttons) CheckPermissionWhenURLAndMethodNotEmpty(user models.UserModel) Buttons {
	btns := make(Buttons, 0)
	for _, b := range b {
		if b.URL() == "" || b.METHOD() == "" || user.CheckPermissionByUrlMethod(b.URL(), b.METHOD(), url.Values{}) {
			btns = append(btns, b)
		}
	}
	return btns
}

// AddNavButton 添加导航按钮
func (b Buttons) AddNavButton(ico, name string, action Action) Buttons {
	if !b.CheckExist(name) {
		return append(b, GetNavButton(language.GetFromHtml(template.HTML(name)), ico, action, name))
	}
	return b
}

// RemoveButtonByName 根据名称移除按钮
func (b Buttons) RemoveButtonByName(name string) Buttons {
	if name == "" {
		return b
	}

	for i := 0; i < len(b); i++ {
		if b[i].GetName() == name {
			b = append(b[:i], b[i+1:]...)
		}
	}
	return b
}

// CheckExist 检查按钮是否存在
func (b Buttons) CheckExist(name string) bool {
	if name == "" {
		return false
	}
	for i := 0; i < len(b); i++ {
		if b[i].GetName() == name {
			return true
		}
	}
	return false
}

// Callbacks 获取所有按钮的回调函数
func (b Buttons) Callbacks() []context.Node {
	cbs := make([]context.Node, 0)
	for _, btn := range b {
		cbs = append(cbs, btn.GetAction().GetCallbacks())
	}
	return cbs
}

const (
	NavBtnSiteName = "site setting"       // 站点设置按钮名称
	NavBtnInfoName = "site info"          // 站点信息按钮名称
	NavBtnToolName = "code generate tool" // 代码生成工具按钮名称
	NavBtnPlugName = "plugins"            // 插件按钮名称
)

// RemoveSiteNavButton 移除站点设置导航按钮
func (b Buttons) RemoveSiteNavButton() Buttons {
	return b.RemoveButtonByName(NavBtnSiteName)
}

// RemoveInfoNavButton 移除站点信息导航按钮
func (b Buttons) RemoveInfoNavButton() Buttons {
	return b.RemoveButtonByName(NavBtnInfoName)
}

// RemoveToolNavButton 移除工具导航按钮
func (b Buttons) RemoveToolNavButton() Buttons {
	return b.RemoveButtonByName(NavBtnToolName)
}

// RemovePlugNavButton 移除插件导航按钮
func (b Buttons) RemovePlugNavButton() Buttons {
	return b.RemoveButtonByName(NavBtnPlugName)
}

// NavButton 是导航按钮结构体
type NavButton struct {
	*BaseButton
	Icon string // 图标
}

// GetNavButton 创建导航按钮
func GetNavButton(title template.HTML, icon string, action Action, names ...string) *NavButton {

	id := btnUUID()
	action.SetBtnId("." + id)
	node := action.GetCallbacks()
	name := ""

	if len(names) > 0 {
		name = names[0]
	}

	return &NavButton{
		BaseButton: &BaseButton{
			Id:     id,
			Title:  title,
			Action: action,
			Url:    node.Path,
			Method: node.Method,
			Name:   name,
		},
		Icon: icon,
	}
}

// Content 生成导航按钮的HTML内容和JavaScript代码
func (n *NavButton) Content(ctx *context.Context) (template.HTML, template.JS) {

	ico := template.HTML("")
	title := template.HTML("")

	if n.Icon != "" {
		ico = template.HTML(`<i class="fa ` + n.Icon + `"></i>`)
	}

	if n.Title != "" {
		title = `<span>` + n.Title + `</span>`
	}

	h := template.HTML(`<li>
    <a class="`+template.HTML(n.Id)+` `+n.Action.BtnClass()+` dropdown-item" `+n.Action.BtnAttribute()+`>
      `+ico+`
      `+title+`
    </a>
</li>`) + n.Action.ExtContent(ctx)
	return h, n.Action.Js()
}

// NavDropDownButton 是导航下拉按钮结构体
type NavDropDownButton struct {
	*BaseButton
	Icon  string                   // 图标
	Items []*NavDropDownItemButton // 下拉菜单项
}

// NavDropDownItemButton 是导航下拉菜单项按钮结构体
type NavDropDownItemButton struct {
	*BaseButton
}

// GetDropDownButton 创建下拉按钮
func GetDropDownButton(title template.HTML, icon string, items []*NavDropDownItemButton, names ...string) *NavDropDownButton {
	id := btnUUID()
	name := ""

	if len(names) > 0 {
		name = names[0]
	}

	return &NavDropDownButton{
		BaseButton: &BaseButton{
			Id:       id,
			Title:    title,
			Name:     name,
			TypeName: ButtonTypeNavDropDown,
			Action:   new(NilAction),
		},
		Items: items,
		Icon:  icon,
	}
}

// SetItems 设置下拉菜单项
func (n *NavDropDownButton) SetItems(items []*NavDropDownItemButton) {
	n.Items = items
}

// AddItem 添加下拉菜单项
func (n *NavDropDownButton) AddItem(item *NavDropDownItemButton) {
	n.Items = append(n.Items, item)
}

// Content 生成下拉按钮的HTML内容和JavaScript代码
func (n *NavDropDownButton) Content(ctx *context.Context) (template.HTML, template.JS) {

	ico := template.HTML("")
	title := template.HTML("")

	if n.Icon != "" {
		ico = template.HTML(`<i class="fa ` + n.Icon + `"></i>`)
	}

	if n.Title != "" {
		title = `<span>` + n.Title + `</span>`
	}

	content := template.HTML("")
	js := template.JS("")

	for _, item := range n.Items {
		c, j := item.Content(ctx)
		content += c
		js += j
	}

	did := utils.Uuid(10)

	h := template.HTML(`<li class="dropdown" id="` + template.HTML(did) + `">
    <a class="` + template.HTML(n.Id) + ` dropdown-toggle" data-toggle="dropdown" style="cursor:pointer;">
      ` + ico + `
      ` + title + `
    </a>
	<ul class="dropdown-menu"  aria-labelledby="` + template.HTML(did) + `">
    	` + content + `
	</ul>
</li>`)

	return h, js
}

const (
	ButtonTypeNavDropDownItem = "navdropdownitem" // 导航下拉菜单项按钮类型
	ButtonTypeNavDropDown     = "navdropdown"     // 导航下拉按钮类型
)

// GetDropDownItemButton 创建下拉菜单项按钮
func GetDropDownItemButton(title template.HTML, action Action, names ...string) *NavDropDownItemButton {
	id := btnUUID()
	action.SetBtnId("." + id)
	node := action.GetCallbacks()
	name := ""

	if len(names) > 0 {
		name = names[0]
	}

	return &NavDropDownItemButton{
		BaseButton: &BaseButton{
			Id:       id,
			Title:    title,
			Action:   action,
			Url:      node.Path,
			Method:   node.Method,
			Name:     name,
			TypeName: ButtonTypeNavDropDownItem,
		},
	}
}

// Content 生成下拉菜单项按钮的HTML内容和JavaScript代码
func (n *NavDropDownItemButton) Content(ctx *context.Context) (template.HTML, template.JS) {

	title := template.HTML("")

	if n.Title != "" {
		title = `<span>` + n.Title + `</span>`
	}

	h := template.HTML(`<li><a class="dropdown-item `+template.HTML(n.Id)+` `+
		n.Action.BtnClass()+`" `+n.Action.BtnAttribute()+`>
      `+title+`
</a></li>`) + n.Action.ExtContent(ctx)
	return h, n.Action.Js()
}

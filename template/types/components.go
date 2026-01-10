// 版权所有 2019 GoAdmin 核心团队。保留所有权利。
// 本源代码的使用受 Apache-2.0 风格许可证管辖
// 该许可证可在 LICENSE 文件中找到。

package types

import (
	"html/template"

	"github.com/purpose168/GoAdmin/modules/menu"
	"github.com/purpose168/GoAdmin/plugins/admin/modules"
	"github.com/purpose168/GoAdmin/template/types/form"
)

// FormAttribute 表单属性接口
type FormAttribute interface {
	SetHeader(value template.HTML) FormAttribute                // 设置表单头部
	SetContent(value FormFields) FormAttribute                  // 设置表单内容
	SetTabContents(value []FormFields) FormAttribute            // 设置标签页内容
	SetTabHeaders(value []string) FormAttribute                 // 设置标签页头部
	SetFooter(value template.HTML) FormAttribute                // 设置表单页脚
	SetPrefix(value string) FormAttribute                       // 设置前缀
	SetUrl(value string) FormAttribute                          // 设置URL
	SetPrimaryKey(value string) FormAttribute                   // 设置主键
	SetHorizontal(value bool) FormAttribute                     // 设置水平布局
	SetId(id string) FormAttribute                              // 设置ID
	SetAjax(successJS, errorJS template.JS) FormAttribute       // 设置AJAX成功和失败的JavaScript
	SetHiddenFields(fields map[string]string) FormAttribute     // 设置隐藏字段
	SetFieldsHTML(html template.HTML) FormAttribute             // 设置字段HTML
	SetMethod(value string) FormAttribute                       // 设置HTTP方法
	SetHeadWidth(width int) FormAttribute                       // 设置头部宽度
	SetInputWidth(width int) FormAttribute                      // 设置输入框宽度
	SetTitle(value template.HTML) FormAttribute                 // 设置标题
	SetLayout(layout form.Layout) FormAttribute                 // 设置布局
	SetOperationFooter(value template.HTML) FormAttribute       // 设置操作页脚
	GetDefaultBoxHeader(hideBack bool) template.HTML            // 获取默认盒子头部
	GetDetailBoxHeader(editUrl, deleteUrl string) template.HTML // 获取详情盒子头部
	GetBoxHeaderNoButton() template.HTML                        // 获取无按钮的盒子头部
	GetContent() template.HTML                                  // 获取内容
}

// BoxAttribute 盒子属性接口
type BoxAttribute interface {
	SetHeader(value template.HTML) BoxAttribute       // 设置头部
	SetBody(value template.HTML) BoxAttribute         // 设置主体
	SetNoPadding() BoxAttribute                       // 设置无内边距
	SetFooter(value template.HTML) BoxAttribute       // 设置页脚
	SetTitle(value template.HTML) BoxAttribute        // 设置标题
	WithHeadBorder() BoxAttribute                     // 设置头部边框
	SetIframeStyle(iframe bool) BoxAttribute          // 设置iframe样式
	SetAttr(attr template.HTMLAttr) BoxAttribute      // 设置属性
	SetStyle(value template.HTMLAttr) BoxAttribute    // 设置样式
	SetHeadColor(value string) BoxAttribute           // 设置头部颜色
	SetClass(value string) BoxAttribute               // 设置类名
	SetTheme(value string) BoxAttribute               // 设置主题
	SetSecondHeader(value template.HTML) BoxAttribute // 设置第二头部
	SetSecondHeadColor(value string) BoxAttribute     // 设置第二头部颜色
	WithSecondHeadBorder() BoxAttribute               // 设置第二头部边框
	SetSecondHeaderClass(value string) BoxAttribute   // 设置第二头部类名
	GetContent() template.HTML                        // 获取内容
}

// ColAttribute 列属性接口
type ColAttribute interface {
	SetSize(value S) ColAttribute                // 设置大小
	SetContent(value template.HTML) ColAttribute // 设置内容
	AddContent(value template.HTML) ColAttribute // 添加内容
	GetContent() template.HTML                   // 获取内容
}

// ImgAttribute 图片属性接口
type ImgAttribute interface {
	SetWidth(value string) ImgAttribute      // 设置宽度
	SetHeight(value string) ImgAttribute     // 设置高度
	WithModal() ImgAttribute                 // 使用模态框
	SetSrc(value template.HTML) ImgAttribute // 设置源地址
	GetContent() template.HTML               // 获取内容
}

// LabelAttribute 标签属性接口
type LabelAttribute interface {
	SetContent(value template.HTML) LabelAttribute // 设置内容
	SetColor(value template.HTML) LabelAttribute   // 设置颜色
	SetType(value string) LabelAttribute           // 设置类型
	GetContent() template.HTML                     // 获取内容
}

// RowAttribute 行属性接口
type RowAttribute interface {
	SetContent(value template.HTML) RowAttribute // 设置内容
	AddContent(value template.HTML) RowAttribute // 添加内容
	GetContent() template.HTML                   // 获取内容
}

// ButtonAttribute 按钮属性接口
type ButtonAttribute interface {
	SetContent(value template.HTML) ButtonAttribute     // 设置内容
	SetOrientationRight() ButtonAttribute               // 设置右对齐
	SetOrientationLeft() ButtonAttribute                // 设置左对齐
	SetMarginLeft(int) ButtonAttribute                  // 设置左边距
	SetMarginRight(int) ButtonAttribute                 // 设置右边距
	SetThemePrimary() ButtonAttribute                   // 设置主主题
	SetSmallSize() ButtonAttribute                      // 设置小尺寸
	AddClass(class string) ButtonAttribute              // 添加类名
	SetID(id string) ButtonAttribute                    // 设置ID
	SetMiddleSize() ButtonAttribute                     // 设置中等尺寸
	SetHref(string) ButtonAttribute                     // 设置链接地址
	SetThemeWarning() ButtonAttribute                   // 设置警告主题
	SetTheme(value string) ButtonAttribute              // 设置主题
	SetLoadingText(value template.HTML) ButtonAttribute // 设置加载文本
	SetThemeDefault() ButtonAttribute                   // 设置默认主题
	SetType(string) ButtonAttribute                     // 设置类型
	GetContent() template.HTML                          // 获取内容
}

// TableAttribute 表格属性接口
type TableAttribute interface {
	SetThead(value Thead) TableAttribute                    // 设置表头
	SetInfoList(value []map[string]InfoItem) TableAttribute // 设置信息列表
	SetType(value string) TableAttribute                    // 设置类型
	SetName(name string) TableAttribute                     // 设置名称
	SetMinWidth(value string) TableAttribute                // 设置最小宽度
	SetHideThead() TableAttribute                           // 设置隐藏表头
	SetSticky(sticky bool) TableAttribute                   // 设置粘性定位
	SetLayout(value string) TableAttribute                  // 设置布局
	SetStyle(style string) TableAttribute                   // 设置样式
	GetContent() template.HTML                              // 获取内容
}

// DataTableAttribute 数据表格属性接口
type DataTableAttribute interface {
	GetDataTableHeader() template.HTML                          // 获取数据表格头部
	SetThead(value Thead) DataTableAttribute                    // 设置表头
	SetInfoList(value []map[string]InfoItem) DataTableAttribute // 设置信息列表
	SetEditUrl(value string) DataTableAttribute                 // 设置编辑URL
	SetDeleteUrl(value string) DataTableAttribute               // 设置删除URL
	SetNewUrl(value string) DataTableAttribute                  // 设置新建URL
	SetPrimaryKey(value string) DataTableAttribute              // 设置主键
	SetStyle(style string) DataTableAttribute                   // 设置样式
	SetAction(action template.HTML) DataTableAttribute          // 设置操作
	SetIsTab(value bool) DataTableAttribute                     // 设置是否为标签页
	SetActionFold(fold bool) DataTableAttribute                 // 设置操作折叠
	SetHideThead() DataTableAttribute                           // 设置隐藏表头
	SetLayout(value string) DataTableAttribute                  // 设置布局
	SetButtons(btns template.HTML) DataTableAttribute           // 设置按钮
	SetSticky(sticky bool) DataTableAttribute                   // 设置粘性定位
	SetHideFilterArea(value bool) DataTableAttribute            // 设置隐藏过滤区域
	SetHideRowSelector(value bool) DataTableAttribute           // 设置隐藏行选择器
	SetActionJs(aj template.JS) DataTableAttribute              // 设置操作JavaScript
	SetNoAction() DataTableAttribute                            // 设置无操作
	SetInfoUrl(value string) DataTableAttribute                 // 设置信息URL
	SetDetailUrl(value string) DataTableAttribute               // 设置详情URL
	SetHasFilter(hasFilter bool) DataTableAttribute             // 设置是否有过滤器
	SetSortUrl(value string) DataTableAttribute                 // 设置排序URL
	SetExportUrl(value string) DataTableAttribute               // 设置导出URL
	SetUpdateUrl(value string) DataTableAttribute               // 设置更新URL
	GetContent() template.HTML                                  // 获取内容
}

// TreeAttribute 树属性接口
type TreeAttribute interface {
	SetTree(value []menu.Item) TreeAttribute // 设置树数据
	SetEditUrl(value string) TreeAttribute   // 设置编辑URL
	SetOrderUrl(value string) TreeAttribute  // 设置排序URL
	SetUrlPrefix(value string) TreeAttribute // 设置URL前缀
	SetDeleteUrl(value string) TreeAttribute // 设置删除URL
	GetContent() template.HTML               // 获取内容
	GetTreeHeader() template.HTML            // 获取树头部
}

// TreeViewAttribute 树视图属性接口
type TreeViewAttribute interface {
	SetTree(value TreeViewData) TreeViewAttribute // 设置树数据
	SetUrlPrefix(value string) TreeViewAttribute  // 设置URL前缀
	SetID(id string) TreeViewAttribute            // 设置ID
	GetContent() template.HTML                    // 获取内容
}

// PaginatorAttribute 分页器属性接口
type PaginatorAttribute interface {
	SetCurPageStartIndex(value string) PaginatorAttribute        // 设置当前页起始索引
	SetCurPageEndIndex(value string) PaginatorAttribute          // 设置当前页结束索引
	SetTotal(value string) PaginatorAttribute                    // 设置总数
	SetPreviousClass(value string) PaginatorAttribute            // 设置上一页类名
	SetPreviousUrl(value string) PaginatorAttribute              // 设置上一页URL
	SetPages(value []map[string]string) PaginatorAttribute       // 设置页码列表
	SetPageSizeList(value []string) PaginatorAttribute           // 设置每页条数列表
	SetNextClass(value string) PaginatorAttribute                // 设置下一页类名
	SetNextUrl(value string) PaginatorAttribute                  // 设置下一页URL
	SetOption(value map[string]template.HTML) PaginatorAttribute // 设置选项
	SetUrl(value string) PaginatorAttribute                      // 设置URL
	SetExtraInfo(value template.HTML) PaginatorAttribute         // 设置额外信息
	SetEntriesInfo(value template.HTML) PaginatorAttribute       // 设置条目信息
	GetContent() template.HTML                                   // 获取内容
}

// TabsAttribute 标签页属性接口
type TabsAttribute interface {
	SetData(value []map[string]template.HTML) TabsAttribute // 设置数据
	GetContent() template.HTML                              // 获取内容
}

// AlertAttribute 警告属性接口
type AlertAttribute interface {
	SetTheme(value string) AlertAttribute          // 设置主题
	SetTitle(value template.HTML) AlertAttribute   // 设置标题
	SetContent(value template.HTML) AlertAttribute // 设置内容
	Warning(msg string) template.HTML              // 显示警告消息
	GetContent() template.HTML                     // 获取内容
}

// LinkAttribute 链接属性接口
type LinkAttribute interface {
	OpenInNewTab() LinkAttribute                        // 在新标签页打开
	SetURL(value string) LinkAttribute                  // 设置URL
	SetAttributes(attr template.HTMLAttr) LinkAttribute // 设置属性
	SetClass(class template.HTML) LinkAttribute         // 设置类名
	NoPjax() LinkAttribute                              // 不使用pjax
	SetTabTitle(value template.HTML) LinkAttribute      // 设置标签页标题
	SetContent(value template.HTML) LinkAttribute       // 设置内容
	GetContent() template.HTML                          // 获取内容
}

// PopupAttribute 弹窗属性接口
type PopupAttribute interface {
	SetID(value string) PopupAttribute                // 设置ID
	SetTitle(value template.HTML) PopupAttribute      // 设置标题
	SetDraggable() PopupAttribute                     // 设置可拖动
	SetHideFooter() PopupAttribute                    // 设置隐藏页脚
	SetWidth(width string) PopupAttribute             // 设置宽度
	SetHeight(height string) PopupAttribute           // 设置高度
	SetFooter(value template.HTML) PopupAttribute     // 设置页脚
	SetFooterHTML(value template.HTML) PopupAttribute // 设置页脚HTML
	SetBody(value template.HTML) PopupAttribute       // 设置主体
	SetSize(value string) PopupAttribute              // 设置尺寸
	GetContent() template.HTML                        // 获取内容
}

// PanelInfo 面板信息结构体
type PanelInfo struct {
	Thead    Thead    `json:"thead"`     // 表头
	InfoList InfoList `json:"info_list"` // 信息列表
}

// Thead 表头类型
type Thead []TheadItem

// TheadItem 表头项结构体
type TheadItem struct {
	Head       string       `json:"head"`        // 表头文本
	Sortable   bool         `json:"sortable"`    // 是否可排序
	Field      string       `json:"field"`       // 字段名
	Hide       bool         `json:"hide"`        // 是否隐藏
	Editable   bool         `json:"editable"`    // 是否可编辑
	EditType   string       `json:"edit_type"`   // 编辑类型
	EditOption FieldOptions `json:"edit_option"` // 编辑选项
	Width      string       `json:"width"`       // 宽度
}

// GroupBy 按分组对表头进行分组
func (t Thead) GroupBy(group [][]string) []Thead {
	var res = make([]Thead, len(group))

	for key, value := range group {
		var newThead = make(Thead, 0)

		for _, info := range t {
			if modules.InArray(value, info.Field) {
				newThead = append(newThead, info)
			}
		}

		res[key] = newThead
	}

	return res
}

// TreeViewData 树视图数据结构体
type TreeViewData struct {
	Data                   TreeViewItems `json:"data,omitempty"`                   // 树视图项列表
	Levels                 int           `json:"levels,omitempty"`                 // 层级数
	BackColor              string        `json:"backColor,omitempty"`              // 背景色
	BorderColor            string        `json:"borderColor,omitempty"`            // 边框颜色
	CheckedIcon            string        `json:"checkedIcon,omitempty"`            // 选中图标
	CollapseIcon           string        `json:"collapseIcon,omitempty"`           // 折叠图标
	Color                  string        `json:"color,omitempty"`                  // 颜色
	EmptyIcon              string        `json:"emptyIcon,omitempty"`              // 空图标
	EnableLinks            bool          `json:"enableLinks,omitempty"`            // 启用链接
	ExpandIcon             string        `json:"expandIcon,omitempty"`             // 展开图标
	MultiSelect            bool          `json:"multiSelect,omitempty"`            // 多选
	NodeIcon               string        `json:"nodeIcon,omitempty"`               // 节点图标
	OnhoverColor           string        `json:"onhoverColor,omitempty"`           // 悬停颜色
	SelectedIcon           string        `json:"selectedIcon,omitempty"`           // 选中图标
	SearchResultColor      string        `json:"searchResultColor,omitempty"`      // 搜索结果颜色
	SelectedBackColor      string        `json:"selectedBackColor,omitempty"`      // 选中背景色
	SelectedColor          string        `json:"selectedColor,omitempty"`          // 选中颜色
	ShowBorder             bool          `json:"showBorder,omitempty"`             // 显示边框
	ShowCheckbox           bool          `json:"showCheckbox,omitempty"`           // 显示复选框
	ShowIcon               bool          `json:"showIcon,omitempty"`               // 显示图标
	ShowTags               bool          `json:"showTags,omitempty"`               // 显示标签
	UncheckedIcon          string        `json:"uncheckedIcon,omitempty"`          // 未选中图标
	SearchResultBackColor  string        `json:"searchResultBackColor,omitempty"`  // 搜索结果背景色
	HighlightSearchResults bool          `json:"highlightSearchResults,omitempty"` // 高亮搜索结果
}

// TreeViewItems 树视图项列表类型
type TreeViewItems []TreeViewItem

// TreeViewItemState 树视图项状态结构体
type TreeViewItemState struct {
	Checked  bool `json:"checked,omitempty"`  // 是否选中
	Disabled bool `json:"disabled,omitempty"` // 是否禁用
	Expanded bool `json:"expanded,omitempty"` // 是否展开
	Selected bool `json:"selected,omitempty"` // 是否被选择
}

// TreeViewItem 树视图项结构体
type TreeViewItem struct {
	Text         string            `json:"text,omitempty"`          // 文本
	Icon         string            `json:"icon,omitempty"`          // 图标
	SelectedIcon string            `json:"selected_icon,omitempty"` // 选中图标
	Color        string            `json:"color,omitempty"`         // 颜色
	BackColor    string            `json:"backColor,omitempty"`     // 背景色
	Href         string            `json:"href,omitempty"`          // 链接地址
	Selectable   bool              `json:"selectable,omitempty"`    // 是否可选择
	State        TreeViewItemState `json:"state,omitempty"`         // 状态
	Tags         []string          `json:"tags,omitempty"`          // 标签列表
	Nodes        TreeViewItems     `json:"nodes,omitempty"`         // 子节点列表
}

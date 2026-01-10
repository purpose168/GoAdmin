package form

import (
	"html/template"

	"github.com/purpose168/GoAdmin/modules/config"
	"github.com/purpose168/GoAdmin/modules/db"
	"github.com/purpose168/GoAdmin/modules/language"
)

// Type 表单字段类型枚举
// 定义了所有支持的表单字段类型
// 使用 uint8 类型，节省内存空间
type Type uint8

const (
	Default Type = iota // 默认类型
	Text               // 文本输入框
	SelectSingle       // 单选下拉框
	Select            // 多选下拉框
	IconPicker        // 图标选择器
	SelectBox         // 选择框（带搜索）
	File              // 文件上传
	Multifile         // 多文件上传
	Password          // 密码输入框
	RichText          // 富文本编辑器
	Datetime          // 日期时间选择器
	DatetimeRange      // 日期时间范围选择器
	Radio             // 单选按钮组
	Checkbox          // 复选框组
	CheckboxStacked    // 堆叠式复选框
	CheckboxSingle     // 单个复选框
	Email             // 邮箱输入框
	Date              // 日期选择器
	DateRange         // 日期范围选择器
	Url               // URL 输入框
	Ip                // IP 地址输入框
	Color             // 颜色选择器
	Array             // 数组输入框
	Currency          // 货币输入框
	Rate              // 评分组件
	Number            // 数字输入框
	Table             // 表格组件
	NumberRange       // 数字范围输入框
	TextArea          // 多行文本输入框
	Custom            // 自定义组件
	Switch            // 开关组件
	Code              // 代码编辑器
	Slider            // 滑块组件
)

// AllType 所有表单字段类型的列表
// 用于类型验证和遍历
var AllType = []Type{Default, Text, Array, SelectSingle, Select, IconPicker, SelectBox, File, Multifile, Password,
	RichText, Datetime, DatetimeRange, Checkbox, CheckboxStacked, Radio, Table, Email, Url, Ip, Color, Currency, Number, NumberRange,
	TextArea, Custom, Switch, Code, Rate, Slider, Date, DateRange, CheckboxSingle}

// CheckType 检查表单类型是否有效
// 如果类型有效，返回该类型；否则返回默认类型
//
// 参数：
//   - t: 要检查的表单类型
//   - def: 默认类型，当 t 无效时返回
//
// 返回值：
//   - Type: 有效的表单类型
//
// 使用示例：
//   // 检查类型是否有效，无效则返回 Text 类型
//   CheckType(unknownType, Text)
func CheckType(t, def Type) Type {
	for _, item := range AllType {
		if t == item {
			return t
		}
	}
	return def
}

// Layout 表单布局类型枚举
// 定义了表单字段的布局方式
type Layout uint8

const (
	LayoutDefault Layout = iota // 默认布局
	LayoutTwoCol              // 两列布局
	LayoutThreeCol            // 三列布局
	LayoutFourCol             // 四列布局
	LayoutFiveCol             // 五列布局
	LayoutSixCol              // 六列布局
	LayoutFlow                // 流式布局
	LayoutTab                 // 标签页布局
	LayoutFilter              // 过滤器布局
)

// Col 返回布局的列数
// 用于确定表单字段在网格系统中的列数
//
// 返回值：
//   - int: 列数（2、3、4、5、6），默认返回 0
func (l Layout) Col() int {
	if l == LayoutTwoCol {
		return 2
	}
	if l == LayoutThreeCol {
		return 3
	}
	if l == LayoutFourCol {
		return 4
	}
	if l == LayoutFiveCol {
		return 5
	}
	if l == LayoutSixCol {
		return 6
	}
	return 0
}

// Filter 判断是否为过滤器布局
// 过滤器布局用于搜索和筛选表单
//
// 返回值：
//   - bool: 如果是过滤器布局返回 true，否则返回 false
func (l Layout) Filter() bool {
	return l == LayoutFilter
}

// Flow 判断是否为流式布局
// 流式布局使表单字段自动换行
//
// 返回值：
//   - bool: 如果是流式布局返回 true，否则返回 false
func (l Layout) Flow() bool {
	return l == LayoutFlow
}

// Default 判断是否为默认布局
//
// 返回值：
//   - bool: 如果是默认布局返回 true，否则返回 false
func (l Layout) Default() bool {
	return l == LayoutDefault
}

// String 返回布局类型的字符串表示
// 用于日志记录和调试
//
// 返回值：
//   - string: 布局类型的名称
func (l Layout) String() string {
	switch l {
	case LayoutDefault:
		return "LayoutDefault"
	case LayoutTwoCol:
		return "LayoutTwoCol"
	case LayoutThreeCol:
		return "LayoutThreeCol"
	case LayoutFourCol:
		return "LayoutFourCol"
	case LayoutFiveCol:
		return "LayoutFiveCol"
	case LayoutSixCol:
		return "LayoutSixCol"
	case LayoutFlow:
		return "LayoutFlow"
	case LayoutTab:
		return "LayoutTab"
	default:
		return "LayoutDefault"
	}
}

// GetLayoutFromString 从字符串获取布局类型
// 用于从配置文件或数据库中读取布局配置
//
// 参数：
//   - s: 布局类型的字符串表示
//
// 返回值：
//   - Layout: 对应的布局类型，无效时返回默认布局
func GetLayoutFromString(s string) Layout {
	switch s {
	case "LayoutDefault":
		return LayoutDefault
	case "LayoutTwoCol":
		return LayoutTwoCol
	case "LayoutThreeCol":
		return LayoutThreeCol
	case "LayoutFourCol":
		return LayoutFourCol
	case "LayoutFiveCol":
		return LayoutFiveCol
	case "LayoutSixCol":
		return LayoutSixCol
	case "LayoutFlow":
		return LayoutFlow
	case "LayoutTab":
		return LayoutTab
	default:
		return LayoutDefault
	}
}

// Name 返回表单类型的名称
// 用于日志记录和调试
//
// 返回值：
//   - string: 表单类型的名称
func (t Type) Name() string {
	switch t {
	case Default:
		return "Default"
	case Text:
		return "Text"
	case SelectSingle:
		return "SelectSingle"
	case Select:
		return "Select"
	case IconPicker:
		return "IconPicker"
	case SelectBox:
		return "SelectBox"
	case File:
		return "File"
	case Table:
		return "Table"
	case Multifile:
		return "Multifile"
	case Password:
		return "Password"
	case RichText:
		return "RichText"
	case Rate:
		return "Rate"
	case Checkbox:
		return "Checkbox"
	case CheckboxStacked:
		return "CheckboxStacked"
	case CheckboxSingle:
		return "CheckboxSingle"
	case Date:
		return "Date"
	case DateRange:
		return "DateRange"
	case Datetime:
		return "Datetime"
	case DatetimeRange:
		return "DatetimeRange"
	case Radio:
		return "Radio"
	case Slider:
		return "Slider"
	case Array:
		return "Array"
	case Email:
		return "Email"
	case Url:
		return "Url"
	case Ip:
		return "Ip"
	case Color:
		return "Color"
	case Currency:
		return "Currency"
	case Number:
		return "Number"
	case NumberRange:
		return "NumberRange"
	case TextArea:
		return "TextArea"
	case Custom:
		return "Custom"
	case Switch:
		return "Switch"
	case Code:
		return "Code"
	default:
		panic("wrong form type")
	}
}

// String 返回表单类型的字符串表示
// 用于生成 HTML 和 JavaScript 代码
//
// 返回值：
//   - string: 表单类型的标识符（小写，使用下划线分隔）
func (t Type) String() string {
	switch t {
	case Default:
		return "default"
	case Text:
		return "text"
	case SelectSingle:
		return "select_single"
	case Select:
		return "select"
	case IconPicker:
		return "iconpicker"
	case SelectBox:
		return "selectbox"
	case File:
		return "file"
	case Table:
		return "table"
	case Multifile:
		return "multi_file"
	case Password:
		return "password"
	case RichText:
		return "richtext"
	case Rate:
		return "rate"
	case Checkbox:
		return "checkbox"
	case CheckboxStacked:
		return "checkbox_stacked"
	case CheckboxSingle:
		return "checkbox_single"
	case Date:
		return "datetime"
	case DateRange:
		return "datetime_range"
	case Datetime:
		return "datetime"
	case DatetimeRange:
		return "datetime_range"
	case Radio:
		return "radio"
	case Slider:
		return "slider"
	case Array:
		return "array"
	case Email:
		return "email"
	case Url:
		return "url"
	case Ip:
		return "ip"
	case Color:
		return "color"
	case Currency:
		return "currency"
	case Number:
		return "number"
	case NumberRange:
		return "number_range"
	case TextArea:
		return "textarea"
	case Custom:
		return "custom"
	case Switch:
		return "switch"
	case Code:
		return "code"
	default:
		panic("wrong form type")
	}
}

// IsSelect 判断是否为选择类字段
// 选择类字段包括下拉框、单选按钮、复选框等
//
// 返回值：
//   - bool: 如果是选择类字段返回 true，否则返回 false
func (t Type) IsSelect() bool {
	return t == Select || t == SelectSingle || t == SelectBox || t == Radio || t == Switch ||
		t == Checkbox || t == CheckboxStacked || t == CheckboxSingle
}

// IsArray 判断是否为数组字段
//
// 返回值：
//   - bool: 如果是数组字段返回 true，否则返回 false
func (t Type) IsArray() bool {
	return t == Array
}

// IsTable 判断是否为表格字段
//
// 返回值：
//   - bool: 如果是表格字段返回 true，否则返回 false
func (t Type) IsTable() bool {
	return t == Table
}

// IsSingleSelect 判断是否为单选字段
// 单选字段只能选择一个选项
//
// 返回值：
//   - bool: 如果是单选字段返回 true，否则返回 false
func (t Type) IsSingleSelect() bool {
	return t == SelectSingle || t == Radio || t == Switch || t == CheckboxSingle
}

// IsMultiSelect 判断是否为多选字段
// 多选字段可以选择多个选项
//
// 返回值：
//   - bool: 如果是多选字段返回 true，否则返回 false
func (t Type) IsMultiSelect() bool {
	return t == Select || t == SelectBox || t == Checkbox || t == CheckboxStacked
}

// IsMultiFile 判断是否为多文件上传字段
//
// 返回值：
//   - bool: 如果是多文件上传字段返回 true，否则返回 false
func (t Type) IsMultiFile() bool {
	return t == Multifile
}

// IsRange 判断是否为范围字段
// 范围字段包括日期时间范围和数字范围
//
// 返回值：
//   - bool: 如果是范围字段返回 true，否则返回 false
func (t Type) IsRange() bool {
	return t == DatetimeRange || t == NumberRange
}

// IsFile 判断是否为文件上传字段
// 包括单文件和多文件上传
//
// 返回值：
//   - bool: 如果是文件上传字段返回 true，否则返回 false
func (t Type) IsFile() bool {
	return t == File || t == Multifile
}

// IsSlider 判断是否为滑块字段
//
// 返回值：
//   - bool: 如果是滑块字段返回 true，否则返回 false
func (t Type) IsSlider() bool {
	return t == Slider
}

// IsDateTime 判断是否为日期时间字段
//
// 返回值：
//   - bool: 如果是日期时间字段返回 true，否则返回 false
func (t Type) IsDateTime() bool {
	return t == Datetime
}

// IsDateTimeRange 判断是否为日期时间范围字段
//
// 返回值：
//   - bool: 如果是日期时间范围字段返回 true，否则返回 false
func (t Type) IsDateTimeRange() bool {
	return t == DatetimeRange
}

// IsDate 判断是否为日期字段
//
// 返回值：
//   - bool: 如果是日期字段返回 true，否则返回 false
func (t Type) IsDate() bool {
	return t == Date
}

// IsDateRange 判断是否为日期范围字段
//
// 返回值：
//   - bool: 如果是日期范围字段返回 true，否则返回 false
func (t Type) IsDateRange() bool {
	return t == DateRange
}

// IsCode 判断是否为代码编辑器字段
//
// 返回值：
//   - bool: 如果是代码编辑器字段返回 true，否则返回 false
func (t Type) IsCode() bool {
	return t == Code
}

// IsRichText 判断是否为富文本编辑器字段
//
// 返回值：
//   - bool: 如果是富文本编辑器字段返回 true，否则返回 false
func (t Type) IsRichText() bool {
	return t == RichText
}

// IsTextarea 判断是否为多行文本输入框字段
//
// 返回值：
//   - bool: 如果是多行文本输入框字段返回 true，否则返回 false
func (t Type) IsTextarea() bool {
	return t == TextArea
}

// IsEditor 判断是否为编辑器字段
// 编辑器字段包括代码编辑器、富文本编辑器和多行文本输入框
//
// 返回值：
//   - bool: 如果是编辑器字段返回 true，否则返回 false
func (t Type) IsEditor() bool {
	return t == TextArea || t == Code || t == RichText
}

// IsCustom 判断是否为自定义字段
//
// 返回值：
//   - bool: 如果是自定义字段返回 true，否则返回 false
func (t Type) IsCustom() bool {
	return t == Custom
}

// FixOptions 修复表单类型的默认选项
// 为特定类型的表单字段设置默认配置
//
// 参数：
//   - m: 选项映射
//
// 返回值：
//   - map[string]interface{}: 修复后的选项映射
//
// 使用示例：
//   // 为滑块字段设置默认选项
//   Slider.FixOptions(options)
func (t Type) FixOptions(m map[string]interface{}) map[string]interface{} {
	if t == Slider {
		// 如果未设置 type，默认为 single
		if _, ok := m["type"]; !ok {
			m["type"] = "single"
		}
		// 如果未设置 prettify，默认为 false
		if _, ok := m["prettify"]; !ok {
			m["prettify"] = false
		}
		// 如果未设置 hasGrid，默认为 true
		if _, ok := m["hasGrid"]; !ok {
			m["hasGrid"] = true
		}
		return m
	}
	return m
}

// SelectedLabel 返回选中状态的 HTML 属性
// 用于生成表单字段的选中状态属性
//
// 参数：
//   - field: 字段名称
//
// 返回值：
//   - []template.HTML: 选中状态的 HTML 属性数组
//     - 第一个元素：选中属性（"selected" 或 "checked"）
//     - 第二个元素：空字符串
func (t Type) SelectedLabel() []template.HTML {
	if t == Select || t == SelectSingle || t == SelectBox {
		return []template.HTML{"selected", ""}
	}
	if t == Radio || t == Switch || t == Checkbox || t == CheckboxStacked || t == CheckboxSingle {
		return []template.HTML{"checked", ""}
	}
	return []template.HTML{"", ""}
}

// GetDefaultOptions 获取表单类型的默认选项
// 为不同类型的表单字段生成默认配置选项
//
// 参数：
//   - field: 字段名称
//
// 返回值：
//   - map[string]interface{}: 主选项映射
//   - map[string]interface{}: 次选项映射（用于范围字段）
//   - template.JS: JavaScript 配置代码
func (t Type) GetDefaultOptions(field string) (map[string]interface{}, map[string]interface{}, template.JS) {
	switch t {
	case File, Multifile:
		// 文件上传字段默认选项
		return map[string]interface{}{
			"overwriteInitial":     true,                        // 覆盖初始预览
			"initialPreviewAsData": true,                        // 初始预览作为数据
			"browseLabel":          language.Get("Browse"),       // 浏览按钮标签
			"showRemove":           false,                       // 不显示删除按钮
			"previewClass":         "preview-" + field,          // 预览类名
			"showUpload":           false,                       // 不显示上传按钮
			"allowedFileTypes":     []string{"image"},          // 允许的文件类型
		}, nil, ""
	case Slider:
		// 滑块字段默认选项
		return map[string]interface{}{
			"type":     "single", // 单个滑块
			"prettify": false,   // 不美化数字
			"hasGrid":  true,    // 显示网格
			"max":      100,     // 最大值
			"min":      1,       // 最小值
			"step":     1,       // 步长
			"postfix":  "",      // 后缀
		}, nil, ""
	case DatetimeRange:
		// 日期时间范围字段默认选项
		op1, op2 := getDateTimeRangeOptions(DatetimeRange)
		return op1, op2, ""
	case Datetime:
		// 日期时间字段默认选项
		return getDateTimeOptions(Datetime), nil, ""
	case Date:
		// 日期字段默认选项
		return getDateTimeOptions(Date), nil, ""
	case DateRange:
		// 日期范围字段默认选项
		op1, op2 := getDateTimeRangeOptions(DateRange)
		return op1, op2, ""
	case Code:
		// 代码编辑器字段默认选项
		return nil, nil, `
	theme = "monokai";
	font_size = 14;
	language = "html";
	options = {useWorker: false};
`
	}

	return nil, nil, ""
}

// getDateTimeOptions 获取日期时间字段的默认选项
//
// 参数：
//   - f: 表单类型（Date 或 Datetime）
//
// 返回值：
//   - map[string]interface{}: 日期时间选项映射
func getDateTimeOptions(f Type) map[string]interface{} {
	// 默认格式：年-月-日 时:分:秒
	format := "YYYY-MM-DD HH:mm:ss"
	if f == Date {
		// 日期格式：年-月-日
		format = "YYYY-MM-DD"
	}
	m := map[string]interface{}{
		"format":           format,           // 日期时间格式
		"locale":           "en",            // 语言环境
		"allowInputToggle": true,            // 允许输入切换
	}
	// 如果配置语言为中文，使用中文语言环境
	if config.GetLanguage() == language.CN || config.GetLanguage() == "cn" {
		m["locale"] = "zh-CN"
	}
	return m
}

// getDateTimeRangeOptions 获取日期时间范围字段的默认选项
//
// 参数：
//   - f: 表单类型（DateRange 或 DatetimeRange）
//
// 返回值：
//   - map[string]interface{}: 起始日期时间选项映射
//   - map[string]interface{}: 结束日期时间选项映射
func getDateTimeRangeOptions(f Type) (map[string]interface{}, map[string]interface{}) {
	// 默认格式：年-月-日 时:分:秒
	format := "YYYY-MM-DD HH:mm:ss"
	if f == DateRange {
		// 日期格式：年-月-日
		format = "YYYY-MM-DD"
	}
	// 起始日期时间选项
	m := map[string]interface{}{
		"format": format,  // 日期时间格式
		"locale": "en",   // 语言环境
	}
	// 结束日期时间选项
	m1 := map[string]interface{}{
		"format":     format,  // 日期时间格式
		"locale":     "en",    // 语言环境
		"useCurrent": false,    // 不使用当前时间
	}
	// 如果配置语言为中文，使用中文语言环境
	if config.GetLanguage() == language.CN || config.GetLanguage() == "cn" {
		m["locale"] = "zh-CN"
		m1["locale"] = "zh-CN"
	}
	return m, m1
}

// GetFormTypeFromFieldType 根据数据库字段类型获取表单类型
// 用于自动推断表单字段的类型
//
// 参数：
//   - typeName: 数据库字段类型
//   - fieldName: 字段名称
//
// 返回值：
//   - string: 表单类型名称
//
// 使用示例：
//   // 根据数据库字段类型自动推断表单类型
//   GetFormTypeFromFieldType(db.Int, "age")
func GetFormTypeFromFieldType(typeName db.DatabaseType, fieldName string) string {

	// 特殊字段名处理
	if fieldName == "password" {
		return "Password"
	}

	if fieldName == "id" {
		return "Default"
	}

	if fieldName == "ip" {
		return "Ip"
	}

	if fieldName == "Url" {
		return "Url"
	}

	if fieldName == "email" {
		return "Email"
	}

	if fieldName == "color" {
		return "Color"
	}

	if fieldName == "money" {
		return "Currency"
	}

	// 根据数据库类型推断表单类型
	switch typeName {
	case db.Int, db.Tinyint, db.Int4, db.Integer, db.Mediumint, db.Smallint,
		db.Numeric, db.Smallserial, db.Serial, db.Bigserial, db.Money, db.Bigint:
		// 整数类型：数字输入框
		return "Number"
	case db.Text, db.Longtext, db.Mediumtext, db.Tinytext:
		// 文本类型：富文本编辑器
		return "RichText"
	case db.Datetime, db.Date, db.Time, db.Timestamp, db.Timestamptz, db.Year:
		// 日期时间类型：日期时间选择器
		return "Datetime"
	}

	// 默认：文本输入框
	return "Text"
}

// DefaultHTML 生成默认的 HTML 包装器
// 用于包装表单字段的 HTML 内容
//
// 参数：
//   - value: 要包装的 HTML 内容
//
// 返回值：
//   - template.HTML: 包装后的 HTML
func DefaultHTML(value string) template.HTML {
	return template.HTML(`<div class="box box-solid box-default no-margin"><div class="box-body" style="min-height: 40px;">` +
		value + `</div></div>`)
}

// HiddenInputHTML 生成隐藏输入框的 HTML
// 用于在表单中存储隐藏数据
//
// 参数：
//   - field: 字段名称
//   - value: 字段值
//
// 返回值：
//   - template.HTML: 隐藏输入框的 HTML
func HiddenInputHTML(field, value string) template.HTML {
	return template.HTML(`<input type="hidden" name="` + field + `" value="` + value + `" class="form-control">`)
}

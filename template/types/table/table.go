package table

// Type 表格列类型枚举
// 定义了表格中不同列的数据类型
// 使用 uint8 类型，节省内存空间
type Type uint8

const (
	Text     Type = iota // 文本类型，用于显示普通文本内容
	Textarea             // 多行文本类型，用于显示较长的文本内容
	Select               // 选择类型，用于显示下拉选择框
	Date                 // 日期类型，用于显示日期
	Datetime             // 日期时间类型，用于显示日期和时间
	Year                 // 年份类型，用于显示年份选择器
	Month                // 月份类型，用于显示月份选择器
	Day                  // 日期类型，用于显示日期选择器
	Switch               // 开关类型，用于显示切换开关
)

// String 返回表格列类型的字符串表示
// 用于生成 HTML 和 JavaScript 代码
//
// 返回值：
//   - string: 表格列类型的标识符（小写）
//
// 使用示例：
//
//	// 获取类型的字符串表示
//	Text.String() // 返回 "text"
//	Select.String() // 返回 "select"
func (t Type) String() string {
	switch t {
	case Text:
		return "text"
	case Select:
		return "select"
	case Textarea:
		return "textarea"
	case Date:
		return "date"
	case Switch:
		return "switch"
	case Year:
		return "year"
	case Month:
		return "month"
	case Day:
		return "day"
	case Datetime:
		return "datetime"
	default:
		panic("wrong form type")
	}
}

// IsSwitch 判断是否为开关类型
// 用于确定列是否应该显示为切换开关
//
// 返回值：
//   - bool: 如果是开关类型返回 true，否则返回 false
//
// 使用示例：
//
//	// 检查列类型是否为开关
//	if columnType.IsSwitch() {
//	    // 处理开关类型的逻辑
//	}
func (t Type) IsSwitch() bool {
	return t == Switch
}

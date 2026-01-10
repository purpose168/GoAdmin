package selection

import "fmt"

// Data 选择框数据结构
// 用于存储选择框的选项数据和分页信息
// 通常用于 AJAX 异步加载选择框选项时的响应数据
type Data struct {
	Results    Options    `json:"results"`    // 选择框选项列表
	Pagination Pagination `json:"pagination"` // 分页信息
}

// Pagination 分页信息结构
// 用于指示是否还有更多数据可以加载
type Pagination struct {
	More bool `json:"more"` // 是否还有更多数据，true 表示可以继续加载
}

// Options 选项列表类型
// 是 Option 结构体的切片，用于存储多个选择框选项
type Options []Option

// Option 选择框选项结构
// 定义了单个选择框选项的属性
type Option struct {
	ID       interface{} `json:"id"`                 // 选项的唯一标识符，可以是任意类型
	Text     string      `json:"text"`               // 选项显示的文本内容
	Selected bool        `json:"selected,omitempty"` // 选项是否被选中，omitempty 表示为空时省略
	Disabled bool        `json:"disabled,omitempty"` // 选项是否禁用，omitempty 表示为空时省略
}

// Configuration 选择框配置结构
// 用于配置选择框的各种行为和样式
// TODO: 需要进一步明确各个配置项的具体用途
type Configuration struct {
	AdaptContainerCssClass string                 `json:"adaptContainerCssClass,omitempty"` // 自适应容器 CSS 类名
	AdaptDropdownCssClass  string                 `json:"adaptDropdownCssClass,omitempty"`  // 自适应下拉框 CSS 类名
	Ajax                   map[string]interface{} `json:"ajax,omitempty"`                   // AJAX 配置选项
	AllowClear             bool                   `json:"allowClear,omitempty"`             // 是否允许清除选择
	AmdBase                string                 `json:"amdBase,omitempty"`                // AMD 模块基础路径
	AmdLanguageBase        string                 `json:"amdLanguageBase,omitempty"`        // AMD 语言包基础路径
	CloseOnSelect          bool                   `json:"closeOnSelect,omitempty"`          // 选择后是否关闭下拉框
	ContainerCss           map[string]interface{} `json:"containerCss,omitempty"`           // 容器 CSS 样式
	ContainerCssClass      string                 `json:"containerCssClass,omitempty"`      // 容器 CSS 类名
	Data                   Options                `json:"data,omitempty"`                   // 静态数据选项
	Debug                  bool                   `json:"debug,omitempty"`                  // 是否启用调试模式
	Disabled               bool                   `json:"disabled,omitempty"`               // 是否禁用选择框

	DropdownAutoWidth bool                   `json:"dropdownAutoWidth,omitempty"` // 下拉框宽度是否自适应
	DropdownCss       map[string]interface{} `json:"dropdownCss,omitempty"`       // 下拉框 CSS 样式
	DropdownCssClass  string                 `json:"dropdownCssClass,omitempty"`  // 下拉框 CSS 类名
	DropdownParent    string                 `json:"dropdownParent,omitempty"`    // 下拉框的父元素选择器

	EscapeMarkup  func()      `json:"escapeMarkup,omitempty"`  // HTML 转义函数
	InitSelection func()      `json:"initSelection,omitempty"` // 初始化选择函数
	Language      interface{} `json:"language,omitempty"`      // 语言设置
	Matcher       func()      `json:"matcher,omitempty"`       // 选项匹配函数

	MaximumInputLength      int `json:"maximumInputLength,omitempty"`      // 最大输入长度
	MaximumSelectionLength  int `json:"maximumSelectionLength,omitempty"`  // 最大选择数量
	MinimumInputLength      int `json:"minimumInputLength,omitempty"`      // 最小输入长度（触发搜索）
	MinimumResultsForSearch int `json:"minimumResultsForSearch,omitempty"` // 显示搜索框的最小结果数

	Multiple      bool        `json:"multiple,omitempty"`      // 是否允许多选
	Placeholder   interface{} `json:"placeholder,omitempty"`   // 占位符文本
	Query         func()      `json:"query,omitempty"`         // 查询函数（用于 AJAX）
	SelectOnClose bool        `json:"selectOnClose,omitempty"` // 关闭时是否选择高亮项
	Sorter        func()      `json:"sorter,omitempty"`        // 排序函数
	Tags          bool        `json:"tags,omitempty"`          // 是否允许创建新标签

	TemplateResultFns    []Function // 结果模板函数列表（内部使用）
	TemplateResult       string     `json:"templateResult,omitempty"` // 结果模板字符串
	TemplateSelectionFns []Function // 选择模板函数列表（内部使用）
	TemplateSelection    string     `json:"templateSelection,omitempty"` // 选择模板字符串

	Theme             string   `json:"theme,omitempty"`             // 主题名称
	Tokenizer         func()   `json:"tokenizer,omitempty"`         // 分词函数
	TokenSeparators   []string `json:"tokenSeparators,omitempty"`   // 分词分隔符
	Width             string   `json:"width,omitempty"`             // 选择框宽度
	ScrollAfterSelect bool     `json:"scrollAfterSelect,omitempty"` // 选择后是否滚动到视图
}

// Function 函数结构
// 用于构建 JavaScript 函数调用链
// 支持条件判断、返回值、变量操作等操作
type Function struct {
	Format string                                            // 函数格式字符串，包含占位符 %s
	Args   []Arg                                             // 函数参数列表
	Next   *Function                                         // 下一个要执行的函数
	P      func(f string, args []Arg, next *Function) string // 函数处理器，用于生成最终的 JavaScript 代码
}

// ArgType 参数类型枚举
// 定义了函数参数的三种类型
type ArgType int

const (
	ArgInt       ArgType = iota // 整数类型参数
	ArgString                   // 字符串类型参数
	ArgOperation                // 操作类型参数（如比较运算符）
)

// Arg 参数接口
// 定义了参数的基本行为
type Arg interface {
	Type() ArgType      // 获取参数类型
	String() string     // 获取参数的字符串表示
	Wrap(string) string // 将参数包装为字符串格式
}

// BaseArg 基础参数类型
// 实现了 Arg 接口的基础功能
type BaseArg string

// String 返回参数的字符串表示
func (b BaseArg) String() string {
	return string(b)
}

// Wrap 将参数包装为字符串格式
// 默认不进行包装，直接返回原始字符串
func (b BaseArg) Wrap(s string) string {
	return s
}

// StringArg 字符串参数类型
// 继承自 BaseArg，表示字符串类型的参数
type StringArg BaseArg

// Type 返回参数类型为字符串类型
func (s StringArg) Type() ArgType {
	return ArgString
}

// Wrap 将字符串参数包装为双引号格式
// 用于在 JavaScript 代码中表示字符串字面量
func (s StringArg) Wrap(ss string) string {
	return `"` + ss + `"`
}

// IntArg 整数参数类型
// 继承自 BaseArg，表示整数类型的参数
type IntArg BaseArg

// Type 返回参数类型为整数类型
func (s IntArg) Type() ArgType {
	return ArgInt
}

// OperationArg 操作参数类型
// 继承自 BaseArg，表示操作符类型的参数（如 ==、!=、>、< 等）
type OperationArg BaseArg

// Type 返回参数类型为操作类型
func (s OperationArg) Type() ArgType {
	return ArgOperation
}

// If 创建条件判断函数
// 生成 JavaScript 的 if 条件语句
//
// 参数：
//   - operation: 操作符参数（如 ==、!=、>、< 等）
//   - arg: 比较的值参数
//   - next: 条件为真时执行的下一个函数
//
// 返回值：
//   - Function: 表示条件判断的函数对象
//
// 使用示例：
//
//	// 生成 if (x == 1) { return x; }
//	If(OperationArg("=="), IntArg("1"), Return())
func If(operation, arg Arg, next *Function) Function {
	return Function{
		// 格式化字符串：if (%s %s %s) { %s }
		Format: `if (%s ` + operation.Wrap("%s") + " " + arg.Wrap("%s") + `) {
	%s
}
`,
		Next: next,
		Args: []Arg{operation, arg},
		// 函数处理器：递归生成完整的 JavaScript 代码
		P: func(f string, args []Arg, next *Function) string {
			return fmt.Sprintf(f, args[0], args[1], args[2],
				next.P(next.Format, append([]Arg{args[0]}, next.Args...), next.Next))
		},
	}
}

// Return 创建返回函数
// 生成 JavaScript 的 return 语句
//
// 返回值：
//   - Function: 表示返回语句的函数对象
//
// 使用示例：
//
//	// 生成 return x;
//	Return()
func Return() Function {
	return Function{
		Format: `return %s`,
		// 函数处理器：生成 return 语句
		P: func(f string, args []Arg, next *Function) string {
			return fmt.Sprintf(f, args[0])
		},
	}
}

// Add 创建追加赋值函数
// 生成 JavaScript 的 += 追加赋值语句
//
// 参数：
//   - arg: 要追加的值参数
//
// 返回值：
//   - Function: 表示追加赋值的函数对象
//
// 使用示例：
//
//	// 生成 x += 1;
//	Add(IntArg("1"))
func Add(arg Arg) Function {
	return Function{
		// 格式化字符串：%s += %s
		Format: `%s += ` + arg.Wrap("%s"),
		Args:   []Arg{arg},
		// 函数处理器：生成 += 语句
		P: func(f string, args []Arg, next *Function) string {
			return fmt.Sprintf(f, args[0], args[1])
		},
	}
}

// AddFront 创建前置追加赋值函数
// 生成 JavaScript 的前置追加赋值语句（将值添加到字符串前面）
//
// 参数：
//   - arg: 要前置追加的值参数
//
// 返回值：
//   - Function: 表示前置追加赋值的函数对象
//
// 使用示例：
//
//	// 生成 x = "prefix" + x;
//	AddFront(StringArg("prefix"))
func AddFront(arg Arg) Function {
	return Function{
		// 格式化字符串：%s = %s + %s
		Format: `%s = ` + arg.Wrap("%s") + ` + %s`,
		Args:   []Arg{arg},
		// 函数处理器：生成前置追加语句
		P: func(f string, args []Arg, next *Function) string {
			return fmt.Sprintf(f, args[0], args[1], args[0])
		},
	}
}

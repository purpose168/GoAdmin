package types

import "html/template"

// FilterOperator 是筛选操作符类型
type FilterOperator string

const (
	FilterOperatorLike           FilterOperator = "like" // 模糊匹配操作符
	FilterOperatorGreater        FilterOperator = ">"    // 大于操作符
	FilterOperatorGreaterOrEqual FilterOperator = ">="   // 大于等于操作符
	FilterOperatorEqual          FilterOperator = "="    // 等于操作符
	FilterOperatorNotEqual       FilterOperator = "!="   // 不等于操作符
	FilterOperatorLess           FilterOperator = "<"    // 小于操作符
	FilterOperatorLessOrEqual    FilterOperator = "<="   // 小于等于操作符
	FilterOperatorFree           FilterOperator = "free" // 自由操作符
)

// GetOperatorFromValue 根据值获取对应的筛选操作符
// 参数:
//   - value: 操作符值字符串
//
// 返回: 对应的筛选操作符
func GetOperatorFromValue(value string) FilterOperator {
	switch value {
	case "like":
		return FilterOperatorLike
	case "gr":
		return FilterOperatorGreater
	case "gq":
		return FilterOperatorGreaterOrEqual
	case "eq":
		return FilterOperatorEqual
	case "ne":
		return FilterOperatorNotEqual
	case "le":
		return FilterOperatorLess
	case "lq":
		return FilterOperatorLessOrEqual
	case "free":
		return FilterOperatorFree
	default:
		return FilterOperatorEqual
	}
}

// Value 返回操作符的值字符串
// 返回: 操作符的值
func (o FilterOperator) Value() string {
	switch o {
	case FilterOperatorLike:
		return "like"
	case FilterOperatorGreater:
		return "gr"
	case FilterOperatorGreaterOrEqual:
		return "gq"
	case FilterOperatorEqual:
		return "eq"
	case FilterOperatorNotEqual:
		return "ne"
	case FilterOperatorLess:
		return "le"
	case FilterOperatorLessOrEqual:
		return "lq"
	case FilterOperatorFree:
		return "free"
	default:
		return "eq"
	}
}

// String 返回操作符的字符串表示
// 返回: 操作符字符串
func (o FilterOperator) String() string {
	return string(o)
}

// Label 返回操作符的标签HTML
// 对于like操作符返回空字符串，其他操作符返回其自身
// 返回: 操作符标签HTML
func (o FilterOperator) Label() template.HTML {
	if o == FilterOperatorLike {
		return ""
	}
	return template.HTML(o)
}

// AddOrNot 判断是否需要添加操作符字段
// 对于空字符串、free和like操作符返回false，其他返回true
// 返回: 是否需要添加操作符字段
func (o FilterOperator) AddOrNot() bool {
	return string(o) != "" && o != FilterOperatorFree && o != FilterOperatorLike
}

// Valid 判断操作符是否有效
// 返回: 操作符是否有效
func (o FilterOperator) Valid() bool {
	switch o {
	case FilterOperatorLike, FilterOperatorGreater, FilterOperatorGreaterOrEqual,
		FilterOperatorLess, FilterOperatorLessOrEqual, FilterOperatorFree:
		return true
	default:
		return false
	}
}

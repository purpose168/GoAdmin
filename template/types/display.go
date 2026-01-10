package types

import (
	"fmt"
	"html"
	"html/template"
	"strings"

	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/modules/config"
	"github.com/purpose168/GoAdmin/template/types/form"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// DisplayFnGenerator 显示函数生成器接口
type DisplayFnGenerator interface {
	Get(ctx *context.Context, args ...interface{}) FieldFilterFn // 获取显示函数
	JS() template.HTML                                           // 返回JavaScript代码
	HTML() template.HTML                                         // 返回HTML代码
}

// BaseDisplayFnGenerator 基础显示函数生成器结构体
type BaseDisplayFnGenerator struct{}

func (base *BaseDisplayFnGenerator) JS() template.HTML   { return "" } // 返回空JavaScript代码
func (base *BaseDisplayFnGenerator) HTML() template.HTML { return "" } // 返回空HTML代码

// displayFnGens 显示函数生成器映射
var displayFnGens = make(map[string]DisplayFnGenerator)

// RegisterDisplayFnGenerator 注册显示函数生成器
func RegisterDisplayFnGenerator(key string, gen DisplayFnGenerator) {
	if _, ok := displayFnGens[key]; ok {
		panic("display function generator has been registered")
	}
	displayFnGens[key] = gen
}

// FieldDisplay 字段显示结构体
type FieldDisplay struct {
	Display              FieldFilterFn          // 显示函数
	DisplayProcessChains DisplayProcessFnChains // 显示处理函数链
}

// ToDisplay 将字段值转换为显示格式
func (f FieldDisplay) ToDisplay(value FieldModel) interface{} {
	val := f.Display(value)

	// 如果存在显示处理链且值不是选择结果
	if len(f.DisplayProcessChains) > 0 && f.IsNotSelectRes(val) {
		valStr := fmt.Sprintf("%v", val)
		for _, process := range f.DisplayProcessChains {
			valStr = fmt.Sprintf("%v", process(FieldModel{
				Row:   value.Row,
				Value: valStr,
				ID:    value.ID,
			}))
		}
		return valStr
	}

	return val
}

// IsNotSelectRes 检查值是否不是选择结果
func (f FieldDisplay) IsNotSelectRes(v interface{}) bool {
	switch v.(type) {
	case template.HTML:
		return false
	case []string:
		return false
	case [][]string:
		return false
	default:
		return true
	}
}

// ToDisplayHTML 将字段值转换为HTML格式
func (f FieldDisplay) ToDisplayHTML(value FieldModel) template.HTML {
	v := f.ToDisplay(value)
	if h, ok := v.(template.HTML); ok {
		return h
	} else if s, ok := v.(string); ok {
		return template.HTML(s)
	} else if arr, ok := v.([]string); ok && len(arr) > 0 {
		return template.HTML(arr[0])
	} else if arr, ok := v.([]template.HTML); ok && len(arr) > 0 {
		return arr[0]
	} else if v != nil {
		return template.HTML(fmt.Sprintf("%v", v))
	} else {
		return ""
	}
}

// ToDisplayString 将字段值转换为字符串格式
func (f FieldDisplay) ToDisplayString(value FieldModel) string {
	v := f.ToDisplay(value)
	if h, ok := v.(template.HTML); ok {
		return string(h)
	} else if s, ok := v.(string); ok {
		return s
	} else if arr, ok := v.([]string); ok && len(arr) > 0 {
		return arr[0]
	} else if arr, ok := v.([]template.HTML); ok && len(arr) > 0 {
		return string(arr[0])
	} else if v != nil {
		return fmt.Sprintf("%v", v)
	} else {
		return ""
	}
}

// ToDisplayStringArray 将字段值转换为字符串数组格式
func (f FieldDisplay) ToDisplayStringArray(value FieldModel) []string {
	v := f.ToDisplay(value)
	if h, ok := v.(template.HTML); ok {
		return []string{string(h)}
	} else if s, ok := v.(string); ok {
		return []string{s}
	} else if arr, ok := v.([]string); ok && len(arr) > 0 {
		return arr
	} else if arr, ok := v.([]template.HTML); ok && len(arr) > 0 {
		ss := make([]string, len(arr))
		for k, a := range arr {
			ss[k] = string(a)
		}
		return ss
	} else if v != nil {
		return []string{fmt.Sprintf("%v", v)}
	} else {
		return []string{}
	}
}

// ToDisplayStringArrayArray 将字段值转换为字符串二维数组格式
func (f FieldDisplay) ToDisplayStringArrayArray(value FieldModel) [][]string {
	v := f.ToDisplay(value)
	if h, ok := v.(template.HTML); ok {
		return [][]string{{string(h)}}
	} else if s, ok := v.(string); ok {
		return [][]string{{s}}
	} else if arr, ok := v.([]string); ok && len(arr) > 0 {
		return [][]string{arr}
	} else if arr, ok := v.([][]string); ok && len(arr) > 0 {
		return arr
	} else if arr, ok := v.([]template.HTML); ok && len(arr) > 0 {
		ss := make([]string, len(arr))
		for k, a := range arr {
			ss[k] = string(a)
		}
		return [][]string{ss}
	} else if v != nil {
		return [][]string{{fmt.Sprintf("%v", v)}}
	} else {
		return [][]string{}
	}
}

// AddLimit 添加长度限制处理函数
func (f FieldDisplay) AddLimit(limit int) DisplayProcessFnChains {
	return f.DisplayProcessChains.Add(func(value FieldModel) interface{} {
		if limit > len(value.Value) {
			return value.Value
		} else if limit < 0 {
			return ""
		} else {
			return value.Value[:limit]
		}
	})
}

// AddTrimSpace 添加去除空格处理函数
func (f FieldDisplay) AddTrimSpace() DisplayProcessFnChains {
	return f.DisplayProcessChains.Add(func(value FieldModel) interface{} {
		return strings.TrimSpace(value.Value)
	})
}

// AddSubstr 添加子字符串处理函数
func (f FieldDisplay) AddSubstr(start int, end int) DisplayProcessFnChains {
	return f.DisplayProcessChains.Add(func(value FieldModel) interface{} {
		if start > end || start > len(value.Value) || end < 0 {
			return ""
		}
		if start < 0 {
			start = 0
		}
		if end > len(value.Value) {
			end = len(value.Value)
		}
		return value.Value[start:end]
	})
}

// AddToTitle 添加标题格式处理函数
func (f FieldDisplay) AddToTitle() DisplayProcessFnChains {
	return f.DisplayProcessChains.Add(func(value FieldModel) interface{} {
		return cases.Title(language.Und).String(value.Value)
	})
}

// AddToUpper 添加大写转换处理函数
func (f FieldDisplay) AddToUpper() DisplayProcessFnChains {
	return f.DisplayProcessChains.Add(func(value FieldModel) interface{} {
		return strings.ToUpper(value.Value)
	})
}

// AddToLower 添加小写转换处理函数
func (f FieldDisplay) AddToLower() DisplayProcessFnChains {
	return f.DisplayProcessChains.Add(func(value FieldModel) interface{} {
		return strings.ToLower(value.Value)
	})
}

// DisplayProcessFnChains 显示处理函数链类型
type DisplayProcessFnChains []FieldFilterFn

// Valid 检查处理链是否有效
func (d DisplayProcessFnChains) Valid() bool {
	return len(d) > 0
}

// Add 添加处理函数到链中
func (d DisplayProcessFnChains) Add(f FieldFilterFn) DisplayProcessFnChains {
	return append(d, f)
}

// Append 追加处理函数链
func (d DisplayProcessFnChains) Append(f DisplayProcessFnChains) DisplayProcessFnChains {
	return append(d, f...)
}

// Copy 复制处理函数链
func (d DisplayProcessFnChains) Copy() DisplayProcessFnChains {
	if len(d) == 0 {
		return make(DisplayProcessFnChains, 0)
	} else {
		var newDisplayProcessFnChains = make(DisplayProcessFnChains, len(d))
		copy(newDisplayProcessFnChains, d)
		return newDisplayProcessFnChains
	}
}

// chooseDisplayProcessChains 选择显示处理函数链
func chooseDisplayProcessChains(internal DisplayProcessFnChains) DisplayProcessFnChains {
	if len(internal) > 0 {
		return internal
	}
	return globalDisplayProcessChains.Copy()
}

// globalDisplayProcessChains 全局显示处理函数链
var globalDisplayProcessChains = make(DisplayProcessFnChains, 0)

// AddGlobalDisplayProcessFn 添加全局显示处理函数
func AddGlobalDisplayProcessFn(f FieldFilterFn) {
	globalDisplayProcessChains = globalDisplayProcessChains.Add(f)
}

// AddLimit 添加长度限制处理函数到全局链
func AddLimit(limit int) DisplayProcessFnChains {
	return addLimit(limit, globalDisplayProcessChains)
}

// AddTrimSpace 添加去除空格处理函数到全局链
func AddTrimSpace() DisplayProcessFnChains {
	return addTrimSpace(globalDisplayProcessChains)
}

// AddSubstr 添加子字符串处理函数到全局链
func AddSubstr(start int, end int) DisplayProcessFnChains {
	return addSubstr(start, end, globalDisplayProcessChains)
}

// AddToTitle 添加标题格式处理函数到全局链
func AddToTitle() DisplayProcessFnChains {
	return addToTitle(globalDisplayProcessChains)
}

// AddToUpper 添加大写转换处理函数到全局链
func AddToUpper() DisplayProcessFnChains {
	return addToUpper(globalDisplayProcessChains)
}

// AddToLower 添加小写转换处理函数到全局链
func AddToLower() DisplayProcessFnChains {
	return addToLower(globalDisplayProcessChains)
}

// AddXssFilter 添加XSS过滤处理函数到全局链
func AddXssFilter() DisplayProcessFnChains {
	return addXssFilter(globalDisplayProcessChains)
}

// AddXssJsFilter 添加XSS JavaScript过滤处理函数到全局链
func AddXssJsFilter() DisplayProcessFnChains {
	return addXssJsFilter(globalDisplayProcessChains)
}

// addLimit 添加长度限制处理函数
func addLimit(limit int, chains DisplayProcessFnChains) DisplayProcessFnChains {
	chains = chains.Add(func(value FieldModel) interface{} {
		if limit > len(value.Value) {
			return value
		} else if limit < 0 {
			return ""
		} else {
			return value.Value[:limit]
		}
	})
	return chains
}

// addTrimSpace 添加去除空格处理函数
func addTrimSpace(chains DisplayProcessFnChains) DisplayProcessFnChains {
	chains = chains.Add(func(value FieldModel) interface{} {
		return strings.TrimSpace(value.Value)
	})
	return chains
}

// addSubstr 添加子字符串处理函数
func addSubstr(start int, end int, chains DisplayProcessFnChains) DisplayProcessFnChains {
	chains = chains.Add(func(value FieldModel) interface{} {
		if start > end || start > len(value.Value) || end < 0 {
			return ""
		}
		if start < 0 {
			start = 0
		}
		if end > len(value.Value) {
			end = len(value.Value)
		}
		return value.Value[start:end]
	})
	return chains
}

// addToTitle 添加标题格式处理函数
func addToTitle(chains DisplayProcessFnChains) DisplayProcessFnChains {
	chains = chains.Add(func(value FieldModel) interface{} {
		return cases.Title(language.Und).String(value.Value)
	})
	return chains
}

// addToUpper 添加大写转换处理函数
func addToUpper(chains DisplayProcessFnChains) DisplayProcessFnChains {
	chains = chains.Add(func(value FieldModel) interface{} {
		return strings.ToUpper(value.Value)
	})
	return chains
}

// addToLower 添加小写转换处理函数
func addToLower(chains DisplayProcessFnChains) DisplayProcessFnChains {
	chains = chains.Add(func(value FieldModel) interface{} {
		return strings.ToLower(value.Value)
	})
	return chains
}

// addXssFilter 添加XSS过滤处理函数
func addXssFilter(chains DisplayProcessFnChains) DisplayProcessFnChains {
	chains = chains.Add(func(value FieldModel) interface{} {
		return html.EscapeString(value.Value)
	})
	return chains
}

// addXssJsFilter 添加XSS JavaScript过滤处理函数
func addXssJsFilter(chains DisplayProcessFnChains) DisplayProcessFnChains {
	chains = chains.Add(func(value FieldModel) interface{} {
		replacer := strings.NewReplacer("<script>", "&lt;script&gt;", "</script>", "&lt;/script&gt;")
		return replacer.Replace(value.Value)
	})
	return chains
}

// setDefaultDisplayFnOfFormType 设置表单类型的默认显示函数
func setDefaultDisplayFnOfFormType(f *FormPanel, typ form.Type) {
	// 如果是多文件类型
	if typ.IsMultiFile() {
		f.FieldList[f.curFieldListIndex].Display = func(value FieldModel) interface{} {
			if value.Value == "" {
				return ""
			}
			arr := strings.Split(value.Value, ",")
			res := "["
			for i, item := range arr {
				if i == len(arr)-1 {
					res += "'" + config.GetStore().URL(item) + "']"
				} else {
					res += "'" + config.GetStore().URL(item) + "',"
				}
			}
			return res
		}
	}
	// 如果是选择类型
	if typ.IsSelect() {
		f.FieldList[f.curFieldListIndex].Display = func(value FieldModel) interface{} {
			return strings.Split(value.Value, ",")
		}
	}
}

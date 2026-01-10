package types

import (
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/purpose168/GoAdmin/modules/config"

	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/modules/db"
	"github.com/purpose168/GoAdmin/modules/errors"
	"github.com/purpose168/GoAdmin/modules/language"
	"github.com/purpose168/GoAdmin/modules/logger"
	"github.com/purpose168/GoAdmin/modules/utils"
	"github.com/purpose168/GoAdmin/plugins/admin/modules"
	"github.com/purpose168/GoAdmin/plugins/admin/modules/parameter"
	"github.com/purpose168/GoAdmin/template/types/form"
	"github.com/purpose168/GoAdmin/template/types/table"
)

// FieldModel is the single query result.
type FieldModel struct {
	// The primaryKey of the table.
	ID string

	// The value of the single query result.
	Value string

	// The current row data.
	Row map[string]interface{}

	// Post type
	PostType PostType
}

type PostType uint8

const (
	PostTypeCreate = iota
	PostTypeUpdate
)

func (m FieldModel) IsCreate() bool {
	return m.PostType == PostTypeCreate
}

func (m FieldModel) IsUpdate() bool {
	return m.PostType == PostTypeUpdate
}

// PostFieldModel 包含单个查询结果的ID和值以及当前行数据
type PostFieldModel struct {
	ID    string            // 记录ID
	Value FieldModelValue   // 字段模型值
	Row   map[string]string // 当前行数据
	// 提交类型
	PostType PostType
}

// IsCreate 判断是否为创建操作
// 返回: 如果是创建操作返回true，否则返回false
func (m PostFieldModel) IsCreate() bool {
	return m.PostType == PostTypeCreate
}

// IsUpdate 判断是否为更新操作
// 返回: 如果是更新操作返回true，否则返回false
func (m PostFieldModel) IsUpdate() bool {
	return m.PostType == PostTypeUpdate
}

// InfoList 是信息项列表的别名
type InfoList []map[string]InfoItem

// InfoItem 表示信息项
type InfoItem struct {
	Content template.HTML `json:"content"` // HTML内容
	Value   string        `json:"value"`   // 值
}

// GroupBy 根据标签组对信息列表进行分组
// 参数:
//   - groups: 标签组
//
// 返回: 分组后的信息列表数组
func (i InfoList) GroupBy(groups TabGroups) []InfoList {

	var res = make([]InfoList, len(groups))

	for key, value := range groups {
		var newInfoList = make(InfoList, len(i))

		for index, info := range i {
			var newRow = make(map[string]InfoItem)
			for mk, m := range info {
				if modules.InArray(value, mk) {
					newRow[mk] = m
				}
			}
			newInfoList[index] = newRow
		}

		res[key] = newInfoList
	}

	return res
}

// Callbacks 是回调节点列表
type Callbacks []context.Node

// AddCallback 添加回调节点
// 参数:
//   - node: 上下文节点
//
// 返回: 更新后的回调列表
func (c Callbacks) AddCallback(node context.Node) Callbacks {
	if node.Path != "" && node.Method != "" && len(node.Handlers) > 0 {
		for _, item := range c {
			if strings.EqualFold(item.Path, node.Path) &&
				strings.EqualFold(item.Method, node.Method) {
				return c
			}
		}
		parr := strings.Split(node.Path, "?")
		if len(parr) > 1 {
			node.Path = parr[0]
			return append(c, node)
		}
		return append(c, node)
	}
	return c
}

// FieldModelValue 是字段模型值类型
type FieldModelValue []string

// Value 获取字段模型的值
// 返回: 第一个值
func (r FieldModelValue) Value() string {
	return r.First()
}

// First 获取第一个值
// 返回: 第一个值，如果不存在则返回空字符串
func (r FieldModelValue) First() string {
	if len(r) > 0 {
		return r[0]
	}
	return ""
}

// FieldFilterFn 是数据过滤函数类型
type FieldFilterFn func(value FieldModel) interface{}

// PostFieldFilterFn 是数据过滤函数类型
type PostFieldFilterFn func(value PostFieldModel) interface{}

// Field 表示表字段
type Field struct {
	Head     string          // 字段标题
	Field    string          // 字段名
	TypeName db.DatabaseType // 数据库类型

	Joins Joins // 关联配置

	Width       int  // 字段宽度
	Sortable    bool // 是否可排序
	EditAble    bool // 是否可编辑
	Fixed       bool // 是否固定
	Filterable  bool // 是否可筛选
	Hide        bool // 是否隐藏
	HideForList bool // 在列表中是否隐藏

	EditType    table.Type   // 编辑类型
	EditOptions FieldOptions // 编辑选项

	FilterFormFields []FilterFormField // 筛选表单字段

	IsEditParam   bool // 是否为编辑参数
	IsDeleteParam bool // 是否为删除参数
	IsDetailParam bool // 是否为详情参数

	FieldDisplay // 字段显示配置
}

// QueryFilterFn 是查询过滤函数类型
type QueryFilterFn func(param parameter.Parameters, conn db.Connection) (ids []string, stopQuery bool)

// UpdateParametersFn 是更新参数函数类型
type UpdateParametersFn func(param *parameter.Parameters)

// FilterFormField 是筛选表单字段结构体
type FilterFormField struct {
	Type        form.Type           // 表单类型
	Options     FieldOptions        // 选项列表
	OptionTable OptionTable         // 选项表配置
	Width       int                 // 字段宽度
	HeadWidth   int                 // 标题宽度
	InputWidth  int                 // 输入框宽度
	Style       template.HTMLAttr   // 样式属性
	Operator    FilterOperator      // 筛选操作符
	OptionExt   template.JS         // 选项扩展配置
	Head        string              // 标题
	Placeholder string              // 占位符
	HelpMsg     template.HTML       // 帮助信息
	NoIcon      bool                // 是否不显示图标
	ProcessFn   func(string) string // 处理函数
}

// GetFilterFormFields 获取筛选表单字段
// 参数:
//   - params: 参数对象
//   - headField: 头字段名
//   - sql: 可选的SQL对象
//
// 返回: 筛选表单字段列表
func (f Field) GetFilterFormFields(params parameter.Parameters, headField string, sql ...*db.SQL) []FormField {

	var (
		filterForm               = make([]FormField, 0)
		value, value2, keySuffix string
	)

	for index, filter := range f.FilterFormFields {

		if index > 0 {
			keySuffix = parameter.FilterParamCountInfix + strconv.Itoa(index)
		}

		if filter.Type.IsRange() {
			value = params.GetFilterFieldValueStart(headField)
			value2 = params.GetFilterFieldValueEnd(headField)
		} else if filter.Type.IsMultiSelect() {
			value = params.GetFieldValuesStr(headField)
		} else {
			if filter.Operator == FilterOperatorFree {
				value2 = GetOperatorFromValue(params.GetFieldOperator(headField, keySuffix)).String()
			}
			value = params.GetFieldValue(headField + keySuffix)
		}

		var (
			optionExt1 = filter.OptionExt
			optionExt2 template.JS
		)

		if filter.OptionExt == template.JS("") {
			op1, op2, js := filter.Type.GetDefaultOptions(headField + keySuffix)
			if op1 != nil {
				s, _ := json.Marshal(op1)
				optionExt1 = template.JS(string(s))
			}
			if op2 != nil {
				s, _ := json.Marshal(op2)
				optionExt2 = template.JS(string(s))
			}
			if js != template.JS("") {
				optionExt1 = js
			}
		}

		field := &FormField{
			Field:       headField + keySuffix,
			FieldClass:  headField + keySuffix,
			Head:        filter.Head,
			TypeName:    f.TypeName,
			HelpMsg:     filter.HelpMsg,
			NoIcon:      filter.NoIcon,
			FormType:    filter.Type,
			Editable:    true,
			Width:       filter.Width,
			HeadWidth:   filter.HeadWidth,
			InputWidth:  filter.InputWidth,
			Style:       filter.Style,
			Placeholder: filter.Placeholder,
			Value:       template.HTML(value),
			Value2:      value2,
			Options:     filter.Options,
			OptionExt:   optionExt1,
			OptionExt2:  optionExt2,
			OptionTable: filter.OptionTable,
			Label:       filter.Operator.Label(),
		}

		field.setOptionsFromSQL(sql[0])

		if filter.Type.IsSingleSelect() {
			field.Options = field.Options.SetSelected(params.GetFieldValue(f.Field), filter.Type.SelectedLabel())
		}

		if filter.Type.IsMultiSelect() {
			field.Options = field.Options.SetSelected(params.GetFieldValues(f.Field), filter.Type.SelectedLabel())
		}

		filterForm = append(filterForm, *field)

		if filter.Operator.AddOrNot() {
			ff := headField + parameter.FilterParamOperatorSuffix + keySuffix
			filterForm = append(filterForm, FormField{
				Field:      ff,
				FieldClass: ff,
				Head:       f.Head,
				TypeName:   f.TypeName,
				Value:      template.HTML(filter.Operator.Value()),
				FormType:   filter.Type,
				Hide:       true,
			})
		}
	}

	return filterForm
}

// Exist 判断字段是否存在
// 返回: 如果字段存在返回true，否则返回false
func (f Field) Exist() bool {
	return f.Field != ""
}

// FieldList 是字段列表类型
type FieldList []Field

// TableInfo 是表信息结构体
type TableInfo struct {
	Table      string // 表名
	PrimaryKey string // 主键
	Delimiter  string // 左分隔符
	Delimiter2 string // 右分隔符
	Driver     string // 数据库驱动
}

// GetTheadAndFilterForm 获取表头和筛选表单
// 参数:
//   - info: 表信息
//   - params: 参数对象
//   - columns: 列名数组
//   - sql: 可选的SQL函数
//
// 返回: 表头、字段、关联字段、关联SQL、筛选表单
func (f FieldList) GetTheadAndFilterForm(info TableInfo, params parameter.Parameters, columns []string,
	sql ...func() *db.SQL) (Thead, string, string, string, []string, []FormField) {
	var (
		thead      = make(Thead, 0)
		fields     = ""
		joinFields = ""
		joins      = ""
		joinTables = make([]string, 0)
		filterForm = make([]FormField, 0)
		tableName  = info.Delimiter + info.Table + info.Delimiter2
	)
	for _, field := range f {
		if field.Field != info.PrimaryKey && modules.InArray(columns, field.Field) &&
			!field.Joins.Valid() {
			fields += tableName + "." + modules.FilterField(field.Field, info.Delimiter, info.Delimiter2) + ","
		}

		headField := field.Field

		if field.Joins.Valid() {
			headField = field.Joins.Last().GetTableName() + parameter.FilterParamJoinInfix + field.Field
			joinFields += db.GetAggregationExpression(info.Driver, field.Joins.Last().GetTableName(info.Delimiter, info.Delimiter2)+"."+
				modules.FilterField(field.Field, info.Delimiter, info.Delimiter2), headField, JoinFieldValueDelimiter) + ","
			for _, join := range field.Joins {
				if !modules.InArray(joinTables, join.GetTableName(info.Delimiter, info.Delimiter2)) {
					joinTables = append(joinTables, join.GetTableName(info.Delimiter, info.Delimiter2))
					if join.BaseTable == "" {
						join.BaseTable = info.Table
					}
					joins += " left join " + modules.FilterField(join.Table, info.Delimiter, info.Delimiter2) + " " + join.TableAlias + " on " +
						join.GetTableName(info.Delimiter, info.Delimiter2) + "." + modules.FilterField(join.JoinField, info.Delimiter, info.Delimiter2) + " = " +
						modules.Delimiter(info.Delimiter, info.Delimiter2, join.BaseTable) + "." + modules.FilterField(join.Field, info.Delimiter, info.Delimiter2)
				}
			}
		}

		if field.Filterable {
			if len(sql) > 0 {
				filterForm = append(filterForm, field.GetFilterFormFields(params, headField, sql[0]())...)
			} else {
				filterForm = append(filterForm, field.GetFilterFormFields(params, headField)...)
			}
		}

		if field.Hide {
			continue
		}
		if field.HideForList {
			continue
		}
		thead = append(thead, TheadItem{
			Head:       field.Head,
			Sortable:   field.Sortable,
			Field:      headField,
			Hide:       !modules.InArrayWithoutEmpty(params.Columns, headField),
			Editable:   field.EditAble,
			EditType:   field.EditType.String(),
			EditOption: field.EditOptions,
			Width:      strconv.Itoa(field.Width) + "px",
		})
	}

	return thead, fields, joinFields, joins, joinTables, filterForm
}

func (f FieldList) GetThead(info TableInfo, params parameter.Parameters, columns []string) (Thead, string, string) {
	var (
		thead      = make(Thead, 0)
		fields     = ""
		joins      = ""
		joinTables = make([]string, 0)
	)
	for _, field := range f {
		if field.Field != info.PrimaryKey && modules.InArray(columns, field.Field) &&
			!field.Joins.Valid() {
			fields += info.Table + "." + modules.FilterField(field.Field, info.Delimiter, info.Delimiter2) + ","
		}

		headField := field.Field

		if field.Joins.Valid() {
			headField = field.Joins.Last().GetTableName(info.Delimiter, info.Delimiter2) + parameter.FilterParamJoinInfix + field.Field
			for _, join := range field.Joins {
				if !modules.InArray(joinTables, join.GetTableName(info.Delimiter, info.Delimiter2)) {
					joinTables = append(joinTables, join.GetTableName(info.Delimiter, info.Delimiter2))
					if join.BaseTable == "" {
						join.BaseTable = info.Table
					}
					joins += " left join " + modules.FilterField(join.Table, info.Delimiter, info.Delimiter2) + " " + join.TableAlias + " on " +
						join.GetTableName(info.Delimiter, info.Delimiter2) + "." + modules.FilterField(join.JoinField, info.Delimiter, info.Delimiter2) + " = " +
						modules.Delimiter(info.Delimiter, info.Delimiter2, join.BaseTable) + "." + modules.FilterField(join.Field, info.Delimiter, info.Delimiter2)
				}
			}
		}

		if field.Hide {
			continue
		}
		thead = append(thead, TheadItem{
			Head:       field.Head,
			Sortable:   field.Sortable,
			Field:      headField,
			Hide:       !modules.InArrayWithoutEmpty(params.Columns, headField),
			Editable:   field.EditAble,
			EditType:   field.EditType.String(),
			EditOption: field.EditOptions,
			Width:      strconv.Itoa(field.Width) + "px",
		})
	}

	return thead, fields, joins
}

func (f FieldList) GetFieldFilterProcessValue(key, value, keyIndex string) string {
	field := f.GetFieldByFieldName(key)
	index := 0
	if keyIndex != "" {
		index, _ = strconv.Atoi(keyIndex)
	}
	if field.FilterFormFields != nil && len(field.FilterFormFields) > index {
		if field.FilterFormFields[index].ProcessFn != nil {
			value = field.FilterFormFields[index].ProcessFn(value)
		}
	}
	return value
}

func (f FieldList) GetFieldJoinTable(key string) string {
	field := f.GetFieldByFieldName(key)
	if field.Exist() {
		return field.Joins.Last().Table
	}
	return ""
}

func (f FieldList) GetFieldByFieldName(name string) Field {
	for _, field := range f {
		if field.Field == name {
			return field
		}
		if JoinField(field.Joins.Last().GetTableName(), field.Field) == name {
			return field
		}
	}
	return Field{}
}

// Join 存储关联表信息。例如:
//
//	Join {
//	    BaseTable:   "users",
//	    Field:       "role_id",
//	    Table:       "roles",
//	    JoinField:   "id",
//	}
//
// 它将生成如下关联表SQL:
//
// ... left join roles on roles.id = users.role_id ...
type Join struct {
	Table      string // 关联表名
	TableAlias string // 表别名
	Field      string // 当前表字段
	JoinField  string // 关联表字段
	BaseTable  string // 基础表名
}

type Joins []Join

func JoinField(table, field string) string {
	return table + parameter.FilterParamJoinInfix + field
}

func GetJoinField(field string) string {
	return strings.Split(field, parameter.FilterParamJoinInfix)[1]
}

func (j Joins) Valid() bool {
	for i := 0; i < len(j); i++ {
		if j[i].Valid() {
			return true
		}
	}
	return false
}

func (j Joins) Last() Join {
	if len(j) > 0 {
		return j[len(j)-1]
	}
	return Join{}
}

func (j Join) Valid() bool {
	return j.Table != "" && j.Field != "" && j.JoinField != ""
}

func (j Join) GetTableName(delimiter ...string) string {
	if j.TableAlias != "" {
		return j.TableAlias
	}
	if len(delimiter) > 0 {
		return delimiter[0] + j.Table + delimiter[1]
	}
	return j.Table
}

var JoinFieldValueDelimiter = utils.Uuid(8)

type TabGroups [][]string

func (t TabGroups) Valid() bool {
	return len(t) > 0
}

func NewTabGroups(items ...string) TabGroups {
	var t = make(TabGroups, 0)
	return append(t, items)
}

func (t TabGroups) AddGroup(items ...string) TabGroups {
	return append(t, items)
}

type TabHeaders []string

func (t TabHeaders) Add(header string) TabHeaders {
	return append(t, header)
}

type GetDataFn func(param parameter.Parameters) ([]map[string]interface{}, int)

type DeleteFn func(ids []string) error
type DeleteFnWithRes func(ids []string, res error) error

type Sort uint8

const (
	SortDesc Sort = iota
	SortAsc
)

type primaryKey struct {
	Type db.DatabaseType
	Name string
}

type ExportProcessFn func(param parameter.Parameters) (PanelInfo, error)

// InfoPanel
type InfoPanel struct {
	Ctx *context.Context

	FieldList         FieldList
	curFieldListIndex int

	Table       string
	Title       string
	Description string

	// Warn: may be deprecated future.
	TabGroups  TabGroups
	TabHeaders TabHeaders

	Sort      Sort
	SortField string

	PageSizeList    []int
	DefaultPageSize int

	ExportType      int
	ExportProcessFn ExportProcessFn

	primaryKey primaryKey

	IsHideNewButton    bool
	IsHideExportButton bool
	IsHideEditButton   bool
	IsHideDeleteButton bool
	IsHideDetailButton bool
	IsHideFilterButton bool
	IsHideRowSelector  bool
	IsHidePagination   bool
	IsHideFilterArea   bool
	IsHideQueryInfo    bool
	FilterFormLayout   form.Layout

	FilterFormHeadWidth  int
	FilterFormInputWidth int

	Wheres    Wheres
	WhereRaws WhereRaw

	Callbacks Callbacks

	Buttons Buttons

	TableLayout string

	DeleteHook  DeleteFn
	PreDeleteFn DeleteFn
	DeleteFn    DeleteFn

	DeleteHookWithRes DeleteFnWithRes

	GetDataFn GetDataFn

	processChains DisplayProcessFnChains

	ActionButtons    Buttons
	ActionButtonFold bool

	DisplayGeneratorRecords map[string]struct{}

	QueryFilterFn       QueryFilterFn
	UpdateParametersFns []UpdateParametersFn

	Wrapper ContentWrapper

	// column operation buttons
	Action     template.HTML
	HeaderHtml template.HTML
	FooterHtml template.HTML

	PageError     errors.PageError
	PageErrorHTML template.HTML

	NoCompress  bool
	HideSideBar bool

	AutoRefresh uint
}

type Where struct {
	Join     string
	Field    string
	Operator string
	Arg      interface{}
}

type Wheres []Where

func (whs Wheres) Statement(wheres, delimiter, delimiter2 string, whereArgs []interface{}, existKeys, columns []string) (string, []interface{}) {
	pwheres := ""
	for k, wh := range whs {

		whFieldArr := strings.Split(wh.Field, ".")
		whField := ""
		whTable := ""
		if len(whFieldArr) > 1 {
			whField = whFieldArr[1]
			whTable = whFieldArr[0]
		} else {
			whField = whFieldArr[0]
		}

		if modules.InArray(existKeys, whField) {
			continue
		}

		// TODO: 支持like操作和关联表
		if modules.InArray(columns, whField) {

			joinMark := ""
			if k != len(whs)-1 {
				joinMark = whs[k+1].Join
			}

			if whTable != "" {
				pwheres += whTable + "." + modules.FilterField(whField, delimiter, delimiter2) + " " + wh.Operator + " ? " + joinMark + " "
			} else {
				pwheres += modules.FilterField(whField, delimiter, delimiter2) + " " + wh.Operator + " ? " + joinMark + " "
			}
			whereArgs = append(whereArgs, wh.Arg)
		}
	}
	if wheres != "" && pwheres != "" {
		wheres += " and "
	}
	return wheres + pwheres, whereArgs
}

// WhereRaw 是原始WHERE条件结构体
type WhereRaw struct {
	Raw  string        // 原始SQL
	Args []interface{} // 参数数组
}

// check 检查原始WHERE条件中是否包含and或or关键字
// 返回: 关键字的位置
func (wh WhereRaw) check() int {
	index := 0
	for i := 0; i < len(wh.Raw); i++ {
		if wh.Raw[i] == ' ' {
			continue
		}
		if wh.Raw[i] == 'a' {
			if len(wh.Raw) < i+3 {
				break
			} else if wh.Raw[i+1] == 'n' && wh.Raw[i+2] == 'd' {
				index = i + 3
			}
		} else if wh.Raw[i] == 'o' {
			if len(wh.Raw) < i+2 {
				break
			} else if wh.Raw[i+1] == 'r' {
				index = i + 2
			}
		} else {
			break
		}
	}
	return index
}

// Statement 生成WHERE语句
// 参数:
//   - wheres: WHERE条件字符串
//   - whereArgs: WHERE参数数组
//
// 返回: WHERE语句和参数数组
func (wh WhereRaw) Statement(wheres string, whereArgs []interface{}) (string, []interface{}) {

	if wh.Raw == "" {
		return wheres, whereArgs
	}

	if wheres != "" {
		if wh.check() != 0 {
			wheres += wh.Raw + " "
		} else {
			wheres += " and " + wh.Raw + " "
		}

		whereArgs = append(whereArgs, wh.Args...)
	} else {
		wheres += wh.Raw[wh.check():] + " "
		whereArgs = append(whereArgs, wh.Args...)
	}

	return wheres, whereArgs
}

// Handler 是处理函数类型
type Handler func(ctx *context.Context) (success bool, msg string, data interface{})

// Wrap 包装处理函数为上下文处理函数
// 返回: 上下文处理函数
func (h Handler) Wrap() context.Handler {
	return func(ctx *context.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(err)
				ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
					"code": 500,
					"data": "",
					"msg":  "错误",
				})
			}
		}()

		code := 0
		s, m, d := h(ctx)

		if !s {
			code = 500
		}
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"code": code,
			"data": d,
			"msg":  m,
		})
	}
}

// ContentWrapper 是内容包装函数类型
type ContentWrapper func(content template.HTML) template.HTML

// Action 是操作接口
type Action interface {
	Js() template.JS                                  // 获取JavaScript
	BtnAttribute() template.HTML                      // 获取按钮属性
	BtnClass() template.HTML                          // 获取按钮类名
	ExtContent(ctx *context.Context) template.HTML    // 获取扩展内容
	FooterContent(ctx *context.Context) template.HTML // 获取底部内容
	SetBtnId(btnId string)                            // 设置按钮ID
	SetBtnData(data interface{})                      // 设置按钮数据
	GetCallbacks() context.Node                       // 获取回调
}

// NilAction 是空操作结构体
type NilAction struct{}

func (def *NilAction) SetBtnId(btnId string)                            {}
func (def *NilAction) SetBtnData(data interface{})                      {}
func (def *NilAction) Js() template.JS                                  { return "" }
func (def *NilAction) BtnAttribute() template.HTML                      { return "" }
func (def *NilAction) BtnClass() template.HTML                          { return "" }
func (def *NilAction) ExtContent(ctx *context.Context) template.HTML    { return "" }
func (def *NilAction) FooterContent(ctx *context.Context) template.HTML { return "" }
func (def *NilAction) GetCallbacks() context.Node                       { return context.Node{} }

// Actions 是操作列表
type Actions []Action

// DefaultAction 是默认操作结构体
type DefaultAction struct {
	Attr   template.HTML // 属性
	JS     template.JS   // JavaScript
	Ext    template.HTML // 扩展内容
	Footer template.HTML // 底部内容
}

// NewDefaultAction 创建默认操作
// 参数:
//   - attr: 属性HTML
//   - ext: 扩展内容HTML
//   - footer: 底部内容HTML
//   - js: JavaScript代码
//
// 返回: 默认操作对象
func NewDefaultAction(attr, ext, footer template.HTML, js template.JS) *DefaultAction {
	return &DefaultAction{Attr: attr, Ext: ext, Footer: footer, JS: js}
}

func (def *DefaultAction) SetBtnId(btnId string)                            {}
func (def *DefaultAction) SetBtnData(data interface{})                      {}
func (def *DefaultAction) Js() template.JS                                  { return def.JS }
func (def *DefaultAction) BtnAttribute() template.HTML                      { return def.Attr }
func (def *DefaultAction) BtnClass() template.HTML                          { return "" }
func (def *DefaultAction) ExtContent(ctx *context.Context) template.HTML    { return def.Ext }
func (def *DefaultAction) FooterContent(ctx *context.Context) template.HTML { return def.Footer }
func (def *DefaultAction) GetCallbacks() context.Node                       { return context.Node{} }

var _ Action = (*DefaultAction)(nil)

// DefaultPageSizeList 是默认页面大小列表
var DefaultPageSizeList = []int{10, 20, 30, 50, 100}

// DefaultPageSize 是默认页面大小
const DefaultPageSize = 10

// NewInfoPanel 创建新的信息面板
// 参数:
//   - ctx: 上下文对象
//   - pk: 主键
//
// 返回: 初始化后的信息面板
func NewInfoPanel(ctx *context.Context, pk string) *InfoPanel {
	return &InfoPanel{
		Ctx:                     ctx,
		curFieldListIndex:       -1,
		PageSizeList:            DefaultPageSizeList,
		DefaultPageSize:         DefaultPageSize,
		processChains:           make(DisplayProcessFnChains, 0),
		Buttons:                 make(Buttons, 0),
		Callbacks:               make(Callbacks, 0),
		DisplayGeneratorRecords: make(map[string]struct{}),
		Wheres:                  make([]Where, 0),
		WhereRaws:               WhereRaw{},
		SortField:               pk,
		TableLayout:             "auto",
		FilterFormInputWidth:    10,
		FilterFormHeadWidth:     2,
		AutoRefresh:             0,
	}
}

// Where 添加WHERE条件
// 参数:
//   - field: 字段名
//   - operator: 操作符
//   - arg: 参数
//
// 返回: 更新后的信息面板
func (i *InfoPanel) Where(field string, operator string, arg interface{}) *InfoPanel {
	i.Wheres = append(i.Wheres, Where{Field: field, Operator: operator, Arg: arg, Join: "and"})
	return i
}

// WhereOr 添加OR条件的WHERE条件
// 参数:
//   - field: 字段名
//   - operator: 操作符
//   - arg: 参数
//
// 返回: 更新后的信息面板
func (i *InfoPanel) WhereOr(field string, operator string, arg interface{}) *InfoPanel {
	i.Wheres = append(i.Wheres, Where{Field: field, Operator: operator, Arg: arg, Join: "or"})
	return i
}

// WhereRaw 添加原始WHERE条件
// 参数:
//   - raw: 原始SQL语句
//   - arg: 参数
//
// 返回: 更新后的信息面板
func (i *InfoPanel) WhereRaw(raw string, arg ...interface{}) *InfoPanel {
	i.WhereRaws.Raw = raw
	i.WhereRaws.Args = arg
	return i
}

// AddSelectBox 添加选择框
// 参数:
//   - ctx: 上下文对象
//   - placeholder: 占位符
//   - options: 选项列表
//   - action: 操作对象
//   - width: 可选的宽度
//
// 返回: 更新后的信息面板
func (i *InfoPanel) AddSelectBox(ctx *context.Context, placeholder string, options FieldOptions, action Action, width ...int) *InfoPanel {
	options = append(FieldOptions{{Value: "", Text: language.Get("All")}}, options...)
	action.SetBtnData(options)
	i.addButton(GetDefaultSelection(placeholder, options, action, width...)).
		addFooterHTML(action.FooterContent(ctx)).
		addCallback(action.GetCallbacks())

	return i
}

// ExportValue 设置导出类型为值
// 返回: 更新后的信息面板
func (i *InfoPanel) ExportValue() *InfoPanel {
	i.ExportType = 1
	return i
}

// IsExportValue 判断是否导出值
// 返回: 如果是导出值返回true，否则返回false
func (i *InfoPanel) IsExportValue() bool {
	return i.ExportType == 1
}

// AddButtonRaw 添加原始按钮
// 参数:
//   - ctx: 上下文对象
//   - btn: 按钮对象
//   - action: 操作对象
//
// 返回: 更新后的信息面板
func (i *InfoPanel) AddButtonRaw(ctx *context.Context, btn Button, action Action) *InfoPanel {
	i.Buttons = append(i.Buttons, btn)
	i.addFooterHTML(action.FooterContent(ctx)).addCallback(action.GetCallbacks())
	return i
}

// AddButton 添加按钮
// 参数:
//   - ctx: 上下文对象
//   - title: 按钮标题
//   - icon: 图标
//   - action: 操作对象
//   - color: 可选的颜色
//
// 返回: 更新后的信息面板
func (i *InfoPanel) AddButton(ctx *context.Context, title template.HTML, icon string, action Action, color ...template.HTML) *InfoPanel {
	i.addButton(GetDefaultButtonGroup(title, icon, action, color...)).
		addFooterHTML(action.FooterContent(ctx)).
		addCallback(action.GetCallbacks())
	return i
}

// AddActionIconButton 添加操作图标按钮
// 参数:
//   - ctx: 上下文对象
//   - icon: 图标
//   - action: 操作对象
//   - ids: 可选的ID列表
//
// 返回: 更新后的信息面板
func (i *InfoPanel) AddActionIconButton(ctx *context.Context, icon string, action Action, ids ...string) *InfoPanel {
	i.addActionButton(GetActionIconButton(icon, action, ids...)).
		addFooterHTML(action.FooterContent(ctx)).
		addCallback(action.GetCallbacks())

	return i
}

// AddActionButtonFront 在前面添加操作按钮
// 参数:
//   - ctx: 上下文对象
//   - title: 按钮标题
//   - action: 操作对象
//   - ids: 可选的ID列表
//
// 返回: 更新后的信息面板
func (i *InfoPanel) AddActionButtonFront(ctx *context.Context, title template.HTML, action Action, ids ...string) *InfoPanel {
	i.SetActionButtonFold()
	i.ActionButtons = append([]Button{GetActionButton(title, action, ids...)}, i.ActionButtons...)
	i.addFooterHTML(action.FooterContent(ctx)).
		addCallback(action.GetCallbacks())
	return i
}

// AddActionButton 添加操作按钮
// 参数:
//   - ctx: 上下文对象
//   - title: 按钮标题
//   - action: 操作对象
//   - ids: 可选的ID列表
//
// 返回: 更新后的信息面板
func (i *InfoPanel) AddActionButton(ctx *context.Context, title template.HTML, action Action, ids ...string) *InfoPanel {
	i.SetActionButtonFold()
	i.addActionButton(GetActionButton(title, action, ids...)).
		addFooterHTML(action.FooterContent(ctx)).
		addCallback(action.GetCallbacks())

	return i
}

// SetActionButtonFold 设置操作按钮折叠
// 返回: 更新后的信息面板
func (i *InfoPanel) SetActionButtonFold() *InfoPanel {
	i.ActionButtonFold = true
	return i
}

// AddLimitFilter 添加长度限制过滤器
// 参数:
//   - limit: 限制长度
//
// 返回: 更新后的信息面板
func (i *InfoPanel) AddLimitFilter(limit int) *InfoPanel {
	i.processChains = addLimit(limit, i.processChains)
	return i
}

// AddTrimSpaceFilter 添加去空格过滤器
// 返回: 更新后的信息面板
func (i *InfoPanel) AddTrimSpaceFilter() *InfoPanel {
	i.processChains = addTrimSpace(i.processChains)
	return i
}

// AddSubstrFilter 添加子字符串过滤器
// 参数:
//   - start: 起始位置
//   - end: 结束位置
//
// 返回: 更新后的信息面板
func (i *InfoPanel) AddSubstrFilter(start int, end int) *InfoPanel {
	i.processChains = addSubstr(start, end, i.processChains)
	return i
}

// AddToTitleFilter 添加标题过滤器
// 返回: 更新后的信息面板
func (i *InfoPanel) AddToTitleFilter() *InfoPanel {
	i.processChains = addToTitle(i.processChains)
	return i
}

// AddToUpperFilter 添加大写过滤器
// 返回: 更新后的信息面板
func (i *InfoPanel) AddToUpperFilter() *InfoPanel {
	i.processChains = addToUpper(i.processChains)
	return i
}

// AddToLowerFilter 添加小写过滤器
// 返回: 更新后的信息面板
func (i *InfoPanel) AddToLowerFilter() *InfoPanel {
	i.processChains = addToLower(i.processChains)
	return i
}

// AddXssFilter 添加XSS过滤器
// 返回: 更新后的信息面板
func (i *InfoPanel) AddXssFilter() *InfoPanel {
	i.processChains = addXssFilter(i.processChains)
	return i
}

// AddXssJsFilter 添加XSS JS过滤器
// 返回: 更新后的信息面板
func (i *InfoPanel) AddXssJsFilter() *InfoPanel {
	i.processChains = addXssJsFilter(i.processChains)
	return i
}

// SetExportProcessFn 设置导出处理函数
// 参数:
//   - fn: 导出处理函数
//
// 返回: 更新后的信息面板
func (i *InfoPanel) SetExportProcessFn(fn ExportProcessFn) *InfoPanel {
	i.ExportProcessFn = fn
	return i
}

// SetDeleteHook 设置删除钩子
// 参数:
//   - fn: 删除函数
//
// 返回: 更新后的信息面板
func (i *InfoPanel) SetDeleteHook(fn DeleteFn) *InfoPanel {
	i.DeleteHook = fn
	return i
}

// SetDeleteHookWithRes 设置带结果的删除钩子
// 参数:
//   - fn: 带结果的删除函数
//
// 返回: 更新后的信息面板
func (i *InfoPanel) SetDeleteHookWithRes(fn DeleteFnWithRes) *InfoPanel {
	i.DeleteHookWithRes = fn
	return i
}

// SetQueryFilterFn 设置查询过滤函数
// 参数:
//   - fn: 查询过滤函数
//
// 返回: 更新后的信息面板
func (i *InfoPanel) SetQueryFilterFn(fn QueryFilterFn) *InfoPanel {
	i.QueryFilterFn = fn
	return i
}

// AddUpdateParametersFn 添加更新参数函数
// 参数:
//   - fn: 更新参数函数
//
// 返回: 更新后的信息面板
func (i *InfoPanel) AddUpdateParametersFn(fn UpdateParametersFn) *InfoPanel {
	i.UpdateParametersFns = append(i.UpdateParametersFns, fn)
	return i
}

// SetWrapper 设置内容包装器
// 参数:
//   - wrapper: 内容包装函数
//
// 返回: 更新后的信息面板
func (i *InfoPanel) SetWrapper(wrapper ContentWrapper) *InfoPanel {
	i.Wrapper = wrapper
	return i
}

// SetPreDeleteFn 设置删除前函数
// 参数:
//   - fn: 删除函数
//
// 返回: 更新后的信息面板
func (i *InfoPanel) SetPreDeleteFn(fn DeleteFn) *InfoPanel {
	i.PreDeleteFn = fn
	return i
}

// SetDeleteFn 设置删除函数
// 参数:
//   - fn: 删除函数
//
// 返回: 更新后的信息面板
func (i *InfoPanel) SetDeleteFn(fn DeleteFn) *InfoPanel {
	i.DeleteFn = fn
	return i
}

// SetGetDataFn 设置获取数据函数
// 参数:
//   - fn: 获取数据函数
//
// 返回: 更新后的信息面板
func (i *InfoPanel) SetGetDataFn(fn GetDataFn) *InfoPanel {
	i.GetDataFn = fn
	return i
}

// SetPrimaryKey 设置主键
// 参数:
//   - name: 主键名
//   - typ: 数据库类型
//
// 返回: 更新后的信息面板
func (i *InfoPanel) SetPrimaryKey(name string, typ db.DatabaseType) *InfoPanel {
	i.primaryKey = primaryKey{Name: name, Type: typ}
	return i
}

// SetTableFixed 设置表格布局为固定
// 返回: 更新后的信息面板
func (i *InfoPanel) SetTableFixed() *InfoPanel {
	i.TableLayout = "fixed"
	return i
}

// AddColumn 添加列
// 参数:
//   - head: 列标题
//   - fun: 字段过滤函数
//
// 返回: 更新后的信息面板
func (i *InfoPanel) AddColumn(head string, fun FieldFilterFn) *InfoPanel {
	i.FieldList = append(i.FieldList, Field{
		Head:     head,
		Field:    head,
		TypeName: db.Varchar,
		Sortable: false,
		EditAble: false,
		EditType: table.Text,
		FieldDisplay: FieldDisplay{
			Display:              fun,
			DisplayProcessChains: chooseDisplayProcessChains(i.processChains),
		},
	})
	i.curFieldListIndex++
	return i
}

// AddColumnButtons 添加按钮列
// 参数:
//   - ctx: 上下文对象
//   - head: 列标题
//   - buttons: 按钮列表
//
// 返回: 更新后的信息面板
func (i *InfoPanel) AddColumnButtons(ctx *context.Context, head string, buttons ...Button) *InfoPanel {
	var content, js template.HTML
	for _, btn := range buttons {
		btn.GetAction().SetBtnId("." + btn.ID())
		btnContent, btnJs := btn.Content(ctx)
		content += btnContent
		js += template.HTML(btnJs)
		i.FooterHtml += template.HTML(ParseTableDataTmpl(btn.GetAction().FooterContent(ctx)))
		i.Callbacks = i.Callbacks.AddCallback(btn.GetAction().GetCallbacks())
	}
	i.FooterHtml += template.HTML("<script>") + template.HTML(ParseTableDataTmpl(js)) + template.HTML("</script>")
	i.FieldList = append(i.FieldList, Field{
		Head:     head,
		Field:    head,
		TypeName: db.Varchar,
		Sortable: false,
		EditAble: false,
		EditType: table.Text,
		FieldDisplay: FieldDisplay{
			Display: func(value FieldModel) interface{} {
				pk := db.GetValueFromDatabaseType(i.primaryKey.Type, value.Row[i.primaryKey.Name], i.isFromJSON())
				var v = make(map[string]InfoItem)
				for key, item := range value.Row {
					itemValue := fmt.Sprintf("%v", item)
					v[key] = InfoItem{Value: itemValue, Content: template.HTML(itemValue)}
				}
				return template.HTML(ParseTableDataTmplWithID(pk.HTML(), string(content), v))
			},
			DisplayProcessChains: chooseDisplayProcessChains(i.processChains),
		},
	})
	i.curFieldListIndex++
	return i
}

// AddFieldTr 添加带翻译的字段
// 参数:
//   - ctx: 上下文对象
//   - head: 字段标题
//   - field: 字段名
//   - typeName: 数据库类型
//
// 返回: 更新后的信息面板
func (i *InfoPanel) AddFieldTr(ctx *context.Context, head, field string, typeName db.DatabaseType) *InfoPanel {
	return i.AddFieldWithTranslation(ctx, head, field, typeName)
}

// AddFieldWithTranslation 添加带翻译的字段
// 参数:
//   - ctx: 上下文对象
//   - head: 字段标题
//   - field: 字段名
//   - typeName: 数据库类型
//
// 返回: 更新后的信息面板
func (i *InfoPanel) AddFieldWithTranslation(ctx *context.Context, head, field string, typeName db.DatabaseType) *InfoPanel {
	return i.AddField(language.GetWithLang(head, ctx.Lang()), field, typeName)
}

// AddField 添加字段
// 参数:
//   - head: 字段标题
//   - field: 字段名
//   - typeName: 数据库类型
//
// 返回: 更新后的信息面板
func (i *InfoPanel) AddField(head, field string, typeName db.DatabaseType) *InfoPanel {
	i.FieldList = append(i.FieldList, Field{
		Head:     head,
		Field:    field,
		TypeName: typeName,
		Sortable: false,
		Joins:    make(Joins, 0),
		EditAble: false,
		EditType: table.Text,
		FieldDisplay: FieldDisplay{
			Display: func(value FieldModel) interface{} {
				return value.Value
			},
			DisplayProcessChains: chooseDisplayProcessChains(i.processChains),
		},
	})
	i.curFieldListIndex++
	return i
}

// AddFilter 添加筛选字段
// 参数:
//   - head: 字段标题
//   - field: 字段名
//   - typeName: 数据库类型
//   - fn: 更新参数函数
//   - filterType: 可选的筛选类型
//
// 返回: 更新后的信息面板
func (i *InfoPanel) AddFilter(head, field string, typeName db.DatabaseType, fn UpdateParametersFn, filterType ...FilterType) *InfoPanel {
	return i.AddField(head, field, typeName).FieldHide().FieldFilterable(filterType...).AddUpdateParametersFn(fn)
}

// 字段属性设置函数
// ====================================================

// FieldDisplay 设置字段显示函数
// 参数:
//   - filter: 字段过滤函数
//
// 返回: 更新后的信息面板
func (i *InfoPanel) FieldDisplay(filter FieldFilterFn) *InfoPanel {
	i.FieldList[i.curFieldListIndex].Display = filter
	return i
}

// FieldLabelParam 是字段标签参数结构体
type FieldLabelParam struct {
	Color template.HTML // 颜色
	Type  string        // 类型
}

// FieldLabel 设置字段标签
// 参数:
//   - args: 标签参数
//
// 返回: 更新后的信息面板
func (i *InfoPanel) FieldLabel(args ...FieldLabelParam) *InfoPanel {
	i.addDisplayChains(displayFnGens["label"].Get(i.Ctx, args))
	return i
}

// FieldImage 设置字段为图片显示
// 参数:
//   - width: 宽度
//   - height: 高度
//   - prefix: 可选的前缀
//
// 返回: 更新后的信息面板
func (i *InfoPanel) FieldImage(width, height string, prefix ...string) *InfoPanel {
	i.addDisplayChains(displayFnGens["image"].Get(i.Ctx, width, height, prefix))
	return i
}

// FieldBool 设置字段为布尔值显示
// 参数:
//   - flags: 标志列表
//
// 返回: 更新后的信息面板
func (i *InfoPanel) FieldBool(flags ...string) *InfoPanel {
	i.addDisplayChains(displayFnGens["bool"].Get(i.Ctx, flags))
	return i
}

// FieldLink 设置字段为链接显示
// 参数:
//   - src: 链接地址
//   - openInNewTab: 是否在新标签页打开
//
// 返回: 更新后的信息面板
func (i *InfoPanel) FieldLink(src string, openInNewTab ...bool) *InfoPanel {
	i.addDisplayChains(displayFnGens["link"].Get(i.Ctx, src, openInNewTab))
	return i
}

// FieldFileSize 设置字段为文件大小显示
// 返回: 更新后的信息面板
func (i *InfoPanel) FieldFileSize() *InfoPanel {
	i.addDisplayChains(displayFnGens["filesize"].Get(i.Ctx))
	return i
}

// FieldDate 设置字段为日期显示
// 参数:
//   - format: 日期格式
//
// 返回: 更新后的信息面板
func (i *InfoPanel) FieldDate(format string) *InfoPanel {
	i.addDisplayChains(displayFnGens["date"].Get(i.Ctx, format))
	return i
}

// FieldIcon 设置字段为图标显示
// 参数:
//   - icons: 图标映射
//   - defaultIcon: 默认图标
//
// 返回: 更新后的信息面板
func (i *InfoPanel) FieldIcon(icons map[string]string, defaultIcon string) *InfoPanel {
	i.addDisplayChains(displayFnGens["link"].Get(i.Ctx, icons, defaultIcon))
	return i
}

// FieldDotColor 是字段点颜色类型
type FieldDotColor string

const (
	FieldDotColorDanger  FieldDotColor = "danger"  // 危险颜色
	FieldDotColorInfo    FieldDotColor = "info"    // 信息颜色
	FieldDotColorPrimary FieldDotColor = "primary" // 主要颜色
	FieldDotColorSuccess FieldDotColor = "success" // 成功颜色
)

// FieldDot 设置字段为点显示
// 参数:
//   - icons: 图标映射
//   - defaultDot: 默认点颜色
//
// 返回: 更新后的信息面板
func (i *InfoPanel) FieldDot(icons map[string]FieldDotColor, defaultDot FieldDotColor) *InfoPanel {
	i.addDisplayChains(displayFnGens["dot"].Get(i.Ctx, icons, defaultDot))
	return i
}

// FieldProgressBarData 是进度条数据结构体
type FieldProgressBarData struct {
	Style string // 样式
	Size  string // 大小
	Max   int    // 最大值
}

// FieldProgressBar 设置字段为进度条显示
// 参数:
//   - data: 进度条数据
//
// 返回: 更新后的信息面板
func (i *InfoPanel) FieldProgressBar(data ...FieldProgressBarData) *InfoPanel {
	i.addDisplayChains(displayFnGens["progressbar"].Get(i.Ctx, data))
	return i
}

// FieldLoading 设置字段为加载中显示
// 参数:
//   - data: 数据列表
//
// 返回: 更新后的信息面板
func (i *InfoPanel) FieldLoading(data []string) *InfoPanel {
	i.addDisplayChains(displayFnGens["loading"].Get(i.Ctx, data))
	return i
}

// FieldDownLoadable 设置字段为可下载显示
// 参数:
//   - prefix: 可选的前缀
//
// 返回: 更新后的信息面板
func (i *InfoPanel) FieldDownLoadable(prefix ...string) *InfoPanel {
	i.addDisplayChains(displayFnGens["downloadable"].Get(i.Ctx, prefix))
	return i
}

// FieldCopyable 设置字段为可复制显示
// 参数:
//   - prefix: 可选的前缀
//
// 返回: 更新后的信息面板
func (i *InfoPanel) FieldCopyable(prefix ...string) *InfoPanel {
	i.addDisplayChains(displayFnGens["copyable"].Get(i.Ctx, prefix))
	if _, ok := i.DisplayGeneratorRecords["copyable"]; !ok {
		i.addFooterHTML(`<script>` + displayFnGens["copyable"].JS() + `</script>`)
		i.DisplayGeneratorRecords["copyable"] = struct{}{}
	}
	return i
}

// FieldGetImgArrFn 是获取图片数组函数类型
type FieldGetImgArrFn func(value string) []string

// FieldCarousel 设置字段为轮播图显示
// 参数:
//   - fn: 获取图片数组函数
//   - size: 可选的大小
//
// 返回: 更新后的信息面板
func (i *InfoPanel) FieldCarousel(fn FieldGetImgArrFn, size ...int) *InfoPanel {
	i.addDisplayChains(displayFnGens["carousel"].Get(i.Ctx, fn, size))
	return i
}

// FieldQrcode 设置字段为二维码显示
// 返回: 更新后的信息面板
func (i *InfoPanel) FieldQrcode() *InfoPanel {
	i.addDisplayChains(displayFnGens["qrcode"].Get(i.Ctx))
	if _, ok := i.DisplayGeneratorRecords["qrcode"]; !ok {
		i.addFooterHTML(`<script>` + displayFnGens["qrcode"].JS() + `</script>`)
		i.DisplayGeneratorRecords["qrcode"] = struct{}{}
	}
	return i
}

// FieldWidth 设置字段宽度
// 参数:
//   - width: 宽度
//
// 返回: 更新后的信息面板
func (i *InfoPanel) FieldWidth(width int) *InfoPanel {
	i.FieldList[i.curFieldListIndex].Width = width
	return i
}

// FieldSortable 设置字段可排序
// 返回: 更新后的信息面板
func (i *InfoPanel) FieldSortable() *InfoPanel {
	i.FieldList[i.curFieldListIndex].Sortable = true
	return i
}

func (i *InfoPanel) FieldEditOptions(options FieldOptions, extra ...map[string]string) *InfoPanel {
	if i.FieldList[i.curFieldListIndex].EditType.IsSwitch() {
		if len(extra) == 0 {
			options[0].Extra = map[string]string{
				"size":     "small",
				"onColor":  "primary",
				"offColor": "default",
			}
		} else {
			if extra[0]["size"] == "" {
				extra[0]["size"] = "small"
			}
			if extra[0]["onColor"] == "" {
				extra[0]["onColor"] = "primary"
			}
			if extra[0]["offColor"] == "" {
				extra[0]["offColor"] = "default"
			}
			options[0].Extra = extra[0]
		}
	}
	i.FieldList[i.curFieldListIndex].EditOptions = options
	return i
}

func (i *InfoPanel) FieldEditAble(editType ...table.Type) *InfoPanel {
	i.FieldList[i.curFieldListIndex].EditAble = true
	if len(editType) > 0 {
		i.FieldList[i.curFieldListIndex].EditType = editType[0]
	}
	return i
}

func (i *InfoPanel) FieldAsEditParam() *InfoPanel {
	i.FieldList[i.curFieldListIndex].IsEditParam = true
	return i
}

func (i *InfoPanel) FieldAsDeleteParam() *InfoPanel {
	i.FieldList[i.curFieldListIndex].IsDeleteParam = true
	return i
}

func (i *InfoPanel) FieldAsDetailParam() *InfoPanel {
	i.FieldList[i.curFieldListIndex].IsDetailParam = true
	return i
}

func (i *InfoPanel) FieldFixed() *InfoPanel {
	i.FieldList[i.curFieldListIndex].Fixed = true
	return i
}

type FilterType struct {
	Options     FieldOptions
	Process     func(string) string
	OptionExt   map[string]interface{}
	FormType    form.Type
	HelpMsg     template.HTML
	Style       template.HTMLAttr
	Operator    FilterOperator
	Head        string
	Placeholder string
	Width       int
	HeadWidth   int
	InputWidth  int
	NoHead      bool
	NoIcon      bool
}

// FieldFilterable set a field filterable which will display in the filter box.
func (i *InfoPanel) FieldFilterable(filterType ...FilterType) *InfoPanel {
	i.FieldList[i.curFieldListIndex].Filterable = true

	if len(filterType) == 0 {
		i.FieldList[i.curFieldListIndex].FilterFormFields = append(i.FieldList[i.curFieldListIndex].FilterFormFields,
			FilterFormField{
				Type:        form.Text,
				Head:        i.FieldList[i.curFieldListIndex].Head,
				Placeholder: language.Get("input") + " " + i.FieldList[i.curFieldListIndex].Head,
			})
	}

	for _, filter := range filterType {
		var ff FilterFormField
		ff.Operator = filter.Operator
		if filter.FormType == form.Default {
			ff.Type = form.Text
		} else {
			ff.Type = filter.FormType
		}
		ff.Head = modules.AorB(!filter.NoHead && filter.Head == "",
			i.FieldList[i.curFieldListIndex].Head, filter.Head)
		ff.Width = filter.Width
		ff.HeadWidth = filter.HeadWidth
		ff.InputWidth = filter.InputWidth
		ff.HelpMsg = filter.HelpMsg
		ff.NoIcon = filter.NoIcon
		ff.Style = filter.Style
		ff.ProcessFn = filter.Process
		ff.Placeholder = modules.AorB(filter.Placeholder == "", language.Get("input")+" "+ff.Head, filter.Placeholder)
		ff.Options = filter.Options
		if len(filter.OptionExt) > 0 {
			s, _ := json.Marshal(filter.OptionExt)
			ff.OptionExt = template.JS(s)
		}
		i.FieldList[i.curFieldListIndex].FilterFormFields = append(i.FieldList[i.curFieldListIndex].FilterFormFields, ff)
	}

	return i
}

// FieldFilterOptions set options for a filterable field to select. It will work when you set the
// FormType of the field to SelectSingle/Select/SelectBox.
func (i *InfoPanel) FieldFilterOptions(options FieldOptions) *InfoPanel {
	i.FieldList[i.curFieldListIndex].FilterFormFields[0].Options = options
	i.FieldList[i.curFieldListIndex].FilterFormFields[0].OptionExt = `{"allowClear": "true"}`
	return i
}

// FieldFilterOptionsFromTable set options for a filterable field to select. The options is from other table.
// For example,
//
//	`FieldFilterOptionsFromTable("roles", "name", "id")`
//
// will generate the sql like:
//
//	`select id, name from roles`.
//
// And the `id` will be the value of options, `name` is the text to be shown.
func (i *InfoPanel) FieldFilterOptionsFromTable(table, textFieldName, valueFieldName string, process ...OptionTableQueryProcessFn) *InfoPanel {
	var fn OptionTableQueryProcessFn
	if len(process) > 0 {
		fn = process[0]
	}
	i.FieldList[i.curFieldListIndex].FilterFormFields[0].OptionTable = OptionTable{
		Table:          table,
		TextField:      textFieldName,
		ValueField:     valueFieldName,
		QueryProcessFn: fn,
	}
	return i
}

// FieldFilterOptionExt set the option extension js of the field.
func (i *InfoPanel) FieldFilterOptionExt(m map[string]interface{}) *InfoPanel {
	s, _ := json.Marshal(m)
	i.FieldList[i.curFieldListIndex].FilterFormFields[0].OptionExt = template.JS(s)
	return i
}

// FieldFilterProcess process the field content.
// For example:
//
//	FieldFilterProcess(func(val string) string {
//			return val + "ms"
//	})
func (i *InfoPanel) FieldFilterProcess(process func(string) string) *InfoPanel {
	i.FieldList[i.curFieldListIndex].FilterFormFields[0].ProcessFn = process
	return i
}

// FieldFilterOnSearch set the url and the corresponding handler which has some backend logic and
// return the options of the field.
// For example:
//
//	FieldFilterOnSearch("/search/city", func(ctx *context.Context) (success bool, msg string, data interface{}) {
//		return true, "ok", selection.Data{
//			Results: selection.Options{
//				{Text: "GuangZhou", ID: "1"},
//				{Text: "ShenZhen", ID: "2"},
//				{Text: "BeiJing", ID: "3"},
//				{Text: "ShangHai", ID: "4"},
//			}
//		}
//	}, 1000)
func (i *InfoPanel) FieldFilterOnSearch(url string, handler Handler, delay ...int) *InfoPanel {
	ext, callback := searchJS(i.FieldList[i.curFieldListIndex].FilterFormFields[0].OptionExt, url, handler, delay...)
	i.FieldList[i.curFieldListIndex].FilterFormFields[0].OptionExt = ext
	i.Callbacks = append(i.Callbacks, callback)
	return i
}

// FieldFilterOnChooseCustom set the js that will be called when filter option be selected.
func (i *InfoPanel) FieldFilterOnChooseCustom(js template.HTML) *InfoPanel {
	i.FooterHtml += chooseCustomJS(i.FieldList[i.curFieldListIndex].Field, js)
	return i
}

// FieldFilterOnChooseMap set the actions that will be taken when filter option be selected.
// For example:
//
//	map[string]types.LinkField{
//	     "men": {Field: "ip", Value:"127.0.0.1"},
//	     "women": {Field: "ip", Hide: true},
//	     "other": {Field: "ip", Disable: true}
//	}
//
// mean when choose men, the value of field ip will be set to 127.0.0.1,
// when choose women, field ip will be hidden, and when choose other, field ip will be disabled.
func (i *InfoPanel) FieldFilterOnChooseMap(m map[string]LinkField) *InfoPanel {
	i.FooterHtml += chooseMapJS(i.FieldList[i.curFieldListIndex].Field, m)
	return i
}

// FieldFilterOnChoose set the given value of the given field when choose the value of val.
func (i *InfoPanel) FieldFilterOnChoose(val, field string, value template.HTML) *InfoPanel {
	i.FooterHtml += chooseJS(i.FieldList[i.curFieldListIndex].Field, field, val, value)
	return i
}

// OperationURL get the operation api url.
func (i *InfoPanel) OperationURL(id string) string {
	return config.Url("/operation/" + utils.WrapURL(id))
}

// FieldFilterOnChooseAjax set the url and handler that will be called when field be choosed.
// The handler will return the option of the field. It will help to link two or more form items.
// For example:
//
//	FieldFilterOnChooseAjax("city", "/search/city", func(ctx *context.Context) (success bool, msg string, data interface{}) {
//		return true, "ok", selection.Data{
//			Results: selection.Options{
//				{Text: "GuangZhou", ID: "1"},
//				{Text: "ShenZhen", ID: "2"},
//				{Text: "BeiJing", ID: "3"},
//				{Text: "ShangHai", ID: "4"},
//			}
//		}
//	})
//
// When you choose the country, it trigger the action of ajax which be sent to the given handler,
// and return the city options to the field city.
func (i *InfoPanel) FieldFilterOnChooseAjax(field, url string, handler Handler) *InfoPanel {
	js, callback := chooseAjax(i.FieldList[i.curFieldListIndex].Field, field, i.OperationURL(url), handler)
	i.FooterHtml += js
	i.Callbacks = append(i.Callbacks, callback)
	return i
}

// FieldFilterOnChooseHide hide the fields when value to be chosen.
func (i *InfoPanel) FieldFilterOnChooseHide(value string, field ...string) *InfoPanel {
	i.FooterHtml += chooseHideJS(i.FieldList[i.curFieldListIndex].Field, []string{value}, field...)
	return i
}

// FieldFilterOnChooseShow display the fields when value to be chosen.
func (i *InfoPanel) FieldFilterOnChooseShow(value string, field ...string) *InfoPanel {
	i.FooterHtml += chooseShowJS(i.FieldList[i.curFieldListIndex].Field, []string{value}, field...)
	return i
}

// FieldFilterOnChooseDisable disable the fields when value to be chosen.
func (i *InfoPanel) FieldFilterOnChooseDisable(value string, field ...string) *InfoPanel {
	i.FooterHtml += chooseDisableJS(i.FieldList[i.curFieldListIndex].Field, []string{value}, field...)
	return i
}

// FieldHide hide field. Include the filter area.
func (i *InfoPanel) FieldHide() *InfoPanel {
	i.FieldList[i.curFieldListIndex].Hide = true
	return i
}

// FieldHide hide field for only the table.
func (i *InfoPanel) FieldHideForList() *InfoPanel {
	i.FieldList[i.curFieldListIndex].HideForList = true
	return i
}

// FieldJoin gets the field of the concatenated table.
//
//	Join {
//	    BaseTable:   "users",
//	    Field:       "role_id",
//	    Table:       "roles",
//	    JoinField:   "id",
//	}
//
// It will generate the join table sql like:
//
// select ... from users left join roles on roles.id = users.role_id
func (i *InfoPanel) FieldJoin(join Join) *InfoPanel {
	i.FieldList[i.curFieldListIndex].Joins = append(i.FieldList[i.curFieldListIndex].Joins, join)
	return i
}

// FieldLimit limit the field length.
func (i *InfoPanel) FieldLimit(limit int) *InfoPanel {
	i.FieldList[i.curFieldListIndex].DisplayProcessChains = i.FieldList[i.curFieldListIndex].AddLimit(limit)
	return i
}

// FieldTrimSpace trim space of the field.
func (i *InfoPanel) FieldTrimSpace() *InfoPanel {
	i.FieldList[i.curFieldListIndex].DisplayProcessChains = i.FieldList[i.curFieldListIndex].AddTrimSpace()
	return i
}

// FieldSubstr intercept string of the field.
func (i *InfoPanel) FieldSubstr(start int, end int) *InfoPanel {
	i.FieldList[i.curFieldListIndex].DisplayProcessChains = i.FieldList[i.curFieldListIndex].AddSubstr(start, end)
	return i
}

// FieldToTitle update the field to a string that begin words mapped to their Unicode title case.
func (i *InfoPanel) FieldToTitle() *InfoPanel {
	i.FieldList[i.curFieldListIndex].DisplayProcessChains = i.FieldList[i.curFieldListIndex].AddToTitle()
	return i
}

// FieldToUpper update the field to a string with all Unicode letters mapped to their upper case.
func (i *InfoPanel) FieldToUpper() *InfoPanel {
	i.FieldList[i.curFieldListIndex].DisplayProcessChains = i.FieldList[i.curFieldListIndex].AddToUpper()
	return i
}

// FieldToLower update the field to a string with all Unicode letters mapped to their lower case.
func (i *InfoPanel) FieldToLower() *InfoPanel {
	i.FieldList[i.curFieldListIndex].DisplayProcessChains = i.FieldList[i.curFieldListIndex].AddToLower()
	return i
}

// FieldXssFilter escape field with html.Escape.
func (i *InfoPanel) FieldXssFilter() *InfoPanel {
	i.FieldList[i.curFieldListIndex].DisplayProcessChains = i.FieldList[i.curFieldListIndex].DisplayProcessChains.
		Add(func(value FieldModel) interface{} {
			return html.EscapeString(value.Value)
		})
	return i
}

// InfoPanel attribute setting functions
// ====================================================

func (i *InfoPanel) SetTable(table string) *InfoPanel {
	i.Table = table
	return i
}

func (i *InfoPanel) SetPageSizeList(pageSizeList []int) *InfoPanel {
	i.PageSizeList = pageSizeList
	return i
}

func (i *InfoPanel) SetDefaultPageSize(defaultPageSize int) *InfoPanel {
	i.DefaultPageSize = defaultPageSize
	return i
}

func (i *InfoPanel) GetPageSizeList() []string {
	var pageSizeList = make([]string, len(i.PageSizeList))
	for j := 0; j < len(i.PageSizeList); j++ {
		pageSizeList[j] = strconv.Itoa(i.PageSizeList[j])
	}
	return pageSizeList
}

func (i *InfoPanel) GetSort() string {
	switch i.Sort {
	case SortAsc:
		return "asc"
	default:
		return "desc"
	}
}

func (i *InfoPanel) SetTitle(title string) *InfoPanel {
	i.Title = title
	return i
}

func (i *InfoPanel) SetTabGroups(groups TabGroups) *InfoPanel {
	i.TabGroups = groups
	return i
}

func (i *InfoPanel) SetTabHeaders(headers ...string) *InfoPanel {
	i.TabHeaders = headers
	return i
}

func (i *InfoPanel) SetDescription(desc string) *InfoPanel {
	i.Description = desc
	return i
}

func (i *InfoPanel) SetFilterFormLayout(layout form.Layout) *InfoPanel {
	i.FilterFormLayout = layout
	return i
}

func (i *InfoPanel) SetFilterFormHeadWidth(w int) *InfoPanel {
	i.FilterFormHeadWidth = w
	return i
}

func (i *InfoPanel) SetFilterFormInputWidth(w int) *InfoPanel {
	i.FilterFormInputWidth = w
	return i
}

func (i *InfoPanel) SetSortField(field string) *InfoPanel {
	i.SortField = field
	return i
}

func (i *InfoPanel) SetSortAsc() *InfoPanel {
	i.Sort = SortAsc
	return i
}

func (i *InfoPanel) SetSortDesc() *InfoPanel {
	i.Sort = SortDesc
	return i
}

func (i *InfoPanel) SetAction(action template.HTML) *InfoPanel {
	i.Action = action
	return i
}

func (i *InfoPanel) SetHeaderHtml(header template.HTML) *InfoPanel {
	i.HeaderHtml += header
	return i
}

func (i *InfoPanel) SetFooterHtml(footer template.HTML) *InfoPanel {
	i.FooterHtml += footer
	return i
}

func (i *InfoPanel) HasError() bool {
	return i.PageError != nil
}

func (i *InfoPanel) SetError(err errors.PageError, content ...template.HTML) *InfoPanel {
	i.PageError = err
	if len(content) > 0 {
		i.PageErrorHTML = content[0]
	}
	return i
}

func (i *InfoPanel) SetNoCompress() *InfoPanel {
	i.NoCompress = true
	return i
}

func (i *InfoPanel) SetHideSideBar() *InfoPanel {
	i.HideSideBar = true
	return i
}

func (i *InfoPanel) SetAutoRefresh(interval uint) *InfoPanel {
	i.AutoRefresh = interval
	return i
}

func (i *InfoPanel) Set404Error(content ...template.HTML) *InfoPanel {
	i.SetError(errors.PageError404, content...)
	return i
}

func (i *InfoPanel) Set403Error(content ...template.HTML) *InfoPanel {
	i.SetError(errors.PageError403, content...)
	return i
}

func (i *InfoPanel) Set400Error(content ...template.HTML) *InfoPanel {
	i.SetError(errors.PageError401, content...)
	return i
}

func (i *InfoPanel) Set500Error(content ...template.HTML) *InfoPanel {
	i.SetError(errors.PageError500, content...)
	return i
}

func (i *InfoPanel) HideNewButton() *InfoPanel {
	i.IsHideNewButton = true
	return i
}

func (i *InfoPanel) HideExportButton() *InfoPanel {
	i.IsHideExportButton = true
	return i
}

func (i *InfoPanel) HideFilterButton() *InfoPanel {
	i.IsHideFilterButton = true
	return i
}

func (i *InfoPanel) HideRowSelector() *InfoPanel {
	i.IsHideRowSelector = true
	return i
}

func (i *InfoPanel) HidePagination() *InfoPanel {
	i.IsHidePagination = true
	return i
}

func (i *InfoPanel) HideFilterArea() *InfoPanel {
	i.IsHideFilterArea = true
	return i
}

func (i *InfoPanel) HideQueryInfo() *InfoPanel {
	i.IsHideQueryInfo = true
	return i
}

func (i *InfoPanel) HideEditButton() *InfoPanel {
	i.IsHideEditButton = true
	return i
}

func (i *InfoPanel) HideDeleteButton() *InfoPanel {
	i.IsHideDeleteButton = true
	return i
}

func (i *InfoPanel) HideDetailButton() *InfoPanel {
	i.IsHideDetailButton = true
	return i
}

func (i *InfoPanel) HideCheckBoxColumn() *InfoPanel {
	return i.HideColumn(1)
}

func (i *InfoPanel) HideColumn(n int) *InfoPanel {
	i.AddCSS(template.CSS(fmt.Sprintf(`
	.box-body table.table tbody tr td:nth-child(%v), .box-body table.table tbody tr th:nth-child(%v) {
		display: none;
	}`, n, n)))
	return i
}

func (i *InfoPanel) addFooterHTML(footer template.HTML) *InfoPanel {
	i.FooterHtml += template.HTML(ParseTableDataTmpl(footer))
	return i
}

func (i *InfoPanel) AddCSS(css template.CSS) *InfoPanel {
	return i.addFooterHTML(template.HTML("<style>" + css + "</style>"))
}

func (i *InfoPanel) AddJS(js template.JS) *InfoPanel {
	return i.addFooterHTML(template.HTML("<script>" + js + "</script>"))
}

func (i *InfoPanel) AddJSModule(js template.JS) *InfoPanel {
	return i.addFooterHTML(template.HTML("<script type='module'>" + js + "</script>"))
}

func (i *InfoPanel) addCallback(node context.Node) *InfoPanel {
	i.Callbacks = i.Callbacks.AddCallback(node)
	return i
}

func (i *InfoPanel) addButton(btn Button) *InfoPanel {
	i.Buttons = append(i.Buttons, btn)
	return i
}

func (i *InfoPanel) addActionButton(btn Button) *InfoPanel {
	i.ActionButtons = append(i.ActionButtons, btn)
	return i
}

func (i *InfoPanel) isFromJSON() bool {
	return i.GetDataFn != nil
}

func (i *InfoPanel) addDisplayChains(fn FieldFilterFn) *InfoPanel {
	i.FieldList[i.curFieldListIndex].DisplayProcessChains =
		i.FieldList[i.curFieldListIndex].DisplayProcessChains.Add(fn)
	return i
}

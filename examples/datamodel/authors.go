package datamodel

import (
	"github.com/purpose168/GoAdmin/context"
	"github.com/purpose168/GoAdmin/modules/db"
	"github.com/purpose168/GoAdmin/plugins/admin/modules/table"
	"github.com/purpose168/GoAdmin/template/types/form"
)

// GetAuthorsTable 返回作者表的数据模型配置
// 该函数创建并配置一个用于管理作者信息的表格，包含表格列表和表单两个部分
// 参数 ctx：上下文对象，用于传递请求上下文信息
// 返回值：配置完成的表格对象，可直接用于后台管理系统
func GetAuthorsTable(ctx *context.Context) (authorsTable table.Table) {

	// 使用默认配置创建表格实例
	// DefaultConfig() 会使用配置文件中设置的默认数据库连接
	authorsTable = table.NewDefaultTable(ctx, table.DefaultConfig())

	// 如果需要使用自定义数据库连接，可以使用以下方式：
	// authorsTable = table.NewDefaultTable(ctx, table.DefaultConfigWithDriverAndConnection("mysql", "admin"))
	// 参数说明：
	//   - "mysql"：数据库驱动类型，支持 mysql、postgres、sqlite、mssql 等
	//   - "admin"：连接名称，需要在配置文件中预先定义

	// 获取表格信息配置对象
	// info 对象用于配置表格列表页面的显示字段和属性
	info := authorsTable.GetInfo()

	// 添加 ID 字段
	// 参数说明：
	//   - "ID"：字段显示名称（中文）
	//   - "id"：数据库字段名
	//   - db.Int：字段类型为整型
	// FieldSortable()：设置该字段可排序，用户可以点击表头进行升序/降序排列
	info.AddField("ID", "id", db.Int).FieldSortable()

	// 添加名字字段
	// db.Varchar：字段类型为可变长度字符串，适合存储文本信息
	info.AddField("名字", "first_name", db.Varchar)

	// 添加姓氏字段
	info.AddField("姓氏", "last_name", db.Varchar)

	// 添加邮箱字段
	info.AddField("邮箱", "email", db.Varchar)

	// 添加出生日期字段
	// db.Date：字段类型为日期，格式为 YYYY-MM-DD
	info.AddField("出生日期", "birthdate", db.Date)

	// 添加添加时间字段
	// db.Timestamp：字段类型为时间戳，记录数据创建或更新的时间
	info.AddField("添加时间", "added", db.Timestamp)

	// 设置表格基本信息
	// SetTable("authors")：指定数据库表名为 authors
	// SetTitle("作者")：设置表格在后台显示的标题
	// SetDescription("作者")：设置表格的描述信息
	info.SetTable("authors").SetTitle("作者").SetDescription("作者")

	// 获取表单配置对象
	// formList 对象用于配置新增和编辑表单的字段和属性
	formList := authorsTable.GetForm()

	// 添加 ID 字段到表单
	// form.Default：使用默认表单控件类型
	// FieldDisplayButCanNotEditWhenUpdate()：在更新时显示但不可编辑
	// FieldDisableWhenCreate()：在创建时禁用该字段（通常由数据库自动生成）
	formList.AddField("ID", "id", db.Int, form.Default).FieldDisplayButCanNotEditWhenUpdate().FieldDisableWhenCreate()

	// 添加名字字段到表单
	// form.Text：使用文本输入框控件
	formList.AddField("名字", "first_name", db.Varchar, form.Text)

	// 添加姓氏字段到表单
	formList.AddField("姓氏", "last_name", db.Varchar, form.Text)

	// 添加邮箱字段到表单
	formList.AddField("邮箱", "email", db.Varchar, form.Text)

	// 添加出生日期字段到表单
	formList.AddField("出生日期", "birthdate", db.Date, form.Text)

	// 添加添加时间字段到表单
	formList.AddField("添加时间", "added", db.Timestamp, form.Text)

	// 设置表单基本信息
	// 与表格信息配置保持一致，使用相同的表名、标题和描述
	formList.SetTable("authors").SetTitle("作者").SetDescription("作者")

	return
}

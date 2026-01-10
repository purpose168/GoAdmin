// Package tables 提供数据库表模型的定义和配置
//
// 本包实现了 GoAdmin 管理后台中使用的各种数据表结构的定义和配置，提供以下功能：
//   - GetAuthorsTable: 获取作者表模型配置
//   - GetUserTable: 获取用户表模型配置
//   - Generators: 表生成器列表，用于初始化管理插件
//
// 核心概念：
//   - 表模型: 定义数据库表的结构和行为
//   - 表单配置: 配置表单字段的显示和验证规则
//   - 列表配置: 配置列表页面的显示和排序
//   - 字段类型: 定义数据库字段的数据类型（Int、Varchar、Date、Timestamp 等）
//   - 表单类型: 定义表单字段的输入类型（Text、Default、Select、Number 等）
//   - 字段选项: 配置字段的特殊行为（可排序、可编辑、禁用等）
//
// 技术栈：
//   - GoAdmin Table: 表格模块，提供表模型定义功能
//   - GoAdmin Form: 表单模块，提供表单配置功能
//   - GoAdmin Context: 上下文模块，提供请求上下文管理
//   - GoAdmin DB: 数据库模块，提供数据库操作功能
//
// 数据库支持：
//   - MySQL: 开源关系型数据库
//   - PostgreSQL: 高级开源关系型数据库
//   - SQLite: 轻量级嵌入式数据库
//   - MSSQL: Microsoft SQL Server 数据库
//
// 使用场景：
//   - 后台管理: 为 GoAdmin 管理后台提供数据表定义
//   - 数据管理: 管理数据库中的数据记录
//   - 表单生成: 自动生成数据录入表单
//   - 列表展示: 在管理界面展示数据列表
//   - 数据验证: 通过表单配置验证用户输入
//
// 配置说明：
//   - 表名: 对应数据库中的实际表名
//   - 字段名: 对应数据库表中的列名
//   - 字段类型: 必须与数据库字段类型匹配
//   - 表单类型: 决定表单字段的输入方式
//   - 字段选项: 控制字段的可编辑性、可见性等
//
// 注意事项：
//   - 需要确保数据库表已创建且结构正确
//   - 字段类型必须与数据库字段类型一致
//   - 表单字段配置必须与列表字段配置对应
//   - ID 字段通常设置为创建时禁用、更新时不可编辑
//   - 时间戳字段通常设置为自动更新
//
// 作者: GoAdmin Team
// 创建日期: 2019-01-01
// 版本: 1.0.0
package tables

import (
	"github.com/purpose168/GoAdmin/context"                     // 上下文包，提供请求上下文管理功能
	"github.com/purpose168/GoAdmin/modules/db"                  // 数据库模块包，提供数据库操作和字段类型定义
	"github.com/purpose168/GoAdmin/plugins/admin/modules/table" // 表格模块包，提供表模型定义和配置功能
	"github.com/purpose168/GoAdmin/template/types/form"         // 表单类型包，提供表单字段类型定义
)

// GetAuthorsTable 获取作者表模型配置
//
// 参数:
//   - ctx: 上下文对象，用于管理请求上下文和数据库连接
//
// 返回:
//   - table.Table: 作者表配置对象，包含列表显示和表单编辑两部分配置
//
// 说明：
//
//	该函数创建并返回一个完整的作者表配置，包括列表显示和表单编辑两部分。
//	使用默认配置连接数据库，也可以自定义数据库连接。
//
// 功能特性：
//   - 创建默认表配置
//   - 定义表的字段（ID、名、姓、邮箱、生日、添加时间）
//   - 配置列表显示字段
//   - 配置表单编辑字段
//   - 设置表的标题和描述
//   - 支持自定义数据库连接
//
// 字段说明：
//   - id: 主键，整数类型，自动递增
//   - first_name: 名，字符串类型，最大长度 255
//   - last_name: 姓，字符串类型，最大长度 255
//   - email: 邮箱，字符串类型，最大长度 255
//   - birthdate: 生日，日期类型
//   - added: 添加时间，时间戳类型
//
// 配置说明：
//   - 数据库表名: authors
//   - 表标题: 作者
//   - 表描述: 作者管理
//   - ID 字段: 可排序，创建时禁用，更新时不可编辑
//   - 其他字段: 可编辑，文本输入
//
// 技术细节：
//   - 使用 table.NewDefaultTable() 创建默认表配置
//   - 使用 table.DefaultConfig() 使用默认数据库配置
//   - 使用 table.DefaultConfigWithDriverAndConnection() 自定义数据库连接
//   - 使用 info.AddField() 添加列表显示字段
//   - 使用 info.SetTable() 设置数据库表名
//   - 使用 info.SetTitle() 设置表标题
//   - 使用 info.SetDescription() 设置表描述
//   - 使用 info.FieldSortable() 设置字段可排序
//   - 使用 formList.AddField() 添加表单字段
//   - 使用 formList.FieldDisplayButCanNotEditWhenUpdate() 设置字段更新时不可编辑
//   - 使用 formList.FieldDisableWhenCreate() 设置字段创建时禁用
//
// 使用场景：
//   - 后台管理: 在管理后台管理作者信息
//   - 数据录入: 录入和编辑作者数据
//   - 数据展示: 在列表页面展示作者信息
//   - 数据查询: 查询和筛选作者数据
//
// 注意事项：
//   - 需要确保数据库表 authors 已创建
//   - 字段类型必须与数据库字段类型一致
//   - ID 字段通常设置为创建时禁用、更新时不可编辑
//   - 时间戳字段通常设置为自动更新
//   - 邮箱字段可以添加邮箱验证规则
//
// 错误处理：
//   - 如果数据库连接失败，会在初始化时返回错误
//   - 如果表结构不匹配，会在运行时返回错误
//
// 示例：
//
//	// 创建上下文
//	ctx := context.NewContext(r)
//
//	// 获取作者表配置
//	authorsTable := GetAuthorsTable(ctx)
//
//	// 使用表配置
//	info := authorsTable.GetInfo()
//	formList := authorsTable.GetForm()
//
//	// 自定义数据库连接示例
//	// authorsTable = table.NewDefaultTable(ctx, table.DefaultConfigWithDriverAndConnection("mysql", "admin"))
func GetAuthorsTable(ctx *context.Context) (authorsTable table.Table) {

	authorsTable = table.NewDefaultTable(ctx, table.DefaultConfig()) // 创建默认表配置，使用默认数据库驱动和连接

	// 如果需要使用自定义数据库连接，可以使用以下代码：
	// 第一个参数是数据库驱动类型（如：mysql、postgres、sqlite等）
	// 第二个参数是数据库连接名称（在配置文件中定义的连接）
	// authorsTable = table.NewDefaultTable(ctx, table.DefaultConfigWithDriverAndConnection("mysql", "admin"))

	info := authorsTable.GetInfo() // 获取表信息配置对象，用于配置表的列表显示部分

	// 添加表字段
	// AddField 参数说明：
	//   - 第一个参数：字段显示名称（在列表中显示的标题）
	//   - 第二个参数：数据库字段名（对应数据库表的列名）
	//   - 第三个参数：字段类型（使用 db 包定义的类型常量）
	info.AddField("ID", "id", db.Int).FieldSortable() // ID 字段，整数类型，可排序
	info.AddField("名", "first_name", db.Varchar)      // 名字段，字符串类型
	info.AddField("姓", "last_name", db.Varchar)       // 姓字段，字符串类型
	info.AddField("邮箱", "email", db.Varchar)          // 邮箱字段，字符串类型
	info.AddField("生日", "birthdate", db.Date)         // 生日字段，日期类型
	info.AddField("添加时间", "added", db.Timestamp)      // 添加时间字段，时间戳类型

	// 设置表的基本信息
	// SetTable: 设置数据库表名
	// SetTitle: 设置表在界面中显示的标题
	// SetDescription: 设置表的描述信息
	info.SetTable("authors").SetTitle("作者").SetDescription("作者管理") // 设置表名为 authors，标题为"作者"，描述为"作者管理"

	formList := authorsTable.GetForm() // 获取表单配置对象，用于配置表的表单编辑部分（创建和编辑表单）

	// 添加表单字段
	// AddField 参数说明：
	//   - 第一个参数：字段显示名称（在表单中显示的标签）
	//   - 第二个参数：数据库字段名（对应数据库表的列名）
	//   - 第三个参数：字段类型（使用 db 包定义的类型常量）
	//   - 第四个参数：表单类型（使用 form 包定义的类型常量）
	formList.AddField("ID", "id", db.Int, form.Default).FieldDisplayButCanNotEditWhenUpdate().FieldDisableWhenCreate() // ID 字段：默认显示，更新时不可编辑，创建时禁用
	formList.AddField("名", "first_name", db.Varchar, form.Text)                                                        // 名字段：文本输入
	formList.AddField("姓", "last_name", db.Varchar, form.Text)                                                         // 姓字段：文本输入
	formList.AddField("邮箱", "email", db.Varchar, form.Text)                                                            // 邮箱字段：文本输入
	formList.AddField("生日", "birthdate", db.Date, form.Text)                                                           // 生日字段：文本输入
	formList.AddField("添加时间", "added", db.Timestamp, form.Text)                                                        // 添加时间字段：文本输入

	// 设置表单的基本信息
	// SetTable: 设置数据库表名
	// SetTitle: 设置表单在界面中显示的标题
	// SetDescription: 设置表单的描述信息
	formList.SetTable("authors").SetTitle("作者").SetDescription("作者管理") // 设置表名为 authors，标题为"作者"，描述为"作者管理"

	return // 返回作者表配置对象
}

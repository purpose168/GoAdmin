// Package tables 提供数据库表模型的定义和配置
//
// 本包实现了 GoAdmin 管理后台中使用的各种数据表结构的定义和配置，提供以下功能：
//   - GetAuthorsTable: 获取作者表模型配置
//   - GetUserTable: 获取用户表模型配置
//   - GetExternalTable: 获取外部数据源表模型配置
//   - Generators: 表生成器列表，用于初始化管理插件
//
// 核心概念：
//   - 表模型: 定义数据库表的结构和行为
//   - 外部数据源: 非数据库的数据源，如 API、缓存、文件等
//   - 表单配置: 配置表单字段的显示和验证规则
//   - 列表配置: 配置列表页面的显示和排序
//   - 详情配置: 配置详情页面的显示
//   - 字段类型: 定义数据库字段的数据类型（Int、Varchar、Date、Timestamp 等）
//   - 表单类型: 定义表单字段的输入类型（Text、Default、Select、Number 等）
//   - 字段选项: 配置字段的特殊行为（可排序、可编辑、禁用等）
//   - 自定义数据获取: 通过 SetGetDataFn 方法自定义数据获取逻辑
//
// 技术栈：
//   - GoAdmin Table: 表格模块，提供表模型定义功能
//   - GoAdmin Form: 表单模块，提供表单配置功能
//   - GoAdmin Context: 上下文模块，提供请求上下文管理
//   - GoAdmin DB: 数据库模块，提供数据库操作功能
//   - GoAdmin Parameter: 参数模块，提供查询参数管理
//
// 数据源类型：
//   - 数据库: MySQL、PostgreSQL、SQLite、MSSQL 等关系型数据库
//   - 外部数据源: API、缓存、文件、内存数据等非数据库数据源
//
// 使用场景：
//   - 后台管理: 为 GoAdmin 管理后台提供数据表定义
//   - 数据管理: 管理数据库中的数据记录
//   - 外部数据集成: 集成外部 API、缓存等数据源
//   - 表单生成: 自动生成数据录入表单
//   - 列表展示: 在管理界面展示数据列表
//   - 数据验证: 通过表单配置验证用户输入
//   - 自定义数据获取: 从非数据库数据源获取数据
//
// 配置说明：
//   - 表名: 对应数据库中的实际表名或虚拟表名
//   - 字段名: 对应数据库表中的列名或数据源的字段名
//   - 字段类型: 必须与数据源字段类型匹配
//   - 表单类型: 决定表单字段的输入方式
//   - 字段选项: 控制字段的可编辑性、可见性等
//   - 数据获取函数: 自定义数据获取逻辑，支持外部数据源
//
// 注意事项：
//   - 需要确保数据源已正确配置
//   - 字段类型必须与数据源字段类型一致
//   - 表单字段配置必须与列表字段配置对应
//   - ID 字段通常设置为创建时禁用、更新时不可编辑
//   - 外部数据源需要正确处理分页、排序和筛选
//   - 自定义数据获取函数需要返回正确的数据格式
//
// 作者: GoAdmin Team
// 创建日期: 2019-01-01
// 版本: 1.0.0
package tables

import (
	"github.com/purpose168/GoAdmin/context"                         // 上下文包，提供请求上下文管理功能
	"github.com/purpose168/GoAdmin/modules/db"                      // 数据库模块包，提供数据库操作和字段类型定义
	"github.com/purpose168/GoAdmin/plugins/admin/modules/parameter" // 参数模块包，提供查询参数管理功能
	"github.com/purpose168/GoAdmin/plugins/admin/modules/table"     // 表格模块包，提供表模型定义和配置功能
	"github.com/purpose168/GoAdmin/template/types/form"             // 表单类型包，提供表单字段类型定义
)

// GetExternalTable 从外部数据源获取表模型
//
// 参数:
//   - ctx: 上下文对象，用于管理请求上下文
//
// 返回:
//   - table.Table: 外部数据源表配置对象，包含列表显示、表单编辑和详情显示三部分配置
//
// 说明：
//
//	该函数展示了如何使用外部数据源（非数据库）创建表模型。
//	通过 SetGetDataFn 方法可以自定义数据获取逻辑，例如从 API、缓存或其他数据源获取数据。
//	返回值包含数据列表和总记录数。
//
// 功能特性：
//   - 创建默认表配置
//   - 定义表的字段（ID、标题）
//   - 配置列表显示字段
//   - 配置表单编辑字段
//   - 配置详情显示字段
//   - 设置表的标题和描述
//   - 使用自定义数据获取函数（SetGetDataFn）从外部数据源获取数据
//   - 配置列表页面的数据获取
//   - 配置详情页面的数据获取
//
// 字段说明：
//   - id: 主键，整数类型，用于唯一标识记录
//   - title: 标题，字符串类型，用于显示数据内容
//
// 配置说明：
//   - 数据库表名: external（虚拟表名）
//   - 表标题: 外部数据
//   - 表描述: 外部数据源管理
//   - ID 字段: 可排序，创建时禁用，更新时不可编辑
//   - 标题字段: 可编辑，文本输入
//
// 技术细节：
//   - 使用 table.NewDefaultTable() 创建默认表配置
//   - 使用 info.AddField() 添加列表显示字段
//   - 使用 info.SetTable() 设置表名
//   - 使用 info.SetTitle() 设置表标题
//   - 使用 info.SetDescription() 设置表描述
//   - 使用 info.SetGetDataFn() 设置自定义数据获取函数
//   - 使用 formList.AddField() 添加表单字段
//   - 使用 formList.FieldDisplayButCanNotEditWhenUpdate() 设置字段更新时不可编辑
//   - 使用 formList.FieldDisableWhenCreate() 设置字段创建时禁用
//   - 使用 detail.SetGetDataFn() 设置详情页数据获取函数
//
// 外部数据源说明：
//   - 数据获取函数参数: parameter.Parameters，包含查询参数（分页、排序、筛选等）
//   - 数据获取函数返回值: ([]map[string]interface{}, int)，第一个是数据列表，第二个是总记录数
//   - 列表数据获取: 用于列表页面，支持分页、排序、筛选
//   - 详情数据获取: 用于详情页面，根据 ID 获取单条数据
//
// 使用场景：
//   - API 集成: 从外部 API 获取数据并在管理后台展示
//   - 缓存集成: 从 Redis、Memcached 等缓存系统获取数据
//   - 文件集成: 从 JSON、XML、CSV 等文件读取数据
//   - 内存数据: 从内存中的数据结构获取数据
//   - 混合数据源: 同时使用数据库和外部数据源
//
// 注意事项：
//   - 需要确保外部数据源已正确配置
//   - 字段类型必须与数据源字段类型一致
//   - ID 字段通常设置为创建时禁用、更新时不可编辑
//   - 外部数据源需要正确处理分页、排序和筛选
//   - 自定义数据获取函数需要返回正确的数据格式
//   - 数据列表中的每个元素必须是 map[string]interface{} 类型
//   - 总记录数必须准确，用于分页计算
//
// 错误处理：
//   - 如果数据获取失败，会在运行时返回错误
//   - 如果数据格式不正确，会在运行时返回错误
//
// 示例：
//
//	// 创建上下文
//	ctx := context.NewContext(r)
//
//	// 获取外部数据源表配置
//	externalTable := GetExternalTable(ctx)
//
//	// 使用表配置
//	info := externalTable.GetInfo()
//	formList := externalTable.GetForm()
//	detail := externalTable.GetDetail()
//
//	// 实际应用中的数据获取函数示例
//	// 从 API 获取数据
//	SetGetDataFn(func(param parameter.Parameters) ([]map[string]interface{}, int) {
//	    // 调用外部 API
//	    response := callExternalAPI(param)
//	    // 返回数据列表和总记录数
//	    return response.Data, response.Total
//	})
func GetExternalTable(ctx *context.Context) (externalTable table.Table) {

	externalTable = table.NewDefaultTable(ctx) // 创建默认表配置

	info := externalTable.GetInfo() // 获取表信息配置对象，用于配置表的列表显示部分

	// 添加表字段
	// AddField 参数说明：
	//   - 第一个参数：字段显示名称（在列表中显示的标题）
	//   - 第二个参数：数据库字段名（对应数据源的字段名）
	//   - 第三个参数：字段类型（使用 db 包定义的类型常量）
	info.AddField("ID", "id", db.Int).FieldSortable() // ID 字段，整数类型，可排序
	info.AddField("标题", "title", db.Varchar)          // 标题字段，字符串类型

	// 设置表的基本信息并配置外部数据源
	// SetTable: 设置表名（可以是虚拟表名）
	// SetTitle: 设置表在界面中显示的标题
	// SetDescription: 设置表的描述信息
	// SetGetDataFn: 设置自定义数据获取函数
	//   参数：parameter.Parameters - 包含查询参数（分页、排序、筛选等）
	//   返回值：
	//     - []map[string]interface{}: 数据列表
	//     - int: 总记录数
	info.SetTable("external").
		SetTitle("外部数据").
		SetDescription("外部数据源管理").
		SetGetDataFn(func(param parameter.Parameters) ([]map[string]interface{}, int) {
			// 这里是模拟从外部数据源获取数据
			// 实际应用中，可以从 API、缓存、文件或其他数据源获取
			return []map[string]interface{}{
				{
					"id":    10,
					"title": "这是一个标题",
				}, {
					"id":    11,
					"title": "这是一个标题2",
				}, {
					"id":    12,
					"title": "这是一个标题3",
				}, {
					"id":    13,
					"title": "这是一个标题4",
				},
			}, 10 // 返回数据列表和总记录数
		})

	formList := externalTable.GetForm() // 获取表单配置对象，用于配置表的表单编辑部分（创建和编辑表单）

	// 添加表单字段
	// AddField 参数说明：
	//   - 第一个参数：字段显示名称（在表单中显示的标签）
	//   - 第二个参数：数据库字段名（对应数据源的字段名）
	//   - 第三个参数：字段类型（使用 db 包定义的类型常量）
	//   - 第四个参数：表单类型（使用 form 包定义的类型常量）
	formList.AddField("ID", "id", db.Int, form.Default).FieldDisplayButCanNotEditWhenUpdate().FieldDisableWhenCreate() // ID 字段：默认显示，更新时不可编辑，创建时禁用
	formList.AddField("标题", "title", db.Varchar, form.Text)                                                            // 标题字段：文本输入

	// 设置表单的基本信息
	// SetTable: 设置表名
	// SetTitle: 设置表单在界面中显示的标题
	// SetDescription: 设置表单的描述信息
	formList.SetTable("external").SetTitle("外部数据").SetDescription("外部数据源管理") // 设置表名为 external，标题为"外部数据"，描述为"外部数据源管理"

	detail := externalTable.GetDetail() // 获取详情页配置对象，用于配置表的详情显示部分

	// 设置详情页的基本信息并配置外部数据源
	// SetTable: 设置表名
	// SetTitle: 设置详情页在界面中显示的标题
	// SetDescription: 设置详情页的描述信息
	// SetGetDataFn: 设置自定义数据获取函数（用于详情页）
	detail.SetTable("external").
		SetTitle("外部数据").
		SetDescription("外部数据源管理").
		SetGetDataFn(func(param parameter.Parameters) ([]map[string]interface{}, int) {
			// 这里是模拟从外部数据源获取单条详情数据
			// 实际应用中，可以根据 param 中的 ID 参数获取对应的数据
			return []map[string]interface{}{
				{
					"id":    10,
					"title": "这是一个标题",
				},
			}, 1 // 返回数据列表和总记录数（详情页通常只有一条数据）
		})

	return // 返回外部数据源表配置对象
}

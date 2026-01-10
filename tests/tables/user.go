// Package tables 提供数据库表模型的定义和配置
//
// 本包实现了 GoAdmin 管理后台中使用的各种数据表结构的定义和配置，提供以下功能：
//   - GetAuthorsTable: 获取作者表模型配置
//   - GetUserTable: 获取用户表模型配置
//   - GetExternalTable: 获取外部数据源表模型配置
//   - GetPostsTable: 获取文章表模型配置
//   - Generators: 表生成器列表，用于初始化管理插件
//
// 核心概念：
//   - 表模型: 定义数据库表的结构和行为
//   - 自定义表配置: 通过 table.Config 自定义表的各种属性
//   - 字段类型: 定义数据库字段的数据类型（Int、Varchar、Tinyint、Timestamp 等）
//   - 表单类型: 定义表单字段的输入类型（Text、Number、Email、Password、RichText、Select、Switch 等）
//   - 字段选项: 配置字段的特殊行为（可排序、可编辑、可筛选、隐藏、关联等）
//   - 操作按钮: 配置行操作按钮和全局按钮（跳转、Ajax、弹窗等）
//   - 筛选器: 配置字段筛选器（下拉框、单选框、文本输入、日期时间范围等）
//   - 表单分组: 将表单字段分组到不同的标签页
//   - Ajax 联动: 通过 Ajax 实现字段之间的联动
//   - 表关联: 通过 FieldJoin 实现表之间的关联查询
//
// 技术栈：
//   - GoAdmin Table: 表格模块，提供表模型定义和配置功能
//   - GoAdmin Form: 表单模块，提供表单配置功能
//   - GoAdmin Context: 上下文模块，提供请求上下文管理
//   - GoAdmin DB: 数据库模块，提供数据库操作功能
//   - GoAdmin Types: 类型模块，提供类型定义和字段模型
//   - GoAdmin Action: 动作模块，提供按钮和动作配置功能
//   - GoAdmin Icon: 图标模块，提供图标定义
//
// 表单类型：
//   - Text: 文本输入框
//   - Number: 数字输入框
//   - Email: 邮箱输入框
//   - Password: 密码输入框
//   - Url: URL 输入框
//   - Ip: IP 地址输入框
//   - Datetime: 日期时间选择器
//   - RichText: 富文本编辑器
//   - Switch: 开关
//   - SelectBox: 下拉框（多选）
//   - SelectSingle: 下拉框（单选）
//   - Select: 下拉框（多选，带后处理）
//   - Radio: 单选框
//   - Multifile: 多文件上传
//   - Currency: 货币输入框
//
// 字段选项：
//   - FieldSortable: 设置字段可排序
//   - FieldEditAble: 设置字段在列表中可编辑
//   - FieldDisplay: 自定义字段显示方式
//   - FieldHide: 隐藏字段（不在列表中显示）
//   - FieldFilterable: 设置字段可筛选
//   - FieldFilterOptions: 设置筛选选项
//   - FieldJoin: 配置表关联
//   - FieldDefault: 设置默认值
//   - FieldHelpMsg: 设置帮助信息
//   - FieldOptions: 设置选项值
//   - FieldOptionExt: 设置扩展选项
//   - FieldOnChooseAjax: 设置 Ajax 联动
//   - FieldOptionInitFn: 设置初始化函数
//   - FieldPostFilterFn: 设置后处理函数
//
// 操作按钮类型：
//   - action.Jump: 跳转按钮，在新标签页打开指定 URL
//   - action.JumpInNewTab: 跳转按钮，在新标签页打开指定 URL（全局按钮）
//   - action.Ajax: Ajax 按钮，发送异步请求
//   - action.PopUp: 弹窗按钮，弹出预览窗口
//
// 筛选器类型：
//   - types.FilterOperatorLike: 模糊匹配
//   - form.SelectSingle: 单选下拉框
//   - form.Radio: 单选框
//   - form.Select: 下拉框
//   - form.DatetimeRange: 日期时间范围
//
// 使用场景：
//   - 后台管理: 为 GoAdmin 管理后台提供数据表定义
//   - 用户管理: 管理用户信息、权限、角色等
//   - 内容管理: 管理文章、新闻、博客等内容
//   - 自定义显示: 自定义字段的显示方式（链接、按钮、图标等）
//   - 富文本编辑: 支持富文本内容编辑，如文章内容、产品描述等
//   - 文件上传: 支持单文件和多文件上传
//   - 字段联动: 通过 Ajax 实现字段之间的联动（如省市联动）
//   - 表关联: 通过表关联查询关联表的数据
//   - 批量操作: 通过选择框实现批量操作
//
// 配置说明：
//   - 表名: 对应数据库中的实际表名
//   - 字段名: 对应数据库表中的列名
//   - 字段类型: 必须与数据库字段类型匹配
//   - 表单类型: 决定表单字段的输入方式
//   - 字段选项: 控制字段的可编辑性、可见性、排序性、筛选性等
//   - 操作按钮: 配置行操作按钮和全局按钮
//   - 筛选器: 配置字段筛选器
//   - 表单分组: 将表单字段分组到不同的标签页
//   - Ajax 联动: 配置字段之间的联动关系
//
// 注意事项：
//   - 需要确保数据库表已正确创建
//   - 字段类型必须与数据库字段类型一致
//   - 表单字段配置必须与列表字段配置对应
//   - ID 字段通常设置为创建时禁用、更新时不可编辑
//   - 富文本编辑器需要正确配置文件上传路径
//   - 自定义字段显示需要正确处理字段值和上下文
//   - Ajax 联动需要正确配置 URL 和处理函数
//   - 表关联需要正确配置关联字段和关联表
//   - 文件上传需要确保上传目录存在且有写入权限
//   - 筛选器选项必须与数据库字段值匹配
//
// 作者: GoAdmin Team
// 创建日期: 2019-01-01
// 版本: 1.0.0
package tables

import (
	"fmt"     // 格式化输出包，提供字符串格式化功能
	"strings" // 字符串处理包，提供字符串分割、连接等功能

	"github.com/purpose168/GoAdmin/context"                              // 上下文包，提供请求上下文管理功能
	"github.com/purpose168/GoAdmin/modules/config"                       // 配置模块包，提供配置管理功能
	"github.com/purpose168/GoAdmin/modules/db"                           // 数据库模块包，提供数据库操作和字段类型定义
	"github.com/purpose168/GoAdmin/plugins/admin/modules/table"          // 表格模块包，提供表模型定义和配置功能
	"github.com/purpose168/GoAdmin/template/icon"                        // 图标包，提供图标定义
	"github.com/purpose168/GoAdmin/template/types"                       // 类型包，提供类型定义和字段模型
	"github.com/purpose168/GoAdmin/template/types/action"                // 动作包，提供按钮和动作配置功能
	"github.com/purpose168/GoAdmin/template/types/form"                  // 表单类型包，提供表单字段类型定义
	selection "github.com/purpose168/GoAdmin/template/types/form/select" // 选择框包，提供选择框类型定义
	editType "github.com/purpose168/GoAdmin/template/types/table"        // 表格类型包，提供表格编辑类型定义
)

// GetUserTable 获取用户表模型
//
// 参数:
//   - ctx: 上下文对象，用于管理请求上下文
//
// 返回:
//   - table.Table: 用户表配置对象，包含列表显示和表单编辑两部分配置
//
// 说明：
//
//	该函数展示了如何创建一个功能完整的用户表模型，包含了各种字段类型、表单类型、操作按钮和筛选器的使用示例。
//	展示了字段关联、自定义显示、Ajax 联动等高级功能。
//
// 功能特性：
//   - 创建自定义表配置，设置各种表属性
//   - 定义表的字段（ID、姓名、性别、经验、饮料、城市、书籍、头像、创建时间、更新时间）
//   - 配置表单字段（包括各种表单类型：文本、数字、邮箱、密码、富文本、选择框等）
//   - 设置表的标题和描述
//   - 配置操作按钮（跳转、Ajax、弹窗等）
//   - 配置筛选器（下拉框、单选框等）
//   - 配置表单分组（将字段分组到不同的标签页）
//   - 配置 Ajax 联动（省市联动）
//   - 配置表关联（书籍关联）
//
// 字段说明：
//   - id: 主键，整数类型，用于唯一标识用户
//   - name: 姓名，字符串类型，用于显示用户姓名
//   - gender: 性别，整数类型，用于标识用户性别（0-男，1-女）
//   - experience: 经验，整数类型，用于标识用户工作经验
//   - drink: 饮料，整数类型，用于标识用户喜欢的饮料
//   - city: 城市，字符串类型，用于标识用户所在城市
//   - name (书籍): 书籍名称，字符串类型，通过表关联查询
//   - avatar: 头像，字符串类型，用于存储用户头像路径
//   - created_at: 创建时间，时间戳类型，用于记录用户创建时间
//   - updated_at: 更新时间，时间戳类型，用于记录用户更新时间
//   - age: 年龄，整数类型，用于标识用户年龄
//   - homepage: 主页，字符串类型，用于存储用户主页 URL
//   - email: 邮箱，字符串类型，用于存储用户邮箱地址
//   - birthday: 生日，字符串类型，用于存储用户生日
//   - password: 密码，字符串类型，用于存储用户密码
//   - ip: IP 地址，字符串类型，用于存储用户 IP 地址
//   - certificate: 证书，字符串类型，用于存储用户证书文件路径
//   - money: 金额，整数类型，用于存储用户金额
//   - resume: 内容，文本类型，用于存储用户简历内容
//   - website: 开关，整数类型，用于标识网站是否开启
//   - fruit: 水果，字符串类型，用于标识用户喜欢的水果
//   - country: 国家，整数类型，用于标识用户所在国家
//
// 配置说明：
//   - 数据库表名: users
//   - 表标题: 用户
//   - 表描述: 用户管理
//   - 筛选表单布局: 三列布局
//   - 默认查询条件: gender=0
//   - ID 字段: 可排序
//   - 姓名字段: 可编辑，可筛选（模糊匹配）
//   - 性别字段: 可编辑（开关），可筛选（单选下拉框），自定义显示
//   - 经验字段: 可筛选（单选框），隐藏
//   - 饮料字段: 可筛选（下拉框），隐藏
//   - 城市字段: 可筛选（文本输入）
//   - 书籍字段: 表关联（user_like_books 表）
//   - 头像字段: 自定义显示
//   - 创建时间字段: 可筛选（日期时间范围）
//   - 更新时间字段: 可编辑（日期时间选择器）
//
// 技术细节：
//   - 使用 table.NewDefaultTable() 创建默认表配置
//   - 使用 table.Config 自定义表配置
//   - 使用 config.GetDatabases().GetDefault().Driver 获取默认数据库驱动
//   - 使用 info.AddField() 添加列表显示字段
//   - 使用 info.FieldSortable() 设置字段可排序
//   - 使用 info.FieldEditAble() 设置字段在列表中可编辑
//   - 使用 info.FieldDisplay() 自定义字段显示
//   - 使用 info.FieldHide() 隐藏字段
//   - 使用 info.FieldFilterable() 设置字段可筛选
//   - 使用 info.FieldFilterOptions() 设置筛选选项
//   - 使用 info.FieldJoin() 配置表关联
//   - 使用 info.AddActionButton() 添加行操作按钮
//   - 使用 info.AddButton() 添加全局按钮
//   - 使用 info.AddSelectBox() 添加选择框
//   - 使用 info.SetFilterFormLayout() 设置筛选表单布局
//   - 使用 info.Where() 设置默认查询条件
//   - 使用 info.SetTable() 设置表名
//   - 使用 info.SetTitle() 设置表标题
//   - 使用 info.SetDescription() 设置表描述
//   - 使用 formList.AddField() 添加表单字段
//   - 使用 formList.FieldDefault() 设置默认值
//   - 使用 formList.FieldHelpMsg() 设置帮助信息
//   - 使用 formList.FieldOptions() 设置选项值
//   - 使用 formList.FieldOptionExt() 设置扩展选项
//   - 使用 formList.FieldDisplay() 自定义显示
//   - 使用 formList.FieldOnChooseAjax() 设置 Ajax 联动
//   - 使用 formList.FieldOptionInitFn() 设置初始化函数
//   - 使用 formList.FieldPostFilterFn() 设置后处理函数
//   - 使用 formList.SetTabGroups() 设置标签页分组
//   - 使用 formList.SetTabHeaders() 设置标签页标题
//   - 使用 formList.SetTable() 设置表名
//   - 使用 formList.SetTitle() 设置表标题
//   - 使用 formList.SetDescription() 设置表描述
//
// 操作按钮说明：
//   - 跳转按钮: 在新标签页打开指定 URL
//   - Ajax 按钮: 发送异步请求，处理函数返回 success、msg、data
//   - 弹窗按钮: 弹出预览窗口，处理函数返回 success、msg、data
//   - 全局跳转按钮: 在页面顶部显示，在新标签页打开指定 URL
//   - 全局弹窗按钮: 在页面顶部显示，弹出预览窗口
//   - 全局 Ajax 按钮: 在页面顶部显示，发送异步请求
//   - 选择框: 用于批量操作，选择后执行指定动作
//
// 筛选器说明：
//   - 模糊匹配: 使用 types.FilterOperatorLike 操作符
//   - 单选下拉框: 使用 form.SelectSingle 表单类型
//   - 单选框: 使用 form.Radio 表单类型
//   - 下拉框: 使用 form.Select 表单类型
//   - 文本输入: 使用默认表单类型
//   - 日期时间范围: 使用 form.DatetimeRange 表单类型
//
// 表单分组说明：
//   - 标签页分组: 使用 SetTabGroups() 方法将字段分组到不同的标签页
//   - 标签页标题: 使用 SetTabHeaders() 方法设置标签页的标题
//   - 分组示例: {"name", "age", "homepage"} 表示这三个字段在同一个标签页
//
// Ajax 联动说明：
//   - FieldOnChooseAjax(): 设置 Ajax 联动
//   - 参数: 目标字段名、Ajax URL、处理函数
//   - 处理函数: 接收上下文对象，返回 success、msg、data
//   - data: 返回选择框的选项数据
//   - 示例: 选择国家后，加载对应的城市选项
//
// 表关联说明：
//   - FieldJoin(): 配置表关联
//   - JoinField: 关联字段（当前表的字段）
//   - Field: 目标字段（关联表的字段）
//   - Table: 关联表名
//   - 示例: 关联 user_like_books 表，查询书籍名称
//
// 使用场景：
//   - 用户管理: 管理用户信息、权限、角色等
//   - 内容管理: 管理文章、新闻、博客等内容
//   - 产品管理: 管理产品信息、分类、库存等
//   - 订单管理: 管理订单信息、状态、物流等
//   - 自定义显示: 将字段显示为链接、按钮、图标等
//   - 富文本编辑: 支持富文本内容编辑
//   - 文件上传: 支持单文件和多文件上传
//   - 字段联动: 通过 Ajax 实现字段之间的联动
//   - 表关联: 通过表关联查询关联表的数据
//   - 批量操作: 通过选择框实现批量操作
//
// 注意事项：
//   - 需要确保数据库表 users 已创建
//   - 字段类型必须与数据库字段类型一致
//   - 表单字段配置必须与列表字段配置对应
//   - ID 字段通常设置为创建时禁用、更新时不可编辑
//   - 富文本编辑器需要正确配置文件上传路径
//   - 自定义字段显示需要正确处理字段值和上下文
//   - Ajax 联动需要正确配置 URL 和处理函数
//   - 表关联需要正确配置关联字段和关联表
//   - 文件上传需要确保上传目录存在且有写入权限
//   - 筛选器选项必须与数据库字段值匹配
//   - 表单分组需要确保所有字段都被分配到标签页
//
// 错误处理：
//   - 如果数据库连接失败，会在初始化时返回错误
//   - 如果表结构不匹配，会在运行时返回错误
//   - 如果文件上传失败，会在运行时返回错误
//   - 如果 Ajax 请求失败，会在运行时返回错误
//   - 如果表关联查询失败，会在运行时返回错误
//
// 示例：
//
//	// 创建上下文
//	ctx := context.NewContext(r)
//
//	// 获取用户表配置
//	userTable := GetUserTable(ctx)
//
//	// 使用表配置
//	info := userTable.GetInfo()
//	formList := userTable.GetForm()
//
//	// 自定义字段显示示例
//	info.AddField("状态", "status", db.Varchar).FieldDisplay(func(value types.FieldModel) interface{} {
//	    if value.Value == "active" {
//	        return `<span class="label label-success">活跃</span>`
//	    }
//	    return `<span class="label label-default">禁用</span>`
//	})
//
//	// Ajax 联动示例
//	formList.AddField("省份", "province", db.Varchar, form.SelectSingle).
//	    FieldOnChooseAjax("city", "/choose/province",
//	        func(ctx *context.Context) (bool, string, interface{}) {
//	            province := ctx.FormValue("value")
//	            // 根据省份查询城市
//	            cities := getCitiesByProvince(province)
//	            return true, "ok", cities
//	        })
//
//	// 表关联示例
//	info.AddField("订单", "order_no", db.Varchar).FieldJoin(types.Join{
//	    JoinField: "user_id",
//	    Field:     "order_no",
//	    Table:     "orders",
//	})
func GetUserTable(ctx *context.Context) (userTable table.Table) {

	// 创建自定义表配置
	// table.Config 参数说明：
	//   - Driver: 数据库驱动类型（从配置文件获取）
	//   - CanAdd: 是否允许添加新记录
	//   - Editable: 是否允许编辑记录
	//   - Deletable: 是否允许删除记录
	//   - Exportable: 是否允许导出数据
	//   - Connection: 数据库连接名称（使用默认连接）
	//   - PrimaryKey: 主键配置（类型和名称）
	userTable = table.NewDefaultTable(ctx, table.Config{
		Driver:     config.GetDatabases().GetDefault().Driver, // 获取默认数据库驱动
		CanAdd:     true,                                      // 允许添加新记录
		Editable:   true,                                      // 允许编辑记录
		Deletable:  true,                                      // 允许删除记录
		Exportable: true,                                      // 允许导出数据
		Connection: table.DefaultConnectionName,               // 使用默认数据库连接
		PrimaryKey: table.PrimaryKey{ // 配置主键
			Type: db.Int,                      // 主键类型为整数
			Name: table.DefaultPrimaryKeyName, // 主键名称为默认名称（id）
		},
	})

	// 获取表信息配置对象
	// SetFilterFormLayout: 设置筛选表单布局（三列布局）
	// Where: 设置默认查询条件（筛选 gender=0 的记录）
	info := userTable.GetInfo().SetFilterFormLayout(form.LayoutThreeCol).Where("gender", "=", 0) // 设置筛选表单布局为三列，默认查询 gender=0 的记录

	// 添加表字段
	// AddField 参数说明：
	//   - 第一个参数：字段显示名称（在列表中显示的标题）
	//   - 第二个参数：数据库字段名（对应数据库表的列名）
	//   - 第三个参数：字段类型（使用 db 包定义的类型常量）
	info.AddField("ID", "id", db.Int).FieldSortable()                     // ID 字段，整数类型，可排序
	info.AddField("姓名", "name", db.Varchar).FieldEditAble(editType.Text). // 姓名字段，字符串类型，列表中可编辑（使用文本输入）
										FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}) // 可筛选（模糊匹配）

	// 添加性别字段，并自定义显示方式
	// FieldDisplay: 自定义字段的显示方式
	// FieldEditAble: 设置字段在列表中可编辑（使用开关类型）
	// FieldEditOptions: 设置编辑选项（图标选项）
	// FieldFilterable: 设置筛选器（单选下拉框）
	// FieldFilterOptions: 设置筛选选项
	info.AddField("性别", "gender", db.Tinyint).FieldDisplay(func(model types.FieldModel) interface{} { // 性别字段，整数类型，自定义显示
		if model.Value == "0" { // 如果值为 0
			return "男" // 返回"男"
		}
		if model.Value == "1" { // 如果值为 1
			return "女" // 返回"女"
		}
		return "未知" // 返回"未知"
	}).FieldEditAble(editType.Switch).FieldEditOptions(types.FieldOptions{ // 列表中可编辑（使用开关），设置编辑选项
		{Value: "0", Text: "👦"}, // 选项 0：男孩图标
		{Value: "1", Text: "👧"}, // 选项 1：女孩图标
	}).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).FieldFilterOptions(types.FieldOptions{ // 可筛选（单选下拉框），设置筛选选项
		{Value: "0", Text: "男"}, // 筛选选项 0：男
		{Value: "1", Text: "女"}, // 筛选选项 1：女
	})

	// 添加经验字段，配置筛选器（单选框）
	// FieldHide: 隐藏该字段（不在列表中显示）
	info.AddField("经验", "experience", db.Tinyint). // 经验字段，整数类型
							FieldFilterable(types.FilterType{FormType: form.Radio}). // 可筛选（单选框）
							FieldFilterOptions(types.FieldOptions{                   // 设置筛选选项
			{Value: "0", Text: "一年"}, // 筛选选项 0：一年
			{Value: "1", Text: "两年"}, // 筛选选项 1：两年
			{Value: "3", Text: "三年"}, // 筛选选项 3：三年
		}).FieldHide() // 隐藏该字段（不在列表中显示）

	// 添加饮料字段，配置筛选器（下拉框）
	info.AddField("饮料", "drink", db.Tinyint). // 饮料字段，整数类型
							FieldFilterable(types.FilterType{FormType: form.Select}). // 可筛选（下拉框）
							FieldFilterOptions(types.FieldOptions{                    // 设置筛选选项
			{Value: "water", Text: "水"},     // 筛选选项：水
			{Value: "juice", Text: "果汁"},    // 筛选选项：果汁
			{Value: "red bull", Text: "红牛"}, // 筛选选项：红牛
		}).FieldHide() // 隐藏该字段（不在列表中显示）

	// 添加城市字段，配置筛选器（文本输入）
	info.AddField("城市", "city", db.Varchar).FieldFilterable() // 城市字段，字符串类型，可筛选（文本输入）

	// 添加书籍字段，使用表关联
	// FieldJoin: 配置表关联
	//   - JoinField: 关联字段（当前表的字段）
	//   - Field: 目标字段（关联表的字段）
	//   - Table: 关联表名
	info.AddField("书籍", "name", db.Varchar).FieldJoin(types.Join{ // 书籍字段，字符串类型，表关联
		JoinField: "user_id",         // 关联字段：user_id（当前表的字段）
		Field:     "id",              // 目标字段：id（关联表的字段）
		Table:     "user_like_books", // 关联表名：user_like_books
	})

	// 添加头像字段，自定义显示
	info.AddField("头像", "avatar", db.Varchar).FieldDisplay(func(value types.FieldModel) interface{} { // 头像字段，字符串类型，自定义显示
		return "1231" // 返回固定值（示例）
	})

	// 添加创建时间字段，配置筛选器（日期时间范围）
	info.AddField("创建时间", "created_at", db.Timestamp). // 创建时间字段，时间戳类型
								FieldFilterable(types.FilterType{FormType: form.DatetimeRange}) // 可筛选（日期时间范围）

	// 添加更新时间字段，设置为可编辑（日期时间选择器）
	info.AddField("更新时间", "updated_at", db.Timestamp).FieldEditAble(editType.Datetime) // 更新时间字段，时间戳类型，列表中可编辑（使用日期时间选择器）

	// ===========================
	// Buttons - 按钮配置部分
	// ===========================

	// 添加跳转按钮（在新标签页打开）
	// AddActionButton: 添加行操作按钮（在每行数据后面显示）
	// 参数：
	//   - ctx: 上下文对象
	//   - 按钮ID: 用于标识按钮
	//   - action.Jump: 跳转动作，指定 URL
	info.AddActionButton(ctx, "google", action.Jump("https://google.com")) // 添加跳转按钮，跳转到 Google

	// 添加 Ajax 按钮（异步请求）
	// action.Ajax: Ajax 动作，指定 URL 和处理函数
	// 处理函数参数：
	//   - ctx: 上下文对象
	// 返回值：
	//   - success: 是否成功
	//   - msg: 消息内容
	//   - data: 返回的数据
	info.AddActionButton(ctx, "审核", action.Ajax("/admin/audit", // 添加 Ajax 按钮，发送审核请求
		func(ctx *context.Context) (success bool, msg string, data interface{}) { // 处理函数
			fmt.Println("PostForm", ctx.PostForm()) // 打印表单数据
			return true, "success", ""              // 返回成功
		}))

	// 添加弹窗按钮（弹出预览窗口）
	// action.PopUp: 弹窗动作，指定 URL、标题和处理函数
	info.AddActionButton(ctx, "预览", action.PopUp("/admin/preview", "预览", // 添加弹窗按钮，弹出预览窗口
		func(ctx *context.Context) (success bool, msg string, data interface{}) { // 处理函数
			return true, "", "<h2>预览内容</h2>" // 返回预览内容
		}))

	// 添加全局跳转按钮（在新标签页打开）
	// AddButton: 添加全局按钮（在页面顶部显示）
	// 参数：
	//   - ctx: 上下文对象
	//   - 按钮ID: 用于标识按钮
	//   - icon: 按钮图标
	//   - action.JumpInNewTab: 在新标签页跳转，指定 URL 和文本
	info.AddButton(ctx, "jump", icon.User, action.JumpInNewTab("/admin/info/authors", "作者")) // 添加全局跳转按钮，跳转到作者页面

	// 添加全局弹窗按钮
	// action.PopUp: 弹窗动作，指定 URL、标题和处理函数
	info.AddButton(ctx, "popup", icon.Terminal, action.PopUp("/admin/popup", "弹窗示例", // 添加全局弹窗按钮，弹出弹窗示例
		func(ctx *context.Context) (success bool, msg string, data interface{}) { // 处理函数
			return true, "", "<h2>你好世界</h2>" // 返回弹窗内容
		}))

	// 添加全局 Ajax 按钮
	// action.Ajax: Ajax 动作，指定 URL 和处理函数
	info.AddButton(ctx, "ajax", icon.Android, action.Ajax("/admin/ajax", // 添加全局 Ajax 按钮，发送 Ajax 请求
		func(ctx *context.Context) (success bool, msg string, data interface{}) { // 处理函数
			return true, "哦，我收到了", "" // 返回成功消息
		}))

	// 添加选择框（用于批量操作）
	// AddSelectBox: 添加选择框，指定字段名、选项和动作
	// action.FieldFilter: 字段筛选动作
	info.AddSelectBox(ctx, "gender", types.FieldOptions{ // 添加选择框，用于批量操作性别字段
		{Value: "0", Text: "男"}, // 选项 0：男
		{Value: "1", Text: "女"}, // 选项 1：女
	}, action.FieldFilter("gender")) // 动作：字段筛选

	// 设置表的基本信息
	info.SetTable("users").SetTitle("用户").SetDescription("用户管理") // 设置表名为 users，标题为"用户"，描述为"用户管理"

	// 获取表单配置对象
	// 该对象用于配置表的表单编辑部分（创建和编辑表单）
	formList := userTable.GetForm() // 获取表单配置对象

	// 添加表单字段
	// AddField 参数说明：
	//   - 第一个参数：字段显示名称（在表单中显示的标签）
	//   - 第二个参数：数据库字段名（对应数据库表的列名）
	//   - 第三个参数：字段类型（使用 db 包定义的类型常量）
	//   - 第四个参数：表单类型（使用 form 包定义的类型常量）

	// 添加姓名字段（文本输入）
	formList.AddField("姓名", "name", db.Varchar, form.Text) // 姓名字段，字符串类型，文本输入

	// 添加年龄字段（数字输入）
	formList.AddField("年龄", "age", db.Int, form.Number) // 年龄字段，整数类型，数字输入

	// 添加主页字段（URL 输入）
	// FieldDefault: 设置默认值
	formList.AddField("主页", "homepage", db.Varchar, form.Url).FieldDefault("http://google.com") // 主页字段，字符串类型，URL 输入，默认值为 http://google.com

	// 添加邮箱字段（邮箱输入）
	formList.AddField("邮箱", "email", db.Varchar, form.Email).FieldDefault("xxxx@xxx.com") // 邮箱字段，字符串类型，邮箱输入，默认值为 xxxx@xxx.com

	// 添加生日字段（日期时间选择器）
	formList.AddField("生日", "birthday", db.Varchar, form.Datetime).FieldDefault("2010-09-05") // 生日字段，字符串类型，日期时间选择器，默认值为 2010-09-05

	// 添加密码字段（密码输入）
	formList.AddField("密码", "password", db.Varchar, form.Password) // 密码字段，字符串类型，密码输入

	// 添加 IP 字段（IP 地址输入）
	formList.AddField("IP", "ip", db.Varchar, form.Ip) // IP 字段，字符串类型，IP 地址输入

	// 添加证书字段（多文件上传）
	// form.Multifile: 多文件上传类型
	// FieldOptionExt: 设置扩展选项（最大文件数量）
	formList.AddField("证书", "certificate", db.Varchar, form.Multifile).FieldOptionExt(map[string]interface{}{ // 证书字段，字符串类型，多文件上传
		"maxFileCount": 10, // 扩展选项：最大文件数量为 10
	})

	// 添加金额字段（货币输入）
	formList.AddField("金额", "money", db.Int, form.Currency) // 金额字段，整数类型，货币输入

	// 添加内容字段（富文本编辑器）
	// form.RichText: 富文本编辑器类型
	formList.AddField("内容", "resume", db.Text, form.RichText). // 内容字段，文本类型，富文本编辑器
									FieldDefault(`<h1>343434</h1><p>34344433434</p><ol><li>23234</li><li>2342342342</li><li>asdfads</li></ol><ul><li>3434334</li><li>34343343434</li><li>44455</li></ul><p><span style="color: rgb(194, 79, 74);">343434</span></p><p><span style="background-color: rgb(194, 79, 74); color: rgb(0, 0, 0);">434434433434</span></p><table border="0" width="100%" cellpadding="0" cellspacing="0"><tbody><tr><td>&nbsp;</td><td>&nbsp;</td><td>&nbsp;</td></tr><tr><td>&nbsp;</td><td>&nbsp;</td><td>&nbsp;</td></tr><tr><td>&nbsp;</td><td>&nbsp;</td><td>&nbsp;</td></tr><tr><td>&nbsp;</td><td>&nbsp;</td><td>&nbsp;</td></tr></tbody></table><p><br></p><p><span style="color: rgb(194, 79, 74);"><br></span></p>`) // 默认值为富文本内容

	// 添加开关字段（网站开关）
	// form.Switch: 开关类型
	// FieldHelpMsg: 设置帮助信息
	// FieldOptions: 设置选项值
	formList.AddField("开关", "website", db.Tinyint, form.Switch). // 开关字段，整数类型，开关
									FieldHelpMsg("当网站关闭时将无法访问").     // 帮助信息：当网站关闭时将无法访问
									FieldOptions(types.FieldOptions{ // 设置选项值
			{Value: "0"}, // 选项 0
			{Value: "1"}, // 选项 1
		})

	// 添加水果字段（下拉框）
	// form.SelectBox: 下拉框类型
	// FieldOptions: 设置选项（文本和值）
	// FieldDisplay: 自定义显示方式
	formList.AddField("水果", "fruit", db.Varchar, form.SelectBox). // 水果字段，字符串类型，下拉框
									FieldOptions(types.FieldOptions{ // 设置选项
			{Text: "苹果", Value: "apple"},      // 选项：苹果
			{Text: "香蕉", Value: "banana"},     // 选项：香蕉
			{Text: "西瓜", Value: "watermelon"}, // 选项：西瓜
			{Text: "梨", Value: "pear"},        // 选项：梨
		}).
		FieldDisplay(func(value types.FieldModel) interface{} { // 自定义显示方式
			return []string{"梨"} // 返回固定值（示例）
		})

	// 添加国家字段（单选下拉框）
	// form.SelectSingle: 单选下拉框类型
	// FieldDefault: 设置默认值
	// FieldOnChooseAjax: 设置 Ajax 联动（当选择国家时，加载对应的城市）
	//   参数：目标字段名、Ajax URL、处理函数
	formList.AddField("国家", "country", db.Tinyint, form.SelectSingle). // 国家字段，整数类型，单选下拉框
										FieldOptions(types.FieldOptions{ // 设置选项
			{Text: "中国", Value: "china"},   // 选项：中国
			{Text: "美国", Value: "america"}, // 选项：美国
			{Text: "英国", Value: "england"}, // 选项：英国
			{Text: "加拿大", Value: "canada"}, // 选项：加拿大
		}).FieldDefault("china").FieldOnChooseAjax("city", "/choose/country", // 默认值为 china，设置 Ajax 联动（当选择国家时，加载对应的城市）
		func(ctx *context.Context) (bool, string, interface{}) { // 处理函数
			country := ctx.FormValue("value")     // 获取选择的国家值
			var data = make(selection.Options, 0) // 创建选项数据
			switch country {                      // 根据国家选择对应的城市
			case "china": // 如果是中国
				data = selection.Options{ // 中国的城市选项
					{Text: "北京", ID: "beijing"},   // 北京
					{Text: "上海", ID: "shanghai"},  // 上海
					{Text: "广州", ID: "guangzhou"}, // 广州
					{Text: "深圳", ID: "shenzhen"},  // 深圳
				}
			case "america": // 如果是美国
				data = selection.Options{ // 美国的城市选项
					{Text: "洛杉矶", ID: "los angeles"},      // 洛杉矶
					{Text: "华盛顿特区", ID: "washington, dc"}, // 华盛顿特区
					{Text: "纽约", ID: "new york"},          // 纽约
					{Text: "拉斯维加斯", ID: "las vegas"},      // 拉斯维加斯
				}
			case "england": // 如果是英国
				data = selection.Options{ // 英国的城市选项
					{Text: "伦敦", ID: "london"},       // 伦敦
					{Text: "剑桥", ID: "cambridge"},    // 剑桥
					{Text: "曼彻斯特", ID: "manchester"}, // 曼彻斯特
					{Text: "利物浦", ID: "liverpool"},   // 利物浦
				}
			case "canada": // 如果是加拿大
				data = selection.Options{ // 加拿大的城市选项
					{Text: "温哥华", ID: "vancouver"}, // 温哥华
					{Text: "多伦多", ID: "toronto"},   // 多伦多
				}
			default: // 默认情况
				data = selection.Options{ // 默认的城市选项
					{Text: "北京", ID: "beijing"},   // 北京
					{Text: "上海", ID: "shangHai"},  // 上海
					{Text: "广州", ID: "guangzhou"}, // 广州
					{Text: "深圳", ID: "shenZhen"},  // 深圳
				}
			}
			return true, "ok", data // 返回成功和选项数据
		})

	// 添加城市字段（单选下拉框）
	// FieldOptionInitFn: 设置初始化函数（根据当前值设置选项）
	formList.AddField("城市", "city", db.Varchar, form.SelectSingle). // 城市字段，字符串类型，单选下拉框
									FieldOptionInitFn(func(val types.FieldModel) types.FieldOptions { // 设置初始化函数
			return types.FieldOptions{ // 返回选项
				{Value: val.Value, Text: val.Value, Selected: true}, // 根据当前值设置选项，并选中
			}
		}).FieldOptions(types.FieldOptions{ // 设置默认选项
		{Text: "北京", Value: "beijing"},   // 北京
		{Text: "上海", Value: "shanghai"},  // 上海
		{Text: "广州", Value: "guangzhou"}, // 广州
		{Text: "深圳", Value: "shenzhen"},  // 深圳
	})

	// 添加性别字段（单选框）
	// form.Radio: 单选框类型
	formList.AddField("性别", "gender", db.Tinyint, form.Radio). // 性别字段，整数类型，单选框
									FieldOptions(types.FieldOptions{ // 设置选项
			{Text: "男孩", Value: "0"}, // 选项 0：男孩
			{Text: "女孩", Value: "1"}, // 选项 1：女孩
		})

	// 添加饮料字段（下拉框）
	// FieldDisplay: 自定义显示方式（将逗号分隔的字符串转换为数组）
	// FieldPostFilterFn: 设置后处理函数（将数组转换为逗号分隔的字符串）
	formList.AddField("饮料", "drink", db.Varchar, form.Select). // 饮料字段，字符串类型，下拉框
									FieldOptions(types.FieldOptions{ // 设置选项
			{Text: "啤酒", Value: "beer"},     // 选项：啤酒
			{Text: "果汁", Value: "juice"},    // 选项：果汁
			{Text: "水", Value: "water"},     // 选项：水
			{Text: "红牛", Value: "red bull"}, // 选项：红牛
		}).
		FieldDefault("beer").                                   // 默认值为 beer
		FieldDisplay(func(value types.FieldModel) interface{} { // 自定义显示方式（将逗号分隔的字符串转换为数组）
			return strings.Split(value.Value, ",") // 将逗号分隔的字符串转换为数组
		}).
		FieldPostFilterFn(func(value types.PostFieldModel) interface{} { // 设置后处理函数（将数组转换为逗号分隔的字符串）
			return strings.Join(value.Value, ",") // 将数组转换为逗号分隔的字符串
		})

	// 添加工作经验字段（单选下拉框）
	formList.AddField("工作经验", "experience", db.Tinyint, form.SelectSingle). // 工作经验字段，整数类型，单选下拉框
										FieldOptions(types.FieldOptions{ // 设置选项
			{Text: "两年", Value: "0"}, // 选项 0：两年
			{Text: "三年", Value: "1"}, // 选项 1：三年
			{Text: "四年", Value: "2"}, // 选项 2：四年
			{Text: "五年", Value: "3"}, // 选项 3：五年
		}).FieldDefault("beer") // 默认值为 beer（示例）

	// 设置表单分组（将字段分组到不同的标签页）
	// SetTabGroups: 设置标签页分组，每个数组元素是一个标签页，包含多个字段名
	formList.SetTabGroups(types.TabGroups{ // 设置标签页分组
		{"name", "age", "homepage", "email", "birthday", "password", "ip", "certificate", "money", "resume"}, // 第一个标签页：输入
		{"website", "fruit", "country", "city", "gender", "drink", "experience"},                             // 第二个标签页：选择
	})

	// 设置标签页标题
	// SetTabHeaders: 设置标签页的标题
	formList.SetTabHeaders("输入", "选择") // 设置标签页标题为"输入"和"选择"

	// 设置表单的基本信息
	// SetTable: 设置数据库表名
	// SetTitle: 设置表单在界面中显示的标题
	// SetDescription: 设置表单的描述信息
	formList.SetTable("users").SetTitle("用户").SetDescription("用户管理") // 设置表名为 users，标题为"用户"，描述为"用户管理"

	return // 返回用户表配置对象
}

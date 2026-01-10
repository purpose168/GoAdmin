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
//   - 自定义字段显示: 通过 FieldDisplay 方法自定义字段的显示方式
//   - 富文本编辑器: 支持富文本编辑的表单字段类型
//   - 文件上传: 在富文本编辑器中插入图片等文件
//   - 列表编辑: 在列表页面直接编辑字段内容
//   - 字段类型: 定义数据库字段的数据类型（Int、Varchar、Date、Timestamp 等）
//   - 表单类型: 定义表单字段的输入类型（Text、RichText、Datetime、Select、Number 等）
//   - 字段选项: 配置字段的特殊行为（可排序、可编辑、禁用、可上传文件等）
//
// 技术栈：
//   - GoAdmin Table: 表格模块，提供表模型定义功能
//   - GoAdmin Form: 表单模块，提供表单配置功能
//   - GoAdmin Template: 模板模块，提供模板渲染和组件功能
//   - GoAdmin Context: 上下文模块，提供请求上下文管理
//   - GoAdmin DB: 数据库模块，提供数据库操作功能
//   - GoAdmin Types: 类型模块，提供类型定义和字段模型
//
// 表单类型：
//   - Default: 默认显示类型
//   - Text: 文本输入框
//   - Textarea: 文本域
//   - RichText: 富文本编辑器
//   - Datetime: 日期时间选择器
//   - Select: 下拉选择框
//   - Number: 数字输入框
//   - Password: 密码输入框
//   - File: 文件上传
//
// 使用场景：
//   - 后台管理: 为 GoAdmin 管理后台提供数据表定义
//   - 内容管理: 管理文章、新闻、博客等内容
//   - 自定义显示: 自定义字段的显示方式（链接、按钮等）
//   - 富文本编辑: 支持富文本内容编辑，如文章内容、产品描述等
//   - 文件管理: 在富文本中插入图片、视频等文件
//   - 列表编辑: 在列表页面直接编辑字段内容，提高效率
//
// 配置说明：
//   - 表名: 对应数据库中的实际表名
//   - 字段名: 对应数据库表中的列名
//   - 字段类型: 必须与数据库字段类型匹配
//   - 表单类型: 决定表单字段的输入方式
//   - 字段选项: 控制字段的可编辑性、可见性、排序性等
//   - 自定义显示: 通过 FieldDisplay 方法自定义字段显示
//   - 文件上传: 通过 FieldEnableFileUpload 启用文件上传功能
//
// 注意事项：
//   - 需要确保数据库表已正确创建
//   - 字段类型必须与数据库字段类型一致
//   - 表单字段配置必须与列表字段配置对应
//   - ID 字段通常设置为创建时禁用、更新时不可编辑
//   - 富文本编辑器需要正确配置文件上传路径
//   - 自定义字段显示需要正确处理字段值和上下文
//   - 列表编辑功能需要谨慎使用，避免误操作
//
// 作者: GoAdmin Team
// 创建日期: 2019-01-01
// 版本: 1.0.0
package tables

import (
	"github.com/purpose168/GoAdmin/context"                       // 上下文包，提供请求上下文管理功能
	"github.com/purpose168/GoAdmin/modules/db"                    // 数据库模块包，提供数据库操作和字段类型定义
	"github.com/purpose168/GoAdmin/plugins/admin/modules/table"   // 表格模块包，提供表模型定义和配置功能
	"github.com/purpose168/GoAdmin/template"                      // 模板包，提供模板渲染和组件功能
	"github.com/purpose168/GoAdmin/template/types"                // 类型包，提供类型定义和字段模型
	"github.com/purpose168/GoAdmin/template/types/form"           // 表单类型包，提供表单字段类型定义
	editType "github.com/purpose168/GoAdmin/template/types/table" // 表格类型包，提供表格编辑类型定义
)

// GetPostsTable 获取文章表模型
//
// 参数:
//   - ctx: 上下文对象，用于管理请求上下文
//
// 返回:
//   - table.Table: 文章表配置对象，包含列表显示和表单编辑两部分配置
//
// 说明：
//
//	该函数展示了如何创建文章表模型，包括自定义字段显示和富文本编辑器的使用。
//	FieldDisplay 方法可以自定义字段的显示方式，例如显示为链接、按钮等。
//	RichText 表单类型提供了富文本编辑功能，支持文件上传。
//
// 功能特性：
//   - 创建默认表配置
//   - 定义表的字段（ID、标题、作者ID、描述、内容、日期）
//   - 配置列表显示字段
//   - 配置表单编辑字段
//   - 设置表的标题和描述
//   - 演示自定义字段显示（作者ID显示为链接）
//   - 演示富文本编辑器的使用
//   - 演示列表编辑功能（内容字段可编辑）
//   - 演示文件上传功能（富文本编辑器支持文件上传）
//
// 字段说明：
//   - id: 主键，整数类型，用于唯一标识文章
//   - title: 标题，字符串类型，用于显示文章标题
//   - author_id: 作者ID，字符串类型，关联作者表，显示为可点击的链接
//   - description: 描述，字符串类型，用于显示文章摘要
//   - content: 内容，字符串类型，用于存储文章正文，使用富文本编辑器
//   - date: 日期，字符串类型，用于记录文章发布日期
//
// 配置说明：
//   - 数据库表名: posts
//   - 表标题: 文章
//   - 表描述: 文章管理
//   - ID 字段: 可排序，创建时禁用，更新时不可编辑
//   - 标题字段: 可编辑，文本输入
//   - 作者ID字段: 自定义显示为链接
//   - 描述字段: 可编辑，文本输入
//   - 内容字段: 列表可编辑（文本域），表单使用富文本编辑器，支持文件上传
//   - 日期字段: 可编辑，日期时间选择器
//
// 技术细节：
//   - 使用 table.NewDefaultTable() 创建默认表配置
//   - 使用 table.DefaultConfig() 使用默认数据库配置
//   - 使用 info.AddField() 添加列表显示字段
//   - 使用 info.FieldDisplay() 自定义字段显示
//   - 使用 info.FieldEditAble() 设置字段在列表中可编辑
//   - 使用 info.SetTable() 设置表名
//   - 使用 info.SetTitle() 设置表标题
//   - 使用 info.SetDescription() 设置表描述
//   - 使用 formList.AddField() 添加表单字段
//   - 使用 formList.FieldDisplayButCanNotEditWhenUpdate() 设置字段更新时不可编辑
//   - 使用 formList.FieldDisableWhenCreate() 设置字段创建时禁用
//   - 使用 formList.FieldEnableFileUpload() 启用文件上传功能
//
// 自定义字段显示说明：
//   - FieldDisplay 方法参数: types.FieldModel，包含字段值和上下文信息
//   - FieldDisplay 方法返回值: interface{}，自定义的显示内容
//   - 使用 template.Default(ctx).Link() 创建链接组件
//   - SetURL() 设置链接地址
//   - SetContent() 设置链接文本
//   - OpenInNewTab() 在新标签页打开链接
//   - SetTabTitle() 设置新标签页标题
//   - GetContent() 获取链接的 HTML 内容
//
// 富文本编辑器说明：
//   - form.RichText: 富文本编辑器类型
//   - FieldEnableFileUpload(): 启用文件上传功能
//   - 支持在富文本中插入图片、视频等文件
//   - 文件上传路径需要在配置中设置
//   - 支持常见的富文本编辑功能（加粗、斜体、列表、链接等）
//
// 列表编辑说明：
//   - FieldEditAble(editType.Textarea): 设置字段在列表中可编辑
//   - editType.Textarea: 使用文本域编辑器
//   - 允许用户在列表页面直接编辑字段内容
//   - 编辑后自动保存到数据库
//   - 适用于需要快速编辑的字段
//
// 使用场景：
//   - 内容管理系统: 管理文章、新闻、博客等内容
//   - 博客系统: 管理博客文章
//   - 新闻网站: 管理新闻内容
//   - 产品管理: 管理产品描述和详情
//   - 文档管理: 管理文档内容
//   - 自定义显示: 将字段显示为链接、按钮等
//   - 富文本编辑: 支持富文本内容编辑
//   - 文件上传: 在富文本中插入图片、视频等文件
//
// 注意事项：
//   - 需要确保数据库表 posts 已创建
//   - 字段类型必须与数据库字段类型一致
//   - ID 字段通常设置为创建时禁用、更新时不可编辑
//   - 富文本编辑器需要正确配置文件上传路径
//   - 自定义字段显示需要正确处理字段值和上下文
//   - 列表编辑功能需要谨慎使用，避免误操作
//   - 文件上传需要确保上传目录存在且有写入权限
//
// 错误处理：
//   - 如果数据库连接失败，会在初始化时返回错误
//   - 如果表结构不匹配，会在运行时返回错误
//   - 如果文件上传失败，会在运行时返回错误
//
// 示例：
//
//	// 创建上下文
//	ctx := context.NewContext(r)
//
//	// 获取文章表配置
//	postsTable := GetPostsTable(ctx)
//
//	// 使用表配置
//	info := postsTable.GetInfo()
//	formList := postsTable.GetForm()
//
//	// 自定义字段显示示例
//	info.AddField("状态", "status", db.Varchar).FieldDisplay(func(value types.FieldModel) interface{} {
//	    if value.Value == "published" {
//	        return `<span class="label label-success">已发布</span>`
//	    }
//	    return `<span class="label label-default">草稿</span>`
//	})
//
//	// 列表编辑示例
//	info.AddField("标题", "title", db.Varchar).FieldEditAble(editType.Text)
func GetPostsTable(ctx *context.Context) (postsTable table.Table) {

	postsTable = table.NewDefaultTable(ctx, table.DefaultConfig()) // 创建默认表配置，使用默认数据库驱动和连接

	info := postsTable.GetInfo() // 获取表信息配置对象，用于配置表的列表显示部分

	// 添加表字段
	// AddField 参数说明：
	//   - 第一个参数：字段显示名称（在列表中显示的标题）
	//   - 第二个参数：数据库字段名（对应数据库表的列名）
	//   - 第三个参数：字段类型（使用 db 包定义的类型常量）
	info.AddField("ID", "id", db.Int).FieldSortable() // ID 字段，整数类型，可排序
	info.AddField("标题", "title", db.Varchar)          // 标题字段，字符串类型

	// 添加作者ID字段，并自定义显示方式
	// FieldDisplay: 自定义字段的显示方式
	// 参数：types.FieldModel - 包含字段值和上下文信息
	// 返回值：interface{} - 自定义的显示内容
	// 这里将作者ID显示为一个可点击的链接，点击后在新标签页打开作者详情
	info.AddField("作者ID", "author_id", db.Varchar).FieldDisplay(func(value types.FieldModel) interface{} {
		return template.Default(ctx).
			Link().
			SetURL("/admin/info/authors/detail?__goadmin_detail_pk=100"). // 设置链接地址
			SetContent("100").                                            // 设置链接文本
			OpenInNewTab().                                               // 在新标签页打开链接
			SetTabTitle("作者详情").                                          // 设置新标签页标题
			GetContent()                                                  // 获取链接的 HTML 内容
	})

	info.AddField("描述", "description", db.Varchar) // 描述字段，字符串类型

	// 添加内容字段，并设置为可编辑的文本域
	// FieldEditAble: 设置字段在列表中可编辑
	// 参数：editType.Textarea - 使用文本域编辑器
	info.AddField("内容", "content", db.Varchar).FieldEditAble(editType.Textarea) // 内容字段，字符串类型，列表中可编辑（使用文本域）

	info.AddField("日期", "date", db.Varchar) // 日期字段，字符串类型

	// 设置表的基本信息
	// SetTable: 设置数据库表名
	// SetTitle: 设置表在界面中显示的标题
	// SetDescription: 设置表的描述信息
	info.SetTable("posts").SetTitle("文章").SetDescription("文章管理") // 设置表名为 posts，标题为"文章"，描述为"文章管理"

	formList := postsTable.GetForm() // 获取表单配置对象，用于配置表的表单编辑部分（创建和编辑表单）

	// 添加表单字段
	// AddField 参数说明：
	//   - 第一个参数：字段显示名称（在表单中显示的标签）
	//   - 第二个参数：数据库字段名（对应数据库表的列名）
	//   - 第三个参数：字段类型（使用 db 包定义的类型常量）
	//   - 第四个参数：表单类型（使用 form 包定义的类型常量）
	formList.AddField("ID", "id", db.Int, form.Default).FieldDisplayButCanNotEditWhenUpdate().FieldDisableWhenCreate() // ID 字段：默认显示，更新时不可编辑，创建时禁用
	formList.AddField("标题", "title", db.Varchar, form.Text)                                                            // 标题字段：文本输入
	formList.AddField("描述", "description", db.Varchar, form.Text)                                                      // 描述字段：文本输入

	// 添加内容字段，使用富文本编辑器
	// form.RichText: 富文本编辑器类型
	// FieldEnableFileUpload: 启用文件上传功能，允许在富文本中插入图片等文件
	formList.AddField("内容", "content", db.Varchar, form.RichText).FieldEnableFileUpload() // 内容字段：富文本编辑器，启用文件上传

	// 添加日期字段，使用日期时间选择器
	formList.AddField("日期", "date", db.Varchar, form.Datetime) // 日期字段：日期时间选择器

	// 设置表单的基本信息
	// SetTable: 设置数据库表名
	// SetTitle: 设置表单在界面中显示的标题
	// SetDescription: 设置表单的描述信息
	formList.SetTable("posts").SetTitle("文章").SetDescription("文章管理") // 设置表名为 posts，标题为"文章"，描述为"文章管理"

	return // 返回文章表配置对象
}

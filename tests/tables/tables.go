// Package tables 提供数据库表模型的定义和配置
//
// 本包实现了 GoAdmin 管理后台中使用的各种数据表结构的定义和配置，提供以下功能：
//   - GetAuthorsTable: 获取作者表模型配置
//   - GetUserTable: 获取用户表模型配置
//   - GetExternalTable: 获取外部数据源表模型配置
//   - GetPostsTable: 获取文章表模型配置
//   - Generators: 表生成器映射，统一管理所有表模型的注册
//
// 核心概念：
//   - 表模型: 定义数据库表的结构和行为
//   - 表生成器: 用于创建表配置对象的函数
//   - 表生成器映射: 将表名与表生成函数关联起来的映射表
//   - 表注册: 通过表生成器映射注册所有需要管理的表
//   - 表查找: 根据表名查找对应的表生成函数
//   - 表创建: 调用表生成函数创建表配置对象
//
// 技术栈：
//   - GoAdmin Table: 表格模块，提供表模型定义和配置功能
//   - GoAdmin Context: 上下文模块，提供请求上下文管理
//   - GoAdmin Generator: 生成器模块，提供表生成器类型定义
//
// 使用场景：
//   - 后台管理: 为 GoAdmin 管理后台提供数据表定义
//   - 表注册: 在应用启动时注册所有需要管理的表
//   - 表查找: 根据表名查找对应的表生成函数
//   - 表创建: 动态创建表配置对象
//   - 表管理: 统一管理所有表模型的注册和创建
//
// 配置说明：
//   - 表名: 对应数据库中的实际表名
//   - 表生成函数: 用于创建表配置对象的函数
//   - 表生成器映射: 将表名与表生成函数关联起来的映射表
//   - 表注册: 通过表生成器映射注册所有需要管理的表
//
// 注意事项：
//   - 需要确保数据库表已正确创建
//   - 表名必须与数据库表名一致
//   - 表生成函数必须返回正确的表配置对象
//   - 表生成器映射中的键值对必须一一对应
//   - 表生成函数必须接收上下文对象作为参数
//   - 表生成函数必须返回 table.Table 类型的对象
//
// 作者: GoAdmin Team
// 创建日期: 2019-01-01
// 版本: 1.0.0
package tables

import "github.com/purpose168/GoAdmin/plugins/admin/modules/table" // 表格模块包，提供表模型定义和配置功能

// Generators 表生成器映射
//
// 类型: map[string]table.Generator
//
// 说明：
//
//	该映射将表名（字符串）与对应的表生成函数（table.Generator）关联起来
//	GoAdmin 框架通过这个映射来查找和创建对应的表模型
//
// table.Generator 类型定义：
//
//	type Generator func(ctx *context.Context) Table
//
//	这是一个函数类型，接收上下文对象，返回表配置对象
//
// 功能特性：
//   - 统一管理所有表模型的注册
//   - 支持动态添加和删除表
//   - 支持根据表名快速查找表生成函数
//   - 支持批量注册表
//   - 支持表生成函数的延迟加载
//
// 已注册的表：
//   - "posts": 文章表，使用 GetPostsTable 函数生成
//   - "authors": 作者表，使用 GetAuthorsTable 函数生成
//   - "external": 外部数据源表，使用 GetExternalTable 函数生成
//
// 使用场景：
//  1. 应用启动时，通过这个映射注册所有的表
//  2. 当用户访问某个表时，框架根据表名查找对应的生成函数
//  3. 调用生成函数创建表配置对象
//  4. 在管理插件中使用这个映射初始化所有表
//
// 技术细节：
//   - 使用 map[string]table.Generator 类型定义表生成器映射
//   - 键（key）: 表名（字符串类型），对应数据库表名
//   - 值（value）: 表生成函数（table.Generator 类型），用于创建表配置对象
//   - 表生成函数接收上下文对象作为参数
//   - 表生成函数返回 table.Table 类型的对象
//   - 映射在包初始化时创建
//   - 映射可以在运行时动态修改
//
// 表生成函数签名：
//
//	func GetPostsTable(ctx *context.Context) table.Table {
//	    // 创建表配置对象
//	    postsTable := table.NewDefaultTable(ctx, table.DefaultConfig())
//	    // 配置表信息
//	    info := postsTable.GetInfo()
//	    // 配置表单
//	    formList := postsTable.GetForm()
//	    // 返回表配置对象
//	    return postsTable
//	}
//
// 使用示例：
//
//	// 示例 1: 查找表生成函数
//	if generator, ok := tables.Generators["posts"]; ok {
//	    // 找到表生成函数
//	    postsTable := generator(ctx)
//	    // 使用表配置对象
//	    info := postsTable.GetInfo()
//	    formList := postsTable.GetForm()
//	}
//
//	// 示例 2: 遍历所有表生成函数
//	for tableName, generator := range tables.Generators {
//	    // 获取表配置对象
//	    tableConfig := generator(ctx)
//	    // 使用表配置对象
//	    fmt.Printf("表名: %s\n", tableName)
//	    fmt.Printf("表标题: %s\n", tableConfig.GetInfo().GetTitle())
//	}
//
//	// 示例 3: 动态添加表
//	tables.Generators["new_table"] = func(ctx *context.Context) table.Table {
//	    // 创建新表配置对象
//	    newTable := table.NewDefaultTable(ctx, table.DefaultConfig())
//	    // 配置表信息
//	    info := newTable.GetInfo()
//	    info.SetTable("new_table").SetTitle("新表").SetDescription("新表管理")
//	    // 返回表配置对象
//	    return newTable
//	}
//
//	// 示例 4: 删除表
//	delete(tables.Generators, "old_table")
//
//	// 示例 5: 检查表是否存在
//	if _, ok := tables.Generators["posts"]; ok {
//	    fmt.Println("posts 表存在")
//	} else {
//	    fmt.Println("posts 表不存在")
//	}
//
// 注意事项：
//   - 表名必须与数据库表名一致
//   - 表生成函数必须返回正确的表配置对象
//   - 表生成器映射中的键值对必须一一对应
//   - 表生成函数必须接收上下文对象作为参数
//   - 表生成函数必须返回 table.Table 类型的对象
//   - 动态添加表时需要确保表生成函数正确实现
//   - 删除表时需要确保该表不再被使用
//   - 表生成函数中应该处理可能的错误
//   - 表生成函数应该遵循 GoAdmin 的表配置规范
//
// 错误处理：
//   - 如果表生成函数返回 nil，会在运行时返回错误
//   - 如果表名不存在，查找操作会返回 false
//   - 如果表生成函数执行失败，会在运行时返回错误
//   - 如果表配置对象创建失败，会在运行时返回错误
//
// 性能优化：
//   - 表生成器映射在包初始化时创建，避免重复创建
//   - 表生成函数可以延迟加载，提高启动速度
//   - 表配置对象可以缓存，避免重复创建
//   - 表生成函数应该尽量轻量，避免耗时操作
//
// 扩展性：
//   - 可以动态添加新的表生成函数
//   - 可以动态删除不需要的表生成函数
//   - 可以修改已存在的表生成函数
//   - 可以支持多数据源的表生成函数
//   - 可以支持自定义表生成逻辑
//
// 最佳实践：
//   - 在包初始化时注册所有表生成函数
//   - 使用常量定义表名，避免硬编码
//   - 表生成函数应该有清晰的命名规范
//   - 表生成函数应该有详细的注释说明
//   - 表生成函数应该处理可能的错误
//   - 表生成函数应该遵循单一职责原则
//   - 表生成函数应该易于测试和维护
//
// 常见问题：
//
//	Q: 如何添加新的表？
//	A: 定义表生成函数，然后在 Generators 映射中添加键值对
//
//	Q: 如何删除已有的表？
//	A: 使用 delete 函数从 Generators 映射中删除对应的键值对
//
//	Q: 如何查找表生成函数？
//	A: 使用表名作为键从 Generators 映射中查找对应的表生成函数
//
//	Q: 表生成函数的参数是什么？
//	A: 表生成函数接收一个上下文对象（*context.Context）作为参数
//
//	Q: 表生成函数的返回值是什么？
//	A: 表生成函数返回一个表配置对象（table.Table）
//
//	Q: 如何动态添加表？
//	A: 直接在 Generators 映射中添加新的键值对即可
//
//	Q: 如何删除表？
//	A: 使用 delete 函数从 Generators 映射中删除对应的键值对
//
//	Q: 如何检查表是否存在？
//	A: 使用表名作为键从 Generators 映射中查找，检查是否找到
//
// 相关类型：
//   - table.Generator: 表生成器类型定义
//   - table.Table: 表配置对象类型定义
//   - context.Context: 上下文对象类型定义
//   - table.Info: 表信息配置对象类型定义
//   - table.FormList: 表单配置对象类型定义
//
// 相关函数：
//   - GetPostsTable: 获取文章表模型配置
//   - GetAuthorsTable: 获取作者表模型配置
//   - GetExternalTable: 获取外部数据源表模型配置
//   - GetUserTable: 获取用户表模型配置
//
// 相关包：
//   - github.com/purpose168/GoAdmin/plugins/admin/modules/table: 表格模块包
//   - github.com/purpose168/GoAdmin/context: 上下文模块包
//   - github.com/purpose168/GoAdmin/modules/db: 数据库模块包
var Generators = map[string]table.Generator{
	"posts":    GetPostsTable,    // 文章表，使用 GetPostsTable 函数生成
	"authors":  GetAuthorsTable,  // 作者表，使用 GetAuthorsTable 函数生成
	"external": GetExternalTable, // 外部数据源表，使用 GetExternalTable 函数生成
}

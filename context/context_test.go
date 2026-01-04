// context包的测试文件
// 用于测试context包中的路由树、路径处理和处理器链功能
//
// 测试内容：
//   - slash函数：测试路径斜杠处理
//   - join函数：测试路径连接
//   - tree结构：测试路由树的路径添加和查找
//
// 测试目的：
//   - 验证路径处理的正确性
//   - 确保路由树能正确匹配路径
//   - 测试处理器链的执行顺序
//
// 运行测试：
//
//	go test ./context -v
//
// 作者: GoAdmin Team
// 创建日期: 2020-01-01
// 版本: 1.0.0
package context

import (
	"fmt"
	"testing"

	"github.com/magiconair/properties/assert"
)

// TestSlash测试slash函数的各种边界情况
//
// slash函数功能：确保路径以斜杠开头，并移除末尾的斜杠
//
// 测试用例：
//  1. 正常路径："/abc" → "/abc"（保持不变）
//  2. 空字符串："" → "/"（转换为根路径）
//  3. 末尾有斜杠："abc/" → "/abc"（添加开头斜杠，移除末尾斜杠）
//  4. 只有斜杠："/" → "/"（保持不变）
//  5. 两端都有斜杠："/abc/" → "/abc"（移除末尾斜杠）
//  6. 双斜杠："//" → "/"（简化为单个斜杠）
//
// 预期结果：所有测试用例都应该通过
func TestSlash(t *testing.T) {
	assert.Equal(t, "/abc", slash("/abc"))
	assert.Equal(t, "/", slash(""))
	assert.Equal(t, "/abc", slash("abc/"))
	assert.Equal(t, "/", slash("/"))
	assert.Equal(t, "/abc", slash("/abc/"))
	assert.Equal(t, "/", slash("//"))
}

// TestJoin测试join函数的路径连接功能
//
// join函数功能：将两个路径连接成一个完整的路径
//
// 测试用例：
//  1. 两个正常路径："/abc" + "/abc" → "/abc/abc"
//  2. 两个根路径："/" + "/" → "/"
//  3. 根路径+正常路径："/" + "/abc" → "/abc"
//  4. 正常路径+根路径："abc/" + "/" → "/abc"
//  5. 末尾有斜杠+根路径："/abc/" + "/" → "/abc"
//
// 预期结果：所有测试用例都应该通过
func TestJoin(t *testing.T) {
	assert.Equal(t, "/abc/abc", join(slash("/abc"), slash("/abc")))
	assert.Equal(t, "/", join(slash("/"), slash("/")))
	assert.Equal(t, "/abc", join(slash("/"), slash("/abc")))
	assert.Equal(t, "/abc", join(slash("abc/"), slash("/")))
	assert.Equal(t, "/abc", join(slash("/abc/"), slash("/")))
}

// TestTree测试路由树的路径添加和查找功能
//
// 路由树功能：
//   - addPath：向路由树添加路径和对应的处理器
//   - findPath：从路由树中查找路径对应的处理器
//
// 测试场景：
//  1. 添加相似路径，测试路由树的精确匹配
//  2. 添加带参数的路径，测试参数提取
//  3. 添加不同HTTP方法的路由，测试方法区分
//  4. 测试路径查找的正确性
//
// 路由结构：
//   - /adm：简单路径
//   - /admi：与/adm相似的路径
//   - /admin：管理后台根路径
//   - /admin/menu/new：菜单新建页面（GET和POST）
//   - /admin/info/:__prefix：带参数的信息页面
//   - /admin/info/:__prefix/detail：带参数的详情页面
//
// 测试步骤：
//  1. 创建路由树
//  2. 添加各种路径和处理器
//  3. 测试路径查找
//  4. 验证处理器链的执行
//  5. 打印路由树结构
func TestTree(t *testing.T) {
	tree := tree()

	// 添加简单路径 /adm，测试基础路由功能
	tree.addPath(stringToArr("/adm"), "GET", []Handler{func(ctx *Context) { fmt.Println(1) }})
	// 添加相似路径 /admi，测试路由树对相似路径的区分能力
	tree.addPath(stringToArr("/admi"), "GET", []Handler{func(ctx *Context) { fmt.Println(1) }})
	// 添加管理后台根路径 /admin
	tree.addPath(stringToArr("/admin"), "GET", []Handler{func(ctx *Context) { fmt.Println(1) }})
	// 添加菜单新建页面的POST方法路由
	tree.addPath(stringToArr("/admin/menu/new"), "POST", []Handler{func(ctx *Context) { fmt.Println(1) }})
	// 添加菜单新建页面的GET方法路由，测试同一路径不同HTTP方法的处理
	tree.addPath(stringToArr("/admin/menu/new"), "GET", []Handler{func(ctx *Context) { fmt.Println(1) }})
	// 添加带参数的信息页面路由，包含三个处理器：auth、init、info
	tree.addPath(stringToArr("/admin/info/:__prefix"), "GET", []Handler{
		func(ctx *Context) { fmt.Println("auth") },
		func(ctx *Context) { fmt.Println("init") },
		func(ctx *Context) { fmt.Println("info") },
	})
	// 添加带参数的详情页面路由，包含两个处理器：auth、detail
	tree.addPath(stringToArr("/admin/info/:__prefix/detail"), "GET", []Handler{
		func(ctx *Context) { fmt.Println("auth") },
		func(ctx *Context) { fmt.Println("detail") },
	})

	// 测试用例1：查找 /admin/menu/new 的GET方法路由
	fmt.Println("/admin/menu/new", "GET")
	h := tree.findPath(stringToArr("/admin/menu/new"), "GET")
	assert.Equal(t, h != nil, true)
	printHandler(h)
	// 测试用例2：查找 /admin/menu/new 的POST方法路由
	fmt.Println("/admin/menu/new", "POST")
	h = tree.findPath(stringToArr("/admin/menu/new"), "POST")
	assert.Equal(t, h != nil, true)
	printHandler(h)
	// 测试用例3：查找不存在的路径 /admin/me/new，预期返回nil
	fmt.Println("/admin/me/new", "POST")
	h = tree.findPath(stringToArr("/admin/me/new"), "POST")
	assert.Equal(t, h == nil, true)
	printHandler(h)
	// 测试用例4：查找带参数的路径 /admin/info/user，测试参数匹配功能
	fmt.Println("/admin/info/user", "GET")
	h = tree.findPath(stringToArr("/admin/info/user"), "GET")
	assert.Equal(t, h != nil, true)
	printHandler(h)
	// 测试用例5：查找带参数的详情路径 /admin/info/user/detail，测试嵌套参数匹配
	fmt.Println("/admin/info/user/detail", "GET")
	h = tree.findPath(stringToArr("/admin/info/user/detail"), "GET")
	assert.Equal(t, h != nil, true)
	printHandler(h)
	// 打印路由树结构，用于调试和验证路由树的构建是否正确
	fmt.Println("=========== printChildren ===========")
	tree.printChildren()
}

// printHandler执行处理器链中的所有处理器
//
// 参数说明：
//   - h: 处理器链，包含多个Handler函数
//
// 工作原理：
//   - 遍历处理器链中的每个处理器
//   - 创建一个空的Context对象
//   - 调用每个处理器并传入Context
//
// 使用场景：
//   - 测试处理器链的执行
//   - 调试处理器行为
//   - 验证处理器顺序
func printHandler(h []Handler) {
	for _, value := range h {
		value(&Context{})
	}
}

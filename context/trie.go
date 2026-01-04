// Copyright 2019 GoAdmin Core Team. All rights reserved.
// Use of this source code is governed by a Apache-2.0 style
// license that can be found in the LICENSE file.

package context

import "fmt"

// node是trie树的节点结构体
//
// 字段说明：
//   - children: 子节点列表
//   - value: 节点值（路径片段）
//   - method: HTTP方法列表（支持同一路径多个HTTP方法）
//   - handle: 处理器链列表（与method一一对应）
//
// 使用场景：
//   - 构建路由树
//   - 高效路径匹配
//   - 支持通配符路由
//
// 工作原理：
//   - 每个节点代表URL路径的一个片段
//   - 通过children指针构建树形结构
//   - 支持通配符（*）匹配任意路径片段
//   - 同一节点可以注册多个HTTP方法的处理器
type node struct {
	children []*node
	value    string
	method   []string
	handle   [][]Handler
}

// tree创建一个新的trie树根节点
//
// 返回值：
//   - *node: 新的根节点，value为"/"
//
// 工作原理：
//   - 创建一个node实例
//   - 初始化children为空切片
//   - 设置value为"/"（根路径）
//   - handle为nil（暂无处理器）
//
// 使用场景：
//   - 初始化路由树
//   - 创建新的路由组
func tree() *node {
	return &node{
		children: make([]*node, 0),
		value:    "/",
		handle:   nil,
	}
}

// hasMethod检查节点是否包含指定的HTTP方法
//
// 参数说明：
//   - method: HTTP方法（如"GET"、"POST"等）
//
// 返回值：
//   - int: 方法在method切片中的索引，如果不存在则返回-1
//
// 工作原理：
//   - 遍历method切片
//   - 查找与参数匹配的方法
//   - 返回匹配的索引位置
//
// 使用场景：
//   - 查找路由处理器
//   - 验证HTTP方法是否已注册
func (n *node) hasMethod(method string) int {
	for k, m := range n.method {
		if m == method {
			return k
		}
	}
	return -1
}

// addMethodAndHandler向节点添加HTTP方法和对应的处理器
//
// 参数说明：
//   - method: HTTP方法（如"GET"、"POST"等）
//   - handler: 处理器链
//
// 工作原理：
//   - 将方法添加到method切片
//   - 将处理器链添加到handle切片
//   - method和handle通过索引一一对应
//
// 使用场景：
//   - 注册路由处理器
//   - 支持同一路径多个HTTP方法
//
// 使用示例：
//
//	node.addMethodAndHandler("GET", handler1)
//	node.addMethodAndHandler("POST", handler2)
func (n *node) addMethodAndHandler(method string, handler []Handler) {
	n.method = append(n.method, method)
	n.handle = append(n.handle, handler)
}

// addChild向节点添加子节点
//
// 参数说明：
//   - child: 要添加的子节点
//
// 工作原理：
//   - 将子节点添加到children切片
//   - 保持子节点的顺序
//
// 使用场景：
//   - 构建路由树
//   - 添加新的路径分支
func (n *node) addChild(child *node) {
	n.children = append(n.children, child)
}

// addContent添加或查找子节点
//
// 参数说明：
//   - value: 子节点的值（路径片段）
//
// 返回值：
//   - *node: 子节点，如果不存在则创建新节点
//
// 工作原理：
//   - 首先搜索是否已存在相同value的子节点
//   - 如果存在则返回该子节点
//   - 如果不存在则创建新节点并添加到children
//   - 支持通配符（*）匹配
//
// 使用场景：
//   - 构建路由树
//   - 添加路径节点
//
// 使用示例：
//
//	child := node.addContent("users")
//	child.addMethodAndHandler("GET", handler)
func (n *node) addContent(value string) *node {
	var child = n.search(value)
	if child == nil {
		child = &node{
			children: make([]*node, 0),
			value:    value,
		}
		n.addChild(child)
	}
	return child
}

// search在子节点中查找指定值的节点
//
// 参数说明：
//   - value: 要查找的节点值（路径片段）
//
// 返回值：
//   - *node: 匹配的子节点，如果不存在则返回nil
//
// 工作原理：
//   - 遍历children切片
//   - 查找value匹配的子节点
//   - 支持通配符（*）匹配任意值
//
// 使用场景：
//   - 路径匹配
//   - 查找路由节点
func (n *node) search(value string) *node {
	for _, child := range n.children {
		if child.value == value || child.value == "*" {
			return child
		}
	}
	return nil
}

// addPath向路由树添加路径和对应的处理器
//
// 参数说明：
//   - paths: 路径片段数组（由stringToArr函数生成）
//   - method: HTTP方法（如"GET"、"POST"等）
//   - handler: 处理器链
//
// 工作原理：
//   - 从根节点开始遍历路径片段
//   - 为每个片段创建或查找对应的子节点
//   - 在最后一个节点上添加HTTP方法和处理器
//
// 使用场景：
//   - 注册路由
//   - 构建路由树
//
// 使用示例：
//
//	paths := stringToArr("/api/users/:id")
//	root.addPath(paths, "GET", handler)
func (n *node) addPath(paths []string, method string, handler []Handler) {
	child := n
	for i := 0; i < len(paths); i++ {
		child = child.addContent(paths[i])
	}
	child.addMethodAndHandler(method, handler)
}

// findPath在路由树中查找路径和HTTP方法对应的处理器
//
// 参数说明：
//   - paths: 路径片段数组（由stringToArr函数生成）
//   - method: HTTP方法（如"GET"、"POST"等）
//
// 返回值：
//   - []Handler: 处理器链，如果未找到则返回nil
//
// 工作原理：
//   - 从根节点开始遍历路径片段
//   - 逐级查找对应的子节点
//   - 如果任一节点不存在则返回nil
//   - 在最后一个节点上查找指定的HTTP方法
//   - 返回对应的处理器链
//
// 使用场景：
//   - 路由匹配
//   - 查找请求处理器
//
// 使用示例：
//
//	paths := stringToArr("/api/users/123")
//	handler := root.findPath(paths, "GET")
//	if handler != nil {
//	    handler(ctx)
//	}
func (n *node) findPath(paths []string, method string) []Handler {
	child := n
	for i := 0; i < len(paths); i++ {
		child = child.search(paths[i])
		if child == nil {
			return nil
		}
	}

	methodIndex := child.hasMethod(method)
	if methodIndex == -1 {
		return nil
	}

	return child.handle[methodIndex]
}

// print打印当前节点的信息
//
// 工作原理：
//   - 使用fmt.Println打印节点
//   - 主要用于调试和测试
func (n *node) print() {
	fmt.Println(n)
}

// printChildren递归打印当前节点及其所有子节点
//
// 工作原理：
//   - 先打印当前节点
//   - 递归打印每个子节点
//   - 采用深度优先遍历方式
//
// 使用场景：
//   - 调试路由树结构
//   - 查看路由注册情况
func (n *node) printChildren() {
	n.print()
	for _, child := range n.children {
		child.printChildren()
	}
}

// stringToArr将URL路径字符串转换为路径片段数组
//
// 参数说明：
//   - path: URL路径字符串（如"/api/users/:id/info"）
//
// 返回值：
//   - []string: 路径片段数组，通配符参数被替换为"*"
//
// 转换规则：
//   - 忽略开头的"/"
//   - 以"/"作为分隔符
//   - 包含":"的路径片段（如":id"）被识别为通配符，替换为"*"
//
// 转换示例：
//   - "/api/users/:id"       => ["api", "users", "*"]
//   - "/api/users/:id/info"  => ["api", "users", "*", "info"]
//   - "/api/users"           => ["api", "users"]
//   - "/"                    => []
//
// 工作原理：
//   - 遍历路径字符串的每个字符
//   - 跳过开头的"/"
//   - 检测":"标记通配符参数
//   - 以"/"作为分隔符提取路径片段
//   - 通配符片段替换为"*"
//
// 使用场景：
//   - 路径解析
//   - 构建路由树
//   - 支持动态路由参数
//
// 使用示例：
//
//	paths := stringToArr("/api/users/:id/posts/:postId")
//	// 返回 ["api", "users", "*", "posts", "*"]
//	root.addPath(paths, "GET", handler)
func stringToArr(path string) []string {
	var (
		paths      = make([]string, 0)
		start      = 0
		end        int
		isWildcard = false
	)
	for i := 0; i < len(path); i++ {
		if i == 0 && path[0] == '/' {
			start = 1
			continue
		}
		if path[i] == ':' {
			isWildcard = true
		}
		if i == len(path)-1 {
			end = i + 1
			if isWildcard {
				paths = append(paths, "*")
			} else {
				paths = append(paths, path[start:end])
			}
		}
		if path[i] == '/' {
			end = i
			if isWildcard {
				paths = append(paths, "*")
			} else {
				paths = append(paths, path[start:end])
			}
			start = i + 1
			isWildcard = false
		}
	}
	return paths
}

package gee

import (
	"net/http"
	"strings"
)

// 路由表，由HTTP节点前缀树表(GET/POST等键)和处理函数表组成
type router struct {
	roots map[string]*node
	// 不同的路由名称对应不同的路由处理函数
	handlers map[string]HandlerFunc
}

// 初始化路由表
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 辅助函数：解析路由模式字符串/p/:lang/doc为["p", ":lang", "doc"]
func parsePattern(pattern string) []string {
	// pattern := "/p/:lang/doc"
	// vs 的值为: ["", "p", ":lang", "doc"]
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		// 如果 item 不是空字符串，则将其添加到 parts 切片中。
		if item != "" {
			parts = append(parts, item)
			// 遇到 item 的第一个字符是 * 时，表示这是一个通配符，停止添加后续部分。
			// 输入"/static/*filepath/more"
			// vs["", "static", "*filepath", "more"]
			// parts(break)["static", "*filepath"]
			// parts(无break)["static", "*filepath", "more"]
			if item[0] == '*' {
				// break跳出for循环
				// 程序执行for循环之后的代码
				break
			}
		}
	}
	return parts
}

// 添加路由时需要路由名字和处理函数
// 输入GET /p/go/doc 处理函数
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	// 输入/p/go/doc 返回["p", "go", "doc"]
	// 输入/static/*filepath/more 返回["static", "*filepath"]
	// 返回字符串切片
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	// 初始化路由树
	_, ok := r.roots[method]
	// 检查当前 HTTP 方法（如GET）是否已经初始化了对应的路由树。
	// 如果没有，则创建一个新的根节点（&node{}），并将其存储到 r.roots[method] 中。
	if !ok {
		r.roots[method] = &node{}
	}
	// 将路由模式插入到前缀树中
	r.roots[method].insert(pattern, parts, 0)

	// 把路由名字和路由处理函数放入路由表
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	// 用于存储提取的动态参数
	params := make(map[string]string)
	// 查找PUT/POST的根节点
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}
	// 在前缀树中查找与 searchParts 匹配的路由节点
	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		// 若path = "/p/go/doc" 则不执行以下代码
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		// 若path = "/p/go/doc" 返回相应的节点和空映射
		return n, params
	}

	return nil, nil
}

// 根据requset的method和path查找对应的处理函数，未找到返回404
func (r *router) handle(c *Context) {
	// n返回叶子节点，params返回对应的路径map
	// n:{pattern:"/static/*filepath",part:"*filepath",}
	// params:{"filepath": "css/style.css"},
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}

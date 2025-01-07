package gee

import (
	"log"
	"net/http"
)

// 路由表
type router struct {
	// 不同的路由名称对应不同的路由处理函数
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	// 初始化一个空的路由表
	return &router{handlers: make(map[string]HandlerFunc)}
}

// 添加路由时需要路由名字和处理函数
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	// 记录路由信息
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	// 把路由名字和路由处理函数放入路由表
	r.handlers[key] = handler
}

// 根据requset的method和path查找对应的处理函数，未找到返回404
func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	// 查路由表，并执行处理函数
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}

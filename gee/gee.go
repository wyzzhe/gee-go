package gee

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Engine struct {
	router map[string]HandlerFunc
}

func New() *Engine {
	// 路由映射表router，[路由名称]处理函数
	// make 初始化映射的哈希表元数据，为哈希表分配内存
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	// 请求方法和静态路由地址 GET-/
	key := method + "-" + pattern
	// 把方法注册至路由名称
	engine.router[key] = handler
}

// 注册GET下的路由的具体方法
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// 将服务注册至具体的端口号
func (engine *Engine) Run(addr string) (err error) {
	// 只有实现了Handler的struct才能传入ListenAndServe
	// 启动一个 HTTP 服务器并监听指定的地址和端口，engine是具体的服务处理器
	return http.ListenAndServe(addr, engine)
}

// 为Engine实现Handler，客户端传入Request
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 解析客户端Request
	key := req.Method + "-" + req.URL.Path
	// 如果客户端Request在路由哈希表中
	if handler, ok := engine.router[key]; ok {
		// type HandlerFunc func(http.ResponseWriter, *http.Request)
		// 执行客户端Request的具体响应方法
		handler(w, req)
	} else {
		// 服务器无法找到客户端请求的资源(页面)
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

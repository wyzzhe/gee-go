package gee

import (
	"net/http"
)

// 函数只接收Context型指针且无返回值
type HandlerFunc func(*Context)

type Engine struct {
	router *router
}

func New() *Engine {
	// 路由映射表router，[路由名称]处理函数
	// make 初始化映射的哈希表元数据，为哈希表分配内存
	return &Engine{router: newRouter()}
}

// 添加路由时需要路由名字和处理函数
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

// 注册GET下的路由的处理函数，此处理函数需要传入*Context
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
	// 当 HTTP 服务器接收到请求时，Go 的 net/http 包会调用 engine 的 ServeHTTP 方法。
	return http.ListenAndServe(addr, engine)
}

// 为Engine实现Handler，客户端传入Request
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 用http的w和req构造上下文结构体并返回指针
	c := newContext(w, req)
	// 路由器处理的实现细节在handle方法中
	// 此处传入handler所需的*Context
	engine.router.handle(c)
}

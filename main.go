package main

import (
	"fmt"
	"net/http"

	"gee"
)

func main() {
	r := gee.New()
	// 注册路由名称和具体服务端响应
	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	})

	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			// 输出格式：Header[k]v
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})

	r.Run(":9999")
}
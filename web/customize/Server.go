// @Author: zqy
// @File: Server.go
// @Date: 2023/5/4 12:06

package customize

import "net/http"

type HandleFunc func(ctx *Context)

type Server interface {
	// Handler  用于处理逻辑handler
	http.Handler
	// Start 用户服务器的启动，方便控制生命周期
	Start(address string) error

	AddRoute(method, path string, handleFunc HandleFunc)
}

var _ Server = &HttpServer{}

type HttpServer struct {
}

// ServeHTTP 核心入口
// Context的构建
// 路由匹配
// 执行业务逻辑
func (h *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := &Context{
		req:  r,
		resp: w,
	}

	h.serve(c)
}
func (h *HttpServer) serve(c *Context) {
	// 查找路由，执行操作
}
func (h *HttpServer) Start(address string) error {

	// 启动前进行的一些操作

	return http.ListenAndServe(address, h)
}

// AddRoute 路由注册
// 这里handleFunc只传入一个，方便进行处理
// 使用可选参数的话 会考虑到很多问题，并且用户还可能不传参数，编译不会检查到，导致程序出错
func (h *HttpServer) AddRoute(method, path string, handleFunc HandleFunc) {

}

// Get 衍生API
func (h *HttpServer) Get(path string, handleFunc HandleFunc) {
	h.AddRoute(http.MethodGet, path, handleFunc)
}

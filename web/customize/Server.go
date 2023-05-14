// @Author: zqy
// @File: Server.go
// @Date: 2023/5/4 12:06

package customize

import (
	"fmt"

	"net/http"
)

type HandleFunc func(ctx *Context)

type Server interface {
	// Handler  用于处理逻辑handler
	http.Handler
	// Start 用户服务器的启动，方便控制生命周期
	Start(address string) error
	//router
	//
	addRoute(method, path string, handleFunc HandleFunc)

	findRoute(method, path string) (*matchInfo, bool)
}

var _ Server = &HttpServer{}

type HttpServer struct {
	//route.
	*router
	log        func(msg string, args ...any)
	middleWare []MiddleWare
}
type HTTPServerOption func(server *HttpServer)

func NewHttpServer(opts ...HTTPServerOption) *HttpServer {
	res := &HttpServer{
		router: newRouter(),
		log: func(msg string, args ...any) {
			fmt.Printf(msg, args...)
		},
	}

	for _, opt := range opts {
		opt(res)
	}
	return res
}

func ServerWithMiddleWare(mdls ...MiddleWare) HTTPServerOption {
	return func(server *HttpServer) {
		server.middleWare = mdls
	}
}

// ServeHTTP 核心入口
// Context的构建
// 路由匹配
// 执行业务逻辑
func (h *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := &Context{
		Req:  r,
		Resp: w,
	}

	root := h.serve
	// 从后向前进行组装链条
	//root = func(func(h.serve(c *Context)))
	for i := len(h.middleWare) - 1; i >= 0; i-- {
		root = h.middleWare[i](root)
	}
	//进行状态数据的刷新

	var m MiddleWare = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			next(ctx)
			h.flashResp(ctx)
		}
	}

	// 最先执行，然后不断执行middleware，最后返回的时候刷新状态
	root = m(root)
	// 所以最后执行root 就是从前往后了
	root(c)
}

func (h *HttpServer) flashResp(ctx *Context) {
	if ctx.RespCode != 0 {
		ctx.Resp.WriteHeader(ctx.RespCode)
	}
	ctx.Resp.Header().Set("Content-Type", "application/json")

	write, err := ctx.Resp.Write(ctx.RespData)
	if err != nil || write != len(ctx.RespData) {
		//return
		//log.Fatal(err)
		h.log("错误信息%v", err)
	}
}
func (h *HttpServer) serve(c *Context) {
	// 查找路由，执行操作
	info, ok := h.findRoute(c.Req.Method, c.Req.URL.Path)
	if !ok || info.n.handler == nil {
		//c.Resp.WriteHeader(http.StatusNotFound)
		c.RespCode = http.StatusNotFound
		c.RespData = []byte("NOT FOUND")
		//todo 404 重定向到某个界面
		return
	}
	c.pathParams = info.pathParam
	c.MatchedRoute = info.n.route
	info.n.handler(c)
}
func (h *HttpServer) Start(address string) error {

	// 启动前进行的一些操作

	return http.ListenAndServe(address, h)
}

// AddRoute 路由注册
// 这里handleFunc只传入一个，方便进行处理
// 使用可选参数的话 会考虑到很多问题，并且用户还可能不传参数，编译不会检查到，导致程序出错
func (h *HttpServer) addRoute(method, path string, handleFunc HandleFunc) {
	h.router.AddRoute(method, path, handleFunc)
}
func (h *HttpServer) findRoute(method string, path string) (*matchInfo, bool) {
	return h.router.findRoute(method, path)
}

// Get 衍生API
func (h *HttpServer) Get(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodGet, path, handleFunc)
}

// Post  衍生API
func (h *HttpServer) Post(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodPost, path, handleFunc)
}

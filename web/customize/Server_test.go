// @Author: zqy
// @File: Server_test.go
// @Date: 2023/5/4 15:00

package customize

import (
	"fmt"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	s := NewHttpServer()
	//
	//s.Get("/", func(ctx *Context) {
	//	fmt.Println("ok")
	//	fmt.Println("ok1")
	//})
	s.Post("/Form", func(ctx *Context) {
		// PostForm
		//ctx.req.ParseForm()
		//asInt64, err := ctx.PathValueAsString("id").AsInt64()
		//if err != nil {
		//	ctx.resp.Write([]byte("id非法值"))
		//} else {
		//	ctx.resp.Write([]byte(fmt.Sprintf("hello,value:%d", asInt64)))
		//}
		//ctx.resp.Write([]byte(ctx.req.URL.Path))
		value, err := ctx.QueryValue("id")
		if err != nil {
			ctx.resp.Write([]byte("id非法值"))
		} else {
			ctx.resp.Write([]byte(fmt.Sprintf("id正常,value:%s", value)))
		}
	})

	s.Start(":9090")
}

func TestHttpServer_ServeHTTP(t *testing.T) {
	server := NewHttpServer()
	server.middleWare = []MiddleWare{
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第一个before")
				next(ctx)
				fmt.Println("第一个after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第二个before")
				next(ctx)
				fmt.Println("第二个after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第三个中断")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第四个看不到")
			}
		},
	}
	server.ServeHTTP(nil, &http.Request{})
}

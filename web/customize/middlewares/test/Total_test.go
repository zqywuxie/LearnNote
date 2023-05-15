// @Author: zqy
// @File: Total_test.go
// @Date: 2023/5/15 14:44
// @Description todo

package test

import (
	"GoCode/web/customize"
	"GoCode/web/customize/middlewares/accesslog"
	recover2 "GoCode/web/customize/middlewares/recover"
	"log"
	"net/http"
	"testing"
)

func TestHttpServer_ServeHTTP1(t *testing.T) {
	build := accesslog.NewMiddlewareBuilder().Build()
	//errhdl := errhdl.NewMiddlewareBuilder().Build()
	recoverMsg := &recover2.MiddlewareBuilder{
		StatusCode: 500,
		ErrMsg:     "出错了",
		LogFunc: func(ctx *customize.Context) {
			log.Println(ctx.Req.URL.Path)
		},
	}
	server := customize.NewHttpServer()
	server.Use(http.MethodGet, "/user/*", recoverMsg.Build(), build)

	server.Get("/user/add", func(ctx *customize.Context) {
		//fmt.Println("hello")
		panic("测试")
	})
	server.Start(":9090")
}

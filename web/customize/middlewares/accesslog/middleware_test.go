// @Author: zqy
// @File: middleware_test.go.go
// @Date: 2023/5/12 14:55
// @Description todo

package accesslog

import (
	"GoCode/web/customize"
	"fmt"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	//build := accessLog{}.Build()
	//accessLo
	//NewMiddlewareBuilder
	build := NewMiddlewareBuilder().Build()
	//build := customize..Build()
	server := customize.NewHttpServer(customize.ServerWithMiddleWare(build))
	//server
	server.Get("/a/b/*", func(ctx *customize.Context) {
		fmt.Println("hello")
	})
	//req, err := http.NewRequest(http.MethodGet, "/a/b/c", nil)
	//req.URL.Host = "localhost"
	//if err != nil {
	//	t.Fatal(err)
	//}
	//server.ServeHTTP(nil, req)
	server.Start(":9090")

}

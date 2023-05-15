// @Author: zqy
// @File: middleware_test.go.go
// @Date: 2023/5/14 16:35
// @Description todo

package recover

import (
	"GoCode/web/customize"
	"log"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := &MiddlewareBuilder{
		StatusCode: 500,
		ErrMsg:     "出错了",
		LogFunc: func(ctx *customize.Context) {
			log.Println(ctx.Req.URL.Path)
		},
	}
	s := customize.NewHttpServer(customize.ServerWithMiddleWare(builder.Build()))
	s.Get("/panic", func(ctx *customize.Context) {
		panic("出错了哈哈")
	})
	s.Start(":9090")
}

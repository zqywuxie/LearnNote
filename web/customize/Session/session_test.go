// @Author: zqy
// @File: session_test.go
// @Date: 2023/5/18 15:03
// @Description todo

package Session

import (
	"GoCode/web/customize"
	"net/http"
	"testing"
)

func TestSession(t *testing.T) {
	var m Manager
	server := customize.NewHttpServer(customize.ServerWithMiddleWare(
		func(next customize.HandleFunc) customize.HandleFunc {
			return func(ctx *customize.Context) {
				if ctx.Req.URL.Path == "/login" {
					next(ctx)
					return
				}
				_, err := m.GetSession(ctx)
				if err != nil {
					ctx.RespCode = http.StatusUnauthorized
					ctx.RespData = []byte("用户未登录")
				}
				next(ctx)
			}
		}))

	server.Get("/user", func(ctx *customize.Context) {
		session, _ := m.GetSession(ctx)
		value, _ := session.Get(ctx.Req.Context(), "zqy")
	})

	server.Start(":9090")
}

// @Author: zqy
// @File: session_test.go
// @Date: 2023/5/18 15:03
// @Description todo

package test

import (
	"GoCode/web/customize"
	"GoCode/web/customize/session"
	"GoCode/web/customize/session/cookie"
	"GoCode/web/customize/session/memory"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestSession(t *testing.T) {
	m := &session.Manager{
		Store:         memory.NewStore(time.Minute * 12),
		Propagator:    cookie.NewPropagator(),
		CtxSessionKey: "sessKey",
	}
	server := customize.NewHttpServer(customize.ServerWithMiddleWare(
		func(next customize.HandleFunc) customize.HandleFunc {
			return func(ctx *customize.Context) {
				if ctx.Req.URL.Path == "/login" {
					fmt.Println("登陆吧")
					next(ctx)
					return
				}
				fmt.Println("ok")
				_, err := m.GetSession(ctx)
				if err != nil {
					ctx.RespCode = http.StatusUnauthorized
					ctx.RespData = []byte("用户未登录")
					return
				}

				_ = m.RefreshSession(ctx)
				next(ctx)
			}
		}))

	server.Get("/user", func(ctx *customize.Context) {
		session, _ := m.GetSession(ctx)
		value, _ := session.Get(ctx.Req.Context(), "username")
		ctx.RespData = []byte(value.(string))

	})
	server.Post("/logout", func(ctx *customize.Context) {
		err := m.RemoveSession(ctx)
		if err != nil {
			ctx.RespCode = http.StatusInternalServerError
			ctx.RespData = []byte("退出失败")
			return
		}
		ctx.RespCode = http.StatusOK
		ctx.RespData = []byte("退出成功")
	})
	server.Post("/login", func(ctx *customize.Context) {
		session, err := m.InitSession(ctx)
		if err != nil {
			ctx.RespCode = http.StatusInternalServerError
			ctx.RespData = []byte("登录失败")
			return
		}
		err = session.Set(ctx.Req.Context(), "username", "zqy")
		if err != nil {
			ctx.RespCode = http.StatusInternalServerError
			ctx.RespData = []byte("session：设置值失败")
			return
		}
		ctx.RespCode = http.StatusOK
		ctx.RespData = []byte("登录成功")

	})

	server.Start(":9090")
}

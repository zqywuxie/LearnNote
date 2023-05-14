// @Author: zqy
// @File: middleware_test.go.go
// @Date: 2023/5/14 15:35
// @Description todo

package Prometheus

import (
	"GoCode/web/customize"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := MiddlewareBuilder{
		Subsystem: "web",
		Name:      "http_response",
		Help:      "my_help",
	}
	server := customize.NewHttpServer(customize.ServerWithMiddleWare(builder.Build()))
	server.Get("/user,", func(ctx *customize.Context) {
		val := rand.Intn(1000) + 1
		time.Sleep(time.Duration(val) * time.Millisecond)
		ctx.RespJson(http.StatusOK, User{Name: "ok"})
	})

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":9008", nil)
	}()
	server.Start(":9090")
}

type User struct {
	Name string
}

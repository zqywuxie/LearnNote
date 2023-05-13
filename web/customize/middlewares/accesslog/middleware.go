// @Author: zqy
// @File: middleware.go
// @Date: 2023/5/12 14:27
// @Description todo

package accesslog

import (
	"GoCode/web/customize"
	"encoding/json"
	"log"
)

type MiddlewareBuilder struct {
	logFunc func(accessLog string)
}

// accessLog 访问日志，记录请求
type accessLog struct {
	Host string `json:"host,omitempty"`
	// 匹配到的路由
	Route      string `json:"route,omitempty"`
	HTTPMethod string `json:"http_method"`
	Path       string `json:"path,omitempty"`
}

func (m MiddlewareBuilder) Build() customize.MiddleWare {
	return func(next customize.HandleFunc) customize.HandleFunc {
		return func(ctx *customize.Context) {
			// 使用defer 确保next里面发生panic，也能将请求记录下来
			// Path：指路由匹配pattern
			// Route：MatchedRoute指next调用完后的最终路由树
			defer func() {
				l := accessLog{
					Host:       ctx.Req.Host,
					Route:      ctx.MatchedRoute,
					HTTPMethod: ctx.Req.Method,
					Path:       ctx.Req.URL.Path,
				}
				data, _ := json.Marshal(l)
				m.logFunc(string(data))
			}()
			next(ctx)
		}
	}
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{

		logFunc: func(accessLog string) {
			log.Println(accessLog)
		},
	}
}

// LogFunc 定义方法打印日志
func (m *MiddlewareBuilder) LogFunc(logFunc func(accessLog string)) *MiddlewareBuilder {
	m.logFunc = logFunc
	return m
}

// @Author: zqy
// @File: middleware.go
// @Date: 2023/5/14 16:02
// @Description todo

package errhdl

import (
	"GoCode/web/customize"
)

type MiddlewareBuilder struct {
	// resp 根据返回码 来返回值
	// 缺陷只能返回固定的值,而不是动态
	// todo 根据路由跳转
	resp map[int][]byte
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{resp: make(map[int][]byte, 64)}
}

// AddData 可以使用链式调用
func (m *MiddlewareBuilder) AddData(code int, data []byte) *MiddlewareBuilder {
	m.resp[code] = data
	return m
}
func (m *MiddlewareBuilder) Build() customize.MiddleWare {
	return func(next customize.HandleFunc) customize.HandleFunc {
		return func(ctx *customize.Context) {
			next(ctx)
			resp, ok := m.resp[ctx.RespCode]
			if ok {
				ctx.RespData = resp
			}

		}
	}
}

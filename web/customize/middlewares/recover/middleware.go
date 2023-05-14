// @Author: zqy
// @File: middleware.go
// @Date: 2023/5/14 16:29
// @Description panic恢复

package recover

import "GoCode/web/customize"

type MiddlewareBuilder struct {
	statusCode int
	ErrMsg     string
	//LogFunc func(err any)
	LogFunc func(ctx *customize.Context)
}

func (m *MiddlewareBuilder) Build() customize.MiddleWare {
	return func(next customize.HandleFunc) customize.HandleFunc {
		return func(ctx *customize.Context) {
			defer func() {
				// 判断是否有panic,有就进行篡改响应数据
				if err := recover(); err != nil {
					ctx.RespCode = m.statusCode
					ctx.RespData = []byte(m.ErrMsg)
					m.LogFunc(ctx)
				}
			}()
			next(ctx)
		}
	}
}

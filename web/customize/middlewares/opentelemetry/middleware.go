// @Author: zqy
// @File: middleware.go
// @Date: 2023/5/12 16:05
// @Description tracing,数据链路跟踪

package opentelemetry

import (
	"GoCode/web/customize"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type MiddlewareBuilder struct {
	Tracer trace.Tracer
}

// //  NewMiddlewareBuilder 如果属性设置为私有，要求用户一定传入
//
//	func NewMiddlewareBuilder(tracer trace.Tracer) *MiddlewareBuilder {
//		return &MiddlewareBuilder{Tracer: tracer}
//	}
const instrumentationName = "GoCode/web/customize/opentelemetry/middleware"

func (m *MiddlewareBuilder) Build() customize.MiddleWare {
	if m.Tracer == nil {
		m.Tracer = otel.GetTracerProvider().Tracer(instrumentationName)
	}
	return func(next customize.HandleFunc) customize.HandleFunc {
		return func(ctx *customize.Context) {

			// 不同进程之间的span如何结合。尝试与客户端的tracer结合在一起 todo
			reqctx := ctx.Req.Context()
			//otel.GetTextMapPropagator().Extract(reqctx, propagation.HeaderCarrier{})
			carrier := &propagation.HeaderCarrier{}
			otel.GetTextMapPropagator().Inject(reqctx, carrier)

			reqctx, span := m.Tracer.Start(reqctx, "unknown")
			defer span.End()

			defer func() {
				span.SetAttributes(attribute.String("http.host", ctx.Req.Host))
				span.SetAttributes(attribute.String("http.method", ctx.Req.Method))
				span.SetAttributes(attribute.String("http.url", ctx.Req.URL.String()))
				span.SetAttributes(attribute.Int("http.status", ctx.RespCode))
				// 只有next执行完才可以得到的匹配路径
				span.SetName(ctx.MatchedRoute)
			}()

			// 想要将后面的span挂载到根路径下
			ctx.Req = ctx.Req.WithContext(reqctx)

			//还可以继续记录一些东西
			next(ctx)

			// 响应数据/码呢？ 拿ctx.Resp强转到某个类型然后拿到响应数据吗？
			// 不推荐，应该实现类有很多，不知道是具体的哪个
		}
	}
}

// @Author: zqy
// @File: middleware_test.go.go
// @Date: 2023/5/13 14:25
// @Description todo

package opentelemetry

import (
	"GoCode/web/customize"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	tracer := otel.GetTracerProvider().Tracer(instrumentationName)
	builder := MiddlewareBuilder{Tracer: tracer}
	server := customize.NewHttpServer(customize.ServerWithMiddleWare(builder.Build()))

	server.Get("/user", func(ctx *customize.Context) {
		c, span := tracer.Start(ctx.Req.Context(), "first_layer")
		defer span.End()
		c, second := tracer.Start(c, "second_layer")
		time.Sleep(time.Second)
		_, third1 := tracer.Start(c, "third1_layer")
		time.Sleep(100 * time.Millisecond)
		third1.End()
		_, third2 := tracer.Start(c, "third2_layer")
		time.Sleep(100 * time.Millisecond)
		third2.End()
		second.End()

		ctx.RespJson(http.StatusOK, User{
			Name: "zqy",
		})
	})
	server.Get("/", func(ctx *customize.Context) {

		ctx.RespJson(http.StatusOK, map[string]string{"zqy": "123"})
	})
	initZipkin(t)
	server.Start(":9091")
}

type User struct {

	//todo 关于json.marshal结构体打印得到空的问题，因为结构体的首字母为小写私有，改成大写即可
	Name string `json:"name,omitempty"`
}

func initZipkin(t *testing.T) {
	exporter, err := zipkin.New(
		"http://localhost:9411/api/v2/spans",
		zipkin.WithLogger(log.New(os.Stderr, "zipkin", log.Llongfile)),
	)
	if err != nil {
		t.Fatal(err)
	}
	batcher := sdktrace.NewBatchSpanProcessor(exporter)
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("zipkin-test"),
		)))
	otel.SetTracerProvider(provider)
}

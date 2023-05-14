// @Author: zqy
// @File: middleware.go
// @Date: 2023/5/14 14:56
// @Description Prometheus 数据统计

package Prometheus

import (
	"GoCode/web/customize"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type MiddlewareBuilder struct {
	Subsystem   string
	Name        string
	Help        string
	ConstLabels map[string]string
	//
}

func (m *MiddlewareBuilder) Build() customize.MiddleWare {
	summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Subsystem:   m.Subsystem,
		Name:        m.Name,
		Help:        m.Help,
		ConstLabels: m.ConstLabels,
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.90:  0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	}, []string{"pattern", "method", "status"})

	// 记得注册一下
	prometheus.MustRegister(summaryVec)
	return func(next customize.HandleFunc) customize.HandleFunc {
		return func(ctx *customize.Context) {
			startTime := time.Now()
			defer func() {
				duration := time.Now().Sub(startTime).Milliseconds()
				pattern := ctx.MatchedRoute
				if pattern == "" {
					pattern = "unknown"
				}
				method := ctx.Req.Method
				code := ctx.RespCode
				summaryVec.WithLabelValues(pattern, method, strconv.Itoa(code)).Observe(float64(duration))
			}()
			next(ctx)
		}

	}
}

func Counter() {
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "my_namespace",
		Subsystem: "my_subsystem",
		Name:      "my_name",
	})

	prometheus.MustRegister(counter)
	counter.Inc()
}

//func Vector() *prometheus.Summary {
//	summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
//		Subsystem: "http_request",
//		Name:      "geekbang",
//		Help:      "The statics info for http request",
//		ConstLabels: map[string]string{
//			"server":  "localhost:9091",
//			"env":     "test",
//			"appname": "test_app",
//		},
//	}, []string{"pattern", "method", "status"})
//
//	// 方便直接建造summary
//	//Observe(128) 响应实践128ms
//	//summaryVec.WithLabelValues("/user/:id", "POST", "200").Observe(128)
//}

// @Author: zqy
// @File: middleware.go
// @Date: 2023/5/8 17:06
// @Description 默认中间件

package customize

// MiddleWare 函数式的责任链模式(双向的)，洋葱模式(无侵入式增强核心功能,在已有的功能上进行封装，不改变原来的功能，比侵入式的性能较低 )
type MiddleWare func(next HandleFunc) HandleFunc

//type MiddleWareV1 interface {
//	Invoke(next HandleFunc) HandleFunc
//}
//
//// MiddleWareV2 拦截器模式，针对不同时间段
//type MiddleWareV2 interface {
//	Before(cxt *Context)
//	After(cxt *Context)
//	Surround(cxt *Context)
//}
//
//// Chain 使用切片，然后使用next进行传递
//type Chain []HandleFuncV1
//
//type HandleFuncV1 func(ctx *Context) (next bool)
//
//type ChainV1 struct {
//	handlers []HandleFuncV1
//	count    int
//}
//
//func (c *ChainV1) Run(ctx *Context) {
//	for _, h := range c.handlers {
//		next := h(ctx)
//		if !next {
//			return
//		}
//		c.count++
//	}
//}

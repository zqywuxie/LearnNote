// @Author: zqy
// @File: template.go
// @Date: 2023/5/15 21:59
// @Description todo

package Render

import (
	"bytes"
	"context"
	"html/template"
)

type TemplateEngine interface {

	// Render 渲染模板
	// tplName 模板引擎
	// data 渲染数据
	Render(ctx context.Context, tplName string, data any) ([]byte, error)

	// 直接将渲染数据输出到writer当中，缺点：其他中间件无法进行修改，比如渲染404错误页面
	//Render(ctx context.Context, tplName string, data any,writer io.Writer)  error

	//AddTemplate等方法，让具体实现管自己的模板
}

type GoTemplateEngine struct {
	T *template.Template
}

func (g *GoTemplateEngine) Render(ctx context.Context, tplName string, data any) ([]byte, error) {
	b := &bytes.Buffer{}
	err := g.T.ExecuteTemplate(b, tplName, data)
	return b.Bytes(), err
}

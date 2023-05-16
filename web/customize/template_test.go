//go:build e2e

// @Author: zqy
// @File: template_test.go
// @Date: 2023/5/15 22:25
// @Description todo

package customize

import (
	"GoCode/web/customize/Render"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"html/template"
	"log"
	"mime/multipart"
	"path/filepath"
	"testing"
)

type User struct {
	Name string
}

func TestHelloWorld(t *testing.T) {
	tpl := template.New("hello_world")
	// 解析模板，.代表的当前作用域的当前对象
	// 切片 index 按照下标索引
	tpl, err := tpl.Parse(`hello,{{index . 1}}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, []string{"Tom", "zqy"})
	require.NoError(t, err)
	assert.Equal(t, "hello,Tom", buffer.String())
}

func TestBasic(t *testing.T) {
	tpl := template.New("hello_world")
	// 解析模板，.代表的当前作用域的当前对象
	// 切片 index 按照下标索引
	tpl, err := tpl.Parse(`hello,{{.}}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, "Tom")
	require.NoError(t, err)
	assert.Equal(t, "hello,Tom", buffer.String())
}

func TestVar(t *testing.T) {
	const serviceTpl = `
	{{-  $service := .GenName -}}
	type {{ $service}} struct {
	Name string
}
`

	tpl := template.New("hello_world")
	// 解析模板，.代表的当前作用域的当前对象
	// 切片 index 按照下标索引
	tpl, err := tpl.Parse(`hello,{{.}}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, "Tom")
	require.NoError(t, err)
	assert.Equal(t, "hello,Tom", buffer.String())
}

type FuncCall struct {
	Slice []string
}

func (f FuncCall) Hello(firstName string, lastName string) string {
	return fmt.Sprintf("Hello,%s:%s", firstName, lastName)
}

// 模板使用内置方法
// .func args1 args2
func TestFuncCall(t *testing.T) {
	tpl := template.New("hello_world")
	parse, err := tpl.Parse(`
切片长度: {{len .Slice}}
say hello: {{.Hello "Tom" "Jerry"}}
打印数字: {{printf "%.2f" 1.234}}
`)
	assert.Nil(t, err)
	bs := &bytes.Buffer{}
	err = parse.Execute(bs,
		&FuncCall{Slice: []string{"Tom", "Jerry"}})

	assert.Nil(t, err)
	assert.Equal(t, `
切片长度: 2
say hello: Hello,Tom:Jerry
打印数字: 1.23
`, bs.String())
}
func TestLoop(t *testing.T) {
	tpl := template.New("loop")
	tpl, err := tpl.Parse(`
{{- range $idx,$ele := .}}
{{- $idx}}
{{- end}}
`)
	assert.Nil(t, err)
	b := &bytes.Buffer{}
	// 简介for...i,make创建一个指定大小的切片
	Slice := make([]int, 20)
	err = tpl.Execute(b, Slice)
	assert.Nil(t, err)
}
func TestIFELSE(t *testing.T) {
	type User struct {
		Age int
	}
	tpl := template.New("ifelse")
	tpl, err := tpl.Parse(`
{{- if and (gt .Age 0) (le .Age 6) -}}
儿童 0<age<6
{{- else if and (ge .Age 12) (le .Age 18) -}}
儿童 12<age<18
{{- else -}}
承认
{{- end -}}
`)
	assert.Nil(t, err)
	bfs := &bytes.Buffer{}
	tpl.Execute(bfs, User{Age: 12})
	assert.Equal(t, "儿童 12<age<18", bfs.String())
}

func TestGoTemplateEngine_Render(t *testing.T) {

	parseGlob, err := template.ParseGlob("testdata/tpls/*.gohtml")
	require.NoError(t, err)
	engine := &Render.GoTemplateEngine{T: parseGlob}
	server := NewHttpServer(ServerWithTemplateEngine(engine))
	server.Get("/user", func(ctx *Context) {
		// 默认是文件名，否则会报undefined
		err = ctx.Render("login.gohtml", nil)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println("hello")
	})
	server.Start(":9090")
}

func TestUpload(t *testing.T) {
	parseGlob, err := template.ParseGlob("testdata/tpls/*.gohtml")
	require.NoError(t, err)
	engine := &Render.GoTemplateEngine{T: parseGlob}
	server := NewHttpServer(ServerWithTemplateEngine(engine))
	server.Get("/upload", func(ctx *Context) {
		// 默认是文件名，否则会报undefined
		err = ctx.Render("upload.gohtml", nil)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println("hello")
	})
	fn := &FileUploader{
		FileField: "myfile",
		DstPathFunc: func(header *multipart.FileHeader) string {
			//fmt.Println(filepath.Join("testdata", "upload", uuid.New().String()))
			return filepath.Join("testdata", "upload", header.Filename)
		},
	}
	server.Post("/upload", fn.Handle())
	server.Start(":9090")
}

func TestDownLoader(t *testing.T) {
	server := NewHttpServer()
	// 浏览器输入localhost:9090/download?file=640.png，浏览器就会从testdata/download里面下载
	server.Get("/download", (&FileDownLoader{Dir: "./testdata/download"}).Handle())
	server.Start(":9090")
}

func TestStaticResource_Handle(t *testing.T) {
	server := NewHttpServer()
	s, err := NewStaticResource(filepath.Join("testdata", "static"))
	require.NoError(t, err)

	// 访问 localhost:9090/static/test.js
	server.Get("/static/:file", s.Handle)
	server.Start(":9090")
}

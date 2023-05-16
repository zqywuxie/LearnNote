// @Author: zqy
// @File:FileOperation.go
// @Date: 2023/5/16 10:31
// @Description 文件操作；目前文件操作使用oss，而不是直接与服务器操作，这样会比较安全

package customize

import (
	lru "github.com/hashicorp/golang-lru"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type FileUploader struct {
	FileField string

	// 存放路径交给用户自己处理
	// DstPathFunc header.Filename 获得文件名称
	DstPathFunc func(*multipart.FileHeader) string
}

type FileDownLoader struct {
	// Dir 下载文件所在地址
	Dir string
}

// StaticResource 处理静态资源
type StaticResource struct {
	dir string

	// 缓存
	cache *lru.Cache

	// 设置文件匹配策略
	extensionContextType map[string]string

	// 设置最大缓存
	maxSize int
}

type fileCacheItem struct {
	contentType string
	fileName    string
	fileSize    int
	data        []byte
}

// cacheFile 封装缓存方法
func (r *StaticResource) cacheFile(item *fileCacheItem) {
	if r.cache != nil && r.maxSize >= item.fileSize {
		r.cache.Add(item.fileName, item.data)
	}
}

func (r *StaticResource) writeItemAsResponse(item *fileCacheItem, writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", item.contentType)
	writer.Header().Set("Content-Length", string(item.data))
	writer.Write([]byte(item.data))
	writer.WriteHeader(http.StatusOK)
}

// 从缓存里面读取配置
func (r *StaticResource) readFileFromData(fileName string) (*fileCacheItem, bool) {
	if r.cache != nil {
		if value, ok := r.cache.Get(fileName); ok {
			return value.(*fileCacheItem), true
		}
	}
	return nil, false
}

// StaticResourceHandlerOption 用户自定义内容
type StaticResourceHandlerOption func(s *StaticResource)

func StaticWithCache(maxFileSizeThreshold int, maxCacheFileCnt int) StaticResourceHandlerOption {
	return func(s *StaticResource) {
		cache, err := lru.New(maxCacheFileCnt)
		if err != nil {
			log.Println("创建缓存失败，不会进行缓存数据")
		}
		s.cache = cache
		s.maxSize = maxFileSizeThreshold
	}
}

func WithMoreExtension(extMap map[string]string) StaticResourceHandlerOption {
	return func(s *StaticResource) {
		for key, value := range extMap {
			s.extensionContextType[key] = value
		}
	}
}

func NewStaticResource(dir string, opts ...StaticResourceHandlerOption) (*StaticResource, error) {
	s := &StaticResource{
		dir: dir,
		extensionContextType: map[string]string{
			"jpg":  "image/jpg",
			"png":  "image/png",
			"jpeg": "image/jpeg",
			"pdf":  "image/pdf",
			"jpe":  "image/jpe",
		}}
	for _, opt := range opts {
		opt(s)
	}
	return s, nil
}

// Handle 静态资源处理
func (r *StaticResource) Handle(ctx *Context) {
	fileName, err := ctx.PathValue("file")
	if err != nil {
		ctx.RespData = []byte("路径错误")
		ctx.RespCode = http.StatusBadRequest
		return
	}
	// 先从缓存里面查询，如果没有再进入磁盘查询

	if data, ok := r.readFileFromData(fileName); ok {
		r.writeItemAsResponse(data, ctx.Resp)
		return
	}

	// 从磁盘里面读取
	path := filepath.Join(r.dir, fileName)
	ext := filepath.Ext(path)[1:]
	contextType := r.extensionContextType[ext]
	file, err := os.ReadFile(path)
	if err != nil {
		ctx.RespData = []byte("服务器故障")
		ctx.RespCode = http.StatusInternalServerError
		return
	}

	item := &fileCacheItem{
		contentType: contextType,
		fileName:    fileName,
		fileSize:    len(file),
		data:        file,
	}
	r.cacheFile(item)
	r.writeItemAsResponse(item, ctx.Resp)
}

// 结合option使用
//func (u *FileUploader) HandleFunc(ctx *Context) {
//}

// Handle 在注册路由的时候 作为HandleFunc进行传入
func (u *FileUploader) Handle() HandleFunc {
	return func(ctx *Context) {
		// 1.获取http请求的数据
		file, header, err := ctx.Req.FormFile(u.FileField)
		//header.Filename
		if err != nil {
			ctx.RespCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败，原因：" + err.Error())
			return
		}
		// 2.找到存放路径
		defer file.Close()
		dstpath := u.DstPathFunc(header)

		// Stat 判断路径是否存在，如果不存在就进行创建目录
		if _, err = os.Stat(dstpath); err != nil {

			// filepath.Dir 获得目录路径
			// filepath.Base 文件名（包含扩展名）
			os.MkdirAll(filepath.Dir(dstpath), 0755)
		}
		openFile, err := os.OpenFile(dstpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
		if err != nil {
			ctx.RespCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败，原因：" + err.Error())
			return
		}
		defer openFile.Close()
		// 3.保存文件

		_, err = io.CopyBuffer(openFile, file, nil)
		if err != nil {
			ctx.RespCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败，原因：" + err.Error())
			return
		}
		// 4.返回响应
		ctx.RespCode = http.StatusOK
		ctx.RespData = []byte("上传成功")
	}
}

// Handle 文件下载
func (f FileDownLoader) Handle() HandleFunc {
	return func(ctx *Context) {
		req, err := ctx.QueryValue("file")
		if err != nil {
			ctx.RespCode = http.StatusBadRequest
			ctx.RespData = []byte("找不到目前文件")
			return
		}
		//filepath.Clean 函数的作用是返回等效的路径名，它会通过纠正路径中的错误或规范化路径分隔符等方式来实现。
		path := filepath.Join(f.Dir, filepath.Clean(req))
		base := filepath.Base(path)
		path, _ = filepath.Abs(path)
		if !strings.Contains(path, f.Dir) {
			ctx.RespCode = http.StatusGatewayTimeout
			ctx.RespData = []byte("错误请求")
			return
		}
		header := ctx.Resp.Header()
		// 设置响应头
		header.Set("Content-Disposition", "attachment;filename="+base)
		header.Set("Content-Description", "File Transfer")
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Content-Transfer-Encoding", "binary")
		header.Set("Expires", "0")
		header.Set("Cache-Control", "must-revalidate")
		header.Set("Pragma", "public")
		http.ServeFile(ctx.Resp, ctx.Req, path)
	}
}

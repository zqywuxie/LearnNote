// @Author: zqy
// @File:FileOperation.go
// @Date: 2023/5/16 10:31
// @Description todo

package customize

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type FileUploader struct {
	FileField string

	// 存放路径交给用户自己处理
	// DstPathFunc header.Filename 获得文件名称
	DstPathFunc func(*multipart.FileHeader) string
}

func (u FileUploader) Handle() HandleFunc {
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

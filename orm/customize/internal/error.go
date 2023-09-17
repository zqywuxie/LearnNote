package internal

import (
	"errors"
	"fmt"
)

var (
	ErrModelTypeSelect = errors.New("只支持一级指针或者结构体")
)

func NewErrorExpression(expr string) error {
	return fmt.Errorf("orm:不支持的表达式 %v", expr)
}
func NewErrorUnknownField(filed string) error {
	return fmt.Errorf("orm:未知的字段 %v", filed)
}

// ErrorResult
// @ErrorResult 返回错误信息
// 解决方案：xxx
//func ErrorResult(code int, Message string) error {
//
//}

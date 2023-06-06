// @Author: zqy
// @File: fields.go
// @Date: 2023/6/6 22:17
// @Description todo

package reflect

import (
	"errors"
	"reflect"
)

// IterateFields 遍历字段
// 注释里面说明 接收xxx类型的数据 （结构体，指针）
func IterateFields(entity any) (map[string]any, error) {
	if entity == nil {
		return nil, errors.New("空指针异常")
	}

	typeOf := reflect.TypeOf(entity)
	valueOf := reflect.ValueOf(entity)

	if valueOf.IsZero() {
		return nil, errors.New("不支持零值")
	}

	// 层层解引用
	for typeOf.Kind() == reflect.Pointer {
		// Elem 获得切片等内部值，或指针指向的值
		typeOf = typeOf.Elem()
		valueOf = valueOf.Elem()
	}

	if typeOf.Kind() != reflect.Struct {
		return nil, errors.New("非法类型")
	}
	// 类型数量
	field := typeOf.NumField()
	res := make(map[string]any, field)

	for i := 0; i < field; i++ {
		fieldType := typeOf.Field(i)
		fieldValue := valueOf.Field(i)
		if fieldType.IsExported() {
			res[fieldType.Name] = fieldValue.Interface()
		} else {
			res[fieldType.Name] = reflect.Zero(fieldType.Type).Interface()
		}
	}

	return res, nil
}

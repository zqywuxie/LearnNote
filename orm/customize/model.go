package customize

import (
	"GoCode/orm/customize/internal"
	"reflect"
	"unicode"
)

// @Description todo
// 反射解析模型
type Model struct {
	TableName string
	FiledMap  map[string]*Filed
}

type Filed struct {
	ColName string
}

var defauleModles = map[reflect.Type]*Model{}

var defaultRegister = &register{}

type register struct {
	models map[reflect.Type]*Model
}

// 先从注册中心获取
func (r *register) get(val any) (*Model, error) {
	typ := reflect.TypeOf(val)
	model, ok := r.models[typ]
	// 缓存中没有就进行parse
	if !ok {
		parseModel, err := r.ParseModel(val)
		if err != nil {
			return nil, err
		}
		r.models[typ] = parseModel
		return parseModel, err
	}
	return model, nil
}

// ParseModel 传结构体
// 限制用户输入一级指针或者结构体，简化开发
func (r *register) ParseModel(val any) (*Model, error) {
	types := reflect.TypeOf(val)
	if types.Kind() != reflect.Pointer && types.Kind() != reflect.Struct {
		return nil, internal.ErrModelTypeSelect
	}
	types = types.Elem()
	numField := types.NumField()
	filedMap := make(map[string]*Filed, numField)
	for i := 0; i < numField; i++ {
		field := types.Field(i)
		filedMap[field.Name] = &Filed{ColName: underscoreName(field.Name)}
	}
	return &Model{
		TableName: underscoreName(types.Name()),
		FiledMap:  filedMap,
	}, nil
}

// 驼峰转字符串
func underscoreName(tableName string) string {
	var buf []byte
	for i, v := range tableName {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}
	}
	return string(buf)
}

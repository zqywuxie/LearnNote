package customize

import (
	"GoCode/orm/customize/internal"
	"reflect"
	"strings"
	"sync"
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

const tagKeyColumn = "column"

var defauleModles = map[reflect.Type]*Model{}

var defaultRegister = &register{}

type register struct {
	// 读写锁
	//lock   sync.RWMutex
	//models map[reflect.Type]*Model
	//sync.Map
	models sync.Map
}

func newRegister() *register {
	return defaultRegister
}

type User struct {
	ID int `orm:"column=id,xxx=xxx"`
}

func (r *register) parseTags(tags reflect.StructTag) (map[string]string, error) {

	ormTag, ok := tags.Lookup("orm")
	if !ok {
		return map[string]string{}, nil
	}
	tagsValue := strings.Split(ormTag, ",")
	res := make(map[string]string, len(tagsValue))
	for _, tagValue := range tagsValue {
		segs := strings.Split(tagValue, "=")
		if len(segs) != 2 {
			return nil, internal.NewInvalidTagContent(tagValue)
		}
		key := segs[0]
		val := segs[1]
		res[key] = val
	}
	return res, nil
}

func (r *register) get(val any) (*Model, error) {
	typ := reflect.TypeOf(val)
	model, ok := r.models.Load(typ)
	if ok {
		return model.(*Model), nil
	}
	if !ok {
		var err error
		if model, err = r.ParseModel(val); err != nil {
			return nil, err
		}
	}
	r.models.Store(typ, model)
	return model.(*Model), nil
}

// 先从注册中心获取
//func (r *register) get(val any) (*Model, error) {
//	typ := reflect.TypeOf(val)
//	// 先读锁
//	r.lock.RLock()
//	model, ok := r.models[typ]
//	r.lock.RUnlock()
//	if ok {
//		return model, nil
//	}
//
//	//写锁
//	r.lock.Lock()
//	//为了确保在此期间没有其他并发操作修改了 r.models,
//	//所以需要再一次获取数据
//	model, ok = r.models[typ]
//	defer r.lock.Unlock()
//	if ok {
//		return model, nil
//	}
//
//	// 缓存中没有就进行parse
//	if !ok {
//		parseModel, err := r.ParseModel(val)
//		if err != nil {
//			return nil, err
//		}
//		r.models[typ] = parseModel
//		return parseModel, err
//	}
//	return model, nil
//}

// ParseModel 传结构体
// 限制用户输入一级指针或者结构体，简化开发
func (r *register) ParseModel(val any) (*Model, error) {
	types := reflect.TypeOf(val)
	if types.Kind() != reflect.Pointer && types.Kind() != reflect.Struct {
		return nil, internal.ErrModelTypeSelect
	}
	if types.Kind() == reflect.Pointer {
		types = types.Elem()
	}
	numField := types.NumField()
	filedMap := make(map[string]*Filed, numField)
	for i := 0; i < numField; i++ {
		field := types.Field(i)
		tags, err := r.parseTags(field.Tag)
		if err != nil {
			return nil, err
		}
		val := tags[tagKeyColumn]
		if val == "" {
			val = underscoreName(field.Name)
		}

		filedMap[field.Name] = &Filed{ColName: val}
	}

	var tableName string
	if name, ok := val.(TableName); ok {
		tableName = name.TableName()
	} else if tableName == "" {
		tableName = underscoreName(types.Name())
	}

	return &Model{
		TableName: tableName,
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

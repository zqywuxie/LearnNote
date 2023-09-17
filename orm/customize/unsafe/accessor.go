package unsafe

import (
	"errors"
	"reflect"
	"unsafe"
)

type UnsafeAccessor struct {
	fields  map[string]FieldMata
	address unsafe.Pointer
}

type FieldMata struct {
	Offset uintptr
	Type   reflect.Type
}

func NewUnsafeAccessor(entity any) UnsafeAccessor {
	typeof := reflect.TypeOf(entity)
	typeof = typeof.Elem()
	numField := typeof.NumField()
	filedMatas := make(map[string]FieldMata, numField)
	for i := 0; i < numField; i++ {
		field := typeof.Field(i)
		filedMatas[field.Name] = FieldMata{Offset: field.Offset, Type: field.Type}
	}
	valueOf := reflect.ValueOf(entity)
	return UnsafeAccessor{
		fields: filedMatas,
		//valueOf.UnsafeAddr() 不用这个是防止被GC等干扰导致出错
		// 使用指针就算地址出错，但还是会指向
		address: valueOf.UnsafePointer(),
	}
}

func (u UnsafeAccessor) Field(filed string) (any, error) {
	val, ok := u.fields[filed]
	if !ok {
		return nil, errors.New("未知字段")
	}
	// 获得了当前地址
	cur := unsafe.Pointer(val.Offset + uintptr(u.address))

	//return *(*int)(cur), nil
	// 一般来讲不会直到值的确切类型 通过上面的反射拿到type
	//reflect.New/NewAt 创建指针 所以使用Elem
	return reflect.NewAt(val.Type, cur).Elem().Interface(), nil
}

func (u UnsafeAccessor) SetField(filed string, value any) error {
	val, ok := u.fields[filed]
	if !ok {
		return errors.New("未知字段")
	}
	// 获得了当前地址
	cur := unsafe.Pointer(val.Offset + uintptr(u.address))

	//*(*int)(cur) = value.(int)

	reflect.NewAt(val.Type, cur).Elem().Set(reflect.ValueOf(value))

	//return *(*int)(cur), nil
	// 一般来讲不会直到值的确切类型 通过上面的反射拿到type
	//reflect.New/NewAt 创建指针 所以使用Elem
	return nil
}

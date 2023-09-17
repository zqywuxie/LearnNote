// @Author: zqy
// @File: func_call.go
// @Date: 2023/6/7 16:39
// @Description todo

// 方法类型的字段和方法，在反射上也是两个不同的东西

package reflect

import (
	"reflect"
)

type UserService struct {
	// 方法类型的字段
	GetById func()
}

// 方法
func (s UserService) getById() {

}

type FuncInfo struct {
	Name string
	// 方法输入参数类型
	InputTypes []reflect.Type
	// 方法输出参数类型
	OutPutTypes []reflect.Type

	//执行结果
	Result []any
}

func IterateFunc(entity any) (map[string]FuncInfo, error) {
	typ := reflect.TypeOf(entity)

	// 获得定义到结构体上的，指针上的得不到
	numMethod := typ.NumMethod()
	res := make(map[string]FuncInfo, numMethod)
	for i := 0; i < numMethod; i++ {
		method := typ.Method(i)
		Func := method.Func
		numIn := Func.Type().NumIn()
		numOut := Func.Type().NumOut()
		InputArgs := make([]reflect.Type, 0, numIn)

		// 获得输入参数信息
		InputValues := make([]reflect.Value, 0, numIn)

		// 所以对于结构体实现接口，参数是结构体本身,所以第一个参数是结构体本身
		// 否则下面只会遍历接入参数
		InputValues = append(InputValues, reflect.ValueOf(entity))
		InputArgs = append(InputArgs, reflect.TypeOf(entity))
		for i := 1; i < numIn; i++ {
			fnInputType := Func.Type().In(i)
			InputArgs = append(InputArgs, fnInputType)
			InputValues = append(InputValues, reflect.Zero(fnInputType))

		}
		// 获得输出参数信息
		OutPutArgs := make([]reflect.Type, 0, numOut)
		for i := 0; i < numOut; i++ {
			OutPutArgs = append(OutPutArgs, Func.Type().Out(i))
		}

		// 调用该方法
		resValues := Func.Call(InputValues)

		result := make([]any, 0, len(resValues))

		// value类型的真实值
		for _, v := range resValues {
			result = append(result, v.Interface())
		}

		res[method.Name] = FuncInfo{
			Name:        method.Name,
			InputTypes:  InputArgs,
			OutPutTypes: OutPutArgs,
			Result:      result,
		}
	}
	return res, nil
}

package unsafe

import (
	"fmt"
	"reflect"
)

func PrintFiledOffset(entity any) {
	typeOf := reflect.TypeOf(entity)
	numField := typeOf.NumField()
	for i := 0; i < numField; i++ {
		filed := typeOf.Field(i)
		fmt.Println(filed.Offset)
	}
}

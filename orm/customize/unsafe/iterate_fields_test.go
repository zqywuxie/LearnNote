// @Author: zqy
// @File: iterate_fields_test.go.go
// @Date: 2023/9/17 15:43
// @Description todo

package unsafe

import "testing"

func TestPrintFiledOffset(t *testing.T) {
	type args struct {
		entity any
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "offset",
			args: args{entity: User{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintFiledOffset(tt.args.entity)
		})
	}
}

type User struct {
	Name string
	Age  int32
	//Alias int32
	Hello string
}

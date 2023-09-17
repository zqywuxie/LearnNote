// @Author: zqy
// @File: func_call_test.go.go
// @Date: 2023/6/7 16:43
// @Description todo

package reflect

import (
	"GoCode/orm/customize/types"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestIterateFunc(t *testing.T) {
	tests := []struct {
		name    string
		entity  any
		wantErr error
		wantRes map[string]FuncInfo
	}{
		{
			name:   "struct",
			entity: types.NewUser("zqy", 12),
			wantRes: map[string]FuncInfo{
				"GetAge": {
					Name: "GetAge",
					InputTypes: []reflect.Type{
						reflect.TypeOf(types.User{}),
					},
					OutPutTypes: []reflect.Type{
						reflect.TypeOf(0),
					},
					Result: []any{
						12,
					},
				},
				//"ChangeName": {
				//	Name: "ChangeName",
				//	InputTypes: []reflect.Type{
				//		reflect.TypeOf(""),
				//	},
				//	OutPutTypes: nil,
				//	Result:      nil,
				//},
			},
		},
		{
			name:   "pointer",
			entity: types.NewUserPtr("zqy", 12),
			wantRes: map[string]FuncInfo{
				"GetAge": {
					Name: "GetAge",
					InputTypes: []reflect.Type{
						reflect.TypeOf(&types.User{}),
					},
					OutPutTypes: []reflect.Type{
						reflect.TypeOf(0),
					},
					Result: []any{
						12,
					},
				},
				"ChangeName": {
					Name: "ChangeName",
					InputTypes: []reflect.Type{
						reflect.TypeOf(&types.User{}),
						reflect.TypeOf(""),
					},
					OutPutTypes: []reflect.Type{},
					Result:      []any{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := IterateFunc(tt.entity)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}

			assert.Equal(t, tt.wantRes, res)

		})
	}
}

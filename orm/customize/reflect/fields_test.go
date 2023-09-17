// @Author: zqy
// @File: fields_test.go.go
// @Date: 2023/6/6 22:18
// @Description todo

package reflect

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type User struct {
	Name string
	age  int
}

func TestIterateFields(t *testing.T) {

	testCases := []struct {
		name    string
		entity  any
		wantErr error
		wantRes map[string]any
	}{
		{
			name:    "hello",
			entity:  User{Name: "Tom", age: 12},
			wantErr: nil,
			wantRes: map[string]any{
				"Name": "Tom",
				"age":  0,
			},
		},
		{
			name:    "hello",
			entity:  User{Name: "Tom", age: 12},
			wantErr: nil,
			wantRes: map[string]any{
				"Name": "Tom",
				// 私有变量 设置为零值
				"age": 0,
			},
		},
		{
			name:    "basic type",
			entity:  12,
			wantErr: errors.New("非法类型"),
		},

		{
			name: "multiple pointer",
			entity: func() **User {
				res := &User{Name: "Tom", age: 12}
				return &res
			}(),
			//wantErr: nil,
			wantRes: map[string]any{
				"Name": "Tom",
				// 私有变量 设置为零值
				"age": 0,
			},
		},
		{
			name:    "nil",
			entity:  nil,
			wantErr: errors.New("空指针异常"),
		},
		{
			name:    "user nil",
			entity:  (*User)(nil),
			wantErr: errors.New("不支持零值"),
			//wantRes: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fields, err := IterateFields(tc.entity)
			assert.Equal(t, err, tc.wantErr)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, fields)
		})
	}
}

func TestSetField(t *testing.T) {

	testCases := []struct {
		name     string
		entity   any
		field    string
		newValue any

		wantErr error

		wantEntity any
	}{
		{
			name: "struct",
			entity: User{
				Name: "Tom",
				age:  18,
			},
			wantEntity: User{
				Name: "ZQY",
				age:  18,
			},
			field:    "Name",
			newValue: "ZQY",
			wantErr:  errors.New("该字段不可被修改"),
		},
		{
			name: "pointer",
			entity: &User{
				Name: "Tom",
				age:  18,
			},
			wantEntity: &User{
				Name: "ZQY",
				age:  18,
			},
			field:    "Name",
			newValue: "ZQY",
			//wantErr:  errors.New("该字段不可被修改"),
		},
		{
			name: "private field",
			entity: &User{
				age: 18,
			},
			wantEntity: &User{
				age: 12,
			},
			field:    "age",
			newValue: 12,
			//wantErr:  errors.New("该字段不可被修改"),
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := SetField(tt.entity, tt.field, tt.newValue)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantEntity, tt.entity)
		})
	}

	//var i = 0
	//iptr := &i
}

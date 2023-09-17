// @Author: zqy
// @File: model_test.go.go
// @Date: 2023/9/17 10:50
// @Description todo

package reflect

import (
	"GoCode/orm/customize"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestModel struct {
	Id        int64
	FirstName string
	Age       int
	LastName  string
}

func Test_parseModel(t *testing.T) {
	tests := []struct {
		name    string
		args    any
		want    *customize.Model
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
		{
			name: "struct",
			args: TestModel{},
			want: &customize.Model{
				TableName: "test_model",
				FiledMap: map[string]*customize.Filed{
					"Id": {
						ColName: "id",
					},
					"FirstName": {
						ColName: "first_name",
					},
					"Age": {
						"age",
					},
					"LastName": {
						"last_name",
					},
				},
			},
		},
		{
			name: "pointer",
			args: &TestModel{},
			want: &customize.Model{
				TableName: "test_model",
				FiledMap: map[string]*customize.Filed{
					"Id": {
						ColName: "id",
					},
					"FirstName": {
						ColName: "first_name",
					},
					"Age": {
						"age",
					},
					"LastName": {
						"last_name",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := customize.ParseModel(tt.args)
			if err != nil {
				panic(err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

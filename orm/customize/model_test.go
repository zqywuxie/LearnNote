// @Author: zqy
// @File: model_test.go.go
// @Date: 2023/9/19 12:58
// @Description todo

package customize

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_register_ParseModel(t *testing.T) {
	tests := []struct {
		name    string
		entity  any
		want    *Model
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "tags",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column=first-name"`
				}
				return &TagTable{}
			}(),
			want: &Model{
				TableName: "tag_table",
				FiledMap: map[string]*Filed{
					"FirstName": &Filed{ColName: "first-name"},
				},
			},
			wantErr: nil,
		},
	}
	r := newRegister()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			get, err := r.get(tt.entity)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, get)
			typ := reflect.TypeOf(tt.entity)
			value, ok := r.models.Load(typ)
			assert.True(t, ok)
			assert.Equal(t, tt.want, value)
		})
	}
}

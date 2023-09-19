// @Author: zqy
// @File: model_test.go.go
// @Date: 2023/9/19 12:58
// @Description todo

package customize

import (
	"GoCode/orm/customize/internal"
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
		{
			name: "empty tag",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column"`
				}
				return &TagTable{}
			}(),
			wantErr: internal.NewInvalidTagContent("column"),
		},
		{
			name: "default columnName",
			entity: func() any {
				type TagTable struct {
					FirstName string
				}
				return &TagTable{}
			}(),
			want: &Model{
				TableName: "tag_table",
				FiledMap: map[string]*Filed{
					"FirstName": &Filed{ColName: "first_name"},
				},
			},
			wantErr: nil,
		},
		{
			name:   "custom tableName",
			entity: CustomTableName{},
			want: &Model{
				TableName: "custom_table_name_t",
				FiledMap: map[string]*Filed{
					"Name": &Filed{ColName: "name"},
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
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, get)
			typ := reflect.TypeOf(tt.entity)
			value, ok := r.models.Load(typ)
			assert.True(t, ok)
			assert.Equal(t, tt.want, value)
		})
	}

}

type CustomTableName struct {
	Name int `orm:"column=name"`
}
type CustomTableNamePtr struct {
}

func (c CustomTableName) TableName() string {
	return "custom_table_name_t"
}

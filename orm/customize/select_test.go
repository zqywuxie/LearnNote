// @Author: zqy
// @File: select_test.go
// @Date: 2023/6/5 16:19
// @Description todo

package customize

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelector_Build(t *testing.T) {
	db, err := NewDB()
	if err != nil {
		panic(err)
	}
	testCases := []struct {
		name      string
		builder   QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name:    "selector without from",
			builder: newSelector[TestModel](db),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model`;",
				args: nil,
			},
			wantErr: nil,
		},
		{
			name:    "selector with from",
			builder: (&Selector[TestModel]{}).From("TestModel"),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel`;",
				args: nil,
			},
			wantErr: nil,
		},
		{
			name: "selector with from",

			builder: (&Selector[TestModel]{}).From("test.model"),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test`.`model`;",
				args: nil,
			},
			wantErr: nil,
		},
		{
			name: "selector with where",

			builder: (&Selector[TestModel]{}).Where(C("id").Eq(12)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE `id` = ?;",
				args: []any{12},
			},
			wantErr: nil,
		},
		{
			name: "selector with Not",

			builder: (&Selector[TestModel]{}).Where(Not(C("Id").Eq(12))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE  NOT (`id` = ?);",
				args: []any{12},
			},
			wantErr: nil,
		},
		{
			name: "selector with And",

			builder: (&Selector[TestModel]{}).Where(C("id").Eq(12).And(C("Age").Eq(12))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE (`id` = ?) AND (`Age` = ?);",
				args: []any{12, 12},
			},
			wantErr: nil,
		},
		{
			name: "selector with OR",

			builder: (&Selector[TestModel]{}).Where(C("id").Eq(12).Or(C("Age").Eq(12))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE (`id` = ?) OR (`Age` = ?);",
				args: []any{12, 12},
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.builder.Build()
			assert.Equal(t, err, tc.wantErr)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

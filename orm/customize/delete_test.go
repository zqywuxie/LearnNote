// @Author: zqy
// @File: delete_test.go.go
// @Date: 2023/9/17 14:33
// @Description todo

package customize

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleter_Build(t *testing.T) {
	testCases := []struct {
		name      string
		builder   QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name:    "simple delete",
			builder: (&Deleter[TestModel]{}),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model`;",
				args: nil,
			},
			wantErr: nil,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actul, err := tt.builder.Build()
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantQuery, actul)
		})
	}
}

// @Author: zqy
// @File: iterate_test.go.go
// @Date: 2023/9/17 10:06
// @Description todo

package reflect

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIterateArray(t *testing.T) {
	tests := []struct {
		name    string
		args    any
		want    []any
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
		{
			name:    "Array",
			args:    [4]int{1, 2, 3, 4},
			want:    []any{1, 2, 3, 4},
			wantErr: nil,
		}, {
			name:    "Slice",
			args:    []int{1, 2, 3, 4},
			want:    []any{1, 2, 3, 4},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IterateArrayOrSlice(tt.args)
			if err != nil {
				panic(err)
			}
			assert.Equalf(t, tt.want, got, "IterateArrayOrSlice(%v)", tt.args)
		})
	}
}

func TestIterateMap(t *testing.T) {
	tests := []struct {
		name    string
		args    any
		want    []any
		want1   []any
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
		{
			name: "Map",
			args: map[string]any{"你好": 123, "hello": "ok"},
			want: []any{
				"你好", "hello",
			},
			want1: []any{
				123, "ok",
			},
		},
		{
			name: "非Map",
			args: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := IterateMap(tt.args)
			if err != nil {
				panic(err)
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

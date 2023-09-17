// @Author: zqy
// @File: accessor_test.go.go
// @Date: 2023/9/17 16:34
// @Description todo

package unsafe

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUnsafeAccessor_Field(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}
	u := &User{
		Name: "123",
		Age:  123,
	}
	accessor := NewUnsafeAccessor(u)
	field, err := accessor.Field("Age")
	require.NoError(t, err)
	assert.Equal(t, 123, field)
	err = accessor.SetField("Age", 12)
	require.NoError(t, err)

}

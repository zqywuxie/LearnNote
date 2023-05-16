// @Author: zqy
// @File: file_test.go
// @Date: 2023/5/16 10:14
// @Description todo

package file_demo

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestFile(t *testing.T) {
	fmt.Println(os.Getwd())
	open, err := os.Open("testdata/my_file.txt")
	os.Create()
	require.NoError(t, err)
	bytes := make([]byte, 64)
	n, err := open.Read(bytes)
	require.NoError(t, err)
	fmt.Println(n, string(bytes[:n]))

	_, err = open.WriteString("你好")
	fmt.Println(err)
	open.Close()
}

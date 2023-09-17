// @Author: zqy
// @File: mapS.go
// @Date: 2023/9/11 17:18
// @Description
// 如何扩充系统类型或者别人的类型
// 1. 定义别名
// type Queue []int
// 2. 使用组合
// 3. 内嵌 省下代码
//type NewNode struct {
//*tree.MyTreeNode //embedding 内嵌
//}

package main

import "fmt"

// "ab-cda"
func lengthOfNonRepeatingSubStr(s string) int {
	lastOccurred := make(map[rune]int)
	start := 0
	maxLength := 0
	for i, ch := range []rune(s) {
		if lastI, ok := lastOccurred[ch]; ok && lastI >= start {
			start = lastOccurred[ch] + 1
		}
		if i-start+1 > maxLength {
			maxLength = i - start + 1
		}
		lastOccurred[ch] = i

	}
	return maxLength
}

func main() {
	//var myMAP []int
	fmt.Println(lengthOfNonRepeatingSubStr("abcdea"))
	fmt.Println(lengthOfNonRepeatingSubStr("a"))
	fmt.Println(lengthOfNonRepeatingSubStr("aaa"))
	fmt.Println(lengthOfNonRepeatingSubStr("avced"))
	fmt.Println(lengthOfNonRepeatingSubStr("你好啊啊啊啊啊你好啊"))
	fmt.Println(lengthOfNonRepeatingSubStr("欧妮基瓦"))
}

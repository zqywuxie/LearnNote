// @Author: zqy
// @File: user.go
// @Date: 2023/6/7 16:49
// @Description todo

package types

import "fmt"

type User struct {
	Name string
	age  int
}

func NewUserPtr(name string, age int) *User {
	return &User{Name: name, age: age}
}

func NewUser(name string, age int) User {
	return User{Name: name, age: age}
}

// 与下面那个等价，所以输入参数为User
//func GetAge(u User) int {
//	return u.age
//}

func (u User) GetAge() int {
	return u.age
}
func (u *User) ChangeName(newName string) {
	u.Name = newName
}

func (u User) private() {
	fmt.Println("private")
}

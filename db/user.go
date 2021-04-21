package db

import (
	"fmt"
)

// 对用户建模
type User struct {
	Id   int    `json:"id" binding:"required"`
	Psw  string `json:"psw" binding:"required"`
	Name string `json:"name"`
}

// 对比结构体当前值与数据库数据
func (u *User) Login() bool {
	var p, n string
	err := u_login.QueryRow(u.Id).Scan(&p, &n)
	return err == nil && p == u.Psw
}

func (u *User) Print() {
	fmt.Printf("id=%d, psw=%s, name=%s\n", u.Id, u.Psw, u.Name)
}

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

// 登录, 对比结构体当前值与数据库数据
func (u *User) Login() bool {
	var p, n string
	err := u_login.QueryRow(u.Id).Scan(&p, &n)
	return err == nil && p == u.Psw
}

// 当前结构体是否与uid为好友
func (u *User) HasFriend(frid int) bool {
	if u.Id == frid {
		return true
	}
	var r int
	err := u_has_fr.QueryRow(u.Id, frid).Scan(&r)
	return err == nil && r == frid
}

// 当前结构体与uid添加好友
// 成功返回true,nil
// 失败返回false;数据库出错返回err,否则err=nil
func (u *User) AddFriend(frid int) (suss bool, err error) {
	if u.Id == frid {
		return
	}
	res, err := u_add_fr.Exec(u.Id, frid)
	if err != nil {
		return
	}
	i, err := res.RowsAffected()
	if i != 1 || err != nil {
		return
	}
	suss = true
	return
}

func (u *User) Print() {
	fmt.Printf("id=%d, psw=%s, name=%s\n", u.Id, u.Psw, u.Name)
}

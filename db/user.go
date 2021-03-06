package db

import (
	"database/sql"
	"log"
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
	if err == nil && p == u.Psw {
		u.Psw = p
		u.Name = n
		return true
	}
	return false
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

// 加载昵称
func (u *User) GetName() bool {
	err := u_get_name.QueryRow(u.Id).Scan(&u.Name)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// 下线时, 存储好友请求
func (u *User) StoreFrReq(frid int) (err error) {
	res, err := set_fr_req.Exec(u.Id, frid)
	if err != nil {
		return
	}
	i, err := res.RowsAffected()
	if i != 1 || err != nil {
		return
	}
	return
}

// 上线时, 删除所有暂存好友请求/回复
func (u *User) DelFrNotice() (rows int) {
	res, err := del_fr_ntc.Exec(u.Id)
	if err != nil {
		return
	}
	i, err := res.RowsAffected()
	if err != nil {
		return
	}
	rows = int(i)
	return
}

// 下线时, 存储好友回复
func (u *User) StoreFrAns(frid int, ans bool, conv_id int) (err error) {
	res, err := set_fr_ans.Exec(u.Id, frid, ans, conv_id)
	if err != nil {
		return
	}
	i, err := res.RowsAffected()
	if i != 1 || err != nil {
		return
	}
	return
}

// 下线时, 存储第一条消息id, 其余消息不影响
func (u *User) StoreOfflineMsg(msg_id int) (err error) {
	_, err = store_fmi.Exec(msg_id, u.Id)
	return
}

// 获取第一条离线消息的msg_id, 若有离线消息，设置has_set_fmi=0
func (u *User) GetOffMsgID() (msg_id int, suss bool) {
	var fmi sql.NullInt32
	err := get_fmi.QueryRow(u.Id).Scan(&fmi)
	if err == sql.ErrNoRows && !fmi.Valid {
		return
	}
	msg_id = int(fmi.Int32)

	row, err := set_has_fmi.Exec(u.Id)
	if err != nil {
		log.Println(err)
		return
	}
	i, err := row.RowsAffected()
	if err != nil {
		log.Println(err)
		return
	}
	if i != 1 {
		log.Println("expected RowsAffected 1, but actual", i)
		return
	}
	suss = true
	return
}

func (u *User) Print() {
	log.Printf("id=%d, psw=%s, name=%s\n", u.Id, u.Psw, u.Name)
}

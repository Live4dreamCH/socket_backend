package db

import "log"

func CreateConv(u1, u2 int) (conv_id int, err error) {
	res, err := new_conv.Exec()
	if err != nil {
		return
	}
	i, err := res.LastInsertId()
	if err != nil {
		return
	}
	conv_id = int(i)
	_, err = add_conv_mem.Exec(conv_id, u1)
	if err != nil {
		return
	}
	_, err = add_conv_mem.Exec(conv_id, u2)
	return
}

func GetOtherConvMems(uid, conv_id int) (mems []int, err error) {
	rows, err := get_conv_mems.Query(conv_id, uid, uid)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var uid int
		err = rows.Scan(&uid)
		if err != nil {
			return
		}
		mems = append(mems, uid)
	}

	return
}

type WsMsg struct {
	Conv_id int    `json:"conv_id"`
	Sender  int    `json:"sender"`
	Time    string `json:"time"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

func (msg *WsMsg) Save() (msg_id int, suss bool) {
	res, err := save_content.Exec(0, msg.Content)
	if err != nil {
		log.Println(err)
		return
	}
	r, err := res.RowsAffected()
	if r != 1 || err != nil {
		log.Println(err)
		return
	}
	con_id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return
	}

	res, err = save_msg.Exec(msg.Sender, msg.Time, con_id, msg.Conv_id)
	if err != nil {
		log.Println(err)
		return
	}
	r, err = res.RowsAffected()
	if r != 1 || err != nil {
		log.Println(err)
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return
	}
	return int(id), true
}

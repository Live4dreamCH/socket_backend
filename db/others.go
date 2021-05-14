package db

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

func (msg *WsMsg) Save() {
	save_content.Exec()
}

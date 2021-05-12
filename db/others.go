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
	rows, err := get_conv_mems.Query(conv_id, uid)
	if err != nil {
		return
	}
	for rows.Next() {
		var uid int
		err = rows.Scan(&uid)
		if err != nil {
			return
		}
		mems = append(mems, uid)
	}
	rows.Close()

	return
}

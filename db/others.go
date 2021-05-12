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

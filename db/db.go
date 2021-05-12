// 对后端事物建模, 并且使用嵌入式sql, 与数据库交换数据
package db

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// type Orm interface {
// 	Load() error
// 	Store() error
// }

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// 数据库
var dbp *sql.DB

// 预编译语句
var (
	u_login    *sql.Stmt
	u_get_name *sql.Stmt

	u_has_fr *sql.Stmt
	u_add_fr *sql.Stmt

	set_fr_req *sql.Stmt
	set_fr_ans *sql.Stmt
	get_fr_req *sql.Stmt
	get_fr_ans *sql.Stmt
	del_fr_ntc *sql.Stmt

	set_fmi     *sql.Stmt
	set_has_fmi *sql.Stmt
	get_has_fmi *sql.Stmt
	get_fmi     *sql.Stmt

	new_conv     *sql.Stmt
	add_conv_mem *sql.Stmt
)

func init() {
	s, err := os.ReadFile("../pwd/local_mysql.txt")
	check(err)
	psw := string(s)

	dbp, err = sql.Open("mysql", "root:"+psw+"@/chat?charset=utf8")
	check(err)
	err = dbp.Ping()
	check(err)

	u_login, err = dbp.Prepare(
		`select psw, u_name
		from users
		where u_id=?;`)
	check(err)

	u_get_name, err = dbp.Prepare(
		`select u_name
		from users
		where u_id=?;`)
	check(err)

	u_has_fr, err = dbp.Prepare(
		`select fr_id
		from friends
		where my_id=? and fr_id=?;`)
	check(err)

	u_add_fr, err = dbp.Prepare(
		`insert into friends (my_id,fr_id)
		values (?, ?);`)
	check(err)

	set_fr_req, err = dbp.Prepare(
		`insert into fr_notices (u_id, fr_id, is_ans)
		values (?, ?, 0);`)
	check(err)

	set_fr_ans, err = dbp.Prepare(
		`insert into fr_notices (u_id, fr_id, is_ans, ans, conv_id)
		values (?, ?, 1, ?, ?);`)
	check(err)

	get_fr_req, err = dbp.Prepare(
		`select fr_id
		from fr_notices
		where u_id=? and is_ans=0;`)
	check(err)

	get_fr_ans, err = dbp.Prepare(
		`select fr_id, ans
		from fr_notices
		where u_id=? and is_ans=1;`)
	check(err)

	del_fr_ntc, err = dbp.Prepare(
		`delete from fr_notices
		where u_id=?;`)
	check(err)

	set_fmi, err = dbp.Prepare(
		`insert into users (first_msg_id)
		values (?);`)
	check(err)

	set_has_fmi, err = dbp.Prepare(
		`insert into users (has_set_fmi)
		values (?);`)
	check(err)

	get_has_fmi, err = dbp.Prepare(
		`select has_set_fmi
		from users
		where u_id=?;`)
	check(err)

	get_fmi, err = dbp.Prepare(
		`select first_msg_id
		from users
		where u_id=?;`)
	check(err)

	new_conv, err = dbp.Prepare(
		`insert into convs (is_group)
		values (0);`)
	check(err)

	add_conv_mem, err = dbp.Prepare(
		`insert into conv_members (conv_id, mem_id)
		values (?,?);`)
	check(err)
}

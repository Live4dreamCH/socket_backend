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

// 严厉检查，让问题在启动时得以发现
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

	store_fmi   *sql.Stmt
	get_fmi     *sql.Stmt
	set_has_fmi *sql.Stmt

	new_conv     *sql.Stmt
	add_conv_mem *sql.Stmt
	//其它成员
	get_conv_mems *sql.Stmt

	Get_fr_list       *sql.Stmt
	Get_conv_list     *sql.Stmt
	Get_conv_mem_list *sql.Stmt

	save_content *sql.Stmt
	save_msg     *sql.Stmt
	Load_msg     *sql.Stmt
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

	// 下线时存储第一条离线消息id
	store_fmi, err = dbp.Prepare(
		`update users
		set first_msg_id = ?,
			has_set_fmi = 1
		where has_set_fmi = 0 and u_id = ?;`)
	check(err)
	// 上线时查询第一条离线消息id
	get_fmi, err = dbp.Prepare(
		`select first_msg_id
		from users
		where u_id = ? and has_set_fmi = 1;`)
	check(err)
	// 上线时修改has_set_fmi
	set_has_fmi, err = dbp.Prepare(
		`update users
		set has_set_fmi = 0
		where u_id = ? and has_set_fmi = 1;`)
	check(err)

	new_conv, err = dbp.Prepare(
		`insert into convs (is_group)
		values (0);`)
	check(err)
	add_conv_mem, err = dbp.Prepare(
		`insert into conv_members (conv_id, mem_id)
		values (?,?);`)
	check(err)
	get_conv_mems, err = dbp.Prepare(
		`select cm1.mem_id
		from conv_members cm1
		where cm1.conv_id = ? and cm1.mem_id != ? and cm1.conv_id in (
			select cm2.conv_id
			from conv_members cm2
			where cm2.mem_id = ?
		);`)
	check(err)

	Get_fr_list, err = dbp.Prepare(
		`select u.u_id, u.u_name
		from friends fr, users u
		where fr.my_id = ? and fr.fr_id = u.u_id;`)
	check(err)
	Get_conv_list, err = dbp.Prepare(
		`select cm1.conv_id, u.u_name
		from conv_members cm1, users u
		where cm1.mem_id != ? and u.u_id = cm1.mem_id and cm1.conv_id in (
			select cm2.conv_id
			from conv_members cm2
			where cm2.mem_id=?
		);`)
	check(err)
	Get_conv_mem_list, err = dbp.Prepare(
		`select u.u_id, u.u_name
		from conv_members cm1, users u
		where cm1.conv_id = ? and cm1.mem_id = u.u_id and cm1.conv_id in (
			select cm2.conv_id
			from conv_members cm2
			where cm2.mem_id = ?
		);`)
	check(err)

	save_content, err = dbp.Prepare(
		`insert into contents(con_type, con)
		values(?, ?);`)
	check(err)
	save_msg, err = dbp.Prepare(
		`insert into msgs(sender, msg_time, con_id, conv_id)
		values(?, ?, ?, ?);`)
	check(err)
	// 选取消息: 用户uid所在的所有会话里, 发送者不为用户uid的, 编号大于等于msg_id的, 消息以及消息内容
	Load_msg, err = dbp.Prepare(
		`select m.sender, m.msg_time, m.conv_id, c.con_type, c.con
		from msgs m, contents c, conv_members cm
		where cm.mem_id = ? and cm.conv_id = m.conv_id and m.sender != ? and m.msg_id >= ? and m.con_id = c.con_id;`)
	check(err)
}

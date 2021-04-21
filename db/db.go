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
var u_login *sql.Stmt

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

}

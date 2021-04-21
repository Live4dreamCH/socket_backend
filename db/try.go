package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Example() {
	s, err := os.ReadFile("../pwd/local_mysql.txt")
	checkErr(err)
	psw := string(s)

	db, err := sql.Open("mysql", "root:"+psw+"@/chat?charset=utf8")
	checkErr(err)

	//插入数据
	stmt, err := db.Prepare(
		`insert into users (psw, u_name) 
		values (?, ?);`)
	checkErr(err)

	res, err := stmt.Exec("fens", "我是谁")
	stmt.Close()
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)

	//更新数据
	stmt, err = db.Prepare(
		`update users 
		set u_name=? 
		where u_id=?`)
	checkErr(err)

	res, err = stmt.Exec("哈哈哈", id)
	stmt.Close()
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	//查询数据
	rows, err := db.Query("SELECT * FROM users")
	checkErr(err)

	for rows.Next() {
		var uid int
		var psw string
		var username string
		err = rows.Scan(&uid, &psw, &username)
		checkErr(err)
		fmt.Println("uid=", uid, "psw=", psw, "name=", username)
	}

	//删除数据
	stmt, err = db.Prepare(
		`delete from users 
		where u_id=?`)
	checkErr(err)

	res, err = stmt.Exec(id)
	stmt.Close()
	checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	db.Close()

}

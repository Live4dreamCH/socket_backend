//服务器主程序
package main

import (
	"github.com/gin-gonic/gin"
)

var sess sID

func init() {
	sess.m = make(map[string]int)
}

func main() {
	r := gin.Default()
	r.POST("/login", login)
	r.POST("/addfriend", addfriend)
	go wsStarter()
	err := r.Run(":43851")
	if err != nil {
		panic(err)
	}
}

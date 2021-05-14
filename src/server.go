//服务器主程序
package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

// 全局变量
var sess sID

var msgRouter wsRouter
var fileRouter wsRouter

// 初始化全局变量
func init() {
	log.SetPrefix("LOG: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	sess.m = make(map[string]int)
	msgRouter.m = make(map[int]*wsLink)
	fileRouter.m = make(map[int]*wsLink)
}

func main() {
	// rand.Seed(1)
	r := gin.Default()
	r.POST("/login", login)
	r.POST("/addfriend", addfriend)
	r.POST("/resfriend", resfriend)
	r.POST("/name", nameService)
	r.POST("/friendlist", friendList)
	r.POST("/convlist", convList)
	r.POST("/convmemlist", convMemList)
	go wsStarter()
	err := r.Run(":43851")
	if err != nil {
		panic(err)
	}
}

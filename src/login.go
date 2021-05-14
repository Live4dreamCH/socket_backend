package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Live4dreamCH/socket_backend/db"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 登录
func login(c *gin.Context) {
	var u db.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(400, gin.H{"res": "NO", "reason": "json bind error"})
		return
	}

	if u.Login() {
		sid := sess.set(u.Id)
		c.JSON(http.StatusOK, gin.H{"res": "OK", "sid": sid, "name": u.Name})
		log.Println("user", u.Id, "http login suss")
	} else {
		c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "wrong password"})
	}
}

// ws登录请求
type wl struct {
	Op  string `binding:"required"`
	Seq int    `binding:"required"`
	Sid string `binding:"required"`
}

// ws登录逻辑
func wsLogin(conn *websocket.Conn) (suss bool, uid int) {
	var req wl
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	for req.Op != "login" {
		err := conn.ReadJSON(&req)
		if err != nil {
			conn.WriteJSON(gin.H{"op": "login", "ack": req.Seq, "res": "NO", "reason": err})
			return
		}
	}
	conn.SetReadDeadline(time.Time{})

	uid, err := sess.get(req.Sid)
	if err != nil {
		conn.WriteJSON(gin.H{"op": "login", "ack": req.Seq, "res": "NO", "reason": "sid wrong"})
		return
	}
	conn.WriteJSON(gin.H{"op": "login", "ack": req.Seq, "res": "OK"})
	suss = true
	return
}

// todo: 登陆时推送所有下线时产生的消息、好友请求与回复
func notice() {

}

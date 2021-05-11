package main

import (
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
	} else {
		if u.Login() {
			sid := sess.set(u.Id)
			c.JSON(http.StatusOK, gin.H{"res": "OK", "sid": sid, "name": u.Name})
		} else {
			c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "wrong password"})
		}
	}
}

// ws登录
type wl struct {
	Op  string `binding:"required"`
	Seq int    `binding:"required"`
	Sid string `binding:"required"`
}

func wsLogin(conn *websocket.Conn) (suss bool, uid int) {
	var req wl
	for req.Op != "login" {
		conn.SetReadDeadline(time.Now().Add(time.Second * 5))
		err := conn.ReadJSON(&req)
		if err != nil {
			conn.WriteJSON(gin.H{"op": "login", "ack": req.Seq, "res": "NO"})
			return
		}
	}
	conn.SetReadDeadline(time.Time{})

	uid, err := sess.get(req.Sid)
	if err != nil {
		conn.WriteJSON(gin.H{"op": "login", "ack": req.Seq, "res": "NO"})
		return
	}
	conn.WriteJSON(gin.H{"op": "login", "ack": req.Seq, "res": "OK"})
	suss = true
	return
}

// todo: 登陆时推送所有下线时产生的消息
func notice() {

}

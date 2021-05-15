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

// todo: 登陆时推送所有下线时产生的好友请求与回复
func msgNotice(uid int) {
	// 找到发送的ws连接
	msgRouter.l.RLock()
	link, ok := msgRouter.m[uid]
	msgRouter.l.RUnlock()
	if !ok {
		log.Println("want to send Offline Msgs to user", uid, ", but conn already close")
		return
	}

	// 获取第一条离线消息id
	u := db.User{Id: uid}
	fmi, suss := u.GetOffMsgID()
	if !suss {
		log.Println("user", uid, "login, but has no offline msgs")
		return
	}

	// 执行查询数据库语句, 获取离线消息
	rows, err := db.Load_msg.Query(uid, uid, fmi)
	if err != nil {
		return
	}
	defer rows.Close()

	// 读取, 发送离线消息
	var pkg struct {
		wsMain
		db.WsMsg
	}
	pkg.Op = "msg"
	var ty int
	num := 0

	for rows.Next() {
		err = rows.Scan(&(pkg.Sender), &(pkg.Time), &(pkg.Conv_id), &ty, &(pkg.Content))
		if err != nil {
			log.Println(err)
			return
		}
		switch ty {
		case 0:
			pkg.Type = "text"
		}

		link.l.Lock()
		pkg.Seq = link.seq_num
		err = link.conn.WriteJSON(pkg)
		link.seq_num++
		link.l.Unlock()
		if err != nil {
			log.Println(err)
			return
		}
		num++
	}
	log.Println("send", num, "Offline Msgs to user", uid)
}

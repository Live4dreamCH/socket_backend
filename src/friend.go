package main

import (
	"log"
	"net/http"

	"github.com/Live4dreamCH/socket_backend/db"
	"github.com/gin-gonic/gin"
)

// addfriend 请求格式
type af struct {
	Sid  string `json:"sid" binding:"required"`
	B_id int    `json:"frid" binding:"required"`
}

// A向B发起好友申请
func addfriend(c *gin.Context) {
	var A db.User
	var req af
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"res": "NO", "reason": "json bind error"})
		return
	}

	A.Id, err = sess.get(req.Sid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "wrong sid"})
		return
	}

	ok := A.HasFriend(req.B_id)
	if ok {
		c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "already friend"})
		return
	}

	ok, err = A.AddFriend(req.B_id)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "db err"})
		log.Println("addfriend: db err:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "OK"})
	pushFrReq(req.B_id, A.Id)
}

// resfriend 请求格式
type rf struct {
	Sid  string `json:"sid" binding:"required"`
	A_id int    `json:"frid" binding:"required"`
	Ans  string `json:"ans" binding:"required"`
}

// B回复A的好友申请
func resfriend(c *gin.Context) {
	var B db.User
	var req rf
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"res": "NO", "reason": "json bind error"})
		return
	}

	B.Id, err = sess.get(req.Sid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "wrong sid"})
		return
	}

	if req.Ans == "refuse" {
		c.JSON(http.StatusOK, gin.H{"res": "OK"})
		pushFrAns(req.A_id, B.Id, false, 0)
		return
	}

	ok := B.HasFriend(req.A_id)
	if ok {
		c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "already friend"})
		return
	}

	ok, err = B.AddFriend(req.A_id)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "db err"})
		log.Println("addfriend: db err:", err)
		return
	}

	conv_id, err := db.CreateConv(req.A_id, B.Id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "db err"})
		log.Println("addfriend:CreateConv: db err:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "OK", "conv_id": conv_id})
	pushFrAns(req.A_id, B.Id, true, conv_id)
}

// 向Bid(B)推送来自Aid(A)的好友请求, uid不在线则在数据库中存下来
func pushFrReq(Bid int, Aid int) {
	var B, A db.User
	B.Id = Bid
	A.Id = Aid
	A.GetName()

	msgRouter.l.RLock()
	link, ok := msgRouter.m[Bid]
	msgRouter.l.RUnlock()

	if !ok {
		err := B.StoreFrReq(Aid)
		if err != nil {
			log.Println(err)
		}
		return
	}

	link.l.Lock()
	err := link.conn.WriteJSON(gin.H{"op": "friend request", "seq": link.seq_num, "frid": Aid, "name": A.Name})
	link.seq_num++
	link.l.Unlock()
	if err != nil {
		log.Println(err)
	}
}

// 向Aid(A)推送来自Bid(B)的加好友回复, uid不在线则在数据库中存下来
func pushFrAns(Aid, Bid int, ans bool, conv_id int) {
	var A, B db.User
	A.Id = Aid
	B.Id = Bid
	B.GetName()

	msgRouter.l.RLock()
	link, ok := msgRouter.m[Aid]
	msgRouter.l.RUnlock()

	if !ok {
		err := A.StoreFrAns(Bid, ans, conv_id)
		if err != nil {
			log.Println(err)
		}
		return
	}

	temp := gin.H{"op": "friend answer", "frid": Bid, "name": B.Name}
	if ans {
		temp["ans"] = "accept"
		temp["conv_id"] = conv_id
	} else {
		temp["ans"] = "refuse"
	}

	link.l.Lock()
	temp["seq"] = link.seq_num
	err := link.conn.WriteJSON(temp)
	link.seq_num++
	link.l.Unlock()
	if err != nil {
		log.Println(err)
	}
}

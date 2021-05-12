package main

import (
	"fmt"
	"net/http"

	"github.com/Live4dreamCH/socket_backend/db"
	"github.com/gin-gonic/gin"
)

func addfriend(c *gin.Context) {
	var u db.User
	var req af
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"res": "NO", "reason": "json bind error"})
	} else {
		uid, err := sess.get(req.Sid)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "wrong sid"})
		} else {
			u.Id = uid
			ok := u.HasFriend(req.Frid)
			if ok {
				c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "already friend"})
			} else {
				ok, err = u.AddFriend(req.Frid)
				if !ok {
					c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "db err"})
					fmt.Println("addfriend: db err:", err)
				} else {
					c.JSON(http.StatusOK, gin.H{"res": "OK"})
					pushFrReq(req.Frid, uid)
				}
			}
		}
	}
}

// addfriend 请求格式
type af struct {
	Sid  string `json:"sid" binding:"required"`
	Frid int    `json:"frid" binding:"required"`
}

func resfriend(c *gin.Context) {
	var u db.User
	var req rf
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"res": "NO", "reason": "json bind error"})
	} else {
		uid, err := sess.get(req.Sid)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "wrong sid"})
		} else {
			if req.Ans == "refuse" {
				c.JSON(http.StatusOK, gin.H{"res": "OK"})
				pushFrAns(req.Frid, uid, false, 0)
			} else {
				u.Id = uid
				ok := u.HasFriend(req.Frid)
				if ok {
					c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "already friend"})
				} else {
					ok, err = u.AddFriend(req.Frid)
					if !ok {
						c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "db err"})
						fmt.Println("addfriend: db err:", err)
					} else {
						conv_id, err := db.CreateConv(req.Frid, uid)
						if err != nil {
							c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "db err"})
							fmt.Println("addfriend:CreateConv: db err:", err)
						} else {
							c.JSON(http.StatusOK, gin.H{"res": "OK", "conv_id": conv_id})
							pushFrAns(req.Frid, uid, true, conv_id)
						}
					}
				}
			}
		}
	}
}

// resfriend 请求格式
type rf struct {
	Sid  string `json:"sid" binding:"required"`
	Frid int    `json:"frid" binding:"required"`
	Ans  string `json:"ans" binding:"required"`
}

// 向uid(B)推送来自frid(A)的好友请求, uid不在线则在数据库中存下来
func pushFrReq(uid int, frid int) {
	var B, A db.User
	B.Id = uid
	A.Id = frid
	A.GetName()
	msgRouter.l.RLock()
	link, ok := msgRouter.m[uid]
	msgRouter.l.RUnlock()
	if !ok {
		B.StoreFrReq(frid)
	} else {
		link.l.Lock()
		link.conn.WriteJSON(gin.H{"op": "friend request", "seq": link.seq_num, "frid": frid, "name": A.Name})
		link.seq_num++
		link.l.Unlock()
	}
}

// 向uid(A)推送来自frid(B)的加好友结果, uid不在线则在数据库中存下来
func pushFrAns(uid, frid int, ans bool, conv_id int) {
	var A, B db.User
	A.Id = uid
	B.Id = frid
	B.GetName()
	msgRouter.l.RLock()
	link, ok := msgRouter.m[uid]
	msgRouter.l.RUnlock()
	if !ok {
		A.StoreFrAns(frid, ans, conv_id)
	} else {
		temp := gin.H{"op": "friend answer", "frid": frid, "name": B.Name}
		if ans {
			temp["ans"] = "accept"
			temp["conv_id"] = conv_id
		} else {
			temp["ans"] = "refuse"
		}
		link.l.Lock()
		temp["seq"] = link.seq_num
		link.conn.WriteJSON(temp)
		link.seq_num++
		link.l.Unlock()
	}
}

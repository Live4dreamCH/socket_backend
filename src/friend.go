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
			u.Id = uid
			if req.Ans == "refuse" {
				// todo: ws、数据库删除好友、上线通知
			}
			// ok := u.HasFriend(req.Frid)
			// if ok {
			// 	c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "already friend"})
			// } else {
			// 	ok, err = u.AddFriend(req.Frid)
			// 	if !ok {
			// 		c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "db err"})
			// 		fmt.Println("addfriend: db err:", err)
			// 	} else {
			// 		c.JSON(http.StatusOK, gin.H{"res": "OK"})
			// 	}
			// }
		}
	}
}

// resfriend 请求格式
type rf struct {
	Sid  string `json:"sid" binding:"required"`
	Frid int    `json:"frid" binding:"required"`
	Ans  string `json:"ans" binding:"required"`
}

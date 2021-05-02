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
	if err := c.BindJSON(&req); err != nil {
		// c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "request json error"})
		fmt.Println("addfriend: json bind error:", err)
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

type af struct {
	Sid  string `json:"sid" binding:"required"`
	Frid int    `json:"frid" binding:"required"`
}

package main

import (
	"net/http"

	"github.com/Live4dreamCH/socket_backend/db"
	"github.com/gin-gonic/gin"
)

// 登录
func login(c *gin.Context) {
	var u db.User
	if err := c.BindJSON(&u); err == nil {
		if u.Login() {
			sid := sess.set(u.Id)
			c.JSON(http.StatusOK, gin.H{"res": "OK", "sessionid": sid})
		} else {
			c.JSON(http.StatusOK, gin.H{"res": "NO"})
		}
	}
}

package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Serves unicode entities
	r.GET("/json", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"html": "<b>Hello, world!</b>",
		})
	})

	// Serves literal characters
	r.GET("/purejson", func(c *gin.Context) {
		c.PureJSON(200, gin.H{
			"html": "<b>Hello, world!</b>",
		})
	})

	r.POST("/login", login)

	// listen and serve on 0.0.0.0:8080
	r.Run(":43851")
}

type User struct {
	Id   int    `json:"id" binding:"required"`
	Psw  string `json:"psw" binding:"required"`
	Name string `json:"name" binding:"required"`
}

func (u *User) Print() {
	var l sync.Mutex
	l.Lock()
	fmt.Printf("id=%d, psw=%s, name=%s\n", u.Id, u.Psw, u.Name)
	l.Unlock()
}

func login(c *gin.Context) {
	var u User
	if err := c.BindJSON(&u); err != nil {
		fmt.Println(err)
	} else {
		u.Print()
		c.JSON(http.StatusOK, u)
	}
}

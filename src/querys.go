package main

import (
	"log"
	"net/http"

	"github.com/Live4dreamCH/socket_backend/db"
	"github.com/gin-gonic/gin"
)

func nameService(c *gin.Context) {
	var req struct {
		Sid string `binding:"required"`
		Id  int    `binding:"required"`
	}
	var u db.User

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"res": "NO", "reason": "json bind error"})
		return
	}

	_, err = sess.get(req.Sid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "wrong sid"})
		return
	}

	u.Id = req.Id
	if !u.GetName() {
		c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "wrong uid"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "OK", "name": u.Name})
}

func friendList(c *gin.Context) {
	var req struct {
		Sid string `binding:"required"`
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"res": "NO", "reason": "json bind error"})
		return
	}

	uid, err := sess.get(req.Sid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "wrong sid"})
		return
	}

	type fl struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	var list []fl
	rows, err := db.Get_fr_list.Query(uid)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "db err"})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "db err"})
			return
		}
		list = append(list, fl{id, name})
	}
	response := gin.H{"res": "OK", "friendlist": list}
	c.JSON(http.StatusOK, response)
}

func convList(c *gin.Context) {
	var req struct {
		Sid string `binding:"required"`
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"res": "NO", "reason": "json bind error"})
		return
	}

	uid, err := sess.get(req.Sid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "wrong sid"})
		return
	}

	type cl struct {
		Conv_id int    `json:"conv_id"`
		Name    string `json:"name"`
	}
	var list []cl
	rows, err := db.Get_conv_list.Query(uid, uid)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "db err"})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var conv_id int
		var name string
		err = rows.Scan(&conv_id, &name)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "db err"})
			return
		}
		list = append(list, cl{conv_id, name})
	}
	response := gin.H{"res": "OK", "convlist": list}
	c.JSON(http.StatusOK, response)
}

func convMemList(c *gin.Context) {
	var req struct {
		Sid     string `binding:"required"`
		Conv_id int    `binding:"required"`
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"res": "NO", "reason": "json bind error"})
		return
	}

	uid, err := sess.get(req.Sid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "wrong sid"})
		return
	}

	type cml struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	var list []cml
	rows, err := db.Get_conv_mem_list.Query(req.Conv_id, uid)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "db err"})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusOK, gin.H{"res": "NO", "reason": "db err"})
			return
		}
		list = append(list, cml{id, name})
	}
	response := gin.H{"res": "OK", "convmemlist": list}
	c.JSON(http.StatusOK, response)
}

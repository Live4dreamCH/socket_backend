package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

func wsStarter() {
	http.HandleFunc("/msg", msgListen)
	http.HandleFunc("/file", nil)

	if err := http.ListenAndServe(":43852", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func msgListen(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	suss, uid := wsLogin(conn)
	if !suss {
		conn.Close()
		return
	}

	msgRead(conn, uid)

	// for {
	// 	messageType, p, err := conn.ReadMessage()
	// 	if err != nil {
	// 		log.Println(err)
	// 		return
	// 	}
	// 	if err := conn.WriteMessage(messageType, p); err != nil {
	// 		log.Println(err)
	// 		return
	// 	}
	// }
}

// 用uid找到ws链接
type wsRouter struct {
	m map[int]*wsLink
	l sync.RWMutex
}

// ws链接
type wsLink struct {
	conn *websocket.Conn
	l    sync.Mutex
}

func msgRead(conn *websocket.Conn, uid int) {

}

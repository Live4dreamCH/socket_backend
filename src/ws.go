package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func wsStarter() {
	h := websocket.Handler(Echo)
	http.Handle("/echouid", h)

	if err := http.ListenAndServe(":43852", h); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func Echo(ws *websocket.Conn) {
	var err error
	var recv string
	var send string
	for {
		if err = websocket.Message.Receive(ws, &recv); err != nil {
			fmt.Println("Can't receive")
			break
		}
		uid, err := sess.get(recv)
		if err != nil {
			send = recv + " binds no uid!"
		} else {
			send = recv + " binds uid " + fmt.Sprint(uid)
		}
		if err = websocket.Message.Send(ws, send); err != nil {
			fmt.Println("Can't send")
			break
		}
	}
}

package ws

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func Echo(ws *websocket.Conn) {
	var err error
	var reply string
	for {
		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("Can't receive")
			break
		}

		msg := "Received:  " + reply
		if err = websocket.Message.Send(ws, msg); err != nil {
			fmt.Println("Can't send")
			break
		}
	}
}

func Example() {
	h := websocket.Handler(Echo)
	http.Handle("/wsecho", h)

	if err := http.ListenAndServe(":43852", h); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

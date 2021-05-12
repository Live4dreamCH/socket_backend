package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/Live4dreamCH/socket_backend/db"
	"github.com/gorilla/websocket"
)

// 开启ws服务
func wsStarter() {
	http.HandleFunc("/msg", msgListen)
	// http.HandleFunc("/file", nil)

	if err := http.ListenAndServe(":43852", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

// 处理新ws连接
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

	msgRouter.l.Lock()
	msgRouter.m[uid] = &wsLink{conn: conn, seq_num: 0}
	msgRouter.l.Unlock()

	msgRead(conn, uid)
}

// 用uid找到ws链接
type wsRouter struct {
	m map[int]*wsLink
	l sync.RWMutex
}

// ws链接
type wsLink struct {
	conn    *websocket.Conn
	seq_num int // 发送时的请求号
	l       sync.Mutex
}

// ws结构的共有属性
type wsMain struct {
	Op  string `binding:"required"`
	Seq int
}

// ws好友申请的属性
type wsFrReq struct {
	Frid int
	Name string
}

// ws好友回复的属性
type wsFrAns struct {
	Frid int
	Name string
}

type wsMsg struct {
	Conv_id int
	Sender  int
	Time    string
	Type    string
	Content string
}

type wsSDP struct {
	From int
	To   int
	Sdp  string
}

// 接收ws包
func msgRead(conn *websocket.Conn, uid int) {
	for {
		ty, b, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			msgRouter.l.Lock()
			delete(msgRouter.m, uid)
			msgRouter.l.Unlock()
			log.Println(err)
			break
		}
		if ty != websocket.TextMessage {
			continue
		}
		var head wsMain
		err = json.Unmarshal(b, &head)
		if err != nil {
			log.Println(err)
			continue
		}
		switch head.Op {
		case "msg":
			log.Println(msgForward(uid, b))
		case "connect":
			var sdp_pkg wsSDP
			json.Unmarshal(b, &sdp_pkg)
			msgCopy(sdp_pkg.To, b)
		}
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

func msgCopy(uid int, b []byte) bool {
	msgRouter.l.RLock()
	link, ok := msgRouter.m[uid]
	msgRouter.l.RUnlock()

	if !ok {
		return false
	}
	link.l.Lock()
	err := link.conn.WriteMessage(websocket.TextMessage, b) // todo：这里seq没有用到
	link.seq_num++
	link.l.Unlock()
	return err == nil
}

func msgForward(uid int, b []byte) (err error) {
	// todo: 第一条离线消息的记录

	var pkg struct {
		wsMain
		wsMsg
	}
	err = json.Unmarshal(b, &pkg)
	if err != nil {
		return
	}
	pkg.Sender = uid
	mems, err := db.GetOtherConvMems(uid, pkg.Conv_id)
	if err != nil {
		return
	}
	for _, i := range mems {
		msgRouter.l.RLock()
		link, ok := msgRouter.m[i]
		msgRouter.l.RUnlock()

		if !ok {
			// todo:记录第一条离线消息
			continue
		}
		link.l.Lock()
		pkg.Seq = link.seq_num
		err := link.conn.WriteJSON(pkg)
		link.seq_num++
		link.l.Unlock()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("forward:", pkg, "from", uid, "to", i)
	}
	return nil
}

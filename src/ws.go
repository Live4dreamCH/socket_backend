package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Live4dreamCH/socket_backend/db"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 开启ws服务
func wsStarter() {
	http.HandleFunc("/msg", msgListen)
	http.HandleFunc("/file", fileListen)

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
	log.Println("user", uid, "msg ws login")

	go msgNotice(uid)

	msgRead(conn, uid)
}

// 处理新ws连接
func fileListen(w http.ResponseWriter, r *http.Request) {
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

	fileRouter.l.Lock()
	fileRouter.m[uid] = &wsLink{conn: conn, seq_num: 0}
	fileRouter.l.Unlock()
	log.Println("user", uid, "file ws login")

	fileRead(conn, uid)
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
	Op  string `json:"op" binding:"required"`
	Seq int    `json:"seq"`
}

type wsSDP struct {
	From int    `json:"from"`
	To   int    `json:"to"`
	Sdp  string `json:"sdp"`
}

type wsConnect struct {
	wsSDP
	wsMain
}

// 接收ws包
func msgRead(conn *websocket.Conn, uid int) {
	for {
		ty, b, err := conn.ReadMessage()
		if err != nil {
			msgRouter.l.Lock()
			delete(msgRouter.m, uid)
			msgRouter.l.Unlock()
			conn.Close()
			log.Println("user", uid, "msg ws close, err:", err)
			return
		}
		if ty != websocket.TextMessage {
			log.Println("user", uid, "a binary msg pkg received")
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
			err = msgForward(uid, b)
			if err != nil {
				log.Println(err)
			}
		case "connect":
			msgCopy(uid, b)
		case "connect response":
			msgCopy(uid, b)
		}
	}
}

type wsFile struct {
	Op      string `binding:"required"`
	Conv_id int    `binding:"required"`
}

// 接收ws包
func fileRead(conn *websocket.Conn, uid int) {
	var mems []int
	for {
		ty, b, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			fileRouter.l.Lock()
			delete(fileRouter.m, uid)
			fileRouter.l.Unlock()
			log.Println("user", uid, "file ws close, err:", err)
			return
		}

		// 确定转发目标
		if ty == websocket.TextMessage {
			var head wsFile
			err = json.Unmarshal(b, &head)
			if err != nil {
				log.Println(err)
				continue
			}
			if head.Op == "start" {
				mems, err = db.GetOtherConvMems(uid, head.Conv_id)
				if err != nil {
					log.Println(err)
				} else {
					log.Println("mems changed into:", mems)
				}
			}
			log.Println("user", uid, "file ws ctrl pkg", string(b))
		}

		// 进行转发
		for _, i := range mems {
			fileRouter.l.RLock()
			link, ok := fileRouter.m[i]
			fileRouter.l.RUnlock()
			if !ok {
				continue
			}

			link.l.Lock()
			err := link.conn.WriteMessage(ty, b)
			// link.seq_num++
			link.l.Unlock()
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println("file pkg forward: from", uid, "to", i)
		}
	}
}

// SDP的转发
func msgCopy(uid int, b []byte) bool {
	var pkg wsConnect
	json.Unmarshal(b, &pkg)

	msgRouter.l.RLock()
	link, ok := msgRouter.m[pkg.To]
	msgRouter.l.RUnlock()
	// 对方不在线，回复发送方
	if !ok {
		msgRouter.l.RLock()
		link, ok = msgRouter.m[uid]
		msgRouter.l.RUnlock()
		if ok {
			link.l.Lock()
			err := link.conn.WriteJSON(gin.H{"op": "conncet error", "seq": link.seq_num, "ack": pkg.Seq, "reason": "offline"})
			link.seq_num++
			link.l.Unlock()
			if err != nil {
				log.Println(err)
			}
		}
		log.Println("SDP from", uid, "to", pkg.To, "but target offline")
		return false
	}

	// 对方在线, 进行转发
	link.l.Lock()
	pkg.Seq = link.seq_num
	err := link.conn.WriteJSON(pkg)
	link.seq_num++
	link.l.Unlock()
	log.Println("SDP from", uid, "to", pkg.To, "forward suss:", string(b))
	return err == nil
}

func msgForward(uid int, b []byte) (err error) {
	// 包格式处理
	var pkg, no_content struct {
		wsMain
		db.WsMsg
	}
	err = json.Unmarshal(b, &pkg)
	if err != nil {
		return
	}
	pkg.Sender = uid
	pkg.Time = time.Now().Format("2006-01-02 15:04:05")
	no_content = pkg
	no_content.Content = ""

	// 寻找目标
	mems, err := db.GetOtherConvMems(uid, pkg.Conv_id)
	if err != nil {
		return
	}
	if len(mems) == 0 {
		log.Println("a msg with no forward target:", string(b))
		return
	}

	// 保存消息
	msg_id, suss := pkg.Save()
	if !suss {
		log.Println(no_content, "save fail!")
	}

	// 转发
	var u db.User
	for _, i := range mems {
		msgRouter.l.RLock()
		link, ok := msgRouter.m[i]
		msgRouter.l.RUnlock()
		// 目标不在线, 设置离线消息提醒
		if !ok {
			u.Id = i
			err = u.StoreOfflineMsg(msg_id)
			if err != nil {
				log.Println(err)
			}
			log.Println("save notice for user:", u.Id, ", of saved msg:", msg_id)
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
		log.Println("msg forward:", no_content, "to", i)
	}
	return nil
}

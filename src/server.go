package src

import (
	"fmt"
	"net"
)

func handleConn(c net.Conn) {
	for {
		b := make([]byte, 1024)
		nr, err := c.Read(b)
		if err != nil {
			fmt.Println("read error:", c.LocalAddr(), c.RemoteAddr(), err)
			break
		}
		_, err = c.Write(b[0:nr])
		if err != nil {
			fmt.Println("write error:", c.LocalAddr(), c.RemoteAddr(), err)
			break
		}
	}
	c.Close()
}

func main1() {
	l, err := net.Listen("tcp4", "0.0.0.0:43853")
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		go handleConn(c)
	}
}

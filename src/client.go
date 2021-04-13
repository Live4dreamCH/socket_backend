package src

import (
	"fmt"
	"net"
)

func main0() {
	c, err := net.Dial("tcp4", "hwc.l2d.top:43853")
	if err != nil {
		fmt.Println(err)
	}
	var s string
	for {
		n, err := fmt.Scanln(&s)
		if err != nil {
			fmt.Println("scan:", n, err)
		}
		n, err = c.Write([]byte(s))
		if err != nil {
			fmt.Println("write:", n, err)
		}
		n, err = c.Read([]byte(s))
		if err != nil {
			fmt.Println("read:", n, err)
		}
		n, err = fmt.Println(s)
		if err != nil {
			fmt.Println("print:", n, err)
		}
	}
}

package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("Client Test... start")

	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:5704")
	if err != nil {
		fmt.Println("client start error, exit")
		return
	}

	for i := 0; i < 1; i++ {
		_, err := conn.Write([]byte("Zinx v03"))
		if err != nil {
			fmt.Println("write error ", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error ", err)
			return
		}

		fmt.Printf("server call back: %s, cnt = %d\n", buf, cnt)

		time.Sleep(1 * time.Second)
	}
}

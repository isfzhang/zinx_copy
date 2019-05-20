package znet

import (
	"fmt"
	"net"
	"testing"
	"time"
	"zinx/ziface"
)

//ping test 自定义路由
type PingRouter struct {
	BaseRouter
}

// //Test PreHandle
// func (*PingRouter) PreHandle(request ziface.IRequest) {
// 	fmt.Println("Call Router PreHandle")
// 	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping ....\n"))
// 	if err != nil {
// 		fmt.Println("call back ping ping ping error")
// 	}
// }

//Test Handle
func (*PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

// //Test PostHandle
// func (*PingRouter) PostHandle(request ziface.IRequest) {
// 	fmt.Println("Call Router PostHandle")
// 	_, err := request.GetConnection().GetTCPConnection().Write([]byte("After ping .....\n"))
// 	if err != nil {
// 		fmt.Println("call back ping ping ping error")
// 	}
// }

// ClientTest 模拟客户端
func ClientTest() {
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

// TestServer 模块的测试函数
func TestServer(t *testing.T) {
	s := NewServer("[zinx v0.3]")
	s.AddRouter(&PingRouter{})
	go ClientTest()

	s.Serve()
}

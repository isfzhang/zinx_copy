package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

// PingRouter 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Handle 方法
func (*PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")

	fmt.Println("recv from client: msgID= ", request.MsgID(), ", data=", string(request.Data()))
	err := request.Connection().SendMsg(0, []byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

// HelloRouter say hi
type HelloRouter struct {
	znet.BaseRouter
}

// Handle 方法
func (*HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloRouter Handle")

	fmt.Println("recv from client: msgID= ", request.MsgID(), ", data=", string(request.Data()))
	err := request.Connection().SendMsg(1, []byte("Hello Zinx Router V06\n"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	s := znet.NewServer()
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	s.Serve()
}

package znet

import (
	"errors"
	"fmt"
	"net"
	"time"
	"zinx/ziface"
)

// Server IServer接口实现， 定义一个Server服务类
type Server struct {
	Name       string
	IPVsersion string
	IP         string
	Port       int
}

// CallBackToClient 定义当前客户端的handle api
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Println("[Conn Handle] CallBackToClient ...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err ", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}

// Start 开启网络服务
func (s *Server) Start() {
	fmt.Printf("[START] Server listenner at IP: %s, Port %d, is starting\n", s.IP, s.Port)

	go func() {
		addr, err := net.ResolveTCPAddr(s.IPVsersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}

		listenner, err := net.ListenTCP(s.IPVsersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVsersion, "err", err)
			return
		}

		fmt.Println("start Zinx server ", s.Name, " succ, now listening...")

		// 生成ID
		var cid uint32
		cid = 0

		for {
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}

			dealConn := NewConnection(conn, cid, CallBackToClient)
			cid++

			go dealConn.Start()
		}
	}()
}

// Stop 停止服务并清理
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server, name ", s.Name)
	// Todo 清理资源
}

// Serve 开启服务
func (s *Server) Serve() {
	s.Start()

	// Todo 启动时处理的其他事情

	// 防止主进程退出
	for {
		time.Sleep(10 * time.Second)
	}
}

// NewServer 创建一个服务器句柄
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       name,
		IPVsersion: "tcp4",
		IP:         "0.0.0.0",
		Port:       5704,
	}

	return s
}

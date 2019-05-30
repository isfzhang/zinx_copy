package znet

import (
	"fmt"
	"net"
	"time"
	"zinx/utils"
	"zinx/ziface"
)

// Server IServer接口实现， 定义一个Server服务类
type Server struct {
	Name       string
	IPVsersion string
	IP         string
	Port       int
	msgHandler ziface.IMsgHandle
}

// NewServer 创建一个服务器句柄
func NewServer() ziface.IServer {
	utils.GlobalObject.Reload()

	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVsersion: "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TCPPort,
		msgHandler: MewMsgHandle(),
	}
	return s
}

// Start 开启网络服务
func (s *Server) Start() {
	fmt.Printf("[START] Server listenner at IP: %s, Port %d, is starting\n", s.IP, s.Port)
	fmt.Println(utils.GlobalObject.Maxconn)

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

			dealConn := NewConnection(conn, cid, s.msgHandler)
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

// AddRouter 路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgID, router)
}

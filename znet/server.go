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
	Name        string
	IPVsersion  string
	IP          string
	Port        int
	msgHandler  ziface.IMsgHandle
	ConnMgr     ziface.IConnManager // 连接管理器
	OnConnStart func(conn ziface.IConnection)
	OnConnStop  func(conn ziface.IConnection)
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
		ConnMgr:    NewConnManager(),
	}
	return s
}

// Start 开启网络服务
func (s *Server) Start() {
	fmt.Printf("[START] Server listenner at IP: %s, Port %d, is starting\n", s.IP, s.Port)
	fmt.Println(utils.GlobalObject.Maxconn)

	go func() {
		// 启动工作池
		s.msgHandler.StartWorkerPool()

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

			// 超过服务器最大连接控制，关闭连接
			if s.ConnMgr.Len() >= utils.GlobalObject.Maxconn {
				conn.Close()
				fmt.Println("连接数过多，关闭新连接")
				continue
			}
			dealConn := NewConnection(s, conn, cid, s.msgHandler)
			cid++

			go dealConn.Start()
		}
	}()
}

// Stop 停止服务并清理
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server, name ", s.Name)
	// Todo 清理资源
	s.ConnMgr.ClearConn()
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

// GetConnMgr 得到连接管理
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// SetOnConnStart 设置连接创建时执行函数
func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// SetOnConnStop 设置连接断开时执行函数
func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// CallOnConnStart 调用
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart...")
		s.OnConnStart(conn)
	}
}

// CallOnConnStop 调用
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStop...")
		s.OnConnStop(conn)
	}
}

package ziface

import (
	"net"
)

// IConnection 定义连接接口
type IConnection interface {
	Start()
	Stop()
	GetTCPConnection() *net.TCPConn
	GetConnID() uint32
	RemoteAddr() net.Addr
}

// HandFunc 统一处理链接业务的接口
type HandFunc func(*net.TCPConn, []byte, int) error

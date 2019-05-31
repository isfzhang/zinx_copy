package ziface

import (
	"net"
)

// IConnection 定义连接接口
type IConnection interface {
	Start()
	Stop()
	TCPConnection() *net.TCPConn
	ConnID() uint32
	RemoteAddr() net.Addr
	SendMsg(msgID uint32, data []byte) error
	SendBuffMsg(msgID uint32, data []byte) error // 带缓冲发送消息接口
}

package znet

import (
	"zinx/ziface"
)

// Request 客户端请求
type Request struct {
	conn ziface.IConnection
	msg  ziface.IMessage
}

// Connection 获取请求连接信息
func (r *Request) Connection() ziface.IConnection {
	return r.conn
}

// Data 获取请求消息的数据
func (r *Request) Data() []byte {
	return r.msg.Data()
}

// MsgID 获取请求消息的ID
func (r *Request) MsgID() uint32 {
	return r.msg.MsgID()
}

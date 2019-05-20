package znet

import (
	"zinx/ziface"
)

// Request 客户端请求
type Request struct {
	conn ziface.IConnection
	data []byte
}

// GetConnection 获取请求连接信息
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// GetData 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.data
}

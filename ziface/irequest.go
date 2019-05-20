package ziface

// IRequest 包装客户端请求的连接信息和请求数据的接口
type IRequest interface {
	GetConnection() IConnection // 获取请求连接信息
	GetData() []byte            // 获取请求消息的数据
}

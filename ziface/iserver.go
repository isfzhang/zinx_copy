package ziface

// IServer 定义服务器接口
type IServer interface {
	Start()
	Stop()
	Serve()                                 // 开启业务服务
	AddRouter(msgID uint32, router IRouter) // 路由功能，当前服务注册一个路由业务方法
	GetConnMgr() IConnManager               // 得到连接管理
	SetOnConnStart(func(IConnection))       // 设置连接创建时hook函数
	SetOnConnStop(func(IConnection))        // 设置连接断开时hook函数
	CallOnConnStart(conn IConnection)
	CallOnConnStop(conn IConnection)
}

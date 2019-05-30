package ziface

// IServer 定义服务器接口
type IServer interface {
	Start()
	Stop()
	Serve()                                 // 开启业务服务
	AddRouter(msgID uint32, router IRouter) // 路由功能，当前服务注册一个路由业务方法
}

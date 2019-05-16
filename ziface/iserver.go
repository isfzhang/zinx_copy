package ziface

// IServer 定义服务器接口
type IServer interface {
	Start()
	Stop()
	// 开启业务服务
	Serve()
}

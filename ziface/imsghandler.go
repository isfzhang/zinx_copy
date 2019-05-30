package ziface

// IMsgHandle 消息管理抽象层
type IMsgHandle interface {
	DoMsgHandler(request IRequest)          // 马上以非阻塞方式处理消息
	AddRouter(msgID uint32, router IRouter) // 为消息添加处理逻辑
	StartWorkerPool()                       //  启动工作池
	SendMsgToTaskQueue(request IRequest)    // 加入消息队列
}

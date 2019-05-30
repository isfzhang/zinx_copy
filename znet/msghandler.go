package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

// MsgHandle 消息管理结构体
type MsgHandle struct {
	Apis         map[uint32]ziface.IRouter
	WorkPoolSize uint32
	TaskQueue    []chan ziface.IRequest
}

// MewMsgHandle 创建结构体
func MewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:         make(map[uint32]ziface.IRouter),
		WorkPoolSize: utils.GlobalObject.WorkPoolSize,
		TaskQueue:    make([]chan ziface.IRequest, utils.GlobalObject.WorkPoolSize),
	}
}

// DoMsgHandler 马上处理消息
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.Apis[request.MsgID()]
	if !ok {
		fmt.Println("api magID = ", request.MsgID(), " is not FOUND")
		return
	}

	// 执行对应处理方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// AddRouter 添加消息的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeated api, msgID = " + strconv.Itoa(int(msgID)))
	}

	// 绑定msg和api
	mh.Apis[msgID] = router
	fmt.Println("Add api msgID = ", msgID)
}

// StartOneWorker 启动一个工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("worker ID = ", workerID, " is started.")

	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// StartWorkerPool 启动工作池
func (mh *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(mh.WorkPoolSize); i++ {
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// SendMsgToTaskQueue 分配消息到不同队列
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	workerID := request.Connection().ConnID() % mh.WorkPoolSize
	fmt.Println("Add ConnID = ", request.Connection().ConnID(), "request msgID = ", request.MsgID(), " to WorkerID = ", workerID)
	mh.TaskQueue[workerID] <- request

}

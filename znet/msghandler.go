package znet

import (
	"fmt"
	"strconv"
	"zinx/ziface"
)

// MsgHandle 消息管理结构体
type MsgHandle struct {
	Apis map[uint32]ziface.IRouter
}

// MewMsgHandle 创建结构体
func MewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
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


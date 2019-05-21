package ziface

// IMessage 封装请求消息
type IMessage interface {
	DataLen() uint32
	MsgID() uint32
	Data() []byte

	SetDataLen(uint32)
	SetMsgID(uint32)
	SetData([]byte)
}

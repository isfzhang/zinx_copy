package znet

// Message 请求消息包
type Message struct {
	id      uint32
	dataLen uint32
	data    []byte
}

// NewMsgPackage 新建消息包
func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		id:      id,
		dataLen: uint32(len(data)),
		data:    data,
	}
}

// MsgID 获取消息ID
func (msg *Message) MsgID() uint32 {
	return msg.id
}

// DataLen 获取消息长度
func (msg *Message) DataLen() uint32 {
	return msg.dataLen
}

// Data 获取消息内容
func (msg *Message) Data() []byte {
	return msg.data
}

// SetMsgID 设置消息ID
func (msg *Message) SetMsgID(id uint32) {
	msg.id = id
}

// SetDataLen 设置消息长度
func (msg *Message) SetDataLen(len uint32) {
	msg.dataLen = len
}

// SetData 设置消息内容
func (msg *Message) SetData(data []byte) {
	msg.data = data
}

package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

// DataPack 封包拆包实例
type DataPack struct{}

// NewDataPack 初始化DataPack实例
func NewDataPack() *DataPack {
	return &DataPack{}
}

// HeadLen 获取包头长度
func (dp *DataPack) HeadLen() uint32 {
	// id-uint32-4 + dataLen-uint32-4
	return 8
}

// Pack 封包方法
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// write dataLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.DataLen()); err != nil {
		return nil, err
	}

	// write MsgID
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.MsgID()); err != nil {
		return nil, err
	}

	// write data
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.Data()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// UnPack 拆包
func (dp *DataPack) UnPack(binaryData []byte) (ziface.IMessage, error) {
	dataBuff := bytes.NewReader(binaryData)

	msg := &Message{}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.dataLen); err != nil {
		return nil, err
	}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.id); err != nil {
		return nil, err
	}

	if utils.GlobalObject.MaxPacketSize > 0 && msg.dataLen > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("Too large msg data received")
	}

	// 这里不包含数据，只有数据长度，需要再次读取
	return msg, nil
}

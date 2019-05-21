package ziface

// IDataPack 数据封包 拆包
type IDataPack interface {
	HeadLen() uint32
	Pack(msg IMessage) ([]byte, error)
	UnPack([]byte) (IMessage, error)
}

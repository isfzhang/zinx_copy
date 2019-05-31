package ziface

// IConnManager 连接管理抽象层
type IConnManager interface {
	Add(conn IConnection)
	Remove(conn IConnection)
	Get(connID uint32) (IConnection, error)
	Len() int
	ClearConn() // 删除并停止所有连接
}

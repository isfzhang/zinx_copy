package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

// ConnManager 连接管理模块
type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex
}

// NewConnManager 创建一个链接管理
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// Add 添加链接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	connMgr.connections[conn.ConnID()] = conn

	fmt.Println("connection add to ConnManager successfully: conn num = ", connMgr.Len())
}

// Remove 删除连接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	delete(connMgr.connections, conn.ConnID())

	fmt.Println("connection Remove ConnID = ", conn.ConnID(), "successfullt: conn num = ", connMgr.Len())
}

// Get 利用ConnID获取链接
func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	}
	return nil, errors.New("connection not found")
}

// Len 获取当前连接个数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

// ClearConn 清除并停止所有连接
func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	for connID, conn := range connMgr.connections {
		conn.Stop()

		delete(connMgr.connections, connID)
	}

	fmt.Println("Clear ALL Connections successfully: conn num = ", connMgr.Len())
}

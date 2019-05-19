package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

// Connection 连接
type Connection struct {
	Conn     *net.TCPConn
	ConnID   uint32
	isClosed bool

	handleAPI    ziface.HandFunc
	ExitBuffChan chan bool
}

// NewConnection 创建连接
func NewConnection(conn *net.TCPConn, connID uint32, callbackAPI ziface.HandFunc) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		handleAPI:    callbackAPI,
		ExitBuffChan: make(chan bool, 1),
	}

	return c
}

// StartReader 处理conn读数据
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()

	for {
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err", err)
			c.ExitBuffChan <- true
			continue
		}

		// 调用当前链接
		if err = c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("connID ", c.ConnID, "handle is error")
			c.ExitBuffChan <- true
			return
		}
	}
}

// Start 启动当前连接
func (c *Connection) Start() {

	go c.StartReader()

	for {
		select {
		case <-c.ExitBuffChan:
			return
		}
	}
}

// Stop 停止当前连接
func (c *Connection) Stop() {
	if c.isClosed == true {
		return
	}

	c.isClosed = true

	c.Conn.Close()
	c.ExitBuffChan <- true
	close(c.ExitBuffChan)
}

// GetTCPConnection 获取原始socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr 获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

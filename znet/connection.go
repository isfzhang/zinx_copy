package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/ziface"
)

// Connection 连接
type Connection struct {
	conn     *net.TCPConn
	connID   uint32
	isClosed bool

	Router       ziface.IRouter
	ExitBuffChan chan bool
}

// NewConnection 创建连接
func NewConnection(conn *net.TCPConn, connID uint32, callbackAPI ziface.IRouter) *Connection {
	c := &Connection{
		conn:         conn,
		connID:       connID,
		isClosed:     false,
		Router:       callbackAPI,
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
		// 创建拆包对象
		dp := NewDataPack()

		headData := make([]byte, dp.HeadLen())
		if _, err := io.ReadFull(c.TCPConnection(), headData); err != nil {
			fmt.Println("read msg head error ", err)
			break
		}

		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			break
		}

		var data []byte
		if msg.DataLen() > 0 {
			data = make([]byte, msg.DataLen())
			if _, err := io.ReadFull(c.TCPConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				break
			}
		}
		msg.SetData(data)

		req := Request{
			conn: c,
			msg:  msg,
		}

		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
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

	c.conn.Close()
	c.ExitBuffChan <- true
	close(c.ExitBuffChan)
}

// TCPConnection 获取原始socket TCPConn
func (c *Connection) TCPConnection() *net.TCPConn {
	return c.conn
}

// ConnID 获取当前连接ID
func (c *Connection) ConnID() uint32 {
	return c.connID
}

// RemoteAddr 获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// SendMsg 发送封包后的数据
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when send msg")
	}

	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgID)
		return err
	}

	if _, err := c.conn.Write(msg); err != nil {
		fmt.Println("write msg id ", msgID, " error")
		c.ExitBuffChan <- true
		return err
	}

	return nil
}

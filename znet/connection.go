package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

// Connection 连接
type Connection struct {
	TCPServer ziface.IServer // server归属
	conn      *net.TCPConn
	connID    uint32
	isClosed  bool

	MsgHandler   ziface.IMsgHandle
	ExitBuffChan chan bool
	msgChan      chan []byte // 读写两个协程之间的通信
	msgBuffChan  chan []byte // 有缓冲通道
}

// NewConnection 创建连接
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TCPServer:    server, // 将隶属的server传递
		conn:         conn,
		connID:       connID,
		isClosed:     false,
		MsgHandler:   msgHandler,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
		msgBuffChan:  make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
	}

	c.TCPServer.GetConnMgr().Add(c) //添加当前连接到连接管理器
	return c
}

// StartWriter 处理数据发送到客户端
func (c *Connection) StartWriter() {
	defer fmt.Println(c.RemoteAddr().String(), " conn Writer exit!")

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.conn.Write(data); err != nil {
				fmt.Println("Send Data error: ", err, " Conn Writer exit")
				return
			}
		case data, ok := <-c.msgBuffChan:
			if ok {
				if _, err := c.conn.Write(data); err != nil {
					fmt.Println("Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is Closed")
				break
			}
		case <-c.ExitBuffChan:
			return
		}
	}
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

		if utils.GlobalObject.WorkPoolSize > 0 {
			// 工作池已启用
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

// Start 启动当前连接
func (c *Connection) Start() {
	go c.StartReader()
	go c.StartWriter()

	// 执行hook方法
	c.TCPServer.CallOnConnStart(c)
}

// Stop 停止当前连接
func (c *Connection) Stop() {
	if c.isClosed == true {
		return
	}

	c.isClosed = true

	// 执行hook方法
	c.TCPServer.CallOnConnStop(c)

	c.conn.Close()
	c.ExitBuffChan <- true

	// 从连接管理器删除
	c.TCPServer.GetConnMgr().Remove(c)
	// 关闭该连接全部管道
	close(c.ExitBuffChan)
	close(c.msgBuffChan)
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

	c.msgChan <- msg

	return nil
}

// SendBuffMsg 有缓存方式发送封包后的数据
func (c *Connection) SendBuffMsg(msgID uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when send msg")
	}

	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgID)
		return err
	}

	c.msgBuffChan <- msg

	return nil
}

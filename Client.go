package main

import (
	"fmt"
	"net"
	"time"
  "io"
  "zinx/znet"
)

func main() {
	fmt.Println("Client Test... start")

	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:5704")
	if err != nil {
		fmt.Println("client start error, exit")
		return
	}

	for {
    dp := znet.NewDataPack()
    msg, _ := dp.Pack(znet.NewMsgPackage(0, []byte("zinx V05 client test message")))
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("write error ", err)
			return
		}

		headData := make([]byte, dp.HeadLen())
		if _, err := io.ReadFull(conn, headData); err != nil {
			fmt.Println("read head error ", err)
			return 
		}

		msgHead, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			return
		}

		if msgHead.DataLen() > 0 {
      msg := msgHead.(*znet.Message)
      
			data := make([]byte, msg.DataLen())
			if _, err := io.ReadFull(conn, data); err != nil {
				fmt.Println("read msg data error ", err)
				return 
			}
      
      fmt.Println("==> Recv Msg: ID =", msg.MsgID(), ", len=", msg.DataLen(), ", data=", string(data))
      
		}

		time.Sleep(1 * time.Second)
	}
}

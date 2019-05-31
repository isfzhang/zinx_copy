package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"zinx/ziface"
)

// GlobalObj 框架的全局参数
type GlobalObj struct {
	// Server
	TCPServr ziface.IServer
	Host     string
	TCPPort  int
	Name     string

	// Zinx
	Version          string
	MaxPacketSize    uint32
	Maxconn          int
	WorkPoolSize     uint32
	MaxWorkerTaskLen uint32 // 任务队列长度
	MaxMsgChanLen    uint32
	// config file path
	ConfFilePath string
}

// GlobalObject 定义全局对象
var GlobalObject *GlobalObj

func init() {
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "V08",
		TCPPort:          5704,
		Host:             "0.0.0.0",
		Maxconn:          12000,
		MaxPacketSize:    4096,
		ConfFilePath:     "conf/zinx.json",
		WorkPoolSize:     10,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    1024,
	}

	GlobalObject.Reload()
}

// Reload 加载配置文件
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		// file not exists return
		fmt.Println("读取配置文件失败", err)
		return
		// panic(err)
	}

	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		fmt.Println("解析配置文件失败", err)
		return
		// panic(err)
	}
}

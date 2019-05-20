package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

// GlobalObj 框架的全局参数
type GlobalObj struct {
	TCPServr      ziface.IServer
	Host          string
	TCPPort       int
	Name          string
	Version       string
	MaxPacketSize uint32
	Maxconn       int
}

// GlobalObject 定义全局对象
var GlobalObject *GlobalObj

func init() {
	GlobalObject = &GlobalObj{
		Name:          "ZinxServerApp",
		Version:       "V04",
		TCPPort:       5704,
		Host:          "0.0.0.0",
		Maxconn:       12000,
		MaxPacketSize: 4096,
	}

	GlobalObject.Reload()
}

// Reload 加载配置文件
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

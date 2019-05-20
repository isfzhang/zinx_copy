package znet

import (
	"zinx/ziface"
)

// BaseRouter 路由基类
type BaseRouter struct{}

// PreHandle 处理前方法
func (br *BaseRouter) PreHandle(req ziface.IRequest) {}

// Handle 处理时方法
func (br *BaseRouter) Handle(req ziface.IRequest) {}

// PostHandle 处理后方法
func (br *BaseRouter) PostHandle(req ziface.IRequest) {}

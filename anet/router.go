package anet

import "Alchemist/iface"

type BaseRouter struct{}

//处理业务之前的方法
func (br *BaseRouter) PreHandle(request iface.IRequest) {}

//处理业务的主方法
func (br *BaseRouter) Handle(request iface.IRequest) {}

//处理业务之后的方法
func (br *BaseRouter) PostHandle(request iface.IRequest) {}

package main

import (
	"Alchemist/anet"
	"Alchemist/iface"
	"fmt"
)

type PingRouter struct {
	anet.BaseRouter
}

//PreHandle
func (this *PingRouter) PreHandle(request iface.IRequest) {
	fmt.Println("Call Router PreHandle")
	request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
}

//Handle
func (this *PingRouter) Handle(request iface.IRequest) {
	fmt.Println("Call Router Handle")
	request.GetConnection().GetTCPConnection().Write([]byte("ping ping ping...\n"))
}

//PostHandle
func (this *PingRouter) PostHandle(request iface.IRequest) {
	fmt.Println("Call Router PostHandle")
	request.GetConnection().GetTCPConnection().Write([]byte("post ping...\n"))
}

//服务端应用程序
func main() {
	//1.创建一个server
	s := anet.NewServer("[v0.3]")

	//2.添加自定义Router
	s.AddRouter(&PingRouter{})

	//3.启动server
	s.Serve()
}

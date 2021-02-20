package main

import (
	"Alchemist/anet"
	"Alchemist/iface"
	"fmt"
)

type PingRouter struct {
	anet.BaseRouter
}

//Handle
func (this *PingRouter) Handle(request iface.IRequest) {
	fmt.Println("Call Router Ping Handle")

	//先读取客户端的数据, 再回写
	fmt.Println("recv from client, msgId:", request.GetMsgID(), "data:", string(request.GetData()))

	err := request.GetConnection().SendMsg(200, []byte("ping..ping..ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloRouter struct {
	anet.BaseRouter
}

//Handle
func (this *HelloRouter) Handle(request iface.IRequest) {
	fmt.Println("Call Router Hello Handle")

	//先读取客户端的数据, 再回写
	fmt.Println("recv from client, msgId:", request.GetMsgID(), "data:", string(request.GetData()))

	err := request.GetConnection().SendMsg(201, []byte("hello.."))
	if err != nil {
		fmt.Println(err)
	}
}

//服务端应用程序
func main() {
	//1.创建一个server
	s := anet.NewServer("[v0.7]")

	//2.添加自定义Router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	//3.启动server
	s.Serve()
}

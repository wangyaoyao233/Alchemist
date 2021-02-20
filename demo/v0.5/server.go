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
	fmt.Println("Call Router Handle")

	//先读取客户端的数据, 再回写
	fmt.Println("recv from client, msgId:", request.GetMsgID(), "data:", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("ping..ping..ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

//服务端应用程序
func main() {
	//1.创建一个server
	s := anet.NewServer("[v0.5]")

	//2.添加自定义Router
	s.AddRouter(&PingRouter{})

	//3.启动server
	s.Serve()
}

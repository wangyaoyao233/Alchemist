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

//创建连接之后的hook函数
func DoConnectionBegin(conn iface.IConnection) {
	fmt.Println("===>DoConnectionBegin is Called..")
	conn.SendMsg(202, []byte("DoConnectionBegin"))

	//给当前连接设置自定义属性
	conn.SetProperty("Name", "AAA")
	conn.SetProperty("Email", "aaa@a.com")
}

//销毁连接之前的hook函数
func DoConnectionLost(conn iface.IConnection) {
	fmt.Println("===>DoConnectionLost is Called..")

	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name: ", name)
	}
	if email, err := conn.GetProperty("Email"); err == nil {
		fmt.Println("Email: ", email)
	}
}

//服务端应用程序
func main() {
	//1.创建一个server
	s := anet.NewServer("[v0.10]")

	//2.注册连接的hook函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//3.添加自定义Router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	//3.启动server
	s.Serve()
}

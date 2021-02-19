package main

import "Alchemist/net"

//服务端应用程序
func main() {
	//1.创建一个server句柄
	s := net.NewServer("[v0.1]")
	//2.启动server
	s.Serve()
}

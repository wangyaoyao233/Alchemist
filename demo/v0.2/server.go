package main

import net "Alchemist/anet"

//服务端应用程序
func main() {
	//1.创建一个server句柄
	s := net.NewServer("[v0.2]")
	//2.启动server
	s.Serve()
}

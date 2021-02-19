package net

import (
	"Alchemist/iface"
	"fmt"
	"net"
)

//Iserver的接口实现, 定义一个Server的服务器模块
type Server struct {
	//服务器的名称
	Name string
	//服务器绑定的ip版本
	IPversion string
	//服务器监听的ip
	IP string
	//服务器监听的port
	Port int
}

//提供一个初始化Server模块方法
func NewServer(name string) iface.IServer {
	s := &Server{
		Name:      name,
		IPversion: "tcp4",
		IP:        "0.0.0.0",
		Port:      9000,
	}
	return s
}

//启动服务器
func (s *Server) Start() {
	fmt.Printf("[Start] Server Listen at IP:%s, Port:%d, is starting...\n", s.IP, s.Port)

	go func() {
		//1.获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPversion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("Resolve tcp addr error", err)
			return
		}

		//2.监听服务器的地址
		listener, err := net.ListenTCP(s.IPversion, addr)
		if err != nil {
			fmt.Println("listen error", err)
			return
		}
		fmt.Println("Start server succ ", s.Name)

		//3.阻塞的等待客户端连接, 处理客户端连接业务(读写)
		for {
			//如果有客户端连接, 阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("AcceptTcp err", err)
				continue
			}

			//已经与客户端建立连接, 做一些业务
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("recv buf err", err)
						continue
					}

					fmt.Printf("recv client buf %s, cnt %d\n", buf, cnt)
					//回显
					if _, err = conn.Write(buf[:cnt]); err != nil {
						fmt.Println("write back buf err", err)
						continue
					}
				}
			}()
		}
	}()

}

func (s *Server) Stop() {
	//TODO 将一些服务器的资源, 状态,或者一些已经开辟的连接信息停止或回收
}

func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()

	//TODO 扩展:可以做一些启动服务器之后的额外功能

	//阻塞状态
	select {}
}

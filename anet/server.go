package anet

import (
	"Alchemist/iface"
	"Alchemist/utils"
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

	//注册的router
	Router iface.IRouter
}

//提供一个初始化Server模块方法
func NewServer(name string) iface.IServer {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPversion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Router:    nil,
	}
	return s
}

//启动服务器
func (s *Server) Start() {
	fmt.Printf("[conf] Server Name:%s, IP:%s, Port:%d\n", utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[conf] Version:%s, MaxConn:%d, MaxPackageSize:%d\n", utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

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

		var cid uint32
		cid = 0
		//3.阻塞的等待客户端连接, 处理客户端连接业务(读写)
		for {
			//如果有客户端连接, 阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("AcceptTcp err", err)
				continue
			}

			//绑定连接的客户端，得到连接模块
			//func(*net.TCPConn, []byte, int)
			dealConn := NewConnection(conn, cid, s.Router)
			cid++

			//启动当前的连接
			go dealConn.Start()
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

//添加路由: 给当前的server注册一个路由方法，供客户端的连接使用
func (s *Server) AddRouter(router iface.IRouter) {
	s.Router = router
	fmt.Println("Add Router succ")
}

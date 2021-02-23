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

	//MsgHandle消息管理模块
	MsgHandler iface.IMsgHandle

	//ConnMgr连接管理模块
	ConnMgr iface.IConnManager

	//hook函数
	OnConnStart func(conn iface.IConnection)
	OnConnStop  func(conn iface.IConnection)
}

//提供一个初始化Server模块方法
func NewServer(name string) iface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPversion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

//启动服务器
func (s *Server) Start() {
	fmt.Printf("[conf] Server Name:%s, IP:%s, Port:%d\n", utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[conf] Version:%s, MaxConn:%d, MaxPackageSize:%d\n", utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

	go func() {
		//0.开启消息队列及Worker工作池
		s.MsgHandler.StartWorkerPool()

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
			fmt.Println("Get Conn remote addr =", conn.RemoteAddr().String())

			//设置最大的连接个数，如果超过最大连接，那么关闭新的连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				//TODO 给客户端响应一个最大连接的错误
				fmt.Println("too many connections MaxConn:", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}
			//绑定连接的客户端，得到连接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			//启动当前的连接
			go dealConn.Start()
		}
	}()

}

func (s *Server) Stop() {
	//将一些服务器的资源, 状态,或者一些已经开辟的连接信息停止或回收
	fmt.Println("[Server Stop] server name:", s.Name)
	//清理连接管理模块
	s.ConnMgr.Clear()
}

func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()

	//TODO 扩展:可以做一些启动服务器之后的额外功能

	//阻塞状态
	select {}
}

//添加路由: 给当前的server注册一个路由方法，供客户端的连接使用
func (s *Server) AddRouter(msgID uint32, router iface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router succ")
}

//获取ConnMgr连接管理器方法
func (s *Server) GetConnMgr() iface.IConnManager {
	return s.ConnMgr
}

//注册连接之间的hook函数
func (s *Server) SetOnConnStart(hookFunc func(conn iface.IConnection)) {
	s.OnConnStart = hookFunc
}

//注册连接销毁之前的hook函数
func (s *Server) SetOnConnStop(hookFunc func(conn iface.IConnection)) {
	s.OnConnStop = hookFunc
}

//调用hook函数
func (s *Server) CallOnConnStart(conn iface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("Call OnConnStart()")
		s.OnConnStart(conn)
	}

}
func (s *Server) CallOnConnStop(conn iface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("Call OnConnStop()")
		s.OnConnStop(conn)
	}
}

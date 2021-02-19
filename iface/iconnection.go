package iface

import "net"

//连接模块的抽象层
type Iconnection interface {
	//启动连接,让当前的连接准备开始工作
	Start()
	//停止连接, 结束当前连接的工作
	Stop()
	//获取当前连接的conn对象(套接字)
	GetTCPConn() *net.TCPConn
	//获取连接ID
	GetConnID() uint32
	//获取客户端连接的地址和端口
	GetRemoteAddr() net.Addr
	//发送数据的方法
	Send(data []byte) error
}

//定义一个处理连接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error

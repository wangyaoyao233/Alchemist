package anet

import (
	"Alchemist/iface"
	"fmt"
	"net"
)

type Connection struct {
	//socket TCP套接字
	Conn *net.TCPConn
	//连接的ID
	ConnID uint32
	//当前的连接状态(是否已经关闭)
	isClosed bool
	//与当前连接说绑定的处理业务的方法
	handleAPI iface.HandleFunc
	//等待连接被动退出的channel
	ExitChan chan bool
}

//初始化方法
func NewConnection(conn *net.TCPConn, connID uint32, callback_api iface.HandleFunc) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		handleAPI: callback_api,
		isClosed:  false,
		ExitChan:  make(chan bool),
	}
	return c
}

func (conn *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")

	defer fmt.Println("connID= ", conn.ConnID, " Reader is exit, remote addr is ", conn.GetRemoteAddr().String())
	defer conn.Stop()

	for {
		//读取客户端的数据到buf中
		buf := make([]byte, 512)
		cnt, err := conn.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err..", err)
			continue
		}

		//调用当前连接说绑定的handleAPI
		if err := conn.handleAPI(conn.Conn, buf, cnt); err != nil {
			fmt.Println("ConnID:", conn.ConnID, "handle is error ", err)
			break
		}
	}
}

//启动连接,让当前的连接准备开始工作
func (conn *Connection) Start() {
	fmt.Println("Conn Start()..ConnID=", conn.ConnID)

	//启动从当前连接的读数据业务
	go conn.StartReader()

	// TODO启动从当前连接的写数据业务
}

//停止连接, 结束当前连接的工作
func (conn *Connection) Stop() {
	fmt.Println("Conn Stop()..ConnID=", conn.ConnID)

	//如果当前连接已经关闭
	if conn.isClosed {
		return
	}
	conn.isClosed = true

	//关闭socket连接
	conn.Conn.Close()
	//回收资源
	close(conn.ExitChan)
}

//获取当前连接的conn对象(套接字)
func (conn *Connection) GetTCPConn() *net.TCPConn {
	return conn.Conn
}

//获取连接ID
func (conn *Connection) GetConnID() uint32 {
	return conn.ConnID
}

//获取客户端连接的地址和端口
func (conn *Connection) GetRemoteAddr() net.Addr {
	return conn.Conn.RemoteAddr()
}

//发送数据的方法
//func (conn *Connection) Send(data []byte) error

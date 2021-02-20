package anet

import (
	"Alchemist/iface"
	"errors"
	"fmt"
	"io"
	"net"
)

type Connection struct {
	//socket TCP套接字
	Conn *net.TCPConn
	//连接的ID
	ConnID uint32
	//当前的连接状态(是否已经关闭)
	isClosed bool
	//等待连接被动退出的channel
	ExitChan chan bool

	//该连接的Router
	Router iface.IRouter
}

//初始化方法
func NewConnection(conn *net.TCPConn, connID uint32, router iface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		ExitChan: make(chan bool),
		Router:   router,
	}
	return c
}

func (conn *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")

	defer fmt.Println("connID= ", conn.ConnID, " Reader is exit, remote addr is ", conn.GetRemoteAddr().String())
	defer conn.Stop()

	for {
		//创建一个拆包解包的对象
		dp := NewDataPack()

		//读取msgHead (8字节)
		headData := make([]byte, int(dp.GetHeadLen()))
		if _, err := io.ReadFull(conn.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error", err)
			break
		}

		//拆包,得到dataLen, msgId, 放在msg中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}

		var data []byte
		//根据dataLen,再次读取data, 放在msg中
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(conn.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error: ", err)
				break
			}
		}
		msg.SetData(data)

		//得到当前连接的Request请求数据
		req := Request{
			conn: conn,
			msg:  msg,
		}
		//调用注册的路由方法
		go func(request iface.IRequest) {
			conn.Router.PreHandle(request)
			conn.Router.Handle(request)
			conn.Router.PostHandle(request)
		}(&req)

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
func (conn *Connection) GetTCPConnection() *net.TCPConn {
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

//发送数据的方法:将我们要发送给客户端的数据,先封包,再发送
func (conn *Connection) SendMsg(msgId uint32, data []byte) error {
	if conn.isClosed {
		return errors.New("connection closed when send msg")
	}

	//将data进行封包: dataLen, msgId, data的顺序
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("pack error msgId: ", msgId)
		return errors.New("pack error")
	}

	//将数据发送给客户端
	if _, err := conn.Conn.Write(binaryMsg); err != nil {
		fmt.Println("write msg id: ", msgId, "error ", err)
		return errors.New("conn write error")
	}

	return nil
}

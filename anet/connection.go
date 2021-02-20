package anet

import (
	"Alchemist/iface"
	"Alchemist/utils"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

type Connection struct {
	//当前Connection属于哪个Server
	TcpServer iface.IServer
	//socket TCP套接字
	Conn *net.TCPConn
	//连接的ID
	ConnID uint32
	//当前的连接状态(是否已经关闭)
	isClosed bool
	//等待连接被动退出的channel
	ExitChan chan bool
	//无缓冲的管道，用于读写Goroutine之间的消息通信
	msgChan chan []byte

	//MsgHandle模块
	MsgHandler iface.IMsgHandle

	//自定义连接属性集合
	property map[string]interface{}
	//保护连接属性的锁
	propertyLock sync.RWMutex
}

//初始化方法
func NewConnection(server iface.IServer, conn *net.TCPConn, connID uint32, msgHandler iface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		ExitChan:   make(chan bool),
		msgChan:    make(chan []byte),
		MsgHandler: msgHandler,
		property:   make(map[string]interface{}),
	}

	//将conn加入到ConnManager中
	c.TcpServer.GetConnMgr().Add(c)

	return c
}

func (conn *Connection) StartReader() {
	fmt.Println("[Reader Goroutine] is starting...")

	defer fmt.Println("[Reader is exit] connID= ", conn.ConnID, " remote addr is ", conn.GetRemoteAddr().String())
	defer conn.Stop()

	for {
		//创建一个拆包解包的对象
		dp := NewDataPack()

		//读取msgHead (8字节)
		headData := make([]byte, int(dp.GetHeadLen()))
		if _, err := io.ReadFull(conn.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error: ", err)
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

		if utils.GlobalObject.WorkerPoolSize > 0 {
			//将request发送给MsgHandler的工作池
			conn.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			//从路由中，找到注册绑定的Conn对应的Router调用
			//从路由中，根据绑定好的MsgID, 找到对应api处理业务执行
			go conn.MsgHandler.DoMsgHandler(&req)
		}

	}
}

func (conn *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine] is starting...")
	defer fmt.Println("[conn Writer exit] ", conn.GetRemoteAddr().String())

	//不断阻塞的等待channel的消息，进行写给客户端
	for {
		select {
		case data := <-conn.msgChan:
			//有数据要写给客户端
			if _, err := conn.Conn.Write(data); err != nil {
				fmt.Println("Send data error: ", err)
				return
			}
		case <-conn.ExitChan:
			//代表Reader已经退出，此时Writer也要退出
			return
		}
	}
}

//启动连接,让当前的连接准备开始工作
func (conn *Connection) Start() {
	fmt.Println("Conn Start()..ConnID=", conn.ConnID)

	//启动从当前连接的读数据业务
	go conn.StartReader()

	//启动从当前连接的写数据业务
	go conn.StartWriter()

	//调用注册的OnConnStart hook函数
	conn.TcpServer.CallOnConnStart(conn)
}

//停止连接, 结束当前连接的工作
func (conn *Connection) Stop() {
	fmt.Println("Conn Stop()..ConnID=", conn.ConnID)

	//如果当前连接已经关闭
	if conn.isClosed {
		return
	}
	conn.isClosed = true

	//调用注册的OnConnStop hook函数
	conn.TcpServer.CallOnConnStop(conn)

	//关闭socket连接
	conn.Conn.Close()

	//告知Writer关闭
	conn.ExitChan <- true

	//将当前连接从ConnMgr中移除
	conn.TcpServer.GetConnMgr().Remove(conn)
	//回收资源
	close(conn.ExitChan)
	close(conn.msgChan)
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

	//将数据发送给Writer
	conn.msgChan <- binaryMsg

	return nil
}

//添加自定义连接属性
func (conn *Connection) SetProperty(key string, value interface{}) {
	conn.propertyLock.Lock()
	defer conn.propertyLock.Unlock()

	conn.property[key] = value
}

//获取自定义连接属性
func (conn *Connection) GetProperty(key string) (interface{}, error) {
	conn.propertyLock.RLock()
	defer conn.propertyLock.RUnlock()

	if value, ok := conn.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

//移除自定义连接属性
func (conn *Connection) RemoveProperty(key string) {
	conn.propertyLock.Lock()
	defer conn.propertyLock.Unlock()

	delete(conn.property, key)
}

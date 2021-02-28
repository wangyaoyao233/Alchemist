package anet

import (
	"Alchemist/iface"
	"Alchemist/utils"
	"context"
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

	//MsgHandle模块
	MsgHandler iface.IMsgHandle
	//告知该链接已经退出/停止的channel
	ctx    context.Context
	cancel context.CancelFunc

	//无缓冲的管道，用于读写Goroutine之间的消息通信
	msgChan chan []byte
	//有缓冲的管道，用于读写Goroutine之间的消息通信
	msgBuffChan chan []byte

	sync.RWMutex
	//自定义连接属性集合
	property map[string]interface{}
	//保护连接属性的锁
	propertyLock sync.Mutex
	//当前连接的关闭状态
	isClosed bool
}

//初始化方法
func NewConnection(server iface.IServer, conn *net.TCPConn, connID uint32, msgHandler iface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:   server,
		Conn:        conn,
		ConnID:      connID,
		isClosed:    false,
		msgChan:     make(chan []byte),
		msgBuffChan: make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
		MsgHandler:  msgHandler,
		property:    make(map[string]interface{}),
	}

	//将conn加入到ConnManager中
	c.TcpServer.GetConnMgr().Add(c)

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine] is starting...")

	defer fmt.Println("[Reader is exit] connID= ", c.ConnID, " remote addr is ", c.GetRemoteAddr().String())
	defer c.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:

			//创建一个拆包解包的对象
			dp := NewDataPack()

			//读取msgHead (8字节)
			headData := make([]byte, dp.GetHeadLen())
			if _, err := io.ReadFull(c.Conn, headData); err != nil {
				fmt.Println("read msg head error: ", err)
				return
			}

			//拆包,得到dataLen, msgId, 放在msg中
			msg, err := dp.UnPack(headData)
			if err != nil {
				fmt.Println("unpack error", err)
				return
			}

			var data []byte
			//根据dataLen,再次读取data, 放在msg中
			if msg.GetDataLen() > 0 {
				data = make([]byte, msg.GetDataLen())
				if _, err := io.ReadFull(c.Conn, data); err != nil {
					fmt.Println("read msg data error: ", err)
					return
				}
			}
			msg.SetData(data)

			//得到当前连接的Request请求数据
			req := Request{
				conn: c,
				msg:  msg,
			}

			if utils.GlobalObject.WorkerPoolSize > 0 {
				//将request发送给MsgHandler的工作池
				c.MsgHandler.SendMsgToTaskQueue(&req)
			} else {
				//从路由中，找到注册绑定的Conn对应的Router调用
				//从路由中，根据绑定好的MsgID, 找到对应api处理业务执行
				go c.MsgHandler.DoMsgHandler(&req)
			}
		}
	}
}

func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine] is starting...")
	defer fmt.Println("[conn Writer exit] ", c.GetRemoteAddr().String())

	//不断阻塞的等待channel的消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error: ", err, "Conn Writer exit")
				return
			}
		case data, ok := <-c.msgBuffChan:
			if ok {
				//有数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send Buff data error: ", err, "Conn Writer exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is Closed")
				break
			}
		case <-c.ctx.Done():
			//代表Reader已经退出，此时Writer也要退出
			return
		}
	}
}

//启动连接,让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()..ConnID=", c.ConnID)
	c.ctx, c.cancel = context.WithCancel(context.Background())

	//启动从当前连接的读数据业务
	go c.StartReader()

	//启动从当前连接的写数据业务
	go c.StartWriter()

	//调用注册的OnConnStart hook函数
	c.TcpServer.CallOnConnStart(c)
}

//停止连接, 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()..ConnID=", c.ConnID)

	c.Lock()
	defer c.Unlock()

	//调用注册的OnConnStop hook函数
	c.TcpServer.CallOnConnStop(c)

	//如果当前连接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//关闭socket连接
	c.Conn.Close()

	//告知Writer关闭
	c.cancel()

	//将当前连接从ConnMgr中移除
	c.TcpServer.GetConnMgr().Remove(c)
	//回收资源
	close(c.msgBuffChan)
}

//获取当前连接的conn对象(套接字)
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取客户端连接的地址和端口
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//发送数据的方法:将我们要发送给客户端的数据,先封包,再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	c.RLock()
	if c.isClosed == true {
		c.RUnlock()
		return errors.New("connection closed when send msg")
	}
	c.RUnlock()

	//将data进行封包: dataLen, msgId, data的顺序
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("pack error msgId: ", msgId)
		return errors.New("pack error")
	}

	//将数据发送给Writer
	c.msgChan <- binaryMsg

	return nil
}

func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	c.RLock()
	if c.isClosed == true {
		c.RUnlock()
		return errors.New("Connection closed when send msg")
	}
	c.RUnlock()

	//将data封包： dataLen, MsgId, dataLen
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("pack error msgId: ", msgId)
		return errors.New("pack error")
	}

	//将数据发送给Writer
	c.msgBuffChan <- binaryMsg

	return nil
}

//添加自定义连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

//获取自定义连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

//移除自定义连接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

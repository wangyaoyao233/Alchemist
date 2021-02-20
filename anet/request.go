package anet

import "Alchemist/iface"

//将连接和请求的数据 包装到一个Request中
type Request struct {
	//和客户端建立的连接
	conn iface.IConnection
	//客户端请求的数据
	msg iface.IMessage
}

//得到连接
func (r *Request) GetConnection() iface.IConnection {
	return r.conn
}

//得到请求的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}

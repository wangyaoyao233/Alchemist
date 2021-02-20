package iface

//将连接和请求的数据 包装到一个Request中
type IRequest interface {
	//得到连接
	GetConnection() IConnection
	//得到请求的数据
	GetData() []byte
	GetMsgID() uint32
}

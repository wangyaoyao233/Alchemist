package iface

//针对Message进行TLV格式封装
type IDataPack interface {
	//获取包的头的长度方法
	GetHeadLen() uint32
	//封包
	Pack(msg IMessage) ([]byte, error)
	//拆包
	UnPack([]byte) (IMessage, error)
}

package iface

//将请求的消息封装到一个Message中
type IMessage interface {
	GetMsgID() uint32
	SetMsgID(uint32)
	GetDataLen() uint32
	SetDataLen(uint32)
	GetData() []byte
	SetData([]byte)
}

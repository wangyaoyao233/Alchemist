package anet

type Message struct {
	DataLen uint32
	Id      uint32
	Data    []byte
}

//创建一个Message信息方法
func NewMessage(id uint32, data []byte) *Message {
	return &Message{
		DataLen: uint32(len(data)),
		Id:      id,
		Data:    data,
	}
}

func (m *Message) GetMsgID() uint32 {
	return m.Id
}
func (m *Message) SetMsgID(id uint32) {
	m.Id = id
}
func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}
func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}
func (m *Message) GetData() []byte {
	return m.Data
}
func (m *Message) SetData(data []byte) {
	m.Data = data
}

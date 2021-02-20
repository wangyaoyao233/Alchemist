package anet

import (
	"Alchemist/iface"
	"Alchemist/utils"
	"bytes"
	"encoding/binary"
	"errors"
)

type DataPack struct {
}

//实例化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

//获取包的头的长度方法
func (dp *DataPack) GetHeadLen() uint32 {
	//DataLen uint32 (4) + MsgID uint32 (4)
	return 8
}

//封包
func (dp *DataPack) Pack(msg iface.IMessage) ([]byte, error) {
	//创建一个存放bytes直接的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//将dataLen写入dataBuff
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	//将msgID写入dataBuff
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}
	//将data写入dataBuff
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

//拆包:只需要将包的Head信息读出来
func (dp *DataPack) UnPack(headBinaryData []byte) (iface.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(headBinaryData)

	msg := &Message{}

	//解压head信息
	//读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	//读msgId
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//判断dataLen是否已经超出允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data recv")
	}

	return msg, nil
}

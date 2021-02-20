package main

import (
	"Alchemist/anet"
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	fmt.Println("client1 start..")
	time.Sleep(1 * time.Second)
	//1.直接连接远程服务器, 得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		fmt.Println("client start err", err)
		return
	}

	for {
		//发送封包的message消息
		dp := anet.NewDataPack()
		binaryMsg, err := dp.Pack(anet.NewMessage(1, []byte("client test ...")))
		if err != nil {
			fmt.Println("client pack error: ", err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("write error: ", err)
			return
		}

		//接受服务器回复的message消息
		//1.先读取流中的head部分(8)
		headData := make([]byte, int(dp.GetHeadLen()))
		if _, err := io.ReadFull(conn, headData); err != nil {
			fmt.Println("read head error: ", err)
			return
		}
		//2.Unpack得到dataLen, msgId
		msgHead, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error: ", err)
			return
		}
		if msgHead.GetDataLen() > 0 {
			//3.再根据dataLen,读取data
			msg := msgHead.(*anet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg error: ", err)
				return
			}
			fmt.Println("recv from server: msgID:", msg.GetMsgID(), "len:", msg.GetDataLen(), "data:", string(msg.GetData()))
		}

		//cpu阻塞
		time.Sleep(1 * time.Second)
	}
}

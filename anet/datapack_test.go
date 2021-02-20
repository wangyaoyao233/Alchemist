package anet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

//datapack unit test
func TestDataPack(t *testing.T) {
	//模拟的服务器
	//1.创建socketTCP
	listener, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("server listen err ", err)
		return
	}

	//创建一个go 曾在从客户端处理业务
	go func() {
		//2.从客户端读取数据, 拆包处理
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accpet err ", err)
				return
			}

			go func(conn net.Conn) {
				//处理客户端的请求
				//--->拆包过程
				dp := NewDataPack()
				for {
					//1.从conn读, 把包的head读出来
					headData := make([]byte, int(dp.GetHeadLen()))
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error: ", err)
						return
					}

					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("unpack error: ", err)
						return
					}
					//2.跟据head中的dataLen再读取data
					if msgHead.GetDataLen() > 0 {
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetDataLen())

						//再次从io流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("unpack data error", err)
							return
						}

						fmt.Println("-->Recv MsgID:", msg.Id, " msg.DataLen:", msg.DataLen, " data:", string(msg.Data))
					}

				}

			}(conn)
		}
	}()

	//模拟客户端
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("client dial err ", err)
		return
	}

	//创建一个封包对象
	dp := NewDataPack()

	//模拟粘包过程, 封装2个msg一起发送
	msg1 := &Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte{'H', 'E', 'L', 'L', 'O'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err", err)
		return
	}

	msg2 := &Message{
		Id:      2,
		DataLen: 2,
		Data:    []byte{'N', 'I'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 err", err)
		return
	}

	sendData1 = append(sendData1, sendData2...)
	conn.Write(sendData1)

	select {}
}

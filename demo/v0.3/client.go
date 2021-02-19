package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("client start..")
	time.Sleep(1 * time.Second)
	//1.直接连接远程服务器, 得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		fmt.Println("client start err", err)
		return
	}

	for {
		//2.连接调用write写数据
		_, err := conn.Write([]byte("Hello v0.3"))
		if err != nil {
			fmt.Println("write conn err", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error", err)
			return
		}

		fmt.Printf("server call back: %s, cnt = %d\n", buf, cnt)

		//cpu阻塞
		time.Sleep(1 * time.Second)
	}
}

package main

import (
	"Alchemist/anet"
	"Alchemist/iface"
	"Alchemist/mmo_game/apis"
	"Alchemist/mmo_game/core"
	"fmt"
)

func OnConnectionAdd(conn iface.IConnection) {
	//创建一个Player对象
	player := core.NewPlayer(conn)

	//给客户端发送MsgID:1的消息:同步当前Player的ID给客户端
	player.SyncPid()
	//给客户端发送MsgID:200的消息:同步当前Player的初始位置给客户端
	player.BroadCastStartPosition()

	//将新玩家加入世界管理模块
	core.WorldMgrObj.AddPlayer(player)

	//将该链接绑定一个pid, 玩家ID的属性
	conn.SetProperty("pid", player.Pid)

	//同步新玩家位置(看见别人+别人看见我)
	player.SyncSurrounding()

	fmt.Println("====>Player ID:", player.Pid, "is arrived..")

}

func OnConnectionLost(conn iface.IConnection) {
	//通过连接属性获得当前连接说绑定的pid
	pid, _ := conn.GetProperty("pid")
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//玩家下线业务(给周围玩家广播MsgId:201)
	if pid != nil {
		player.Offline()
	}
	fmt.Println("====>Player pid:", pid, "offline...")
}

func main() {
	//创建server
	s := anet.NewServer("MMO Game Server")

	//注册连接的创建和销毁的hook函数
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)

	//注册路由
	s.AddRouter(2, &apis.WorldChatApi{})
	s.AddRouter(3, &apis.MoveApi{})

	//启动
	s.Serve()
}

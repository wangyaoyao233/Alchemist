package apis

import (
	"Alchemist/anet"
	"Alchemist/iface"
	"Alchemist/mmo_game/core"
	"Alchemist/mmo_game/pb"
	"fmt"

	"google.golang.org/protobuf/proto"
)

//世界聊天 路由业务
type WorldChatApi struct {
	anet.BaseRouter
}

func (*WorldChatApi) Handle(request iface.IRequest) {
	//1.解析客户端传进来的proto协议
	proto_msg := &pb.Talk{}
	err := proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		fmt.Println("talk Unmarshal error,", err)
		return
	}
	//2.当前的聊天数据属于哪个玩家发送的
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty pid error,", err)
		request.GetConnection().Stop()
		return
	}
	//3.根据pid得到对应的player对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//4.将这个消息广播给其他全部在线玩家
	player.Talk(proto_msg.Content)
}

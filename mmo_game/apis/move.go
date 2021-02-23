package apis

import (
	"Alchemist/anet"
	"Alchemist/iface"
	"Alchemist/mmo_game/core"
	"Alchemist/mmo_game/pb"
	"fmt"

	"google.golang.org/protobuf/proto"
)

//玩家移动
type MoveApi struct {
	anet.BaseRouter
}

func (m *MoveApi) Handle(request iface.IRequest) {
	//1.解析客户端传递进来的proto协议
	proto_msg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		fmt.Println("Move: Position Unmarshal error", err)
		return
	}

	//2.得到当前发送位置的事哪个玩家
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty pid error", err)
		request.GetConnection().Stop()
		return
	}
	fmt.Printf("Player pid:%d,move(%f,%f,%f,%f)\n", pid, proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)

	//将当前玩家的位置广播给周围玩家
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	player.UpdatePos(proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)
}

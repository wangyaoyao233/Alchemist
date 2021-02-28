package core

import (
	"Alchemist/iface"
	"Alchemist/mmo_game/pb"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
)

type Player struct {
	//玩家ID
	Pid int32
	//当前玩家的连接(用于和客户端的连接)
	Conn iface.IConnection
	//位置信息
	X float32
	Y float32
	Z float32
	V float32 //旋转的0-360角度
}

//PlayerID 生成器
var PidGen int32 = 1  //用来生成玩家ID的计数器
var IdLock sync.Mutex //保护PidGen的Mutex

//创建一个玩家
func NewPlayer(conn iface.IConnection) *Player {
	//生成一个玩家ID
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()

	p := &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(100 + rand.Intn(10)),
		Y:    0,
		Z:    float32(120 + rand.Intn(20)),
		V:    0,
	}
	return p
}

//发送给客户端消息的方法
//主要是将pb的protobuf数据序列化之后, 再调用框架的Connection.SendMsg方法
func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	//1.将proto Message结构体序列化 转换成二进制
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err,", err)
		return
	}
	//2.将二进制文件 通过框架的SendMsg进行TLV格式发送给客户端
	if p.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}

	if err := p.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("player sendmsg error,", err)
		return
	}

	return
}

//告知客户端pid,同步已经生成的玩家ID给客户端
func (p *Player) SyncPid() {
	//组建MsgID:1 的proto数据
	data := &pb.SyncPid{
		Pid: p.Pid,
	}
	//发送消息
	p.SendMsg(1, data)
}

//同步玩家自己的出生地点
func (p *Player) BroadCastStartPosition() {
	//组建MsgID:200 的proto数据
	data := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2, //2代表广播的位置坐标
		Data: &pb.BroadCast_P{
			//Position
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	//同步玩家自己的初始位置
	p.SendMsg(200, data)
}

//玩家广播世界聊天消息
func (p *Player) Talk(content string) {
	//1.组建一个MsgId:200 proto数据
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  1, //tp-1代表聊天广播
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}
	//2.得到当前世界所有的在线玩家
	players := WorldMgrObj.GetAllPlayers()

	//3.向所有玩家广播聊天信息
	for _, player := range players {
		//player分别给对应的客户端发送消息
		player.SendMsg(200, proto_msg)
	}

}

//同步玩家新上线的位置
func (p *Player) SyncSurrounding() {
	//1.获取当前玩家周围的玩家
	players := p.GetSurroundingPlayers()

	//2.将新玩家位置发给周围玩家
	//2.1组建MsgId:200 proto数据
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	//2.2发送给周围玩家
	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}

	//3.将周围玩家位置发给新玩家
	//3.1组建MsgId:202 proto数据
	//3.1.1制作pb.Player slice
	players_proto_msg := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		//制作一个message Player
		p_msg := &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		players_proto_msg = append(players_proto_msg, p_msg)
	}

	//3.1.2将玩家集合放在message SyncPlayer
	sync_players_proto_msg := &pb.SyncPlayers{
		Ps: players_proto_msg[:],
	}
	//3.2发送给新玩家
	p.SendMsg(202, sync_players_proto_msg)
}

//广播并更新当前玩家的位置
func (p *Player) UpdatePos(x, y, z, v float32) {

	//触发消失视野和添加视野
	//计算旧格子gid
	oldGid := WorldMgrObj.AoiMgr.GetGidByPos(p.X, p.Z)
	//计算新格子gid
	newGid := WorldMgrObj.AoiMgr.GetGidByPos(x, z)

	//1.更新当前玩家player对象的坐标
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v

	if oldGid != newGid {
		//触发grid切换
		//把pid从旧的aoi格子中删除
		WorldMgrObj.AoiMgr.RemovePidFromGrid(int(p.Pid), oldGid)
		//把pid添加到新的aoi格子中
		WorldMgrObj.AoiMgr.AddPidToGrid(int(p.Pid), newGid)

		_ = p.OnExchangeAoiGrid(oldGid, newGid)
	}

	//2.组建MsgId:200,Tp=4的广播消息
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  4,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	//获取当前玩家的周围玩家
	players := p.GetSurroundingPlayers()
	//给周围每个玩家发送广播消息
	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}
}

func (p *Player) OnExchangeAoiGrid(oldGid, newGid int) error {
	//获取旧的九宫格成员
	oldGrids := WorldMgrObj.AoiMgr.GetSurroundGridsByGid(oldGid)

	//为旧的九宫格成员建立哈希表,用来快速查找
	oldGridsMap := make(map[int]bool, len(oldGrids))
	for _, grid := range oldGrids {
		oldGridsMap[grid.GID] = true
	}

	//获取新的九宫格成员
	newGrids := WorldMgrObj.AoiMgr.GetSurroundGridsByGid(newGid)
	//为新的九宫格成员建立哈希表,用来快速查找
	newGridsMap := make(map[int]bool, len(newGrids))
	for _, grid := range newGrids {
		newGridsMap[grid.GID] = true
	}

	//-------处理视野消失
	offlineMsg := &pb.SyncPid{
		Pid: p.Pid,
	}

	//找到旧的九宫格中出现,且在新的九宫格中没有出现的格子
	leavingGrids := make([]*Grid, 0)
	for _, grid := range oldGrids {
		if _, ok := newGridsMap[grid.GID]; !ok {
			leavingGrids = append(leavingGrids, grid)
		}
	}

	//获取需要消失的格子中的全部玩家
	for _, grid := range leavingGrids {
		players := WorldMgrObj.GetPlayersByGid(grid.GID)
		for _, player := range players {
			//让自己在其他玩家的客户端消失
			player.SendMsg(201, offlineMsg)
			//将其他玩家信息 在自己的客户端中消失
			anotherOfflineMsg := &pb.SyncPid{
				Pid: player.Pid,
			}
			p.SendMsg(201, anotherOfflineMsg)
			time.Sleep(200 * time.Millisecond)
		}
	}

	//--------处理视野出现
	//找到在小的九宫格出现, 但是没有在旧的九宫格出现的格子
	enteringGrids := make([]*Grid, 0)
	for _, grid := range newGrids {
		if _, ok := oldGridsMap[grid.GID]; !ok {
			enteringGrids = append(enteringGrids, grid)
		}
	}

	onlineMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	//获取需要显示格子的全部玩家
	for _, grid := range enteringGrids {
		players := WorldMgrObj.GetPlayersByGid(grid.GID)

		for _, player := range players {
			//让自己出现在其他人视野中
			player.SendMsg(200, onlineMsg)

			//让其他人出现在自己的视野
			anotherOnlineMsg := &pb.BroadCast{
				Pid: player.Pid,
				Tp:  2,
				Data: &pb.BroadCast_P{
					P: &pb.Position{
						X: player.X,
						Y: player.Y,
						Z: player.Z,
						V: player.V,
					},
				},
			}

			time.Sleep(200 * time.Millisecond)
			p.SendMsg(200, anotherOnlineMsg)
		}
	}

	return nil
}

//获取当前玩家的周围九宫格玩家
func (p *Player) GetSurroundingPlayers() []*Player {
	//得到当前九宫格内所有玩家PID
	pids := WorldMgrObj.AoiMgr.GetPidsByPos(p.X, p.Z)

	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}

	return players
}

func (p *Player) Offline() {
	//得到当前玩家周围的玩家
	players := p.GetSurroundingPlayers()
	//给周围玩家广播MsgId:201
	proto_msg := &pb.SyncPid{
		Pid: p.Pid,
	}

	for _, player := range players {
		player.SendMsg(201, proto_msg)
	}

	//将当前玩家从AOI管理模块删除
	WorldMgrObj.AoiMgr.RemoveFromGridByPos(int(p.Pid), p.X, p.Z)
	//将当前玩家从世界管理器中删除
	WorldMgrObj.RemovePlayerByPid(p.Pid)
}

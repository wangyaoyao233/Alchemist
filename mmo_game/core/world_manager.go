package core

import "sync"

//格子信息的宏
const (
	AOI_MIN_X  int = 0
	AOI_MAX_X  int = 500
	AOI_CNTS_X int = 10
	AOI_MIN_Y  int = 0
	AOI_MAX_Y  int = 500
	AOI_CNTS_Y int = 10
)

//当前游戏的世界总管理模块
type WorldManager struct {
	// AOIManager 当前世界地图AOI的管理模块
	AoiMgr *AOIManager
	// 当前全部在线的Player集合 key-Pid, value-Player
	Players map[int32]*Player
	// 保护Player集合的锁
	pLock sync.RWMutex
}

//提供一个对外的世界管理模块的句柄(全局)
var WorldMgrObj *WorldManager

//初始化方法
func init() {
	WorldMgrObj = &WorldManager{
		AoiMgr: NewAOIManager(
			AOI_MIN_X,
			AOI_MAX_X,
			AOI_CNTS_X,
			AOI_MIN_Y,
			AOI_MAX_Y,
			AOI_CNTS_Y,
		),
		Players: make(map[int32]*Player),
	}
}

//添加玩家
func (wm *WorldManager) AddPlayer(p *Player) {
	wm.pLock.Lock()
	wm.Players[p.Pid] = p
	wm.pLock.Unlock()

	//将player添加到AOIManager中
	wm.AoiMgr.AddToGridByPos(int(p.Pid), p.X, p.Z)
}

//删除玩家,通过Pid
func (wm *WorldManager) RemovePlayerByPid(pid int32) {
	wm.pLock.Lock()
	delete(wm.Players, pid)
	wm.pLock.Unlock()
}

//通过玩家ID查询Player对象
func (wm *WorldManager) GetPlayerByPid(pid int32) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	return wm.Players[pid]
}

//获取全部在线玩家
func (wm *WorldManager) GetAllPlayers() []*Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	players := make([]*Player, 0)
	for _, p := range wm.Players {
		players = append(players, p)
	}

	return players
}

//获取指定gid中的所有player信息
func (wm *WorldManager) GetPlayersByGid(gid int) []*Player {
	//通过gid获取对应格子中的所有pid
	pids := wm.AoiMgr.grids[gid].GetPlayerIDs()

	//通过pid找到对应的player对象
	players := make([]*Player, 0, len(pids))
	wm.pLock.RLock()
	for _, pid := range pids {
		players = append(players, wm.Players[int32(pid)])
	}
	wm.pLock.RUnlock()

	return players
}

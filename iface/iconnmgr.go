package iface

//连接管理模块
type IConnManager interface {
	Add(conn IConnection)
	Remove(conn IConnection)
	Get(connID int) (IConnection, error)
	Len() int
	Clear()
}

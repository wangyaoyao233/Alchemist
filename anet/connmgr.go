package anet

import (
	"Alchemist/iface"
	"errors"
	"fmt"
	"sync"
)

type ConnManager struct {
	//管理的连接集合
	connections map[int]iface.IConnection
	//保护连接集合的读写锁
	connLock sync.RWMutex
}

//创建当前管理模块方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[int]iface.IConnection),
	}
}

func (connMgr *ConnManager) Add(conn iface.IConnection) {
	//保护共享资源map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//将conn加入map
	connMgr.connections[int(conn.GetConnID())] = conn

	fmt.Println("connID =", conn.GetConnID(), "add to ConnManager successfully: conn num:", connMgr.Len())
}

func (connMgr *ConnManager) Remove(conn iface.IConnection) {
	//保护共享资源map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除连接信息
	delete(connMgr.connections, int(conn.GetConnID()))

	fmt.Println("connID =", conn.GetConnID(), "remove from ConnManager succ: conn num:", connMgr.Len())
}

func (connMgr *ConnManager) Get(connID int) (iface.IConnection, error) {
	//保护共享资源map，加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

func (connMgr *ConnManager) Clear() {
	//保护共享资源map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除conn并停止conn的工作
	for connID, conn := range connMgr.connections {
		//停止
		conn.Stop()
		//删除
		delete(connMgr.connections, connID)
	}

	fmt.Println("Clear All connections succ, conn num=", connMgr.Len())
}

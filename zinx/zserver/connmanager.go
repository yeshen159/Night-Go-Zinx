package zserver

import (
	"Night/zinx/ziserver"
	"errors"
	"fmt"
	"sync"
)

/*
  链接管理模块
*/

type ConnManager struct {
	connections map[uint32]ziserver.IConnection // 管理的链接集合
	connLock    sync.RWMutex                    // 保护链接集合的读写锁
}

//创建当前链接的方法

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziserver.IConnection),
	}
}

//添加链接
func (connMgr *ConnManager) Add(conn ziserver.IConnection) {
	//保护共享资源map, 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//将conn加入到ConnManager中
	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("connID = ", conn.GetConnID(), " add to ConnManager successfully: conn num =", connMgr.Len())

}

//删除链接
func (connMgr *ConnManager) Remove(conn ziserver.IConnection) {
	//保护共享资源map, 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("connID = ", conn.GetConnID(), " remove from ConnManager successfully: conn num =", connMgr.Len())

}

//根据connID获取链接
func (connMgr *ConnManager) Get(connID uint32) (ziserver.IConnection, error) {
	//保护共享资源map, 加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not FOUND!")
	}

}

//得到当前链接总数
func (connMgr *ConnManager) Len() int {
	//Add方法调用了len方法,加锁会死锁
	return len(connMgr.connections)
}

//清除并终止所有链接
func (connMgr *ConnManager) ClearConn() {
	//保护共享资源map, 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除conn并停止conn工作
	for connID, conn := range connMgr.connections {
		//停止
		conn.Stop()
		//删除
		delete(connMgr.connections, connID)
	}

	fmt.Println("Clear All connections succ! conn num = ", connMgr.Len())
}

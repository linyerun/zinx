package znet

import (
	"errors"
	"fmt"
	"github.com/linyerun/zinx/ziface"
	"sync"
)

type ConnsManager struct {
	conns    map[uint32]ziface.IConnection
	connLock sync.RWMutex
}

func newConnsManager() ziface.IConnsManager {
	return &ConnsManager{
		conns: map[uint32]ziface.IConnection{},
	}
}

func (c *ConnsManager) Add(conn ziface.IConnection) {
	c.connLock.Lock()
	c.conns[conn.GetConnID()] = conn
	c.connLock.Unlock()
	fmt.Println("add conn successfully,connId =", conn.GetConnID())
}

func (c *ConnsManager) Remove(connId uint32) {
	c.connLock.Lock()
	delete(c.conns, connId)
	c.connLock.Unlock()
	fmt.Println("remove conn from conns map where connId =", connId)
}

func (c *ConnsManager) Get(connId uint32) (ziface.IConnection, error) {
	conn, ok := c.conns[connId]
	if ok {
		return conn, nil
	}
	return nil, errors.New(fmt.Sprintf("can not find conn by connId = %d", connId))
}

func (c *ConnsManager) Len() int {
	c.connLock.RLock()
	defer c.connLock.RUnlock()
	return len(c.conns)
}

func (c *ConnsManager) ClearConns() {
	c.connLock.Lock()
	for _, conn := range c.conns {
		//关闭连接
		//这个可以作为一个信号，叫connection读协程关闭
		if err := conn.GetTCPConnection().Close(); err != nil {
			return
		}
		//delete(c.conns, conn.GetConnID()) //这一步交给connection做吧
	}
	c.connLock.Unlock()
}

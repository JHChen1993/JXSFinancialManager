package service

import (
	"net"
	"sync"
)

// ClientConn 客户端连接类
type ClientConn struct {
	netId  int64
	netObj net.Conn
	svPtr  *Server

	one *sync.Once
	wg  *sync.WaitGroup
}

// Start 启动
func (ct *ClientConn) Start() {

}

func (ct *ClientConn) Close() {
	ct.one.Do(func() {
		ct.netObj.Close()
		ct.wg.Wait()
		ct.svPtr.wg.Done()
	})
}

// NewClient 新建客户端类
func NewClientConn(cnid int64, conn net.Conn, sv *Server) *ClientConn {
	cl := ClientConn{
		netId: cnid,
		svPtr: sv,
		svPtr: sv,
		one:   snc.Once{},
		wg:    &sync.WaitGroup{},
	}
	return &cl
}

// readLoop 读数据循环
func readLoop(ct interface{}, wg *sync.WaitGroup) {
	clt := ct.(ClientConn)
	
	for{
		select {
			case <-
		}
	}
}

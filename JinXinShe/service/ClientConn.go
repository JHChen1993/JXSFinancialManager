package service

import (
	"fmt"
	"net"
	"sync"
)

// ClientConn 客户端连接类
type ClientConn struct {
	netId  int64
	netObj net.Conn
	svPtr  *Server

	one     *sync.Once
	wg      *sync.WaitGroup
	closeCh chan struct{}
}

// Start 启动
func (ct *ClientConn) Start() {
	loopers := []func(interface{}, *sync.WaitGroup){readLoop, writeLoop}
	for _, l := range loopers {
		ct.wg.Add(1)
		l(ct, ct.wg)
	}
}

func (ct *ClientConn) Close() {
	ct.one.Do(func() {
		close(ct.closeCh)
		ct.netObj.Close()
		ct.wg.Wait()
		ct.svPtr.wg.Done()
	})
}

// NewClient 新建客户端类
func NewClientConn(cnid int64, conn net.Conn, sv *Server) *ClientConn {
	cl := ClientConn{
		netId:   cnid,
		netObj:  conn,
		svPtr:   sv,
		one:     &sync.Once{},
		wg:      &sync.WaitGroup{},
		closeCh: make(chan struct{}),
	}
	return &cl
}

// readLoop 读数据循环
func readLoop(ct interface{}, wg *sync.WaitGroup) {
	clt := ct.(ClientConn)

	defer func() {
		clt.wg.Done()
		clt.Close()
	}()

	for {
		select {
		case <-clt.closeCh:
			fmt.Fprintln("recive clientconn close signal ")
			return
		case <-clt.svPtr.clsChan:
			fmt.Fprintln("recive Server Close singal")
			return
		default:
			clt.netObj.Read()
		}
	}
}

// writeLoop 写入数据循环
func writeLoop(ct interface{}, wg *sync.WaitGroup) {

}

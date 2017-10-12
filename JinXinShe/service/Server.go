// package service 服务器包
package service

import (
	"net"
	"sync"
)

var (
	clientId *AutomicInt64
)

// init 初始化
func init() {
	clientId = NewAutomicInt64(0)
}

// Server 服务器对象
type Server struct {
	clsChan     chan struct{} // 通知conn对象关闭
	ClientConns sync.Map
	wg          *sync.WaitGroup
}

// Start 启动服务器
func (sv *Server) Start(ls net.Listener) error {
	// <TODO: 应该做一些配置信息检查>
	defer func() {
		ls.Close()
	}()
	for {
		con, err := ls.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				// <TODO: 发生临时错误应该sleep几秒>
				continue
			}
			return err
		}
		// 通过con建立连接客户端
		sv.newConnect(con)
	} // loop
}

// newConnect 新建客户端连接
func (sv *Server) newConnect(ccon net.Conn) {
	cnnid := clientId.GetAndIncrement()
	scc := NewClientConn(cnnid, ccon, sv.clsChan)
	sv.ClientConns.Store(connid, scc)
	// 连接开始工作
	scc.Start()
}

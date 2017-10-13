// package service 服务器包
package service

import (
	"fmt"
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
	endSever    chan struct{} // 通知conn对象关闭
	clientConns *sync.Map
	wg          *sync.WaitGroup
	ls          net.Listener
}

// NewServer 创建服务器对象
func NewServer() *Server {
	sv := Server{
		endSever:    make(chan struct{}),
		clientConns: &sync.Map{},
		wg:          &sync.WaitGroup{},
	}
	return &sv
}
func (sv *Server) Stop() {
	close(sv.endSever)
	sv.ls.Close()
	sv.wg.Wait()
	fmt.Println("service stop！")
}

// Start 启动服务器
func (sv *Server) Start(ls net.Listener) error {
	sv.wg.Add(1)
	sv.ls = ls
	// <TODO: 应该做一些配置信息检查>
	defer func() {
		sv.wg.Done()
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
	scc := NewClientConn(cnnid, ccon, sv)
	sv.clientConns.Store(cnnid, scc)
	// 连接开始工作
	scc.Start()
}

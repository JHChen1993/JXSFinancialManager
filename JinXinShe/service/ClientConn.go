package service

import (
	"fmt"
	"net"
	"sync"
)

// ClientConn 客户端连接类
type ClientConn struct {
	netId   int64
	netConn net.Conn
	svPtr   *Server

	one      *sync.Once
	wg       *sync.WaitGroup
	endConn  chan struct{}
	sendMsg  chan []byte
	transMsg chan []byte
}

// Start 启动
func (c *ClientConn) Start() {
	loopers := []func(WriteCloser, *sync.WaitGroup){readLoop, writeLoop}
	for _, l := range loopers {
		c.wg.Add(1)
		l(c, c.wg)
	}
}

func (c *ClientConn) Close() {
	c.one.Do(func() {
		close(c.endConn)
		c.netConn.Close()
		c.wg.Wait()
		c.svPtr.wg.Done()
	})
}

// NewClient 新建客户端类
func NewClientConn(cnid int64, conn net.Conn, sv *Server) *ClientConn {
	cl := ClientConn{
		netId:   cnid,
		netConn: conn,
		svPtr:   sv,
		one:     &sync.Once{},
		wg:      &sync.WaitGroup{},
		endConn: make(chan struct{}),
	}
	return &cl
}

// readLoop 读数据循环
func readLoop(c WriteCloser, wg *sync.WaitGroup) {
	var (
		netConn  net.Conn
		transMsg chan<- []byte
		cEnd     <-chan struct{}
		sEnd     <-chan struct{}

		readBuff []byte
	)
	switch c := c.(type) {
	case *ClientConn:
		netConn = c.netConn
		transMsg = c.transMsg
		cEnd = c.endConn
		sEnd = c.svPtr.endSever
		readBuff = make([]byte, 4096)
	}
	defer func() {
		wg.Done()
		c.Close()
	}()

	for {
		select {
		case <-cEnd:
			fmt.Println("recive clientconn close Signal ")
			return
		case <-sEnd:
			fmt.Println("recive Server Close Singal")
			return
		default:
			n, err := netConn.Read(readBuff)
			if err != nil {
				fmt.Println(err)
				return
			}
			transMsg <- readBuff[:n]
		}
	}
}

// writeLoop 写入数据循环
func writeLoop(c WriteCloser, wg *sync.WaitGroup) {
	var (
		netConn net.Conn
		sendMsg <-chan []byte
		cEnd    <-chan struct{}
		sEnd    <-chan struct{}
	)
	switch c := c.(type) {
	case *ClientConn:
		netConn = c.netConn
		sendMsg = c.sendMsg
		cEnd = c.endConn
		sendMsg = c.svPtr.endSever
	}

	defer func() {
		for msg := range <-sendMsg {
			netConn.Write(msg)
		}
		wg.Done()
		c.Close()
	}()

	for {
		select {
		case <-cEnd:
			fmt.Println("writeLoop Recive Connect Close Signal")
			return
		case <-sEnd:
			fmt.Println("writeLoop Recive Server Close Signal")
			return
		case msg := <-sendMsg:
			wnb, err := netConn.Write(msg)
			if err != nil {
				fmt.Println("writeLoop Write Msg has failt")
				return
			}
		}
	}
}

// handLoop 消息处理和超时处理<TODO:>
func handLoop(c writeLoop, wg *sync.WaitGroup) {
	var (
		netId    int64
		transMsg <-chan []byte
		cEnd     <-chan struct{}
		sEnd     <-chan struct{}
	)
	switch c := c.(type) {
	case *ClientConn:
		netId = c.netId
		transMsg = c.transMsg
		cEnd = c.endConn
		sEnd = c.svPtr.endSever
	}

	defer func() {
		wg.Done()
		c.Close()
	}()

	for {
		select {
		case <-cEnd:
			fmt.Println("handLoop Recive Connect Close Signal")
			return
		case <-sEnd:
			fmt.Println("handLoop Recive Server Close Signal")
			return
		case msg := <-transMsg:
		}
	}
}

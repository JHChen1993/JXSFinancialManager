package main

import (
	"fmt"
	"net"
	"os"

	"github.com/JH_Chen1992/JinXinShe/service"
)

func main() {
	sv := service.NewServer()
	ls, err := net.Listen("tcp", "127.0.0.1:6588")
	if err != nil {
		fmt.Println("listen failt")
		return
	}
	go sv.Start(ls)
	fmt.Println("输入xx关闭：")
	os.Stdin.Read(make([]byte, 1))
	sv.Stop()
}

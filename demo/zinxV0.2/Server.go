package main

import (
	"Night/zinx/zserver"
	"fmt"
)

/*
	基于Zinx框架来开发的 服务器端应用程序
*/

func main() {
	fmt.Println("vim-go")

	//1 创建一个server句柄,使用Zinx的api
	s := zserver.NewServer("[nightzinx V0.2]")

	//2 启动server
	s.Serve()
}

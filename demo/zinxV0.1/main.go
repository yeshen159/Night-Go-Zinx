package main

import (
	"Night/zinx/zserver"
	"fmt"
)

func main() {
	fmt.Println("vim-go")

	//1 创建一个server句柄,使用Zinx的api
	s := zserver.NewServer("[nightzinx V0.1]")

	//2 启动server
	s.Serve()
}

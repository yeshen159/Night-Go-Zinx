package main

import (
	"Night/zinx/ziserver"
	"Night/zinx/zserver"
	"fmt"
)

/*
	基于Zinx框架来开发的 服务器端应用程序
*/

//ping test 自定义路由
type PingRouter struct {
	zserver.BaseRouter
}

//Test PreHandle
func (this *PingRouter) PreHandle(request ziserver.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping..."))

	if err != nil {
		fmt.Println("call back before ping error")
	}
}

//Test Handle
func (this *PingRouter) Handle(request ziserver.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping... ping... ping..."))

	if err != nil {
		fmt.Println("call back ping ping ping error")
	}

}

//Test PostHandle
func (this *PingRouter) PostHandle(request ziserver.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping..."))

	if err != nil {
		fmt.Println("call back after ping error")
	}

}

func main() {
	//1 创建一个server句柄,使用Zinx的api
	s := zserver.NewServer("[nightzinx V0.3]")

	//2 给当前zinx框架添加一个自定义的router
	s.AddRouter(&PingRouter{})

	//3 启动server
	s.Serve()
}

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

//Test Handle
func (this *PingRouter) Handle(request ziserver.IRequest) {
	fmt.Println("Call Ping Router Handle...")

	//先读取客户端的数据,在回写ping..ping..ping

	fmt.Println("recv from client: msgID =", request.GetMsgID(), ",data =", string(request.GetData()))

	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

//hello ZinxRouter test 自定义路由
type HelloZinxRouter struct {
	zserver.BaseRouter
}

func (this *HelloZinxRouter) Handle(request ziserver.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle...")

	//先读取客户端的数据,在回写ping..ping..ping

	fmt.Println("recv from client: msgID =", request.GetMsgID(), ",data =", string(request.GetData()))

	err := request.GetConnection().SendMsg(201, []byte("Hello Welcome to Zinx!!"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	//1 创建一个server句柄,使用Zinx的api
	s := zserver.NewServer("[nightzinx V0.6]")

	//2 给当前zinx框架添加一个自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})

	//3 启动server
	s.Serve()
}

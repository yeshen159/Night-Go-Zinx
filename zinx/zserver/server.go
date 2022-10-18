package zserver

import (
	"Night/zinx/utils"
	"Night/zinx/ziserver"
	"fmt"
	"net"
)

//server.go的接口实现,定义一个Server的服务器模块
type Server struct {
	Name       string                //服务器的名称
	IPVersion  string                //服务器绑定的ip版本
	Ip         string                //服务器监听的ip
	Port       int                   //服务器监听的端口
	MsgHandler ziserver.IMsgHandle   //当前server的消息管理模块, 用来绑定MsgID和对应的处理业务API关系
	ConnMgr    ziserver.IConnManager //当前server的链接管理器
	//Router    ziserver.IRouter   //当前的Server添加一个router，server注册的链接对应的处理业务

	//该Server创建链接之后自动调用Hook函数--OnConnStart
	OnConnStart func(conn ziserver.IConnection)
	//该Server创建链接之后自动调用Hook函数--OnConnStop
	OnConnStop func(conn ziserver.IConnection)
}

//定义当前客户端链接的所绑定handle api(目前这个handle是写死的，以后优化应该由用户自定义handle方法)
//func CallBackToClient(conn *net.TCPConn, date []byte, cnt int) error {
//	//回显业务
//	fmt.Println("[Conn Handle] CallBackToClient...")
//
//	if _, err := conn.Write(date[:cnt]); err != nil {
//
//		fmt.Println("write back buf err", err)
//
//		return errors.New("CallBackToClient error")
//	}
//
//	return nil
//}

//初始化Server模块的方法
func NewServer(name string) ziserver.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		Ip:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
		//Router:    nil,
	}

	return s
}

//路由功能: 给当前的服务注册一个路由方法，供客户端的链接处理使用
func (s *Server) AddRouter(msgID uint32, router ziserver.IRouter) {
	//s.Router = router

	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Success!!")
}

func (s *Server) GetConnMgr() ziserver.IConnManager {
	return s.ConnMgr
}

//启动服务器
func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name %s Listenner at IP: %s, Port: %d, is starting...\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d, MaxPackageSize: %d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)
	//fmt.Printf("[start] Server Listenner at IP: %s, Port %d, is starting\n", s.Ip, s.Port)

	go func() {
		//0 开启消息队列及Worker工作池
		s.MsgHandler.StartWorkerPool()

		//1 获取一个TCP的addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.Ip, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addt error:", err)
			return
		}
		//2 监听服务器的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}

		fmt.Println("start Zinx server succ,", s.Name, "succ,Listenning...")

		var cid uint32
		cid = 0

		//3 阻塞的等待客户端连接,处理客户端链接业务(读写)
		for {
			//如果有客户端链接过来,阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accpet err", err)
				continue
			}

			//设置最大链接个数的判断, 如果超过最大链接数量, 那么关闭此新的链接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				//TODO 给客户端响应一个超出最大链接的错误包
				fmt.Println("---------------->Too Many Connections MaxConn = ", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			//将处理新链接的业务方法 和 conn 进行绑定 得到我们的链接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			//dealConn := NewConnection(conn, cid, s.)

			cid++

			//启动当前的链接业务处理

			go dealConn.Start()
			//已经与客户端建立链接,做一些业务
			//go func() {
			//	for {
			//		buf := make([]byte, 512)
			//		cnt, err := conn.Read(buf)
			//		if err != nil {
			//			fmt.Println("recv buf err", err)
			//			continue
			//		}

			//		fmt.Printf("recv client buf %s, cnt %d\n", buf, cnt)

			//		//回显功能
			//		if _, err := conn.Write(buf[:cnt]); err != nil {
			//			fmt.Println("write back buf err", err)
			//			continue
			//		}
			//		fmt.Println(buf[:cnt])
			//	}
			//}()
		}
	}()
}

//停止服务器
func (s *Server) Stop() {
	//将一些服务器的资源.状态或者一些已经开辟的链接信息 进行停止或者回收
	fmt.Println("[STOP] Zinx server name", s.Name)
	s.ConnMgr.ClearConn()
}

//运行服务器
func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()

	select {}
}

//注册OnConnStart 钩子函数的方法
func (s *Server) SetOnConnStart(hookFunc func(connection ziserver.IConnection)) {
	s.OnConnStart = hookFunc
}

//注册OnConnStop 钩子函数的方法
func (s *Server) SetOnConnStop(hookFunc func(connection ziserver.IConnection)) {
	s.OnConnStop = hookFunc
}

//调用OnConnStart 钩子函数的方法
func (s *Server) CallOnConnStart(conn ziserver.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("-->call OnConnStart()...")
		s.OnConnStart(conn)
	}
}

//调用OnConnStop 钩子函数的方法
func (s *Server) CallOnConnStop(conn ziserver.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("-->call OnConnStop()...")
		s.OnConnStop(conn)
	}

}

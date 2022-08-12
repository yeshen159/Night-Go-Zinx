package zserver

import (
	"Night/zinx/ziserver"
	"errors"
	"fmt"
	"net"
)

//server.go的接口实现,定义一个Server的服务器模块
type Server struct {
	Name      string //服务器的名称
	IPVersion string //服务器绑定的ip版本
	Ip        string //服务器监听的ip
	Port      int    //服务器监听的端口
}

//定义当前客户端链接的所绑定handle api(目前这个handle是写死的，以后优化应该由用户自定义handle方法)
func CallBackToClient(conn *net.TCPConn, date []byte, cnt int) error {
	//回显业务
	fmt.Println("[Conn Handle] CallBackToClient...")

	if _, err := conn.Write(date[:cnt]); err != nil {

		fmt.Println("write back buf err", err)

		return errors.New("CallBackToClient error")
	}

	return nil
}

//初始化Server模块的方法
func NewServer(name string) ziserver.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		Ip:        "0.0.0.0",
		Port:      8999,
	}

	return s
}

//启动服务器
func (s *Server) Start() {
	fmt.Printf("[start] Server Listenner at IP: %s, Port %d, is starting\n", s.Ip, s.Port)

	go func() {
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
			}

			//将处理新链接的业务方法 和 conn 进行板顶 得到我们的链接模块
			dealConn := NewConnection(conn, cid, CallBackToClient)

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

}

//运行服务器
func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()

	select {}
}

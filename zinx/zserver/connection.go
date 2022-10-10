package zserver

import (
	"Night/zinx/ziserver"
	"errors"
	"fmt"
	"io"
	"net"
)

type Connection struct {
	//当前链接的socket TCP套接字
	Conn *net.TCPConn

	//链接的ID
	ConnID uint32

	//当前的连接状态
	isClosed bool

	//当前链接锁绑定的处理业务方法的API
	//handleAPI ziserver.HandleFunc

	//告知当前链接已经退出的/停止的 channel
	ExitChan chan bool

	//该链接处理的方法Router
	Router ziserver.IRouter
}

func NewConnection(conn *net.TCPConn, connID uint32, router ziserver.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
		//handleAPI: callback_pai,
	}

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")

	defer fmt.Println("connID =", c.ConnID, "Reader is exit, remote add is", c.RemoteAddr().String)
	defer c.Stop()

	for {
		//读取客户端的数据到buf中
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("recv buf err", err)
		//	continue
		//}

		//创建一个拆包解包对象
		dp := NewDataPack()

		//读取客户端的Msg Head二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error", err)
			break
		}
		//拆包，得到msgID 和 msgDataLen 放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}

		//根据datalen 再次读取Data 放在msg.Data中
		var data []byte

		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data err", err)
				break
			}
		}

		msg.SetData(data)

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		//执行注册的路由方法
		go func(request ziserver.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

		//调用当前链接所绑定的HandleAPI
		//if c.handleAPI(c.Conn, buf, cnt); err != nil {
		//	fmt.Println("ConnID", c.ConnID, "handle is error", err)
		//	break
		//}
	}
}

//启动链接 让当前的链接准备开始工作
func (c *Connection) Start() {

	fmt.Println("Conn Start().. ConnID = ", c.ConnID)
	//启动从当前链接的读数据的业务
	go c.StartReader()
}

//停止链接 结束当前链接的工作
func (c *Connection) Stop() {

	fmt.Println("Conn Stop().. ConnID = ", c.ConnID)

	//如果当前链接已经关闭
	if c.isClosed == true {
		return
	}

	c.isClosed = true

	//关闭socket链接
	c.Conn.Close()

	//回收资源
	close(c.ExitChan)
}

//获取当前链接的绑定socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {

	return c.Conn
}

func (c *Connection) GetConnID() uint32 {

	return c.ConnID
}

//获取远程客户端的TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr {

	return c.Conn.RemoteAddr()
}

//提供一个SendMsg方法 将我们要发送给客户端的数据，先进行封包，在发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}

	//将data进行封包 MsgDataLen|MsgID|Data
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg")
	}

	//将数据发送给客户端
	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("Write msg id", msgId, "error :", err)
		return errors.New("conn Wrire error")
	}

	return nil
}

package zserver

import (
	//"Night/zinx/utils"
	//"Night/zinx/ziserver"

	"github.com/yeshen159/Night-Go-Zinx/zinx/utils"
	"github.com/yeshen159/Night-Go-Zinx/zinx/ziserver"

	//"github.com/yeshen159/Night-Go-Zinx/zinx/ziserver"
	//"github.com/yeshen159/Night-Go-Zinx/zinx/utils"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

type Connection struct {
	//当前Conn隶属于哪个Server
	TcpServer ziserver.IServer

	//当前链接的socket TCP套接字
	Conn *net.TCPConn

	//链接的ID
	ConnID uint32

	//当前的连接状态
	isClosed bool

	//当前链接锁绑定的处理业务方法的API
	//handleAPI ziserver.HandleFunc

	//告知当前链接已经退出的/停止的 channel(由Redear告知Writer)
	ExitChan chan bool

	//无缓冲的管道,用于读.写Goroutine之间的消息通信
	msgChan chan []byte

	//该链接处理的方法Router
	//Router ziserver.IRouter

	//消息的管理MsgID 和对应的处理业务API关系
	MsgHandler ziserver.IMsgHandle

	//链接属性合集
	property map[string]interface{}

	//保护链接属性的锁
	propertyLock sync.RWMutex
}

func NewConnection(server ziserver.IServer, conn *net.TCPConn, connID uint32, msgHandler ziserver.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
		property:   make(map[string]interface{}),
		//Router:   router,
		//handleAPI: callback_pai,
	}

	//将conn加入到ConnMananger中
	c.TcpServer.GetConnMgr().Add(c)

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]")
	//fmt.Println("Reader Goroutine is running...")

	defer fmt.Println("[Reader is exit!], connID =", c.ConnID, "remote add is", c.RemoteAddr().String())
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

		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经开启了工作池机制, 将消息发送给Worker工作池处理即可
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			//从路由中,找到注册绑定的Conn对应的router调用
			//根据绑定好的MsgID 找到对应处理api业务 执行
			go c.MsgHandler.DoMsgHandler(&req)
		}

		//执行注册的路由方法
		//go func(request ziserver.IRequest) {
		//	c.Router.PreHandle(request)
		//	c.Router.Handle(request)
		//	c.Router.PostHandle(request)
		//}(&req)

		//调用当前链接所绑定的HandleAPI
		//if c.handleAPI(c.Conn, buf, cnt); err != nil {
		//	fmt.Println("ConnID", c.ConnID, "handle is error", err)
		//	break
		//}
	}
}

/*
	写消息Goroutine，专门发送给客户端消息的模块
*/
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println("[conn Writer exit!]", c.RemoteAddr().String())

	//不断地阻塞的等待channel的消息,进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			//有数据要写
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error,", err)
				return
			}
		case <-c.ExitChan:
			//代表Reader已经退出,此时Writer也要退出
			return
		}
	}

}

//启动链接 让当前的链接准备开始工作
func (c *Connection) Start() {

	fmt.Println("Conn Start().. ConnID = ", c.ConnID)
	//启动从当前链接的读数据的业务
	go c.StartReader()
	//启动从当前链接的写数据的业务
	go c.StartWriter()

	//按照开发者传递进来的 创建链接之后需要调用的处理业务, 执行对应的Hook函数
	c.TcpServer.CallOnConnStart(c)
}

//停止链接 结束当前链接的工作
func (c *Connection) Stop() {

	fmt.Println("Conn Stop().. ConnID = ", c.ConnID)

	//如果当前链接已经关闭
	if c.isClosed == true {
		return
	}

	c.isClosed = true

	//调用开发者注册的 销毁链接之前 需要执行的业务Hook函数
	c.TcpServer.CallOnConnStop(c)

	//关闭socket链接
	c.Conn.Close()

	//告知Writer关闭
	c.ExitChan <- true

	//将当前链接从ConnMgr中摘除掉
	c.TcpServer.GetConnMgr().Remove(c)

	//回收资源
	close(c.ExitChan)
	close(c.msgChan)
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
	c.msgChan <- binaryMsg
	//if _, err := c.Conn.Write(binaryMsg); err != nil {
	//	fmt.Println("Write msg id", msgId, "error :", err)
	//	return errors.New("conn Wrire error")
	//}

	return nil
}

//设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	//添加一个链接属性
	c.property[key] = value
}

//获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	//读取属性
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

//移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	//删除属性
	delete(c.property, key)
}

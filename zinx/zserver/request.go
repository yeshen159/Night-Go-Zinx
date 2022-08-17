package zserver

import "Night/zinx/ziserver"

type Request struct {
	//已经和客户端建立好的链接
	conn ziserver.IConnection

	//客户端请求的
	data []byte
}

//得到当前链接
func (r *Request) GetConnection() ziserver.IConnection {

	return r.conn
}

//得到请求的消息数据
func (r *Request) GetData() []byte {

	return r.data
}

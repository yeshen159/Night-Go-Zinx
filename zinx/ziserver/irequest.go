package ziserver

/*
	IRequest接口:
	实际上是把客户端请求的链接信息, 和请求的数据 包装到了一个Request请求中
*/

type IRequest interface {
	//得到当前链接
	GetConnection() IConnection

	//得到得到得消息数据
	GetData() []byte
}

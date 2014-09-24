package mqtt

type OnConnectHandler func()
type OnDisconnectHandler func()

type Callbacks struct {
	OnConnect    OnConnectHandler
	OnDisconnect OnDisconnectHandler
}

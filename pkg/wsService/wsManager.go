package wsService

import (
	"blog/pkg/l"
	"go.uber.org/zap"
	"sync"
)

var (
	gWsClientManager *wsClientManager
	once             sync.Once
)

type wsClientManager struct {
	clients    map[*WsClient]bool
	broadCast  chan *wsMsgWithDst
	register   chan *WsClient
	unregister chan *WsClient
}

func GetWsCliManager() *wsClientManager {
	return gWsClientManager
}

func startWebSocket(bfDataChan chan []byte) {
	once.Do(func() {
		gWsClientManager = &wsClientManager{
			clients:    make(map[*WsClient]bool),
			broadCast:  make(chan *wsMsgWithDst),
			register:   make(chan *WsClient),
			unregister: make(chan *WsClient),
		}
		go gWsClientManager.getBFQueryStr(bfDataChan)
		go gWsClientManager.start()
	})
}

func (c *wsClientManager) AddClient(client *WsClient) {
	l.GetLogger().Info("AddClient", zap.Any("client", client.Socket.RemoteAddr().String()))
	c.register <- client

	go client.read()
	go client.write()
}

func (c *wsClientManager) getBFQueryStr(bfDataCh chan []byte) {
	for msg := range bfDataCh {
		msgPkt := wsMsgWithDst{
			Message: msg,
		}
		c.broadCast <- &msgPkt
	}
}

func (c *wsClientManager) start() {
	for {
		select {
		case client := <-c.register:
			l.GetLogger().Info("wsClientManager receive one register")
			c.clients[client] = true
			l.GetLogger().Info("wsClientManager total size", zap.Int("length", len(c.clients)))
		case client := <-c.unregister:
			l.GetLogger().Info("wsClientManager receive one unregister")
			if _, ok := c.clients[client]; ok {
				close(client.SendCh)
				delete(c.clients, client)
			}
		case msgPkt := <-c.broadCast:
			l.GetLogger().Info("wsClientManager receive a message")
			for conn := range c.clients {
				//peerId为空时，广播给所有客户端；否则单播给特定客户端
				if msgPkt.PeerId != "" && msgPkt.PeerId != conn.PeerId {
					continue
				}
				select {
				case conn.SendCh <- msgPkt.Message:
				default:
					close(conn.SendCh)
					delete(c.clients, conn)
				}
			}
		}
	}
}

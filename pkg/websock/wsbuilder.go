package websock

import (
	"blog/pkg/l"
	"github.com/gorilla/websocket"
)

type wsMsgWithDst struct {
	PeerId  string
	Message []byte
}

type WsClient struct {
	PeerId string //对端标记，用于区分不同的Web在线用户
	Socket *websocket.Conn
	SendCh chan []byte
}

func GetRegisterFunc() func() chan []byte {
	bfDataChan := make(chan []byte)
	startWebSocket(bfDataChan)
	return func() chan []byte {
		return bfDataChan
	}
}

func (c *WsClient) read() {
	defer func() {
		gWsClientManager.unregister <- c
		_ = c.Socket.Close()
	}()

	for {
		_, _, err := c.Socket.ReadMessage()
		if err != nil {
			gWsClientManager.unregister <- c
			_ = c.Socket.Close()
			break
		}
	}
}

func (c *WsClient) write() {
	defer c.Socket.Close()
	for {
		select {
		case message, ok := <-c.SendCh:
			if !ok {
				l.Logger().Warn("sendChan is closed")
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			l.Logger().Debug("send a message to websocket")
			_ = c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

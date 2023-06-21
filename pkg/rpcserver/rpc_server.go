package rpcserver

import (
	"blog/pkg/internal/rpcpackage"
	"blog/pkg/l"
	"go.uber.org/zap"
	"net"
	"sync"
)

type UserChecker interface {
	Identify(string) (int, error)
	GetToken(string) (string, error)
}

func listenAndServe(bindAddr string, checker UserChecker, rpcSessionStatusReceiveCh chan<- *SessionStatus, listenStateCh chan<- error) {
	l.Logger().Info("RpcServer ListenAndServe:", zap.String("listen", bindAddr))
	socket, err := net.Listen("tcp", bindAddr)
	if err != nil {
		l.Logger().Error("RpcServer ListenAndServe:", zap.String("error", err.Error()))
		listenStateCh <- err
		return
	}
	listenStateCh <- nil
	listenStateCh = nil

	defer socket.Close()

	var connId = 1

	for {
		conn, err := socket.Accept()
		if err != nil {
			continue
		}

		l.Logger().Info("RpcServer accept one connection", zap.Int("connId", connId))

		go connectionLoop(conn, connId, checker, rpcSessionStatusReceiveCh)

		connId++
	}
}

func connectionLoop(conn net.Conn, connId int, checker UserChecker, rpcSessionStatusReceiveCh chan<- *SessionStatus) {
	var once sync.Once
	doClean := func() {
		once.Do(func() {
			l.Logger().Info("RPCServer conn Close", zap.Int("connId", connId))
			conn.Close()
		})
	}
	defer doClean()

	session := newRPCSession(checker, rpcSessionStatusReceiveCh)
	defer session.Close()
	sendChan := session.Open()

	// send loop
	go func() {
		for {
			msgs, ok := <-sendChan
			if !ok {
				l.Logger().Info("send loop sendChan closed!")
				doClean()
				break
			}

			for _, msg := range msgs {
				_, err := writeToConn(conn, msg)
				if err != nil {
					l.Logger().Info("send loop writeToConn error!")
					break
				}
			}
		}
		l.Logger().Info("send loop exit!", zap.Int("connId", connId))
	}()

	// receive loop
	builder := rpcpackage.CreateRPCMsgBuilder(conn)
	for {
		payload, err := builder.DoRead()
		if err != nil {
			l.Logger().Warn("receive loop read", zap.Error(err), zap.Int("connId", connId))
			break
		}
		//-----------如果binary部分是空，就填入ip地址，复用--------------
		if len(payload.BinaryRequest) == 0 {
			payload.BinaryRequest = []byte(conn.RemoteAddr().String())
		}

		//-----------如果binary部分是空，就填入ip地址，复用--------------

		err = session.OnMessage(payload.JsonRequest, payload.BinaryRequest)
		if err != nil {
			l.Logger().Warn("receive loop OnMessage", zap.Error(err), zap.Int("connId", connId))
			break
		}
	}
	l.Logger().Info("receive loop exit!", zap.Int("connId", connId))
}

func writeToConn(conn net.Conn, data []byte) (int, error) {
	dataLen := len(data)
	nWrites := 0
	for {
		n, err := conn.Write(data[nWrites:])
		if err != nil {
			return n, err
		}

		nWrites += n
		if nWrites >= dataLen {
			break
		}
	}

	return nWrites, nil
}

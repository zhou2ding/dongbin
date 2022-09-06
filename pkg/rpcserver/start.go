package rpcserver

import "blog/pkg/v"

func StartRPC(addr string, registerFn func() (checker UserChecker, registers []func(), sessionCh chan<- *SessionStatus)) error {
	ck, reg, ch := registerFn()
	for _, r := range reg {
		r()
	}

	alive := v.GetViper().GetInt64("rpc_max_alive")
	if alive == 0 {
		alive = 30
	}
	GetSessionMgr().Start(alive)

	listen := make(chan error, 1)
	go listenAndServe(addr, ck, ch, listen)
	return <-listen
}

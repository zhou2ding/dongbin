package gb28181

import (
	"context"
	"github.com/ghettovoice/gosip"
	"github.com/ghettovoice/gosip/sip"
	"github.com/husanpao/ip"
	"strings"
)

var (
	gServer gosip.Server
)

func StartServer(ctx context.Context) error {
	// todo NewServer的logger后续增加，补充注册和MESSAGE函数
	initRoutes()
	cfg := gosip.ServerConfig{}
	gServer = gosip.NewServer(cfg, nil, nil, nil)
	err := gServer.OnRequest(sip.REGISTER, onRegister)
	if err != nil {
		return err
	}
	err = gServer.OnRequest(sip.MESSAGE, onMessage)
	if err != nil {
		return err
	}
	err = gServer.OnRequest(sip.BYE, onBye)
	if err != nil {
		return err
	}
	err = gServer.Listen("", "")
	if err != nil {
		return err
	}
	return nil
}

func initRoutes() map[string]string {
	routes := make(map[string]string)
	for k, v := range myip.LocalAndInternalIPs() {
		routes[k] = v
		if dot := strings.LastIndex(k, "."); dot >= 0 {
			routes[k[0:dot]] = k
		}
	}
	return routes
}

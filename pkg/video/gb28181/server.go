package gb28181

import (
	"context"
	"github.com/ghettovoice/gosip"
	"github.com/ghettovoice/gosip/sip"
	"github.com/husanpao/ip"
	"strings"
)

func StartServer(ctx context.Context) error {
	// todo NewServer的logger后续增加，补充注册和MESSAGE函数
	initRoutes()
	cfg := gosip.ServerConfig{}
	srv := gosip.NewServer(cfg, nil, nil, nil)
	err := srv.OnRequest(sip.REGISTER, nil)
	if err != nil {
		return err
	}
	err = srv.OnRequest(sip.MESSAGE, nil)
	if err != nil {
		return err
	}
	err = srv.OnRequest(sip.BYE, onBye)
	if err != nil {
		return err
	}
	err = srv.Listen("", "")
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

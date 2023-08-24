package gb28181

import (
	"context"
	"github.com/ghettovoice/gosip"
)

func StartServer(ctx context.Context) {
	cfg := gosip.ServerConfig{}
	gosip.NewServer(cfg, nil, nil, nil) // todo logger后续增加
}

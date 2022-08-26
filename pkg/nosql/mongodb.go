package nosql

import (
	"blog/pkg/l"
	"blog/pkg/v"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"time"
)

func initMongo() error {
	user := v.GetViper().GetString("nosql.user")
	pwd := v.GetViper().GetString("nosql.pwd")
	host := v.GetViper().GetString("nosql.host")
	port := v.GetViper().GetString("nosql.port")

	url := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?retryWrites=fales", user, pwd, host, port, "admin")
	opts := options.Client().ApplyURI(url)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cli, err := mongo.Connect(ctx, opts)
	if err != nil {
		l.GetLogger().Error("connect to mongodb failed", zap.Error(err))
		return err
	}

	err = cli.Ping(ctx, readpref.Primary())
	if err != nil {
		l.GetLogger().Error("ping mongodb failed", zap.Error(err))
		return err
	}

	return nil
}

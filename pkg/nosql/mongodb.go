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
	"sync"
	"time"
)

var (
	once   sync.Once
	inst   *NoSQL
	dbName string
	cli    *mongo.Client
)

type NoSQL struct {
}

func GetNoSQL() *NoSQL {
	once.Do(func() {
		inst = &NoSQL{}
	})
	return inst
}

func (n *NoSQL) GetDB() *mongo.Database {
	return cli.Database(dbName)
}

func initMongo() error {
	user := v.GetViper().GetString("nosql.user")
	pwd := v.GetViper().GetString("nosql.pwd")
	host := v.GetViper().GetString("nosql.host")
	port := v.GetViper().GetString("nosql.port")

	//url := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?retryWrites=fales", user, pwd, host, port, "admin")
	url := fmt.Sprintf("mongodb://%s:%s", host, port)
	opt := &options.ClientOptions{}
	opt.SetAuth(options.Credential{AuthMechanism: "SCRAM-SHA-1", AuthSource: "admin", Username: user, Password: pwd}) // 不把用户密码写入连接url
	opt.SetMaxPoolSize(5)                                                                                             // 连接池大小

	opts := options.Client().ApplyURI(url)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, err := mongo.Connect(ctx, opts)
	if err != nil {
		l.GetLogger().Error("connect to mongodb failed", zap.Error(err))
		return err
	}

	err = c.Ping(ctx, readpref.Primary())
	if err != nil {
		l.GetLogger().Error("ping mongodb failed", zap.Error(err))
		return err
	}

	dbName = v.GetViper().GetString("nosql.database_name")
	l.GetLogger().Info("connect to mongodb success", zap.String("host", url))
	return nil
}

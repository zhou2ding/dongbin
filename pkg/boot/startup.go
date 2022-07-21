package boot

import (
	"blog/pkg/v"
	"blog/pkg/l"
	"blog/pkg/version"
	"flag"
	"fmt"
	"go.uber.org/zap"
)

func StartUp(appName string) {
	var cfgFile string
	flag.StringVar(&cfgFile, "config", "", "config file(.toml)")
	flag.StringVar(&cfgFile, "c", "", "config file(.toml)")
	flag.Parse()

	if cfgFile == "" {
		cfgFile = "app.toml"
	}
	fmt.Println("configFile path is " + cfgFile)

	v.InitViper()
	if err := v.LoadConfig(cfgFile); err != nil {
		panic(err)
	}

	if err := l.InitLogger(appName); err != nil {
		panic(err)
	}

	l.GetLogger().Info("start service", zap.String("service", appName), zap.Any("version", version.GetVersionInfo()))
}

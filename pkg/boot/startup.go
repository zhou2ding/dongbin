package boot

import (
	"blog/pkg/cfg"
	"blog/pkg/logger"
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
		cfgFile = "../config/app.toml"
	}
	fmt.Println("configFile path is " + cfgFile)

	cfg.InitViper()
	if err := cfg.LoadConfig(cfgFile); err != nil {
		panic(err)
	}

	if err := logger.InitLogger(appName); err != nil {
		panic(err)
	}

	logger.GetLogger().Info("start service", zap.String("service", appName), zap.Any("version", version.GetVersionInfo()))
}

package util

import (
	"flag"
	"fmt"
)

func StartUp() {
	var cfgFile string
	flag.StringVar(&cfgFile, "config", "", "配置文件(.toml)")
	flag.StringVar(&cfgFile, "c", "", "配置文件(.toml)")

	flag.Parse()
	if cfgFile == "" {
		cfgFile = "../config/app.toml"
	}

	fmt.Println("configFile path is " + cfgFile)
}

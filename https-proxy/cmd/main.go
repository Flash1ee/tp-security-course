package main

import (
	"flag"
	"log"

	"http-proxy/cfg"
	"http-proxy/internal/proxy/server"
	"http-proxy/pkg/utils"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "cfg/config.toml", "path to config file")
}
func main() {
	flag.Parse()
	config := cfg.NewConfig()
	_, err := toml.DecodeFile(configPath, &config)
	if err != nil {
		logrus.Fatal(err)
	}
	srv := proxy.New(utils.GetLogger(config), config.BindAddr)
	if err := srv.Start(config); err != nil {
		log.Fatal("srv.Start err:", err)
	}
}

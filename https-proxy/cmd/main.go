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
	conn := utils.NewPostgresConn(config.DatabaseURL)
	if conn == nil {
		logrus.Fatal("database conn is nil")
	}

	srv := proxy.New(utils.GetLogger(config), conn, config.BindAddr)
	if err := srv.Start(config); err != nil {
		log.Fatal("srv.Start err:", err)
	}
}

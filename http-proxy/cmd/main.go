package main

import (
	"flag"
	"log"
	"net/http"

	"http-proxy/cfg"
	"http-proxy/pkg/proxy"
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
	srv := proxy.New(utils.GetLogger(config))
	if err := http.ListenAndServe(config.BindAddr, srv); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

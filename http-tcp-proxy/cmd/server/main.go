package main

import (
	"flag"
	"log"

	"http-tcp-proxy/cfg"
	"http-tcp-proxy/pkg/proxy"

	"github.com/BurntSushi/toml"
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
		log.Fatal(err)
	}
	proxy.Run(config)
}

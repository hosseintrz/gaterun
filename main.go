package main

import (
	"flag"
	"log"

	"github.com/hosseintrz/gaterun/config"
	"github.com/hosseintrz/gaterun/proxy"
	"github.com/hosseintrz/gaterun/router/gorilla"
)

func main() {
	port := flag.Int("p", 7800, "Port of the service")
	configFile := flag.String("c", "./config.json", "Path to the configuration filename")
	flag.Parse()

	serviceConfig, err := config.Parse(*configFile)
	if err != nil {
		log.Fatal("ERROR:", err.Error())
	}

	if *port != 0 {
		serviceConfig.Port = *port
	}

	routerFactory := gorilla.DefaultFactory(proxy.NewDefaultFactory())
	routerFactory.New().Run(serviceConfig)
}

package main

import (
	"backend_gen/config"
	"backend_gen/internal/server"
	"flag"
	"log"
)

func main() {
	cfgPath := flag.String("c", "config/config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.ReadConfig(*cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	s, err := server.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	s.Run()
}

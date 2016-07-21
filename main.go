package main

import (
	"os"

	"github.com/NYTimes/gizmo/config"
	"github.com/NYTimes/gizmo/server"
	"github.com/NYTimes/mock-ec2-metadata/service"
)

func main() {
	var cfg *service.Config

	if _, err := os.Stat("./mock-ec2-metadata-config.json"); err == nil {
		config.LoadJSONFile("./mock-ec2-metadata-config.json", &cfg)
	} else if _, err := os.Stat("/etc/mock-ec2-metadata-config.json"); err == nil {
		config.LoadJSONFile("/etc/mock-ec2-metadata-config.json", &cfg)
	} else {
		server.Log.Fatal("unable to locate config file")
	}

	server.Init("mock-ec2-metadata", cfg.Server)
	err := server.Register(service.NewMetadataService(cfg))
	if err != nil {
		server.Log.Fatal("unable to register service: ", err)
	}

	err = server.Run()
	if err != nil {
		server.Log.Fatal("server encountered a fatal error: ", err)
	}
}

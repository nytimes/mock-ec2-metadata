package main

import (
    "github.com/NYTimes/mock-ec2-metadata/service"
    "github.com/NYTimes/gizmo/config"
    "github.com/NYTimes/gizmo/server"
)

func main() {
    // showing 1 way of managing gizmo/config: importing from a local file
    var cfg *service.Config
    config.LoadJSONFile("./config.json", &cfg)

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
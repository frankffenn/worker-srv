package main

import (
	"log"
	"time"
	"os"

	pb "github.com/frankffenn/worker-srv/proto"
	micro "github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	srv "github.com/micro/go-micro/v2/registry/service"
)

func main() {
	addr := os.Getenv("REGISTRY_ADDR")
	reg := srv.NewRegistry(registry.Addrs(addr))
	service := micro.NewService(
		micro.Name("worker.srv"),
		micro.RegisterInterval(10*time.Second),
		micro.RegisterTTL(30*time.Second),
		micro.Registry(reg),
	)

	service.Init()

	pb.RegisterWorkerServiceHandler(service.Server(), NewServer())

	// Run the server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

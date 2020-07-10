package main

import (
	"log"

	pb "github.com/frankffenn/worker-srv/registry/service/proto"
	micro "github.com/micro/go-micro/v2"
)

func main() {
	service := micro.NewService(
		micro.Name("registry.center"),
		micro.Address(":8000"),
	)

	service.Init()

	pb.RegisterRegistryHandler(service.Server(), NewServer())

	// Run the server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"
	"time"

	pb "github.com/frankffenn/worker-srv/proto"
	micro "github.com/micro/go-micro/v2"
)

func main() {
	service := micro.NewService(
		micro.Name("worker.srv"),
		micro.RegisterInterval(10*time.Second),
		micro.RegisterTTL(30*time.Second),
	)

	service.Init()

	pb.RegisterWorkerServiceHandler(service.Server(), NewServer())

	// Run the server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"fmt"
	"time"

	pb "github.com/frankffenn/worker-srv/proto"
	"github.com/frankffenn/worker-srv/src"
	"github.com/micro/go-micro"
)

func main() {
	service := micro.NewService(
		micro.Name("worker"),
		micro.RegisterInterval(5*time.Second),
		micro.RegisterTTL(30*time.Second),
	)

	service.Init()

	srv := src.NewServer()
	pb.RegisterWorkerServiceHandler(service.Server(), srv)

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

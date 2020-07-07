package main

import (
	"log"
	"time"

	pb "github.com/frankffenn/worker-srv/proto"
	"github.com/micro/cli/v2"
	micro "github.com/micro/go-micro/v2"
)

func main() {
	service := micro.NewService(
		micro.Name("worker.srv"),
		micro.RegisterInterval(10*time.Second),
		micro.RegisterTTL(30*time.Second),
		micro.Flags(
			&cli.Uint64Flag{
				Name:  "max",
				Usage: "the maximum number of concurrent tasks",
			},
		),
	)

	var max uint64 = 1

	service.Init(
		micro.Action(func(c *cli.Context) error {
			if c.Uint64("max") > 0 {
				max = c.Uint64("max")
			}
			return nil
		}),
	)

	pb.RegisterWorkerServiceHandler(service.Server(), NewServer(max))

	// Run the server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

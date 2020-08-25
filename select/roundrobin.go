package roundrobin

import (
	"log"
	"os"
	"sync/atomic"

	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"

	"context"

	rpb "github.com/frankffenn/worker-srv/registry/service/proto"
)

var defaultService = "go.micro.srv.registry"

type roundrobin struct {
	rr map[string]int
	client.Client
	registry rpb.RegistryService

	busy int32
}

func (s *roundrobin) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	var node *registry.Node
	nOpts := append(opts, client.WithSelectOption(
		// create a selector strategy
		selector.WithStrategy(func(services []*registry.Service) selector.Next {
			if atomic.CompareAndSwapInt32(&s.busy, 0, 1) {
				defer atomic.StoreInt32(&s.busy, 0)

				idles := make(map[string]uint64, 0)
				rsp, err := s.registry.GetServices(ctx, &rpb.GetServicesRequest{}, callOpts()...)
				if err != nil {
					log.Println("get all services failed", err)
				} else {
					idles = rsp.Services
				}

				// flatten
				var nodes []*registry.Node
				for _, service := range services {
					for _, node := range service.Nodes {
						if val := idles[node.Id]; val <= 0 {
							continue
						}
						nodes = append(nodes, node)
					}
				}

				if len(nodes) > 0 {
					rr := s.rr[req.Service()]
					node = nodes[rr%len(nodes)]
					rr++
					s.rr[req.Service()] = rr

					_, err = s.registry.Mark(ctx, &rpb.MarkRequest{Id: node.Id}, callOpts()...)
					if err != nil {
						log.Println("mark service failed", err)
					}
				}

			}

			// create the next func that always returns our node
			return func() (*registry.Node, error) {
				if node == nil {
					return nil, selector.ErrNoneAvailable
				}
				return node, nil
			}
		}),
	))

	defer func() {
		if node == nil {
			return
		}
		log.Println("reset service", node.Id)
		if _, err := s.registry.Reset(ctx, &rpb.ResetRequest{Id: node.Id}, callOpts()...); err != nil {
			log.Println("reset request error", err)
		}
	}()

	if err := s.Client.Call(ctx, req, rsp, nOpts...); err != nil {
		return err
	}

	return nil
}

// NewClientWrapper is a wrapper which roundrobins requests
func NewClientWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		registry := rpb.NewRegistryService(defaultService, grpc.NewClient())
		return &roundrobin{
			rr:       make(map[string]int),
			Client:   c,
			registry: registry,
		}
	}
}

func callOpts() []client.CallOption {
	var opts []client.CallOption
	addr := os.Getenv("REGISTRY_ADDR")
	if addr != "" {
		opts = append(opts, client.WithAddress(addr))
	}
	return opts
}

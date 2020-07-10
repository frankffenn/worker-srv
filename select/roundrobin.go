package roundrobin

import (
	"log"
	"os"
	"sync"

	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"

	"context"

	rpb "github.com/frankffenn/worker-srv/registry/service/proto"
)

var DefaultName = "registry.center"

type roundrobin struct {
	sync.Mutex
	rr map[string]int
	client.Client
	registry rpb.RegistryService
	done     chan struct{}
}

func (s *roundrobin) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	var nodeId string
	nOpts := append(opts, client.WithSelectOption(
		// create a selector strategy
		selector.WithStrategy(func(services []*registry.Service) selector.Next {
			idles := make(map[string]uint64, 0)
			rsp, err := s.registry.GetServices(ctx, &rpb.GetServicesRequest{}, opt)
			if err != nil {
				log.Println("GetService error", err)
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
				// nodes = append(nodes, service.Nodes...)
			}

			// create the next func that always returns our node
			return func() (*registry.Node, error) {
				if len(nodes) == 0 {
					return nil, selector.ErrNoneAvailable
				}
				s.Lock()
				// get counter
				rr := s.rr[req.Service()]
				// get node
				node := nodes[rr%len(nodes)]
				// increment
				rr++
				// save
				s.rr[req.Service()] = rr
				s.Unlock()

				_, err := s.registry.Mark(ctx, &rpb.MarkRequest{Id: node.Id}, opt)
				if err != nil {
					log.Println("mark request error", err)
				}

				nodeId = node.Id
				return node, nil
			}
		}),
	))
	if err := s.Client.Call(ctx, req, rsp, nOpts...); err != nil {
		return err
	}

	if _, err := s.registry.Reset(ctx, &rpb.ResetRequest{Id: nodeId}, opt); err != nil {
		log.Println("reset request error", err)
	}

	return nil
}

// NewClientWrapper is a wrapper which roundrobins requests
func NewClientWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		registry := rpb.NewRegistryService(DefaultName, grpc.NewClient())
		return &roundrobin{
			rr:       make(map[string]int),
			Client:   c,
			registry: registry,
		}
	}
}

var opt = func(opt *client.CallOptions) {
	addr := os.Getenv("REGIRSTY_ADDR")
	if addr == "" {
		addr = ":8000"
	}
	opt.Address = []string{addr}
}

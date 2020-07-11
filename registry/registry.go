package main

import (
	"context"
	"sync"
	"time"

	pb "github.com/frankffenn/worker-srv/registry/service/proto"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/memory"
)

var (
	interval   = 5 * time.Second
	defaultTTL = 1 * time.Minute
)

type Server struct {
	sync.RWMutex

	working  map[string]int64
	service  map[string]uint64
	registry registry.Registry
}

func NewServer(opts ...registry.Option) *Server {
	s := &Server{
		service:  make(map[string]uint64, 0),
		working:  make(map[string]int64, 0),
		registry: memory.NewRegistry(opts...),
	}

	go s.check()

	return s
}

func (s *Server) check() {
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			s.Lock()
			for k, v := range s.working {
				start := time.Unix(v, 0)
				if time.Since(start) > defaultTTL {
					delete(s.working, k)
					s.service[k]++
				}
			}
			s.Unlock()
		}
	}
}

func (s *Server) GetService(ctx context.Context, req *pb.GetRequest, rsp *pb.GetResponse) error {
	services, err := s.registry.GetService(req.Service)
	if err != nil {
		return err
	}

	for _, service := range services {
		rsp.Services = append(rsp.Services, ToProto(service))
	}

	return nil
}

func (s *Server) Register(ctx context.Context, req *pb.Service, rsp *pb.EmptyResponse) error {
	s.Lock()
	defer s.Unlock()

	for _, node := range req.Nodes {
		if _, ok := s.service[node.Id]; !ok {
			s.service[node.Id] = 1
		}
	}
	return s.registry.Register(ToService(req))
}

func (s *Server) Deregister(ctx context.Context, req *pb.Service, rsp *pb.EmptyResponse) error {
	s.Lock()
	defer s.Unlock()

	for _, node := range req.Nodes {
		delete(s.service, node.Id)
	}

	return s.registry.Deregister(ToService(req))
}

func (s *Server) ListServices(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	services, err := s.registry.ListServices()
	if err != nil {
		return err
	}

	for _, service := range services {
		rsp.Services = append(rsp.Services, ToProto(service))
	}

	return nil
}

func (s *Server) Watch(ctx context.Context, req *pb.WatchRequest, rsp pb.Registry_WatchStream) error {
	_, err := s.registry.Watch()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Mark(ctx context.Context, req *pb.MarkRequest, rsp *pb.EmptyResponse) error {
	s.Lock()
	defer s.Unlock()

	if s.service[req.Id] <= 0 {
		return nil
	}

	s.working[req.Id] = time.Now().Unix()
	s.service[req.Id]--
	return nil
}

func (s *Server) Reset(ctx context.Context, req *pb.ResetRequest, rsp *pb.EmptyResponse) error {
	s.Lock()
	defer s.Unlock()

	delete(s.working, req.Id)
	s.service[req.Id]++
	return nil
}

func (s *Server) GetServices(ctx context.Context, req *pb.GetServicesRequest, rsp *pb.GetServicesResponse) error {
	s.Lock()
	defer s.Unlock()
	rsp.Services = s.service
	return nil
}

var _ pb.RegistryHandler = &Server{}

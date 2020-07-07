package main

import (
	"context"
	"log"
	"sync/atomic"
	"time"

	ffi "github.com/filecoin-project/filecoin-ffi"
	"github.com/filecoin-project/specs-actors/actors/abi"
	storage2 "github.com/filecoin-project/specs-storage/storage"
	pb "github.com/frankffenn/worker-srv/proto"
	"golang.org/x/xerrors"
)

var ErrServerBusy = xerrors.New("server busy")

type Server struct {
	max uint64
}

func NewServer(max uint64) *Server {
	return &Server{
		max: max,
	}
}

func (s *Server) SealCommit2(ctx context.Context, req *pb.SealCommit2Request, rsp *pb.SealCommit2Response) error {
	log.Println("sealcommit2 start")
	old := atomic.LoadUint64(&s.max)
	if old <= 0 {
		log.Println("--------- server busy >>>>>")
		return ErrServerBusy
	}

	atomic.CompareAndSwapUint64(&s.max, old, old-1)
	defer func() {
		old := atomic.LoadUint64(&s.max)
		atomic.CompareAndSwapUint64(&s.max, old, old+1)
	}()

	phase1Out := storage2.Commit1Out(req.Commit1Out)
	sector := abi.SectorID{
		Number: abi.SectorNumber(req.Sector.Number),
		Miner:  abi.ActorID(req.Sector.Miner),
	}
	ret, err := ffi.SealCommitPhase2(phase1Out, sector.Number, sector.Miner)
	if err != nil {
		return err
	}

	rsp.Proof = ret
	log.Printf("done, but wait a minute ...(idle:%d)", old)
	<-time.After(1 * time.Minute)
	return nil
}

var _ pb.WorkerServiceHandler = &Server{}

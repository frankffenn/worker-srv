package main

import (
	"context"
	"fmt"
	"log"

	ffi "github.com/filecoin-project/filecoin-ffi"
	"github.com/filecoin-project/specs-actors/actors/abi"
	storage2 "github.com/filecoin-project/specs-storage/storage"
	pb "github.com/frankffenn/worker-srv/proto"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) SealCommit2(ctx context.Context, req *pb.SealCommit2Request, rsp *pb.SealCommit2Response) error {
	phase1Out := storage2.Commit1Out(req.Commit1Out)
	sector := abi.SectorID{
		Number: abi.SectorNumber(req.Sector.Number),
		Miner:  abi.ActorID(req.Sector.Miner),
	}
	log.Printf("SealCommit2 start, %s\n", s.SectorName(sector))
	ret, err := ffi.SealCommitPhase2(phase1Out, sector.Number, sector.Miner)
	if err != nil {
		return err
	}

	rsp.Proof = ret
	log.Printf("SealCommit2 end, %s\n", s.SectorName(sector))
	return nil
}

func (s *Server) SectorName(sid abi.SectorID) string {
	return fmt.Sprintf("s-t0%d-%d", sid.Miner, sid.Number)
}

var _ pb.WorkerServiceHandler = &Server{}

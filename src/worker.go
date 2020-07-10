package main

import (
	"context"
	"log"

	ffi "github.com/filecoin-project/filecoin-ffi"
	"github.com/filecoin-project/specs-actors/actors/abi"
	storage2 "github.com/filecoin-project/specs-storage/storage"
	pb "github.com/frankffenn/worker-srv/proto"
)

var (
	DefaultService = "registry.center"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) SealCommit2(ctx context.Context, req *pb.SealCommit2Request, rsp *pb.SealCommit2Response) error {
	log.Println("----------------->>> start SealCommit2")
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
	log.Println("----------------->>> finish SealCommit2")
	return nil
}

var _ pb.WorkerServiceHandler = &Server{}

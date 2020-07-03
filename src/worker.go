package src

import (
	"context"

	ffi "github.com/filecoin-project/filecoin-ffi"
	"github.com/filecoin-project/specs-actors/actors/abi"
	storage2 "github.com/filecoin-project/specs-storage/storage"
	pb "github.com/frankffenn/worker-srv/proto"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) SealCommit2(ctx context.Context, req *pb.SealCommitRequest, rsp *pb.SealCommitResponse) error {
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
	return nil
}

var _ pb.WorkerServiceHandler = &Server{}

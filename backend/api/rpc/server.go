package rpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "powervoting-server/api/rpc/proto"
	"powervoting-server/service"
)

type BackendRpc struct {
	*pb.UnimplementedBackendServer
	rpcSrv *service.RpcService
}

func NewBackendRpc(rpcSrv *service.RpcService) *BackendRpc {
	return &BackendRpc{
		rpcSrv: rpcSrv,
	}
}

func (b *BackendRpc) GetAllVoterAddresss(ctx context.Context, req *pb.GetAllVoterAddressRequest) (*pb.GetAllVoterAddressResponse, error) {
	addresss, err := b.rpcSrv.GetAllVoterAddresss(ctx, req.ChainId)
	if err != nil {
		return &pb.GetAllVoterAddressResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetAllVoterAddressResponse{
		Addresses: addresss,
	}, err
}

func (b *BackendRpc) GetVoterInfo(ctx context.Context, req *pb.GetVoterInfoRequest) (*pb.GetVoterInfoResponse, error) {
	voteInfo, err := b.rpcSrv.GetVoterInfoByAddress(ctx, req.Address)
	if err != nil {
		return &pb.GetVoterInfoResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetVoterInfoResponse{
		MinerIds:      voteInfo.MinerIds,
		ActorId:       voteInfo.OwnerId,
		GithubAccount: voteInfo.GithubId,
	}, nil

}

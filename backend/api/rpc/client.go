package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "powervoting-server/api/rpc/proto"
	"powervoting-server/config"
	"powervoting-server/constant"
	"powervoting-server/model"
)

var (
	snapshotClient pb.SnapshotClient
	clientOnce     sync.Once
)

// getClient returns a singleton gRPC client instance.
func getClient() pb.SnapshotClient {
	clientOnce.Do(func() {

		conn, err := grpc.Dial(
			config.Client.Snapshot.Rpc,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		if err != nil {
			zap.L().Error("failed to connect to gRPC server", zap.Error(err))
		}
		snapshotClient = pb.NewSnapshotClient(conn)
	})
	return snapshotClient
}

// GetAddressPower fetches power information from the gRPC server.
func GetAddressPower(netId int64, address string, randomNum int32) (model.Power, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	grpcReq := &pb.AddressPowerRequest{
		NetId:     netId,
		Address:   address,
		RandomNum: randomNum,
	}

	grpcRes, err := getClient().GetAddressPower(ctx, grpcReq)
	if err != nil {
		return model.Power{}, fmt.Errorf("failed to get address power: %v", err)
	}

	power, err := ParseAddressPowerResponse(grpcRes)
	if err != nil {
		return model.Power{}, fmt.Errorf("failed to parse address power response: %v", err)
	}

	return power, nil
}

// ParseAddressPowerResponse parses gRPC response into model.Power.
func ParseAddressPowerResponse(res *pb.AddressPowerResponse) (model.Power, error) {
	var power model.Power

	spPower := new(big.Int)
	clientPower := new(big.Int)
	tokenHolderPower := new(big.Int)
	developerPower := new(big.Int)

	if _, ok := spPower.SetString(res.SpPower, 10); !ok {
		return model.Power{}, fmt.Errorf("failed to parse SpPower: %s", res.SpPower)
	}
	if _, ok := clientPower.SetString(res.ClientPower, 10); !ok {
		return model.Power{}, fmt.Errorf("failed to parse ClientPower: %s", res.ClientPower)
	}
	if _, ok := tokenHolderPower.SetString(res.TokenHolderPower, 10); !ok {
		return model.Power{}, fmt.Errorf("failed to parse TokenHolderPower: %s", res.TokenHolderPower)
	}
	if _, ok := developerPower.SetString(res.DeveloperPower, 10); !ok {
		return model.Power{}, fmt.Errorf("failed to parse DeveloperPower: %s", res.DeveloperPower)
	}

	power.SpPower = spPower
	power.ClientPower = clientPower
	power.TokenHolderPower = tokenHolderPower
	power.DeveloperPower = developerPower
	power.BlockHeight = new(big.Int).SetInt64(res.BlockHeight)

	return power, nil
}

func GetDataHeight(netId int64, day string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.DataHeightRequest{
		NetId: netId,
		Day:   day,
	}

	res, err := getClient().GetDataHeight(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("failed to get address power: %v", err)
	}

	return res.Height, nil
}

func GetAddressPowerByDay(netId int64, address, day string) (model.Power, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constant.RequestTimeout)
	defer cancel()

	grpcReq := &pb.AddressPowerByDayRequest{
		NetId:   netId,
		Address: address,
		Day:     day,
	}

	grpcRes, err := getClient().GetAddressPowerByDay(ctx, grpcReq)
	if err != nil {
		return model.Power{}, fmt.Errorf("failed to get address power: %v", err)
	}

	power, err := ParseAddressPowerResponse(grpcRes)
	if err != nil {
		return model.Power{}, fmt.Errorf("failed to parse address power response: %v", err)
	}

	return power, nil
}

func GetAllAddressPowerByDay(chainId int64, snapshotDay string) (model.SnapshotAllPower, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constant.RequestTimeout)
	defer cancel()

	grpcReq := &pb.GetAllAddrPowerByDayRequest{
		NetId: chainId,
		Day:   snapshotDay,
	}

	grpcRes, err := getClient().GetAllAddrPowerByDay(ctx, grpcReq)
	if err != nil {
		return model.SnapshotAllPower{}, fmt.Errorf("failed to get all address power: %v", err)
	}

	// Declare a variable to hold the unmarshalled power information.
	var allPower model.SnapshotAllPower
	if err := json.Unmarshal([]byte(grpcRes.Info), &allPower); err != nil {
		return model.SnapshotAllPower{}, err
	}

	return allPower, nil
}

func UploadSnapshotInfo(chainId int64, snapshotDay string) (model.SnapshotHeight, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constant.RequestTimeout)
	defer cancel()

	grpcReq := &pb.UploadSnapshotInfoByDayRequest{
		Day:   snapshotDay,
		NetId: chainId,
	}

	grpcRes, err := getClient().UploadSnapshotInfoByDay(ctx, grpcReq)
	if err != nil {
		return model.SnapshotHeight{}, fmt.Errorf("failed to sync snapshot: %v", err)
	}

	return model.SnapshotHeight{
		Height: grpcRes.Height,
		Day:    grpcRes.Day,
	}, nil
}

func SyncAddressPower(chainId int64, ethAddr string) {
	ctx, cancel := context.WithTimeout(context.Background(), constant.RequestTimeout)
	defer cancel()

	grpcReq := &pb.SyncAddrPowerRequest{}
	if _, err := getClient().SyncAddrPower(ctx, grpcReq); err != nil {
		zap.L().Error("failed to sync address power", zap.Error(err))
	}

}

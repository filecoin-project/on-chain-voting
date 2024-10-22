package client

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"powervoting-server/config"
	"powervoting-server/model"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "powervoting-server/client/proto"
)

var (
	clientInstance pb.SnapshotClient
	clientOnce     sync.Once
)

// getClient returns a singleton gRPC client instance.
func getClient() pb.SnapshotClient {
	clientOnce.Do(func() {
		conn, err := grpc.Dial(config.Client.Snapshot.Rpc, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err != nil {
			log.Fatalf("failed to connect to gRPC server: %v", err)
		}
		clientInstance = pb.NewSnapshotClient(conn)
	})
	return clientInstance
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

func GetAddressPowerByDay(netId int64, address string, day string) (model.Power, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
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

func GetAllAddressPowerByDay(netId int64, day string) (model.AllPowerByDay, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcReq := &pb.GetAllAddrPowerByDayRequest{
		NetId: netId,
		Day:   day,
	}

	grpcRes, err := getClient().GetAllAddrPowerByDay(ctx, grpcReq)
	if err != nil {
		zap.L().Error("failed to get all address power: ", zap.Error(err))
		return model.AllPowerByDay{}, err
	}

	powerInfos := model.AllPowerByDay{
		Id:        netId,
		Day:       day,
		PowerInfo: grpcRes.Info,
	}

	return powerInfos, nil
}

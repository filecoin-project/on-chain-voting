// Copyright (C) 2023-2024 StorSwift Inc.
// This file is part of the PowerVoting library.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handler

import (
	"context"
	"encoding/json"

	"github.com/golang-module/carbon"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "power-snapshot/api/proto"
	"power-snapshot/internal/service"
	"power-snapshot/utils"
)

type Snapshot struct {
	*pb.UnimplementedSnapshotServer
	querySrv *service.QueryService
	syncSrv  *service.SyncService
}

var _ pb.SnapshotServer = (*Snapshot)(nil)

func NewSnapshot(query *service.QueryService, sync *service.SyncService) *Snapshot {
	return &Snapshot{
		querySrv: query,
		syncSrv:  sync,
	}
}

// GetAddressPower retrieves the power information for a given address.
// It takes a context and an AddressPowerRequest as input and returns an AddressPowerResponse or an error.
// The function queries the address power using the provided network ID, address, and random number.
// If the query fails, it returns an internal error with the corresponding error message.
// Otherwise, it constructs and returns an AddressPowerResponse with the retrieved power details.
func (s *Snapshot) GetAddressPower(ctx context.Context, req *pb.AddressPowerRequest) (*pb.AddressPowerResponse, error) {
	m, err := s.querySrv.GetAddressPower(ctx, req.GetNetId(), utils.EthStandardAddressToHex(req.GetAddress()), req.GetRandomNum())
	if err != nil {
		return &pb.AddressPowerResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &pb.AddressPowerResponse{
		Address:          m.Address,
		SpPower:          m.SpPower.String(),
		ClientPower:      m.ClientPower.String(),
		TokenHolderPower: m.TokenHolderPower.String(),
		DeveloperPower:   m.DeveloperPower.String(),
		BlockHeight:      m.BlockHeight,
		DateStr:          m.DateStr,
	}, nil
}

// GetAddressPowerByDay retrieves the power information for a given address on a specific day.
// It takes a context and an AddressPowerByDayRequest as input and returns an AddressPowerResponse or an error.
// The function parses the provided day string into a timestamp representing the end of the day.
// It then queries the address power using the network ID, address, day, and the parsed timestamp.
// If the query fails, it returns an internal error with the corresponding error message.
// Otherwise, it constructs and returns an AddressPowerResponse with the retrieved power details.
func (s *Snapshot) GetAddressPowerByDay(ctx context.Context, req *pb.AddressPowerByDayRequest) (*pb.AddressPowerResponse, error) {
	day := req.GetDay()
	dayTime := carbon.Parse(day).EndOfDay().ToStdTime()
	m, err := s.querySrv.GetAddressPowerByDay(ctx, req.GetNetId(), utils.EthStandardAddressToHex(req.GetAddress()), day, dayTime)
	if err != nil {
		return &pb.AddressPowerResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &pb.AddressPowerResponse{
		Address:          m.Address,
		SpPower:          m.SpPower.String(),
		ClientPower:      m.ClientPower.String(),
		TokenHolderPower: m.TokenHolderPower.String(),
		DeveloperPower:   m.DeveloperPower.String(),
		BlockHeight:      m.BlockHeight,
		DateStr:          m.DateStr,
	}, nil
}

// GetDataHeight retrieves the data height for a given network ID and day.
// It takes a DataHeightRequest as input and returns a DataHeightResponse or an error.
// The function queries the data height using the provided network ID and day.
// If the query fails, it returns an internal error with the corresponding error message.
// Otherwise, it constructs and returns a DataHeightResponse containing the day and the retrieved height.
func (s *Snapshot) GetDataHeight(_ context.Context, req *pb.DataHeightRequest) (*pb.DataHeightResponse, error) {
	height, err := s.querySrv.GetDataHeight(context.Background(), req.GetNetId(), req.GetDay())
	if err != nil {
		return &pb.DataHeightResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &pb.DataHeightResponse{
		Day:    req.GetDay(),
		Height: height,
	}, nil
}

// GetAllAddrPowerByDay retrieves the power information for all addresses on a specific day.
// It takes a GetAllAddrPowerByDayRequest as input and returns a GetAllAddrPowerByDayResponse or an error.
// The function queries the power information for all addresses using the provided network ID and day.
// If the query fails, it returns an internal error with the corresponding error message.
// The retrieved address powers are then marshaled into a JSON string.
// If marshaling fails, it returns an internal error with the corresponding error message.
// Otherwise, it constructs and returns a GetAllAddrPowerByDayResponse containing the day, network ID, and the JSON-encoded power information.
func (s *Snapshot) GetAllAddrPowerByDay(_ context.Context, req *pb.GetAllAddrPowerByDayRequest) (*pb.GetAllAddrPowerByDayResponse, error) {
	addrPowers, err := s.querySrv.GetAllAddressPowerByDay(context.Background(), req.GetNetId(), req.GetDay())
	if err != nil {
		return &pb.GetAllAddrPowerByDayResponse{}, status.Error(codes.Internal, err.Error())
	}

	jsonRes, err := json.Marshal(addrPowers)
	if err != nil {
		return &pb.GetAllAddrPowerByDayResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetAllAddrPowerByDayResponse{
		Day:   req.Day,
		Info:  string(jsonRes),
		NetId: req.NetId,
	}, nil
}

// SyncDateHeight synchronizes the data height for a given network ID.
// It takes a SyncDateHeightRequest as input and returns a SyncDateHeightResponse or an error.
// The function triggers a synchronization process for the data height using the provided network ID.
// If the synchronization fails, it returns an internal error with the corresponding error message.
// Otherwise, it returns an empty SyncDateHeightResponse indicating success.
func (s *Snapshot) SyncDateHeight(_ context.Context, req *pb.SyncDateHeightRequest) (*pb.SyncDateHeightResponse, error) {
	err := s.syncSrv.SyncDateHeight(context.Background(), req.GetNetId())
	if err != nil {
		return &pb.SyncDateHeightResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &pb.SyncDateHeightResponse{}, nil
}

// SyncAddrPower synchronizes the power information for a specific address on a given network.
// It takes a SyncAddrPowerRequest as input and returns a SyncAddrPowerResponse or an error.
// The function triggers a synchronization process for the address power using the provided network ID and address.
// If the synchronization fails, it returns an internal error with the corresponding error message.
// Otherwise, it returns an empty SyncAddrPowerResponse indicating success.
func (s *Snapshot) SyncAddrPower(_ context.Context, req *pb.SyncAddrPowerRequest) (*pb.SyncAddrPowerResponse, error) {
	err := s.syncSrv.SyncDateHeight(context.Background(), req.GetNetId())
	if err != nil {
		return &pb.SyncAddrPowerResponse{}, status.Error(codes.Internal, err.Error())
	}

	err = s.syncSrv.AddAddrPowerTaskToMQ(context.Background(), req.GetNetId(), utils.EthStandardAddressToHex(req.GetAddress()))
	if err != nil {
		return &pb.SyncAddrPowerResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &pb.SyncAddrPowerResponse{}, nil
}

// SyncAllAddrPower synchronizes the power information for all addresses on a given network.
// It takes a SyncAllAddrPowerRequest as input and returns a SyncAllAddrPowerResponse or an error.
// The function triggers a synchronization process for the power information of all addresses using the provided network ID.
// If the synchronization fails, it returns an internal error with the corresponding error message.
// Otherwise, it returns an empty SyncAllAddrPowerResponse indicating success.
func (s *Snapshot) SyncAllAddrPower(_ context.Context, req *pb.SyncAllAddrPowerRequest) (*pb.SyncAllAddrPowerResponse, error) {
	err := s.syncSrv.SyncDateHeight(context.Background(), req.GetNetId())
	if err != nil {
		return &pb.SyncAllAddrPowerResponse{}, status.Error(codes.Internal, err.Error())
	}

	err = s.syncSrv.SyncAllAddrPower(context.Background(), req.GetNetId())
	if err != nil {
		return &pb.SyncAllAddrPowerResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &pb.SyncAllAddrPowerResponse{}, nil
}

// UploadSnapshotInfoByDay upload the snapshot data for a specific day and network.
// It takes a SyncSnapshotRequest as input and returns a SyncSnapshotResponse or an error.
// The function retrieves all address power information for the given day and network ID.
// If the retrieval fails, it returns an internal error with the corresponding error message.
// Otherwise, it triggers a synchronization process for the snapshot using the retrieved data, day, and network ID.
// If the synchronization fails, it returns an internal error with the corresponding error message.
// On success, it returns a SyncSnapshotResponse containing the day and the snapshot height.
func (s *Snapshot) UploadSnapshotInfoByDay(_ context.Context, req *pb.UploadSnapshotInfoByDayRequest) (*pb.UploadSnapshotInfoByDayResponse, error) {
	allPower, err := s.querySrv.GetAllAddressPowerByDay(context.Background(), req.GetNetId(), req.GetDay())
	if err != nil {
		zap.L().Error("get all addr power by day error", zap.Error(err))
	}

	snapshotHeight, err := s.syncSrv.UploadSnapshotInfoByDay(context.Background(), allPower, req.GetDay(), req.GetNetId())
	if err != nil {
		return &pb.UploadSnapshotInfoByDayResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &pb.UploadSnapshotInfoByDayResponse{
		Day:    req.Day,
		Height: snapshotHeight,
	}, nil
}

func (s *Snapshot) SyncAllDeveloperWeight(context.Context, *pb.SyncAllDeveloperWeightRequest) (*pb.SyncAllDeveloperWeightResponse, error) {
	err := s.syncSrv.SyncLatestDeveloperWeight(context.Background())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.SyncAllDeveloperWeightResponse{}, nil
}

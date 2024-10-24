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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "power-snapshot/api/proto"
	"power-snapshot/internal/service"
)

type Snapshot struct {
	*pb.UnimplementedSnapshotServer
	querySrv *service.QueryService
	syncSrv  *service.SyncService
}

func NewSnapshot(query *service.QueryService, sync *service.SyncService) *Snapshot {
	return &Snapshot{
		querySrv: query,
		syncSrv:  sync,
	}
}

func (s *Snapshot) GetAddressPower(ctx context.Context, req *pb.AddressPowerRequest) (*pb.AddressPowerResponse, error) {
	m, err := s.querySrv.GetAddressPower(ctx, req.GetNetId(), req.GetAddress(), req.GetRandomNum())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
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

func (s *Snapshot) GetAddressPowerByDay(ctx context.Context, req *pb.AddressPowerByDayRequest) (*pb.AddressPowerResponse, error) {
	day := req.GetDay()
	dayTime := carbon.Parse(day).EndOfDay().ToStdTime()
	m, err := s.querySrv.GetAddressPowerByDay(ctx, req.GetNetId(), req.GetAddress(), day, dayTime)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
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

func (s *Snapshot) GetDataHeight(_ context.Context, req *pb.DataHeightRequest) (*pb.DataHeightResponse, error) {
	height, err := s.querySrv.GetDataHeight(context.Background(), req.GetNetId(), req.GetDay())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DataHeightResponse{
		Day:    req.GetDay(),
		Height: height,
	}, nil
}

func (s *Snapshot) GetAllAddrPowerByDay(_ context.Context, req *pb.GetAllAddrPowerByDayRequest) (*pb.GetAllAddrPowerByDayResponse, error) {
	addrPowers, err := s.querySrv.GetAllAddressPowerByDay(context.Background(), req.GetNetId(), req.GetDay())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	jsonRes, err := json.Marshal(addrPowers)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetAllAddrPowerByDayResponse{
		Day:   req.Day,
		Info:  string(jsonRes),
		NetId: req.NetId,
	}, nil
}

func (s *Snapshot) SyncDateHeight(_ context.Context, req *pb.SyncDateHeightRequest) (*pb.SyncDateHeightResponse, error) {
	err := s.syncSrv.SyncDateHeight(context.Background(), req.GetNetId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.SyncDateHeightResponse{}, nil
}

func (s *Snapshot) SyncAddrPower(_ context.Context, req *pb.SyncAddrPowerRequest) (*pb.SyncAddrPowerResponse, error) {
	err := s.syncSrv.SyncAddrPower(context.Background(), req.GetNetId(), req.GetAddress())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.SyncAddrPowerResponse{}, nil
}

func (s *Snapshot) SyncAllAddrPower(_ context.Context, req *pb.SyncAllAddrPowerRequest) (*pb.SyncAllAddrPowerResponse, error) {
	err := s.syncSrv.SyncAllAddrPower(context.Background(), req.GetNetId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.SyncAllAddrPowerResponse{}, nil
}

func (s *Snapshot) SyncAllDeveloperWeight(context.Context, *pb.SyncAllDeveloperWeightRequest) (*pb.SyncAllDeveloperWeightResponse, error) {
	err := s.syncSrv.SyncAllDeveloperWeight(context.Background())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.SyncAllDeveloperWeightResponse{}, nil
}

func (s *Snapshot) SyncDeveloperWeight(ctx context.Context, req *pb.SyncDeveloperWeightRequest) (*pb.SyncDeveloperWeightResponse, error) {
	err := s.syncSrv.SyncDeveloperWeight(context.Background(), req.GetDateStr())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.SyncDeveloperWeightResponse{}, nil
}

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

package data

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"

	"power-snapshot/config"
)

type StreamClient struct {
	nc *nats.Conn
	jetstream.JetStream
}

func NewJetstreamClient() (*StreamClient, error) {
	nc, err := nats.Connect(config.Client.Nats.URI)
	if err != nil {
		zap.S().Error("Failed to connect to NATS", zap.Error(err))
		return nil, err
	}

	js, err := jetstream.New(nc)
	if err != nil {
		zap.S().Error("Failed to init jetstream", zap.Error(err))
	}

	return &StreamClient{
		nc, js,
	}, nil
}

// Drain is a method on the StreamClient struct that is responsible for draining the connection.
func (s *StreamClient) Drain() error {
	err := s.nc.Drain()
	if err != nil {
		zap.S().Error(err)
		return err
	}

	return nil
}

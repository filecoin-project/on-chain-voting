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
	"bytes"
	"errors"
	"net/http"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ipfs/go-cid"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/multiformats/go-multihash"
	"github.com/web3-storage/go-ucanto/core/delegation"
	"github.com/web3-storage/go-ucanto/did"
	"github.com/web3-storage/go-ucanto/principal"
	"github.com/web3-storage/go-ucanto/principal/ed25519/signer"
	"github.com/web3-storage/go-w3up/capability/storeadd"
	"github.com/web3-storage/go-w3up/client"
	godelegation "github.com/web3-storage/go-w3up/delegation"
	"go.uber.org/zap"

	"power-snapshot/config"
	models "power-snapshot/internal/model"
)

var W3 W3Client

type W3Client struct {
	Space  did.DID
	Issuer principal.Signer
	Proof  delegation.Delegation
}

func NewW3Client() *W3Client {
	space, err := did.Parse(config.Client.W3Client.Space)
	if err != nil {
		zap.L().Error("w3client parse space error:", zap.Error(err))
		return &W3Client{}
	}

	issuer, err := signer.Parse(config.Client.W3Client.PrivateKey)
	if err != nil {
		zap.L().Error("w3client parse private error:", zap.Error(err))
		return &W3Client{}
	}

	prfbytes, err := os.ReadFile(config.Client.W3Client.Proof)
	if err != nil {
		zap.L().Error("read proof.ucan file error:", zap.Error(err))
		return &W3Client{}
	}

	proof, err := godelegation.ExtractProof(prfbytes)
	if err != nil {
		zap.L().Error("init w3storage error:", zap.Error(err))
		return &W3Client{}
	}

	return &W3Client{
		Space:  space,
		Issuer: issuer,
		Proof:  proof,
	}
}

func (w *W3Client) UploadByte(data []byte) (string, error) {
	// generate the CID for the CAR
	mh, _ := multihash.Sum(data, multihash.SHA2_256, -1)
	link := cidlink.Link{Cid: cid.NewCidV1(0x0202, mh)}
	rcpt, _ := client.StoreAdd(
		w.Issuer,
		w.Space,
		&storeadd.Caveat{Link: link, Size: uint64(len(data))},
		client.WithProofs([]delegation.Delegation{w.Proof}),
	)

	if rcpt.Out().Ok().Status == "upload" {
		hr, _ := http.NewRequest("PUT", *rcpt.Out().Ok().Url, bytes.NewReader(data))

		hdr := map[string][]string{}
		for k, v := range rcpt.Out().Ok().Headers.Values {
			hdr[k] = []string{v}
		}

		hr.Header = hdr
		hr.ContentLength = int64(len(data))
		httpClient := http.Client{}
		res, _ := httpClient.Do(hr)
		res.Body.Close()
	}

	return link.String(), nil
}


type GoEthClientManager struct {
	lock        sync.Mutex
	instanceMap map[int64]models.GoEthClient
}

func NewGoEthClientManager(network []models.Network ) (*GoEthClientManager, error) {
	instMap := make(map[int64]models.GoEthClient)
	for _, net := range network {
		client, err := getGoEthClient(net)
		if err != nil {
			zap.L().Error("init eth client failed", zap.Error(err))
			return nil, err
		}
		instMap[net.ChainId] = client
	}

	return &GoEthClientManager{
		lock:        sync.Mutex{},
		instanceMap: instMap,
	}, nil
}

func (g *GoEthClientManager) GetClient(netId int64) (models.GoEthClient, error) {
	g.lock.Lock()
	defer g.lock.Unlock()
	client, ok := g.instanceMap[netId]
	if !ok {
		zap.L().Error("client not exist", zap.Int64("netId", netId))
		return models.GoEthClient{}, errors.New("client not exist")
	}

	return client, nil
}

func (g *GoEthClientManager) ListClientNetWork() []int64 {
	g.lock.Lock()
	defer g.lock.Unlock()
	result := make([]int64, 0, len(g.instanceMap))
	for netId := range g.instanceMap {
		result = append(result, netId)
	}

	return result
}

// getGoEthClient initializes a Go-ethereum client with the provided configuration.
func getGoEthClient(network models.Network) (models.GoEthClient, error) {
	rpcs := make([]*ethclient.Client, 0)
	for _, rpc := range network.QueryRpc {
		queryClient, err := ethclient.Dial(rpc)
		if err != nil {
			zap.L().Warn("ethClient.Dial error: ",
				zap.String("network", network.Name),
				zap.String("rpc", rpc), zap.Error(err))
			continue
		}
		rpcs = append(rpcs, queryClient)
	}

	// generate goEthClient
	goEthClient := models.GoEthClient{
		ChainId:     network.ChainId,
		Name:        network.Name,
		IdPrefix:    network.IdPrefix,
		QueryClient: rpcs,
		QueryRpc:    network.QueryRpc,
	}
	return goEthClient, nil
}

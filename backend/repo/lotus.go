package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/ybbus/jsonrpc/v3"

	"powervoting-server/config"
	"powervoting-server/model"
	"powervoting-server/utils/types"
)

type LotusRPCRepo struct {
	client jsonrpc.RPCClient
}

func NewLotusRPCRepo() *LotusRPCRepo {
	return &LotusRPCRepo{
		client: jsonrpc.NewClientWithOpts(config.Client.Network.Rpc, &jsonrpc.RPCClientOpts{}),
	}
}

func (l *LotusRPCRepo) FilecoinAddressToID(ctx context.Context, addr string) (string, error) {
	resp, err := l.client.Call(ctx, "Filecoin.StateLookupID", addr, types.TipSetKey{})
	if err != nil {
		return "", err
	}

	if resp.Error != nil {
		return "", resp.Error
	}

	return resp.Result.(string), nil
}

func (l *LotusRPCRepo) EthAddrToFilcoinAddr(ctx context.Context, addr string) (string, error) {
	if !strings.HasPrefix(addr, "0x") {
		return addr, nil
	}

	resp, err := l.client.Call(ctx, "Filecoin.EthAddressToFilecoinAddress", addr)
	if err != nil {
		return "", err
	}

	if resp.Error != nil {
		return "", resp.Error
	}

	return resp.Result.(string), nil
}

func (l *LotusRPCRepo) getOwnerByMinerId(ctx context.Context, minerId string) (string, error) {
	resp, err := l.client.Call(ctx, "Filecoin.StateMinerInfo", minerId, types.TipSetKey{})
	if err != nil {
		return "", err
	}

	if resp.Error != nil {
		return "", resp.Error
	}

	owner := resp.Result.(map[string]any)["Owner"]
	return owner.(string), nil
}

func (l *LotusRPCRepo) GetActorIdByAddress(ctx context.Context, addr string) (string, error) {
	filcoinAddr, err := l.EthAddrToFilcoinAddr(ctx, addr)
	if err != nil {
		return "", err
	}

	return l.FilecoinAddressToID(ctx, filcoinAddr)
}

func (l *LotusRPCRepo) GetValidMinerIds(ctx context.Context, actorId string, minerIds []uint64) (model.StringSlice, error) {
	var ids model.StringSlice
	for _, id := range minerIds {
		minerIdStr := fmt.Sprintf("%s%d", config.Client.Network.MinerIdPrefix, id)
		owner, err := l.getOwnerByMinerId(ctx, minerIdStr)
		if err != nil {
			return nil, err
		}

		if owner == actorId {
			ids = append(ids, minerIdStr)
		}
	}

	return ids, nil
}

func (l *LotusRPCRepo) FilecoinAddrToEthAddr(ctx context.Context, addr string) (string, error) {
	if strings.HasPrefix(addr, "0x") {
		return addr, nil
	}

	resp, err := l.client.Call(ctx, "Filecoin.FilecoinAddressToEthAddress", addr)
	if err != nil {
		return "", err
	}

	if resp.Error != nil {
		return "", resp.Error
	}

	return resp.Result.(string), nil
}

package repo_test

import (
	"context"
	"testing"

	filecoinAddress "github.com/filecoin-project/go-address"
	"github.com/stretchr/testify/assert"

	"powervoting-server/config"
	"powervoting-server/repo"
)

func TestGetActorIdByAddress(t *testing.T) {
	config.GetDefaultConfig()

	lotusRepo := repo.NewLotusRPCRepo()

	res, err := lotusRepo.GetActorIdByAddress(context.Background(), "0x763D410594a24048537990dde6ca81c38CfF566a")
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}

func TestFilecoinAddressToID(t *testing.T) {
	config.GetDefaultConfig()

	lotusRepo := repo.NewLotusRPCRepo()
	res, err := lotusRepo.GetValidMinerIds(context.Background(), "t017386", []uint64{
		17387,
		28064,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}

func TestEthAddrToFilcoinAddr(t *testing.T) {
	config.GetDefaultConfig()

	lotusRepo := repo.NewLotusRPCRepo()
	res, err := lotusRepo.EthAddrToFilcoinAddr(context.Background(), "0xfF000000000000000000000000000000000278bc")
	assert.NoError(t, err)
	assert.NotEmpty(t, res)

	assert.Equal(t, "t0161980", res)
	addressStr, err := filecoinAddress.NewFromString(res)
	assert.NoError(t, err)
	assert.Equal(t, "f0161980", addressStr.String())
}

// func TestGetFilcoinAddr(t *testing.T) {
// 	config.Client.Network.Rpc = "http://192.168.11.139:1235/rpc/v1"

// 	client := jsonrpc.NewClientWithOpts(config.Client.Network.Rpc, &jsonrpc.RPCClientOpts{})
// 	resp, err := client.Call(context.Background(), "Filecoin.ChainGetTipSetByHeight", 2599912, types.TipSetKey{})
// 	assert.NoError(t, err)

// 	assert.NoError(t, resp.Error)

// 	rMap, ok := resp.Result.(map[string]interface{})
// 	assert.True(t, ok)

// 	resTipSet := rMap["Cids"].([]interface{})
// 	resp, err = client.Call(context.Background(), "Filecoin.StateGetActor", "t0161980", resTipSet)
// 	assert.NoError(t, err)
// 	assert.NoError(t, resp.Error)
// 	fmt.Printf("resp: %s\n", resp.Result)
// }

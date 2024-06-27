package utils

import (
	"power-snapshot/config"
	models "power-snapshot/internal/model"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGetVoterAddresses(t *testing.T) {
	config.InitLogger()
	err := config.InitConfig("../")
	assert.Nil(t, err)
	manager, err := NewGoEthClientManager(config.Client.Network)
	assert.Nil(t, err)
	client, err := manager.GetClient(314159)

	ethAddressList, err := GetVoterAddresses(client)
	assert.Nil(t, err)

	assert.NotEmpty(t, ethAddressList)
	zap.L().Info("result", zap.Any("ethAddressList", ethAddressList))
	zap.L().Info("result length", zap.Any("ethAddressListLength", len(ethAddressList)))
}

func TestGetVoterInfo(t *testing.T) {
	config.InitLogger()
	err := config.InitConfig("../")
	assert.Nil(t, err)
	manager, err := NewGoEthClientManager(config.Client.Network)
	assert.Nil(t, err)
	client, err := manager.GetClient(314159)

	voterInfo, err := GetVoterInfo("0x763D410594a24048537990dde6ca81c38CfF566a", client)
	assert.Nil(t, err)

	testVoterInfo := models.VoterInfo{
		ActorIds:      []uint64{35363},
		MinerIds:      []uint64{},
		GithubAccount: "",
		EthAddress:    common.HexToAddress("0x763D410594a24048537990dde6ca81c38CfF566a"),
		UcanCid:       "",
	}
	assert.Equal(t, testVoterInfo, voterInfo)

	zap.L().Info("result", zap.Any("voterInfo", voterInfo))
}

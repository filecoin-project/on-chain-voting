package config

import (
	"encoding/json"
	"powervoting-server/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	InitConfig("../")
	c := model.Config{
		Server: model.Server{
			Port: ":9999",
		},
		Mysql: model.Mysql{
			Url:      "localhost:3306",
			Username: "root",
			Password: "root",
		},
		Drand: model.Drand{
			Url: []string{
				"https://api.drand.secureweb3.com:6875",
				"https://api.drand.sh/",
				"https://api2.drand.sh/",
				"https://api3.drand.sh/",
				"https://drand.cloudflare.com/",
			},
			ChainHash: "",
		},
		Network: []model.Network{
			{
				Id:                  314159,
				Name:                "FileCoin-Calibration",
				Rpc:                 "https://api.calibration.node.glif.io/rpc/v1",
				PowerVotingAbi:      "power-voting.json",
				OracleAbi:           "oracle.json",
				PowerVotingContract: "0x1000000000000000000000000000000000000000",
				OracleContract:      "0x1000000000000000000000000000000000000000",
			},
		},
	}

	cjson, _ := json.Marshal(c)
	initC, _ := json.Marshal(Client)
	assert.Equal(t, string(cjson), string(initC))
}

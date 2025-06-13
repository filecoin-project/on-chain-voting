package config

import (
	"fil-vote/model"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
)

var Client model.Config
var PowerVotingAbi abi.ABI
var OracleAbi abi.ABI

// InitConfig initializes the configuration by reading from a YAML file located at the specified path.
func InitConfig(path string) error {
	// configuration file name
	viper.SetConfigName("configuration")

	viper.AddConfigPath(path)

	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		zap.L().Error("read config file error:", zap.Error(err))
		return err
	}

	err = viper.Unmarshal(&Client)
	if err != nil {
		zap.L().Error("unmarshal error:", zap.Error(err))
		return err
	}
	// open abi file and parse json
	powerVotingFile, err := os.Open(Client.ABIPath.PowerVotingABI)
	if err != nil {
		zap.L().Error("open power voting abi file error: ", zap.Error(err))
		return err
	}
	PowerVotingAbi, err = abi.JSON(powerVotingFile)
	if err != nil {
		zap.L().Error("abi.JSON error: ", zap.Error(err))
		return err
	}

	oracleFile, err := os.Open(Client.ABIPath.OracleABI)
	if err != nil {
		zap.L().Error("open oracle abi file error: ", zap.Error(err))
		return err
	}
	OracleAbi, err = abi.JSON(oracleFile)
	if err != nil {
		zap.L().Error("abi.JSON error: ", zap.Error(err))
		return err
	}
	return nil
}

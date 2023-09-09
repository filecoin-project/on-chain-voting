package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server  Server
	Drand   Drand
	Ipfs    Ipfs
	Network []Network
}

type Server struct {
	Port string
}

type Drand struct {
	Url       []string
	ChainHash string
}

type Ipfs struct {
	Token      string
	UploadPath string
}

type Network struct {
	Id              int64
	Name            string
	Rpc             string
	SubgraphUrl     string
	AbiPath         string
	ContractAddress string
	PrivateKey      string
	WalletAddress   string
	GasLimit        int64
}

// export client
var Client Config

// InitConfig initialization configuration
func InitConfig() {
	// configuration file name
	viper.SetConfigName("configuration")

	viper.AddConfigPath("./")

	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("read config file error:", err)
		return
	}

	err = viper.Unmarshal(&Client)
	if err != nil {
		fmt.Println("unmarshal error:", err)
		return
	}

}

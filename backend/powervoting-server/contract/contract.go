package contract

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"os"
	"powervoting-server/config"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GoEthClient go-ethereum client
type GoEthClient struct {
	Id              int64
	Name            string
	Client          *ethclient.Client
	Abi             abi.ABI
	Amount          *big.Int
	GasLimit        uint64
	ChainID         *big.Int
	ContractAddress common.Address
	PrivateKey      *ecdsa.PrivateKey
	WalletAddress   common.Address
}

// Proposal contract proposal
type Proposal struct {
	Cid         string
	ExpTime     *big.Int
	IsCounted   bool
	OptionIds   []*big.Int
	VoteResults []VoteResult
}

// VoteResult vote result
type VoteResult struct {
	OptionId *big.Int
	Votes    *big.Int
}

// ClientConfig config for get go-ethereum client
type ClientConfig struct {
	Id              int64
	Name            string
	Rpc             string
	ContractAddress string
	PrivateKey      string
	WalletAddress   string
	AbiPath         string
	GasLimit        int64
}

var (
	lock        sync.Mutex
	instanceMap map[int64]GoEthClient
)

func init() {
	instanceMap = make(map[int64]GoEthClient)
}

// GetClient get go eth client
func GetClient(id int64) (GoEthClient, error) {
	client, ok := instanceMap[id]
	if ok {
		return client, nil
	}
	networkList := config.Client.Network
	var clientConfig ClientConfig
	for _, network := range networkList {
		if network.Id == id {
			clientConfig = ClientConfig{
				Id:              network.Id,
				Name:            network.Name,
				Rpc:             network.Rpc,
				ContractAddress: network.ContractAddress,
				PrivateKey:      network.PrivateKey,
				WalletAddress:   network.WalletAddress,
				AbiPath:         network.AbiPath,
				GasLimit:        network.GasLimit,
			}
			break
		}
	}
	ethClient, err := getGoEthClient(clientConfig)
	if err != nil {
		return ethClient, err
	}
	instanceMap[id] = ethClient
	log.Printf("network init, net id: %d", id)
	return ethClient, nil
}

func GetEthMainClient() (GoEthClient, error) {
	client, ok := instanceMap[1]
	if ok {
		log.Println("get eth main client")
		return client, nil
	}
	ethClient, err := ethclient.Dial("https://ethereum.publicnode.com")
	if err != nil {
		log.Panic(err)
	}
	goEthClient := GoEthClient{
		Client: ethClient,
	}
	instanceMap[1] = goEthClient
	log.Println("init eth main client")
	return goEthClient, err
}

func GetGoerliClient() (GoEthClient, error) {
	client, ok := instanceMap[5]
	if ok {
		log.Println("get goerli client")
		return client, nil
	}
	ethClient, err := ethclient.Dial("https://ethereum-goerli.publicnode.com")
	if err != nil {
		log.Panic(err)
	}
	goEthClient := GoEthClient{
		Client: ethClient,
	}
	instanceMap[5] = goEthClient
	log.Println("init goerli client")
	return goEthClient, err
}

// getGoEthClient get go-ethereum client
func getGoEthClient(clientConfig ClientConfig) (GoEthClient, error) {
	client, err := ethclient.Dial(clientConfig.Rpc)
	if err != nil {
		log.Println("ethclient.Dial error: ", err)
		return GoEthClient{}, err
	}

	// contract address, wallet private key , wallet address
	contractAddress := common.HexToAddress(clientConfig.ContractAddress)
	privateKey, err := crypto.HexToECDSA(clientConfig.PrivateKey)
	walletAddress := common.HexToAddress(clientConfig.WalletAddress)

	// open abi file and parse json
	open, err := os.Open(clientConfig.AbiPath)
	if err != nil {
		log.Println("open abi file error: ", err)
		return GoEthClient{}, err
	}
	contractAbi, err := abi.JSON(open)
	if err != nil {
		log.Println("abi.JSON error: ", err)
		return GoEthClient{}, err
	}

	// transfer amount, if no set zero
	amount := big.NewInt(0)
	// gas limit
	gasLimit := uint64(clientConfig.GasLimit)
	// get chain id
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Println("client.NetworkID error: ", err)
		return GoEthClient{}, err
	}
	// generate goEthClient
	goEthClient := GoEthClient{
		Id:              clientConfig.Id,
		Name:            clientConfig.Name,
		Client:          client,
		Abi:             contractAbi,
		Amount:          amount,
		GasLimit:        gasLimit,
		ChainID:         chainID,
		ContractAddress: contractAddress,
		PrivateKey:      privateKey,
		WalletAddress:   walletAddress,
	}
	return goEthClient, nil
}

// CountContract contract count
func CountContract(id *big.Int, voteResult []VoteResult, voteListCid string, ethClient GoEthClient) error {
	// avoiding nonce duplication in multiple threads
	lock.Lock()
	defer lock.Unlock()

	// pack method and param
	var data []byte
	var err error
	data, err = ethClient.Abi.Pack("count", id, voteResult, voteListCid)
	if err != nil {
		log.Println("contractAbi.Pack error: ", err)
		return err
	}

	// get transaction nonce
	nonce, err := ethClient.Client.PendingNonceAt(context.Background(), ethClient.WalletAddress)
	if err != nil {
		log.Println("PendingNonceAt error: ", err)
		return err
	}

	// get suggest gas price
	gasPrice, err := ethClient.Client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println("client.SuggestGasPrice error: ", err)
		return err
	}
	log.Println("nonce: ", nonce)
	log.Println("GasPrice: ", gasPrice)

	gasFeeCap := gasPrice.Mul(gasPrice, big.NewInt(5))

	// create transaction
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   ethClient.ChainID,
		Nonce:     nonce,
		GasTipCap: gasPrice,
		GasFeeCap: gasFeeCap,
		Gas:       ethClient.GasLimit,
		To:        &ethClient.ContractAddress,
		Value:     ethClient.Amount,
		Data:      data,
	})
	// sign with private key
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(ethClient.ChainID), ethClient.PrivateKey)
	if err != nil {
		log.Println("types.SignTx error: ", err)
		return err
	}
	// send transaction
	err = ethClient.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Println("client.SendTransaction error: ", err)
		return err
	}

	log.Printf("net work：%s, transaction id: %s", ethClient.Name, signedTx.Hash().Hex())
	return nil
}

package scheduler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"io"
	"log"
	"math/big"
	"net/http"
	"powervoting-server/config"
	"powervoting-server/contract"
	"powervoting-server/model"
	"strings"
	"time"

	"github.com/storswiftlabs/tlock"
	drandhttp "github.com/storswiftlabs/tlock/networks/http"

	"github.com/robfig/cron/v3"
)

func TaskScheduler() {
	// create a new scheduler
	crontab := cron.New(cron.WithSeconds())
	// task function
	votingCountTask := VotingCountFunc

	// 1m
	votingCountSpec := "0 0/1 * * * ? "
	//votingCountSpec := "0/5 * * * * ? "

	// add task to scheduler
	_, err := crontab.AddFunc(votingCountSpec, votingCountTask)
	if err != nil {
		log.Println("add vote count task error: ", err)
	}

	// start
	crontab.Start()

	select {}
}

// VotingCountFunc vote count
func VotingCountFunc() {
	networkList := config.Client.Network
	for _, network := range networkList {
		ethClient, err := contract.GetClient(network.Id)
		if err != nil {
			log.Println("get go-eth client error:", err)
		}
		VotingCountHandler(network, ethClient)
	}
}

func getTimestamp(ethClient contract.GoEthClient) (int64, error) {
	ctx := context.Background()
	number, err := ethClient.Client.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}
	block, err := ethClient.Client.BlockByNumber(ctx, big.NewInt(int64(number)))
	if err != nil {
		return 0, err
	}
	now := int64(block.Time())
	return now, nil
}

var countMap map[string]bool

func init() {
	countMap = make(map[string]bool)
}

// VotingCountHandler voting count
func VotingCountHandler(network config.Network, ethClient contract.GoEthClient) {
	var now int64
	now, err := getTimestamp(ethClient)
	if err != nil {
		now = time.Now().Unix() - 20
	}
	proposals, err := SubgraphProposal(network.SubgraphUrl, network.Id, now)
	if err != nil {
		log.Println("get proposal from subgraph error:", err)
	}
	for _, proposal := range proposals {
		if countMap[fmt.Sprintf("%d-%d", network.Id, proposal.ProposalID)] {
			continue
		}
		voteInfos, err := SubgraphVote(network.SubgraphUrl, proposal.ProposalID)
		if err != nil {
			log.Println("get vote info from subgraph error:", err)
			continue
		}
		var voteList []model.Vote
		var result = make(map[int64]int64, 5)
		fmt.Println("voteInfos: ", voteInfos)
		for _, voteInfo := range voteInfos {
			ipfs, _ := GetIpfs(voteInfo.VoteInfo)
			decrypt, err := Decrypt(ipfs)
			if err != nil {
				log.Println("decrypt error:", err)
			}
			fmt.Println("decrypt string: ", string(decrypt))
			var mapData [][]string
			err = json.Unmarshal(decrypt, &mapData)
			if err != nil {
				log.Println("Unmarshal error：", err)
			}

			for _, value := range mapData {
				optionIdInt := new(big.Int)
				optionIdInt.SetString(value[0], 10)
				valueInt := new(big.Int)
				valueInt.SetString(value[1], 10)
				vote := model.Vote{
					OptionId:        optionIdInt,
					Votes:           valueInt,
					TransactionHash: voteInfo.TransactionHash,
					Address:         voteInfo.Address,
				}
				voteList = append(voteList, vote)
			}
		}
		fmt.Println("voteList: ", voteList)
		var voteHistoryList []model.Vote
		for _, vote := range voteList {
			balance, err := getBalance(vote.Address, ethClient)
			if err != nil {
				log.Println("get balance error:", err)
			}
			var votes float64
			var votePercent = float64(vote.Votes.Int64()) / 100
			if vote.Votes.Int64() != 0 {
				votes = votePercent * balance
			}
			var ethVotes float64
			if network.Id == 280 || network.Id == 324 {
				ethVotes, err = getEthVotes(vote.Address, votePercent, network.Id)
				if err != nil {
					log.Println("get eth votes error: ", err)
				}
			}
			votes += ethVotes
			voteHistory := model.Vote{
				OptionId:        vote.OptionId,
				Votes:           big.NewInt(int64(votes)),
				Address:         vote.Address,
				TransactionHash: vote.TransactionHash,
			}
			voteHistoryList = append(voteHistoryList, voteHistory)
			if _, ok := result[vote.OptionId.Int64()]; ok {
				result[vote.OptionId.Int64()] += int64(votes)
			} else {
				result[vote.OptionId.Int64()] = int64(votes)
			}
		}
		var voteResultList []contract.VoteResult
		for k, v := range result {
			voteResult := contract.VoteResult{
				OptionId: big.NewInt(k),
				Votes:    big.NewInt(v),
			}
			voteResultList = append(voteResultList, voteResult)
		}

		// save to IPFS
		cid, err := nftStorage(voteHistoryList)
		if err != nil {
			log.Println("save vote list to ipfs error:", err)
		}

		//  contract
		fmt.Println(cid, voteResultList, proposal.ProposalID)
		err = contract.CountContract(big.NewInt(proposal.ProposalID), voteResultList, cid, ethClient)
		if err != nil {
			log.Println("call contract count function error:", err)
		}
		countMap[fmt.Sprintf("%d-%d", network.Id, proposal.ProposalID)] = true
	}

}

func getEthVotes(address string, percent float64, networkId int64) (float64, error) {
	// zkSync test net
	var client contract.GoEthClient
	var err error
	if networkId == 280 {
		client, err = contract.GetGoerliClient()
		if err != nil {
			log.Println("get goerli client error: ", err)
			return 0, err
		}

	}
	// zkSync main net
	if networkId == 324 {
		client, err = contract.GetEthMainClient()
		if err != nil {
			log.Println("get eth main net client error: ", err)
			return 0, err
		}
	}
	balance, err := getBalance(address, client)
	if err != nil {
		log.Println("get balance error:", err)
		return 0, err
	}
	return balance * percent, nil
}

func getBalance(address string, ethClient contract.GoEthClient) (float64, error) {
	ctx := context.Background()
	addr := common.HexToAddress(address[0:42])
	number, err := ethClient.Client.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}
	balance, err := ethClient.Client.BalanceAt(ctx, addr, big.NewInt(int64(number)))
	if err != nil {
		return 0, err
	}
	balanceFloat := float64(balance.Int64())
	return balanceFloat / float64(1000000000000000000), nil
}

// SubgraphProposal Subgraph gets Proposals
func SubgraphProposal(url string, chainId, expTime int64) ([]model.ProposalRes, error) {
	//get proposal from subgraph, status == 0 && chainId == network.chainId && expTime_lte == expTime
	query := fmt.Sprintf(`
	query {
		proposals(where: { status: 0, chainId: %d,expTime_lte:%d }) {
		  proposalId
		}
	  }
	`, chainId, expTime)
	requestBody := model.GraphQLRequestBody{
		Query: query,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err

	}
	var graphQLResponse model.GraphQLProposalResponse
	err = json.Unmarshal(body, &graphQLResponse)
	if err != nil {
		return nil, err
	}

	if len(graphQLResponse.Errors) > 0 {
		fmt.Println("GraphQL errors:")
		for _, err := range graphQLResponse.Errors {
			fmt.Println(err.Message)
		}
		return nil, err

	}
	return graphQLResponse.Data.Proposals, nil
}

// SubgraphVote Subgraph gets Voteinfo
func SubgraphVote(url string, proposalsId int64) ([]model.VoteInfo, error) {
	query := fmt.Sprintf(`
		query {
			votes(where: {proposalId: %d}) {
				id
				voteInfo
				transactionHash
			  }
			}
		`, proposalsId)
	requestBody := model.GraphQLRequestBody{
		Query: query,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var graphQLVoteResponse model.GraphQLVoteResponse
	err = json.Unmarshal(body, &graphQLVoteResponse)
	if err != nil {
		return nil, err
	}

	if len(graphQLVoteResponse.Errors) > 0 {
		fmt.Println("GraphQL errors:")
		for _, err := range graphQLVoteResponse.Errors {
			fmt.Println(err.Message)
		}
		return nil, err
	}

	return graphQLVoteResponse.Data.Votes, nil
}

func GetIpfs(votinfo string) (string, error) {
	url := fmt.Sprintf("https://%s.ipfs.nftstorage.link/", votinfo)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err

	}

	return string(body), err
}

func Decrypt(ipfs string) ([]byte, error) {
	// Construct a network that can talk to a drand network. Example using the mainnet fastnet network.
	replace := strings.ReplaceAll(ipfs, "\\n", "\n")
	replace2 := strings.ReplaceAll(replace, "\"", "")

	var network *drandhttp.Network
	var err error
	for _, url := range config.Client.Drand.Url {
		network, err = drandhttp.NewNetwork(url, config.Client.Drand.ChainHash)
		if err == nil {
			break
		}
	}
	reader := strings.NewReader(replace2)
	// Write the encrypted file data to this buffer.
	var cipherData bytes.Buffer
	// Encrypt the data for the given round.
	if err := tlock.New(network).Decrypt(&cipherData, reader); err != nil {
		log.Fatalf("Decrypt: %v", err)
		return nil, err
	}

	data := make([]byte, cipherData.Len())
	_, err = cipherData.Read(data)
	if err != nil {
		log.Fatalf("read: %v", err)
		return nil, err
	}
	return data, nil
}

func nftStorage(voteList []model.Vote) (string, error) {
	marshal, err := json.Marshal(voteList)
	if err != nil {
		log.Println("marshal error:", err)
		return "", err
	}
	url := config.Client.Ipfs.UploadPath
	payload := bytes.NewBuffer(marshal)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		log.Println("Error creating request:", err)
		return "", err
	}

	// 设置请求头，包括你的 NFT.Storage API 密钥
	req.Header.Set("Authorization", config.Client.Ipfs.Token)
	req.Header.Set("Content-Type", "text/plain")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// 检查响应
	if resp.StatusCode != 200 {
		log.Println("Error uploading data. Status:", resp.Status)
		return "", err
	}

	// 从响应中提取 CID（Content Identifier）
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Error decoding response:", err)
		return "", err
	}

	cid := result["value"].(map[string]interface{})["cid"].(string)
	return cid, nil
}

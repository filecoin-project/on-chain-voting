package service

import (
	"encoding/json"
	"fil-vote/config"
	"fil-vote/model"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
)

// Helper function for making GET requests
func makeGETRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		// Log the error using zap
		zap.L().Error("failed to execute GET request", zap.String("url", url), zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// Log the error using zap
		zap.L().Error("failed to read response body", zap.Error(err))
		return nil, err
	}
	return body, nil
}

// GetProposalList sends a GET request to retrieve a list of proposals and returns a list of Proposal models
func GetProposalList(page, pageSize int) (model.ProposalsResponse, error) {
	url := fmt.Sprintf(config.Client.Network.PowerBackendURL+model.BaseProposalAPIPath+"list?status=0&searchKey=&chainId=%d&page=%d&pageSize=%d", config.Client.Network.ChainID, page, pageSize)
	body, err := makeGETRequest(url)
	if err != nil {
		return model.ProposalsResponse{}, err
	}

	var result model.ProposalsResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		zap.L().Error("failed to parse JSON response", zap.Error(err))
		return model.ProposalsResponse{}, err
	}
	return result, nil
}

// GetProposalByID retrieves a proposal by its ID
func GetProposalByID(proposalId int64) (model.Proposal, error) {
	url := fmt.Sprintf(config.Client.Network.PowerBackendURL+model.BaseProposalAPIPath+"details?chainId=%d&proposalId=%d", config.Client.Network.ChainID, proposalId)
	body, err := makeGETRequest(url)
	if err != nil {
		return model.Proposal{}, err
	}

	var result struct {
		Proposals model.Proposal `json:"data"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		zap.L().Error("failed to parse JSON response", zap.Error(err))
		return model.Proposal{}, err
	}

	return result.Proposals, nil
}

// GetVotesByProposalID retrieves votes for a specific proposal by its ID
func GetVotesByProposalID(proposalId int64) ([]model.Vote, error) {
	url := fmt.Sprintf(config.Client.Network.PowerBackendURL+model.BaseProposalAPIPath+"votes?chainId=%d&proposalId=%d", config.Client.Network.ChainID, proposalId)
	body, err := makeGETRequest(url)
	if err != nil {
		return nil, err
	}

	var result model.VoteResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		zap.L().Error("failed to parse JSON response", zap.Error(err))
		return nil, err
	}

	return result.Data, nil
}

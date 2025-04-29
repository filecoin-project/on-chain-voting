package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fil-vote/config"
	"fil-vote/model"
	"fil-vote/utils"
	"fmt"
	"go.uber.org/zap"
	"math/big"
	"os"

	"github.com/filecoin-project/lotus/chain/types"
	cbor "github.com/whyrusleeping/cbor/go"
)

// Vote allows a user to vote on a proposal (approve or reject).
func Vote(client *RPCClient, from, action string, proposalId int64) (string, error) {
	// Fetch proposal details using the provided proposal ID
	proposalDetail, err := GetProposalByID(proposalId)
	if err != nil {
		zap.L().Error("Failed to fetch proposal details", zap.Int64("proposalId", proposalId), zap.Error(err))
		return "", err
	}

	if proposalDetail.ProposalId == 0 {
		return "", errors.New("invalid Proposal ID")
	}

	if proposalDetail.Status != 2 {
		return "", errors.New("the proposal is not in an 'in-progress' state, voting cannot proceed")
	}

	// Encrypt the vote result (approve or reject) using proposal's end time
	result, err := utils.EncryptVoteResult([][]string{{action}}, proposalDetail.EndTime)
	if err != nil {
		zap.L().Error("Failed to encrypt vote result", zap.Error(err))
		return "", err
	}

	// Encode the encrypted result and proposal ID into ABI format
	encodedData, err := config.PowerVotingAbi.Pack("vote", big.NewInt(proposalId), result)
	if err != nil {
		zap.L().Error("Failed to ABI encode vote data", zap.Error(err))
		return "", err
	}

	// Serialize the ABI-encoded data to CBOR format
	var buf bytes.Buffer
	enc := cbor.NewEncoder(&buf)
	if err := enc.Encode(encodedData); err != nil {
		zap.L().Error("Failed to serialize data to CBOR", zap.Error(err))
		return "", err
	}

	// Convert the CBOR data to a base64-encoded string
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Retrieve the nonce for the sender address to avoid replay attacks
	nonce, err := client.GetNonce(context.Background(), from)
	if err != nil {
		zap.L().Error("Failed to get nonce for sender", zap.Error(err))
		return "", err
	}

	// Prepare the message with the encoded vote data for sending
	estimateMessage := model.Message{
		Version:    0,
		To:         config.Client.Network.PowerVotingContract,
		From:       from,
		Nonce:      nonce,
		Value:      "0",                  // No value is transferred in the vote message
		GasLimit:   0,                    // Gas limit for the transaction (to be estimated)
		GasFeeCap:  "0",                  // Maximum fee cap
		GasPremium: "0",                  // Gas premium
		Method:     model.InvokeContract, // Method ID for the "vote" function
		Params:     base64Str,            // Base64-encoded vote data
		Cid:        types.EmptyTSK,       // Empty CID (will be generated during transaction)
	}

	// Estimate gas usage for the message
	estimatedMsg, err := client.EstimateGas(context.Background(), estimateMessage)
	if err != nil {
		zap.L().Error("Failed to estimate gas for the message", zap.Error(err))
		return "", err
	}

	// Display the estimated gas cost to the user
	gasFeeCapInNanoFIL := new(big.Rat).SetFrac(estimatedMsg.GasFeeCap.Int, big.NewInt(1000000000))
	// Output the result in nanoFIL with a more user-friendly message
	gasFeeCapInFILFloat, _ := gasFeeCapInNanoFIL.Float64()

	decimalPlaces := len(estimatedMsg.GasFeeCap.String())
	fmt.Printf("The total fee rate set by sender: %.*f nanoFIL\n", decimalPlaces, gasFeeCapInFILFloat)

	// Prompt the user for confirmation to proceed with the vote
	for {
		fmt.Print("Do you want to proceed with the vote? Type 'yes' to confirm or 'no' to cancel: ")

		// Read user input
		reader := bufio.NewReader(os.Stdin)
		confirmation, _ := reader.ReadString('\n')
		confirmation = confirmation[:len(confirmation)-1] // Remove the newline character

		if confirmation == "yes" {
			// Send the vote transaction after estimating gas
			messageHash, err := client.SendMessage(context.Background(), estimatedMsg)
			if err != nil {
				zap.L().Error("Failed to send vote message", zap.Error(err))
				return "", err
			}
			// Return the message hash after successful submission
			return messageHash, nil
		} else if confirmation == "no" {
			// If the user doesn't confirm, exit without sending the transaction
			return "", errors.New("vote not confirmed, transaction not sent")
		} else {
			// Handle invalid input
			fmt.Println("Invalid input, please type 'yes' or 'no'.")
		}
	}
}

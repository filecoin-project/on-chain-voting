// Copyright (C) 2023-2024 StorSwift Inc.
// This file is part of the PowerVoting library.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"backend/models"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Res struct {
	JsonRpc string `json:"jsonRpc"`
	Result  string `json:"result"`
	Id      int    `json:"id"`
}

// GetContractMapping retrieves mapping result for a given method from a smart contract.
func GetContractMapping(methodName string, ethClient models.GoEthClient, args []interface{}) (map[string]interface{}, error) {
	methodAbi := ethClient.Abi.Methods[methodName]

	arguments, err := methodAbi.Inputs.Pack(args...)
	if err != nil {
		return nil, err
	}

	data := append(methodAbi.ID, arguments...)

	ethCall, err := EthCall(hex.EncodeToString(data), ethClient)
	if err != nil {
		return nil, err
	}

	hexString := ethCall.Result[2:]
	byteData, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}

	resultMap := make(map[string]interface{}, len(methodAbi.Outputs))
	err = methodAbi.Outputs.UnpackIntoMap(resultMap, byteData)
	if err != nil {
		return nil, err
	}

	return resultMap, err
}

// EthCall executes an Ethereum RPC call to retrieve information from the blockchain.
func EthCall(data string, ethClient models.GoEthClient) (Res, error) {

	payload := strings.NewReader(fmt.Sprintf(`{
	  "method": "eth_call",
	  "params": [
	      {
	          "to": "%s",
	          "data": "%s"
	      },
	      "latest"
	  ],
	  "id": 1,
	  "jsonrpc": "2.0"
	}`, ethClient.ContractAddress.String(), data))

	client := &http.Client{}

	req, err := http.NewRequest("POST", ethClient.Rpc, payload)
	if err != nil {
		return Res{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		return Res{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Res{}, err
	}

	var resque Res
	err = json.Unmarshal(body, &resque)
	if err != nil {
		return Res{}, err
	}

	return resque, nil
}

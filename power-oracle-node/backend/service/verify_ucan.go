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
	"backend/utils"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/storyicon/sigverify"
	"github.com/ybbus/jsonrpc/v3"
	"go.uber.org/zap"
	"strings"
)

// VerifyUCAN verifies the authenticity of a User Controlled Authorization Network (UCAN) token.
// It splits the UCAN into its components, verifies the signature, and performs additional validation
// depending on whether it is a single or double UCAN.
func VerifyUCAN(ucan string, lotusRpcClient jsonrpc.RPCClient) (string, string, string, bool, error) {
	_, payload, signatureBytes, _, err := ucanSplit(ucan)
	if err != nil {
		zap.L().Error("failed to split ucan into parts", zap.Error(err))
		return "", "", "", false, err
	}

	if payload.Prf != "" {
		return VerifyDoubleUcan(payload, signatureBytes, ucan, lotusRpcClient)
	}

	if payload.Prf == "" {
		parts := strings.Split(ucan, ".")
		if len(parts) != 3 {
			return "", "", "", false, errors.New("invalid ucan format")
		}

		headerStr := parts[0]
		payloadStr := parts[1]
		verificationData := []byte(fmt.Sprintf("%s.%s", headerStr, payloadStr))

		if _, err := VerifyEthUcan(payload.Iss, string(signatureBytes), verificationData); err != nil {
			zap.L().Error("failed to verify eth ucan", zap.Error(err))
			return "", "", "", false, err
		}

		return payload.Iss, payload.Aud, payload.Act, true, nil
	}

	return "", "", "", false, errors.New("ucan has no prf")
}

// VerifyDoubleUcan verifies the authenticity of a double UCAN token.
// It verifies both the Ethereum and Filecoin signatures of the double UCAN.
func VerifyDoubleUcan(payload models.Payload, signatureBytes []byte, ucan string, lotusRpcClient jsonrpc.RPCClient) (string, string, string, bool, error) {
	parts := strings.Split(ucan, ".")
	if len(parts) != 3 {
		return "", "", "", false, errors.New("invalid ucan format")
	}

	headerStr, payloadStr := parts[0], parts[1]

	ethVerificationData := []byte(fmt.Sprintf("%s.%s", headerStr, payloadStr))
	ethVerify, err := VerifyEthUcan(payload.Iss, string(signatureBytes), ethVerificationData)
	if err != nil {
		zap.L().Error("Failed to verify eth ucan signature", zap.Error(err))
		return "", "", "", false, err
	}

	zap.L().Info("eth verification result", zap.Bool("eth verify", ethVerify))

	if ethVerify {
		parts := strings.Split(payload.Prf, ".")
		if len(parts) != 3 {
			return "", "", "", false, errors.New("invalid ucan format")
		}

		headerStr, payloadStr := parts[0], parts[1]
		fileCoinHeader, fileCoinPayload, fileCoinSignatureBytes, _, err := ucanSplit(payload.Prf)
		if err != nil {
			zap.L().Error("failed to split filecoin ucan", zap.Error(err))
			return "", "", "", false, err
		}

		filecoinSignature, err := createFilecoinSignature(fileCoinHeader.Alg, fileCoinSignatureBytes)
		if err != nil {
			zap.L().Error("failed to create filecoin signature", zap.Error(err))
			return "", "", "", false, err
		}

		fileCoinVerify, err := VerifyFilecoinUcan(fileCoinPayload.Iss, filecoinSignature, []byte(fmt.Sprintf("%s.%s", headerStr, payloadStr)), lotusRpcClient)
		if err != nil {
			zap.L().Error("failed to verify filecoin ucan signature", zap.Error(err))
			return "", "", "", false, err
		}

		zap.L().Info("filecoin verification result", zap.Bool("filecoin verify", fileCoinVerify))

		if fileCoinVerify {
			if payload.Iss == fileCoinPayload.Aud && payload.Aud == fileCoinPayload.Iss && payload.Act == fileCoinPayload.Act {
				return payload.Iss, payload.Aud, payload.Act, false, nil
			} else {
				return "", "", "", false, errors.New("the signature contents of eth and filecoin do not match")
			}
		}
	}

	return "", "", "", false, errors.New("verification failed")
}

// createFilecoinSignature creates a Filecoin signature based on the provided algorithm and signature bytes.
func createFilecoinSignature(alg string, signatureBytes []byte) (crypto.Signature, error) {
	switch alg {

	case models.KTSecp256k1:
		return crypto.Signature{
			Type: crypto.SigType(models.SigTypeSecp256k1),
			Data: signatureBytes,
		}, nil

	case models.KTBLS:
		return crypto.Signature{
			Type: crypto.SigType(models.SigTypeBLS),
			Data: signatureBytes,
		}, nil
	default:
		return crypto.Signature{}, errors.New("unsupported filecoin signature algorithm")
	}
}

// VerifyFilecoinUcan verifies the authenticity of a Filecoin UCAN signature.
func VerifyFilecoinUcan(address string, signature crypto.Signature, msgData []byte, lotusRpcClient jsonrpc.RPCClient) (bool, error) {
	verify, err := utils.WalletVerify(context.Background(), lotusRpcClient, address, signature, msgData)
	if err != nil {
		return verify, err
	}
	return verify, nil
}

// VerifyEthUcan verifies the authenticity of an Ethereum UCAN signature.
func VerifyEthUcan(address string, signature string, msgData []byte) (bool, error) {
	valid, err := sigverify.VerifyEllipticCurveHexSignatureEx(
		common.HexToAddress(address),
		msgData,
		signature,
	)
	if err != nil {
		return valid, err
	}
	return valid, nil
}

// ucanSplit splits a UCAN token into its header, payload, and signature components.
func ucanSplit(ucan string) (models.Header, models.Payload, []byte, []byte, error) {
	parts := strings.Split(ucan, ".")
	if len(parts) != 3 {
		return models.Header{}, models.Payload{}, nil, nil, errors.New("invalid ucan format")
	}

	headerStr := parts[0]
	payloadStr := parts[1]
	sig := parts[2]

	decodeBase64 := func(encoded string) ([]byte, error) {
		decoded, err := base64.RawURLEncoding.DecodeString(encoded)
		if err != nil {
			zap.L().Error("failed to decode base64 string", zap.Error(err))
			return nil, err
		}
		return decoded, nil
	}

	headerBytes, err := decodeBase64(headerStr)
	if err != nil {
		zap.L().Error("ailed to decode header", zap.Error(err))
		return models.Header{}, models.Payload{}, nil, nil, err
	}

	payloadBytes, err := decodeBase64(payloadStr)
	if err != nil {
		zap.L().Error("failed to decode payload", zap.Error(err))
		return models.Header{}, models.Payload{}, nil, nil, err
	}

	signatureBytes, err := decodeBase64(sig)
	if err != nil {
		zap.L().Error("failed to decode signature", zap.Error(err))
		return models.Header{}, models.Payload{}, nil, nil, err
	}

	var header models.Header
	err = json.Unmarshal(headerBytes, &header)
	if err != nil {
		zap.L().Error("failed to parse header", zap.Error(err))
		return models.Header{}, models.Payload{}, nil, nil, err
	}

	var payload models.Payload
	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		zap.L().Error("failed to parse payload", zap.Error(err))
		return models.Header{}, models.Payload{}, nil, nil, err
	}

	return header, payload, signatureBytes, payloadBytes, nil
}

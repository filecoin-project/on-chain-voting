package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fil-vote/config"
	"fmt"
	"github.com/drand/tlock"
	drandhttp "github.com/drand/tlock/networks/http"
	"go.uber.org/zap"
	"strings"
	"time"
)

func EncryptVoteResult(voteInfo [][]string, endTime int64) (string, error) {
	jsonData, err := json.Marshal(voteInfo)
	if err != nil {
		zap.L().Error("Error marshalling data to JSON", zap.Error(err))
		return "", err
	}
	reader := bytes.NewReader(jsonData)

	// Retry encrypting for a few times
	for i := 0; i < 5; i++ {
		encryptedData, err := encrypt(reader, endTime)
		if err == nil {
			return encryptedData, nil
		}

		// Log retry attempt
		zap.L().Warn(fmt.Sprintf("Encrypt failed: %v, retrying %d/%d", err, i+1, 5))
		if i == 4 { // Last retry
			zap.L().Error("Final retry failed", zap.Error(err))
		}
	}

	return "", fmt.Errorf("failed to encrypt data after 5 attempts")
}

func encrypt(dataToEncrypt *bytes.Reader, endTime int64) (string, error) {
	var network *drandhttp.Network
	var err error

	// Try to create a network using the provided URLs
	for _, url := range config.Client.Drand.Urls {
		network, err = drandhttp.NewNetwork(url, config.Client.Drand.ChainHash)
		if err == nil {
			break
		}
	}
	if err != nil {
		return "", fmt.Errorf("failed to create network: %v", err)
	}

	// Use the round number for encryption
	roundNumber := network.RoundNumber(time.Unix(endTime, 0))

	// Encrypt the data
	var encryptedData bytes.Buffer
	err = tlock.New(network).Encrypt(&encryptedData, dataToEncrypt, roundNumber)
	if err != nil {
		return "", fmt.Errorf("error encrypting data: %v", err)
	}

	// Convert the encrypted data to AGE format
	return convertToAgeFormat(encryptedData.Bytes()), nil
}

func convertToAgeFormat(encryptedData []byte) string {
	// Base64 encode the encrypted data
	encodedData := base64.StdEncoding.EncodeToString(encryptedData)
	// Split the encoded data into chunks of 64 characters (AGE format requirement)
	var chunkedData strings.Builder
	for i := 0; i < len(encodedData); i += 64 {
		end := i + 64
		if end > len(encodedData) {
			end = len(encodedData)
		}
		chunkedData.WriteString(encodedData[i:end] + "\n")
	}

	// Format it in AGE style
	ageFormatted := fmt.Sprintf("-----BEGIN AGE ENCRYPTED FILE-----\n%s-----END AGE ENCRYPTED FILE-----", chunkedData.String())
	return ageFormatted
}

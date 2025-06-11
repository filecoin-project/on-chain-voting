package developer

import (
	"context"
	"encoding/hex"
	"fil-vote/model"
	"fil-vote/service"
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"time"
)

// GenerateCmd creates a new command that generates GitHub proof for the wallet address
func GenerateCmd(client *service.RPCClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-github-proof",
		Short: "Generates a GitHub proof for the given wallet address",
		Run: func(cmd *cobra.Command, args []string) {
			// Retrieve the wallet address
			from, _ := cmd.Flags().GetString("from")

			// If the "from" flag is empty, retrieve the default wallet address
			if from == "" {
				var err error
				from, err = client.WalletDefaultAddress(context.Background())
				if err != nil {
					zap.L().Error("Failed to retrieve default wallet address", zap.Error(err))
					return
				}
			}

			// Retrieve the GitHub username
			githubName, err := cmd.Flags().GetString("githubName")
			if err != nil || githubName == "" {
				zap.L().Error("GitHub username is required", zap.Error(err))
				return
			}

			// Check if GitHub username exists
			exists, err := CheckGitHubUserExistence(githubName)
			if err != nil {
				zap.L().Error("Failed to check GitHub username existence", zap.Error(err))
				return
			}
			if !exists {
				zap.L().Error("GitHub username does not exist", zap.String("githubName", githubName))
				return
			}

			// Prepare the signature data
			data := model.SignatureData{
				WalletAddress: from,
				GithubName:    githubName,
				Timestamp:     time.Now().Unix(),
			}

			// Sign the data
			signature, err := client.WalletSign(context.Background(), data)
			if err != nil {
				zap.L().Error("Failed to sign wallet data", zap.Error(err))
				return
			}
			hexString := hex.EncodeToString(signature.Data)

			sigType := int(signature.Type)
			str := strconv.Itoa(sigType)

			hexString = "0" + str + hexString

			proof := fmt.Sprintf(
				`
  I hereby claim:

    * I am %s on Github.
    * I control %s (Filecoin wallet address).

  To claim this, I am signing this object

  {"walletAddress":"%s","githubName":"%s","timestamp":%d}

  with my Filecoin wallet's private key, yielding the signature:

  %s

  And finally, I am proving ownership of the github account by posting this as a gist.`,
				githubName, data.WalletAddress, data.WalletAddress, data.GithubName, data.Timestamp, hexString)

			fmt.Println(proof)
			fmt.Println("\nPlease paste the above text first and obtain gistId: https://gist.github.com/")

		},
	}

	cmd.Flags().String("from", "", "Wallet address to generate the proof for (optional)")
	cmd.Flags().String("githubName", "", "GitHub username to generate the proof for (required)")
	cmd.MarkFlagRequired("githubName")
	return cmd
}

// CheckGitHubUserExistence checks whether a GitHub username exists
func CheckGitHubUserExistence(username string) (bool, error) {
	// Send a GET request to GitHub API to check the user
	resp, err := http.Get(model.GithubAPI + username)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// If the status code is 404, the user does not exist
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	// If status code is 200, the user exists
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	body, _ := io.ReadAll(resp.Body)
	return false, fmt.Errorf("unexpected response from GitHub API: %s", string(body))
}

package wallet

import (
	"context"
	"fil-vote/service"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// AddCmd 返回一个 Cobra 命令，用于导入钱包并显示状态。
func AddCmd(client *service.RPCClient) *cobra.Command {
	return &cobra.Command{
		Use:   "add [walletType] [privateKey]", // 命令名称，包含钱包类型和私钥作为参数
		Short: "导入钱包",                          // 简短的命令描述
		Args:  cobra.ExactArgs(2),              // 确保提供两个参数
		Run: func(cmd *cobra.Command, args []string) {
			// 获取钱包类型和私钥
			walletType := args[0]
			privateKey := args[1]

			// 调用 WalletImport 函数导入钱包
			_, err := client.WalletImport(context.Background(), walletType, privateKey)
			if err != nil {
				// Enhanced logging with more context, including walletType
				zap.L().Error("Failed to import wallet",
					zap.String("WalletType", walletType),
					zap.String("PrivateKey", privateKey),
					zap.Error(err))

			}
			//else {
			//	displayWalletStatus(walletAddress, "Success", "")
			//}
		},
	}
}

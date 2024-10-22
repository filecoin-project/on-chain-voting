package task

import (
	"powervoting-server/client"
	"powervoting-server/config"
	"powervoting-server/db"
	"testing"

	"go.uber.org/zap"
)

func TestBackupPowerHandler(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	config.InitConfig("../")
	client.InitW3Client()
	db.InitMysql()

	BackupPowerHandler()
}

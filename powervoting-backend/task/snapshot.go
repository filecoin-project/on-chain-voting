package task

import (
	"powervoting-server/client"
	"powervoting-server/config"
	"powervoting-server/contract"
	"powervoting-server/db"
	"powervoting-server/model"
	"sync"
	"time"

	"go.uber.org/zap"
)

func BackupPowerHandler() {
	zap.L().Info("backup power start: ", zap.Int64("timestamp", time.Now().Unix()))
	wg := sync.WaitGroup{}
	errList := make([]error, 0, len(config.Client.Network))
	mu := &sync.Mutex{}

	for _, network := range config.Client.Network {
		network := network
		ethClient, err := contract.GetClient(network.Id)
		if err != nil {
			zap.L().Error("get go-eth client error:", zap.Error(err))
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := BackupPower(ethClient, db.Engine); err != nil {
				mu.Lock()
				errList = append(errList, err)
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	if len(errList) != 0 {
		zap.L().Error("backup power finished with err:", zap.Errors("errors", errList))
	}
	zap.L().Info("backup power finished: ", zap.Int64("timestamp", time.Now().Unix()))
}

func BackupPower(ethClient model.GoEthClient, db db.DataRepo) error {
	list, err := db.GetSnapshotList(ethClient.Id)
	if err != nil {
		zap.L().Error("fail to get snapshot list ", zap.Error(err))
		return err
	}

	zap.L().Info("snapshot list: ", zap.Any("list", list))
	for _, item := range list {
		zap.L().Info("backup snapshot start", zap.Any("info", item))

		allPower, err := client.GetAllAddressPowerByDay(ethClient.Id, item.Day)
		if err != nil {
			zap.L().Error("failed to get all power: ", zap.Error(err), zap.Any("day", item.Day))
			continue
		}

		zap.L().Info("all power: ", zap.Any("allPower", allPower))
		// item.PowerInfo = allPower.PowerInfo
		// fileName := filepath.Join("./backup", fmt.Sprintf("%d.txt", item.Id))
		// jsonData, err := json.Marshal(allPower.PowerInfo)
		// if err != nil {
		// 	zap.L().Error("failed to marshal content to JSON: ", zap.Error(err))
		// 	continue
		// }

		// err = os.WriteFile(fileName, jsonData, 0644)
		// if err != nil {
		// 	zap.L().Error("failed to write JSON to file: ", zap.Error(err))
		// 	continue
		// }

		// zap.L().Info("Successfully saved content to: ", zap.String("file", fileName))

		// absPath, err := filepath.Abs(fileName)
		// if err != nil {
		// 	zap.L().Error("failed to get absolute path: ", zap.Error(err))
		// 	continue
		// }

		// zap.L().Info("upload with w3,absPath:", zap.String("absPath", absPath))

		// cid, err := client.W3.Upload(absPath)
		// if err != nil {
		// 	os.Remove(absPath)
		// 	zap.L().Error("failed to upload file: ", zap.Error(err))
		// 	continue
		// }

		// zap.L().Info("upload with w3,cid:", zap.String("cid", cid))
		item.Cid = "test123"

		zap.L().Info("update snapshot ", zap.Any("id", item.Id))
		if err = db.UpdateSnapshot(item); err != nil {
			// err := os.Remove(absPath)
			if err != nil {
				zap.L().Error("failed to remove file: ", zap.Error(err))
				continue
			}
			zap.L().Error("update snapshot error: ", zap.Error(err))
			continue
		}

		// os.Remove(absPath)
	}

	return nil
}

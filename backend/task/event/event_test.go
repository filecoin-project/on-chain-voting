package event

import (
	"powervoting-server/config"
	"powervoting-server/data"
	"powervoting-server/repo"
	"powervoting-server/service"
)

func getSyncService() *service.SyncService {
	// config.GetDefaultConfig()
	config.InitConfig("../../")
	config.Client.ABIPath.PowerVotingAbi = "../../abi/power-voting.json"
	config.Client.ABIPath.FipAbi = "../../abi/power-voting-fip.json"
	config.Client.ABIPath.OracleAbi = "../../abi/oracle.json"

	return service.NewSyncService(
		repo.NewSyncRepo(data.NewMysql()),
		repo.NewVoteRepo(data.NewMysql()),
		repo.NewProposalRepo(data.NewMysql()),
		repo.NewFipRepo(data.NewMysql()),
		repo.NewLotusRPCRepo(),
	)
}

// func TestFetchMatchingEventLogs(t *testing.T) {
// 	syncService := getSyncService()
// 	gethClient, err := data.GetClient(syncService, 314159)
// 	assert.NoError(t, err)
// 	ev := &Event{
// 		Client:      gethClient,
// 		SyncService: syncService,
// 		Network:     &config.Client.Network,
// 	}

// 	logs, err := ev.FetchMatchingEventLogs(context.Background(), big.NewInt(2539215), big.NewInt(2539220))
// 	assert.NoError(t, err)

// 	ev.ProcessingEventLogs(context.Background(), logs)
// }

// func TestFetchEventFromRPC(t *testing.T) {
// 	syncService := getSyncService()
// 	gethClient, err := data.GetClient(syncService, 314159)
// 	assert.NoError(t, err)
// 	ev := &Event{
// 		Client:      gethClient,
// 		SyncService: syncService,
// 		Network:     &config.Client.Network,
// 	}

// 	var (
// 		url = "https://calibration.filfox.info/api/v1/address/0xBd87133b920e3147705362512f3897465f35f50E/events"
// 	)

// 	logs, err := FetchEventFromRPC(url)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, logs)

// 	ev.ProcessingEventLogs(context.Background(), logs)
// }

// func TestUpdateMinerIds(t *testing.T) {
// 	syncService := getSyncService()
// 	gethClient, err := data.GetClient(syncService, 314159)
// 	assert.NoError(t, err)
// 	ev := &Event{
// 		Client:      gethClient,
// 		SyncService: syncService,
// 		Network:     &config.Client.Network,
// 	}

// 	err = ev.SyncService.UpdateVoterByMinerIds(
// 		context.Background(),
// 		"0xfF000000000000000000000000000000000278bc",
// 		[]uint64{161842},
// 	)

// 	assert.NoError(t, err)
// }

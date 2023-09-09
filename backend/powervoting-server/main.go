package main

import (
	"powervoting-server/config"
	"powervoting-server/scheduler"
)

func main() {
	// initialization configuration
	config.InitConfig()
	scheduler.TaskScheduler()
	//scheduler.VotingCountFunc()
	//select {}
}

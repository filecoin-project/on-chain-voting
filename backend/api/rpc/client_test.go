package rpc

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"powervoting-server/config"
)

func TestGetAddressPowerByDay(t *testing.T) {
	config.InitConfig("../../")
	config.InitLogger()

	client := getClient()
	assert.NotNil(t, client)

	powers, err := GetAddressPowerByDay(314159, "0xBc27ca842D22cD5BdBC41B27A571EC1FbB559307", "20250211")
	assert.Nil(t, err)
	assert.NotNil(t, powers)
}

func TestUploadSnapshotInfo(t *testing.T) {
	config.InitConfig("../../")
	config.InitLogger()

	client := getClient()
	assert.NotNil(t, client)

	res, err := UploadSnapshotInfo(314159, "20250323")
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

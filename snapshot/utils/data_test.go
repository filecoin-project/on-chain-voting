package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"power-snapshot/config"
	"power-snapshot/utils"
)

func TestCalMissDates(t *testing.T) {
	config.InitConfig("../")
	res := utils.CalMissDates([]string{})
	assert.NotEmpty(t, res)
}
